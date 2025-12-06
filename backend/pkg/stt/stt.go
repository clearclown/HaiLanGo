package stt

import (
	"context"
	"io"
)

// STTClient はSTT (Speech-to-Text) APIのインターフェース
type STTClient interface {
	// Transcribe は音声をテキストに変換する
	Transcribe(ctx context.Context, audio io.Reader, language string) (*TranscriptionResult, error)

	// GetSupportedLanguages はサポートされている言語一覧を取得
	GetSupportedLanguages(ctx context.Context) ([]string, error)

	// GetName はプロバイダー名を返す
	GetName() string
}

// PronunciationEvaluator は発音評価のインターフェース
// HaiLanGoの価値提案: マイナー言語でも発音評価を提供
type PronunciationEvaluator interface {
	// EvaluatePronunciation は発音を評価する
	// Whisper + LLM方式でマイナー言語も対応可能
	EvaluatePronunciation(ctx context.Context, audio io.Reader, expectedText string, language string) (*PronunciationResult, error)
}

// TranscriptionResult は音声認識結果
type TranscriptionResult struct {
	Text             string              `json:"text"`               // 認識テキスト
	Language         string              `json:"language"`           // 検出/指定言語
	Confidence       float64             `json:"confidence"`         // 信頼度 (0.0-1.0)
	Duration         int64               `json:"duration_ms"`        // 音声長（ミリ秒）
	Words            []WordTiming        `json:"words,omitempty"`    // 単語タイミング
	Segments         []TranscriptSegment `json:"segments,omitempty"` // セグメント
	Provider         string              `json:"provider"`           // 使用プロバイダー
	ProcessingTimeMs int64               `json:"processing_time_ms"` // 処理時間
}

// WordTiming は単語ごとのタイミング情報
type WordTiming struct {
	Word       string  `json:"word"`
	Start      float64 `json:"start"`      // 開始時間（秒）
	End        float64 `json:"end"`        // 終了時間（秒）
	Confidence float64 `json:"confidence"` // 信頼度
}

// TranscriptSegment はセグメント情報
type TranscriptSegment struct {
	ID         int     `json:"id"`
	Start      float64 `json:"start"`
	End        float64 `json:"end"`
	Text       string  `json:"text"`
	Confidence float64 `json:"confidence"`
}

// PronunciationResult は発音評価結果
// HaiLanGoの価値提案: ドメイン特化の評価が可能
type PronunciationResult struct {
	OverallScore     float64           `json:"overall_score"`      // 総合スコア (0-100)
	AccuracyScore    float64           `json:"accuracy_score"`     // 正確性スコア
	FluencyScore     float64           `json:"fluency_score"`      // 流暢性スコア
	ProsodyScore     float64           `json:"prosody_score"`      // 韻律スコア
	RecognizedText   string            `json:"recognized_text"`    // 認識されたテキスト
	ExpectedText     string            `json:"expected_text"`      // 期待されたテキスト
	WordScores       []WordScore       `json:"word_scores"`        // 単語別スコア
	Feedback         string            `json:"feedback"`           // フィードバック
	ImprovementTips  []string          `json:"improvement_tips"`   // 改善アドバイス
	DomainNotes      string            `json:"domain_notes"`       // ドメイン特化コメント
	EvaluationMethod string            `json:"evaluation_method"`  // 評価方法 (azure_native, whisper_llm)
}

// WordScore は単語ごとの発音スコア
type WordScore struct {
	Word           string   `json:"word"`
	Score          float64  `json:"score"`
	PhonemeScores  []Phoneme `json:"phoneme_scores,omitempty"`
	ErrorType      string   `json:"error_type,omitempty"` // mispronunciation, omission, insertion
	Suggestion     string   `json:"suggestion,omitempty"`
}

// Phoneme は音素情報
type Phoneme struct {
	Symbol string  `json:"symbol"`
	Score  float64 `json:"score"`
}

// STTProvider はSTTプロバイダーの種類
type STTProvider string

const (
	// ProviderWhisper はOpenAI Whisper API（推奨: 99言語対応、マイナー言語最強）
	ProviderWhisper STTProvider = "whisper"

	// ProviderWhisperLocal はローカルWhisper（whisper.cpp）
	ProviderWhisperLocal STTProvider = "whisper_local"

	// ProviderAzureSpeech はAzure Speech Services
	ProviderAzureSpeech STTProvider = "azure"

	// ProviderGoogleSTT はGoogle Cloud Speech-to-Text
	ProviderGoogleSTT STTProvider = "google"

	// ProviderDeepgram はDeepgram（リアルタイム特化）
	ProviderDeepgram STTProvider = "deepgram"
)

// PronunciationEvalMethod は発音評価の方法
type PronunciationEvalMethod string

const (
	// EvalMethodWhisperLLM はWhisper + LLMによる評価（マイナー言語対応）
	EvalMethodWhisperLLM PronunciationEvalMethod = "whisper_llm"

	// EvalMethodAzureNative はAzure Speech発音評価（主要言語高精度）
	EvalMethodAzureNative PronunciationEvalMethod = "azure_native"
)

// WhisperSupportedLanguages はWhisperがサポートする99言語
// HaiLanGoの価値提案: マイナー言語サポートの根幹
var WhisperSupportedLanguages = []string{
	// 主要言語
	"en", "zh", "de", "es", "ru", "ko", "fr", "ja", "pt", "tr",
	"pl", "ca", "nl", "ar", "sv", "it", "id", "hi", "fi", "vi",
	"he", "uk", "el", "ms", "cs", "ro", "da", "hu", "ta", "no",
	"th", "ur", "hr", "bg", "lt", "la", "mi", "ml", "cy", "sk",
	"te", "fa", "lv", "bn", "sr", "az", "sl", "kn", "et", "mk",
	"br", "eu", "is", "hy", "ne", "mn", "bs", "kk", "sq", "sw",
	"gl", "mr", "pa", "si", "km", "sn", "yo", "so", "af", "oc",
	"ka", "be", "tg", "sd", "gu", "am", "yi", "lo", "uz", "fo",
	"ht", "ps", "tk", "nn", "mt", "sa", "lb", "my", "bo", "tl",
	"mg", "as", "tt", "haw", "ln", "ha", "ba", "jw", "su",
	// クルド語も対応（ku）
}

// GetRecommendedSTTProvider は言語に基づいて推奨プロバイダーを返す
func GetRecommendedSTTProvider(languageCode string) STTProvider {
	// Whisperは99言語対応で、マイナー言語も高精度
	// Azure/Googleは主要言語でリアルタイム性が優れる

	majorLanguages := map[string]bool{
		"en": true, "ja": true, "zh": true, "es": true,
		"fr": true, "de": true, "pt": true, "it": true,
		"ko": true, "ru": true,
	}

	if majorLanguages[languageCode] {
		// 主要言語: Azure（リアルタイム発音評価あり）
		return ProviderAzureSpeech
	}

	// マイナー言語: Whisper（99言語対応）
	return ProviderWhisper
}

// GetRecommendedEvalMethod は言語に基づいて推奨発音評価方法を返す
func GetRecommendedEvalMethod(languageCode string) PronunciationEvalMethod {
	// Azure発音評価は主要言語のみ対応
	azureSupportedLanguages := map[string]bool{
		"en": true, "zh": true, "de": true, "es": true,
		"fr": true, "ja": true, "pt": true, "ko": true,
	}

	if azureSupportedLanguages[languageCode] {
		return EvalMethodAzureNative
	}

	// マイナー言語: Whisper + LLMによる評価
	return EvalMethodWhisperLLM
}
