package models

import "time"

// NotificationType represents the type of notification
type NotificationType string

const (
	// NotificationTypeOCRProgress represents OCR processing progress
	NotificationTypeOCRProgress NotificationType = "ocr_progress"
	// NotificationTypeTTSProgress represents TTS generation progress
	NotificationTypeTTSProgress NotificationType = "tts_progress"
	// NotificationTypeLearningUpdate represents learning progress update
	NotificationTypeLearningUpdate NotificationType = "learning_update"
	// NotificationTypeError represents an error notification
	NotificationTypeError NotificationType = "error"
	// NotificationTypePing represents a ping message
	NotificationTypePing NotificationType = "ping"
	// NotificationTypePong represents a pong message
	NotificationTypePong NotificationType = "pong"
)

// Notification represents a WebSocket notification message
type Notification struct {
	Type      NotificationType `json:"type"`
	Data      interface{}      `json:"data"`
	Timestamp time.Time        `json:"timestamp"`
}

// OCRProgressData represents OCR processing progress data
type OCRProgressData struct {
	BookID          string  `json:"book_id"`
	TotalPages      int     `json:"total_pages"`
	ProcessedPages  int     `json:"processed_pages"`
	CurrentPage     int     `json:"current_page"`
	Progress        float64 `json:"progress"` // 0-100
	EstimatedTimeMS int64   `json:"estimated_time_ms"`
	Status          string  `json:"status"` // "processing", "completed", "failed"
}

// TTSProgressData represents TTS generation progress data
type TTSProgressData struct {
	BookID          string  `json:"book_id"`
	PageNumber      int     `json:"page_number"`
	TotalSegments   int     `json:"total_segments"`
	ProcessedSegments int   `json:"processed_segments"`
	Progress        float64 `json:"progress"` // 0-100
	Status          string  `json:"status"` // "processing", "completed", "failed"
}

// LearningUpdateData represents learning progress update data
type LearningUpdateData struct {
	UserID         string `json:"user_id"`
	BookID         string `json:"book_id"`
	PageNumber     int    `json:"page_number"`
	CompletedPages int    `json:"completed_pages"`
	TotalPages     int    `json:"total_pages"`
	LearnedWords   int    `json:"learned_words"`
	StudyTimeMS    int64  `json:"study_time_ms"`
}

// ErrorData represents error notification data
type ErrorData struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}
