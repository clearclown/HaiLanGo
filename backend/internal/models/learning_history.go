package models

import (
	"time"

	"github.com/google/uuid"
)

// LearningHistory は学習履歴を表すモデル
type LearningHistory struct {
	ID          uuid.UUID  `json:"id" db:"id"`
	UserID      uuid.UUID  `json:"user_id" db:"user_id"`
	BookID      uuid.UUID  `json:"book_id" db:"book_id"`
	PageID      uuid.UUID  `json:"page_id" db:"page_id"`
	PageNumber  int        `json:"page_number" db:"page_number"`
	IsCompleted bool       `json:"is_completed" db:"is_completed"`
	StudyTime   int        `json:"study_time" db:"study_time"` // 秒単位
	CompletedAt *time.Time `json:"completed_at,omitempty" db:"completed_at"`
	CreatedAt   time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at" db:"updated_at"`
}

// LearningProgress は書籍全体の学習進捗を表すモデル
type LearningProgress struct {
	BookID         uuid.UUID `json:"book_id" db:"book_id"`
	TotalPages     int       `json:"total_pages" db:"total_pages"`
	CompletedPages int       `json:"completed_pages" db:"completed_pages"`
	Progress       float64   `json:"progress"` // 0-100のパーセンテージ
	TotalStudyTime int       `json:"total_study_time" db:"total_study_time"` // 秒単位
}
