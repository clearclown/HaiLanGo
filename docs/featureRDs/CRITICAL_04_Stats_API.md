# CRITICAL_04: Stats API実装

**優先度**: P0（最高優先度）
**担当者**: 未割当
**見積時間**: 4-6時間
**ブロッカー**: フロントエンド（ホーム画面・統計ダッシュボード）がこのAPIを待っている

---

## ⚠️ PM指示

**現状**: フロントエンドはStats APIを呼び出しているが、バックエンドが404を返している。
**期限**: 48時間以内に実装完了すること。
**言い訳は不要**: 技術的な難易度は高くない。既存のReview APIを参考にすれば実装可能。

---

## 📋 実装要件

### エンドポイント仕様

#### 1. GET /api/v1/stats/dashboard
**説明**: ダッシュボード用の統計サマリーを取得

**Request**:
```http
GET /api/v1/stats/dashboard
Authorization: Bearer <JWT_TOKEN>
```

**Response** (200 OK):
```json
{
  "learning_time_today": 45,
  "learning_time_this_week": 180,
  "total_learning_time": 3420,
  "current_streak": 7,
  "longest_streak": 15,
  "completed_pages": 45,
  "total_pages": 150,
  "mastered_words": 230,
  "mastered_phrases": 45,
  "completed_books": 1,
  "total_books": 3,
  "average_pronunciation_score": 85.5
}
```

#### 2. GET /api/v1/stats/learning-time
**説明**: 学習時間の推移データを取得

**Request**:
```http
GET /api/v1/stats/learning-time?period=week
Authorization: Bearer <JWT_TOKEN>

Query Parameters:
- period: day | week | month | year
```

**Response** (200 OK):
```json
{
  "period": "week",
  "data": [
    {"date": "2025-11-08", "minutes": 30},
    {"date": "2025-11-09", "minutes": 25},
    {"date": "2025-11-10", "minutes": 40},
    {"date": "2025-11-11", "minutes": 20},
    {"date": "2025-11-12", "minutes": 35},
    {"date": "2025-11-13", "minutes": 30},
    {"date": "2025-11-14", "minutes": 0}
  ],
  "total_minutes": 180,
  "average_minutes": 25.7
}
```

#### 3. GET /api/v1/stats/progress
**説明**: 学習進捗の推移データを取得

**Request**:
```http
GET /api/v1/stats/progress?period=month
Authorization: Bearer <JWT_TOKEN>

Query Parameters:
- period: week | month | year
```

**Response** (200 OK):
```json
{
  "period": "month",
  "words": [
    {"date": "2025-10-14", "count": 50},
    {"date": "2025-10-21", "count": 85},
    {"date": "2025-10-28", "count": 120},
    {"date": "2025-11-04", "count": 165},
    {"date": "2025-11-11", "count": 230}
  ],
  "phrases": [
    {"date": "2025-10-14", "count": 10},
    {"date": "2025-10-21", "count": 18},
    {"date": "2025-10-28", "count": 25},
    {"date": "2025-11-04", "count": 35},
    {"date": "2025-11-11", "count": 45}
  ],
  "pages": [
    {"date": "2025-10-14", "count": 5},
    {"date": "2025-10-21", "count": 12},
    {"date": "2025-10-28", "count": 20},
    {"date": "2025-11-04", "count": 32},
    {"date": "2025-11-11", "count": 45}
  ]
}
```

#### 4. GET /api/v1/stats/weak-points
**説明**: 弱点分析データを取得（苦手な単語・フレーズ）

**Request**:
```http
GET /api/v1/stats/weak-points?limit=10
Authorization: Bearer <JWT_TOKEN>
```

**Response** (200 OK):
```json
{
  "weak_words": [
    {
      "word": "Здравствуйте",
      "language": "ru",
      "attempts": 15,
      "average_score": 45,
      "last_attempt": "2025-11-14T10:30:00Z"
    }
  ],
  "weak_phrases": [
    {
      "phrase": "Как дела?",
      "language": "ru",
      "attempts": 8,
      "average_score": 52,
      "last_attempt": "2025-11-13T15:20:00Z"
    }
  ]
}
```

