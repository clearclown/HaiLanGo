package llm

import (
	"context"
)

// LLMClient はLLM APIのインターフェース
// ドメイン特化やプロンプト生成はLLMの性能に任せる（過剰設計を避ける）
type LLMClient interface {
	// Generate はプロンプトに基づいてテキストを生成する
	Generate(ctx context.Context, prompt string, options *GenerateOptions) (*GenerateResult, error)

	// Chat はチャット形式のやり取りを行う
	Chat(ctx context.Context, messages []Message, options *GenerateOptions) (*GenerateResult, error)

	// GetName はプロバイダー名を返す
	GetName() string
}

// Message はチャットメッセージ
type Message struct {
	Role    string `json:"role"`    // system, user, assistant
	Content string `json:"content"` // メッセージ内容
}

// GenerateOptions は生成オプション
type GenerateOptions struct {
	MaxTokens    int      `json:"max_tokens"`    // 最大トークン数
	Temperature  float64  `json:"temperature"`   // 温度 (0.0-2.0)
	TopP         float64  `json:"top_p"`         // Top-P サンプリング
	StopWords    []string `json:"stop_words"`    // 停止ワード
	SystemPrompt string   `json:"system_prompt"` // システムプロンプト
}

// GenerateResult は生成結果
type GenerateResult struct {
	Content      string `json:"content"`       // 生成テキスト
	FinishReason string `json:"finish_reason"` // 終了理由
	TokensUsed   int    `json:"tokens_used"`   // 使用トークン数
	Provider     string `json:"provider"`      // 使用プロバイダー
	Model        string `json:"model"`         // 使用モデル
}

// LLMProvider はLLMプロバイダーの種類
type LLMProvider string

const (
	// ProviderClaude はAnthropic Claude（推奨: 多言語対応に優れる）
	ProviderClaude LLMProvider = "claude"

	// ProviderOpenAI はOpenAI GPT
	ProviderOpenAI LLMProvider = "openai"

	// ProviderGemini はGoogle Gemini
	ProviderGemini LLMProvider = "gemini"

	// ProviderLocal はローカルLLM（Llama, Mistral等）
	ProviderLocal LLMProvider = "local"
)

// DefaultGenerateOptions はデフォルトの生成オプション
func DefaultGenerateOptions() *GenerateOptions {
	return &GenerateOptions{
		MaxTokens:   1024,
		Temperature: 0.7,
		TopP:        0.9,
	}
}
