package models

import (
	"time"

	"github.com/google/uuid"
)

// Priority は復習項目の優先度
type Priority string

const (
	PriorityUrgent      Priority = "urgent"      // 緊急（期限切れ）
	PriorityRecommended Priority = "recommended" // 推奨（今日が期限）
	PriorityRelaxed     Priority = "relaxed"     // 余裕あり（期限まで余裕）
)

// ReviewItem は復習項目
type ReviewItem struct {
	ID             uuid.UUID  `json:"id"`
	UserID         uuid.UUID  `json:"user_id"`
	BookID         uuid.UUID  `json:"book_id"`
	PageNumber     int        `json:"page_number"`
	ItemType       string     `json:"item_type"` // word, phrase
	Content        string     `json:"content"`
	Translation    string     `json:"translation"`
	Language       string     `json:"language"`
	MasteryLevel   int        `json:"mastery_level"`
	IntervalDays   int        `json:"-"`
	EaseFactor     float64    `json:"-"`
	LastScore      int        `json:"last_score"`
	LastReviewDate *time.Time `json:"last_review_date"`
	NextReviewDate *time.Time `json:"next_review_date"`
	ReviewCount    int        `json:"review_count"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
}

// GetPriority は現在時刻に基づいて優先度を返す
func (r *ReviewItem) GetPriority(now time.Time) Priority {
	if r.NextReviewDate == nil {
		return PriorityUrgent
	}

	daysUntilDue := r.NextReviewDate.Sub(now).Hours() / 24

	if daysUntilDue < 0 {
		// 期限切れ
		return PriorityUrgent
	} else if daysUntilDue < 1 {
		// 今日が期限
		return PriorityRecommended
	}
	// 余裕あり
	return PriorityRelaxed
}

// ReviewStats は復習統計
type ReviewStats struct {
	UrgentCount          int     `json:"urgent_count"`
	RecommendedCount     int     `json:"recommended_count"`
	OptionalCount        int     `json:"optional_count"`
	TotalCompletedToday  int     `json:"total_completed_today"`
	WeeklyCompletionRate float64 `json:"weekly_completion_rate"`
}

// ReviewResult は復習結果
type ReviewResult struct {
	ItemID      uuid.UUID `json:"item_id" binding:"required"`
	Score       int       `json:"score" binding:"required,min=0,max=100"`
	CompletedAt time.Time `json:"completed_at" binding:"required"`
}

// ReviewHistory は復習履歴
type ReviewHistory struct {
	ID           uuid.UUID `json:"id"`
	ReviewItemID uuid.UUID `json:"review_item_id"`
	UserID       uuid.UUID `json:"user_id"`
	Score        int       `json:"score"`
	TimeSpentSec int       `json:"time_spent_sec"`
	ReviewedAt   time.Time `json:"reviewed_at"`
}
