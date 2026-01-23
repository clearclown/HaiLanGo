package tts

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
)

// QwenTTSClient はQwen-TTS APIのクライアント
// ローカルGPUベースのTTSサービスを使用
// 対応言語: 中国語、英語、日本語、韓国語、ドイツ語、フランス語、ロシア語、ポルトガル語、スペイン語、イタリア語
type QwenTTSClient struct {
	baseURL string
}

// NewQwenTTSClient は新しいQwen-TTSクライアントを作成する
func NewQwenTTSClient(baseURL string) *QwenTTSClient {
	if baseURL == "" {
		baseURL = "http://localhost:8001"
	}
	// Remove trailing slash
	baseURL = strings.TrimSuffix(baseURL, "/")
	return &QwenTTSClient{
		baseURL: baseURL,
	}
}

// qwenSynthesizeRequest はQwen-TTS APIリクエスト
type qwenSynthesizeRequest struct {
	Text     string `json:"text"`
	Language string `json:"language"`
	Voice    string `json:"voice"`
	Instruct string `json:"instruct,omitempty"`
	Format   string `json:"format"`
}

// qwenSynthesizeResponse はQwen-TTS APIレスポンス
type qwenSynthesizeResponse struct {
	Audio      string `json:"audio"` // base64 encoded
	Format     string `json:"format"`
	SampleRate int    `json:"sample_rate"`
	DurationMs int    `json:"duration_ms"`
}

// qwenVoiceInfo はQwen-TTS音声情報
type qwenVoiceInfo struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Language    string `json:"language"`
	Gender      string `json:"gender"`
	Description string `json:"description"`
	IsNative    bool   `json:"is_native"`
}

// qwenLanguageInfo はQwen-TTS言語情報
type qwenLanguageInfo struct {
	Code       string `json:"code"`
	Name       string `json:"name"`
	NativeName string `json:"native_name"`
}

// Synthesize はテキストを音声に変換する
func (c *QwenTTSClient) Synthesize(ctx context.Context, text string, language string, voice string) (io.ReadCloser, error) {
	// デフォルトの音声を設定
	if voice == "" {
		voice = c.getDefaultVoice(language)
	}

	// 言語コードを変換
	qwenLang := c.mapLanguageCode(language)

	reqBody := qwenSynthesizeRequest{
		Text:     text,
		Language: qwenLang,
		Voice:    voice,
		Format:   "wav",
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", c.baseURL+"/synthesize", bytes.NewReader(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 120 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		resp.Body.Close()
		return nil, fmt.Errorf("Qwen-TTS API error (status %d): %s", resp.StatusCode, string(body))
	}

	// Parse response
	var qwenResp qwenSynthesizeResponse
	if err := json.NewDecoder(resp.Body).Decode(&qwenResp); err != nil {
		resp.Body.Close()
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}
	resp.Body.Close()

	// Decode base64 audio
	audioData, err := base64.StdEncoding.DecodeString(qwenResp.Audio)
	if err != nil {
		return nil, fmt.Errorf("failed to decode audio: %w", err)
	}

	return io.NopCloser(bytes.NewReader(audioData)), nil
}

// mapLanguageCode はHaiLanGoの言語コードをQwen-TTSの言語名にマップする
func (c *QwenTTSClient) mapLanguageCode(code string) string {
	langMap := map[string]string{
		"zh":    "Chinese",
		"zh-CN": "Chinese",
		"zh-TW": "Chinese",
		"en":    "English",
		"en-US": "English",
		"en-GB": "English",
		"ja":    "Japanese",
		"ja-JP": "Japanese",
		"ko":    "Korean",
		"ko-KR": "Korean",
		"de":    "German",
		"de-DE": "German",
		"fr":    "French",
		"fr-FR": "French",
		"ru":    "Russian",
		"ru-RU": "Russian",
		"pt":    "Portuguese",
		"pt-BR": "Portuguese",
		"pt-PT": "Portuguese",
		"es":    "Spanish",
		"es-ES": "Spanish",
		"it":    "Italian",
		"it-IT": "Italian",
	}

	if lang, ok := langMap[code]; ok {
		return lang
	}
	// Try short code
	shortCode := strings.Split(code, "-")[0]
	if lang, ok := langMap[shortCode]; ok {
		return lang
	}
	return "Auto"
}

