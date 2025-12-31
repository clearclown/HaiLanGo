package llm

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
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
	apiKey     string
	model      string
	httpClient *http.Client
}

// NewClaudeClient は新しいClaude LLMクライアントを作成する
func NewClaudeClient(apiKey string) *ClaudeClient {
	return &ClaudeClient{
		apiKey: apiKey,
		model:  "claude-sonnet-4-20250514",
		httpClient: &http.Client{
			Timeout: 120 * time.Second,
		},
	}
}

// claudeRequest はAnthropic APIリクエスト
type claudeRequest struct {
	Model     string          `json:"model"`
	MaxTokens int             `json:"max_tokens"`
	Messages  []claudeMessage `json:"messages"`
	System    string          `json:"system,omitempty"`
}

type claudeMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type claudeResponse struct {
	ID      string `json:"id"`
	Type    string `json:"type"`
	Role    string `json:"role"`
	Content []struct {
		Type string `json:"type"`
		Text string `json:"text"`
	} `json:"content"`
	StopReason string `json:"stop_reason"`
	Usage      struct {
		InputTokens  int `json:"input_tokens"`
		OutputTokens int `json:"output_tokens"`
	} `json:"usage"`
}

// Generate はプロンプトに基づいてテキストを生成する
func (c *ClaudeClient) Generate(ctx context.Context, prompt string, options *GenerateOptions) (*GenerateResult, error) {
	messages := []Message{{Role: "user", Content: prompt}}
	return c.Chat(ctx, messages, options)
}

// Chat はチャット形式のやり取りを行う
func (c *ClaudeClient) Chat(ctx context.Context, messages []Message, options *GenerateOptions) (*GenerateResult, error) {
	if options == nil {
		options = DefaultGenerateOptions()
	}

	// メッセージ変換
	claudeMessages := make([]claudeMessage, 0, len(messages))
	for _, m := range messages {
		if m.Role != "system" {
			claudeMessages = append(claudeMessages, claudeMessage{
				Role:    m.Role,
				Content: m.Content,
			})
		}
	}

	req := claudeRequest{
		Model:     c.model,
		MaxTokens: options.MaxTokens,
		Messages:  claudeMessages,
		System:    options.SystemPrompt,
	}

	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", "https://api.anthropic.com/v1/messages", bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("x-api-key", c.apiKey)
	httpReq.Header.Set("anthropic-version", "2023-06-01")

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(respBody))
	}

	var claudeResp claudeResponse
	if err := json.Unmarshal(respBody, &claudeResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	content := ""
	if len(claudeResp.Content) > 0 {
		content = claudeResp.Content[0].Text
	}

	return &GenerateResult{
		Content:      content,
		FinishReason: claudeResp.StopReason,
		TokensUsed:   claudeResp.Usage.InputTokens + claudeResp.Usage.OutputTokens,
		Provider:     "claude",
		Model:        c.model,
	}, nil
}

// GetName はプロバイダー名を返す
func (c *ClaudeClient) GetName() string {
	return "claude"
}

// Ensure ClaudeClient implements LLMClient interface
var _ LLMClient = (*ClaudeClient)(nil)

// OpenAIClient はOpenAI APIクライアント
type OpenAIClient struct {
	apiKey     string
	model      string
	httpClient *http.Client
}

// NewOpenAIClient は新しいOpenAI LLMクライアントを作成する
func NewOpenAIClient(apiKey string) *OpenAIClient {
	return &OpenAIClient{
		apiKey: apiKey,
		model:  "gpt-4o",
		httpClient: &http.Client{
			Timeout: 120 * time.Second,
		},
	}
}

// openaiRequest はOpenAI APIリクエスト
type openaiRequest struct {
	Model       string          `json:"model"`
	Messages    []openaiMessage `json:"messages"`
	MaxTokens   int             `json:"max_tokens,omitempty"`
	Temperature float64         `json:"temperature,omitempty"`
	TopP        float64         `json:"top_p,omitempty"`
}

type openaiMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type openaiResponse struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Created int64  `json:"created"`
	Model   string `json:"model"`
	Choices []struct {
		Index   int `json:"index"`
		Message struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		} `json:"message"`
		FinishReason string `json:"finish_reason"`
	} `json:"choices"`
	Usage struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage"`
}

// Generate はプロンプトに基づいてテキストを生成する
func (c *OpenAIClient) Generate(ctx context.Context, prompt string, options *GenerateOptions) (*GenerateResult, error) {
	messages := []Message{{Role: "user", Content: prompt}}
	return c.Chat(ctx, messages, options)
}

