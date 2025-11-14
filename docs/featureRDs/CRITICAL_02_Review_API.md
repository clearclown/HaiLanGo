# 🚨 CRITICAL - Review API Implementation (SRS - Spaced Repetition System)

**優先度**: P0 - CRITICAL
**担当者**: Backend Engineer
**見積もり**: 6-8時間
**期限**: 即座
**ブロッカー**: フロントエンドが実装済みで現在失敗中

## 現状の問題

❌ **Review APIが未実装のため、フロントエンドの復習機能が完全に動作していない**
- フロントエンドは `/api/v1/review/stats` を呼び出すが404エラー
- 復習ページがエラー状態を表示
- E2Eテストが60%失敗（review.spec.ts）

## 実装要件

### 1. APIエンドポイント

#### 1.1 復習統計取得
```
GET /api/v1/review/stats
Headers: Authorization: Bearer {token}
Response 200:
{
  "urgent_count": 3,
  "recommended_count": 5,
  "optional_count": 4,
  "total_completed_today": 2,
  "weekly_completion_rate": 65.5
}
```

#### 1.2 復習アイテム取得
```
GET /api/v1/review/items?priority={urgent|recommended|optional}
Headers: Authorization: Bearer {token}
Query Params:
  - priority (optional): urgent, recommended, optional
Response 200:
{
  "items": [
    {
      "id": "uuid",
      "type": "word",
      "text": "Здравствуйте",
      "translation": "こんにちは",
      "language": "ru",
      "mastery_level": 45,
      "last_reviewed": "2025-11-13T10:00:00Z",
      "next_review": "2025-11-14T10:00:00Z",
      "priority": "urgent"
    }
  ]
}
```

#### 1.3 復習結果送信
```
POST /api/v1/review/submit
Headers: Authorization: Bearer {token}
Content-Type: application/json
Body:
{
  "item_id": "uuid",
  "score": 100,
  "completed_at": "2025-11-14T10:30:00Z"
}
Response 200:
{
  "success": true,
  "next_review": "2025-11-16T10:00:00Z"
}
```

### 2. データベーススキーマ

```sql
-- 復習アイテムテーブル
CREATE TABLE IF NOT EXISTS review_items (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    book_id UUID NOT NULL REFERENCES books(id) ON DELETE CASCADE,
    page_number INTEGER NOT NULL,
    type VARCHAR(50) NOT NULL CHECK (type IN ('word', 'phrase')),
    text TEXT NOT NULL,
    translation TEXT NOT NULL,
    language VARCHAR(10) NOT NULL,
    mastery_level INTEGER DEFAULT 0 CHECK (mastery_level >= 0 AND mastery_level <= 100),
    interval_days INTEGER DEFAULT 1,
    ease_factor DECIMAL(3,2) DEFAULT 2.5,
    last_reviewed TIMESTAMP,
    next_review TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    review_count INTEGER DEFAULT 0,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_user_next_review (user_id, next_review),
    INDEX idx_user_mastery (user_id, mastery_level),
    INDEX idx_book_id (book_id)
);

-- 復習履歴テーブル
CREATE TABLE IF NOT EXISTS review_history (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    review_item_id UUID NOT NULL REFERENCES review_items(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    score INTEGER NOT NULL CHECK (score >= 0 AND score <= 100),
    reviewed_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_user_reviewed (user_id, reviewed_at),
    INDEX idx_item_id (review_item_id)
);
```

### 3. SRSアルゴリズム実装

