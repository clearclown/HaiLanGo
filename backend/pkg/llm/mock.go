package llm

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"
)

// MockLLMClient はモックLLMクライアント
type MockLLMClient struct {
	// カスタムレスポンスを設定可能（テスト用）
	CustomResponse *GenerateResult
	ShouldFail     bool
	FailError      error
}

// NewMockLLMClient は新しいモックLLMクライアントを作成する
func NewMockLLMClient() *MockLLMClient {
	return &MockLLMClient{}
}

// Generate はプロンプトに基づいてテキストを生成する（モック）
func (m *MockLLMClient) Generate(ctx context.Context, prompt string, options *GenerateOptions) (*GenerateResult, error) {
	// エラーシミュレーション
	if m.ShouldFail {
		if m.FailError != nil {
			return nil, m.FailError
		}
		return nil, errors.New("mock error")
	}

	// 空のプロンプトはエラー
	if prompt == "" {
		return nil, errors.New("prompt cannot be empty")
	}

	// カスタムレスポンスがあればそれを返す
	if m.CustomResponse != nil {
		return m.CustomResponse, nil
	}

	// オプションがnilの場合はデフォルトを使用
	if options == nil {
		options = DefaultGenerateOptions()
	}

	// モックレスポンスを生成
	response := m.generateMockResponse(prompt, options)

	return &GenerateResult{
		Content:      response,
		FinishReason: "stop",
		TokensUsed:   len(strings.Fields(response)) * 2, // 簡易的なトークン数計算
		Provider:     "mock",
		Model:        "mock-model-v1",
	}, nil
}

// Chat はチャット形式のやり取りを行う（モック）
func (m *MockLLMClient) Chat(ctx context.Context, messages []Message, options *GenerateOptions) (*GenerateResult, error) {
	// エラーシミュレーション
	if m.ShouldFail {
		if m.FailError != nil {
			return nil, m.FailError
		}
		return nil, errors.New("mock error")
	}

	// 空のメッセージはエラー
	if len(messages) == 0 {
		return nil, errors.New("messages cannot be empty")
	}

	// オプションがnilの場合はデフォルトを使用
	if options == nil {
		options = DefaultGenerateOptions()
	}

	// 最後のユーザーメッセージを取得
	var lastUserMessage string
	for i := len(messages) - 1; i >= 0; i-- {
		if messages[i].Role == "user" {
			lastUserMessage = messages[i].Content
			break
		}
	}

	// モックレスポンスを生成
	response := m.generateChatResponse(messages, lastUserMessage, options)

	return &GenerateResult{
		Content:      response,
		FinishReason: "stop",
		TokensUsed:   len(strings.Fields(response)) * 2,
		Provider:     "mock",
		Model:        "mock-chat-v1",
	}, nil
}

// GetName はプロバイダー名を返す
func (m *MockLLMClient) GetName() string {
	return "mock"
}

// generateMockResponse はプロンプトに基づいたモックレスポンスを生成
func (m *MockLLMClient) generateMockResponse(prompt string, options *GenerateOptions) string {
	// 特定のキーワードに対するレスポンス
	promptLower := strings.ToLower(prompt)

	switch {
	case strings.Contains(promptLower, "hello") || strings.Contains(promptLower, "hi"):
		return "Hello! How can I help you today?"
	case strings.Contains(promptLower, "translate"):
		return "Here is the translation of your text."
	case strings.Contains(promptLower, "explain"):
		return "Let me explain this concept in detail."
	case strings.Contains(promptLower, "write"):
		return "Here is the content you requested."
	case strings.Contains(promptLower, "pronunciation") || strings.Contains(promptLower, "発音"):
		return m.generatePronunciationFeedback()
	default:
		return fmt.Sprintf("This is a mock response to: %s", truncateString(prompt, 50))
	}
}

// generateChatResponse はチャット履歴に基づいたモックレスポンスを生成
func (m *MockLLMClient) generateChatResponse(messages []Message, lastUserMessage string, options *GenerateOptions) string {
	// システムプロンプトを確認
	var systemPrompt string
	for _, msg := range messages {
		if msg.Role == "system" {
			systemPrompt = msg.Content
			break
		}
	}

	// 発音評価のシステムプロンプトかどうか
	if strings.Contains(systemPrompt, "pronunciation") || strings.Contains(systemPrompt, "発音") {
		return m.generatePronunciationFeedback()
	}

	// 通常のチャットレスポンス
	return m.generateMockResponse(lastUserMessage, options)
}

// generatePronunciationFeedback は発音評価のモックフィードバックを生成
func (m *MockLLMClient) generatePronunciationFeedback() string {
	return `{
  "score": 85,
  "overall_feedback": "Good pronunciation! Your intonation is natural.",
  "details": {
    "accuracy": 88,
    "fluency": 82,
    "prosody": 85
  },
  "suggestions": [
    "Try to pronounce the 'r' sound more clearly",
    "Pay attention to word stress"
  ]
}`
}

// truncateString は文字列を指定した長さに切り詰める
func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}

// SetCustomResponse はカスタムレスポンスを設定する（テスト用）
func (m *MockLLMClient) SetCustomResponse(response *GenerateResult) {
	m.CustomResponse = response
}

// SetShouldFail はエラーを発生させるかどうかを設定する（テスト用）
func (m *MockLLMClient) SetShouldFail(fail bool, err error) {
	m.ShouldFail = fail
	m.FailError = err
}

// SimulateLatency はレイテンシをシミュレートする（テスト用）
func (m *MockLLMClient) SimulateLatency(duration time.Duration) {
	time.Sleep(duration)
}

// Ensure MockLLMClient implements LLMClient interface
var _ LLMClient = (*MockLLMClient)(nil)
