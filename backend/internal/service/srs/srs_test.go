package srs

import (
	"context"
	"testing"
	"time"

	"github.com/clearclown/HaiLanGo/backend/internal/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestGetReviewItemsWithPriority は優先度付き復習項目の取得をテスト
func TestGetReviewItemsWithPriority(t *testing.T) {
	ctx := context.Background()
	service := NewMockSRSService()

	userID := uuid.New()
	bookID := uuid.New()
	now := time.Now()

	// 異なる優先度の項目を作成
	items := []*models.ReviewItem{
		{
			ID:             uuid.New(),
			UserID:         userID,
			BookID:         bookID,
			PageNumber:     1,
			Content:        "Urgent 1",
			NextReviewDate: &[]time.Time{now.AddDate(0, 0, -2)}[0],
		},
		{
			ID:             uuid.New(),
			UserID:         userID,
			BookID:         bookID,
			PageNumber:     2,
			Content:        "Urgent 2",
			NextReviewDate: &[]time.Time{now.AddDate(0, 0, -1)}[0],
		},
		{
			ID:             uuid.New(),
			UserID:         userID,
			BookID:         bookID,
			PageNumber:     3,
			Content:        "Recommended",
			NextReviewDate: &[]time.Time{now.AddDate(0, 0, 1)}[0],
		},
		{
			ID:             uuid.New(),
			UserID:         userID,
			BookID:         bookID,
			PageNumber:     4,
			Content:        "Relaxed",
			NextReviewDate: &[]time.Time{now.AddDate(0, 0, 3)}[0],
		},
	}

	for _, item := range items {
		service.items[item.ID] = item
	}

	result, err := service.GetReviewItems(ctx, userID, now)
	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result.UrgentItems, 2)
	assert.Len(t, result.RecommendedItems, 1)
	assert.Len(t, result.RelaxedItems, 1)
}

// TestCompleteReview は復習完了処理をテスト
func TestCompleteReview(t *testing.T) {
	ctx := context.Background()

	userID := uuid.New()
	bookID := uuid.New()
	now := time.Now()

	tests := []struct {
		name               string
		score              int
		expectedReviewCount int
		checkNextReviewDate bool
	}{
		{"高得点（90点）", 90, 1, true},
		{"中得点（75点）", 75, 1, true},
		{"低得点（60点）", 60, 1, true},
		{"最低得点（40点）", 40, 1, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 各テストケースで新しいserviceとitemを作成
			service := NewMockSRSService()

			item := &models.ReviewItem{
				ID:             uuid.New(),
				UserID:         userID,
				BookID:         bookID,
				PageNumber:     1,
				Content:        "Test",
				ReviewCount:    0,
				LastReviewDate: &now,
			}

			service.items[item.ID] = item

			err := service.CompleteReview(ctx, item.ID, tt.score, 30)
			require.NoError(t, err)

			updated := service.items[item.ID]
			assert.Equal(t, tt.expectedReviewCount, updated.ReviewCount)
			assert.Equal(t, tt.score, updated.LastScore)

			if tt.checkNextReviewDate {
				assert.NotNil(t, updated.NextReviewDate)
				assert.True(t, updated.NextReviewDate.After(now))
			}
		})
	}
}

// TestGetReviewStats は統計情報の取得をテスト
func TestGetReviewStats(t *testing.T) {
	ctx := context.Background()
	service := NewMockSRSService()

	userID := uuid.New()
	bookID := uuid.New()
	now := time.Now()

	// テストデータを作成
	for i := 0; i < 10; i++ {
		var nextReviewDate time.Time
		if i < 3 {
			// 緊急項目
			nextReviewDate = now.AddDate(0, 0, -1)
		} else if i < 7 {
			// 推奨項目
			nextReviewDate = now.AddDate(0, 0, 1)
		} else {
			// 余裕あり項目
			nextReviewDate = now.AddDate(0, 0, 3)
		}

		item := &models.ReviewItem{
			ID:             uuid.New(),
			UserID:         userID,
			BookID:         bookID,
			PageNumber:     i + 1,
			Content:        "Test",
			NextReviewDate: &nextReviewDate,
		}
		service.items[item.ID] = item
	}

	stats, err := service.GetStats(ctx, userID, now)
	require.NoError(t, err)
	assert.Equal(t, 10, stats.TotalReviewItems)
	assert.Equal(t, 3, stats.UrgentItems)
	assert.Equal(t, 4, stats.RecommendedItems)
	assert.Equal(t, 3, stats.RelaxedItems)
}

// TestCreateReviewItemFromPhrase はフレーズから復習項目を作成をテスト
func TestCreateReviewItemFromPhrase(t *testing.T) {
	ctx := context.Background()
	service := NewMockSRSService()

	userID := uuid.New()
	bookID := uuid.New()

	phraseData := &PhraseData{
		UserID:      userID,
		BookID:      bookID,
		PageNumber:  1,
		Content:     "Здравствуйте!",
		Translation: "こんにちは",
	}

	itemID, err := service.CreateReviewItem(ctx, phraseData)
	require.NoError(t, err)
	assert.NotEqual(t, uuid.Nil, itemID)

	// 作成された項目を確認
	item := service.items[itemID]
	assert.NotNil(t, item)
	assert.Equal(t, phraseData.Content, item.Content)
	assert.Equal(t, 0, item.ReviewCount)
	assert.Nil(t, item.NextReviewDate) // 初回はまだ復習日未設定
}

