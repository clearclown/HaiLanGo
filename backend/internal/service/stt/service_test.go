package stt

import (
	"context"
	"testing"

	"github.com/clearclown/HaiLanGo/backend/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestRecognizeSpeech はSTT音声認識機能のテスト
func TestRecognizeSpeech(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name          string
		audioData     []byte
		lang          string
		expectedText  string
		expectError   bool
	}{
		{
			name:         "正常な英語音声",
			audioData:    []byte("test audio data"),
			lang:         "en",
			expectedText: "Hello, world!",
			expectError:  false,
		},
		{
			name:         "正常なロシア語音声",
			audioData:    []byte("test russian audio"),
			lang:         "ru",
			expectedText: "Здравствуйте",
			expectError:  false,
		},
		{
			name:        "空の音声データ",
			audioData:   []byte{},
			lang:        "en",
			expectError: true,
		},
		{
			name:        "無効な言語コード",
			audioData:   []byte("test audio"),
			lang:        "invalid",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := NewSTTService()
			result, err := service.RecognizeSpeech(ctx, tt.audioData, tt.lang)

			if tt.expectError {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.NotNil(t, result)
			assert.NotEmpty(t, result.Text)
			assert.Equal(t, tt.lang, result.Language)
			assert.GreaterOrEqual(t, result.Confidence, 0.0)
			assert.LessOrEqual(t, result.Confidence, 1.0)
		})
	}
}

// TestEvaluatePronunciation は発音評価機能のテスト
func TestEvaluatePronunciation(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name           string
		expectedText   string
		audioData      []byte
		lang           string
		minScore       int
		expectError    bool
	}{
		{
			name:         "完璧な発音",
			expectedText: "Hello",
			audioData:    []byte("perfect pronunciation"),
			lang:         "en",
			minScore:     90,
			expectError:  false,
		},
		{
			name:         "良好な発音",
			expectedText: "Здравствуйте",
			audioData:    []byte("good pronunciation"),
			lang:         "ru",
			minScore:     70,
			expectError:  false,
		},
		{
			name:         "改善が必要な発音",
			expectedText: "Hello",
			audioData:    []byte("poor pronunciation"),
			lang:         "en",
			minScore:     30,
			expectError:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := NewSTTService()
			score, err := service.EvaluatePronunciation(ctx, tt.expectedText, tt.audioData, tt.lang)

			if tt.expectError {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.NotNil(t, score)
			assert.GreaterOrEqual(t, score.TotalScore, 0)
			assert.LessOrEqual(t, score.TotalScore, 100)
			assert.NotNil(t, score.Feedback)
			assert.NotEmpty(t, score.Feedback.Message)
		})
	}
}

// TestCalculateAccuracyScore は正確性スコアの計算テスト
func TestCalculateAccuracyScore(t *testing.T) {
	tests := []struct {
		name           string
		expected       string
		recognized     string
		expectedScore  int
	}{
		{
			name:          "完全一致",
			expected:      "Hello",
			recognized:    "Hello",
			expectedScore: 100,
		},
		{
			name:          "部分的な一致",
			expected:      "Hello",
			recognized:    "Hallo",
			expectedScore: 80,
		},
		{
			name:          "大文字小文字の違い",
			expected:      "Hello",
			recognized:    "hello",
			expectedScore: 100,
		},
		{
			name:          "完全に異なる",
			expected:      "Hello",
			recognized:    "Goodbye",
			expectedScore: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			score := CalculateAccuracyScore(tt.expected, tt.recognized)
			assert.Equal(t, tt.expectedScore, score)
		})
	}
}