```go
package service

import (
	"math"
	"time"
)

// SM2Algorithm (SuperMemo 2) アルゴリズム実装
type SM2Algorithm struct{}

func NewSM2Algorithm() *SM2Algorithm {
	return &SM2Algorithm{}
}

// CalculateNextReview は次の復習日時を計算する
// score: 0-100 (30=思い出せない, 70=少し時間がかかった, 100=完璧)
func (s *SM2Algorithm) CalculateNextReview(
	currentEaseFactor float64,
	currentInterval int,
	score int,
) (nextInterval int, nextEaseFactor float64, nextReview time.Time) {

	// スコアを0-5の品質スケールに変換
	quality := s.scoreToQuality(score)

	// 新しい容易度係数を計算
	nextEaseFactor = currentEaseFactor + (0.1 - (5-quality)*(0.08+(5-quality)*0.02))
	if nextEaseFactor < 1.3 {
		nextEaseFactor = 1.3
	}

	// 次の間隔を計算
	if quality < 3 {
		// 失敗：最初からやり直し
		nextInterval = 1
	} else {
		if currentInterval == 0 {
			nextInterval = 1
		} else if currentInterval == 1 {
			nextInterval = 6
		} else {
			nextInterval = int(math.Round(float64(currentInterval) * nextEaseFactor))
		}
	}

	// 次の復習日時
	nextReview = time.Now().Add(time.Duration(nextInterval) * 24 * time.Hour)

	return nextInterval, nextEaseFactor, nextReview
}

func (s *SM2Algorithm) scoreToQuality(score int) int {
	switch {
	case score >= 90:
		return 5 // 完璧
	case score >= 70:
		return 4 // 正解だが努力が必要
	case score >= 50:
		return 3 // かろうじて正解
	case score >= 30:
		return 2 // 不正解だが覚えていた
	default:
		return 0 // 完全に忘れた
	}
}

// CalculatePriority は復習の優先度を計算する
func (s *SM2Algorithm) CalculatePriority(nextReview time.Time) string {
	now := time.Now()
	hoursUntil := nextReview.Sub(now).Hours()

	if hoursUntil <= 0 {
		return "urgent" // 期限切れ
	} else if hoursUntil <= 24 {
		return "urgent" // 今日中
	} else if hoursUntil <= 48 {
		return "recommended" // 明日まで
	} else {
		return "optional" // 余裕あり
	}
}
```

### 4. Handler Implementation