// TestBulkCreateReviewItems は一括作成をテスト
func TestBulkCreateReviewItems(t *testing.T) {
	ctx := context.Background()
	service := NewMockSRSService()

	userID := uuid.New()
	bookID := uuid.New()

	phrases := []*PhraseData{
		{UserID: userID, BookID: bookID, PageNumber: 1, Content: "Hello", Translation: "こんにちは"},
		{UserID: userID, BookID: bookID, PageNumber: 1, Content: "Goodbye", Translation: "さようなら"},
		{UserID: userID, BookID: bookID, PageNumber: 2, Content: "Thank you", Translation: "ありがとう"},
	}

	itemIDs, err := service.BulkCreateReviewItems(ctx, phrases)
	require.NoError(t, err)
	assert.Len(t, itemIDs, 3)

	// すべての項目が作成されたことを確認
	for _, id := range itemIDs {
		assert.NotNil(t, service.items[id])
	}
}

// MockSRSService はテスト用のモックサービス
type MockSRSService struct {
	items     map[uuid.UUID]*models.ReviewItem
	histories []models.ReviewHistory
}

func NewMockSRSService() *MockSRSService {
	return &MockSRSService{
		items:     make(map[uuid.UUID]*models.ReviewItem),
		histories: make([]models.ReviewHistory, 0),
	}
}

func (m *MockSRSService) GetReviewItems(ctx context.Context, userID uuid.UUID, now time.Time) (*ReviewItemsByPriority, error) {
	result := &ReviewItemsByPriority{
		UrgentItems:      make([]*models.ReviewItem, 0),
		RecommendedItems: make([]*models.ReviewItem, 0),
		RelaxedItems:     make([]*models.ReviewItem, 0),
	}

	for _, item := range m.items {
		if item.UserID == userID {
			priority := item.GetPriority(now)
			switch priority {
			case models.PriorityUrgent:
				result.UrgentItems = append(result.UrgentItems, item)
			case models.PriorityRecommended:
				result.RecommendedItems = append(result.RecommendedItems, item)
			case models.PriorityRelaxed:
				result.RelaxedItems = append(result.RelaxedItems, item)
			}
		}
	}

	return result, nil
}

func (m *MockSRSService) CompleteReview(ctx context.Context, itemID uuid.UUID, score int, timeSpentSec int) error {
	item, ok := m.items[itemID]
	if !ok {
		return nil
	}

	now := time.Now()

	// 次回復習日を計算（srs.CalculateNextReviewDate を使用）
	// ここではモックなので簡易的な実装
	item.ReviewCount++
	item.LastScore = score
	item.LastReviewDate = &now

	nextDate := now.AddDate(0, 0, 3) // 簡易的に3日後
	item.NextReviewDate = &nextDate

	// 履歴を記録
	history := models.ReviewHistory{
		ID:           uuid.New(),
		ReviewItemID: itemID,
		UserID:       item.UserID,
		Score:        score,
		ReviewedAt:   now,
		TimeSpentSec: timeSpentSec,
	}
	m.histories = append(m.histories, history)

	return nil
}

func (m *MockSRSService) GetStats(ctx context.Context, userID uuid.UUID, now time.Time) (*models.ReviewStats, error) {
	stats := &models.ReviewStats{
		UserID: userID,
	}

	for _, item := range m.items {
		if item.UserID == userID {
			stats.TotalReviewItems++

			if item.NextReviewDate != nil {
				priority := item.GetPriority(now)
				switch priority {
				case models.PriorityUrgent:
					stats.UrgentItems++
				case models.PriorityRecommended:
					stats.RecommendedItems++
				case models.PriorityRelaxed:
					stats.RelaxedItems++
				}
			}
		}
	}

	return stats, nil
}

func (m *MockSRSService) CreateReviewItem(ctx context.Context, data *PhraseData) (uuid.UUID, error) {
	itemID := uuid.New()
	item := &models.ReviewItem{
		ID:          itemID,
		UserID:      data.UserID,
		BookID:      data.BookID,
		PageNumber:  data.PageNumber,
		ItemType:    "phrase",
		Content:     data.Content,
		Translation: data.Translation,
		ReviewCount: 0,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	m.items[itemID] = item
	return itemID, nil
}

func (m *MockSRSService) BulkCreateReviewItems(ctx context.Context, phrases []*PhraseData) ([]uuid.UUID, error) {
	itemIDs := make([]uuid.UUID, 0, len(phrases))

	for _, phrase := range phrases {
		itemID, err := m.CreateReviewItem(ctx, phrase)
		if err != nil {
			return nil, err
		}
		itemIDs = append(itemIDs, itemID)
	}

	return itemIDs, nil
}
