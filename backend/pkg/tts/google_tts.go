package tts

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/clearclown/HaiLanGo/backend/pkg/language"
)

// LegacyTTSClient は旧TTS APIクライアントのインターフェース（後方互換性のため保持）
// 新規実装はTTSClient (tts.go) を使用してください
type LegacyTTSClient interface {
	Generate(ctx context.Context, text string, lang string, quality string, speed float64) ([]byte, error)
}

// GoogleTTSClient はGoogle Cloud TTSクライアント
// TTSClient インターフェースと LegacyTTSClient インターフェースの両方を実装
type GoogleTTSClient struct {
	apiKey string
}

// NewGoogleTTSClient は新しいGoogle Cloud TTSクライアントを作成
func NewGoogleTTSClient(apiKey string) *GoogleTTSClient {
	return &GoogleTTSClient{
		apiKey: apiKey,
	}
}

// googleTTSRequest はGoogle Cloud TTS APIのリクエスト構造
type googleTTSRequest struct {
	Input       googleTTSInput       `json:"input"`
	Voice       googleTTSVoice       `json:"voice"`
	AudioConfig googleTTSAudioConfig `json:"audioConfig"`
}

type googleTTSInput struct {
	Text string `json:"text,omitempty"`
	SSML string `json:"ssml,omitempty"`
}

type googleTTSVoice struct {
	LanguageCode string `json:"languageCode"`
	Name         string `json:"name,omitempty"`
	SsmlGender   string `json:"ssmlGender,omitempty"`
}

type googleTTSAudioConfig struct {
	AudioEncoding   string  `json:"audioEncoding"`
	SpeakingRate    float64 `json:"speakingRate,omitempty"`
	Pitch           float64 `json:"pitch,omitempty"`
	VolumeGainDb    float64 `json:"volumeGainDb,omitempty"`
	SampleRateHertz int     `json:"sampleRateHertz,omitempty"`
}

// googleTTSResponse はGoogle Cloud TTS APIのレスポンス構造
type googleTTSResponse struct {
	AudioContent string `json:"audioContent"`
}

// Synthesize はテキストを音声に変換する（TTSClient インターフェース実装）
func (c *GoogleTTSClient) Synthesize(ctx context.Context, text string, lang string, voice string) (io.ReadCloser, error) {
	if text == "" {
		return nil, errors.New("text cannot be empty")
	}

	// デフォルトの音声を設定
	if voice == "" {
		voice = c.getDefaultVoice(lang)
	}

	// 言語コードを正規化
	languageCode := c.normalizeLanguageCode(lang)

	// リクエストボディを構築
	reqBody := googleTTSRequest{
		Input: googleTTSInput{
			Text: text,
		},
		Voice: googleTTSVoice{
			LanguageCode: languageCode,
			Name:         voice,
		},
		AudioConfig: googleTTSAudioConfig{
			AudioEncoding:   "MP3",
			SpeakingRate:    1.0,
			Pitch:           0.0,
			SampleRateHertz: 24000,
		},
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	// Google Cloud TTS APIエンドポイント
	url := fmt.Sprintf("https://texttospeech.googleapis.com/v1/text:synthesize?key=%s", c.apiKey)

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 60 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		resp.Body.Close()
		return nil, fmt.Errorf("Google TTS API error (status %d): %s", resp.StatusCode, string(body))
	}

	// レスポンスをパース
	respBody, err := io.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var ttsResp googleTTSResponse
	if err := json.Unmarshal(respBody, &ttsResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	// Base64デコード
	audioData, err := base64.StdEncoding.DecodeString(ttsResp.AudioContent)
	if err != nil {
		return nil, fmt.Errorf("failed to decode audio content: %w", err)
	}

	return io.NopCloser(bytes.NewReader(audioData)), nil
}

// normalizeLanguageCode は言語コードをGoogle TTS形式に正規化する
func (c *GoogleTTSClient) normalizeLanguageCode(lang string) string {
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
		"bn": "bn-IN",
		"ta": "ta-IN",
		"te": "te-IN",
		"mr": "mr-IN",
		"gu": "gu-IN",
	}

	if long, ok := shortToLong[lang]; ok {
		return long
	}
	return lang
}

