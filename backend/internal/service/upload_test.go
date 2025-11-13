package service

import (
	"bytes"
	"context"
	"mime/multipart"
	"testing"

	"github.com/clearclown/HaiLanGo/internal/models"
	"github.com/clearclown/HaiLanGo/pkg/storage"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestUploadService_CreateBook は書籍の作成をテストする
func TestUploadService_CreateBook(t *testing.T) {
	tempDir := t.TempDir()
	store := storage.NewLocalStorage(tempDir)
	service := NewUploadService(store)

	ctx := context.Background()
	userID := uuid.New()

	metadata := models.BookMetadata{
		Title:          "ロシア語入門",
		TargetLanguage: "ru",
		NativeLanguage: "ja",
	}

	book, err := service.CreateBook(ctx, userID, metadata)
	require.NoError(t, err)

	assert.NotEqual(t, uuid.Nil, book.ID)
	assert.Equal(t, userID, book.UserID)
	assert.Equal(t, metadata.Title, book.Title)
	assert.Equal(t, metadata.TargetLanguage, book.TargetLanguage)
	assert.Equal(t, metadata.NativeLanguage, book.NativeLanguage)
	assert.Equal(t, models.BookStatusUploading, book.Status)
}

// TestUploadService_UploadFile はファイルアップロードをテストする
func TestUploadService_UploadFile(t *testing.T) {
	tempDir := t.TempDir()
	store := storage.NewLocalStorage(tempDir)
	service := NewUploadService(store)

	ctx := context.Background()
	userID := uuid.New()
	bookID := uuid.New()

	tests := []struct {
		name        string
		fileName    string
		fileSize    int64
		contentType string
		data        []byte
		wantErr     bool
		errMsg      string
	}{
		{
			name:        "正常なPDFファイルアップロード",
			fileName:    "book.pdf",
			fileSize:    1024,
			contentType: "application/pdf",
			data:        []byte("PDF content"),
			wantErr:     false,
		},
		{
			name:        "正常なPNGファイルアップロード",
			fileName:    "page1.png",
			fileSize:    2048,
			contentType: "image/png",
			data:        []byte("PNG content"),
			wantErr:     false,
		},
		{
			name:        "サイズ超過ファイル",
			fileName:    "large.pdf",
			fileSize:    150 * 1024 * 1024, // 150MB
			contentType: "application/pdf",
			data:        []byte("large content"),
			wantErr:     true,
			errMsg:      "file size exceeds limit",
		},
		{
			name:        "不正なファイル形式",
			fileName:    "malware.exe",
			fileSize:    1024,
			contentType: "application/x-msdownload",
			data:        []byte("malware"),
			wantErr:     true,
			errMsg:      "unsupported file type",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// モックの multipart.FileHeader を作成
			fileHeader := &multipart.FileHeader{
				Filename: tt.fileName,
				Size:     tt.fileSize,
				Header:   make(map[string][]string),
			}
			fileHeader.Header.Set("Content-Type", tt.contentType)

			bookFile, err := service.UploadFile(ctx, userID, bookID, fileHeader, bytes.NewReader(tt.data))

			if tt.wantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
				return
			}

			require.NoError(t, err)
			assert.NotEqual(t, uuid.Nil, bookFile.ID)
			assert.Equal(t, bookID, bookFile.BookID)
			assert.Equal(t, tt.fileName, bookFile.FileName)
			assert.Equal(t, tt.fileSize, bookFile.FileSize)
		})
	}
}

// TestUploadService_ValidateFiles は複数ファイルの検証をテストする
func TestUploadService_ValidateFiles(t *testing.T) {
	tempDir := t.TempDir()
	store := storage.NewLocalStorage(tempDir)
	service := NewUploadService(store)

	tests := []struct {
		name    string
		files   []*multipart.FileHeader
		wantErr bool
		errMsg  string
	}{
		{
			name: "正常なファイル群",
			files: []*multipart.FileHeader{
				{
					Filename: "page1.png",
					Size:     1024,
					Header:   makeHeader("image/png"),
				},
				{
					Filename: "page2.png",
					Size:     2048,
					Header:   makeHeader("image/png"),
				},
			},
			wantErr: false,
		},
		{
			name: "ファイル数が0",
			files: []*multipart.FileHeader{},
			wantErr: true,
			errMsg:  "no files provided",
		},
		{
			name: "不正なファイルを含む",
			files: []*multipart.FileHeader{
				{
					Filename: "page1.png",
					Size:     1024,
					Header:   makeHeader("image/png"),
				},
				{
					Filename: "malware.exe",
					Size:     1024,
					Header:   makeHeader("application/x-msdownload"),
				},
			},
			wantErr: true,
			errMsg:  "unsupported file type",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := service.ValidateFiles(tt.files)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// makeHeader はテスト用のヘッダーを作成するヘルパー関数
func makeHeader(contentType string) map[string][]string {
	header := make(map[string][]string)
	header["Content-Type"] = []string{contentType}
	return header
}

// TestUploadService_GetUploadProgress はアップロード進捗の取得をテストする
func TestUploadService_GetUploadProgress(t *testing.T) {
	tempDir := t.TempDir()
	store := storage.NewLocalStorage(tempDir)
	service := NewUploadService(store)

	ctx := context.Background()
	bookID := uuid.New()

	// 進捗を更新
	progress := &models.UploadProgress{
		BookID:        bookID,
		TotalFiles:    5,
		UploadedFiles: 3,
		TotalBytes:    5000,
		UploadedBytes: 3000,
		Status:        "uploading",
	}

	err := service.UpdateUploadProgress(ctx, progress)
	require.NoError(t, err)

	// 進捗を取得
	retrieved, err := service.GetUploadProgress(ctx, bookID)
	require.NoError(t, err)

	assert.Equal(t, bookID, retrieved.BookID)
	assert.Equal(t, 5, retrieved.TotalFiles)
	assert.Equal(t, 3, retrieved.UploadedFiles)
	assert.Equal(t, int64(5000), retrieved.TotalBytes)
	assert.Equal(t, int64(3000), retrieved.UploadedBytes)
	assert.Equal(t, "uploading", retrieved.Status)
}
