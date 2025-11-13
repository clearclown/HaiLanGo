package models

import (
	"time"

	"github.com/google/uuid"
)

// PatternType represents the type of conversation pattern
type PatternType string

const (
	PatternTypeGreeting     PatternType = "greeting"
	PatternTypeQuestion     PatternType = "question"
	PatternTypeResponse     PatternType = "response"
	PatternTypeRequest      PatternType = "request"
	PatternTypeConfirmation PatternType = "confirmation"
	PatternTypeOther        PatternType = "other"
)

// Pattern represents a conversation pattern extracted from books
type Pattern struct {
	ID          uuid.UUID   `json:"id" db:"id"`
	BookID      uuid.UUID   `json:"book_id" db:"book_id"`
	Type        PatternType `json:"type" db:"type"`
	Pattern     string      `json:"pattern" db:"pattern"`
	Translation string      `json:"translation" db:"translation"`
	Frequency   int         `json:"frequency" db:"frequency"`
	CreatedAt   time.Time   `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time   `json:"updated_at" db:"updated_at"`
}

// PatternExample represents an example usage of a pattern
type PatternExample struct {
	ID            uuid.UUID `json:"id" db:"id"`
	PatternID     uuid.UUID `json:"pattern_id" db:"pattern_id"`
	PageNumber    int       `json:"page_number" db:"page_number"`
	OriginalText  string    `json:"original_text" db:"original_text"`
	TranslatedText string   `json:"translated_text" db:"translated_text"`
	Context       string    `json:"context" db:"context"`
	CreatedAt     time.Time `json:"created_at" db:"created_at"`
}

// PatternPractice represents a practice exercise for a pattern
type PatternPractice struct {
	ID                uuid.UUID `json:"id" db:"id"`
	PatternID         uuid.UUID `json:"pattern_id" db:"pattern_id"`
	Question          string    `json:"question" db:"question"`
	CorrectAnswer     string    `json:"correct_answer" db:"correct_answer"`
	AlternativeAnswers []string  `json:"alternative_answers" db:"alternative_answers"`
	Difficulty        int       `json:"difficulty" db:"difficulty"`
	CreatedAt         time.Time `json:"created_at" db:"created_at"`
}

// PatternProgress represents a user's progress on learning a pattern
type PatternProgress struct {
	ID              uuid.UUID  `json:"id" db:"id"`
	UserID          uuid.UUID  `json:"user_id" db:"user_id"`
	PatternID       uuid.UUID  `json:"pattern_id" db:"pattern_id"`
	MasteryLevel    int        `json:"mastery_level" db:"mastery_level"` // 0-100
	PracticeCount   int        `json:"practice_count" db:"practice_count"`
	CorrectCount    int        `json:"correct_count" db:"correct_count"`
	LastPracticedAt *time.Time `json:"last_practiced_at" db:"last_practiced_at"`
	CreatedAt       time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at" db:"updated_at"`
}

// PatternExtractionRequest represents a request to extract patterns
type PatternExtractionRequest struct {
	BookID      uuid.UUID `json:"book_id"`
	PageStart   int       `json:"page_start"`
	PageEnd     int       `json:"page_end"`
	MinFrequency int      `json:"min_frequency"`
}

// PatternExtractionResponse represents the response from pattern extraction
type PatternExtractionResponse struct {
	Patterns      []Pattern        `json:"patterns"`
	TotalFound    int              `json:"total_found"`
	ProcessedPages int             `json:"processed_pages"`
	Duration      time.Duration    `json:"duration"`
}