// getDefaultVoice は言語に対応するデフォルト音声を返す
func (c *GoogleTTSClient) getDefaultVoice(lang string) string {
	// Google Cloud TTS WaveNet/Neural2 voices（高品質）
	voiceMap := map[string]string{
		"ja":    "ja-JP-Neural2-B",
		"ja-JP": "ja-JP-Neural2-B",
		"en":    "en-US-Neural2-F",
		"en-US": "en-US-Neural2-F",
		"en-GB": "en-GB-Neural2-F",
		"zh":    "zh-CN-Neural2-A",
		"zh-CN": "zh-CN-Neural2-A",
		"zh-TW": "zh-TW-Neural2-A",
		"ru":    "ru-RU-Neural2-A",
		"ru-RU": "ru-RU-Neural2-A",
		"es":    "es-ES-Neural2-A",
		"es-ES": "es-ES-Neural2-A",
		"fr":    "fr-FR-Neural2-A",
		"fr-FR": "fr-FR-Neural2-A",
		"de":    "de-DE-Neural2-A",
		"de-DE": "de-DE-Neural2-A",
		"pt":    "pt-BR-Neural2-A",
		"pt-BR": "pt-BR-Neural2-A",
		"it":    "it-IT-Neural2-A",
		"it-IT": "it-IT-Neural2-A",
		"tr":    "tr-TR-Neural2-A",
		"tr-TR": "tr-TR-Neural2-A",
		"ko":    "ko-KR-Neural2-A",
		"ko-KR": "ko-KR-Neural2-A",
		"ar":    "ar-XA-Neural2-A",
		"ar-XA": "ar-XA-Neural2-A",
		"he":    "he-IL-Neural2-A",
		"he-IL": "he-IL-Neural2-A",
		"hi":    "hi-IN-Neural2-A",
		"hi-IN": "hi-IN-Neural2-A",
		"th":    "th-TH-Neural2-C",
		"th-TH": "th-TH-Neural2-C",
		"vi":    "vi-VN-Neural2-A",
		"vi-VN": "vi-VN-Neural2-A",
		"id":    "id-ID-Neural2-A",
		"id-ID": "id-ID-Neural2-A",
		"nl":    "nl-NL-Neural2-A",
		"nl-NL": "nl-NL-Neural2-A",
		"pl":    "pl-PL-Neural2-A",
		"pl-PL": "pl-PL-Neural2-A",
		"sv":    "sv-SE-Neural2-A",
		"sv-SE": "sv-SE-Neural2-A",
		"da":    "da-DK-Neural2-D",
		"da-DK": "da-DK-Neural2-D",
		"no":    "nb-NO-Neural2-A",
		"nb-NO": "nb-NO-Neural2-A",
		"fi":    "fi-FI-Neural2-A",
		"fi-FI": "fi-FI-Neural2-A",
	}

	if voice, ok := voiceMap[lang]; ok {
		return voice
	}
	// デフォルトは英語
	return "en-US-Neural2-F"
}

// googleVoiceListResponse はGoogle TTS Voices APIのレスポンス構造
type googleVoiceListResponse struct {
	Voices []googleVoiceInfo `json:"voices"`
}

type googleVoiceInfo struct {
	LanguageCodes          []string `json:"languageCodes"`
	Name                   string   `json:"name"`
	SsmlGender             string   `json:"ssmlGender"`
	NaturalSampleRateHertz int      `json:"naturalSampleRateHertz"`
}

