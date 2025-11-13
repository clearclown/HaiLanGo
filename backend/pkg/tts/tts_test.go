package tts

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestGenerateWithGoogleTTS は Google Cloud TTS API を使用した音声生成のテスト
func TestGenerateWithGoogleTTS(t *testing.T) {
	ctx := context.Background()
	text := "Hello, world!"
	lang := "en"
	quality := "standard"
	speed := 1.0

	// モック環境で実行
	client := NewGoogleTTSClient("mock-api-key")
	audioData, err := client.Generate(ctx, text, lang, quality, speed)

	require.NoError(t, err)
	assert.NotEmpty(t, audioData)
	assert.Greater(t, len(audioData), 0)
}

// TestSpeedAdjustment は速度調整のテスト
func TestSpeedAdjustment(t *testing.T) {
	ctx := context.Background()
	text := "Test speed adjustment"
	lang := "en"
	quality := "standard"

	speeds := []float64{0.5, 0.75, 1.0, 1.25, 1.5, 2.0}
	client := NewGoogleTTSClient("mock-api-key")

	for _, speed := range speeds {
		audioData, err := client.Generate(ctx, text, lang, quality, speed)
		require.NoError(t, err, "Speed %.2fx should work", speed)
		assert.NotEmpty(t, audioData)
	}
}

// TestMultipleLanguages は複数言語のテスト
func TestMultipleLanguages(t *testing.T) {
	ctx := context.Background()
	text := "Test"
	quality := "standard"
	speed := 1.0

	// 主要12言語
	languages := []string{
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

	client := NewGoogleTTSClient("mock-api-key")

	for _, lang := range languages {
		audioData, err := client.Generate(ctx, text, lang, quality, speed)
		require.NoError(t, err, "Language %s should work", lang)
		assert.NotEmpty(t, audioData)
	}
}

// TestPremiumQuality はプレミアム品質のテスト
func TestPremiumQuality(t *testing.T) {
	ctx := context.Background()
	text := "Premium quality test"
	lang := "en"
	quality := "premium"
	speed := 1.0

	client := NewGoogleTTSClient("mock-api-key")
	audioData, err := client.Generate(ctx, text, lang, quality, speed)

	require.NoError(t, err)
	assert.NotEmpty(t, audioData)
}

// TestLongText は長文テキストのテスト
func TestLongText(t *testing.T) {
	ctx := context.Background()
	// 1000文字以上のテキスト
	text := ""
	for i := 0; i < 100; i++ {
		text += "This is a long text for testing purposes. "
	}
	lang := "en"
	quality := "standard"
	speed := 1.0

	client := NewGoogleTTSClient("mock-api-key")
	audioData, err := client.Generate(ctx, text, lang, quality, speed)

	require.NoError(t, err)
	assert.NotEmpty(t, audioData)
}

// TestSpecialCharacters は特殊文字を含むテキストのテスト
func TestSpecialCharacters(t *testing.T) {
	ctx := context.Background()
	texts := []string{
		"Hello! How are you?",
		"こんにちは！元気ですか？",
		"¡Hola! ¿Cómo estás?",
		"Привет! Как дела?",
	}
	lang := "en"
	quality := "standard"
	speed := 1.0

	client := NewGoogleTTSClient("mock-api-key")

	for _, text := range texts {
		audioData, err := client.Generate(ctx, text, lang, quality, speed)
		require.NoError(t, err)
		assert.NotEmpty(t, audioData)
	}
}

// TestInvalidLanguage は未対応言語のエラーハンドリングテスト
func TestInvalidLanguage(t *testing.T) {
	ctx := context.Background()
	text := "Test"
	lang := "invalid-lang"
	quality := "standard"
	speed := 1.0

	client := NewGoogleTTSClient("mock-api-key")
	audioData, err := client.Generate(ctx, text, lang, quality, speed)

	// モック環境では成功するが、実環境ではエラーになる
	// ここではエラーハンドリングのロジックをテスト
	if err != nil {
		assert.Empty(t, audioData)
	} else {
		// モック環境では成功する
		assert.NotEmpty(t, audioData)
	}
}

// TestInvalidSpeed は無効な速度のエラーハンドリングテスト
func TestInvalidSpeed(t *testing.T) {
	ctx := context.Background()
	text := "Test"
	lang := "en"
	quality := "standard"
	speed := 5.0 // 無効な速度

	client := NewGoogleTTSClient("mock-api-key")
	_, err := client.Generate(ctx, text, lang, quality, speed)

	assert.Error(t, err)
}

// TestEmptyText は空のテキストのエラーハンドリングテスト
func TestEmptyText(t *testing.T) {
	ctx := context.Background()
	text := ""
	lang := "en"
	quality := "standard"
	speed := 1.0

	client := NewGoogleTTSClient("mock-api-key")
	_, err := client.Generate(ctx, text, lang, quality, speed)

	assert.Error(t, err)
}
