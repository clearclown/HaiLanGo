package tts

import (
	"context"
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMain(m *testing.M) {
	// テスト実行時は自動的にモックを使用
	os.Setenv("TEST_USE_MOCKS", "true")
	code := m.Run()
	os.Exit(code)
}

// ============================================================
// Factory Tests (TDD Red Phase - will fail until factory.go implemented)
// ============================================================

func TestNewTTSClient_WithMocks(t *testing.T) {
	os.Setenv("USE_MOCK_APIS", "true")
	defer os.Unsetenv("USE_MOCK_APIS")

	client, err := NewTTSClient()
	require.NoError(t, err)
	require.NotNil(t, client)

	// モッククライアントが返されることを確認
	_, ok := client.(*MockTTSClient)
	assert.True(t, ok, "expected MockTTSClient")
}

func TestNewTTSClient_WithoutAPIKey(t *testing.T) {
	os.Setenv("TTS_PROVIDER", "azure")
	os.Unsetenv("AZURE_SPEECH_KEY")
	os.Unsetenv("USE_MOCK_APIS")
	defer os.Unsetenv("TTS_PROVIDER")

	client, err := NewTTSClient()
	require.NoError(t, err)
	require.NotNil(t, client)

	// APIキーがない場合はモッククライアントが返されることを確認
	_, ok := client.(*MockTTSClient)
	assert.True(t, ok, "expected MockTTSClient when no API key is provided")
}

// ============================================================
// Mock Client Tests (TDD Red Phase - will fail until mock.go implemented)
// ============================================================

func TestMockTTSClient_Synthesize(t *testing.T) {
	ctx := context.Background()
	client := NewMockTTSClient()

	tests := []struct {
		name     string
		text     string
		language string
		voice    string
	}{
		{
			name:     "Russian text",
			text:     "Здравствуйте!",
			language: "ru",
			voice:    "ru-RU-DmitryNeural",
		},
		{
			name:     "Japanese text",
			text:     "こんにちは",
			language: "ja",
			voice:    "ja-JP-NanamiNeural",
		},
		{
			name:     "English text",
			text:     "Hello, world!",
			language: "en",
			voice:    "en-US-JennyNeural",
		},
		{
			name:     "Chinese text",
			text:     "你好",
			language: "zh",
			voice:    "zh-CN-XiaoxiaoNeural",
		},
		{
			name:     "Minor language - Kurdish",
			text:     "Silav",
			language: "ku",
			voice:    "ku-Arab-FahdNeural",
		},
		{
			name:     "Minor language - Amharic",
			text:     "ሰላም",
			language: "am",
			voice:    "am-ET-MekdesNeural",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			audioReader, err := client.Synthesize(ctx, tt.text, tt.language, tt.voice)
			require.NoError(t, err)
			require.NotNil(t, audioReader)
			defer audioReader.Close()

			// 音声データが読み取れることを確認
			audioData, err := io.ReadAll(audioReader)
			require.NoError(t, err)
			assert.NotEmpty(t, audioData, "audio data should not be empty")
		})
	}
}

func TestMockTTSClient_Synthesize_EmptyText(t *testing.T) {
	ctx := context.Background()
	client := NewMockTTSClient()

	// 空のテキストはエラーになるべき
	_, err := client.Synthesize(ctx, "", "en", "en-US-JennyNeural")
	assert.Error(t, err, "empty text should return error")
}

func TestMockTTSClient_GetVoices(t *testing.T) {
	ctx := context.Background()
	client := NewMockTTSClient()

	tests := []struct {
		name        string
		language    string
		minExpected int // 最低限期待する音声数
	}{
		{
			name:        "English voices",
			language:    "en",
			minExpected: 2,
		},
		{
			name:        "Japanese voices",
			language:    "ja",
			minExpected: 2,
		},
		{
			name:        "Russian voices",
			language:    "ru",
			minExpected: 1,
		},
		{
			name:        "All voices (empty language)",
			language:    "",
			minExpected: 5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			voices, err := client.GetVoices(ctx, tt.language)
			require.NoError(t, err)
			require.NotNil(t, voices)
			assert.GreaterOrEqual(t, len(voices), tt.minExpected)

			// 各音声の基本情報を検証
			for _, voice := range voices {
				assert.NotEmpty(t, voice.ID)
				assert.NotEmpty(t, voice.Name)
				assert.NotEmpty(t, voice.Language)
				assert.NotEmpty(t, voice.Gender)
			}
		})
	}
}

