package storage

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

// AudioStorage は音声ストレージのインターフェース
type AudioStorage interface {
	Save(audioData []byte, filename string) (string, error)
	Get(url string) ([]byte, error)
	Delete(url string) error
	GenerateFilename(text string, lang string, quality string, speed float64) string
}

// LocalAudioStorage はローカルファイルシステム音声ストレージ
type LocalAudioStorage struct {
	basePath string
	baseURL  string
	mu       sync.RWMutex
	data     map[string][]byte // モック用インメモリストレージ
	useMock  bool
}

// NewAudioStorage は新しい音声ストレージを作成
func NewAudioStorage() AudioStorage {
	useMock := os.Getenv("USE_MOCK_APIS") == "true" ||
		os.Getenv("TEST_USE_MOCKS") == "true"

	basePath := os.Getenv("AUDIO_STORAGE_PATH")
	if basePath == "" {
		basePath = "./storage/audio"
	}

	baseURL := os.Getenv("AUDIO_BASE_URL")
	if baseURL == "" {
		baseURL = "http://localhost:8080/audio"
	}

	storage := &LocalAudioStorage{
		basePath: basePath,
		baseURL:  baseURL,
		data:     make(map[string][]byte),
		useMock:  useMock,
	}

	// ディレクトリを作成
	if !useMock {
		// ディレクトリ作成のエラーは無視（Save時に再度作成を試みる）
		_ = os.MkdirAll(basePath, 0755)
	}

	return storage
}

// Save は音声データを保存
func (s *LocalAudioStorage) Save(audioData []byte, filename string) (string, error) {
	if len(audioData) == 0 {
		return "", errors.New("audio data is empty")
	}

	if s.useMock {
		// モック環境：インメモリに保存
		s.mu.Lock()
		defer s.mu.Unlock()

		url := fmt.Sprintf("%s/%s", s.baseURL, filename)
		s.data[url] = audioData
		return url, nil
	}

	// 実環境：ファイルシステムに保存
	filePath := filepath.Join(s.basePath, filename)

	// ディレクトリを作成
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", fmt.Errorf("failed to create directory: %w", err)
	}

	// ファイルに書き込み
	if err := os.WriteFile(filePath, audioData, 0644); err != nil {
		return "", fmt.Errorf("failed to write file: %w", err)
	}

	url := fmt.Sprintf("%s/%s", s.baseURL, filename)
	return url, nil
}

// Get は音声データを取得
func (s *LocalAudioStorage) Get(url string) ([]byte, error) {
	if s.useMock {
		// モック環境：インメモリから取得
		s.mu.RLock()
		defer s.mu.RUnlock()

		data, found := s.data[url]
		if !found {
			return nil, errors.New("audio not found")
		}
		return data, nil
	}

	// 実環境：ファイルシステムから取得
	// URLからファイル名を抽出
	filename := filepath.Base(url)
	filePath := filepath.Join(s.basePath, filename)

	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	return data, nil
}

// Delete は音声データを削除
func (s *LocalAudioStorage) Delete(url string) error {
	if s.useMock {
		// モック環境：インメモリから削除
		s.mu.Lock()
		defer s.mu.Unlock()

		delete(s.data, url)
		return nil
	}

	// 実環境：ファイルシステムから削除
	filename := filepath.Base(url)
	filePath := filepath.Join(s.basePath, filename)

	if err := os.Remove(filePath); err != nil {
		return fmt.Errorf("failed to delete file: %w", err)
	}

	return nil
}

// GenerateFilename はファイル名を生成
func (s *LocalAudioStorage) GenerateFilename(text string, lang string, quality string, speed float64) string {
	// ハッシュを使用してファイル名を生成
	data := fmt.Sprintf("%s:%s:%s:%.2f", text, lang, quality, speed)
	hash := sha256.Sum256([]byte(data))
	hashStr := hex.EncodeToString(hash[:])

	return fmt.Sprintf("%s.mp3", hashStr[:16])
}
