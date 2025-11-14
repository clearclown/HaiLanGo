# CRITICAL_05: Learning API実装（ページバイページ学習モード）

**優先度**: P0（最高優先度）
**担当者**: 未割当
**見積時間**: 6-8時間
**ブロッカー**: フロントエンド（学習画面）がこのAPIを待っている

---

## ⚠️ PM指示

**現状**: ページバイページ学習モードのバックエンドが完全に欠落している。
**期限**: 72時間以内に実装完了すること。
**重要**: この機能はMVPの核心。遅延は許されない。

---

## 📋 実装要件

### エンドポイント仕様

#### 1. GET /api/v1/learning/books/:bookId/pages/:pageNumber
**説明**: 学習ページのデータを取得

**Request**:
```http
GET /api/v1/learning/books/550e8400-e29b-41d4-a716-446655440000/pages/12
Authorization: Bearer <JWT_TOKEN>
```

**Response** (200 OK):
```json
{
  "page": {
    "id": "page-uuid",
    "book_id": "550e8400-e29b-41d4-a716-446655440000",
    "page_number": 12,
    "image_url": "/storage/books/550e8400/pages/12.jpg",
    "ocr_text": "Здравствуйте! Как дела?",
    "translation": "こんにちは！元気ですか？",
    "language": "ru",
    "has_audio": true,
    "audio_url": "/storage/books/550e8400/pages/12/audio.mp3"
  },
  "progress": {
    "is_completed": false,
    "completed_at": null,
    "study_time": 0,
    "review_count": 0,
    "last_studied_at": null
  },
  "phrases": [
    {
      "id": "phrase-1",
      "text": "Здравствуйте!",
      "translation": "こんにちは",
      "pronunciation": "zdravstvuyte",
      "audio_url": "/storage/books/550e8400/pages/12/phrase-1.mp3"
    },
    {
      "id": "phrase-2",
      "text": "Как дела?",
      "translation": "元気ですか？",
      "pronunciation": "kak dela",
      "audio_url": "/storage/books/550e8400/pages/12/phrase-2.mp3"
    }
  ],
  "vocabulary": [
    {
      "word": "Здравствуйте",
      "translation": "こんにちは",
      "part_of_speech": "interjection",
      "frequency": "common"
    }
  ],
  "navigation": {
    "has_previous": true,
    "has_next": true,
    "total_pages": 150,
    "current_page": 12
  }
}
```

#### 2. POST /api/v1/learning/books/:bookId/pages/:pageNumber/complete
**説明**: ページ学習を完了としてマーク

**Request**:
```http
POST /api/v1/learning/books/550e8400/pages/12/complete
Authorization: Bearer <JWT_TOKEN>
Content-Type: application/json

{
  "study_time": 180,
  "notes": "難しかったが理解できた"
}
```

**Response** (200 OK):
```json
{
  "message": "Page marked as completed",
  "progress": {
    "is_completed": true,
    "completed_at": "2025-11-14T10:30:00Z",
    "study_time": 180,
    "review_count": 1
  },
  "next_page": 13
}
```

#### 3. POST /api/v1/learning/books/:bookId/pages/:pageNumber/session
**説明**: 学習セッションを記録（学習開始・終了時）

**Request**:
```http
POST /api/v1/learning/books/550e8400/pages/12/session
Authorization: Bearer <JWT_TOKEN>
Content-Type: application/json

{
  "action": "start",
  "timestamp": "2025-11-14T10:25:00Z"
}
```

**Response** (200 OK):
```json
{
  "session_id": "session-uuid",
  "started_at": "2025-11-14T10:25:00Z"
}
```

#### 4. GET /api/v1/learning/books/:bookId/progress
**説明**: 書籍全体の学習進捗を取得

**Request**:
```http
GET /api/v1/learning/books/550e8400/progress
Authorization: Bearer <JWT_TOKEN>
```

