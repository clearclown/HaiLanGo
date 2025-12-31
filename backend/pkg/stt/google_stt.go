package stt

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/clearclown/HaiLanGo/backend/internal/models"
)

// GoogleSTTClient はGoogle Cloud Speech-to-Text APIクライアント
// LegacySTTClient と STTClient の両方を実装
type GoogleSTTClient struct {
	apiKey string
}

// NewGoogleSTTClient は新しいGoogle STTクライアントを作成する
func NewGoogleSTTClient(apiKey string) *GoogleSTTClient {
	if apiKey == "" {
		apiKey = os.Getenv("GOOGLE_CLOUD_STT_API_KEY")
	}
	return &GoogleSTTClient{
		apiKey: apiKey,
	}
}

// googleSTTRequest はGoogle Cloud STT APIのリクエスト構造
type googleSTTRequest struct {
	Config googleSTTConfig `json:"config"`
	Audio  googleSTTAudio  `json:"audio"`
}

type googleSTTConfig struct {
	Encoding                   string   `json:"encoding"`
	SampleRateHertz            int      `json:"sampleRateHertz,omitempty"`
	LanguageCode               string   `json:"languageCode"`
	EnableAutomaticPunctuation bool     `json:"enableAutomaticPunctuation"`
	EnableWordTimeOffsets      bool     `json:"enableWordTimeOffsets"`
	Model                      string   `json:"model,omitempty"`
	AlternativeLanguageCodes   []string `json:"alternativeLanguageCodes,omitempty"`
}

type googleSTTAudio struct {
	Content string `json:"content"`
}

// googleSTTResponse はGoogle Cloud STT APIのレスポンス構造
type googleSTTResponse struct {
	Results []googleSTTResult `json:"results"`
}

type googleSTTResult struct {
	Alternatives  []googleSTTAlternative `json:"alternatives"`
	ResultEndTime string                 `json:"resultEndTime,omitempty"`
	LanguageCode  string                 `json:"languageCode,omitempty"`
}

type googleSTTAlternative struct {
	Transcript string              `json:"transcript"`
	Confidence float64             `json:"confidence"`
	Words      []googleSTTWordInfo `json:"words,omitempty"`
}

type googleSTTWordInfo struct {
	Word       string  `json:"word"`
	StartTime  string  `json:"startTime"`
	EndTime    string  `json:"endTime"`
	Confidence float64 `json:"confidence,omitempty"`
}

// ============================================================
// STTClient interface implementation (new interface)
// ============================================================

