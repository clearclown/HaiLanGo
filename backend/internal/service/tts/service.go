package tts

import (
	"context"
	"errors"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/clearclown/HaiLanGo/backend/internal/service/cache"
	"github.com/clearclown/HaiLanGo/backend/pkg/storage"
	"github.com/clearclown/HaiLanGo/backend/pkg/tts"
)

// TTSService はTTS音声生成サービス
type TTSService struct {
	ttsClient tts.TTSClient
	cache     cache.AudioCache
	storage   storage.AudioStorage
	useMock   bool
}

// NewTTSService は新しいTTSサービスを作成
func NewTTSService() *TTSService {
	useMock := os.Getenv("USE_MOCK_APIS") == "true" ||
		os.Getenv("TEST_USE_MOCKS") == "true"

	apiKey := os.Getenv("GOOGLE_CLOUD_TTS_API_KEY")
	if apiKey == "" && !useMock {
		// APIキーがない場合は自動的にモックを使用
		useMock = true
	}

	return &TTSService{
		ttsClient: tts.NewGoogleTTSClient(apiKey),
		cache:     cache.NewAudioCache(),
		storage:   storage.NewAudioStorage(),
		useMock:   useMock,
	}
}

// GenerateAudio はテキストから音声を生成してURLを返す
func (s *TTSService) GenerateAudio(ctx context.Context, text string, lang string, quality string, speed float64) (string, error) {
	// バリデーション
	if err := s.validate(text, lang, quality, speed); err != nil {
		return "", err
	}

	// 1. キャッシュチェック
	cacheKey := s.cache.GenerateKey(text, lang, quality, speed)
	if audioURL, found := s.cache.Get(cacheKey); found {
		return audioURL, nil
	}

	// 2. TTS API呼び出し
	audioData, err := s.ttsClient.Generate(ctx, text, lang, quality, speed)
	if err != nil {
		return "", fmt.Errorf("failed to generate audio: %w", err)
	}

	// 3. ストレージに保存
	filename := s.storage.GenerateFilename(text, lang, quality, speed)
	audioURL, err := s.storage.Save(audioData, filename)
	if err != nil {
		return "", fmt.Errorf("failed to save audio: %w", err)
	}

	// 4. キャッシュに保存（7日間）
	ttl := 7 * 24 * time.Hour
	if err := s.cache.Set(cacheKey, audioURL, ttl); err != nil {
		// キャッシュ保存失敗はエラーとしない（ログのみ）
		fmt.Printf("Warning: failed to cache audio URL: %v\n", err)
	}

	return audioURL, nil
}

// BatchGenerate は複数のテキストから音声を一括生成
func (s *TTSService) BatchGenerate(ctx context.Context, texts []string, lang string, quality string, speed float64) ([]string, error) {
	if len(texts) == 0 {
		return nil, errors.New("texts cannot be empty")
	}

	audioURLs := make([]string, len(texts))
	var wg sync.WaitGroup
	errChan := make(chan error, len(texts))

	for i, text := range texts {
		wg.Add(1)
		go func(index int, txt string) {
			defer wg.Done()

			url, err := s.GenerateAudio(ctx, txt, lang, quality, speed)
			if err != nil {
				errChan <- err
				return
			}
			audioURLs[index] = url
		}(i, text)
	}

	wg.Wait()
	close(errChan)

	// エラーチェック
	if len(errChan) > 0 {
		return nil, <-errChan
	}

	return audioURLs, nil
}

// validate は入力パラメータの検証
func (s *TTSService) validate(text string, lang string, quality string, speed float64) error {
	if text == "" {
		return errors.New("text cannot be empty")
	}

	if speed < 0.5 || speed > 2.0 {
		return fmt.Errorf("speed must be between 0.5 and 2.0, got %.2f", speed)
	}

	if quality != "standard" && quality != "premium" {
		return fmt.Errorf("quality must be 'standard' or 'premium', got '%s'", quality)
	}

	return nil
}

// SupportedLanguages は対応言語のリストを返す
func (s *TTSService) SupportedLanguages() []string {
	return []string{
		"ja", // 日本語
		"zh", // 中国語
		"en", // 英語
		"ru", // ロシア語
		"fa", // ペルシャ語
		"he", // ヘブライ語
		"es", // スペイン語
		"fr", // フランス語
		"pt", // ポルトガル語
		"de", // ドイツ語
		"it", // イタリア語
		"tr", // トルコ語
	}
}