// GetVoices は指定言語で利用可能な音声一覧を取得
func (c *GoogleTTSClient) GetVoices(ctx context.Context, lang string) ([]Voice, error) {
	// Google Cloud TTS Voices APIエンドポイント
	url := fmt.Sprintf("https://texttospeech.googleapis.com/v1/voices?key=%s", c.apiKey)
	if lang != "" {
		url += "&languageCode=" + c.normalizeLanguageCode(lang)
	}

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("Google Voices API error (status %d): %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var voicesResp googleVoiceListResponse
	if err := json.Unmarshal(body, &voicesResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	// 高品質音声のみをフィルタリング（Neural2/WaveNet）
	var voices []Voice
	for _, gv := range voicesResp.Voices {
		// Neural2またはWaveNet音声のみを返す
		if !strings.Contains(gv.Name, "Neural2") && !strings.Contains(gv.Name, "Wavenet") {
			continue
		}

		// 言語コードの最初のものを使用
		voiceLang := ""
		if len(gv.LanguageCodes) > 0 {
			voiceLang = gv.LanguageCodes[0]
		}

		// 性別を正規化
		gender := strings.ToLower(gv.SsmlGender)

		// 音声タイプを判定
		voiceType := "wavenet"
		if strings.Contains(gv.Name, "Neural2") {
			voiceType = "neural"
		}

		voices = append(voices, Voice{
			ID:         gv.Name,
			Name:       gv.Name,
			Language:   voiceLang,
			Gender:     gender,
			Type:       voiceType,
			SampleRate: gv.NaturalSampleRateHertz,
		})
	}

	return voices, nil
}

// GetSupportedLanguages はサポートされている言語一覧を取得
func (c *GoogleTTSClient) GetSupportedLanguages(ctx context.Context) ([]string, error) {
	// すべての音声を取得
	voices, err := c.GetVoices(ctx, "")
	if err != nil {
		return nil, fmt.Errorf("failed to get voices: %w", err)
	}

	// ユニークな言語コードを抽出
	langMap := make(map[string]bool)
	for _, v := range voices {
		langMap[v.Language] = true
	}

	// マップをスライスに変換
	languages := make([]string, 0, len(langMap))
	for voiceLang := range langMap {
		languages = append(languages, voiceLang)
	}

	return languages, nil
}

// GetName はプロバイダー名を返す
func (c *GoogleTTSClient) GetName() string {
	return "google"
}

// =====================================
// Legacy API Support (後方互換性)
// =====================================

// Generate はテキストから音声データを生成（LegacyTTSClient インターフェース実装）
func (c *GoogleTTSClient) Generate(ctx context.Context, text string, lang string, quality string, speed float64) ([]byte, error) {
	// バリデーション
	if err := c.validateLegacy(text, lang, quality, speed); err != nil {
		return nil, err
	}

	// モック環境の場合
	useMock := os.Getenv("USE_MOCK_APIS") == "true" ||
		os.Getenv("TEST_USE_MOCKS") == "true" ||
		c.apiKey == "" ||
		c.apiKey == "mock-api-key"

	if useMock {
		return c.generateMock(text, lang, quality, speed)
	}

	// 実際のGoogle Cloud TTS API呼び出し
	reader, err := c.Synthesize(ctx, text, lang, "")
	if err != nil {
		return nil, err
	}
	defer reader.Close()

	return io.ReadAll(reader)
}

// validateLegacy は入力パラメータの検証（Legacy API用）
func (c *GoogleTTSClient) validateLegacy(text string, lang string, quality string, speed float64) error {
	if text == "" {
		return errors.New("text cannot be empty")
	}

	if speed < 0.5 || speed > 2.0 {
		return fmt.Errorf("speed must be between 0.5 and 2.0, got %.2f", speed)
	}

	if quality != "standard" && quality != "premium" {
		return fmt.Errorf("quality must be 'standard' or 'premium', got '%s'", quality)
	}

	return nil
}

// generateMock はモック音声データを生成
func (c *GoogleTTSClient) generateMock(text string, lang string, quality string, speed float64) ([]byte, error) {
	// モックデータの生成（実際のMP3データの代わりに疑似データを返す）
	// ハッシュを使用して決定論的なデータを生成
	hash := sha256.Sum256([]byte(fmt.Sprintf("%s:%s:%s:%.2f", text, lang, quality, speed)))
	hashStr := hex.EncodeToString(hash[:])

	// 疑似音声データ（実際にはMP3ヘッダーとダミーデータ）
	mockData := []byte(fmt.Sprintf("MOCK_AUDIO_DATA:%s:text=%s:lang=%s:quality=%s:speed=%.2f",
		hashStr[:16], text, lang, quality, speed))

	return mockData, nil
}

// SupportedLanguages は対応言語のリストを返す
// Dynamic language support: returns all registered languages (verified + supported)
// Note: ANY valid language code can be used - experimental languages are allowed
func (c *GoogleTTSClient) SupportedLanguages() []string {
	registry := language.GetRegistry()
	allLangs := registry.GetAll()

	codes := make([]string, 0, len(allLangs))
	for _, lang := range allLangs {
		codes = append(codes, lang.Code)
	}
	return codes
}

// IsLanguageSupported は言語がサポートされているかチェック
// Dynamic language support: ANY valid language code is considered supported
// The underlying LLM/TTS API will handle it if it's truly supported
func (c *GoogleTTSClient) IsLanguageSupported(lang string) bool {
	// With dynamic language support, any valid language code is allowed
	// The registry will return an experimental entry for unknown languages
	return language.IsValidCode(lang)
}

// Ensure GoogleTTSClient implements TTSClient interface
var _ TTSClient = (*GoogleTTSClient)(nil)