---

## 🗃️ データベーススキーマ

### learning_sessions テーブル
```sql
CREATE TABLE learning_sessions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    book_id UUID REFERENCES books(id) ON DELETE SET NULL,
    page_number INT,
    started_at TIMESTAMP NOT NULL DEFAULT NOW(),
    ended_at TIMESTAMP,
    duration_minutes INT,
    activity_type VARCHAR(50) NOT NULL, -- 'reading', 'listening', 'speaking', 'review'
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_learning_sessions_user_id ON learning_sessions(user_id);
CREATE INDEX idx_learning_sessions_started_at ON learning_sessions(started_at);
```

### user_progress テーブル
```sql
CREATE TABLE user_progress (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    date DATE NOT NULL,
    completed_pages INT DEFAULT 0,
    mastered_words INT DEFAULT 0,
    mastered_phrases INT DEFAULT 0,
    learning_minutes INT DEFAULT 0,
    pronunciation_attempts INT DEFAULT 0,
    pronunciation_total_score INT DEFAULT 0,
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    UNIQUE(user_id, date)
);

CREATE INDEX idx_user_progress_user_date ON user_progress(user_id, date);
```

### learning_streaks テーブル
```sql
CREATE TABLE learning_streaks (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL UNIQUE REFERENCES users(id) ON DELETE CASCADE,
    current_streak INT DEFAULT 0,
    longest_streak INT DEFAULT 0,
    last_activity_date DATE,
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_learning_streaks_user_id ON learning_streaks(user_id);
```

---

## 🏗️ 実装ステップ

### Step 1: データモデル作成 (30分)

**ファイル**: `backend/internal/models/stats.go`

```go
package models

import "time"

// DashboardStats はダッシュボード統計
type DashboardStats struct {
	LearningTimeToday          int     `json:"learning_time_today"`
	LearningTimeThisWeek       int     `json:"learning_time_this_week"`
	TotalLearningTime          int     `json:"total_learning_time"`
	CurrentStreak              int     `json:"current_streak"`
	LongestStreak              int     `json:"longest_streak"`
	CompletedPages             int     `json:"completed_pages"`
	TotalPages                 int     `json:"total_pages"`
	MasteredWords              int     `json:"mastered_words"`
	MasteredPhrases            int     `json:"mastered_phrases"`
	CompletedBooks             int     `json:"completed_books"`
	TotalBooks                 int     `json:"total_books"`
	AveragePronunciationScore  float64 `json:"average_pronunciation_score"`
}

// LearningTimeData は学習時間データ
type LearningTimeData struct {
	Period         string              `json:"period"`
	Data           []DailyLearningTime `json:"data"`
	TotalMinutes   int                 `json:"total_minutes"`
	AverageMinutes float64             `json:"average_minutes"`
}

type DailyLearningTime struct {
	Date    string `json:"date"`
	Minutes int    `json:"minutes"`
}

// ProgressData は進捗データ
type ProgressData struct {
	Period  string            `json:"period"`
	Words   []TimeSeriesData  `json:"words"`
	Phrases []TimeSeriesData  `json:"phrases"`
	Pages   []TimeSeriesData  `json:"pages"`
}

type TimeSeriesData struct {
	Date  string `json:"date"`
	Count int    `json:"count"`
}

// WeakPointsData は弱点分析データ
type WeakPointsData struct {
	WeakWords   []WeakItem `json:"weak_words"`
	WeakPhrases []WeakItem `json:"weak_phrases"`
}

type WeakItem struct {
	Word         string    `json:"word,omitempty"`
	Phrase       string    `json:"phrase,omitempty"`
	Language     string    `json:"language"`
	Attempts     int       `json:"attempts"`
	AverageScore float64   `json:"average_score"`
	LastAttempt  time.Time `json:"last_attempt"`
}

// LearningSession は学習セッション
type LearningSession struct {
	ID              string    `json:"id"`
	UserID          string    `json:"user_id"`
	BookID          *string   `json:"book_id"`
	PageNumber      *int      `json:"page_number"`
	StartedAt       time.Time `json:"started_at"`
	EndedAt         *time.Time `json:"ended_at"`
	DurationMinutes *int      `json:"duration_minutes"`
	ActivityType    string    `json:"activity_type"` // reading, listening, speaking, review
	CreatedAt       time.Time `json:"created_at"`
}

// UserProgress はユーザー進捗
type UserProgress struct {
	ID                        string    `json:"id"`
	UserID                    string    `json:"user_id"`
	Date                      time.Time `json:"date"`
	CompletedPages            int       `json:"completed_pages"`
	MasteredWords             int       `json:"mastered_words"`
	MasteredPhrases           int       `json:"mastered_phrases"`
	LearningMinutes           int       `json:"learning_minutes"`
	PronunciationAttempts     int       `json:"pronunciation_attempts"`
	PronunciationTotalScore   int       `json:"pronunciation_total_score"`
	UpdatedAt                 time.Time `json:"updated_at"`
}

// LearningStreak はストリーク情報
type LearningStreak struct {
	ID               string    `json:"id"`
	UserID           string    `json:"user_id"`
	CurrentStreak    int       `json:"current_streak"`
	LongestStreak    int       `json:"longest_streak"`
	LastActivityDate time.Time `json:"last_activity_date"`
	UpdatedAt        time.Time `json:"updated_at"`
}
```