**Response** (200 OK):
```json
{
  "book_id": "550e8400-e29b-41d4-a716-446655440000",
  "total_pages": 150,
  "completed_pages": 45,
  "completion_percentage": 30.0,
  "total_study_time": 2700,
  "average_time_per_page": 60,
  "current_page": 46,
  "last_studied_at": "2025-11-14T10:30:00Z",
  "pages": [
    {
      "page_number": 1,
      "is_completed": true,
      "study_time": 120,
      "review_count": 2
    }
  ]
}
```

---

## 🗃️ データベーススキーマ

### page_progress テーブル
```sql
CREATE TABLE page_progress (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    book_id UUID NOT NULL REFERENCES books(id) ON DELETE CASCADE,
    page_number INT NOT NULL,
    is_completed BOOLEAN DEFAULT FALSE,
    completed_at TIMESTAMP,
    study_time INT DEFAULT 0, -- 秒単位
    review_count INT DEFAULT 0,
    last_studied_at TIMESTAMP,
    notes TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    UNIQUE(user_id, book_id, page_number)
);

CREATE INDEX idx_page_progress_user_book ON page_progress(user_id, book_id);
CREATE INDEX idx_page_progress_completed ON page_progress(user_id, is_completed);
```

### phrases テーブル
```sql
CREATE TABLE phrases (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    book_id UUID NOT NULL REFERENCES books(id) ON DELETE CASCADE,
    page_number INT NOT NULL,
    text TEXT NOT NULL,
    translation TEXT,
    pronunciation TEXT,
    audio_url TEXT,
    start_position INT,
    end_position INT,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_phrases_book_page ON phrases(book_id, page_number);
```

---

## 🏗️ 実装コード

### models/learning.go

```go
package models

import "time"

// PageLearning はページ学習データ
type PageLearning struct {
	Page       PageWithOCR       `json:"page"`
	Progress   PageProgress      `json:"progress"`
	Phrases    []Phrase          `json:"phrases"`
	Vocabulary []VocabularyItem  `json:"vocabulary"`
	Navigation NavigationInfo    `json:"navigation"`
}

// PageWithOCR はOCRデータを含むページ
type PageWithOCR struct {
	ID          string  `json:"id"`
	BookID      string  `json:"book_id"`
	PageNumber  int     `json:"page_number"`
	ImageURL    string  `json:"image_url"`
	OCRText     string  `json:"ocr_text"`
	Translation string  `json:"translation"`
	Language    string  `json:"language"`
	HasAudio    bool    `json:"has_audio"`
	AudioURL    string  `json:"audio_url,omitempty"`
}

// PageProgress はページ進捗
type PageProgress struct {
	ID            string     `json:"id"`
	UserID        string     `json:"user_id"`
	BookID        string     `json:"book_id"`
	PageNumber    int        `json:"page_number"`
	IsCompleted   bool       `json:"is_completed"`
	CompletedAt   *time.Time `json:"completed_at"`
	StudyTime     int        `json:"study_time"` // 秒
	ReviewCount   int        `json:"review_count"`
	LastStudiedAt *time.Time `json:"last_studied_at"`
	Notes         string     `json:"notes,omitempty"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
}

// Phrase はフレーズ
type Phrase struct {
	ID            string `json:"id"`
	BookID        string `json:"book_id"`
	PageNumber    int    `json:"page_number"`
	Text          string `json:"text"`
	Translation   string `json:"translation"`
	Pronunciation string `json:"pronunciation,omitempty"`
	AudioURL      string `json:"audio_url,omitempty"`
	StartPosition int    `json:"start_position,omitempty"`
	EndPosition   int    `json:"end_position,omitempty"`
}

// VocabularyItem は単語
type VocabularyItem struct {
	Word         string `json:"word"`
	Translation  string `json:"translation"`
	PartOfSpeech string `json:"part_of_speech,omitempty"`
	Frequency    string `json:"frequency,omitempty"`
}

// NavigationInfo はナビゲーション情報
type NavigationInfo struct {
	HasPrevious bool `json:"has_previous"`
	HasNext     bool `json:"has_next"`
	TotalPages  int  `json:"total_pages"`
	CurrentPage int  `json:"current_page"`
}

// CompletePageRequest はページ完了リクエスト
type CompletePageRequest struct {
	StudyTime int    `json:"study_time" binding:"required,min=1"`
	Notes     string `json:"notes"`
}

