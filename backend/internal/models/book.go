package models

import (
	"time"

	"github.com/google/uuid"
)

// Book は書籍の情報を保持する
type Book struct {
	ID                uuid.UUID  `json:"id" db:"id"`
	UserID            uuid.UUID  `json:"user_id" db:"user_id"`
	Title             string     `json:"title" db:"title" binding:"required"`
	TargetLanguage    string     `json:"target_language" db:"target_language" binding:"required"`   // 学習先言語
	NativeLanguage    string     `json:"native_language" db:"native_language" binding:"required"`   // 母国語
	ReferenceLanguage string     `json:"reference_language,omitempty" db:"reference_language"`      // 参照言語（本に使用されている言語）
	CoverImageURL     string     `json:"cover_image_url,omitempty" db:"cover_image_url"`
	TotalPages        int        `json:"total_pages" db:"total_pages"`
	ProcessedPages    int        `json:"processed_pages" db:"processed_pages"`
	Status            BookStatus `json:"status" db:"status"`
	OCRStatus         OCRStatus  `json:"ocr_status" db:"ocr_status"`
	CreatedAt         time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at" db:"updated_at"`
}

// BookStatus は書籍の処理状態を表す
type BookStatus string

const (
	BookStatusUploading  BookStatus = "uploading"   // アップロード中
	BookStatusProcessing BookStatus = "processing"  // OCR処理中
	BookStatusReady      BookStatus = "ready"       // 学習可能
	BookStatusFailed     BookStatus = "failed"      // 処理失敗
)

// OCRStatus はOCR処理の状態を表す
type OCRStatus string

const (
	OCRStatusPending    OCRStatus = "pending"     // 処理待ち
	OCRStatusProcessing OCRStatus = "processing"  // 処理中
	OCRStatusCompleted  OCRStatus = "completed"   // 完了
	OCRStatusFailed     OCRStatus = "failed"      // 失敗
)

// BookMetadata は書籍のメタデータを保持する
type BookMetadata struct {
	Title             string `json:"title" binding:"required"`
	TargetLanguage    string `json:"target_language" binding:"required"`
	NativeLanguage    string `json:"native_language" binding:"required"`
	ReferenceLanguage string `json:"reference_language,omitempty"`
}

// BookFile は書籍のファイル情報を保持する
type BookFile struct {
	ID          uuid.UUID `json:"id" db:"id"`
	BookID      uuid.UUID `json:"book_id" db:"book_id"`
	FileName    string    `json:"file_name" db:"file_name"`
	FileType    string    `json:"file_type" db:"file_type"` // pdf, png, jpg, heic
	FileSize    int64     `json:"file_size" db:"file_size"` // バイト単位
	StoragePath string    `json:"storage_path" db:"storage_path"`
	PageNumber  int       `json:"page_number,omitempty" db:"page_number"` // PDFの場合は0、画像の場合はページ番号
	UploadedAt  time.Time `json:"uploaded_at" db:"uploaded_at"`
}

// ChunkUpload はチャンクアップロードの情報を保持する
type ChunkUpload struct {
	ID             uuid.UUID `json:"id" db:"id"`
	BookID         uuid.UUID `json:"book_id" db:"book_id"`
	FileName       string    `json:"file_name" db:"file_name"`
	TotalChunks    int       `json:"total_chunks" db:"total_chunks"`
	UploadedChunks int       `json:"uploaded_chunks" db:"uploaded_chunks"`
	Status         string    `json:"status" db:"status"` // pending, uploading, completed, failed
	CreatedAt      time.Time `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time `json:"updated_at" db:"updated_at"`
}

// UploadProgress はアップロードの進捗情報を保持する
type UploadProgress struct {
	BookID        uuid.UUID `json:"book_id"`
	TotalFiles    int       `json:"total_files"`
	UploadedFiles int       `json:"uploaded_files"`
	TotalBytes    int64     `json:"total_bytes"`
	UploadedBytes int64     `json:"uploaded_bytes"`
	Status        string    `json:"status"`
	Message       string    `json:"message,omitempty"`
}

// Page は書籍のページ情報を表すモデル
type Page struct {
	ID            uuid.UUID  `json:"id" db:"id"`
	BookID        uuid.UUID  `json:"book_id" db:"book_id"`
	PageNumber    int        `json:"page_number" db:"page_number"`
	ImageURL      string     `json:"image_url" db:"image_url"`
	OCRText       string     `json:"ocr_text" db:"ocr_text"`
	OCRConfidence float64    `json:"ocr_confidence" db:"ocr_confidence"`
	DetectedLang  string     `json:"detected_lang" db:"detected_lang"`
	OCRStatus     OCRStatus  `json:"ocr_status" db:"ocr_status"`
	OCRError      *string    `json:"ocr_error,omitempty" db:"ocr_error"`
	CreatedAt     time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at" db:"updated_at"`
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

// PageWithProgress はページ情報と学習進捗を含む
type PageWithProgress struct {
	*Page
	IsCompleted bool       `json:"is_completed"`
	CompletedAt *time.Time `json:"completed_at,omitempty"`
	StudyTime   int        `json:"study_time"` // 秒単位
}
