package llm

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestNewLLMClient_WithMocks はモック環境でのLLMクライアント生成テスト
func TestNewLLMClient_WithMocks(t *testing.T) {
	os.Setenv("USE_MOCK_APIS", "true")
	defer os.Unsetenv("USE_MOCK_APIS")

	client, err := NewLLMClient()

	require.NoError(t, err)
	assert.NotNil(t, client)
	assert.Equal(t, "mock", client.GetName())
}

// TestNewLLMClient_WithoutAPIKey はAPIキーなしでモックが使用されることをテスト
func TestNewLLMClient_WithoutAPIKey(t *testing.T) {
	os.Unsetenv("USE_MOCK_APIS")
	os.Unsetenv("TEST_USE_MOCKS")
	os.Unsetenv("ANTHROPIC_API_KEY")
	os.Unsetenv("OPENAI_API_KEY")
	os.Unsetenv("LLM_PROVIDER")

	client, err := NewLLMClient()

	require.NoError(t, err)
	assert.NotNil(t, client)
	assert.Equal(t, "mock", client.GetName())
}

// TestNewLLMClient_WithProvider はプロバイダー指定のテスト
func TestNewLLMClient_WithProvider(t *testing.T) {
	tests := []struct {
		name       string
		provider   string
		expectName string
	}{
		{
			name:       "Claudeプロバイダー（デフォルト）",
			provider:   "claude",
			expectName: "mock", // APIキーなしでモック
		},
		{
			name:       "OpenAIプロバイダー",
			provider:   "openai",
			expectName: "mock", // APIキーなしでモック
		},
		{
			name:       "Geminiプロバイダー",
			provider:   "gemini",
			expectName: "mock", // APIキーなしでモック
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Unsetenv("USE_MOCK_APIS")
			os.Unsetenv("ANTHROPIC_API_KEY")
			os.Unsetenv("OPENAI_API_KEY")
			os.Unsetenv("GOOGLE_AI_KEY")
			os.Setenv("LLM_PROVIDER", tt.provider)
			defer os.Unsetenv("LLM_PROVIDER")

			client, err := NewLLMClient()

			require.NoError(t, err)
			assert.NotNil(t, client)
			assert.Equal(t, tt.expectName, client.GetName())
		})
	}
}

// TestMockLLMClient_Generate はモックLLMクライアントのGenerateテスト
func TestMockLLMClient_Generate(t *testing.T) {
	ctx := context.Background()
	client := NewMockLLMClient()

	tests := []struct {
		name        string
		prompt      string
		options     *GenerateOptions
		expectError bool
	}{
		{
			name:        "正常な生成",
			prompt:      "Hello, how are you?",
			options:     DefaultGenerateOptions(),
			expectError: false,
		},
		{
			name:        "空のプロンプト",
			prompt:      "",
			options:     DefaultGenerateOptions(),
			expectError: true,
		},
		{
			name:        "nilオプション（デフォルト使用）",
			prompt:      "Test prompt",
			options:     nil,
			expectError: false,
		},
		{
			name:   "カスタムオプション",
			prompt: "Write a story",
			options: &GenerateOptions{
				MaxTokens:   2048,
				Temperature: 0.9,
				TopP:        0.95,
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := client.Generate(ctx, tt.prompt, tt.options)

			if tt.expectError {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.NotNil(t, result)
			assert.NotEmpty(t, result.Content)
			assert.Equal(t, "mock", result.Provider)
			assert.NotEmpty(t, result.Model)
			assert.Greater(t, result.TokensUsed, 0)
		})
	}
}

// TestMockLLMClient_Chat はモックLLMクライアントのChatテスト
func TestMockLLMClient_Chat(t *testing.T) {
	ctx := context.Background()
	client := NewMockLLMClient()

	tests := []struct {
		name        string
		messages    []Message
		options     *GenerateOptions
		expectError bool
	}{
		{
			name: "正常なチャット",
			messages: []Message{
				{Role: "user", Content: "Hello!"},
			},
			options:     DefaultGenerateOptions(),
			expectError: false,
		},
		{
			name: "システムプロンプト付きチャット",
			messages: []Message{
				{Role: "system", Content: "You are a helpful assistant."},
				{Role: "user", Content: "What is the weather?"},
			},
			options:     DefaultGenerateOptions(),
			expectError: false,
		},
		{
			name: "複数ターンのチャット",
			messages: []Message{
				{Role: "user", Content: "Hi"},
				{Role: "assistant", Content: "Hello! How can I help you?"},
				{Role: "user", Content: "Tell me a joke"},
			},
			options:     DefaultGenerateOptions(),
			expectError: false,
		},
		{
			name:        "空のメッセージ",
			messages:    []Message{},
			options:     DefaultGenerateOptions(),
			expectError: true,
		},
		{
			name:        "nilメッセージ",
			messages:    nil,
			options:     DefaultGenerateOptions(),
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := client.Chat(ctx, tt.messages, tt.options)

			if tt.expectError {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.NotNil(t, result)
			assert.NotEmpty(t, result.Content)
			assert.Equal(t, "mock", result.Provider)
		})
	}
}

// TestDefaultGenerateOptions はデフォルトオプションのテスト
func TestDefaultGenerateOptions(t *testing.T) {
	options := DefaultGenerateOptions()

	assert.NotNil(t, options)
	assert.Equal(t, 1024, options.MaxTokens)
	assert.Equal(t, 0.7, options.Temperature)
	assert.Equal(t, 0.9, options.TopP)
}

// TestMockLLMClient_MultiLanguageGeneration は多言語生成テスト
func TestMockLLMClient_MultiLanguageGeneration(t *testing.T) {
	ctx := context.Background()
	client := NewMockLLMClient()

	tests := []struct {
		name   string
		prompt string
	}{
		{
			name:   "日本語のプロンプト",
			prompt: "こんにちは、元気ですか？",
		},
		{
			name:   "中国語のプロンプト",
			prompt: "你好，今天天气怎么样？",
		},
		{
			name:   "ロシア語のプロンプト",
			prompt: "Здравствуйте, как дела?",
		},
		{
			name:   "アラビア語のプロンプト",
			prompt: "مرحبا، كيف حالك؟",
		},
		{
			name:   "ペルシャ語のプロンプト",
			prompt: "سلام، حال شما چطور است؟",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := client.Generate(ctx, tt.prompt, DefaultGenerateOptions())

			require.NoError(t, err)
			assert.NotNil(t, result)
			assert.NotEmpty(t, result.Content)
		})
	}
}

// TestMockLLMClient_PronunciationEvaluation は発音評価用のLLM使用テスト
func TestMockLLMClient_PronunciationEvaluation(t *testing.T) {
	ctx := context.Background()
	client := NewMockLLMClient()

	// 発音評価のためのプロンプト例
	messages := []Message{
		{
			Role: "system",
			Content: `You are a language pronunciation evaluator.
Analyze the following speech-to-text result and provide a pronunciation score (0-100).
Focus on:
1. Accuracy: Did the user say the expected words correctly?
2. Missing words: Were any words omitted?
3. Extra words: Were any extra words added?
Provide feedback in the user's native language.`,
		},
		{
			Role:    "user",
			Content: "Expected: 'Hello, how are you?' Recognized: 'Hello, how you?'",
		},
	}

	result, err := client.Chat(ctx, messages, DefaultGenerateOptions())

	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.NotEmpty(t, result.Content)
}
