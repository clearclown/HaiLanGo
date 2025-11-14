package repository

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/clearclown/HaiLanGo/backend/internal/models"
	"github.com/google/uuid"
)

// LearningRepositoryInterface はLearning APIのインターフェース
type LearningRepositoryInterface interface {
	GetPageLearning(ctx context.Context, userID, bookID uuid.UUID, pageNumber int) (*models.PageLearning, error)
	CompletePage(ctx context.Context, userID, bookID uuid.UUID, pageNumber int, req *models.CompletePageRequest) (*models.PageProgressDetail, error)
	RecordSession(ctx context.Context, userID, bookID uuid.UUID, pageNumber int, req *models.SessionRequest) (*models.SessionResponse, error)
	GetBookProgress(ctx context.Context, userID, bookID uuid.UUID) (*models.BookProgressSummary, error)
}

// InMemoryLearningRepository はInMemory実装
type InMemoryLearningRepository struct {
	pages    map[string]*models.PageWithOCR
	progress map[string]*models.PageProgressRecord // key: userID:bookID:pageNumber
	phrases  map[string][]models.Phrase             // key: bookID:pageNumber
	sessions map[string]*models.SessionResponse
	mu       sync.RWMutex
}

// NewInMemoryLearningRepository は新しいInMemoryLearningRepositoryを作成
func NewInMemoryLearningRepository() *InMemoryLearningRepository {
	repo := &InMemoryLearningRepository{
		pages:    make(map[string]*models.PageWithOCR),
		progress: make(map[string]*models.PageProgressRecord),
		phrases:  make(map[string][]models.Phrase),
		sessions: make(map[string]*models.SessionResponse),
	}

	// サンプルデータ初期化
	repo.initSampleData()

	return repo
}

func (r *InMemoryLearningRepository) initSampleData() {
	testBookID := "550e8400-e29b-41d4-a716-446655440000"
	testUserID := "550e8400-e29b-41d4-a716-446655440000"

	// サンプルページデータ（1-150ページ）
	for i := 1; i <= 150; i++ {
		pageID := uuid.New().String()
		pageKey := fmt.Sprintf("%s:%d", testBookID, i)

		r.pages[pageKey] = &models.PageWithOCR{
			ID:          pageID,
			BookID:      testBookID,
			PageNumber:  i,
			ImageURL:    fmt.Sprintf("/storage/books/%s/pages/%d.jpg", testBookID, i),
			OCRText:     "Здравствуйте! Как дела?",
			Translation: "こんにちは！元気ですか？",
			Language:    "ru",
			HasAudio:    true,
			AudioURL:    fmt.Sprintf("/storage/books/%s/pages/%d/audio.mp3", testBookID, i),
		}

		// フレーズデータ
		phrasesKey := fmt.Sprintf("%s:%d", testBookID, i)
		r.phrases[phrasesKey] = []models.Phrase{
			{
				ID:            uuid.New().String(),
				Text:          "Здравствуйте!",
				Translation:   "こんにちは",
				Pronunciation: "zdravstvuyte",
				AudioURL:      fmt.Sprintf("/storage/books/%s/pages/%d/phrase-1.mp3", testBookID, i),
			},
			{
				ID:            uuid.New().String(),
				Text:          "Как дела?",
				Translation:   "元気ですか？",
				Pronunciation: "kak dela",
				AudioURL:      fmt.Sprintf("/storage/books/%s/pages/%d/phrase-2.mp3", testBookID, i),
			},
		}

		// 進捗データ（最初の45ページは完了済み）
		if i <= 45 {
			progressKey := fmt.Sprintf("%s:%s:%d", testUserID, testBookID, i)
			now := time.Now()
			r.progress[progressKey] = &models.PageProgressRecord{
				ID:            uuid.New().String(),
				UserID:        testUserID,
				BookID:        testBookID,
				PageNumber:    i,
				IsCompleted:   true,
				CompletedAt:   &now,
				StudyTime:     60 + i, // 各ページ60秒+α
				ReviewCount:   1,
				LastStudiedAt: &now,
				CreatedAt:     now,
				UpdatedAt:     now,
			}
		}
	}
}

func (r *InMemoryLearningRepository) GetPageLearning(ctx context.Context, userID, bookID uuid.UUID, pageNumber int) (*models.PageLearning, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	// ページデータを取得
	pageKey := fmt.Sprintf("%s:%d", bookID.String(), pageNumber)
	page, exists := r.pages[pageKey]
	if !exists {
		return nil, fmt.Errorf("page not found")
	}

	// 進捗データを取得
	progressKey := fmt.Sprintf("%s:%s:%d", userID.String(), bookID.String(), pageNumber)
	progress := r.progress[progressKey]

	progressDetail := models.PageProgressDetail{
		IsCompleted:   false,
		CompletedAt:   nil,
		StudyTime:     0,
		ReviewCount:   0,
		LastStudiedAt: nil,
	}

	if progress != nil {
		progressDetail = models.PageProgressDetail{
			IsCompleted:   progress.IsCompleted,
			CompletedAt:   progress.CompletedAt,
			StudyTime:     progress.StudyTime,
			ReviewCount:   progress.ReviewCount,
			LastStudiedAt: progress.LastStudiedAt,
		}
	}

	// フレーズデータを取得
	phrasesKey := fmt.Sprintf("%s:%d", bookID.String(), pageNumber)
	phrases := r.phrases[phrasesKey]
	if phrases == nil {
		phrases = []models.Phrase{}
	}

	// 単語データ（サンプル）
	vocabulary := []models.VocabularyItem{
		{
			Word:         "Здравствуйте",
			Translation:  "こんにちは",
			PartOfSpeech: "interjection",
			Frequency:    "common",
		},
		{
			Word:         "дела",
			Translation:  "事柄、状況",
			PartOfSpeech: "noun",
			Frequency:    "common",
		},
	}

	// ナビゲーション情報
	navigation := models.NavigationInfo{
		HasPrevious: pageNumber > 1,
		HasNext:     pageNumber < 150,
		TotalPages:  150,
		CurrentPage: pageNumber,
	}

	return &models.PageLearning{
		Page:       *page,
		Progress:   progressDetail,
		Phrases:    phrases,
		Vocabulary: vocabulary,
		Navigation: navigation,
	}, nil
}

