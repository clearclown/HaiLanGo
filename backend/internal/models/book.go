package models

import (
	"time"

	"github.com/google/uuid"
)

// Book は書籍を表すモデル
type Book struct {
	ID              uuid.UUID  `json:"id" db:"id"`
	UserID          uuid.UUID  `json:"userId" db:"user_id"`
	Title           string     `json:"title" db:"title"`
	TargetLanguage  string     `json:"targetLanguage" db:"target_language"`
	NativeLanguage  string     `json:"nativeLanguage" db:"native_language"`
	ReferenceLanguage string   `json:"referenceLanguage" db:"reference_language"`
	TotalPages      int        `json:"totalPages" db:"total_pages"`
	CompletedPages  int        `json:"completedPages" db:"completed_pages"`
	CreatedAt       time.Time  `json:"createdAt" db:"created_at"`
	UpdatedAt       time.Time  `json:"updatedAt" db:"updated_at"`
}