// SessionRequest は学習セッションリクエスト
type SessionRequest struct {
	Action    string    `json:"action" binding:"required,oneof=start end"`
	Timestamp time.Time `json:"timestamp" binding:"required"`
}

// SessionResponse は学習セッションレスポンス
type SessionResponse struct {
	SessionID string    `json:"session_id"`
	StartedAt time.Time `json:"started_at"`
	EndedAt   *time.Time `json:"ended_at,omitempty"`
}

// BookProgressSummary は書籍進捗サマリー
type BookProgressSummary struct {
	BookID               string               `json:"book_id"`
	TotalPages           int                  `json:"total_pages"`
	CompletedPages       int                  `json:"completed_pages"`
	CompletionPercentage float64              `json:"completion_percentage"`
	TotalStudyTime       int                  `json:"total_study_time"`
	AverageTimePerPage   float64              `json:"average_time_per_page"`
	CurrentPage          int                  `json:"current_page"`
	LastStudiedAt        *time.Time           `json:"last_studied_at"`
	Pages                []PageProgressSummary `json:"pages"`
}

// PageProgressSummary はページ進捗サマリー
type PageProgressSummary struct {
	PageNumber  int  `json:"page_number"`
	IsCompleted bool `json:"is_completed"`
	StudyTime   int  `json:"study_time"`
	ReviewCount int  `json:"review_count"`
}
```

### repository/learning_inmemory.go

```go
package repository

import (
	"context"
	"sync"
	"time"

	"github.com/clearclown/HaiLanGo/backend/internal/models"
	"github.com/google/uuid"
)

type LearningRepository interface {
	GetPageLearning(ctx context.Context, userID, bookID uuid.UUID, pageNumber int) (*models.PageLearning, error)
	CompletePageまたは(ctx context.Context, userID, bookID uuid.UUID, pageNumber int, req *models.CompletePageRequest) (*models.PageProgress, error)
	RecordSession(ctx context.Context, userID, bookID uuid.UUID, pageNumber int, req *models.SessionRequest) (*models.SessionResponse, error)
	GetBookProgress(ctx context.Context, userID, bookID uuid.UUID) (*models.BookProgressSummary, error)
	GetPageProgress(ctx context.Context, userID, bookID uuid.UUID, pageNumber int) (*models.PageProgress, error)
	UpdatePageProgress(ctx context.Context, progress *models.PageProgress) error
}

type InMemoryLearningRepository struct {
	pages      map[string]*models.PageWithOCR
	progress   map[string]*models.PageProgress // key: userID:bookID:pageNumber
	phrases    map[string][]models.Phrase      // key: bookID:pageNumber
	sessions   map[string]*models.SessionResponse
	mu         sync.RWMutex
}

