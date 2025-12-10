package tts

import (
	"context"
	"fmt"
	"io"
	"os"
)

// NewTTSClient は環境変数に基づいて適切なTTSクライアントを返す
func NewTTSClient() (TTSClient, error) {
	// モック使用の判定
	useMocks := os.Getenv("USE_MOCK_APIS") == "true" ||
		os.Getenv("TEST_USE_MOCKS") == "true"

	if useMocks {
		return NewMockTTSClient(), nil
	}

	// プロバイダーの選択
	provider := TTSProvider(os.Getenv("TTS_PROVIDER"))
	if provider == "" {
		provider = ProviderAzureTTS // デフォルト（140言語サポート）
	}

	switch provider {
	case ProviderAzureTTS:
		subscriptionKey := os.Getenv("AZURE_SPEECH_KEY")
		region := os.Getenv("AZURE_SPEECH_REGION")
		if subscriptionKey == "" {
			// APIキーがない場合は自動的にモックを使用
			return NewMockTTSClient(), nil
		}
		if region == "" {
			region = "eastasia" // デフォルトリージョン
		}
		return NewAzureTTSClient(subscriptionKey, region), nil

	case ProviderGoogleTTS:
		apiKey := os.Getenv("GOOGLE_TTS_KEY")
		if apiKey == "" {
			// APIキーがない場合は自動的にモックを使用
			return NewMockTTSClient(), nil
		}
		// TODO: GoogleTTSClient（新インターフェース対応）を実装
		return NewMockTTSClient(), nil

	case ProviderElevenLabs:
		apiKey := os.Getenv("ELEVENLABS_API_KEY")
		if apiKey == "" {
			return NewMockTTSClient(), nil
		}
		// TODO: ElevenLabsClient を実装
		return NewMockTTSClient(), nil

	case ProviderEdgeTTS:
		// edge-ttsは無料、APIキー不要
		// TODO: EdgeTTSClient を実装
		return NewMockTTSClient(), nil

	case ProviderCoquiTTS:
		// Coqui TTSはオープンソース
		// TODO: CoquiTTSClient を実装
		return NewMockTTSClient(), nil

	default:
		return nil, fmt.Errorf("unsupported TTS provider: %s", provider)
	}
}

// AzureTTSClient はAzure Speech ServicesのTTSクライアント
type AzureTTSClient struct {
	subscriptionKey string
	region          string
}

// NewAzureTTSClient は新しいAzure TTSクライアントを作成する
func NewAzureTTSClient(subscriptionKey, region string) *AzureTTSClient {
	return &AzureTTSClient{
		subscriptionKey: subscriptionKey,
		region:          region,
	}
}

// Synthesize はテキストを音声に変換する
// TODO: 実際のAzure Speech Services API呼び出しを実装
func (c *AzureTTSClient) Synthesize(ctx context.Context, text string, language string, voice string) (io.ReadCloser, error) {
	// 現在はモックを返す（実装予定）
	mock := NewMockTTSClient()
	return mock.Synthesize(ctx, text, language, voice)
}

// GetVoices は指定言語で利用可能な音声一覧を取得
// TODO: 実際のAzure Speech Services API呼び出しを実装
func (c *AzureTTSClient) GetVoices(ctx context.Context, language string) ([]Voice, error) {
	mock := NewMockTTSClient()
	return mock.GetVoices(ctx, language)
}

// GetSupportedLanguages はサポートされている言語一覧を取得
// TODO: 実際のAzure Speech Services API呼び出しを実装
func (c *AzureTTSClient) GetSupportedLanguages(ctx context.Context) ([]string, error) {
	mock := NewMockTTSClient()
	return mock.GetSupportedLanguages(ctx)
}

// GetName はプロバイダー名を返す
func (c *AzureTTSClient) GetName() string {
	return "azure"
}

// Ensure AzureTTSClient implements TTSClient interface
var _ TTSClient = (*AzureTTSClient)(nil)
