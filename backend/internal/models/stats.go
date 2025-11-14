package models

import (
	"time"

	"github.com/google/uuid"
)

// LearningSession represents a single learning session
type LearningSession struct {
	ID        uuid.UUID `json:"id" db:"id"`
	UserID    uuid.UUID `json:"user_id" db:"user_id"`
	BookID    uuid.UUID `json:"book_id" db:"book_id"`
	PageID    uuid.UUID `json:"page_id" db:"page_id"`
	StartTime time.Time `json:"start_time" db:"start_time"`
	EndTime   time.Time `json:"end_time" db:"end_time"`
	Duration  int       `json:"duration" db:"duration"` // in seconds
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

// VocabularyProgress represents vocabulary learning progress
type VocabularyProgress struct {
	ID           uuid.UUID `json:"id" db:"id"`
	UserID       uuid.UUID `json:"user_id" db:"user_id"`
	Word         string    `json:"word" db:"word"`
	Language     string    `json:"language" db:"language"`
	MasteryLevel int       `json:"mastery_level" db:"mastery_level"` // 0-100
	LastReviewed time.Time `json:"last_reviewed" db:"last_reviewed"`
	ReviewCount  int       `json:"review_count" db:"review_count"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}

// PhraseProgress represents phrase learning progress
type PhraseProgress struct {
	ID           uuid.UUID `json:"id" db:"id"`
	UserID       uuid.UUID `json:"user_id" db:"user_id"`
	Phrase       string    `json:"phrase" db:"phrase"`
	Language     string    `json:"language" db:"language"`
	MasteryLevel int       `json:"mastery_level" db:"mastery_level"` // 0-100
	LastReviewed time.Time `json:"last_reviewed" db:"last_reviewed"`
	ReviewCount  int       `json:"review_count" db:"review_count"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}

// PronunciationScore represents pronunciation evaluation results
type PronunciationScoreRecord struct {
	ID        uuid.UUID `json:"id" db:"id"`
	UserID    uuid.UUID `json:"user_id" db:"user_id"`
	Text      string    `json:"text" db:"text"`
	Language  string    `json:"language" db:"language"`
	Score     float64   `json:"score" db:"score"` // 0-100
	Accuracy  float64   `json:"accuracy" db:"accuracy"`
	Fluency   float64   `json:"fluency" db:"fluency"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

// DashboardStats represents the statistics for dashboard
type DashboardStats struct {
	LearningTime      LearningTimeStats      `json:"learning_time"`
	Progress          ProgressStats          `json:"progress"`
	Streak            StreakStats            `json:"streak"`
	PronunciationAvg  float64                `json:"pronunciation_avg"`
	WeakWords         []string               `json:"weak_words"`
	LearningTimeChart []LearningTimeDataPoint `json:"learning_time_chart"`
	ProgressChart     []ProgressDataPoint     `json:"progress_chart"`
}

// LearningTimeStats represents learning time statistics
type LearningTimeStats struct {
	TotalSeconds  int     `json:"total_seconds"`
	TotalHours    float64 `json:"total_hours"`
	DailyAverage  int     `json:"daily_average"`   // in seconds
	WeeklyAverage int     `json:"weekly_average"`  // in seconds
	MonthlyAverage int    `json:"monthly_average"` // in seconds
}

// ProgressStats represents progress statistics
type ProgressStats struct {
	CompletedPages  int `json:"completed_pages"`
	MasteredWords   int `json:"mastered_words"`
	MasteredPhrases int `json:"mastered_phrases"`
	CompletedBooks  int `json:"completed_books"`
}

// StreakStats represents streak statistics
type StreakStats struct {
	CurrentStreak int       `json:"current_streak"`
	LongestStreak int       `json:"longest_streak"`
	LastStudyDate time.Time `json:"last_study_date"`
}

// LearningTimeDataPoint represents a data point for learning time chart
type LearningTimeDataPoint struct {
	Date    time.Time `json:"date"`
	Seconds int       `json:"seconds"`
}

// ProgressDataPoint represents a data point for progress chart
type ProgressDataPoint struct {
	Date   time.Time `json:"date"`
	Words  int       `json:"words"`
	Phrases int      `json:"phrases"`
	Pages  int       `json:"pages"`
}
