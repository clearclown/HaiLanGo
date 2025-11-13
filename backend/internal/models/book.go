package models

import (
	"time"

	"github.com/google/uuid"
)

// Book は書籍情報を表すモデル
type Book struct {
	ID              uuid.UUID  `json:"id" db:"id"`
	UserID          uuid.UUID  `json:"user_id" db:"user_id"`
	Title           string     `json:"title" db:"title"`
	TargetLanguage  string     `json:"target_language" db:"target_language"`   // 学習先言語
	NativeLanguage  string     `json:"native_language" db:"native_language"`   // 母国語
	ReferenceLanguage string   `json:"reference_language" db:"reference_language"` // 参照言語（本に使用されている言語）
	TotalPages      int        `json:"total_pages" db:"total_pages"`
	ProcessedPages  int        `json:"processed_pages" db:"processed_pages"`
	OCRStatus       OCRStatus  `json:"ocr_status" db:"ocr_status"`
	CreatedAt       time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at" db:"updated_at"`
}

// OCRStatus はOCR処理の状態を表す
type OCRStatus string

const (
	OCRStatusPending    OCRStatus = "pending"     // 処理待ち
	OCRStatusProcessing OCRStatus = "processing"  // 処理中
	OCRStatusCompleted  OCRStatus = "completed"   // 完了
	OCRStatusFailed     OCRStatus = "failed"      // 失敗
)

// Page は書籍のページ情報を表すモデル
type Page struct {
	ID             uuid.UUID  `json:"id" db:"id"`
	BookID         uuid.UUID  `json:"book_id" db:"book_id"`
	PageNumber     int        `json:"page_number" db:"page_number"`
	ImageURL       string     `json:"image_url" db:"image_url"`
	OCRText        string     `json:"ocr_text" db:"ocr_text"`
	OCRConfidence  float64    `json:"ocr_confidence" db:"ocr_confidence"`
	DetectedLang   string     `json:"detected_lang" db:"detected_lang"`
	OCRStatus      OCRStatus  `json:"ocr_status" db:"ocr_status"`
	OCRError       *string    `json:"ocr_error,omitempty" db:"ocr_error"`
	CreatedAt      time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at" db:"updated_at"`
}

// OCRJob はOCR処理ジョブを表すモデル
type OCRJob struct {
	ID        uuid.UUID  `json:"id" db:"id"`
	BookID    uuid.UUID  `json:"book_id" db:"book_id"`
	PageID    *uuid.UUID `json:"page_id,omitempty" db:"page_id"` // nilの場合は全ページ
	Status    OCRStatus  `json:"status" db:"status"`
	Progress  int        `json:"progress" db:"progress"` // 0-100
	Error     *string    `json:"error,omitempty" db:"error"`
	CreatedAt time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt time.Time  `json:"updated_at" db:"updated_at"`
}