```go
package handler

import (
	"net/http"
	"time"
	"github.com/gin-gonic/gin"
	"github.com/clearclown/HaiLanGo/backend/internal/models"
	"github.com/clearclown/HaiLanGo/backend/internal/repository"
	"github.com/clearclown/HaiLanGo/backend/internal/service"
)

type ReviewHandler struct {
	repo      repository.ReviewRepository
	srsAlgo   *service.SM2Algorithm
}

func NewReviewHandler(repo repository.ReviewRepository) *ReviewHandler {
	return &ReviewHandler{
		repo:    repo,
		srsAlgo: service.NewSM2Algorithm(),
	}
}

// GetStats godoc
// @Summary Get review statistics
// @Tags review
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} models.ReviewStats
// @Router /api/v1/review/stats [get]
func (h *ReviewHandler) GetStats(c *gin.Context) {
	userID := c.GetString("user_id")

	// 今日の開始時刻
	todayStart := time.Now().Truncate(24 * time.Hour)

	// すべての復習アイテムを取得
	items, err := h.repo.FindByUserID(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch review items"})
		return
	}

	stats := models.ReviewStats{
		UrgentCount:      0,
		RecommendedCount: 0,
		OptionalCount:    0,
	}

	// 優先度別にカウント
	for _, item := range items {
		priority := h.srsAlgo.CalculatePriority(item.NextReview)
		switch priority {
		case "urgent":
			stats.UrgentCount++
		case "recommended":
			stats.RecommendedCount++
		case "optional":
			stats.OptionalCount++
		}
	}

	// 今日完了した復習数を取得
	stats.TotalCompletedToday, err = h.repo.CountCompletedToday(c.Request.Context(), userID, todayStart)
	if err != nil {
		stats.TotalCompletedToday = 0
	}

	// 今週の完了率を計算
	weekStart := todayStart.Add(-7 * 24 * time.Hour)
	weeklyCompleted, _ := h.repo.CountCompletedSince(c.Request.Context(), userID, weekStart)
	weeklyTarget := len(items) * 7 // 1日1回 × 7日
	if weeklyTarget > 0 {
		stats.WeeklyCompletionRate = float64(weeklyCompleted) / float64(weeklyTarget) * 100
	}

	c.JSON(http.StatusOK, stats)
}

// GetItems godoc
// @Summary Get review items by priority
// @Tags review
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param priority query string false "Priority filter (urgent, recommended, optional)"
// @Success 200 {object} map[string][]models.ReviewItem
// @Router /api/v1/review/items [get]
func (h *ReviewHandler) GetItems(c *gin.Context) {
	userID := c.GetString("user_id")
	priorityFilter := c.Query("priority")

	items, err := h.repo.FindByUserID(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch review items"})
		return
	}

	// 優先度でフィルタリング
	var filteredItems []*models.ReviewItem
	for _, item := range items {
		priority := h.srsAlgo.CalculatePriority(item.NextReview)
		item.Priority = priority

		if priorityFilter == "" || priority == priorityFilter {
			filteredItems = append(filteredItems, item)
		}
	}

	c.JSON(http.StatusOK, gin.H{"items": filteredItems})
}

// SubmitReview godoc
// @Summary Submit review result
// @Tags review
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param result body models.ReviewResult true "Review result"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/review/submit [post]
func (h *ReviewHandler) SubmitReview(c *gin.Context) {
	var result models.ReviewResult
	if err := c.ShouldBindJSON(&result); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	userID := c.GetString("user_id")

	// 復習アイテムを取得
	item, err := h.repo.FindByID(c.Request.Context(), result.ItemID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Review item not found"})
		return
	}

	// 所有権チェック
	if item.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden"})
		return
	}

	// SRSアルゴリズムで次の復習日時を計算
	nextInterval, nextEaseFactor, nextReview := h.srsAlgo.CalculateNextReview(
		item.EaseFactor,
		item.IntervalDays,
		result.Score,
	)

	// 習熟度を更新
	newMasteryLevel := item.MasteryLevel
	if result.Score >= 70 {
		newMasteryLevel += 10
		if newMasteryLevel > 100 {
			newMasteryLevel = 100
		}
	} else if result.Score < 50 {
		newMasteryLevel -= 5
		if newMasteryLevel < 0 {
			newMasteryLevel = 0
		}
	}

	// アイテムを更新
	item.MasteryLevel = newMasteryLevel
	item.IntervalDays = nextInterval
	item.EaseFactor = nextEaseFactor
	item.LastReviewed = time.Now()
	item.NextReview = nextReview
	item.ReviewCount++

	if err := h.repo.Update(c.Request.Context(), item); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update review item"})
		return
	}

	// 履歴を保存
	history := &models.ReviewHistory{
		ReviewItemID: item.ID,
		UserID:       userID,
		Score:        result.Score,
		ReviewedAt:   time.Now(),
	}

	if err := h.repo.SaveHistory(c.Request.Context(), history); err != nil {
		// エラーログは出すが、レスポンスは成功を返す
		c.JSON(http.StatusOK, gin.H{
			"success":     true,
			"next_review": nextReview.Format(time.RFC3339),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success":     true,
		"next_review": nextReview.Format(time.RFC3339),
	})
}

// RegisterRoutes registers review routes
func (h *ReviewHandler) RegisterRoutes(rg *gin.RouterGroup) {
	review := rg.Group("/review")
	review.Use(middleware.AuthRequired())
	{
		review.GET("/stats", h.GetStats)
		review.GET("/items", h.GetItems)
		review.POST("/submit", h.SubmitReview)
	}
}
```

### 5. Repository Interface

