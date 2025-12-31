package stt

import (
	"bytes"
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestGoogleSTTClient はGoogle Cloud STTクライアントのテスト
func TestGoogleSTTClient(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name        string
		audioData   []byte
		lang        string
		expectError bool
	}{
		{
			name:        "正常な認識",
			audioData:   []byte("test audio"),
			lang:        "en-US",
			expectError: false,
		},
		{
			name:        "空の音声データ",
			audioData:   []byte{},
			lang:        "en-US",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := NewGoogleSTTClient("") // Empty API key will use mock
			result, err := client.Recognize(ctx, tt.audioData, tt.lang)

			if tt.expectError {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.NotNil(t, result)
			assert.NotEmpty(t, result.Text)
		})
	}
}

// TestWhisperClient はWhisper APIクライアントのテスト
func TestWhisperClient(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name        string
		audioData   []byte
		lang        string
		expectError bool
	}{
		{
			name:        "正常な認識",
			audioData:   []byte("test audio"),
			lang:        "en",
			expectError: false,
		},
		{
			name:        "空の音声データ",
			audioData:   []byte{},
			lang:        "en",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := NewWhisperClient()
			result, err := client.Recognize(ctx, tt.audioData, tt.lang)

			if tt.expectError {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.NotNil(t, result)
			assert.NotEmpty(t, result.Text)
		})
	}
}

// TestMockSTTClient はモックSTTクライアントのテスト
func TestMockSTTClient(t *testing.T) {
	ctx := context.Background()

	client := NewMockSTTClient()
	audioData := []byte("test audio")
	lang := "en"

	result, err := client.Recognize(ctx, audioData, lang)

	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.NotEmpty(t, result.Text)
	assert.Equal(t, lang, result.Language)
}

// TestLegacySTTClientFactory はレガシーSTTクライアントファクトリーのテスト（後方互換性）
func TestLegacySTTClientFactory(t *testing.T) {
	tests := []struct {
		name         string
		useMock      bool
		apiKey       string
		expectType   string
	}{
		{
			name:       "モックを使用",
			useMock:    true,
			apiKey:     "",
			expectType: "mock",
		},
		{
			name:       "Google STTを使用",
			useMock:    false,
			apiKey:     "test-api-key",
			expectType: "google",
		},
		{
			name:       "APIキーなしでモックを使用",
			useMock:    false,
			apiKey:     "",
			expectType: "mock",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := NewLegacySTTClient(tt.useMock, tt.apiKey)
			assert.NotNil(t, client)
		})
	}
}

// ============================================================
// 新インターフェース用テスト
// ============================================================

// TestNewSTTClient_WithMocks はモック環境でのSTTクライアント生成テスト
func TestNewSTTClient_WithMocks(t *testing.T) {
	// 環境変数を設定
	os.Setenv("USE_MOCK_APIS", "true")
	defer os.Unsetenv("USE_MOCK_APIS")

	client, err := NewSTTClient()

	require.NoError(t, err)
	assert.NotNil(t, client)
	assert.Equal(t, "mock", client.GetName())
}

// TestNewSTTClient_WithoutAPIKey はAPIキーなしでモックが使用されることをテスト
func TestNewSTTClient_WithoutAPIKey(t *testing.T) {
	// 環境変数をクリア
	os.Unsetenv("USE_MOCK_APIS")
	os.Unsetenv("TEST_USE_MOCKS")
	os.Unsetenv("OPENAI_API_KEY")
	os.Unsetenv("STT_PROVIDER")

	client, err := NewSTTClient()

	require.NoError(t, err)
	assert.NotNil(t, client)
	// APIキーがないのでモックが使用される
	assert.Equal(t, "mock", client.GetName())
}

// TestNewSTTClient_WithProvider はプロバイダー指定のテスト
func TestNewSTTClient_WithProvider(t *testing.T) {
	tests := []struct {
		name       string
		provider   string
		expectName string
	}{
		{
			name:       "Whisperプロバイダー（デフォルト）",
			provider:   "whisper",
			expectName: "mock", // APIキーなしでモック
		},
		{
			name:       "Azureプロバイダー",
			provider:   "azure",
			expectName: "mock", // APIキーなしでモック
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Unsetenv("USE_MOCK_APIS")
			os.Unsetenv("OPENAI_API_KEY")
			os.Unsetenv("AZURE_SPEECH_KEY")
			os.Setenv("STT_PROVIDER", tt.provider)
			defer os.Unsetenv("STT_PROVIDER")

			client, err := NewSTTClient()

			require.NoError(t, err)
			assert.NotNil(t, client)
			assert.Equal(t, tt.expectName, client.GetName())
		})
	}
}

// TestMockSTTClientNew_Transcribe は新モックSTTクライアントのTranscribeテスト
func TestMockSTTClientNew_Transcribe(t *testing.T) {
	ctx := context.Background()
	client := NewMockSTTClientNew()

	tests := []struct {
		name        string
		audio       []byte
		language    string
		expectText  string
		expectError bool
	}{
		{
			name:        "英語の認識",
			audio:       []byte("test audio"),
			language:    "en",
			expectText:  "Hello, world!",
			expectError: false,
		},
		{
			name:        "日本語の認識",
			audio:       []byte("test audio"),
			language:    "ja",
			expectText:  "こんにちは",
			expectError: false,
		},
		{
			name:        "ロシア語の認識",
			audio:       []byte("test audio"),
			language:    "ru",
			expectText:  "Здравствуйте",
			expectError: false,
		},
		{
			name:        "クルド語の認識（マイナー言語）",
			audio:       []byte("test audio"),
			language:    "ku",
			expectText:  "Silav",
			expectError: false,
		},
		{
			name:        "空の音声データ",
			audio:       []byte{},
			language:    "en",
			expectText:  "",
			expectError: true,
		},
		{
			name:        "nil音声データ",
			audio:       nil,
			language:    "en",
			expectText:  "",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var reader *bytes.Reader
			if tt.audio != nil {
				reader = bytes.NewReader(tt.audio)
			}

			var result *TranscriptionResult
			var err error
			if reader == nil {
				result, err = client.Transcribe(ctx, nil, tt.language)
			} else {
				result, err = client.Transcribe(ctx, reader, tt.language)
			}

			if tt.expectError {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.NotNil(t, result)
			assert.Equal(t, tt.expectText, result.Text)
			assert.Equal(t, tt.language, result.Language)
			assert.Equal(t, "mock", result.Provider)
			assert.Greater(t, result.Confidence, 0.0)
		})
	}
}

