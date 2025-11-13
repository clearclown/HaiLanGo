package service

import (
	"context"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"sync"
	"time"

	"github.com/clearclown/HaiLanGo/internal/models"
	"github.com/clearclown/HaiLanGo/pkg/file"
	"github.com/clearclown/HaiLanGo/pkg/storage"
	"github.com/google/uuid"
)

var (
	// ErrNoFilesProvided はファイルが提供されていないエラー
	ErrNoFilesProvided = errors.New("no files provided")
	// ErrBookNotFound は書籍が見つからないエラー
	ErrBookNotFound = errors.New("book not found")
)

// UploadService はファイルアップロードのビジネスロジックを提供する
type UploadService struct {
	storage       storage.Storage
	progress      map[uuid.UUID]*models.UploadProgress // 進捗管理（実際の実装ではRedisなどを使用）
	mu            sync.RWMutex                         // 進捗マップのロック
	chunkService  *ChunkUploadService                  // チャンクアップロードサービス
	chunkUploads  map[uuid.UUID]*models.ChunkUpload    // チャンクアップロード管理
}

// NewUploadService はUploadServiceの新しいインスタンスを作成する
func NewUploadService(storage storage.Storage, tempDir string) *UploadService {
	return &UploadService{
		storage:      storage,
		progress:     make(map[uuid.UUID]*models.UploadProgress),
		chunkService: NewChunkUploadService(storage, tempDir),
		chunkUploads: make(map[uuid.UUID]*models.ChunkUpload),
	}
}

// InitiateChunkUpload はチャンクアップロードを開始する
func (s *UploadService) InitiateChunkUpload(ctx context.Context, bookID uuid.UUID, fileName string, totalChunks int, fileSize int64) (*models.ChunkUpload, error) {
	return s.chunkService.InitiateChunkUpload(ctx, bookID, fileName, totalChunks, fileSize)
}

// UploadChunk はチャンクをアップロードする
func (s *UploadService) UploadChunk(ctx context.Context, uploadID uuid.UUID, chunkNumber int, reader io.Reader) error {
	return s.chunkService.UploadChunk(ctx, uploadID, chunkNumber, reader)
}

// GetChunkUpload はチャンクアップロード情報を取得する
func (s *UploadService) GetChunkUpload(ctx context.Context, uploadID uuid.UUID) (*models.ChunkUpload, error) {
	return s.chunkService.GetChunkUpload(ctx, uploadID)
}

// CreateBook は新しい書籍を作成する
func (s *UploadService) CreateBook(ctx context.Context, userID uuid.UUID, metadata models.BookMetadata) (*models.Book, error) {
	// メタデータを検証（validator パッケージは後で import）
	// if err := validator.ValidateBookMetadata(metadata); err != nil {
	// 	return nil, fmt.Errorf("invalid metadata: %w", err)
	// }

	book := &models.Book{
		ID:                uuid.New(),
		UserID:            userID,
		Title:             metadata.Title,
		TargetLanguage:    metadata.TargetLanguage,
		NativeLanguage:    metadata.NativeLanguage,
		ReferenceLanguage: metadata.ReferenceLanguage,
		Status:            models.BookStatusUploading,
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
	}

	// 実際の実装ではデータベースに保存
	// repository.SaveBook(ctx, book)

	return book, nil
}

// UploadFile はファイルをアップロードする
func (s *UploadService) UploadFile(ctx context.Context, userID, bookID uuid.UUID, fileHeader *multipart.FileHeader, reader io.Reader) (*models.BookFile, error) {
	// ファイルを検証
	if err := file.ValidateFile(fileHeader); err != nil {
		return nil, fmt.Errorf("file validation failed: %w", err)
	}

	// ストレージに保存
	storagePath, err := s.storage.SaveFile(ctx, userID, bookID, fileHeader.Filename, reader)
	if err != nil {
		return nil, fmt.Errorf("failed to save file: %w", err)
	}

	// BookFileモデルを作成
	bookFile := &models.BookFile{
		ID:          uuid.New(),
		BookID:      bookID,
		FileName:    fileHeader.Filename,
		FileType:    getFileType(fileHeader.Filename),
		FileSize:    fileHeader.Size,
		StoragePath: storagePath,
		UploadedAt:  time.Now(),
	}

	// 実際の実装ではデータベースに保存
	// repository.SaveBookFile(ctx, bookFile)

	return bookFile, nil
}