// TestCalculateFluencyScore は流暢性スコアの計算テスト
func TestCalculateFluencyScore(t *testing.T) {
	tests := []struct {
		name          string
		words         []models.WordInfo
		duration      float64
		expectedScore int
	}{
		{
			name: "自然な速度",
			words: []models.WordInfo{
				{Word: "Hello", StartTime: 0.0, EndTime: 0.5},
				{Word: "world", StartTime: 0.6, EndTime: 1.1},
			},
			duration:      1.5,
			expectedScore: 90,
		},
		{
			name: "遅すぎる速度",
			words: []models.WordInfo{
				{Word: "Hello", StartTime: 0.0, EndTime: 2.0},
				{Word: "world", StartTime: 3.0, EndTime: 5.0},
			},
			duration:      6.0,
			expectedScore: 50,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			score := CalculateFluencyScore(tt.words, tt.duration)
			assert.GreaterOrEqual(t, score, 0)
			assert.LessOrEqual(t, score, 100)
		})
	}
}

// TestGenerateFeedback はフィードバック生成のテスト
func TestGenerateFeedback(t *testing.T) {
	tests := []struct {
		name          string
		score         *models.PronunciationScore
		expectLevel   string
	}{
		{
			name: "優秀なスコア",
			score: &models.PronunciationScore{
				TotalScore:     95,
				AccuracyScore:  98,
				FluencyScore:   92,
				PronuncScore:   95,
				ExpectedText:   "Hello",
				RecognizedText: "Hello",
			},
			expectLevel: "excellent",
		},
		{
			name: "良好なスコア",
			score: &models.PronunciationScore{
				TotalScore:     75,
				AccuracyScore:  78,
				FluencyScore:   72,
				PronuncScore:   75,
				ExpectedText:   "Hello",
				RecognizedText: "Hallo",
			},
			expectLevel: "good",
		},
		{
			name: "改善が必要なスコア",
			score: &models.PronunciationScore{
				TotalScore:     45,
				AccuracyScore:  50,
				FluencyScore:   40,
				PronuncScore:   45,
				ExpectedText:   "Hello",
				RecognizedText: "Helo",
			},
			expectLevel: "fair",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			feedback := GenerateFeedback(tt.score)
			assert.NotNil(t, feedback)
			assert.Equal(t, tt.expectLevel, feedback.Level)
			assert.NotEmpty(t, feedback.Message)

			if tt.score.TotalScore >= 70 {
				assert.NotEmpty(t, feedback.PositivePoints)
			}

			if tt.score.TotalScore < 90 {
				assert.NotEmpty(t, feedback.Improvements)
			}
		})
	}
}

// TestNoiseReduction はノイズ除去のテスト
func TestNoiseReduction(t *testing.T) {
	tests := []struct {
		name        string
		audioData   []byte
		expectError bool
	}{
		{
			name:        "正常な音声データ",
			audioData:   []byte("test audio data with noise"),
			expectError: false,
		},
		{
			name:        "空の音声データ",
			audioData:   []byte{},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cleaned, err := ReduceNoise(tt.audioData)

			if tt.expectError {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.NotNil(t, cleaned)
		})
	}
}

// TestNormalizeVolume は音量正規化のテスト
func TestNormalizeVolume(t *testing.T) {
	audioData := []byte("test audio data")
	normalized, err := NormalizeVolume(audioData)

	require.NoError(t, err)
	assert.NotNil(t, normalized)
}

// TestConvertSampleRate はサンプリングレート変換のテスト
func TestConvertSampleRate(t *testing.T) {
	audioData := []byte("test audio data")
	targetRate := 16000

	converted, err := ConvertSampleRate(audioData, targetRate)

	require.NoError(t, err)
	assert.NotNil(t, converted)
}

// TestValidateAudioQuality は音声品質検証のテスト
func TestValidateAudioQuality(t *testing.T) {
	tests := []struct {
		name          string
		audioData     []byte
		expectValid   bool
	}{
		{
			name:        "高品質音声",
			audioData:   []byte("high quality audio"),
			expectValid: true,
		},
		{
			name:        "低品質音声",
			audioData:   []byte("low quality"),
			expectValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			isValid, err := ValidateAudioQuality(tt.audioData)
			require.NoError(t, err)
			assert.Equal(t, tt.expectValid, isValid)
		})
	}
}