func NewInMemoryLearningRepository() *InMemoryLearningRepository {
	repo := &InMemoryLearningRepository{
		pages:    make(map[string]*models.PageWithOCR),
		progress: make(map[string]*models.PageProgress),
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

	// サンプルページデータ（1-5ページ）
	for i := 1; i <= 5; i++ {
		pageID := uuid.New().String()
		key := pageID

		r.pages[key] = &models.PageWithOCR{
			ID:          pageID,
			BookID:      testBookID,
			PageNumber:  i,
			ImageURL:    "/storage/books/" + testBookID + "/pages/" + string(rune(i)) + ".jpg",
			OCRText:     "Здравствуйте! Как дела?",
			Translation: "こんにちは！元気ですか？",
			Language:    "ru",
			HasAudio:    true,
			AudioURL:    "/storage/books/" + testBookID + "/pages/" + string(rune(i)) + "/audio.mp3",
		}

		// フレーズデータ
		phrasesKey := testBookID + ":" + string(rune(i))
		r.phrases[phrasesKey] = []models.Phrase{
			{
				ID:            uuid.New().String(),
				BookID:        testBookID,
				PageNumber:    i,
				Text:          "Здравствуйте!",
				Translation:   "こんにちは",
				Pronunciation: "zdravstvuyte",
				AudioURL:      "/storage/books/" + testBookID + "/pages/" + string(rune(i)) + "/phrase-1.mp3",
			},
			{
				ID:            uuid.New().String(),
				BookID:        testBookID,
				PageNumber:    i,
				Text:          "Как дела?",
				Translation:   "元気ですか？",
				Pronunciation: "kak dela",
				AudioURL:      "/storage/books/" + testBookID + "/pages/" + string(rune(i)) + "/phrase-2.mp3",
			},
		}
	}

	// サンプル進捗データ（ページ1は完了済み）
	progressKey := testUserID + ":" + testBookID + ":1"
	completedAt := time.Now().Add(-24 * time.Hour)
	lastStudied := time.Now().Add(-24 * time.Hour)

	r.progress[progressKey] = &models.PageProgress{
		ID:            uuid.New().String(),
		UserID:        testUserID,
		BookID:        testBookID,
		PageNumber:    1,
		IsCompleted:   true,
		CompletedAt:   &completedAt,
		StudyTime:     120,
		ReviewCount:   2,
		LastStudiedAt: &lastStudied,
		CreatedAt:     time.Now().Add(-48 * time.Hour),
		UpdatedAt:     time.Now().Add(-24 * time.Hour),
	}
}

func (r *InMemoryLearningRepository) GetPageLearning(ctx context.Context, userID, bookID uuid.UUID, pageNumber int) (*models.PageLearning, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	// ページデータを検索
	var page *models.PageWithOCR
	for _, p := range r.pages {
		if p.BookID == bookID.String() && p.PageNumber == pageNumber {
			page = p
			break
		}
	}

	if page == nil {
		return nil, &RepositoryError{Code: "PAGE_NOT_FOUND", Message: "page not found"}
	}

	// 進捗データを取得
	progressKey := userID.String() + ":" + bookID.String() + ":" + string(rune(pageNumber))
	progress, exists := r.progress[progressKey]
	if !exists {
		// 進捗データがなければデフォルト
		progress = &models.PageProgress{
			UserID:      userID.String(),
			BookID:      bookID.String(),
			PageNumber:  pageNumber,
			IsCompleted: false,
		}
	}

	// フレーズデータを取得
	phrasesKey := bookID.String() + ":" + string(rune(pageNumber))
	phrases, _ := r.phrases[phrasesKey]
	if phrases == nil {
		phrases = []models.Phrase{}
	}

	// ナビゲーション情報
	totalPages := 150 // 仮の値
	navigation := models.NavigationInfo{
		HasPrevious: pageNumber > 1,
		HasNext:     pageNumber < totalPages,
		TotalPages:  totalPages,
		CurrentPage: pageNumber,
	}

	return &models.PageLearning{
		Page:       *page,
		Progress:   *progress,
		Phrases:    phrases,
		Vocabulary: []models.VocabularyItem{}, // TODO: 辞書APIと統合
		Navigation: navigation,
	}, nil
}

func (r *InMemoryLearningRepository) CompletePage(ctx context.Context, userID, bookID uuid.UUID, pageNumber int, req *models.CompletePageRequest) (*models.PageProgress, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	progressKey := userID.String() + ":" + bookID.String() + ":" + string(rune(pageNumber))
	progress, exists := r.progress[progressKey]

	now := time.Now()

	if !exists {
		progress = &models.PageProgress{
			ID:         uuid.New().String(),
			UserID:     userID.String(),
			BookID:     bookID.String(),
			PageNumber: pageNumber,
			CreatedAt:  now,
		}
	}

	progress.IsCompleted = true
	progress.CompletedAt = &now
	progress.StudyTime += req.StudyTime
	progress.ReviewCount++
	progress.LastStudiedAt = &now
	progress.Notes = req.Notes
	progress.UpdatedAt = now

	r.progress[progressKey] = progress

	return progress, nil
}

func (r *InMemoryLearningRepository) RecordSession(ctx context.Context, userID, bookID uuid.UUID, pageNumber int, req *models.SessionRequest) (*models.SessionResponse, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	sessionID := uuid.New().String()

	response := &models.SessionResponse{
		SessionID: sessionID,
		StartedAt: req.Timestamp,
	}

	if req.Action == "end" {
		response.EndedAt = &req.Timestamp
	}

	r.sessions[sessionID] = response

	return response, nil
}

func (r *InMemoryLearningRepository) GetBookProgress(ctx context.Context, userID, bookID uuid.UUID) (*models.BookProgressSummary, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	totalPages := 150
	completedPages := 0
	totalStudyTime := 0
	var lastStudied *time.Time
	currentPage := 1

	pages := make([]models.PageProgressSummary, 0)

	for key, prog := range r.progress {
		if prog.UserID == userID.String() && prog.BookID == bookID.String() {
			if prog.IsCompleted {
				completedPages++
			}
			totalStudyTime += prog.StudyTime

			if prog.LastStudiedAt != nil && (lastStudied == nil || prog.LastStudiedAt.After(*lastStudied)) {
				lastStudied = prog.LastStudiedAt
				currentPage = prog.PageNumber + 1
			}

			pages = append(pages, models.PageProgressSummary{
				PageNumber:  prog.PageNumber,
				IsCompleted: prog.IsCompleted,
				StudyTime:   prog.StudyTime,
				ReviewCount: prog.ReviewCount,
			})
		}
	}

	completionPercentage := 0.0
	if totalPages > 0 {
		completionPercentage = float64(completedPages) / float64(totalPages) * 100
	}

	averageTimePerPage := 0.0
	if completedPages > 0 {
		averageTimePerPage = float64(totalStudyTime) / float64(completedPages)
	}

	return &models.BookProgressSummary{
		BookID:               bookID.String(),
		TotalPages:           totalPages,
		CompletedPages:       completedPages,
		CompletionPercentage: completionPercentage,
		TotalStudyTime:       totalStudyTime,
		AverageTimePerPage:   averageTimePerPage,
		CurrentPage:          currentPage,
		LastStudiedAt:        lastStudied,
		Pages:                pages,
	}, nil
}

func (r *InMemoryLearningRepository) GetPageProgress(ctx context.Context, userID, bookID uuid.UUID, pageNumber int) (*models.PageProgress, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	progressKey := userID.String() + ":" + bookID.String() + ":" + string(rune(pageNumber))
	progress, exists := r.progress[progressKey]

	if !exists {
		return &models.PageProgress{
			UserID:      userID.String(),
			BookID:      bookID.String(),
			PageNumber:  pageNumber,
			IsCompleted: false,
		}, nil
	}

	return progress, nil
}

func (r *InMemoryLearningRepository) UpdatePageProgress(ctx context.Context, progress *models.PageProgress) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	progressKey := progress.UserID + ":" + progress.BookID + ":" + string(rune(progress.PageNumber))
	progress.UpdatedAt = time.Now()
	r.progress[progressKey] = progress

	return nil
}
```

### handler/learning.go

```go
package handler