func (r *InMemoryLearningRepository) CompletePage(ctx context.Context, userID, bookID uuid.UUID, pageNumber int, req *models.CompletePageRequest) (*models.PageProgressDetail, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	progressKey := fmt.Sprintf("%s:%s:%d", userID.String(), bookID.String(), pageNumber)
	progress := r.progress[progressKey]

	now := time.Now()
	if progress == nil {
		// 新規作成
		progress = &models.PageProgressRecord{
			ID:            uuid.New().String(),
			UserID:        userID.String(),
			BookID:        bookID.String(),
			PageNumber:    pageNumber,
			IsCompleted:   true,
			CompletedAt:   &now,
			StudyTime:     req.StudyTime,
			ReviewCount:   1,
			LastStudiedAt: &now,
			Notes:         req.Notes,
			CreatedAt:     now,
			UpdatedAt:     now,
		}
		r.progress[progressKey] = progress
	} else {
		// 更新
		progress.IsCompleted = true
		progress.CompletedAt = &now
		progress.StudyTime += req.StudyTime
		progress.ReviewCount++
		progress.LastStudiedAt = &now
		progress.Notes = req.Notes
		progress.UpdatedAt = now
	}

	return &models.PageProgressDetail{
		IsCompleted:   progress.IsCompleted,
		CompletedAt:   progress.CompletedAt,
		StudyTime:     progress.StudyTime,
		ReviewCount:   progress.ReviewCount,
		LastStudiedAt: progress.LastStudiedAt,
	}, nil
}

func (r *InMemoryLearningRepository) RecordSession(ctx context.Context, userID, bookID uuid.UUID, pageNumber int, req *models.SessionRequest) (*models.SessionResponse, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	sessionKey := uuid.New().String()

	if req.Action == "start" {
		session := &models.SessionResponse{
			SessionID: sessionKey,
			StartedAt: req.Timestamp,
			EndedAt:   nil,
		}
		r.sessions[sessionKey] = session
		return session, nil
	}

	// end の場合
	// 実際にはsession_idをリクエストで受け取るべきだが、簡易実装
	return &models.SessionResponse{
		SessionID: sessionKey,
		StartedAt: req.Timestamp.Add(-5 * time.Minute),
		EndedAt:   &req.Timestamp,
	}, nil
}

func (r *InMemoryLearningRepository) GetBookProgress(ctx context.Context, userID, bookID uuid.UUID) (*models.BookProgressSummary, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	totalPages := 150
	completedPages := 0
	totalStudyTime := 0
	var lastStudiedAt *time.Time
	pages := []models.PageProgressSummaryItem{}

	for i := 1; i <= totalPages; i++ {
		progressKey := fmt.Sprintf("%s:%s:%d", userID.String(), bookID.String(), i)
		progress := r.progress[progressKey]

		if progress != nil {
			if progress.IsCompleted {
				completedPages++
			}
			totalStudyTime += progress.StudyTime

			if lastStudiedAt == nil || (progress.LastStudiedAt != nil && progress.LastStudiedAt.After(*lastStudiedAt)) {
				lastStudiedAt = progress.LastStudiedAt
			}

			pages = append(pages, models.PageProgressSummaryItem{
				PageNumber:  i,
				IsCompleted: progress.IsCompleted,
				StudyTime:   progress.StudyTime,
				ReviewCount: progress.ReviewCount,
			})
		} else {
			pages = append(pages, models.PageProgressSummaryItem{
				PageNumber:  i,
				IsCompleted: false,
				StudyTime:   0,
				ReviewCount: 0,
			})
		}
	}

	completionPercentage := float64(completedPages) / float64(totalPages) * 100
	averageTimePerPage := 0.0
	if completedPages > 0 {
		averageTimePerPage = float64(totalStudyTime) / float64(completedPages)
	}

	currentPage := completedPages + 1
	if currentPage > totalPages {
		currentPage = totalPages
	}

	return &models.BookProgressSummary{
		BookID:               bookID.String(),
		TotalPages:           totalPages,
		CompletedPages:       completedPages,
		CompletionPercentage: completionPercentage,
		TotalStudyTime:       totalStudyTime,
		AverageTimePerPage:   averageTimePerPage,
		CurrentPage:          currentPage,
		LastStudiedAt:        lastStudiedAt,
		Pages:                pages,
	}, nil
}