func TestMockTTSClient_GetSupportedLanguages(t *testing.T) {
	ctx := context.Background()
	client := NewMockTTSClient()

	languages, err := client.GetSupportedLanguages(ctx)
	require.NoError(t, err)
	require.NotNil(t, languages)

	// Azureは140言語をサポート（モックでは主要言語をシミュレート）
	assert.GreaterOrEqual(t, len(languages), 10, "should support at least 10 languages")

	// 主要言語が含まれることを確認
	majorLanguages := []string{"en", "ja", "zh", "ru", "es", "fr", "de", "pt", "it"}
	languageSet := make(map[string]bool)
	for _, lang := range languages {
		languageSet[lang] = true
	}

	for _, major := range majorLanguages {
		assert.True(t, languageSet[major], "should support major language: %s", major)
	}
}

func TestMockTTSClient_GetName(t *testing.T) {
	client := NewMockTTSClient()
	name := client.GetName()
	assert.Equal(t, "mock", name)
}

// ============================================================
// GetRecommendedProvider Tests
// ============================================================

func TestGetRecommendedProvider(t *testing.T) {
	tests := []struct {
		name         string
		languageCode string
		expected     TTSProvider
	}{
		{
			name:         "English - major language",
			languageCode: "en",
			expected:     ProviderAzureTTS,
		},
		{
			name:         "Japanese - major language",
			languageCode: "ja",
			expected:     ProviderAzureTTS,
		},
		{
			name:         "Russian - major language",
			languageCode: "ru",
			expected:     ProviderAzureTTS,
		},
		{
			name:         "Kurdish - minor language",
			languageCode: "ku",
			expected:     ProviderAzureTTS,
		},
		{
			name:         "Amharic - minor language",
			languageCode: "am",
			expected:     ProviderAzureTTS,
		},
		{
			name:         "Tibetan - minor language",
			languageCode: "bo",
			expected:     ProviderAzureTTS,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			provider := GetRecommendedProvider(tt.languageCode)
			assert.Equal(t, tt.expected, provider)
		})
	}
}

// ============================================================
// Struct Validation Tests
// ============================================================

func TestVoiceStruct(t *testing.T) {
	voice := Voice{
		ID:          "en-US-JennyNeural",
		Name:        "Jenny",
		Language:    "en-US",
		Gender:      "female",
		Type:        "neural",
		SampleRate:  24000,
		Description: "American English female voice",
		Styles:      []string{"cheerful", "sad", "angry"},
	}

	assert.Equal(t, "en-US-JennyNeural", voice.ID)
	assert.Equal(t, "Jenny", voice.Name)
	assert.Equal(t, "en-US", voice.Language)
	assert.Equal(t, "female", voice.Gender)
	assert.Equal(t, "neural", voice.Type)
	assert.Equal(t, 24000, voice.SampleRate)
	assert.Len(t, voice.Styles, 3)
}

func TestSynthesizeOptions(t *testing.T) {
	opts := SynthesizeOptions{
		Speed:      1.0,
		Pitch:      0.0,
		VolumeGain: 0.0,
		Style:      "cheerful",
		Format:     "mp3",
	}

	assert.Equal(t, 1.0, opts.Speed)
	assert.Equal(t, 0.0, opts.Pitch)
	assert.Equal(t, 0.0, opts.VolumeGain)
	assert.Equal(t, "cheerful", opts.Style)
	assert.Equal(t, "mp3", opts.Format)
}

func TestTTSResult(t *testing.T) {
	result := TTSResult{
		Audio:     []byte{0x49, 0x44, 0x33}, // MP3 header bytes
		Format:    "mp3",
		Duration:  1500,
		CharCount: 10,
		Provider:  "azure",
		Voice:     "en-US-JennyNeural",
		CacheKey:  "abc123",
		IsCached:  false,
	}

	assert.NotEmpty(t, result.Audio)
	assert.Equal(t, "mp3", result.Format)
	assert.Equal(t, int64(1500), result.Duration)
	assert.Equal(t, 10, result.CharCount)
	assert.Equal(t, "azure", result.Provider)
	assert.Equal(t, "en-US-JennyNeural", result.Voice)
	assert.Equal(t, "abc123", result.CacheKey)
	assert.False(t, result.IsCached)
}

// ============================================================
// Mock SetMockAudio Test
// ============================================================

func TestMockTTSClient_SetMockAudio(t *testing.T) {
	client := NewMockTTSClient()

	// カスタム音声データを設定
	customAudio := []byte{0x49, 0x44, 0x33, 0x00, 0x00, 0x00}
	err := client.SetMockAudio(customAudio)
	require.NoError(t, err)

	// 設定した音声データが返されることを確認
	ctx := context.Background()
	audioReader, err := client.Synthesize(ctx, "Test", "en", "en-US-JennyNeural")
	require.NoError(t, err)
	defer audioReader.Close()

	audioData, err := io.ReadAll(audioReader)
	require.NoError(t, err)
	assert.Equal(t, customAudio, audioData)
}

// ============================================================
// Legacy GoogleTTSClient Tests (Backward Compatibility)
// ============================================================

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
