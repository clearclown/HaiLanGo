package tts

import (
	"bytes"
	"context"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
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
func (c *AzureTTSClient) Synthesize(ctx context.Context, text string, language string, voice string) (io.ReadCloser, error) {
	// デフォルトの音声を設定
	if voice == "" {
		voice = c.getDefaultVoice(language)
	}

	// SSMLを構築
	ssml := c.buildSSML(text, language, voice)

	// Azure Speech Services TTS APIエンドポイント
	url := fmt.Sprintf("https://%s.tts.speech.microsoft.com/cognitiveservices/v1", c.region)

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader([]byte(ssml)))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Ocp-Apim-Subscription-Key", c.subscriptionKey)
	req.Header.Set("Content-Type", "application/ssml+xml")
	req.Header.Set("X-Microsoft-OutputFormat", "audio-16khz-128kbitrate-mono-mp3")
	req.Header.Set("User-Agent", "HaiLanGo")

	client := &http.Client{Timeout: 60 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		resp.Body.Close()
		return nil, fmt.Errorf("Azure TTS API error (status %d): %s", resp.StatusCode, string(body))
	}

	return resp.Body, nil
}

// buildSSML はSSMLドキュメントを構築する
func (c *AzureTTSClient) buildSSML(text string, language string, voice string) string {
	// XMLエスケープ
	escapedText := escapeXML(text)

	return fmt.Sprintf(`<speak version='1.0' xmlns='http://www.w3.org/2001/10/synthesis' xml:lang='%s'>
    <voice name='%s'>
        %s
    </voice>
</speak>`, language, voice, escapedText)
}

// escapeXML はXML特殊文字をエスケープする
func escapeXML(s string) string {
	var buf bytes.Buffer
	if err := xml.EscapeText(&buf, []byte(s)); err != nil {
		return s
	}
	return buf.String()
}

// getDefaultVoice は言語に対応するデフォルト音声を返す
func (c *AzureTTSClient) getDefaultVoice(language string) string {
	// Azure Neural Voices（高品質）
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
	}

	if voice, ok := voiceMap[language]; ok {
		return voice
	}
	// デフォルトは英語
	return "en-US-JennyNeural"
}

// azureVoiceResponse はAzure Voices APIのレスポンス構造
type azureVoiceResponse struct {
	Name            string   `json:"Name"`
	DisplayName     string   `json:"DisplayName"`
	LocalName       string   `json:"LocalName"`
	ShortName       string   `json:"ShortName"`
	Gender          string   `json:"Gender"`
	Locale          string   `json:"Locale"`
	LocaleName      string   `json:"LocaleName"`
	StyleList       []string `json:"StyleList,omitempty"`
	VoiceType       string   `json:"VoiceType"`
	Status          string   `json:"Status"`
	WordsPerMinute  string   `json:"WordsPerMinute,omitempty"`
}

// GetVoices は指定言語で利用可能な音声一覧を取得
func (c *AzureTTSClient) GetVoices(ctx context.Context, language string) ([]Voice, error) {
	// Azure Speech Services Voices APIエンドポイント
	url := fmt.Sprintf("https://%s.tts.speech.microsoft.com/cognitiveservices/voices/list", c.region)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Ocp-Apim-Subscription-Key", c.subscriptionKey)

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("Azure Voices API error (status %d): %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var azureVoices []azureVoiceResponse
	if err := json.Unmarshal(body, &azureVoices); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	// 指定言語でフィルタリング
	var voices []Voice
	for _, av := range azureVoices {
		// 言語コードでフィルタリング（例: "ja-JP" または "ja"）
		if language != "" && !strings.HasPrefix(av.Locale, language) {
			continue
		}

		// Neural voicesのみを返す（高品質）
		if av.VoiceType != "Neural" {
			continue
		}

		voices = append(voices, Voice{
			ID:       av.ShortName,
			Name:     av.DisplayName,
			Language: av.Locale,
			Gender:   av.Gender,
		})
	}

	return voices, nil
}

// GetSupportedLanguages はサポートされている言語一覧を取得
func (c *AzureTTSClient) GetSupportedLanguages(ctx context.Context) ([]string, error) {
	// すべての音声を取得
	voices, err := c.GetVoices(ctx, "")
	if err != nil {
		return nil, fmt.Errorf("failed to get voices: %w", err)
	}

	// ユニークな言語コードを抽出
	langMap := make(map[string]bool)
	for _, v := range voices {
		// "ja-JP" から "ja-JP" を取得（フル言語コード）
		langMap[v.Language] = true
	}

	// マップをスライスに変換
	languages := make([]string, 0, len(langMap))
	for lang := range langMap {
		languages = append(languages, lang)
	}

	return languages, nil
}

// GetName はプロバイダー名を返す
func (c *AzureTTSClient) GetName() string {
	return "azure"
}

// Ensure AzureTTSClient implements TTSClient interface
var _ TTSClient = (*AzureTTSClient)(nil)
