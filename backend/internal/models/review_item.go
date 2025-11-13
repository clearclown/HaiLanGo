package models

import (
	"time"

	"github.com/google/uuid"
)

// ReviewItem は復習項目を表すモデル
type ReviewItem struct {
	ID             uuid.UUID  `json:"id" db:"id"`
	UserID         uuid.UUID  `json:"user_id" db:"user_id"`
	BookID         uuid.UUID  `json:"book_id" db:"book_id"`
	PageNumber     int        `json:"page_number" db:"page_number"`
	ItemType       string     `json:"item_type" db:"item_type"` // "phrase" or "word"
	Content        string     `json:"content" db:"content"`
	Translation    string     `json:"translation" db:"translation"`
	ReviewCount    int        `json:"review_count" db:"review_count"`
	LastReviewDate *time.Time `json:"last_review_date" db:"last_review_date"`
	NextReviewDate *time.Time `json:"next_review_date" db:"next_review_date"`
	LastScore      int        `json:"last_score" db:"last_score"`
	CreatedAt      time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at" db:"updated_at"`
}

// ReviewHistory は復習履歴を表すモデル
type ReviewHistory struct {
	ID           uuid.UUID `json:"id" db:"id"`
	ReviewItemID uuid.UUID `json:"review_item_id" db:"review_item_id"`
	UserID       uuid.UUID `json:"user_id" db:"user_id"`
	Score        int       `json:"score" db:"score"`
	ReviewedAt   time.Time `json:"reviewed_at" db:"reviewed_at"`
	TimeSpentSec int       `json:"time_spent_sec" db:"time_spent_sec"`
}

// ReviewStats は復習統計を表すモデル
type ReviewStats struct {
	UserID              uuid.UUID `json:"user_id"`
	TotalReviewItems    int       `json:"total_review_items"`
	UrgentItems         int       `json:"urgent_items"`
	RecommendedItems    int       `json:"recommended_items"`
	RelaxedItems        int       `json:"relaxed_items"`
	WeeklyReviewCount   int       `json:"weekly_review_count"`
	CurrentStreak       int       `json:"current_streak"`
	LongestStreak       int       `json:"longest_streak"`
	AverageScore        float64   `json:"average_score"`
}

// ReviewPriority は復習優先度
type ReviewPriority string

const (
	PriorityUrgent      ReviewPriority = "urgent"      // 今日中に復習が必要
	PriorityRecommended ReviewPriority = "recommended" // 今日復習すると効果的
	PriorityRelaxed     ReviewPriority = "relaxed"     // 明日以降でもOK
)

// GetPriority は次回復習日に基づいて優先度を返す
func (r *ReviewItem) GetPriority(now time.Time) ReviewPriority {
	if r.NextReviewDate == nil {
		return PriorityUrgent
	}

	daysDiff := r.NextReviewDate.Sub(now).Hours() / 24

	if daysDiff < 0 {
		// 過去の日付 = 緊急
		return PriorityUrgent
	} else if daysDiff <= 1 {
		// 今日〜明日 = 推奨
		return PriorityRecommended
	} else {
		// 明日以降 = 余裕あり
		return PriorityRelaxed
	}
}
