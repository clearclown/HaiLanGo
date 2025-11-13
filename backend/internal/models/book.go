package models

import (
	"time"

	"github.com/google/uuid"
)

// Book は書籍の情報を保持する
type Book struct {
	ID               uuid.UUID  `json:"id" db:"id"`
	UserID           uuid.UUID  `json:"user_id" db:"user_id"`
	Title            string     `json:"title" db:"title" binding:"required"`
	TargetLanguage   string     `json:"target_language" db:"target_language" binding:"required"` // 学習先言語
	NativeLanguage   string     `json:"native_language" db:"native_language" binding:"required"` // 母国語
	ReferenceLanguage string    `json:"reference_language,omitempty" db:"reference_language"`    // 参照言語（オプション）
	CoverImageURL    string     `json:"cover_image_url,omitempty" db:"cover_image_url"`
	TotalPages       int        `json:"total_pages" db:"total_pages"`
	Status           BookStatus `json:"status" db:"status"`
	CreatedAt        time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at" db:"updated_at"`
}

// BookStatus は書籍の処理状態を表す
type BookStatus string

const (
	BookStatusUploading  BookStatus = "uploading"   // アップロード中
	BookStatusProcessing BookStatus = "processing"  // OCR処理中
	BookStatusReady      BookStatus = "ready"       // 学習可能
	BookStatusFailed     BookStatus = "failed"      // 処理失敗
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
	ID           uuid.UUID `json:"id" db:"id"`
	BookID       uuid.UUID `json:"book_id" db:"book_id"`
	FileName     string    `json:"file_name" db:"file_name"`
	FileType     string    `json:"file_type" db:"file_type"` // pdf, png, jpg, heic
	FileSize     int64     `json:"file_size" db:"file_size"` // バイト単位
	StoragePath  string    `json:"storage_path" db:"storage_path"`
	PageNumber   int       `json:"page_number,omitempty" db:"page_number"` // PDFの場合は0、画像の場合はページ番号
	UploadedAt   time.Time `json:"uploaded_at" db:"uploaded_at"`
}

// ChunkUpload はチャンクアップロードの情報を保持する
type ChunkUpload struct {
	ID            uuid.UUID `json:"id" db:"id"`
	BookID        uuid.UUID `json:"book_id" db:"book_id"`
	FileName      string    `json:"file_name" db:"file_name"`
	TotalChunks   int       `json:"total_chunks" db:"total_chunks"`
	UploadedChunks int      `json:"uploaded_chunks" db:"uploaded_chunks"`
	Status        string    `json:"status" db:"status"` // pending, uploading, completed, failed
	CreatedAt     time.Time `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time `json:"updated_at" db:"updated_at"`
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
