package stt

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"time"

	"github.com/clearclown/HaiLanGo/backend/pkg/llm"
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
		return NewGoogleSTTClient(apiKey), nil

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
		// STTクライアントを取得
		sttClient, err := NewSTTClient()
		if err != nil {
			return NewMockPronunciationEvaluator(), nil
		}

		// LLMクライアントを取得
		llmClient, err := llm.NewLLMClient()
		if err != nil {
			return NewMockPronunciationEvaluator(), nil
		}

		return NewWhisperLLMEvaluator(sttClient, llmClient), nil

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

// whisperResponse はOpenAI Whisper APIのレスポンス構造
type whisperResponse struct {
	Text     string          `json:"text"`
	Task     string          `json:"task,omitempty"`
	Language string          `json:"language,omitempty"`
	Duration float64         `json:"duration,omitempty"`
	Segments []whisperSegment `json:"segments,omitempty"`
}

type whisperSegment struct {
	ID               int     `json:"id"`
	Seek             int     `json:"seek"`
	Start            float64 `json:"start"`
	End              float64 `json:"end"`
	Text             string  `json:"text"`
	AvgLogprob       float64 `json:"avg_logprob"`
	CompressionRatio float64 `json:"compression_ratio"`
	NoSpeechProb     float64 `json:"no_speech_prob"`
}

