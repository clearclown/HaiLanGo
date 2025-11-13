package tts

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
)

// TTSClient はTTS APIクライアントのインターフェース
type TTSClient interface {
	Generate(ctx context.Context, text string, lang string, quality string, speed float64) ([]byte, error)
}

// GoogleTTSClient はGoogle Cloud TTSクライアント
type GoogleTTSClient struct {
	apiKey  string
	useMock bool
}

// NewGoogleTTSClient は新しいGoogle Cloud TTSクライアントを作成
func NewGoogleTTSClient(apiKey string) *GoogleTTSClient {
	useMock := os.Getenv("USE_MOCK_APIS") == "true" ||
		os.Getenv("TEST_USE_MOCKS") == "true" ||
		apiKey == "" ||
		apiKey == "mock-api-key"

	return &GoogleTTSClient{
		apiKey:  apiKey,
		useMock: useMock,
	}
}

// Generate はテキストから音声データを生成
func (c *GoogleTTSClient) Generate(ctx context.Context, text string, lang string, quality string, speed float64) ([]byte, error) {
	// バリデーション
	if err := c.validate(text, lang, quality, speed); err != nil {
		return nil, err
	}

	// モック環境の場合
	if c.useMock {
		return c.generateMock(text, lang, quality, speed)
	}

	// 実際のGoogle Cloud TTS API呼び出し
	return c.generateReal(ctx, text, lang, quality, speed)
}

// validate は入力パラメータの検証
func (c *GoogleTTSClient) validate(text string, lang string, quality string, speed float64) error {
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

// generateMock はモック音声データを生成
func (c *GoogleTTSClient) generateMock(text string, lang string, quality string, speed float64) ([]byte, error) {
	// モックデータの生成（実際のMP3データの代わりに疑似データを返す）
	// ハッシュを使用して決定論的なデータを生成
	hash := sha256.Sum256([]byte(fmt.Sprintf("%s:%s:%s:%.2f", text, lang, quality, speed)))
	hashStr := hex.EncodeToString(hash[:])

	// 疑似音声データ（実際にはMP3ヘッダーとダミーデータ）
	mockData := []byte(fmt.Sprintf("MOCK_AUDIO_DATA:%s:text=%s:lang=%s:quality=%s:speed=%.2f",
		hashStr[:16], text, lang, quality, speed))

	return mockData, nil
}

// generateReal は実際のGoogle Cloud TTS APIを呼び出し
func (c *GoogleTTSClient) generateReal(ctx context.Context, text string, lang string, quality string, speed float64) ([]byte, error) {
	// TODO: 実際のGoogle Cloud TTS API実装
	// 現在はモックを返す
	return c.generateMock(text, lang, quality, speed)
}

// SupportedLanguages は対応言語のリストを返す
func (c *GoogleTTSClient) SupportedLanguages() []string {
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

// IsLanguageSupported は言語がサポートされているかチェック
func (c *GoogleTTSClient) IsLanguageSupported(lang string) bool {
	for _, supported := range c.SupportedLanguages() {
		if supported == lang {
			return true
		}
	}
	return false
}
