package repository

import (
	"context"
	"testing"
	"time"

	"github.com/clearclown/HaiLanGo/backend/internal/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestCreateReviewItem は復習項目の作成をテスト
func TestCreateReviewItem(t *testing.T) {
	ctx := context.Background()
	repo := NewMockReviewItemRepository()

	userID := uuid.New()
	bookID := uuid.New()

	item := &models.ReviewItem{
		ID:          uuid.New(),
		UserID:      userID,
		BookID:      bookID,
		PageNumber:  1,
		ItemType:    "phrase",
		Content:     "Здравствуйте!",
		Translation: "こんにちは",
		ReviewCount: 0,
	}

	err := repo.Create(ctx, item)
	require.NoError(t, err)

	// 作成された項目を取得
	created, err := repo.GetByID(ctx, item.ID)
	require.NoError(t, err)
	assert.Equal(t, item.ID, created.ID)
	assert.Equal(t, item.Content, created.Content)
	assert.Equal(t, 0, created.ReviewCount)
}

// TestGetReviewItemsByUserID はユーザーの復習項目取得をテスト
func TestGetReviewItemsByUserID(t *testing.T) {
	ctx := context.Background()
	repo := NewMockReviewItemRepository()

	userID := uuid.New()
	bookID := uuid.New()

	// 複数の復習項目を作成
	for i := 0; i < 5; i++ {
		item := &models.ReviewItem{
			ID:          uuid.New(),
			UserID:      userID,
			BookID:      bookID,
			PageNumber:  i + 1,
			ItemType:    "phrase",
			Content:     "Test content",
			Translation: "テスト内容",
			ReviewCount: 0,
		}
		err := repo.Create(ctx, item)
		require.NoError(t, err)
	}

	// ユーザーの復習項目を取得
	items, err := repo.GetByUserID(ctx, userID)
	require.NoError(t, err)
	assert.Len(t, items, 5)
}

// TestGetReviewItemsDueToday は今日復習すべき項目の取得をテスト
func TestGetReviewItemsDueToday(t *testing.T) {
	ctx := context.Background()
	repo := NewMockReviewItemRepository()

	userID := uuid.New()
	bookID := uuid.New()
	now := time.Now()

	// 過去の復習日の項目（復習が必要）
	pastDueItem := &models.ReviewItem{
		ID:             uuid.New(),
		UserID:         userID,
		BookID:         bookID,
		PageNumber:     1,
		ItemType:       "phrase",
		Content:        "Past due",
		Translation:    "期限切れ",
		ReviewCount:    1,
		NextReviewDate: &[]time.Time{now.AddDate(0, 0, -1)}[0],
	}

	// 今日の復習日の項目（復習が必要）
	todayItem := &models.ReviewItem{
		ID:             uuid.New(),
		UserID:         userID,
		BookID:         bookID,
		PageNumber:     2,
		ItemType:       "phrase",
		Content:        "Today",
		Translation:    "今日",
		ReviewCount:    1,
		NextReviewDate: &now,
	}

	// 未来の復習日の項目（復習不要）
	futureItem := &models.ReviewItem{
		ID:             uuid.New(),
		UserID:         userID,
		BookID:         bookID,
		PageNumber:     3,
		ItemType:       "phrase",
		Content:        "Future",
		Translation:    "未来",
		ReviewCount:    1,
		NextReviewDate: &[]time.Time{now.AddDate(0, 0, 3)}[0],
	}

	require.NoError(t, repo.Create(ctx, pastDueItem))
	require.NoError(t, repo.Create(ctx, todayItem))
	require.NoError(t, repo.Create(ctx, futureItem))

	// 今日復習すべき項目を取得
	dueItems, err := repo.GetDueItems(ctx, userID, now)
	require.NoError(t, err)
	assert.Len(t, dueItems, 2) // 過去と今日の2つ
}

// TestUpdateReviewItem は復習項目の更新をテスト
func TestUpdateReviewItem(t *testing.T) {
	ctx := context.Background()
	repo := NewMockReviewItemRepository()

	userID := uuid.New()
	bookID := uuid.New()

	item := &models.ReviewItem{
		ID:          uuid.New(),
		UserID:      userID,
		BookID:      bookID,
		PageNumber:  1,
		ItemType:    "phrase",
		Content:     "Original",
		Translation: "オリジナル",
		ReviewCount: 0,
	}

	require.NoError(t, repo.Create(ctx, item))

	// 復習完了後の更新
	now := time.Now()
	nextReviewDate := now.AddDate(0, 0, 3)
	item.ReviewCount = 1
	item.LastReviewDate = &now
	item.NextReviewDate = &nextReviewDate
	item.LastScore = 85

	err := repo.Update(ctx, item)
	require.NoError(t, err)

	// 更新された項目を取得
	updated, err := repo.GetByID(ctx, item.ID)
	require.NoError(t, err)
	assert.Equal(t, 1, updated.ReviewCount)
	assert.Equal(t, 85, updated.LastScore)
	assert.NotNil(t, updated.NextReviewDate)
}

// TestAddReviewHistory は復習履歴の追加をテスト
func TestAddReviewHistory(t *testing.T) {
	ctx := context.Background()
	repo := NewMockReviewItemRepository()

	reviewItemID := uuid.New()
	userID := uuid.New()

	history := &models.ReviewHistory{
		ID:           uuid.New(),
		ReviewItemID: reviewItemID,
		UserID:       userID,
		Score:        85,
		ReviewedAt:   time.Now(),
		TimeSpentSec: 30,
	}

	err := repo.AddHistory(ctx, history)
	require.NoError(t, err)

	// 履歴を取得
	histories, err := repo.GetHistoriesByItemID(ctx, reviewItemID)
	require.NoError(t, err)
	assert.Len(t, histories, 1)
	assert.Equal(t, 85, histories[0].Score)
}

// TestGetReviewStats は統計情報の取得をテスト
func TestGetReviewStats(t *testing.T) {
	ctx := context.Background()
	repo := NewMockReviewItemRepository()

	userID := uuid.New()
	bookID := uuid.New()
	now := time.Now()

	// 緊急項目を3つ作成
	for i := 0; i < 3; i++ {
		item := &models.ReviewItem{
			ID:             uuid.New(),
			UserID:         userID,
			BookID:         bookID,
			PageNumber:     i + 1,
			ItemType:       "phrase",
			Content:        "Urgent",
			Translation:    "緊急",
			ReviewCount:    1,
			NextReviewDate: &[]time.Time{now.AddDate(0, 0, -1)}[0],
		}
		require.NoError(t, repo.Create(ctx, item))
	}

	// 推奨項目を2つ作成
	for i := 0; i < 2; i++ {
		item := &models.ReviewItem{
			ID:             uuid.New(),
			UserID:         userID,
			BookID:         bookID,
			PageNumber:     i + 10,
			ItemType:       "phrase",
			Content:        "Recommended",
			Translation:    "推奨",
			ReviewCount:    1,
			NextReviewDate: &[]time.Time{now.AddDate(0, 0, 1)}[0],
		}
		require.NoError(t, repo.Create(ctx, item))
	}

	// 統計情報を取得
	stats, err := repo.GetStats(ctx, userID, now)
	require.NoError(t, err)
	assert.Equal(t, 5, stats.TotalReviewItems)
	assert.Equal(t, 3, stats.UrgentItems)
	assert.Equal(t, 2, stats.RecommendedItems)
}
