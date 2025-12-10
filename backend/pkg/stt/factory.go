package stt

import (
	"context"
	"fmt"
	"io"
	"os"
)

// NewSTTClient は環境変数に基づいて適切なSTTクライアントを返す
func NewSTTClient() (STTClient, error) {
	// モック使用の判定
	useMocks := os.Getenv("USE_MOCK_APIS") == "true" ||
		os.Getenv("TEST_USE_MOCKS") == "true"

	if useMocks {
		return NewMockSTTClientNew(), nil
	}

	// プロバイダーの選択
	provider := STTProvider(os.Getenv("STT_PROVIDER"))
	if provider == "" {
		provider = ProviderWhisper // デフォルト（99言語サポート、マイナー言語最強）
	}

	switch provider {
	case ProviderWhisper:
		apiKey := os.Getenv("OPENAI_API_KEY")
		if apiKey == "" {
			// APIキーがない場合は自動的にモックを使用
			return NewMockSTTClientNew(), nil
		}
		return NewWhisperSTTClient(apiKey), nil

	case ProviderAzureSpeech:
		subscriptionKey := os.Getenv("AZURE_SPEECH_KEY")
		region := os.Getenv("AZURE_SPEECH_REGION")
		if subscriptionKey == "" {
			// APIキーがない場合は自動的にモックを使用
			return NewMockSTTClientNew(), nil
		}
		if region == "" {
			region = "eastasia" // デフォルトリージョン
		}
		return NewAzureSTTClient(subscriptionKey, region), nil

	case ProviderGoogleSTT:
		apiKey := os.Getenv("GOOGLE_STT_KEY")
		if apiKey == "" {
			// APIキーがない場合は自動的にモックを使用
			return NewMockSTTClientNew(), nil
		}
		// TODO: GoogleSTTClient（新インターフェース対応）を実装
		return NewMockSTTClientNew(), nil

	case ProviderDeepgram:
		apiKey := os.Getenv("DEEPGRAM_API_KEY")
		if apiKey == "" {
			return NewMockSTTClientNew(), nil
		}
		// TODO: DeepgramClient を実装
		return NewMockSTTClientNew(), nil

	case ProviderWhisperLocal:
		// whisper.cppはAPIキー不要
		// TODO: WhisperLocalClient を実装
		return NewMockSTTClientNew(), nil

	default:
		return nil, fmt.Errorf("unsupported STT provider: %s", provider)
	}
}

// NewPronunciationEvaluator は環境変数に基づいて適切な発音評価クライアントを返す
func NewPronunciationEvaluator() (PronunciationEvaluator, error) {
	// モック使用の判定
	useMocks := os.Getenv("USE_MOCK_APIS") == "true" ||
		os.Getenv("TEST_USE_MOCKS") == "true"

	if useMocks {
		return NewMockPronunciationEvaluator(), nil
	}

	// 評価方法の選択
	method := PronunciationEvalMethod(os.Getenv("PRONUNCIATION_EVAL_METHOD"))
	if method == "" {
		method = EvalMethodWhisperLLM // デフォルト（マイナー言語対応）
	}

	switch method {
	case EvalMethodWhisperLLM:
		// Whisper + LLMによる評価
		openaiKey := os.Getenv("OPENAI_API_KEY")
		if openaiKey == "" {
			return NewMockPronunciationEvaluator(), nil
		}
		// TODO: WhisperLLMEvaluator を実装
		return NewMockPronunciationEvaluator(), nil

	case EvalMethodAzureNative:
		// Azure Speech発音評価
		subscriptionKey := os.Getenv("AZURE_SPEECH_KEY")
		if subscriptionKey == "" {
			return NewMockPronunciationEvaluator(), nil
		}
		// TODO: AzurePronunciationEvaluator を実装
		return NewMockPronunciationEvaluator(), nil

	default:
		return nil, fmt.Errorf("unsupported pronunciation evaluation method: %s", method)
	}
}

// WhisperSTTClient はOpenAI Whisper APIクライアント
type WhisperSTTClient struct {
	apiKey string
}

// NewWhisperSTTClient は新しいWhisper STTクライアントを作成する
func NewWhisperSTTClient(apiKey string) *WhisperSTTClient {
	return &WhisperSTTClient{
		apiKey: apiKey,
	}
}

// Transcribe は音声をテキストに変換する
// TODO: 実際のOpenAI Whisper API呼び出しを実装
func (c *WhisperSTTClient) Transcribe(ctx context.Context, audio io.Reader, language string) (*TranscriptionResult, error) {
	// 現在はモックを返す（実装予定）
	mock := NewMockSTTClientNew()
	return mock.Transcribe(ctx, audio, language)
}

// GetSupportedLanguages はサポートされている言語一覧を取得
func (c *WhisperSTTClient) GetSupportedLanguages(ctx context.Context) ([]string, error) {
	return WhisperSupportedLanguages, nil
}

// GetName はプロバイダー名を返す
func (c *WhisperSTTClient) GetName() string {
	return "whisper"
}

// Ensure WhisperSTTClient implements STTClient interface
var _ STTClient = (*WhisperSTTClient)(nil)

// AzureSTTClient はAzure Speech ServicesのSTTクライアント
type AzureSTTClient struct {
	subscriptionKey string
	region          string
}

// NewAzureSTTClient は新しいAzure STTクライアントを作成する
func NewAzureSTTClient(subscriptionKey, region string) *AzureSTTClient {
	return &AzureSTTClient{
		subscriptionKey: subscriptionKey,
		region:          region,
	}
}

// Transcribe は音声をテキストに変換する
// TODO: 実際のAzure Speech Services API呼び出しを実装
func (c *AzureSTTClient) Transcribe(ctx context.Context, audio io.Reader, language string) (*TranscriptionResult, error) {
	// 現在はモックを返す（実装予定）
	mock := NewMockSTTClientNew()
	return mock.Transcribe(ctx, audio, language)
}

// GetSupportedLanguages はサポートされている言語一覧を取得
// TODO: 実際のAzure Speech Services API呼び出しを実装
func (c *AzureSTTClient) GetSupportedLanguages(ctx context.Context) ([]string, error) {
	mock := NewMockSTTClientNew()
	return mock.GetSupportedLanguages(ctx)
}

// GetName はプロバイダー名を返す
func (c *AzureSTTClient) GetName() string {
	return "azure"
}

// Ensure AzureSTTClient implements STTClient interface
var _ STTClient = (*AzureSTTClient)(nil)
