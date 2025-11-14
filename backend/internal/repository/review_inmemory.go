package repository

import (
	"context"
	"sync"
	"time"

	"github.com/clearclown/HaiLanGo/backend/internal/models"
	"github.com/google/uuid"
)

// InMemoryReviewRepository はインメモリの復習リポジトリ
type InMemoryReviewRepository struct {
	items     map[string]*models.ReviewItem
	histories map[string]*models.ReviewHistory
	mu        sync.RWMutex
}

// NewInMemoryReviewRepository は新しいInMemoryReviewRepositoryを作成
func NewInMemoryReviewRepository() *InMemoryReviewRepository {
	repo := &InMemoryReviewRepository{
		items:     make(map[string]*models.ReviewItem),
		histories: make(map[string]*models.ReviewHistory),
	}

	// サンプルデータを初期化
	repo.initSampleData()

	return repo
}

// initSampleData はサンプルデータを初期化
func (r *InMemoryReviewRepository) initSampleData() {
	// サンプル復習アイテムを作成
	sampleUserID := "550e8400-e29b-41d4-a716-446655440001"
	sampleBookID := "550e8400-e29b-41d4-a716-446655440000"

	now := time.Now()

	// 緊急: 期限切れアイテム (3個)
	for i := 0; i < 3; i++ {
		id := uuid.New().String()
		r.items[id] = &models.ReviewItem{
			ID:           id,
			UserID:       sampleUserID,
			BookID:       sampleBookID,
			PageNumber:   i + 1,
			Type:         "word",
			Text:         []string{"Здравствуйте", "Спасибо", "До свидания"}[i],
			Translation:  []string{"こんにちは", "ありがとう", "さようなら"}[i],
			Language:     "ru",
			MasteryLevel: 30 + i*10,
			IntervalDays: 1,
			EaseFactor:   2.5,
			LastReviewed: now.Add(-25 * time.Hour), // 昨日より前
			NextReview:   now.Add(-1 * time.Hour),  // 1時間前 (期限切れ)
			ReviewCount:  i + 1,
			CreatedAt:    now.Add(-48 * time.Hour),
			UpdatedAt:    now,
		}
	}

	// 推奨: 明日までのアイテム (5個)
	for i := 0; i < 5; i++ {
		id := uuid.New().String()
		r.items[id] = &models.ReviewItem{
			ID:           id,
			UserID:       sampleUserID,
			BookID:       sampleBookID,
			PageNumber:   i + 10,
			Type:         "word",
			Text:         []string{"Привет", "Пока", "Да", "Нет", "Может быть"}[i],
			Translation:  []string{"やあ", "じゃあね", "はい", "いいえ", "たぶん"}[i],
			Language:     "ru",
			MasteryLevel: 50 + i*5,
			IntervalDays: 2,
			EaseFactor:   2.6,
			LastReviewed: now.Add(-20 * time.Hour),
			NextReview:   now.Add(36 * time.Hour), // 明日
			ReviewCount:  i + 2,
			CreatedAt:    now.Add(-72 * time.Hour),
			UpdatedAt:    now,
		}
	}

	// 余裕あり: 2日以上先のアイテム (4個)
	for i := 0; i < 4; i++ {
		id := uuid.New().String()
		r.items[id] = &models.ReviewItem{
			ID:           id,
			UserID:       sampleUserID,
			BookID:       sampleBookID,
			PageNumber:   i + 20,
			Type:         "phrase",
			Text:         []string{"Как дела?", "Всё хорошо", "Извините", "Пожалуйста"}[i],
			Translation:  []string{"元気ですか？", "元気です", "すみません", "どういたしまして"}[i],
			Language:     "ru",
			MasteryLevel: 70 + i*5,
			IntervalDays: 7,
			EaseFactor:   2.8,
			LastReviewed: now.Add(-100 * time.Hour),
			NextReview:   now.Add(72 * time.Hour), // 3日後
			ReviewCount:  i + 5,
			CreatedAt:    now.Add(-168 * time.Hour),
			UpdatedAt:    now,
		}
	}
}

func (r *InMemoryReviewRepository) Create(ctx context.Context, item *models.ReviewItem) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if item.ID == "" {
		item.ID = uuid.New().String()
	}

	now := time.Now()
	item.CreatedAt = now
	item.UpdatedAt = now

	r.items[item.ID] = item
	return nil
}

func (r *InMemoryReviewRepository) FindByID(ctx context.Context, id string) (*models.ReviewItem, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	item, exists := r.items[id]
	if !exists {
		return nil, ErrReviewItemNotFound
	}

	return item, nil
}

func (r *InMemoryReviewRepository) FindByUserID(ctx context.Context, userID string) ([]*models.ReviewItem, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var items []*models.ReviewItem
	for _, item := range r.items {
		if item.UserID == userID {
			items = append(items, item)
		}
	}

	return items, nil
}

func (r *InMemoryReviewRepository) Update(ctx context.Context, item *models.ReviewItem) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.items[item.ID]; !exists {
		return ErrReviewItemNotFound
	}

	item.UpdatedAt = time.Now()
	r.items[item.ID] = item
	return nil
}

func (r *InMemoryReviewRepository) Delete(ctx context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.items[id]; !exists {
		return ErrReviewItemNotFound
	}

	delete(r.items, id)
	return nil
}

func (r *InMemoryReviewRepository) CountCompletedToday(ctx context.Context, userID string, since time.Time) (int, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	count := 0
	for _, history := range r.histories {
		if history.UserID == userID && history.ReviewedAt.After(since) {
			count++
		}
	}

	return count, nil
}

func (r *InMemoryReviewRepository) CountCompletedSince(ctx context.Context, userID string, since time.Time) (int, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	count := 0
	for _, history := range r.histories {
		if history.UserID == userID && history.ReviewedAt.After(since) {
			count++
		}
	}

	return count, nil
}

func (r *InMemoryReviewRepository) SaveHistory(ctx context.Context, history *models.ReviewHistory) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if history.ID == "" {
		history.ID = uuid.New().String()
	}

	r.histories[history.ID] = history
	return nil
}
