package stt

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/clearclown/HaiLanGo/backend/internal/models"
	"github.com/clearclown/HaiLanGo/backend/internal/service/audio"
	"github.com/clearclown/HaiLanGo/backend/pkg/stt"
)

// STTService は音声認識・発音評価サービス
type STTService struct {
	sttClient      stt.STTClient
	audioProcessor *audio.AudioProcessor
}

// NewSTTService は新しいSTTサービスを作成する
func NewSTTService() *STTService {
	useMock := os.Getenv("USE_MOCK_APIS") == "true" || os.Getenv("TEST_USE_MOCKS") == "true"
	apiKey := os.Getenv("GOOGLE_CLOUD_STT_API_KEY")

	return &STTService{
		sttClient:      stt.NewSTTClient(useMock, apiKey),
		audioProcessor: audio.NewAudioProcessor(),
	}
}

// RecognizeSpeech は音声データをテキストに変換する
func (s *STTService) RecognizeSpeech(ctx context.Context, audioData []byte, language string) (*models.STTResult, error) {
	if len(audioData) == 0 {
		return nil, fmt.Errorf(ErrEmptyAudioData)
	}

	// 言語コードの検証
	if !isValidLanguageCode(language) {
		return nil, fmt.Errorf("%s: %s", ErrInvalidLanguage, language)
	}

	// 音声データを前処理
	processedResult, err := s.audioProcessor.Process(audioData)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", ErrAudioProcessing, err)
	}

	// 低品質音声の警告
	if processedResult.IsLowQuality {
		// 警告をログに記録（実装は省略）
	}

	// STT APIを呼び出し
	result, err := s.sttClient.Recognize(ctx, processedResult.ProcessedAudio, language)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", ErrSTTRecognition, err)
	}

	return result, nil
}

// EvaluatePronunciation は発音を評価する
func (s *STTService) EvaluatePronunciation(ctx context.Context, expectedText string, audioData []byte, language string) (*models.PronunciationScore, error) {
	if expectedText == "" {
		return nil, fmt.Errorf(ErrEmptyExpectedText)
	}

	if len(audioData) == 0 {
		return nil, fmt.Errorf(ErrEmptyAudioData)
	}

	// 音声認識を実行
	sttResult, err := s.RecognizeSpeech(ctx, audioData, language)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", ErrSTTRecognition, err)
	}

	// 正確性スコアを計算
	accuracyScore := CalculateAccuracyScore(expectedText, sttResult.Text)

	// 流暢性スコアを計算
	fluencyScore := CalculateFluencyScore(sttResult.Words, sttResult.Duration)

	// 期待される単語を分割
	expectedWords := parseWords(expectedText)

	// 発音スコアを計算
	pronunciationScore := CalculatePronunciationScore(expectedWords, sttResult.Words)

	// 単語ごとのスコアを計算
	wordScores := calculateWordScores(expectedWords, sttResult.Words)

	// 総合スコアを計算（重み付け平均）
	totalScore := (accuracyScore*AccuracyWeight + fluencyScore*FluencyWeight + pronunciationScore*PronunciationWeight) / 100

	// PronunciationScoreを作成
	score := &models.PronunciationScore{
		TotalScore:     totalScore,
		AccuracyScore:  accuracyScore,
		FluencyScore:   fluencyScore,
		PronuncScore:   pronunciationScore,
		WordScores:     wordScores,
		ExpectedText:   expectedText,
		RecognizedText: sttResult.Text,
		EvaluationID:   generateEvaluationID(),
		CreatedAt:      time.Now(),
	}

	// フィードバックを生成
	score.Feedback = GenerateFeedback(score)

	return score, nil
}

// isValidLanguageCode は言語コードが有効かどうかを検証する
func isValidLanguageCode(code string) bool {
	validCodes := []string{
		"en", "en-US", "en-GB",
		"ja", "ja-JP",
		"zh", "zh-CN",
		"ru", "ru-RU",
		"es", "es-ES",
		"fr", "fr-FR",
		"de", "de-DE",
		"it", "it-IT",
		"pt", "pt-PT",
		"ar", "ar-SA",
		"he", "he-IL",
		"tr", "tr-TR",
	}

	for _, valid := range validCodes {
		if code == valid {
			return true
		}
	}

	return false
}

// parseWords はテキストを単語に分割する
func parseWords(text string) []models.WordInfo {
	words := strings.Fields(text)
	wordInfos := make([]models.WordInfo, 0, len(words))

	for i, word := range words {
		wordInfo := models.WordInfo{
			Word:       word,
			StartTime:  float64(i) * 0.5,      // 仮の開始時間
			EndTime:    float64(i+1) * 0.5,    // 仮の終了時間
			Confidence: 1.0,
		}
		wordInfos = append(wordInfos, wordInfo)
	}

	return wordInfos
}

// calculateWordScores は単語ごとのスコアを計算する
func calculateWordScores(expectedWords, recognizedWords []models.WordInfo) []models.WordScore {
	wordScores := make([]models.WordScore, 0)

	maxLen := len(expectedWords)
	if len(recognizedWords) > maxLen {
		maxLen = len(recognizedWords)
	}

	for i := 0; i < maxLen; i++ {
		var expectedWord, recognizedWord string

		if i < len(expectedWords) {
			expectedWord = expectedWords[i].Word
		}

		if i < len(recognizedWords) {
			recognizedWord = recognizedWords[i].Word
		}

		score := CalculateAccuracyScore(expectedWord, recognizedWord)
		isCorrect := score >= WordScoreCorrectMin

		wordScore := models.WordScore{
			Word:           expectedWord,
			Score:          score,
			IsCorrect:      isCorrect,
			ExpectedWord:   expectedWord,
			RecognizedWord: recognizedWord,
		}

		wordScores = append(wordScores, wordScore)
	}

	return wordScores
}

// generateEvaluationID は評価IDを生成する
func generateEvaluationID() string {
	return fmt.Sprintf("eval_%d", time.Now().UnixNano())
}

// ReduceNoise はノイズ除去を行う（テスト用のグローバル関数）
func ReduceNoise(audioData []byte) ([]byte, error) {
	if len(audioData) == 0 {
		return nil, fmt.Errorf(ErrEmptyAudioData)
	}

	processor := audio.NewAudioProcessor()
	return processor.ApplyNoiseReduction(audioData)
}

// NormalizeVolume は音量正規化を行う（テスト用のグローバル関数）
func NormalizeVolume(audioData []byte) ([]byte, error) {
	processor := audio.NewAudioProcessor()
	return processor.NormalizeVolume(audioData)
}

// ConvertSampleRate はサンプリングレート変換を行う（テスト用のグローバル関数）
func ConvertSampleRate(audioData []byte, targetRate int) ([]byte, error) {
	processor := audio.NewAudioProcessor()
	return processor.ConvertSampleRate(audioData, targetRate)
}

// ValidateAudioQuality は音声品質を検証する（テスト用のグローバル関数）
func ValidateAudioQuality(audioData []byte) (bool, error) {
	processor := audio.NewAudioProcessor()
	return processor.ValidateFormat(audioData)
}
