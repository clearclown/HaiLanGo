package storage

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/google/uuid"
)

// Storage はファイルストレージのインターフェース
type Storage interface {
	// SaveFile はファイルを保存する
	SaveFile(ctx context.Context, userID, bookID uuid.UUID, fileName string, reader io.Reader) (string, error)

	// GetFile はファイルを取得する
	GetFile(ctx context.Context, path string) (io.ReadCloser, error)

	// DeleteFile はファイルを削除する
	DeleteFile(ctx context.Context, path string) error

	// FileExists はファイルの存在をチェックする
	FileExists(ctx context.Context, path string) (bool, error)
}

// LocalStorage はローカルファイルシステムを使用したストレージ実装
type LocalStorage struct {
	basePath string
}

// NewLocalStorage はLocalStorageの新しいインスタンスを作成する
func NewLocalStorage(basePath string) *LocalStorage {
	return &LocalStorage{
		basePath: basePath,
	}
}

// SaveFile はファイルをローカルファイルシステムに保存する
func (s *LocalStorage) SaveFile(ctx context.Context, userID, bookID uuid.UUID, fileName string, reader io.Reader) (string, error) {
	// ストレージパスを生成
	relativePath := GenerateStoragePath(userID, bookID, fileName)
	fullPath := filepath.Join(s.basePath, relativePath)

	// ディレクトリを作成
	dir := filepath.Dir(fullPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", fmt.Errorf("failed to create directory: %w", err)
	}

	// ファイルを作成
	file, err := os.Create(fullPath)
	if err != nil {
		return "", fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	// データをコピー
	if _, err := io.Copy(file, reader); err != nil {
		return "", fmt.Errorf("failed to write file: %w", err)
	}

	return relativePath, nil
}

// GetFile はファイルをローカルファイルシステムから取得する
func (s *LocalStorage) GetFile(ctx context.Context, path string) (io.ReadCloser, error) {
	fullPath := filepath.Join(s.basePath, path)

	file, err := os.Open(fullPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("file not found: %w", err)
		}
		return nil, fmt.Errorf("failed to open file: %w", err)
	}

	return file, nil
}

// DeleteFile はファイルをローカルファイルシステムから削除する
func (s *LocalStorage) DeleteFile(ctx context.Context, path string) error {
	fullPath := filepath.Join(s.basePath, path)

	if err := os.Remove(fullPath); err != nil {
		if os.IsNotExist(err) {
			return nil // ファイルが存在しない場合は成功とみなす
		}
		return fmt.Errorf("failed to delete file: %w", err)
	}

	return nil
}

// FileExists はファイルの存在をチェックする
func (s *LocalStorage) FileExists(ctx context.Context, path string) (bool, error) {
	fullPath := filepath.Join(s.basePath, path)

	_, err := os.Stat(fullPath)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, fmt.Errorf("failed to check file existence: %w", err)
	}

	return true, nil
}

// GenerateStoragePath はユーザーIDと書籍IDからストレージパスを生成する
// パス形式: users/{userID}/books/{bookID}/{sanitized_filename}_{uuid}.{ext}
func GenerateStoragePath(userID, bookID uuid.UUID, fileName string) string {
	// ファイル名をサニタイズ
	sanitized := SanitizeFileName(fileName)

	// 拡張子を取得
	ext := filepath.Ext(sanitized)
	nameWithoutExt := strings.TrimSuffix(sanitized, ext)

	// ユニークなファイル名を生成
	uniqueID := uuid.New().String()[:8]
	uniqueFileName := fmt.Sprintf("%s_%s%s", nameWithoutExt, uniqueID, ext)

	// パスを生成
	return filepath.Join(
		"users",
		userID.String(),
		"books",
		bookID.String(),
		uniqueFileName,
	)
}

// SanitizeFileName はファイル名をサニタイズする
// - スペースをアンダースコアに変換
// - 危険な文字のみを削除（パス区切り文字や制御文字など）
func SanitizeFileName(fileName string) string {
	// スペースをアンダースコアに変換
	sanitized := strings.ReplaceAll(fileName, " ", "_")

	// 危険な文字を削除（パス区切り文字、制御文字など）
	// パス区切り文字（/, \）、制御文字、その他の危険な文字のみを削除
	re := regexp.MustCompile(`[<>:"/\\|?*\x00-\x1f#!]`)
	sanitized = re.ReplaceAllString(sanitized, "")

	return sanitized
}