// ValidateFiles は複数のファイルを検証する
func (s *UploadService) ValidateFiles(files []*multipart.FileHeader) error {
	if len(files) == 0 {
		return ErrNoFilesProvided
	}

	for _, fileHeader := range files {
		if err := file.ValidateFile(fileHeader); err != nil {
			return fmt.Errorf("file %s validation failed: %w", fileHeader.Filename, err)
		}
	}

	return nil
}

// UpdateUploadProgress はアップロード進捗を更新する
func (s *UploadService) UpdateUploadProgress(ctx context.Context, progress *models.UploadProgress) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.progress[progress.BookID] = progress

	// 実際の実装ではRedisなどに保存
	// redis.Set(ctx, fmt.Sprintf("upload:progress:%s", progress.BookID), progress, 24*time.Hour)

	return nil
}

// GetUploadProgress はアップロード進捗を取得する
func (s *UploadService) GetUploadProgress(ctx context.Context, bookID uuid.UUID) (*models.UploadProgress, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	progress, ok := s.progress[bookID]
	if !ok {
		return nil, ErrBookNotFound
	}

	return progress, nil
}

// UploadMultipleFiles は複数のファイルを一括でアップロードする
func (s *UploadService) UploadMultipleFiles(ctx context.Context, userID, bookID uuid.UUID, files []*multipart.FileHeader) ([]*models.BookFile, error) {
	// ファイルを検証
	if err := s.ValidateFiles(files); err != nil {
		return nil, err
	}

	// 進捗を初期化
	totalBytes := int64(0)
	for _, fileHeader := range files {
		totalBytes += fileHeader.Size
	}

	progress := &models.UploadProgress{
		BookID:        bookID,
		TotalFiles:    len(files),
		UploadedFiles: 0,
		TotalBytes:    totalBytes,
		UploadedBytes: 0,
		Status:        "uploading",
	}
	s.UpdateUploadProgress(ctx, progress)

	// ファイルを順次アップロード
	bookFiles := make([]*models.BookFile, 0, len(files))
	for i, fileHeader := range files {
		// ファイルを開く
		file, err := fileHeader.Open()
		if err != nil {
			progress.Status = "failed"
			progress.Message = fmt.Sprintf("failed to open file: %s", fileHeader.Filename)
			s.UpdateUploadProgress(ctx, progress)
			return nil, fmt.Errorf("failed to open file %s: %w", fileHeader.Filename, err)
		}
		defer file.Close()

		// アップロード
		bookFile, err := s.UploadFile(ctx, userID, bookID, fileHeader, file)
		if err != nil {
			progress.Status = "failed"
			progress.Message = fmt.Sprintf("failed to upload file: %s", fileHeader.Filename)
			s.UpdateUploadProgress(ctx, progress)
			return nil, err
		}

		bookFiles = append(bookFiles, bookFile)

		// 進捗を更新
		progress.UploadedFiles = i + 1
		progress.UploadedBytes += fileHeader.Size
		s.UpdateUploadProgress(ctx, progress)
	}

	// 完了
	progress.Status = "completed"
	s.UpdateUploadProgress(ctx, progress)

	return bookFiles, nil
}

// getFileType はファイル名から拡張子を取得する
func getFileType(fileName string) string {
	ext := file.GetFileExtension(fileName)
	if ext != "" && ext[0] == '.' {
		return ext[1:] // 先頭のドットを削除
	}
	return ext
}
