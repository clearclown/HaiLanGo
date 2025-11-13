package repository

import (
	"context"
	"time"

	"github.com/clearclown/HaiLanGo/backend/internal/models"
	"github.com/google/uuid"
)

// MockReviewItemRepository はテスト用のモックリポジトリ
type MockReviewItemRepository struct {
	items     map[uuid.UUID]*models.ReviewItem
	histories []models.ReviewHistory
}

// NewMockReviewItemRepository はモックリポジトリを作成
func NewMockReviewItemRepository() *MockReviewItemRepository {
	return &MockReviewItemRepository{
		items:     make(map[uuid.UUID]*models.ReviewItem),
		histories: make([]models.ReviewHistory, 0),
	}
}

// Create は復習項目を作成
func (m *MockReviewItemRepository) Create(ctx context.Context, item *models.ReviewItem) error {
	m.items[item.ID] = item
	return nil
}

// GetByID はIDで復習項目を取得
func (m *MockReviewItemRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.ReviewItem, error) {
	item, ok := m.items[id]
	if !ok {
		return nil, nil
	}
	return item, nil
}

// GetByUserID はユーザーIDで復習項目を取得
func (m *MockReviewItemRepository) GetByUserID(ctx context.Context, userID uuid.UUID) ([]*models.ReviewItem, error) {
	items := make([]*models.ReviewItem, 0)
	for _, item := range m.items {
		if item.UserID == userID {
			items = append(items, item)
		}
	}
	return items, nil
}

// GetDueItems は今日復習すべき項目を取得
func (m *MockReviewItemRepository) GetDueItems(ctx context.Context, userID uuid.UUID, now time.Time) ([]*models.ReviewItem, error) {
	items := make([]*models.ReviewItem, 0)
	for _, item := range m.items {
		if item.UserID == userID && item.NextReviewDate != nil {
			if item.NextReviewDate.Before(now) || isSameDay(*item.NextReviewDate, now) {
				items = append(items, item)
			}
		}
	}
	return items, nil
}

// Update は復習項目を更新
func (m *MockReviewItemRepository) Update(ctx context.Context, item *models.ReviewItem) error {
	m.items[item.ID] = item
	return nil
}

// AddHistory は復習履歴を追加
func (m *MockReviewItemRepository) AddHistory(ctx context.Context, history *models.ReviewHistory) error {
	m.histories = append(m.histories, *history)
	return nil
}

// GetHistoriesByItemID は復習項目の履歴を取得
func (m *MockReviewItemRepository) GetHistoriesByItemID(ctx context.Context, itemID uuid.UUID) ([]models.ReviewHistory, error) {
	histories := make([]models.ReviewHistory, 0)
	for _, h := range m.histories {
		if h.ReviewItemID == itemID {
			histories = append(histories, h)
		}
	}
	return histories, nil
}

// GetStats は統計情報を取得
func (m *MockReviewItemRepository) GetStats(ctx context.Context, userID uuid.UUID, now time.Time) (*models.ReviewStats, error) {
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

// isSameDay は2つの時刻が同じ日かを判定する
func isSameDay(t1, t2 time.Time) bool {
	y1, m1, d1 := t1.Date()
	y2, m2, d2 := t2.Date()
	return y1 == y2 && m1 == m2 && d1 == d2
}
