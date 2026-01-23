package tts

import (
	"context"
	"io"
)

// TTSClient はTTS APIのインターフェース
type TTSClient interface {
	// Synthesize はテキストを音声に変換する
	Synthesize(ctx context.Context, text string, language string, voice string) (io.ReadCloser, error)

	// GetVoices は指定言語で利用可能な音声一覧を取得
	GetVoices(ctx context.Context, language string) ([]Voice, error)

	// GetSupportedLanguages はサポートされている言語一覧を取得
	GetSupportedLanguages(ctx context.Context) ([]string, error)

	// GetName はプロバイダー名を返す
	GetName() string
}

// Voice は音声情報を表す
type Voice struct {
	ID          string   `json:"id"`           // 音声ID
	Name        string   `json:"name"`         // 表示名
	Language    string   `json:"language"`     // 言語コード
	Gender      string   `json:"gender"`       // male, female, neutral
	Type        string   `json:"type"`         // standard, neural, premium
	SampleRate  int      `json:"sample_rate"`  // サンプリングレート
	Description string   `json:"description"`  // 説明
	Styles      []string `json:"styles"`       // 対応スタイル (cheerful, sad等)
}

// SynthesizeOptions は音声合成のオプション
type SynthesizeOptions struct {
	Speed      float64 `json:"speed"`       // 速度 (0.5-2.0)
	Pitch      float64 `json:"pitch"`       // ピッチ (-20.0 to 20.0)
	VolumeGain float64 `json:"volume_gain"` // 音量 (-96.0 to 16.0)
	Style      string  `json:"style"`       // 音声スタイル
	Format     string  `json:"format"`      // 出力形式 (mp3, wav, ogg)
}

// TTSProvider はTTSプロバイダーの種類
type TTSProvider string

const (
	ProviderAzureTTS    TTSProvider = "azure"      // Azure Speech Services（推奨: 140言語）
	ProviderGoogleTTS   TTSProvider = "google"     // Google Cloud TTS
	ProviderElevenLabs  TTSProvider = "elevenlabs" // ElevenLabs（プレミアム）
	ProviderEdgeTTS     TTSProvider = "edge"       // edge-tts（無料、50言語）
	ProviderCoquiTTS    TTSProvider = "coqui"      // Coqui TTS（オープンソース）
	ProviderQwenTTS     TTSProvider = "qwen"       // Qwen-TTS（ローカルGPU、10言語、無料）
)

// TTSResult は音声合成結果
type TTSResult struct {
	Audio       []byte `json:"audio"`        // 音声データ
	Format      string `json:"format"`       // 形式 (mp3, wav等)
	Duration    int64  `json:"duration_ms"`  // 長さ（ミリ秒）
	CharCount   int    `json:"char_count"`   // 文字数
	Provider    string `json:"provider"`     // 使用したプロバイダー
	Voice       string `json:"voice"`        // 使用した音声
	CacheKey    string `json:"cache_key"`    // キャッシュキー
	IsCached    bool   `json:"is_cached"`    // キャッシュから取得したか
}

// LanguageCoverage は言語カバレッジ情報
type LanguageCoverage struct {
	Azure      bool `json:"azure"`      // Azure対応
	Google     bool `json:"google"`     // Google対応
	ElevenLabs bool `json:"elevenlabs"` // ElevenLabs対応
	Edge       bool `json:"edge"`       // edge-tts対応
}

// GetRecommendedProvider は言語に基づいて推奨プロバイダーを返す
// HaiLanGoの価値提案: マイナー言語サポートを最大化
func GetRecommendedProvider(languageCode string) TTSProvider {
	// マイナー言語はAzureが最もカバレッジが広い（140言語）
	// 主要言語ではElevenLabsがプレミアム品質

	majorLanguages := map[string]bool{
		"en": true, "ja": true, "zh": true, "es": true,
		"fr": true, "de": true, "pt": true, "it": true,
		"ko": true, "ru": true,
	}

	if majorLanguages[languageCode] {
		// 主要言語: プレミアムプランではElevenLabs、通常はAzure
		return ProviderAzureTTS
	}

	// マイナー言語: Azureが最もカバレッジが広い
	return ProviderAzureTTS
}
