package models

import "time"

type ReviewItem struct {
	ID           string    `json:"id"`
	UserID       string    `json:"user_id"`
	BookID       string    `json:"book_id"`
	PageNumber   int       `json:"page_number"`
	Type         string    `json:"type"` // word, phrase
	Text         string    `json:"text"`
	Translation  string    `json:"translation"`
	Language     string    `json:"language"`
	MasteryLevel int       `json:"mastery_level"`
	IntervalDays int       `json:"-"`
	EaseFactor   float64   `json:"-"`
	LastReviewed time.Time `json:"last_reviewed"`
	NextReview   time.Time `json:"next_review"`
	ReviewCount  int       `json:"-"`
	Priority     string    `json:"priority"` // urgent, recommended, optional
	CreatedAt    time.Time `json:"-"`
	UpdatedAt    time.Time `json:"-"`
}

type ReviewStats struct {
	UrgentCount          int     `json:"urgent_count"`
	RecommendedCount     int     `json:"recommended_count"`
	OptionalCount        int     `json:"optional_count"`
	TotalCompletedToday  int     `json:"total_completed_today"`
	WeeklyCompletionRate float64 `json:"weekly_completion_rate"`
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