// TestMockSTTClientNew_GetSupportedLanguages はサポート言語取得テスト
func TestMockSTTClientNew_GetSupportedLanguages(t *testing.T) {
	ctx := context.Background()
	client := NewMockSTTClientNew()

	languages, err := client.GetSupportedLanguages(ctx)

	require.NoError(t, err)
	assert.NotEmpty(t, languages)
	// Whisperは99言語対応
	assert.GreaterOrEqual(t, len(languages), 90)
	// 主要言語が含まれていること
	assert.Contains(t, languages, "en")
	assert.Contains(t, languages, "ja")
	assert.Contains(t, languages, "zh")
	assert.Contains(t, languages, "ru")
	// マイナー言語も含まれていること
	assert.Contains(t, languages, "am") // アムハラ語
	assert.Contains(t, languages, "bo") // チベット語
}

// TestNewPronunciationEvaluator_WithMocks はモック環境での発音評価クライアント生成テスト
func TestNewPronunciationEvaluator_WithMocks(t *testing.T) {
	os.Setenv("USE_MOCK_APIS", "true")
	defer os.Unsetenv("USE_MOCK_APIS")

	evaluator, err := NewPronunciationEvaluator()

	require.NoError(t, err)
	assert.NotNil(t, evaluator)
}

// TestMockPronunciationEvaluator_EvaluatePronunciation は発音評価テスト
func TestMockPronunciationEvaluator_EvaluatePronunciation(t *testing.T) {
	ctx := context.Background()
	evaluator := NewMockPronunciationEvaluator()

	tests := []struct {
		name         string
		audio        []byte
		expectedText string
		language     string
		expectError  bool
	}{
		{
			name:         "正常な評価",
			audio:        []byte("test audio"),
			expectedText: "Hello world",
			language:     "en",
			expectError:  false,
		},
		{
			name:         "日本語の評価",
			audio:        []byte("test audio"),
			expectedText: "こんにちは",
			language:     "ja",
			expectError:  false,
		},
		{
			name:         "空の音声データ",
			audio:        []byte{},
			expectedText: "Hello",
			language:     "en",
			expectError:  true,
		},
		{
			name:         "空の期待テキスト",
			audio:        []byte("test audio"),
			expectedText: "",
			language:     "en",
			expectError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := bytes.NewReader(tt.audio)
			result, err := evaluator.EvaluatePronunciation(ctx, reader, tt.expectedText, tt.language)

			if tt.expectError {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.NotNil(t, result)
			assert.Greater(t, result.OverallScore, 0.0)
			assert.LessOrEqual(t, result.OverallScore, 100.0)
			assert.Equal(t, tt.expectedText, result.ExpectedText)
			assert.NotEmpty(t, result.Feedback)
			assert.NotEmpty(t, result.WordScores)
		})
	}
}

// TestGetRecommendedSTTProvider は推奨プロバイダー取得テスト
func TestGetRecommendedSTTProvider(t *testing.T) {
	tests := []struct {
		name         string
		languageCode string
		expected     STTProvider
	}{
		{
			name:         "英語はAzure推奨",
			languageCode: "en",
			expected:     ProviderAzureSpeech,
		},
		{
			name:         "日本語はAzure推奨",
			languageCode: "ja",
			expected:     ProviderAzureSpeech,
		},
		{
			name:         "クルド語はWhisper推奨",
			languageCode: "ku",
			expected:     ProviderWhisper,
		},
		{
			name:         "アムハラ語はWhisper推奨",
			languageCode: "am",
			expected:     ProviderWhisper,
		},
		{
			name:         "ペルシャ語はWhisper推奨",
			languageCode: "fa",
			expected:     ProviderWhisper,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			provider := GetRecommendedSTTProvider(tt.languageCode)
			assert.Equal(t, tt.expected, provider)
		})
	}
}

// TestGetRecommendedEvalMethod は推奨発音評価方法取得テスト
func TestGetRecommendedEvalMethod(t *testing.T) {
	tests := []struct {
		name         string
		languageCode string
		expected     PronunciationEvalMethod
	}{
		{
			name:         "英語はAzure Native",
			languageCode: "en",
			expected:     EvalMethodAzureNative,
		},
		{
			name:         "日本語はAzure Native",
			languageCode: "ja",
			expected:     EvalMethodAzureNative,
		},
		{
			name:         "クルド語はWhisper+LLM",
			languageCode: "ku",
			expected:     EvalMethodWhisperLLM,
		},
		{
			name:         "アムハラ語はWhisper+LLM",
			languageCode: "am",
			expected:     EvalMethodWhisperLLM,
		},
		{
			name:         "チベット語はWhisper+LLM",
			languageCode: "bo",
			expected:     EvalMethodWhisperLLM,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			method := GetRecommendedEvalMethod(tt.languageCode)
			assert.Equal(t, tt.expected, method)
		})
	}
}