// Chat はチャット形式のやり取りを行う
func (c *OpenAIClient) Chat(ctx context.Context, messages []Message, options *GenerateOptions) (*GenerateResult, error) {
	if options == nil {
		options = DefaultGenerateOptions()
	}

	// メッセージ変換
	openaiMessages := make([]openaiMessage, 0, len(messages))
	if options.SystemPrompt != "" {
		openaiMessages = append(openaiMessages, openaiMessage{
			Role:    "system",
			Content: options.SystemPrompt,
		})
	}
	for _, m := range messages {
		openaiMessages = append(openaiMessages, openaiMessage{
			Role:    m.Role,
			Content: m.Content,
		})
	}

	req := openaiRequest{
		Model:       c.model,
		Messages:    openaiMessages,
		MaxTokens:   options.MaxTokens,
		Temperature: options.Temperature,
		TopP:        options.TopP,
	}

	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", "https://api.openai.com/v1/chat/completions", bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+c.apiKey)

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(respBody))
	}

	var openaiResp openaiResponse
	if err := json.Unmarshal(respBody, &openaiResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	content := ""
	finishReason := ""
	if len(openaiResp.Choices) > 0 {
		content = openaiResp.Choices[0].Message.Content
		finishReason = openaiResp.Choices[0].FinishReason
	}

	return &GenerateResult{
		Content:      content,
		FinishReason: finishReason,
		TokensUsed:   openaiResp.Usage.TotalTokens,
		Provider:     "openai",
		Model:        c.model,
	}, nil
}

// GetName はプロバイダー名を返す
func (c *OpenAIClient) GetName() string {
	return "openai"
}

// Ensure OpenAIClient implements LLMClient interface
var _ LLMClient = (*OpenAIClient)(nil)

// GeminiClient はGoogle Gemini APIクライアント
type GeminiClient struct {
	apiKey     string
	model      string
	httpClient *http.Client
}

// NewGeminiClient は新しいGemini LLMクライアントを作成する
func NewGeminiClient(apiKey string) *GeminiClient {
	return &GeminiClient{
		apiKey: apiKey,
		model:  "gemini-1.5-pro",
		httpClient: &http.Client{
			Timeout: 120 * time.Second,
		},
	}
}

// geminiRequest はGemini APIリクエスト
type geminiRequest struct {
	Contents         []geminiContent `json:"contents"`
	GenerationConfig *geminiConfig   `json:"generationConfig,omitempty"`
	SystemInstruction *geminiContent `json:"systemInstruction,omitempty"`
}

type geminiContent struct {
	Role  string       `json:"role,omitempty"`
	Parts []geminiPart `json:"parts"`
}

type geminiPart struct {
	Text string `json:"text"`
}

type geminiConfig struct {
	MaxOutputTokens int     `json:"maxOutputTokens,omitempty"`
	Temperature     float64 `json:"temperature,omitempty"`
	TopP            float64 `json:"topP,omitempty"`
}

type geminiResponse struct {
	Candidates []struct {
		Content struct {
			Parts []struct {
				Text string `json:"text"`
			} `json:"parts"`
			Role string `json:"role"`
		} `json:"content"`
		FinishReason string `json:"finishReason"`
	} `json:"candidates"`
	UsageMetadata struct {
		PromptTokenCount     int `json:"promptTokenCount"`
		CandidatesTokenCount int `json:"candidatesTokenCount"`
		TotalTokenCount      int `json:"totalTokenCount"`
	} `json:"usageMetadata"`
}

// Generate はプロンプトに基づいてテキストを生成する
func (c *GeminiClient) Generate(ctx context.Context, prompt string, options *GenerateOptions) (*GenerateResult, error) {
	messages := []Message{{Role: "user", Content: prompt}}
	return c.Chat(ctx, messages, options)
}

// Chat はチャット形式のやり取りを行う
func (c *GeminiClient) Chat(ctx context.Context, messages []Message, options *GenerateOptions) (*GenerateResult, error) {
	if options == nil {
		options = DefaultGenerateOptions()
	}

	// メッセージ変換
	geminiContents := make([]geminiContent, 0, len(messages))
	for _, m := range messages {
		role := m.Role
		if role == "assistant" {
			role = "model"
		}
		if role != "system" {
			geminiContents = append(geminiContents, geminiContent{
				Role:  role,
				Parts: []geminiPart{{Text: m.Content}},
			})
		}
	}

	req := geminiRequest{
		Contents: geminiContents,
		GenerationConfig: &geminiConfig{
			MaxOutputTokens: options.MaxTokens,
			Temperature:     options.Temperature,
			TopP:            options.TopP,
		},
	}

	// システムプロンプトがある場合
	if options.SystemPrompt != "" {
		req.SystemInstruction = &geminiContent{
			Parts: []geminiPart{{Text: options.SystemPrompt}},
		}
	}

	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	url := fmt.Sprintf("https://generativelanguage.googleapis.com/v1beta/models/%s:generateContent?key=%s", c.model, c.apiKey)
	httpReq, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(respBody))
	}

	var geminiResp geminiResponse
	if err := json.Unmarshal(respBody, &geminiResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	content := ""
	finishReason := ""
	if len(geminiResp.Candidates) > 0 && len(geminiResp.Candidates[0].Content.Parts) > 0 {
		content = geminiResp.Candidates[0].Content.Parts[0].Text
		finishReason = geminiResp.Candidates[0].FinishReason
	}

	return &GenerateResult{
		Content:      content,
		FinishReason: finishReason,
		TokensUsed:   geminiResp.UsageMetadata.TotalTokenCount,
		Provider:     "gemini",
		Model:        c.model,
	}, nil
}

// GetName はプロバイダー名を返す
func (c *GeminiClient) GetName() string {
	return "gemini"
}

// Ensure GeminiClient implements LLMClient interface
var _ LLMClient = (*GeminiClient)(nil)