---

### Step 2: リポジトリ実装 (1-1.5時間)

**ファイル**: `backend/internal/repository/stats.go`

```go
package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/clearclown/HaiLanGo/backend/internal/models"
	"github.com/google/uuid"
)

type StatsRepository interface {
	GetDashboardStats(ctx context.Context, userID uuid.UUID) (*models.DashboardStats, error)
	GetLearningTimeData(ctx context.Context, userID uuid.UUID, period string) (*models.LearningTimeData, error)
	GetProgressData(ctx context.Context, userID uuid.UUID, period string) (*models.ProgressData, error)
	GetWeakPoints(ctx context.Context, userID uuid.UUID, limit int) (*models.WeakPointsData, error)
	RecordLearningSession(ctx context.Context, session *models.LearningSession) error
	UpdateUserProgress(ctx context.Context, progress *models.UserProgress) error
	UpdateStreak(ctx context.Context, userID uuid.UUID, activityDate time.Time) error
}

type StatsRepositoryPostgres struct {
	db *sql.DB
}

func NewStatsRepositoryPostgres(db *sql.DB) *StatsRepositoryPostgres {
	return &StatsRepositoryPostgres{db: db}
}

func (r *StatsRepositoryPostgres) GetDashboardStats(ctx context.Context, userID uuid.UUID) (*models.DashboardStats, error) {
	// TODO: 実装
	// 複数のテーブルからデータを集計
	return nil, nil
}

// ... 他のメソッド実装
```

**ファイル**: `backend/internal/repository/stats_inmemory.go`