// getDefaultVoice は言語に対応するデフォルト音声を返す
func (c *QwenTTSClient) getDefaultVoice(language string) string {
	// 言語ごとのネイティブスピーカーをデフォルトに
	voiceMap := map[string]string{
		"zh":    "Vivian",
		"zh-CN": "Vivian",
		"zh-TW": "Vivian",
		"en":    "Ryan",
		"en-US": "Ryan",
		"en-GB": "Aiden",
		"ja":    "Ono_Anna",
		"ja-JP": "Ono_Anna",
		"ko":    "Sohee",
		"ko-KR": "Sohee",
		// 他の言語もVivian（中国語ネイティブ）をデフォルトに
		// Qwen-TTSは全スピーカーが全言語を話せる
		"de":    "Ryan",
		"de-DE": "Ryan",
		"fr":    "Serena",
		"fr-FR": "Serena",
		"ru":    "Ryan",
		"ru-RU": "Ryan",
		"pt":    "Ryan",
		"pt-BR": "Ryan",
		"es":    "Ryan",
		"es-ES": "Ryan",
		"it":    "Ryan",
		"it-IT": "Ryan",
	}

	if voice, ok := voiceMap[language]; ok {
		return voice
	}
	return "Ryan" // デフォルト
}

// GetVoices は指定言語で利用可能な音声一覧を取得
func (c *QwenTTSClient) GetVoices(ctx context.Context, language string) ([]Voice, error) {
	url := c.baseURL + "/voices"
	if language != "" {
		url = fmt.Sprintf("%s?language=%s", url, language)
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
		return nil, fmt.Errorf("Qwen-TTS voices API error (status %d): %s", resp.StatusCode, string(body))
	}

	var result struct {
		Voices []qwenVoiceInfo `json:"voices"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	voices := make([]Voice, len(result.Voices))
	for i, v := range result.Voices {
		voices[i] = Voice{
			ID:          v.ID,
			Name:        v.Name,
			Language:    v.Language,
			Gender:      v.Gender,
			Type:        "neural",
			Description: v.Description,
		}
	}

	return voices, nil
}

// GetSupportedLanguages はサポートされている言語一覧を取得
func (c *QwenTTSClient) GetSupportedLanguages(ctx context.Context) ([]string, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", c.baseURL+"/languages", nil)
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
		return nil, fmt.Errorf("Qwen-TTS languages API error (status %d): %s", resp.StatusCode, string(body))
	}

	var result struct {
		Languages []qwenLanguageInfo `json:"languages"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	languages := make([]string, len(result.Languages))
	for i, l := range result.Languages {
		languages[i] = l.Code
	}

	return languages, nil
}

// GetName はプロバイダー名を返す
func (c *QwenTTSClient) GetName() string {
	return "qwen"
}

// IsAvailable はQwen-TTSサービスが利用可能かチェックする
func (c *QwenTTSClient) IsAvailable(ctx context.Context) bool {
	req, err := http.NewRequestWithContext(ctx, "GET", c.baseURL+"/health", nil)
	if err != nil {
		return false
	}

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	return resp.StatusCode == http.StatusOK
}

// NewQwenTTSClientFromEnv は環境変数からQwen-TTSクライアントを作成する
func NewQwenTTSClientFromEnv() *QwenTTSClient {
	baseURL := os.Getenv("QWEN_TTS_URL")
	if baseURL == "" {
		baseURL = "http://localhost:8001"
	}
	return NewQwenTTSClient(baseURL)
}

// Ensure QwenTTSClient implements TTSClient interface
var _ TTSClient = (*QwenTTSClient)(nil)
