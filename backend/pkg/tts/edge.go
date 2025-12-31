package tts

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// EdgeTTSClient はMicrosoft Edge TTSクライアント（無料、APIキー不要）
// edge-tts Pythonライブラリと同等の機能を提供
type EdgeTTSClient struct {
	// voicesはキャッシュされた音声リスト
	voices     []Voice
	voicesMu   sync.RWMutex
	voicesOnce sync.Once
}

// NewEdgeTTSClient は新しいEdge TTSクライアントを作成する
func NewEdgeTTSClient() *EdgeTTSClient {
	return &EdgeTTSClient{}
}

const (
	edgeTTSWSURL = "wss://speech.platform.bing.com/consumer/speech/synthesize/readaloud/edge/v1"
	edgeTTSVoiceListURL = "https://speech.platform.bing.com/consumer/speech/synthesize/readaloud/voices/list"
	trustedClientToken = "6A5AA1D4EAFF4E9FB37E23D68491D6F4"
)

// Synthesize はテキストを音声に変換する
func (c *EdgeTTSClient) Synthesize(ctx context.Context, text string, language string, voice string) (io.ReadCloser, error) {
	if text == "" {
		return nil, fmt.Errorf("text cannot be empty")
	}

	// デフォルトの音声を設定
	if voice == "" {
		voice = c.getDefaultVoice(language)
	}

	// ランダムなリクエストIDを生成
	requestID := generateRequestID()

	// WebSocket URLを構築
	wsURL := fmt.Sprintf("%s?TrustedClientToken=%s&ConnectionId=%s",
		edgeTTSWSURL, trustedClientToken, requestID)

	// WebSocket接続
	dialer := websocket.Dialer{
		HandshakeTimeout: 10 * time.Second,
	}

	headers := http.Header{
		"User-Agent": []string{"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36 Edg/120.0.0.0"},
		"Origin":     []string{"chrome-extension://jdiccldimpdaibmpdkjnbmckianbfold"},
	}

	conn, _, err := dialer.DialContext(ctx, wsURL, headers)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Edge TTS: %w", err)
	}
	defer conn.Close()

	// 設定メッセージを送信
	configMessage := buildConfigMessage(requestID)
	if err := conn.WriteMessage(websocket.TextMessage, []byte(configMessage)); err != nil {
		return nil, fmt.Errorf("failed to send config message: %w", err)
	}

	// SSMLメッセージを送信
	ssmlMessage := buildSSMLMessage(requestID, text, voice, language)
	if err := conn.WriteMessage(websocket.TextMessage, []byte(ssmlMessage)); err != nil {
		return nil, fmt.Errorf("failed to send SSML message: %w", err)
	}

	// 音声データを受信
	var audioBuffer bytes.Buffer
	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		messageType, message, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway) {
				break
			}
			return nil, fmt.Errorf("failed to read message: %w", err)
		}

		if messageType == websocket.TextMessage {
			msgStr := string(message)
			// 終了メッセージを確認
			if strings.Contains(msgStr, "turn.end") {
				break
			}
		} else if messageType == websocket.BinaryMessage {
			// バイナリメッセージから音声データを抽出
			audioData := extractAudioData(message)
			if len(audioData) > 0 {
				audioBuffer.Write(audioData)
			}
		}
	}

	if audioBuffer.Len() == 0 {
		return nil, fmt.Errorf("no audio data received")
	}

	return io.NopCloser(bytes.NewReader(audioBuffer.Bytes())), nil
}

// generateRequestID はランダムなリクエストIDを生成する
func generateRequestID() string {
	bytes := make([]byte, 16)
	rand.Read(bytes)
	return strings.ToUpper(hex.EncodeToString(bytes))
}

// buildConfigMessage は設定メッセージを構築する
func buildConfigMessage(requestID string) string {
	timestamp := time.Now().UTC().Format("2006-01-02T15:04:05.000Z")
	return fmt.Sprintf(`X-Timestamp:%s
Content-Type:application/json; charset=utf-8
Path:speech.config

{"context":{"synthesis":{"audio":{"metadataoptions":{"sentenceBoundaryEnabled":"false","wordBoundaryEnabled":"false"},"outputFormat":"audio-24khz-48kbitrate-mono-mp3"}}}}`, timestamp)
}