```go
package repository

import (
	"context"
	"sync"
	"time"

	"github.com/clearclown/HaiLanGo/backend/internal/models"
	"github.com/google/uuid"
)

type InMemoryStatsRepository struct {
	sessions  map[string]*models.LearningSession
	progress  map[string]map[string]*models.UserProgress // userID -> date -> progress
	streaks   map[string]*models.LearningStreak
	mu        sync.RWMutex
}

func NewInMemoryStatsRepository() *InMemoryStatsRepository {
	repo := &InMemoryStatsRepository{
		sessions: make(map[string]*models.LearningSession),
		progress: make(map[string]map[string]*models.UserProgress),
		streaks:  make(map[string]*models.LearningStreak),
	}

	// サンプルデータ初期化
	repo.initSampleData()

	return repo
}

func (r *InMemoryStatsRepository) initSampleData() {
	// テストユーザーのストリークデータ
	testUserID := "550e8400-e29b-41d4-a716-446655440000"
	r.streaks[testUserID] = &models.LearningStreak{
		ID:               uuid.New().String(),
		UserID:           testUserID,
		CurrentStreak:    7,
		LongestStreak:    15,
		LastActivityDate: time.Now(),
		UpdatedAt:        time.Now(),
	}

	// 過去7日間の進捗データ
	r.progress[testUserID] = make(map[string]*models.UserProgress)
	for i := 0; i < 7; i++ {
		date := time.Now().AddDate(0, 0, -i)
		dateStr := date.Format("2006-01-02")
		r.progress[testUserID][dateStr] = &models.UserProgress{
			ID:                      uuid.New().String(),
			UserID:                  testUserID,
			Date:                    date,
			CompletedPages:          i * 2,
			MasteredWords:           i * 10,
			MasteredPhrases:         i * 2,
			LearningMinutes:         25 + i*5,
			PronunciationAttempts:   i * 3,
			PronunciationTotalScore: i * 250,
			UpdatedAt:               time.Now(),
		}
	}
}

func (r *InMemoryStatsRepository) GetDashboardStats(ctx context.Context, userID uuid.UUID) (*models.DashboardStats, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	streak, _ := r.streaks[userID.String()]

	// 進捗データを集計
	totalPages := 0
	totalWords := 0
	totalPhrases := 0
	todayMinutes := 0
	weekMinutes := 0
	totalMinutes := 0

	userProgress, exists := r.progress[userID.String()]
	if exists {
		for dateStr, prog := range userProgress {
			date, _ := time.Parse("2006-01-02", dateStr)

			totalPages += prog.CompletedPages
			totalWords += prog.MasteredWords
			totalPhrases += prog.MasteredPhrases
			totalMinutes += prog.LearningMinutes

			if isToday(date) {
				todayMinutes = prog.LearningMinutes
			}
			if isThisWeek(date) {
				weekMinutes += prog.LearningMinutes
			}
		}
	}

	stats := &models.DashboardStats{
		LearningTimeToday:         todayMinutes,
		LearningTimeThisWeek:      weekMinutes,
		TotalLearningTime:         totalMinutes,
		CurrentStreak:             0,
		LongestStreak:             0,
		CompletedPages:            totalPages,
		TotalPages:                150, // 仮の値
		MasteredWords:             totalWords,
		MasteredPhrases:           totalPhrases,
		CompletedBooks:            1,
		TotalBooks:                3,
		AveragePronunciationScore: 85.5,
	}

	if streak != nil {
		stats.CurrentStreak = streak.CurrentStreak
		stats.LongestStreak = streak.LongestStreak
	}

	return stats, nil
}

func (r *InMemoryStatsRepository) GetLearningTimeData(ctx context.Context, userID uuid.UUID, period string) (*models.LearningTimeData, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	data := &models.LearningTimeData{
		Period: period,
		Data:   []models.DailyLearningTime{},
	}

	userProgress, exists := r.progress[userID.String()]
	if !exists {
		return data, nil
	}

	days := getDaysForPeriod(period)
	totalMinutes := 0

	for i := days - 1; i >= 0; i-- {
		date := time.Now().AddDate(0, 0, -i)
		dateStr := date.Format("2006-01-02")

		minutes := 0
		if prog, ok := userProgress[dateStr]; ok {
			minutes = prog.LearningMinutes
			totalMinutes += minutes
		}

		data.Data = append(data.Data, models.DailyLearningTime{
			Date:    dateStr,
			Minutes: minutes,
		})
	}

	data.TotalMinutes = totalMinutes
	if len(data.Data) > 0 {
		data.AverageMinutes = float64(totalMinutes) / float64(len(data.Data))
	}

	return data, nil
}

func (r *InMemoryStatsRepository) GetProgressData(ctx context.Context, userID uuid.UUID, period string) (*models.ProgressData, error) {
	// TODO: 実装
	return &models.ProgressData{
		Period:  period,
		Words:   []models.TimeSeriesData{},
		Phrases: []models.TimeSeriesData{},
		Pages:   []models.TimeSeriesData{},
	}, nil
}

func (r *InMemoryStatsRepository) GetWeakPoints(ctx context.Context, userID uuid.UUID, limit int) (*models.WeakPointsData, error) {
	// TODO: 実装（STT/発音データが必要）
	return &models.WeakPointsData{
		WeakWords:   []models.WeakItem{},
		WeakPhrases: []models.WeakItem{},
	}, nil
}

func (r *InMemoryStatsRepository) RecordLearningSession(ctx context.Context, session *models.LearningSession) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.sessions[session.ID] = session
	return nil
}

func (r *InMemoryStatsRepository) UpdateUserProgress(ctx context.Context, progress *models.UserProgress) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.progress[progress.UserID]; !exists {
		r.progress[progress.UserID] = make(map[string]*models.UserProgress)
	}

	dateStr := progress.Date.Format("2006-01-02")
	r.progress[progress.UserID][dateStr] = progress

	return nil
}

func (r *InMemoryStatsRepository) UpdateStreak(ctx context.Context, userID uuid.UUID, activityDate time.Time) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	userIDStr := userID.String()
	streak, exists := r.streaks[userIDStr]

	if !exists {
		streak = &models.LearningStreak{
			ID:               uuid.New().String(),
			UserID:           userIDStr,
			CurrentStreak:    1,
			LongestStreak:    1,
			LastActivityDate: activityDate,
			UpdatedAt:        time.Now(),
		}
		r.streaks[userIDStr] = streak
		return nil
	}

	// ストリーク計算ロジック
	daysDiff := int(activityDate.Sub(streak.LastActivityDate).Hours() / 24)

	if daysDiff == 1 {
		// 連続
		streak.CurrentStreak++
		if streak.CurrentStreak > streak.LongestStreak {
			streak.LongestStreak = streak.CurrentStreak
		}
	} else if daysDiff > 1 {
		// 途切れた
		streak.CurrentStreak = 1
	}
	// daysDiff == 0 なら同じ日なので何もしない

	streak.LastActivityDate = activityDate
	streak.UpdatedAt = time.Now()

	return nil
}

// ヘルパー関数
func isToday(date time.Time) bool {
	now := time.Now()
	return date.Year() == now.Year() && date.YearDay() == now.YearDay()
}

func isThisWeek(date time.Time) bool {
	now := time.Now()
	_, week := now.ISOWeek()
	_, dateWeek := date.ISOWeek()
	return week == dateWeek && now.Year() == date.Year()
}

func getDaysForPeriod(period string) int {
	switch period {
	case "day":
		return 1
	case "week":
		return 7
	case "month":
		return 30
	case "year":
		return 365
	default:
		return 7
	}
}
```

