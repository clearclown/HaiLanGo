package models

import (
	"time"

	"github.com/google/uuid"
)

// Page はページを表すモデル
type Page struct {
	ID          uuid.UUID  `json:"id" db:"id"`
	BookID      uuid.UUID  `json:"bookId" db:"book_id"`
	PageNumber  int        `json:"pageNumber" db:"page_number"`
	ImageURL    string     `json:"imageUrl" db:"image_url"`
	OCRText     string     `json:"ocrText" db:"ocr_text"`
	Translation string     `json:"translation" db:"translation"`
	AudioURL    string     `json:"audioUrl" db:"audio_url"`
	CreatedAt   time.Time  `json:"createdAt" db:"created_at"`
	UpdatedAt   time.Time  `json:"updatedAt" db:"updated_at"`
}

// PageWithProgress はページと学習進捗を含むモデル
type PageWithProgress struct {
	Page
	IsCompleted bool      `json:"isCompleted" db:"is_completed"`
	CompletedAt *time.Time `json:"completedAt,omitempty" db:"completed_at"`
}
