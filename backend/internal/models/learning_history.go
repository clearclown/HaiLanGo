package models

import (
	"time"

	"github.com/google/uuid"
)

// LearningHistory は学習履歴を表すモデル
type LearningHistory struct {
	ID          uuid.UUID  `json:"id" db:"id"`
	UserID      uuid.UUID  `json:"userId" db:"user_id"`
	BookID      uuid.UUID  `json:"bookId" db:"book_id"`
	PageID      uuid.UUID  `json:"pageId" db:"page_id"`
	PageNumber  int        `json:"pageNumber" db:"page_number"`
	IsCompleted bool       `json:"isCompleted" db:"is_completed"`
	StudyTime   int        `json:"studyTime" db:"study_time"` // 秒単位
	CompletedAt *time.Time `json:"completedAt,omitempty" db:"completed_at"`
	CreatedAt   time.Time  `json:"createdAt" db:"created_at"`
	UpdatedAt   time.Time  `json:"updatedAt" db:"updated_at"`
}

// LearningProgress は書籍全体の学習進捗を表すモデル
type LearningProgress struct {
	BookID         uuid.UUID `json:"bookId" db:"book_id"`
	TotalPages     int       `json:"totalPages" db:"total_pages"`
	CompletedPages int       `json:"completedPages" db:"completed_pages"`
	Progress       float64   `json:"progress"` // 0-100のパーセンテージ
	TotalStudyTime int       `json:"totalStudyTime" db:"total_study_time"` // 秒単位
}