---

### Step 3: ハンドラー実装 (1時間)

**ファイル**: `backend/internal/api/handler/stats.go`

```go
package handler

import (
	"net/http"
	"strconv"

	"github.com/clearclown/HaiLanGo/backend/internal/repository"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type StatsHandler struct {
	repo repository.StatsRepository
}

func NewStatsHandler(repo repository.StatsRepository) *StatsHandler {
	return &StatsHandler{repo: repo}
}

// GetDashboard godoc
// @Summary Get dashboard statistics
// @Description Get overall learning statistics for dashboard
// @Tags stats
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} models.DashboardStats
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/stats/dashboard [get]
func (h *StatsHandler) GetDashboard(c *gin.Context) {
	userIDStr, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	stats, err := h.repo.GetDashboardStats(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get dashboard stats"})
		return
	}

	c.JSON(http.StatusOK, stats)
}

// GetLearningTime godoc
// @Summary Get learning time data
// @Description Get learning time data for specified period
// @Tags stats
// @Accept json
// @Produce json
// @Param period query string false "Period (day|week|month|year)" default(week)
// @Security BearerAuth
// @Success 200 {object} models.LearningTimeData
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/stats/learning-time [get]
func (h *StatsHandler) GetLearningTime(c *gin.Context) {
	userIDStr, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	period := c.DefaultQuery("period", "week")

	data, err := h.repo.GetLearningTimeData(c.Request.Context(), userID, period)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get learning time data"})
		return
	}

	c.JSON(http.StatusOK, data)
}

// GetProgress godoc
// @Summary Get progress data
// @Description Get learning progress data for specified period
// @Tags stats
// @Accept json
// @Produce json
// @Param period query string false "Period (week|month|year)" default(month)
// @Security BearerAuth
// @Success 200 {object} models.ProgressData
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/stats/progress [get]
func (h *StatsHandler) GetProgress(c *gin.Context) {
	userIDStr, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	period := c.DefaultQuery("period", "month")

	data, err := h.repo.GetProgressData(c.Request.Context(), userID, period)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get progress data"})
		return
	}

	c.JSON(http.StatusOK, data)
}

// GetWeakPoints godoc
// @Summary Get weak points analysis
// @Description Get weak points (words/phrases with low scores)
// @Tags stats
// @Accept json
// @Produce json
// @Param limit query int false "Limit" default(10)
// @Security BearerAuth
// @Success 200 {object} models.WeakPointsData
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/stats/weak-points [get]
func (h *StatsHandler) GetWeakPoints(c *gin.Context) {
	userIDStr, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	limit := 10
	if limitStr := c.Query("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	data, err := h.repo.GetWeakPoints(c.Request.Context(), userID, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get weak points"})
		return
	}

	c.JSON(http.StatusOK, data)
}

// RegisterRoutes registers stats routes
func (h *StatsHandler) RegisterRoutes(rg *gin.RouterGroup) {
	stats := rg.Group("/stats")
	{
		stats.GET("/dashboard", h.GetDashboard)
		stats.GET("/learning-time", h.GetLearningTime)
		stats.GET("/progress", h.GetProgress)
		stats.GET("/weak-points", h.GetWeakPoints)
	}
}
```