import (
	"net/http"
	"strconv"

	"github.com/clearclown/HaiLanGo/backend/internal/models"
	"github.com/clearclown/HaiLanGo/backend/internal/repository"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type LearningHandler struct {
	repo repository.LearningRepository
}

func NewLearningHandler(repo repository.LearningRepository) *LearningHandler {
	return &LearningHandler{repo: repo}
}

// GetPageLearning godoc
// @Summary Get page learning data
// @Description Get learning data for a specific page
// @Tags learning
// @Accept json
// @Produce json
// @Param bookId path string true "Book ID"
// @Param pageNumber path int true "Page Number"
// @Security BearerAuth
// @Success 200 {object} models.PageLearning
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /api/v1/learning/books/{bookId}/pages/{pageNumber} [get]
func (h *LearningHandler) GetPageLearning(c *gin.Context) {
	userIDStr, _ := c.Get("user_id")
	userID, _ := uuid.Parse(userIDStr.(string))

	bookIDStr := c.Param("bookId")
	bookID, err := uuid.Parse(bookIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid book ID"})
		return
	}

	pageNumberStr := c.Param("pageNumber")
	pageNumber, err := strconv.Atoi(pageNumberStr)
	if err != nil || pageNumber < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid page number"})
		return
	}

	data, err := h.repo.GetPageLearning(c.Request.Context(), userID, bookID, pageNumber)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Page not found"})
		return
	}

	c.JSON(http.StatusOK, data)
}

