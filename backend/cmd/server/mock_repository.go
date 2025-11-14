package main

import (
	"context"
	"fmt"
	"time"

	"github.com/clearclown/HaiLanGo/backend/internal/models"
	"github.com/clearclown/HaiLanGo/backend/internal/service/learning"
	"github.com/google/uuid"
)

// MockRepository はテスト用のモックリポジトリ
type MockRepository struct {
	pages      map[string]*models.PageWithProgress
	progress   map[string]*models.LearningProgress
	histories  map[string]*models.LearningHistory
}

// NewMockRepository は新しいMockRepositoryを作成
func NewMockRepository() *MockRepository {
	repo := &MockRepository{
		pages:     make(map[string]*models.PageWithProgress),
		progress:  make(map[string]*models.LearningProgress),
		histories: make(map[string]*models.LearningHistory),
	}

	// サンプルデータの初期化
	repo.initSampleData()

	return repo
}

// initSampleData はサンプルデータを初期化
func (r *MockRepository) initSampleData() {
	bookID := uuid.MustParse("550e8400-e29b-41d4-a716-446655440000")

	// ページ1-10のサンプルデータ
	for i := 1; i <= 10; i++ {
		pageID := uuid.New()
		key := fmt.Sprintf("%s-%d", bookID.String(), i)

		page := &models.Page{
			ID:            pageID,
			BookID:        bookID,
			PageNumber:    i,
			ImageURL:      fmt.Sprintf("https://example.com/page%d.png", i),
			OCRText:       fmt.Sprintf("Здравствуйте! Это страница %d.", i),
			OCRConfidence: 0.95,
			DetectedLang:  "ru",
			OCRStatus:     models.OCRStatusCompleted,
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
		}
		r.pages[key] = &models.PageWithProgress{
			Page:        page,
			IsCompleted: false,
			CompletedAt: nil,
		}
	}

	// 進捗データ
	r.progress[bookID.String()] = &models.LearningProgress{
		BookID:         bookID,
		TotalPages:     150,
		CompletedPages: 0,
		Progress:       0.0,
		TotalStudyTime: 0,
	}
}

// GetPage はページを取得
func (r *MockRepository) GetPage(ctx context.Context, bookID uuid.UUID, pageNumber int) (*models.PageWithProgress, error) {
	key := fmt.Sprintf("%s-%d", bookID.String(), pageNumber)

	page, exists := r.pages[key]
	if !exists {
		return nil, learning.ErrPageNotFound
	}

	return page, nil
}

// MarkPageCompleted はページを完了としてマーク
func (r *MockRepository) MarkPageCompleted(ctx context.Context, userID, bookID uuid.UUID, pageNumber int, studyTime int) error {
	key := fmt.Sprintf("%s-%d", bookID.String(), pageNumber)

	page, exists := r.pages[key]
	if !exists {
		return learning.ErrPageNotFound
	}

	// ページを完了としてマーク
	now := time.Now()
	page.IsCompleted = true
	page.CompletedAt = &now

	// 進捗を更新
	if progress, exists := r.progress[bookID.String()]; exists {
		progress.CompletedPages++
		progress.TotalStudyTime += studyTime
		progress.Progress = float64(progress.CompletedPages) / float64(progress.TotalPages) * 100
	}

	// 学習履歴を保存
	historyKey := fmt.Sprintf("%s-%s-%d", userID.String(), bookID.String(), pageNumber)
	r.histories[historyKey] = &models.LearningHistory{
		ID:          uuid.New(),
		UserID:      userID,
		BookID:      bookID,
		PageID:      page.ID,
		PageNumber:  pageNumber,
		IsCompleted: true,
		StudyTime:   studyTime,
		CompletedAt: &now,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	return nil
}

// GetProgress は学習進捗を取得
func (r *MockRepository) GetProgress(ctx context.Context, userID, bookID uuid.UUID) (*models.LearningProgress, error) {
	progress, exists := r.progress[bookID.String()]
	if !exists {
		return &models.LearningProgress{
			BookID:         bookID,
			TotalPages:     150,
			CompletedPages: 0,
			Progress:       0.0,
			TotalStudyTime: 0,
		}, nil
	}

	return progress, nil
}
