package srs

import (
	"context"
	"time"

	"github.com/clearclown/HaiLanGo/backend/internal/models"
	"github.com/clearclown/HaiLanGo/backend/pkg/srs"
	"github.com/google/uuid"
)

// SRSService は間隔反復学習サービス
type SRSService struct {
	repo ReviewItemRepository
}

// ReviewItemRepository は復習項目のリポジトリインターフェース
type ReviewItemRepository interface {
	Create(ctx context.Context, item *models.ReviewItem) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.ReviewItem, error)
	GetByUserID(ctx context.Context, userID uuid.UUID) ([]*models.ReviewItem, error)
	GetDueItems(ctx context.Context, userID uuid.UUID, now time.Time) ([]*models.ReviewItem, error)
	Update(ctx context.Context, item *models.ReviewItem) error
	AddHistory(ctx context.Context, history *models.ReviewHistory) error
	GetHistoriesByItemID(ctx context.Context, itemID uuid.UUID) ([]models.ReviewHistory, error)
	GetStats(ctx context.Context, userID uuid.UUID, now time.Time) (*models.ReviewStats, error)
}

// ReviewItemsByPriority は優先度別の復習項目
type ReviewItemsByPriority struct {
	UrgentItems      []*models.ReviewItem `json:"urgent_items"`
	RecommendedItems []*models.ReviewItem `json:"recommended_items"`
	RelaxedItems     []*models.ReviewItem `json:"relaxed_items"`
}

// PhraseData はフレーズデータ
type PhraseData struct {
	UserID      uuid.UUID
	BookID      uuid.UUID
	PageNumber  int
	Content     string
	Translation string
}

// NewSRSService は新しいSRSServiceを作成
func NewSRSService(repo ReviewItemRepository) *SRSService {
	return &SRSService{
		repo: repo,
	}
}

// GetReviewItems は優先度別に復習項目を取得
func (s *SRSService) GetReviewItems(ctx context.Context, userID uuid.UUID, now time.Time) (*ReviewItemsByPriority, error) {
	// ユーザーのすべての復習項目を取得
	items, err := s.repo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	result := &ReviewItemsByPriority{
		UrgentItems:      make([]*models.ReviewItem, 0),
		RecommendedItems: make([]*models.ReviewItem, 0),
		RelaxedItems:     make([]*models.ReviewItem, 0),
	}

	// 優先度で分類
	for _, item := range items {
		if item.NextReviewDate == nil {
			// 次回復習日が未設定の項目は緊急として扱う
			result.UrgentItems = append(result.UrgentItems, item)
			continue
		}

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

	return result, nil
}

// CompleteReview は復習完了処理
func (s *SRSService) CompleteReview(ctx context.Context, itemID uuid.UUID, score int, timeSpentSec int) error {
	// 復習項目を取得
	item, err := s.repo.GetByID(ctx, itemID)
	if err != nil {
		return err
	}
	if item == nil {
		return nil // 項目が見つからない場合は何もしない
	}

	now := time.Now()

	// 次回復習日を計算
	nextReviewDate := srs.CalculateNextReviewDate(item.ReviewCount, score, now)

	// 復習項目を更新
	item.ReviewCount++
	item.LastScore = score
	item.LastReviewDate = &now
	item.NextReviewDate = &nextReviewDate
	item.UpdatedAt = now

	err = s.repo.Update(ctx, item)
	if err != nil {
		return err
	}

	// 復習履歴を記録
	history := &models.ReviewHistory{
		ID:           uuid.New(),
		ReviewItemID: itemID,
		UserID:       item.UserID,
		Score:        score,
		ReviewedAt:   now,
		TimeSpentSec: timeSpentSec,
	}

	return s.repo.AddHistory(ctx, history)
}

// GetStats は復習統計を取得
func (s *SRSService) GetStats(ctx context.Context, userID uuid.UUID, now time.Time) (*models.ReviewStats, error) {
	return s.repo.GetStats(ctx, userID, now)
}

// CreateReviewItem はフレーズから復習項目を作成
func (s *SRSService) CreateReviewItem(ctx context.Context, data *PhraseData) (uuid.UUID, error) {
	now := time.Now()

	item := &models.ReviewItem{
		ID:          uuid.New(),
		UserID:      data.UserID,
		BookID:      data.BookID,
		PageNumber:  data.PageNumber,
		ItemType:    "phrase",
		Content:     data.Content,
		Translation: data.Translation,
		ReviewCount: 0,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	err := s.repo.Create(ctx, item)
	if err != nil {
		return uuid.Nil, err
	}

	return item.ID, nil
}

// BulkCreateReviewItems は複数の復習項目を一括作成
func (s *SRSService) BulkCreateReviewItems(ctx context.Context, phrases []*PhraseData) ([]uuid.UUID, error) {
	itemIDs := make([]uuid.UUID, 0, len(phrases))

	for _, phrase := range phrases {
		itemID, err := s.CreateReviewItem(ctx, phrase)
		if err != nil {
			return nil, err
		}
		itemIDs = append(itemIDs, itemID)
	}

	return itemIDs, nil
}

// GetDueItems は今日復習すべき項目を取得
func (s *SRSService) GetDueItems(ctx context.Context, userID uuid.UUID, now time.Time) ([]*models.ReviewItem, error) {
	return s.repo.GetDueItems(ctx, userID, now)
}

// GetReviewHistory は復習履歴を取得
func (s *SRSService) GetReviewHistory(ctx context.Context, itemID uuid.UUID) ([]models.ReviewHistory, error) {
	return s.repo.GetHistoriesByItemID(ctx, itemID)
}
