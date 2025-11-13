package models

import (
	"time"

	"github.com/google/uuid"
)

// OCRTextCorrection represents a manual correction to OCR text
type OCRTextCorrection struct {
	ID           uuid.UUID `json:"id" db:"id"`
	BookID       uuid.UUID `json:"book_id" db:"book_id"`
	PageID       uuid.UUID `json:"page_id" db:"page_id"`
	OriginalText string    `json:"original_text" db:"original_text"`
	CorrectedText string   `json:"corrected_text" db:"corrected_text"`
	UserID       uuid.UUID `json:"user_id" db:"user_id"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}

// OCRCorrectionHistory represents the history of corrections for a page
type OCRCorrectionHistory struct {
	PageID      uuid.UUID           `json:"page_id"`
	Corrections []OCRTextCorrection `json:"corrections"`
	TotalCount  int                 `json:"total_count"`
}

// UpdateOCRTextRequest represents a request to update OCR text
type UpdateOCRTextRequest struct {
	CorrectedText string `json:"corrected_text" binding:"required,max=10000"`
}

// UpdateOCRTextResponse represents the response after updating OCR text
type UpdateOCRTextResponse struct {
	Success    bool              `json:"success"`
	Correction OCRTextCorrection `json:"correction"`
	Message    string            `json:"message,omitempty"`
}