---

### Step 4: router.goに統合 (15分)

**ファイル**: `backend/internal/api/router/router.go`

```go
// ========================================
// リポジトリの初期化
// ========================================
bookRepo := repository.NewBookRepositoryPostgres(db)
reviewRepo := repository.NewInMemoryReviewRepository()
statsRepo := repository.NewInMemoryStatsRepository() // 追加

// ========================================
// ハンドラーの初期化
// ========================================
uploadHandler := handler.NewUploadHandler(uploadService)
booksHandler := handler.NewBooksHandler(bookRepo)
reviewHandler := handler.NewReviewHandler(reviewRepo)
statsHandler := handler.NewStatsHandler(statsRepo) // 追加

// ========================================
// ルート登録（認証必須グループ内）
// ========================================
authenticated := v1.Group("")
authenticated.Use(middleware.AuthRequired())
{
	booksHandler.RegisterRoutes(authenticated)
	uploadHandler.RegisterRoutes(authenticated)
	reviewHandler.RegisterRoutes(authenticated)
	statsHandler.RegisterRoutes(authenticated) // 追加

	// 他のハンドラー...
}
```

---

### Step 5: テスト作成 (1-1.5時間)

**ファイル**: `backend/internal/api/handler/stats_test.go`

```go
package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/clearclown/HaiLanGo/backend/internal/models"
	"github.com/clearclown/HaiLanGo/backend/internal/repository"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setupStatsTestRouter() (*gin.Engine, *repository.InMemoryStatsRepository) {
	gin.SetMode(gin.TestMode)

	repo := repository.NewInMemoryStatsRepository()
	handler := NewStatsHandler(repo)

	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("user_id", "550e8400-e29b-41d4-a716-446655440000")
		c.Next()
	})

	handler.RegisterRoutes(r.Group("/api/v1"))

	return r, repo
}

func TestGetDashboard(t *testing.T) {
	router, _ := setupStatsTestRouter()

	req, _ := http.NewRequest(http.MethodGet, "/api/v1/stats/dashboard", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var stats models.DashboardStats
	err := json.Unmarshal(w.Body.Bytes(), &stats)
	assert.NoError(t, err)

	assert.GreaterOrEqual(t, stats.CurrentStreak, 0)
	assert.GreaterOrEqual(t, stats.LongestStreak, 0)
	assert.GreaterOrEqual(t, stats.MasteredWords, 0)
}

func TestGetLearningTime(t *testing.T) {
	router, _ := setupStatsTestRouter()

	req, _ := http.NewRequest(http.MethodGet, "/api/v1/stats/learning-time?period=week", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var data models.LearningTimeData
	err := json.Unmarshal(w.Body.Bytes(), &data)
	assert.NoError(t, err)

	assert.Equal(t, "week", data.Period)
	assert.Len(t, data.Data, 7)
}

func TestGetProgress(t *testing.T) {
	router, _ := setupStatsTestRouter()

	req, _ := http.NewRequest(http.MethodGet, "/api/v1/stats/progress?period=month", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var data models.ProgressData
	err := json.Unmarshal(w.Body.Bytes(), &data)
	assert.NoError(t, err)

	assert.Equal(t, "month", data.Period)
}

func TestGetWeakPoints(t *testing.T) {
	router, _ := setupStatsTestRouter()

	req, _ := http.NewRequest(http.MethodGet, "/api/v1/stats/weak-points?limit=5", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var data models.WeakPointsData
	err := json.Unmarshal(w.Body.Bytes(), &data)
	assert.NoError(t, err)
}
```