// Transcribe は音声をテキストに変換する
func (c *WhisperSTTClient) Transcribe(ctx context.Context, audio io.Reader, language string) (*TranscriptionResult, error) {
	startTime := time.Now()

	// 音声データを読み込み
	audioData, err := io.ReadAll(audio)
	if err != nil {
		return nil, fmt.Errorf("failed to read audio data: %w", err)
	}

	// マルチパートフォームを作成
	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)

	// ファイルパートを追加
	part, err := writer.CreateFormFile("file", "audio.wav")
	if err != nil {
		return nil, fmt.Errorf("failed to create form file: %w", err)
	}
	if _, err := part.Write(audioData); err != nil {
		return nil, fmt.Errorf("failed to write audio data: %w", err)
	}

	// モデルパートを追加
	if err := writer.WriteField("model", "whisper-1"); err != nil {
		return nil, fmt.Errorf("failed to write model field: %w", err)
	}

	// 言語パートを追加（指定されている場合）
	if language != "" {
		if err := writer.WriteField("language", language); err != nil {
			return nil, fmt.Errorf("failed to write language field: %w", err)
		}
	}

	// 詳細なレスポンスを要求
	if err := writer.WriteField("response_format", "verbose_json"); err != nil {
		return nil, fmt.Errorf("failed to write response_format field: %w", err)
	}

	// タイムスタンプの粒度を設定
	if err := writer.WriteField("timestamp_granularities[]", "segment"); err != nil {
		return nil, fmt.Errorf("failed to write timestamp_granularities field: %w", err)
	}

	if err := writer.Close(); err != nil {
		return nil, fmt.Errorf("failed to close multipart writer: %w", err)
	}

	// HTTPリクエストを作成
	req, err := http.NewRequestWithContext(ctx, "POST", "https://api.openai.com/v1/audio/transcriptions", &buf)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	// リクエストを送信
	client := &http.Client{Timeout: 120 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Whisper API error (status %d): %s", resp.StatusCode, string(respBody))
	}

	// レスポンスをパース
	var whisperResp whisperResponse
	if err := json.Unmarshal(respBody, &whisperResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	// セグメントを変換
	segments := make([]TranscriptSegment, len(whisperResp.Segments))
	for i, seg := range whisperResp.Segments {
		segments[i] = TranscriptSegment{
			ID:         seg.ID,
			Start:      seg.Start,
			End:        seg.End,
			Text:       seg.Text,
			Confidence: 1.0 - seg.NoSpeechProb, // no_speech_probを信頼度に変換
		}
	}

	processingTime := time.Since(startTime).Milliseconds()

	return &TranscriptionResult{
		Text:             whisperResp.Text,
		Language:         whisperResp.Language,
		Confidence:       0.95, // Whisperはデフォルトで高精度
		Duration:         int64(whisperResp.Duration * 1000),
		Segments:         segments,
		Provider:         "whisper",
		ProcessingTimeMs: processingTime,
	}, nil
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

// azureSTTResponse はAzure Speech STT APIのレスポンス構造
type azureSTTResponse struct {
	RecognitionStatus string `json:"RecognitionStatus"`
	DisplayText       string `json:"DisplayText"`
	Offset            int64  `json:"Offset"`
	Duration          int64  `json:"Duration"`
}

// Transcribe は音声をテキストに変換する
func (c *AzureSTTClient) Transcribe(ctx context.Context, audio io.Reader, language string) (*TranscriptionResult, error) {
	startTime := time.Now()

	// デフォルト言語を設定
	if language == "" {
		language = "en-US"
	}

	// 音声データを読み込み
	audioData, err := io.ReadAll(audio)
	if err != nil {
		return nil, fmt.Errorf("failed to read audio data: %w", err)
	}

	// Azure Speech STT APIエンドポイント
	url := fmt.Sprintf("https://%s.stt.speech.microsoft.com/speech/recognition/conversation/cognitiveservices/v1?language=%s",
		c.region, language)

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(audioData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Ocp-Apim-Subscription-Key", c.subscriptionKey)
	req.Header.Set("Content-Type", "audio/wav; codecs=audio/pcm; samplerate=16000")
	req.Header.Set("Accept", "application/json")

	client := &http.Client{Timeout: 60 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Azure STT API error (status %d): %s", resp.StatusCode, string(respBody))
	}

	// レスポンスをパース
	var azureResp azureSTTResponse
	if err := json.Unmarshal(respBody, &azureResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	// 認識ステータスをチェック
	if azureResp.RecognitionStatus != "Success" {
		return nil, fmt.Errorf("recognition failed: %s", azureResp.RecognitionStatus)
	}

	processingTime := time.Since(startTime).Milliseconds()

	return &TranscriptionResult{
		Text:             azureResp.DisplayText,
		Language:         language,
		Confidence:       0.9, // Azureはデフォルトで高精度
		Duration:         azureResp.Duration / 10000, // 100ナノ秒単位からミリ秒に変換
		Provider:         "azure",
		ProcessingTimeMs: processingTime,
	}, nil
}

// AzureSTTSupportedLanguages はAzure Speech-to-Textがサポートする言語
var AzureSTTSupportedLanguages = []string{
	"ar-AE", "ar-BH", "ar-EG", "ar-IL", "ar-IQ", "ar-JO", "ar-KW", "ar-LB", "ar-LY", "ar-MA",
	"ar-OM", "ar-PS", "ar-QA", "ar-SA", "ar-SY", "ar-TN", "ar-YE",
	"bg-BG", "ca-ES", "cs-CZ", "cy-GB", "da-DK", "de-AT", "de-CH", "de-DE",
	"el-GR", "en-AU", "en-CA", "en-GB", "en-GH", "en-HK", "en-IE", "en-IN", "en-KE",
	"en-NG", "en-NZ", "en-PH", "en-SG", "en-TZ", "en-US", "en-ZA",
	"es-AR", "es-BO", "es-CL", "es-CO", "es-CR", "es-CU", "es-DO", "es-EC", "es-ES",
	"es-GQ", "es-GT", "es-HN", "es-MX", "es-NI", "es-PA", "es-PE", "es-PR", "es-PY",
	"es-SV", "es-US", "es-UY", "es-VE",
	"et-EE", "eu-ES", "fa-IR", "fi-FI", "fil-PH", "fr-BE", "fr-CA", "fr-CH", "fr-FR",
	"ga-IE", "gl-ES", "gu-IN", "he-IL", "hi-IN", "hr-HR", "hu-HU", "hy-AM", "id-ID",
	"is-IS", "it-CH", "it-IT", "ja-JP", "jv-ID", "ka-GE", "kk-KZ", "km-KH", "kn-IN",
	"ko-KR", "lo-LA", "lt-LT", "lv-LV", "mk-MK", "ml-IN", "mn-MN", "mr-IN", "ms-MY",
	"mt-MT", "my-MM", "nb-NO", "ne-NP", "nl-BE", "nl-NL", "pl-PL", "ps-AF", "pt-BR",
	"pt-PT", "ro-RO", "ru-RU", "si-LK", "sk-SK", "sl-SI", "so-SO", "sq-AL", "sr-RS",
	"su-ID", "sv-SE", "sw-KE", "sw-TZ", "ta-IN", "ta-LK", "ta-MY", "ta-SG", "te-IN",
	"th-TH", "tr-TR", "uk-UA", "ur-IN", "ur-PK", "uz-UZ", "vi-VN", "wuu-CN",
	"yue-CN", "zh-CN", "zh-CN-sichuan", "zh-HK", "zh-TW", "zu-ZA",
}

// GetSupportedLanguages はサポートされている言語一覧を取得
func (c *AzureSTTClient) GetSupportedLanguages(ctx context.Context) ([]string, error) {
	return AzureSTTSupportedLanguages, nil
}

// GetName はプロバイダー名を返す
func (c *AzureSTTClient) GetName() string {
	return "azure"
}

// Ensure AzureSTTClient implements STTClient interface
var _ STTClient = (*AzureSTTClient)(nil)
