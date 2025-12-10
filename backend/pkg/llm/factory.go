package llm

import (
	"context"
	"fmt"
	"os"
)

// NewLLMClient は環境変数に基づいて適切なLLMクライアントを返す
func NewLLMClient() (LLMClient, error) {
	// モック使用の判定
	useMocks := os.Getenv("USE_MOCK_APIS") == "true" ||
		os.Getenv("TEST_USE_MOCKS") == "true"

	if useMocks {
		return NewMockLLMClient(), nil
	}

	// プロバイダーの選択
	provider := LLMProvider(os.Getenv("LLM_PROVIDER"))
	if provider == "" {
		provider = ProviderClaude // デフォルト（多言語対応、長文理解に優れる）
	}

	switch provider {
	case ProviderClaude:
		apiKey := os.Getenv("ANTHROPIC_API_KEY")
		if apiKey == "" {
			// APIキーがない場合は自動的にモックを使用
			return NewMockLLMClient(), nil
		}
		return NewClaudeClient(apiKey), nil

	case ProviderOpenAI:
		apiKey := os.Getenv("OPENAI_API_KEY")
		if apiKey == "" {
			// APIキーがない場合は自動的にモックを使用
			return NewMockLLMClient(), nil
		}
		return NewOpenAIClient(apiKey), nil

	case ProviderGemini:
		apiKey := os.Getenv("GOOGLE_AI_KEY")
		if apiKey == "" {
			// APIキーがない場合は自動的にモックを使用
			return NewMockLLMClient(), nil
		}
		return NewGeminiClient(apiKey), nil

	default:
		return nil, fmt.Errorf("unsupported LLM provider: %s", provider)
	}
}

// ClaudeClient はAnthropic Claude APIクライアント
type ClaudeClient struct {
	apiKey string
	model  string
}

// NewClaudeClient は新しいClaude LLMクライアントを作成する
func NewClaudeClient(apiKey string) *ClaudeClient {
	return &ClaudeClient{
		apiKey: apiKey,
		model:  "claude-3-5-sonnet-20241022",
	}
}

// Generate はプロンプトに基づいてテキストを生成する
// TODO: 実際のAnthropic API呼び出しを実装
func (c *ClaudeClient) Generate(ctx context.Context, prompt string, options *GenerateOptions) (*GenerateResult, error) {
	// 現在はモックを返す（実装予定）
	mock := NewMockLLMClient()
	return mock.Generate(ctx, prompt, options)
}

// Chat はチャット形式のやり取りを行う
// TODO: 実際のAnthropic API呼び出しを実装
func (c *ClaudeClient) Chat(ctx context.Context, messages []Message, options *GenerateOptions) (*GenerateResult, error) {
	// 現在はモックを返す（実装予定）
	mock := NewMockLLMClient()
	return mock.Chat(ctx, messages, options)
}

// GetName はプロバイダー名を返す
func (c *ClaudeClient) GetName() string {
	return "claude"
}

// Ensure ClaudeClient implements LLMClient interface
var _ LLMClient = (*ClaudeClient)(nil)

// OpenAIClient はOpenAI APIクライアント
type OpenAIClient struct {
	apiKey string
	model  string
}

// NewOpenAIClient は新しいOpenAI LLMクライアントを作成する
func NewOpenAIClient(apiKey string) *OpenAIClient {
	return &OpenAIClient{
		apiKey: apiKey,
		model:  "gpt-4o",
	}
}

// Generate はプロンプトに基づいてテキストを生成する
// TODO: 実際のOpenAI API呼び出しを実装
func (c *OpenAIClient) Generate(ctx context.Context, prompt string, options *GenerateOptions) (*GenerateResult, error) {
	// 現在はモックを返す（実装予定）
	mock := NewMockLLMClient()
	return mock.Generate(ctx, prompt, options)
}

// Chat はチャット形式のやり取りを行う
// TODO: 実際のOpenAI API呼び出しを実装
func (c *OpenAIClient) Chat(ctx context.Context, messages []Message, options *GenerateOptions) (*GenerateResult, error) {
	// 現在はモックを返す（実装予定）
	mock := NewMockLLMClient()
	return mock.Chat(ctx, messages, options)
}

// GetName はプロバイダー名を返す
func (c *OpenAIClient) GetName() string {
	return "openai"
}

// Ensure OpenAIClient implements LLMClient interface
var _ LLMClient = (*OpenAIClient)(nil)

// GeminiClient はGoogle Gemini APIクライアント
type GeminiClient struct {
	apiKey string
	model  string
}

// NewGeminiClient は新しいGemini LLMクライアントを作成する
func NewGeminiClient(apiKey string) *GeminiClient {
	return &GeminiClient{
		apiKey: apiKey,
		model:  "gemini-1.5-pro",
	}
}

// Generate はプロンプトに基づいてテキストを生成する
// TODO: 実際のGoogle AI API呼び出しを実装
func (c *GeminiClient) Generate(ctx context.Context, prompt string, options *GenerateOptions) (*GenerateResult, error) {
	// 現在はモックを返す（実装予定）
	mock := NewMockLLMClient()
	return mock.Generate(ctx, prompt, options)
}

// Chat はチャット形式のやり取りを行う
// TODO: 実際のGoogle AI API呼び出しを実装
func (c *GeminiClient) Chat(ctx context.Context, messages []Message, options *GenerateOptions) (*GenerateResult, error) {
	// 現在はモックを返す（実装予定）
	mock := NewMockLLMClient()
	return mock.Chat(ctx, messages, options)
}

// GetName はプロバイダー名を返す
func (c *GeminiClient) GetName() string {
	return "gemini"
}

// Ensure GeminiClient implements LLMClient interface
var _ LLMClient = (*GeminiClient)(nil)
