package models

import (
	"time"

	"github.com/google/uuid"
)

// Book represents a language learning book
type Book struct {
	ID               uuid.UUID `json:"id" db:"id"`
	UserID           uuid.UUID `json:"user_id" db:"user_id"`
	Title            string    `json:"title" db:"title"`
	TargetLanguage   string    `json:"target_language" db:"target_language"`
	NativeLanguage   string    `json:"native_language" db:"native_language"`
	ReferenceLanguage string   `json:"reference_language" db:"reference_language"`
	TotalPages       int       `json:"total_pages" db:"total_pages"`
	CreatedAt        time.Time `json:"created_at" db:"created_at"`
	UpdatedAt        time.Time `json:"updated_at" db:"updated_at"`
}

// Page represents a page in a book
type Page struct {
	ID            uuid.UUID `json:"id" db:"id"`
	BookID        uuid.UUID `json:"book_id" db:"book_id"`
	PageNumber    int       `json:"page_number" db:"page_number"`
	ImageURL      string    `json:"image_url" db:"image_url"`
	OCRText       string    `json:"ocr_text" db:"ocr_text"`
	CorrectedText *string   `json:"corrected_text,omitempty" db:"corrected_text"`
	IsCompleted   bool      `json:"is_completed" db:"is_completed"`
	CreatedAt     time.Time `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time `json:"updated_at" db:"updated_at"`
}
