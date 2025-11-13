package service

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/clearclown/HaiLanGo/backend/internal/models"
	"github.com/clearclown/HaiLanGo/backend/pkg/storage"
	"github.com/google/uuid"
)

var (
	// ErrChunkUploadNotFound はチャンクアップロードが見つからないエラー
	ErrChunkUploadNotFound = errors.New("chunk upload not found")
	// ErrInvalidChunkNumber は無効なチャンク番号のエラー
	ErrInvalidChunkNumber = errors.New("invalid chunk number")
	// ErrChunkAlreadyCompleted はチャンクアップロードがすでに完了しているエラー
	ErrChunkAlreadyCompleted = errors.New("chunk upload already completed")
)

const (
	// DefaultChunkSize はデフォルトのチャンクサイズ（5MB）
	DefaultChunkSize = 5 * 1024 * 1024
	// MaxChunkSize は最大チャンクサイズ（10MB）
	MaxChunkSize = 10 * 1024 * 1024
)

// ChunkUploadService はチャンクアップロードのビジネスロジックを提供する
type ChunkUploadService struct {
	storage       storage.Storage
	tempDir       string
	chunkUploads  map[uuid.UUID]*models.ChunkUpload // 実際の実装ではRedisやDBを使用
}

// NewChunkUploadService はChunkUploadServiceの新しいインスタンスを作成する
func NewChunkUploadService(storage storage.Storage, tempDir string) *ChunkUploadService {
	return &ChunkUploadService{
		storage:      storage,
		tempDir:      tempDir,
		chunkUploads: make(map[uuid.UUID]*models.ChunkUpload),
	}
}

// InitiateChunkUpload はチャンクアップロードを開始する
func (s *ChunkUploadService) InitiateChunkUpload(ctx context.Context, bookID uuid.UUID, fileName string, totalChunks int, fileSize int64) (*models.ChunkUpload, error) {
	if totalChunks <= 0 {
		return nil, fmt.Errorf("total chunks must be greater than 0")
	}

	chunkUpload := &models.ChunkUpload{
		ID:             uuid.New(),
		BookID:         bookID,
		FileName:       fileName,
		TotalChunks:    totalChunks,
		UploadedChunks: 0,
		Status:         "pending",
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	// チャンクアップロード情報を保存
	s.chunkUploads[chunkUpload.ID] = chunkUpload

	// 一時ディレクトリを作成
	chunkDir := s.getChunkDir(chunkUpload.ID)
	if err := os.MkdirAll(chunkDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create chunk directory: %w", err)
	}

	return chunkUpload, nil
}

// UploadChunk はチャンクをアップロードする
func (s *ChunkUploadService) UploadChunk(ctx context.Context, uploadID uuid.UUID, chunkNumber int, data io.Reader) error {
	// チャンクアップロード情報を取得
	chunkUpload, ok := s.chunkUploads[uploadID]
	if !ok {
		return ErrChunkUploadNotFound
	}

	// チャンク番号の検証
	if chunkNumber < 0 || chunkNumber >= chunkUpload.TotalChunks {
		return ErrInvalidChunkNumber
	}

	// ステータスチェック
	if chunkUpload.Status == "completed" {
		return ErrChunkAlreadyCompleted
	}

	// チャンクファイルを保存
	chunkPath := s.getChunkPath(uploadID, chunkNumber)
	file, err := os.Create(chunkPath)
	if err != nil {
		return fmt.Errorf("failed to create chunk file: %w", err)
	}
	defer file.Close()

	if _, err := io.Copy(file, data); err != nil {
		return fmt.Errorf("failed to write chunk: %w", err)
	}

	// アップロード済みチャンク数をインクリメント
	chunkUpload.UploadedChunks++
	chunkUpload.UpdatedAt = time.Now()
	chunkUpload.Status = "uploading"

	// すべてのチャンクがアップロードされたか確認
	if chunkUpload.UploadedChunks == chunkUpload.TotalChunks {
		// チャンクを結合
		if err := s.mergeChunks(ctx, chunkUpload); err != nil {
			chunkUpload.Status = "failed"
			return fmt.Errorf("failed to merge chunks: %w", err)
		}
		chunkUpload.Status = "completed"
	}

	return nil
}

// GetChunkUpload はチャンクアップロード情報を取得する
func (s *ChunkUploadService) GetChunkUpload(ctx context.Context, uploadID uuid.UUID) (*models.ChunkUpload, error) {
	chunkUpload, ok := s.chunkUploads[uploadID]
	if !ok {
		return nil, ErrChunkUploadNotFound
	}
	return chunkUpload, nil
}

// mergeChunks はチャンクを結合して最終ファイルを作成する
func (s *ChunkUploadService) mergeChunks(ctx context.Context, chunkUpload *models.ChunkUpload) error {
	// 結合ファイルを作成
	mergedPath := s.getMergedPath(chunkUpload.ID)
	mergedFile, err := os.Create(mergedPath)
	if err != nil {
		return fmt.Errorf("failed to create merged file: %w", err)
	}
	defer mergedFile.Close()

	// チャンクを順番に結合
	for i := 0; i < chunkUpload.TotalChunks; i++ {
		chunkPath := s.getChunkPath(chunkUpload.ID, i)
		chunkFile, err := os.Open(chunkPath)
		if err != nil {
			return fmt.Errorf("failed to open chunk %d: %w", i, err)
		}

		if _, err := io.Copy(mergedFile, chunkFile); err != nil {
			chunkFile.Close()
			return fmt.Errorf("failed to copy chunk %d: %w", i, err)
		}
		chunkFile.Close()
	}

	// 結合したファイルをストレージに保存
	// TODO: userIDを適切に取得する
	userID := uuid.New()
	mergedFile.Seek(0, 0)

	_, err = s.storage.SaveFile(ctx, userID, chunkUpload.BookID, chunkUpload.FileName, mergedFile)
	if err != nil {
		return fmt.Errorf("failed to save merged file to storage: %w", err)
	}

	// 一時ファイルをクリーンアップ
	s.cleanupChunks(chunkUpload.ID)

	return nil
}

// cleanupChunks はチャンクファイルと一時ディレクトリをクリーンアップする
func (s *ChunkUploadService) cleanupChunks(uploadID uuid.UUID) {
	chunkDir := s.getChunkDir(uploadID)
	os.RemoveAll(chunkDir)
}

// getChunkDir はチャンク保存ディレクトリを取得する
func (s *ChunkUploadService) getChunkDir(uploadID uuid.UUID) string {
	return filepath.Join(s.tempDir, "chunks", uploadID.String())
}

// getChunkPath はチャンクファイルのパスを取得する
func (s *ChunkUploadService) getChunkPath(uploadID uuid.UUID, chunkNumber int) string {
	return filepath.Join(s.getChunkDir(uploadID), fmt.Sprintf("chunk_%d", chunkNumber))
}

// getMergedPath は結合ファイルのパスを取得する
func (s *ChunkUploadService) getMergedPath(uploadID uuid.UUID) string {
	return filepath.Join(s.getChunkDir(uploadID), "merged")
}

// CalculateChecksum はデータのMD5チェックサムを計算する
func CalculateChecksum(data io.Reader) (string, error) {
	hash := md5.New()
	if _, err := io.Copy(hash, data); err != nil {
		return "", err
	}
	return hex.EncodeToString(hash.Sum(nil)), nil
}