// CompletePage godoc
// @Summary Complete page learning
// @Description Mark a page as completed
// @Tags learning
// @Accept json
// @Produce json
// @Param bookId path string true "Book ID"
// @Param pageNumber path int true "Page Number"
// @Param request body models.CompletePageRequest true "Complete Page Request"
// @Security BearerAuth
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Router /api/v1/learning/books/{bookId}/pages/{pageNumber}/complete [post]
func (h *LearningHandler) CompletePage(c *gin.Context) {
	userIDStr, _ := c.Get("user_id")
	userID, _ := uuid.Parse(userIDStr.(string))

	bookIDStr := c.Param("bookId")
	bookID, _ := uuid.Parse(bookIDStr)

	pageNumberStr := c.Param("pageNumber")
	pageNumber, _ := strconv.Atoi(pageNumberStr)

	var req models.CompletePageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	progress, err := h.repo.CompletePage(c.Request.Context(), userID, bookID, pageNumber, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to complete page"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":   "Page marked as completed",
		"progress":  progress,
		"next_page": pageNumber + 1,
	})
}

// RecordSession godoc
// @Summary Record learning session
// @Description Record start/end of learning session
// @Tags learning
// @Accept json
// @Produce json
// @Param bookId path string true "Book ID"
// @Param pageNumber path int true "Page Number"
// @Param request body models.SessionRequest true "Session Request"
// @Security BearerAuth
// @Success 200 {object} models.SessionResponse
// @Failure 400 {object} map[string]string
// @Router /api/v1/learning/books/{bookId}/pages/{pageNumber}/session [post]
func (h *LearningHandler) RecordSession(c *gin.Context) {
	userIDStr, _ := c.Get("user_id")
	userID, _ := uuid.Parse(userIDStr.(string))

	bookIDStr := c.Param("bookId")
	bookID, _ := uuid.Parse(bookIDStr)

	pageNumberStr := c.Param("pageNumber")
	pageNumber, _ := strconv.Atoi(pageNumberStr)

	var req models.SessionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := h.repo.RecordSession(c.Request.Context(), userID, bookID, pageNumber, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to record session"})
		return
	}

	c.JSON(http.StatusOK, response)
}

// GetBookProgress godoc
// @Summary Get book progress
// @Description Get overall progress for a book
// @Tags learning
// @Accept json
// @Produce json
// @Param bookId path string true "Book ID"
// @Security BearerAuth
// @Success 200 {object} models.BookProgressSummary
// @Failure 400 {object} map[string]string
// @Router /api/v1/learning/books/{bookId}/progress [get]
func (h *LearningHandler) GetBookProgress(c *gin.Context) {
	userIDStr, _ := c.Get("user_id")
	userID, _ := uuid.Parse(userIDStr.(string))

	bookIDStr := c.Param("bookId")
	bookID, _ := uuid.Parse(bookIDStr)

	summary, err := h.repo.GetBookProgress(c.Request.Context(), userID, bookID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get book progress"})
		return
	}

	c.JSON(http.StatusOK, summary)
}

// RegisterRoutes registers learning routes
func (h *LearningHandler) RegisterRoutes(rg *gin.RouterGroup) {
	learning := rg.Group("/learning")
	{
		learning.GET("/books/:bookId/pages/:pageNumber", h.GetPageLearning)
		learning.POST("/books/:bookId/pages/:pageNumber/complete", h.CompletePage)
		learning.POST("/books/:bookId/pages/:pageNumber/session", h.RecordSession)
		learning.GET("/books/:bookId/progress", h.GetBookProgress)
	}
}
```

---

## ✅ 完了条件

- [ ] すべてのエンドポイントが実装されている
- [ ] InMemoryリポジトリにサンプルデータが含まれている
- [ ] ルーターに登録されている
- [ ] テストがすべてパスする

---

**期限**: 72時間以内
**次のタスク**: CRITICAL_06 (OCR API)