// buildSSMLMessage はSSMLメッセージを構築する
func buildSSMLMessage(requestID, text, voice, language string) string {
	timestamp := time.Now().UTC().Format("2006-01-02T15:04:05.000Z")

	// テキストをエスケープ
	escapedText := escapeXML(text)

	// 言語コードを正規化
	if !strings.Contains(language, "-") {
		language = normalizeEdgeLanguage(language)
	}

	ssml := fmt.Sprintf(`<speak version='1.0' xmlns='http://www.w3.org/2001/10/synthesis' xml:lang='%s'><voice name='%s'><prosody pitch='+0Hz' rate='+0%%' volume='+0%%'>%s</prosody></voice></speak>`,
		language, voice, escapedText)

	return fmt.Sprintf(`X-RequestId:%s
Content-Type:application/ssml+xml
X-Timestamp:%s
Path:ssml

%s`, requestID, timestamp, ssml)
}

// normalizeEdgeLanguage は言語コードをEdge TTS形式に正規化する
func normalizeEdgeLanguage(language string) string {
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
		"ar": "ar-SA",
		"he": "he-IL",
		"fa": "fa-IR",
		"tr": "tr-TR",
	}

	if long, ok := shortToLong[language]; ok {
		return long
	}
	return language + "-" + strings.ToUpper(language)
}

// extractAudioData はバイナリメッセージから音声データを抽出する
func extractAudioData(message []byte) []byte {
	// "Path:audio" ヘッダーの後の音声データを抽出
	headerEnd := bytes.Index(message, []byte("\r\n\r\n"))
	if headerEnd == -1 {
		return nil
	}
	return message[headerEnd+4:]
}

// getDefaultVoice は言語に対応するデフォルト音声を返す
func (c *EdgeTTSClient) getDefaultVoice(language string) string {
	// Microsoft Edge Neural Voices
	voiceMap := map[string]string{
		"ja":    "ja-JP-NanamiNeural",
		"ja-JP": "ja-JP-NanamiNeural",
		"en":    "en-US-JennyNeural",
		"en-US": "en-US-JennyNeural",
		"en-GB": "en-GB-SoniaNeural",
		"zh":    "zh-CN-XiaoxiaoNeural",
		"zh-CN": "zh-CN-XiaoxiaoNeural",
		"zh-TW": "zh-TW-HsiaoChenNeural",
		"ru":    "ru-RU-SvetlanaNeural",
		"ru-RU": "ru-RU-SvetlanaNeural",
		"es":    "es-ES-ElviraNeural",
		"es-ES": "es-ES-ElviraNeural",
		"fr":    "fr-FR-DeniseNeural",
		"fr-FR": "fr-FR-DeniseNeural",
		"de":    "de-DE-KatjaNeural",
		"de-DE": "de-DE-KatjaNeural",
		"pt":    "pt-BR-FranciscaNeural",
		"pt-BR": "pt-BR-FranciscaNeural",
		"it":    "it-IT-ElsaNeural",
		"it-IT": "it-IT-ElsaNeural",
		"tr":    "tr-TR-EmelNeural",
		"tr-TR": "tr-TR-EmelNeural",
		"fa":    "fa-IR-DilaraNeural",
		"fa-IR": "fa-IR-DilaraNeural",
		"he":    "he-IL-HilaNeural",
		"he-IL": "he-IL-HilaNeural",
		"ar":    "ar-SA-ZariyahNeural",
		"ar-SA": "ar-SA-ZariyahNeural",
		"ko":    "ko-KR-SunHiNeural",
		"ko-KR": "ko-KR-SunHiNeural",
		"hi":    "hi-IN-SwaraNeural",
		"hi-IN": "hi-IN-SwaraNeural",
		"th":    "th-TH-PremwadeeNeural",
		"th-TH": "th-TH-PremwadeeNeural",
		"vi":    "vi-VN-HoaiMyNeural",
		"vi-VN": "vi-VN-HoaiMyNeural",
		"id":    "id-ID-GadisNeural",
		"id-ID": "id-ID-GadisNeural",
	}

	if voice, ok := voiceMap[language]; ok {
		return voice
	}
	// デフォルトは英語
	return "en-US-JennyNeural"
}

// edgeVoiceResponse はEdge TTS Voices APIのレスポンス構造
type edgeVoiceResponse struct {
	Name           string `json:"Name"`
	ShortName      string `json:"ShortName"`
	Gender         string `json:"Gender"`
	Locale         string `json:"Locale"`
	SuggestedCodec string `json:"SuggestedCodec"`
	FriendlyName   string `json:"FriendlyName"`
	Status         string `json:"Status"`
}