```go
package repository

import (
	"context"
	"time"
	"github.com/clearclown/HaiLanGo/backend/internal/models"
)

type ReviewRepository interface {
	Create(ctx context.Context, item *models.ReviewItem) error
	FindByID(ctx context.Context, id string) (*models.ReviewItem, error)
	FindByUserID(ctx context.Context, userID string) ([]*models.ReviewItem, error)
	Update(ctx context.Context, item *models.ReviewItem) error
	Delete(ctx context.Context, id string) error

	// 統計用
	CountCompletedToday(ctx context.Context, userID string, since time.Time) (int, error)
	CountCompletedSince(ctx context.Context, userID string, since time.Time) (int, error)

	// 履歴
	SaveHistory(ctx context.Context, history *models.ReviewHistory) error
}
```

### 6. Models

```go
package models

import "time"

type ReviewItem struct {
	ID            string    `json:"id"`
	UserID        string    `json:"user_id"`
	BookID        string    `json:"book_id"`
	PageNumber    int       `json:"page_number"`
	Type          string    `json:"type"` // word, phrase
	Text          string    `json:"text"`
	Translation   string    `json:"translation"`
	Language      string    `json:"language"`
	MasteryLevel  int       `json:"mastery_level"`
	IntervalDays  int       `json:"-"`
	EaseFactor    float64   `json:"-"`
	LastReviewed  time.Time `json:"last_reviewed"`
	NextReview    time.Time `json:"next_review"`
	ReviewCount   int       `json:"-"`
	Priority      string    `json:"priority"` // urgent, recommended, optional
	CreatedAt     time.Time `json:"-"`
	UpdatedAt     time.Time `json:"-"`
}

type ReviewStats struct {
	UrgentCount           int     `json:"urgent_count"`
	RecommendedCount      int     `json:"recommended_count"`
	OptionalCount         int     `json:"optional_count"`
	TotalCompletedToday   int     `json:"total_completed_today"`
	WeeklyCompletionRate  float64 `json:"weekly_completion_rate"`
}

type ReviewResult struct {
	ItemID      string    `json:"item_id" binding:"required"`
	Score       int       `json:"score" binding:"required,min=0,max=100"`
	CompletedAt time.Time `json:"completed_at" binding:"required"`
}

type ReviewHistory struct {
	ID           string    `json:"id"`
	ReviewItemID string    `json:"review_item_id"`
	UserID       string    `json:"user_id"`
	Score        int       `json:"score"`
	ReviewedAt   time.Time `json:"reviewed_at"`
}
```

## 完了条件（Definition of Done）

- [ ] データベーススキーマが適用されている
- [ ] SRSアルゴリズムが実装されている（SM2）
- [ ] `handler/review.go` が実装されている
- [ ] `repository/review.go` が実装されている
- [ ] `router/router.go` にルートが登録されている
- [ ] すべてのエンドポイントが動作する
- [ ] ユニットテストが書かれ、すべてパスする
- [ ] フロントエンドの Review ページが正常に動作する
- [ ] E2Eテストが成功する（review.spec.ts）

## 検証方法

### 1. cURL テスト
```bash
# 統計取得
curl -X GET http://localhost:8080/api/v1/review/stats \
  -H "Authorization: Bearer {token}"

# アイテム取得
curl -X GET "http://localhost:8080/api/v1/review/items?priority=urgent" \
  -H "Authorization: Bearer {token}"

# 復習送信
curl -X POST http://localhost:8080/api/v1/review/submit \
  -H "Authorization: Bearer {token}" \
  -H "Content-Type: application/json" \
  -d '{"item_id":"xxx","score":100,"completed_at":"2025-11-14T10:30:00Z"}'
```

### 2. フロントエンド動作確認
- http://localhost:3000/review にアクセス
- 統計が表示される
- 復習カードが表示される
- 復習セッションが動作する

### 3. E2Eテスト実行
```bash
cd frontend/web
pnpm playwright test review.spec.ts
# すべてのテストがパスすること
```

## 参考資料

- [SM2 Algorithm (SuperMemo)](https://www.supermemo.com/en/archives1990-2015/english/ol/sm2)
- Anki SRS実装

## 注意事項

**重要:** このAPIはフロントエンドが既に実装済み。仕様を勝手に変更しないこと。