// Transcribe は音声をテキストに変換する (STTClient interface)
func (c *GoogleSTTClient) Transcribe(ctx context.Context, audio io.Reader, language string) (*TranscriptionResult, error) {
	startTime := time.Now()

	// 音声データを読み込み
	audioData, err := io.ReadAll(audio)
	if err != nil {
		return nil, fmt.Errorf("failed to read audio data: %w", err)
	}

	// デフォルト言語を設定
	if language == "" {
		language = "en-US"
	}

	// 言語コードを正規化
	language = c.normalizeLanguageCode(language)

	// リクエストボディを構築
	reqBody := googleSTTRequest{
		Config: googleSTTConfig{
			Encoding:                   "LINEAR16",
			LanguageCode:               language,
			EnableAutomaticPunctuation: true,
			EnableWordTimeOffsets:      true,
			Model:                      "default",
		},
		Audio: googleSTTAudio{
			Content: base64.StdEncoding.EncodeToString(audioData),
		},
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	// Google Cloud STT APIエンドポイント
	url := fmt.Sprintf("https://speech.googleapis.com/v1/speech:recognize?key=%s", c.apiKey)

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

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
		return nil, fmt.Errorf("Google STT API error (status %d): %s", resp.StatusCode, string(respBody))
	}

	// レスポンスをパース
	var sttResp googleSTTResponse
	if err := json.Unmarshal(respBody, &sttResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	// 結果を構築
	var fullText strings.Builder
	var words []WordTiming
	var confidence float64
	var wordCount int

	for _, result := range sttResp.Results {
		if len(result.Alternatives) == 0 {
			continue
		}

		alt := result.Alternatives[0]
		fullText.WriteString(alt.Transcript)
		confidence += alt.Confidence
		wordCount++

		for _, w := range alt.Words {
			words = append(words, WordTiming{
				Word:       w.Word,
				Start:      parseDuration(w.StartTime),
				End:        parseDuration(w.EndTime),
				Confidence: w.Confidence,
			})
		}
	}

	// 平均信頼度を計算
	if wordCount > 0 {
		confidence /= float64(wordCount)
	}

	processingTime := time.Since(startTime).Milliseconds()

	return &TranscriptionResult{
		Text:             fullText.String(),
		Language:         language,
		Confidence:       confidence,
		Words:            words,
		Provider:         "google",
		ProcessingTimeMs: processingTime,
	}, nil
}

// parseDuration はGoogle APIの時間形式（例: "1.500s"）をfloat64に変換する
func parseDuration(s string) float64 {
	s = strings.TrimSuffix(s, "s")
	var duration float64
	fmt.Sscanf(s, "%f", &duration)
	return duration
}

// normalizeLanguageCode は言語コードをGoogle STT形式に正規化する
func (c *GoogleSTTClient) normalizeLanguageCode(language string) string {
	// 短い形式（"ja"）を長い形式（"ja-JP"）に変換
	shortToLong := map[string]string{
		"ja": "ja-JP",
		"en": "en-US",
		"zh": "zh-CN",
		"ru": "ru-RU",
		"es": "es-ES",
		"fr": "fr-FR",
		"de": "de-DE",
		"pt": "pt-BR",
		"it": "it-IT",
		"ko": "ko-KR",
		"ar": "ar-XA",
		"he": "he-IL",
		"fa": "fa-IR",
		"tr": "tr-TR",
		"nl": "nl-NL",
		"pl": "pl-PL",
		"sv": "sv-SE",
		"da": "da-DK",
		"no": "nb-NO",
		"fi": "fi-FI",
		"th": "th-TH",
		"vi": "vi-VN",
		"id": "id-ID",
		"ms": "ms-MY",
		"hi": "hi-IN",
		"bn": "bn-BD",
		"ta": "ta-IN",
		"te": "te-IN",
		"mr": "mr-IN",
		"gu": "gu-IN",
		"uk": "uk-UA",
		"cs": "cs-CZ",
		"el": "el-GR",
		"hu": "hu-HU",
		"ro": "ro-RO",
		"bg": "bg-BG",
		"hr": "hr-HR",
		"sk": "sk-SK",
		"sl": "sl-SI",
		"lt": "lt-LT",
		"lv": "lv-LV",
		"et": "et-EE",
	}

	if long, ok := shortToLong[language]; ok {
		return long
	}
	return language
}

// GoogleSTTSupportedLanguages はGoogle Cloud Speech-to-Textがサポートする言語
var GoogleSTTSupportedLanguages = []string{
	"af-ZA", "am-ET", "ar-AE", "ar-BH", "ar-DZ", "ar-EG", "ar-IL", "ar-IQ", "ar-JO", "ar-KW",
	"ar-LB", "ar-MA", "ar-MR", "ar-OM", "ar-PS", "ar-QA", "ar-SA", "ar-SY", "ar-TN", "ar-YE",
	"az-AZ", "bg-BG", "bn-BD", "bn-IN", "bs-BA", "ca-ES", "cs-CZ", "cy-GB", "da-DK", "de-AT",
	"de-CH", "de-DE", "el-GR", "en-AU", "en-CA", "en-GB", "en-GH", "en-HK", "en-IE", "en-IN",
	"en-KE", "en-NG", "en-NZ", "en-PH", "en-PK", "en-SG", "en-TZ", "en-US", "en-ZA", "es-AR",
	"es-BO", "es-CL", "es-CO", "es-CR", "es-CU", "es-DO", "es-EC", "es-ES", "es-GQ", "es-GT",
	"es-HN", "es-MX", "es-NI", "es-PA", "es-PE", "es-PR", "es-PY", "es-SV", "es-US", "es-UY",
	"es-VE", "et-EE", "eu-ES", "fa-IR", "fi-FI", "fil-PH", "fr-BE", "fr-CA", "fr-CH", "fr-FR",
	"ga-IE", "gl-ES", "gu-IN", "he-IL", "hi-IN", "hr-HR", "hu-HU", "hy-AM", "id-ID", "is-IS",
	"it-CH", "it-IT", "ja-JP", "jv-ID", "ka-GE", "kk-KZ", "km-KH", "kn-IN", "ko-KR", "lo-LA",
	"lt-LT", "lv-LV", "mk-MK", "ml-IN", "mn-MN", "mr-IN", "ms-MY", "mt-MT", "my-MM", "nb-NO",
	"ne-NP", "nl-BE", "nl-NL", "pa-IN", "pl-PL", "ps-AF", "pt-BR", "pt-PT", "ro-RO", "ru-RU",
	"si-LK", "sk-SK", "sl-SI", "so-SO", "sq-AL", "sr-RS", "su-ID", "sv-SE", "sw-KE", "sw-TZ",
	"ta-IN", "ta-LK", "ta-MY", "ta-SG", "te-IN", "th-TH", "tr-TR", "uk-UA", "ur-IN", "ur-PK",
	"uz-UZ", "vi-VN", "wuu-CN", "yue-Hant-HK", "zh-CN", "zh-TW", "zu-ZA",
}

// GetSupportedLanguages はサポートされている言語一覧を取得 (STTClient interface)
func (c *GoogleSTTClient) GetSupportedLanguages(ctx context.Context) ([]string, error) {
	return GoogleSTTSupportedLanguages, nil
}

// GetName はプロバイダー名を返す (STTClient interface)
func (c *GoogleSTTClient) GetName() string {
	return "google"
}

// ============================================================
// LegacySTTClient interface implementation (for backward compatibility)
// ============================================================

// Recognize は音声データをテキストに変換する (LegacySTTClient interface)
func (c *GoogleSTTClient) Recognize(ctx context.Context, audioData []byte, language string) (*models.STTResult, error) {
	if len(audioData) == 0 {
		return nil, fmt.Errorf("音声データが空です")
	}

	// 実際のAPIキーがない場合はモックを返す
	if c.apiKey == "" || os.Getenv("USE_MOCK_APIS") == "true" {
		mockClient := NewMockSTTClient()
		return mockClient.Recognize(ctx, audioData, language)
	}

	// 新しいTranscribeメソッドを使用
	result, err := c.Transcribe(ctx, bytes.NewReader(audioData), language)
	if err != nil {
		return nil, err
	}

	// 結果をmodels.STTResultに変換
	words := make([]models.WordInfo, len(result.Words))
	for i, w := range result.Words {
		words[i] = models.WordInfo{
			Word:       w.Word,
			StartTime:  w.Start,
			EndTime:    w.End,
			Confidence: w.Confidence,
		}
	}

	return &models.STTResult{
		Text:       result.Text,
		Language:   result.Language,
		Confidence: result.Confidence,
		Duration:   float64(result.Duration) / 1000.0,
		Words:      words,
		CreatedAt:  time.Now(),
	}, nil
}

// Ensure GoogleSTTClient implements both interfaces
var _ STTClient = (*GoogleSTTClient)(nil)
var _ LegacySTTClient = (*GoogleSTTClient)(nil)