// GetVoices は指定言語で利用可能な音声一覧を取得
func (c *EdgeTTSClient) GetVoices(ctx context.Context, language string) ([]Voice, error) {
	// 音声リストを取得（キャッシュがあれば使用）
	c.voicesMu.RLock()
	if len(c.voices) > 0 {
		voices := c.filterVoicesByLanguage(c.voices, language)
		c.voicesMu.RUnlock()
		return voices, nil
	}
	c.voicesMu.RUnlock()

	// APIから音声リストを取得
	url := fmt.Sprintf("%s?trustedclienttoken=%s", edgeTTSVoiceListURL, trustedClientToken)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36")

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("Edge Voices API error (status %d): %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var edgeVoices []edgeVoiceResponse
	// レスポンスはプレーンなJSON配列
	if err := parseEdgeVoicesResponse(body, &edgeVoices); err != nil {
		return nil, fmt.Errorf("failed to parse voices: %w", err)
	}

	// Voiceに変換
	var voices []Voice
	for _, ev := range edgeVoices {
		// Neuralボイスのみを含める
		if !strings.Contains(ev.ShortName, "Neural") {
			continue
		}

		voices = append(voices, Voice{
			ID:       ev.ShortName,
			Name:     ev.FriendlyName,
			Language: ev.Locale,
			Gender:   strings.ToLower(ev.Gender),
			Type:     "neural",
		})
	}

	// キャッシュに保存
	c.voicesMu.Lock()
	c.voices = voices
	c.voicesMu.Unlock()

	return c.filterVoicesByLanguage(voices, language), nil
}

// parseEdgeVoicesResponse はEdge TTSの音声リストレスポンスをパースする
func parseEdgeVoicesResponse(body []byte, voices *[]edgeVoiceResponse) error {
	// JSONレスポンスをパース
	// Edge TTSはJSON配列を返す
	re := regexp.MustCompile(`\{[^{}]+\}`)
	matches := re.FindAll(body, -1)

	for _, match := range matches {
		var voice edgeVoiceResponse
		if err := parseVoiceJSON(match, &voice); err != nil {
			continue
		}
		if voice.ShortName != "" {
			*voices = append(*voices, voice)
		}
	}

	return nil
}

// parseVoiceJSON は単一の音声JSONをパースする
func parseVoiceJSON(data []byte, voice *edgeVoiceResponse) error {
	// シンプルなJSONパース
	str := string(data)

	voice.Name = extractJSONValue(str, "Name")
	voice.ShortName = extractJSONValue(str, "ShortName")
	voice.Gender = extractJSONValue(str, "Gender")
	voice.Locale = extractJSONValue(str, "Locale")
	voice.FriendlyName = extractJSONValue(str, "FriendlyName")
	voice.Status = extractJSONValue(str, "Status")

	return nil
}

// extractJSONValue はJSON文字列から値を抽出する
func extractJSONValue(json, key string) string {
	pattern := fmt.Sprintf(`"%s"\s*:\s*"([^"]*)"`, key)
	re := regexp.MustCompile(pattern)
	matches := re.FindStringSubmatch(json)
	if len(matches) >= 2 {
		return matches[1]
	}
	return ""
}

// filterVoicesByLanguage は言語でフィルタリングする
func (c *EdgeTTSClient) filterVoicesByLanguage(voices []Voice, language string) []Voice {
	if language == "" {
		return voices
	}

	var filtered []Voice
	for _, v := range voices {
		if strings.HasPrefix(v.Language, language) {
			filtered = append(filtered, v)
		}
	}
	return filtered
}

// GetSupportedLanguages はサポートされている言語一覧を取得
func (c *EdgeTTSClient) GetSupportedLanguages(ctx context.Context) ([]string, error) {
	voices, err := c.GetVoices(ctx, "")
	if err != nil {
		return nil, fmt.Errorf("failed to get voices: %w", err)
	}

	langMap := make(map[string]bool)
	for _, v := range voices {
		langMap[v.Language] = true
	}

	languages := make([]string, 0, len(langMap))
	for lang := range langMap {
		languages = append(languages, lang)
	}

	return languages, nil
}

// GetName はプロバイダー名を返す
func (c *EdgeTTSClient) GetName() string {
	return "edge"
}

// Ensure EdgeTTSClient implements TTSClient interface
var _ TTSClient = (*EdgeTTSClient)(nil)
