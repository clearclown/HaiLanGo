package stt

import (
	"context"
	"errors"
	"io"
	"strings"
	"time"
)

// MockSTTClientNew は新インターフェース用モックSTTクライアント
type MockSTTClientNew struct {
	dataDir   string
	mockAudio []byte
}

// NewMockSTTClientNew は新しいモックSTTクライアントを作成する
func NewMockSTTClientNew() *MockSTTClientNew {
	return &MockSTTClientNew{
		dataDir:   "./mocks/data",
		mockAudio: nil,
	}
}

// Transcribe は音声をテキストに変換する（モック）
func (m *MockSTTClientNew) Transcribe(ctx context.Context, audio io.Reader, language string) (*TranscriptionResult, error) {
	if audio == nil {
		return nil, errors.New("audio cannot be nil")
	}

	// 音声データを読み取り
	audioData, err := io.ReadAll(audio)
	if err != nil {
		return nil, err
	}

	if len(audioData) == 0 {
		return nil, errors.New("audio data cannot be empty")
	}

	// 言語ごとのモックレスポンスを生成
	text := m.generateMockText(language)

	startTime := time.Now()
	result := &TranscriptionResult{
		Text:       text,
		Language:   language,
		Confidence: 0.95,
		Duration:   1500, // 1.5秒
		Words:      m.generateMockWords(text, language),
		Segments:   m.generateMockSegments(text),
		Provider:   "mock",
		ProcessingTimeMs: time.Since(startTime).Milliseconds(),
	}

	return result, nil
}

// generateMockText は言語に基づいたモックテキストを生成
func (m *MockSTTClientNew) generateMockText(language string) string {
	// 言語コードの正規化（en-US -> en）
	lang := language
	if idx := strings.Index(language, "-"); idx != -1 {
		lang = language[:idx]
	}

	mockTexts := map[string]string{
		// 主要言語
		"en": "Hello, world!",
		"ja": "こんにちは",
		"zh": "你好",
		"ru": "Здравствуйте",
		"es": "¡Hola, mundo!",
		"fr": "Bonjour le monde!",
		"de": "Hallo Welt!",
		"pt": "Olá mundo!",
		"it": "Ciao mondo!",
		"ko": "안녕하세요",
		// マイナー言語（HaiLanGoの価値提案）
		"fa": "سلام",           // ペルシャ語
		"he": "שלום",           // ヘブライ語
		"ar": "مرحبا",          // アラビア語
		"tr": "Merhaba dünya!", // トルコ語
		"ku": "Silav",          // クルド語
		"am": "ሰላም",            // アムハラ語
		"bo": "བཀྲ་ཤིས་བདེ་ལེགས།",  // チベット語
	}

	if text, ok := mockTexts[lang]; ok {
		return text
	}
	return "Hello, world!"
}

// generateMockWords は単語タイミング情報を生成
func (m *MockSTTClientNew) generateMockWords(text string, language string) []WordTiming {
	words := strings.Fields(text)
	result := make([]WordTiming, len(words))

	for i, word := range words {
		startTime := float64(i) * 0.5
		result[i] = WordTiming{
			Word:       word,
			Start:      startTime,
			End:        startTime + 0.4,
			Confidence: 0.90 + float64(i%10)*0.01,
		}
	}

	return result
}

// generateMockSegments はセグメント情報を生成
func (m *MockSTTClientNew) generateMockSegments(text string) []TranscriptSegment {
	return []TranscriptSegment{
		{
			ID:         0,
			Start:      0.0,
			End:        1.5,
			Text:       text,
			Confidence: 0.95,
		},
	}
}

// GetSupportedLanguages はサポートされている言語一覧を取得（モック）
func (m *MockSTTClientNew) GetSupportedLanguages(ctx context.Context) ([]string, error) {
	// Whisperは99言語をサポート（モックでは主要言語をシミュレート）
	return WhisperSupportedLanguages, nil
}

// GetName はプロバイダー名を返す
func (m *MockSTTClientNew) GetName() string {
	return "mock"
}

// SetMockAudio はカスタム音声データを設定する（テスト用）
func (m *MockSTTClientNew) SetMockAudio(audioData []byte) error {
	m.mockAudio = audioData
	return nil
}

// Ensure MockSTTClientNew implements STTClient interface
var _ STTClient = (*MockSTTClientNew)(nil)

// ============================================================
// 発音評価モック
// ============================================================

// MockPronunciationEvaluator はモック発音評価クライアント
type MockPronunciationEvaluator struct{}

// NewMockPronunciationEvaluator は新しいモック発音評価クライアントを作成
func NewMockPronunciationEvaluator() *MockPronunciationEvaluator {
	return &MockPronunciationEvaluator{}
}

// EvaluatePronunciation は発音を評価する（モック）
func (m *MockPronunciationEvaluator) EvaluatePronunciation(ctx context.Context, audio io.Reader, expectedText string, language string) (*PronunciationResult, error) {
	if audio == nil {
		return nil, errors.New("audio cannot be nil")
	}

	audioData, err := io.ReadAll(audio)
	if err != nil {
		return nil, err
	}

	if len(audioData) == 0 {
		return nil, errors.New("audio data cannot be empty")
	}

	if expectedText == "" {
		return nil, errors.New("expected text cannot be empty")
	}

	// モック評価結果を生成
	words := strings.Fields(expectedText)
	wordScores := make([]WordScore, len(words))

	for i, word := range words {
		score := 85.0 + float64(i%15)
		wordScores[i] = WordScore{
			Word:      word,
			Score:     score,
			ErrorType: "",
		}
		if score < 70 {
			wordScores[i].ErrorType = "mispronunciation"
			wordScores[i].Suggestion = "Try pronouncing more clearly"
		}
	}

	result := &PronunciationResult{
		OverallScore:   85.0,
		AccuracyScore:  88.0,
		FluencyScore:   82.0,
		ProsodyScore:   80.0,
		RecognizedText: expectedText, // モックでは完全一致
		ExpectedText:   expectedText,
		WordScores:     wordScores,
		Feedback:       "Great job! Your pronunciation is very good.",
		ImprovementTips: []string{
			"Focus on speaking more smoothly",
			"Pay attention to word stress",
		},
		DomainNotes:      "",
		EvaluationMethod: "whisper_llm",
	}

	return result, nil
}

// Ensure MockPronunciationEvaluator implements PronunciationEvaluator interface
var _ PronunciationEvaluator = (*MockPronunciationEvaluator)(nil)