---

## ✅ 完了条件

- [ ] すべてのエンドポイントが実装され、ルーターに登録されている
- [ ] InMemoryリポジトリにサンプルデータが含まれている
- [ ] すべてのテストがパスする（`go test ./internal/api/handler -run Stats`）
- [ ] サーバー起動時にルートが登録される
- [ ] フロントエンドからのリクエストが200/401を返す（404ではない）

---

## 📝 実装チェックリスト

### コード実装
- [ ] `internal/models/stats.go` 作成
- [ ] `internal/repository/stats.go` 作成
- [ ] `internal/repository/stats_inmemory.go` 作成
- [ ] `internal/api/handler/stats.go` 作成
- [ ] `internal/api/router/router.go` 修正（StatsHandler登録）

### テスト
- [ ] `internal/api/handler/stats_test.go` 作成
- [ ] すべてのテストがパス

### 動作確認
- [ ] サーバー起動
- [ ] `GET /api/v1/stats/dashboard` → 401 Unauthorized
- [ ] `GET /api/v1/stats/learning-time?period=week` → 401 Unauthorized
- [ ] `GET /api/v1/stats/progress?period=month` → 401 Unauthorized
- [ ] `GET /api/v1/stats/weak-points` → 401 Unauthorized

---

## 🚨 注意事項

1. **PostgreSQL実装は後回し**: まずInMemoryで動作させること
2. **サンプルデータ必須**: フロントエンドが表示できるように最低限のデータを用意
3. **エラーハンドリング**: すべてのエンドポイントで適切なHTTPステータスコードを返す
4. **パフォーマンス**: 後でPostgreSQL実装時にインデックスを追加

---

## 🎯 成果物

実装完了後、以下を提出：

1. **コミット**: `feat(backend): Stats API実装`
2. **テスト結果**: `go test -v ./internal/api/handler -run Stats` の出力
3. **動作確認**: curl/Postmanでの各エンドポイントのレスポンス

---

**期限**: 48時間以内
**次のタスク**: CRITICAL_05 (Learning API)
