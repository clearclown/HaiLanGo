package audio

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestProcessAudio は音声処理のテスト
func TestProcessAudio(t *testing.T) {
	tests := []struct {
		name        string
		audioData   []byte
		expectError bool
	}{
		{
			name:        "正常な音声データ",
			audioData:   []byte("test audio data"),
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
			processor := NewAudioProcessor()
			result, err := processor.Process(tt.audioData)

			if tt.expectError {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.NotNil(t, result)
			assert.NotEmpty(t, result.ProcessedAudio)
		})
	}
}

// TestDetectNoiseLevel はノイズレベル検出のテスト
func TestDetectNoiseLevel(t *testing.T) {
	audioData := []byte("test audio with varying noise levels")
	processor := NewAudioProcessor()

	noiseLevel, err := processor.DetectNoiseLevel(audioData)

	require.NoError(t, err)
	assert.GreaterOrEqual(t, noiseLevel, 0.0)
	assert.LessOrEqual(t, noiseLevel, 1.0)
}

// TestApplyNoiseReduction はノイズ除去適用のテスト
func TestApplyNoiseReduction(t *testing.T) {
	audioData := []byte("noisy audio data")
	processor := NewAudioProcessor()

	cleaned, err := processor.ApplyNoiseReduction(audioData)

	require.NoError(t, err)
	assert.NotNil(t, cleaned)
	assert.NotEmpty(t, cleaned)
}

// TestNormalizeAudioVolume は音量正規化のテスト
func TestNormalizeAudioVolume(t *testing.T) {
	audioData := []byte("audio with varying volume")
	processor := NewAudioProcessor()

	normalized, err := processor.NormalizeVolume(audioData)

	require.NoError(t, err)
	assert.NotNil(t, normalized)
	assert.NotEmpty(t, normalized)
}

// TestConvertToTargetSampleRate はサンプリングレート変換のテスト
func TestConvertToTargetSampleRate(t *testing.T) {
	audioData := []byte("audio data")
	targetRate := 16000
	processor := NewAudioProcessor()

	converted, err := processor.ConvertSampleRate(audioData, targetRate)

	require.NoError(t, err)
	assert.NotNil(t, converted)
}

// TestValidateAudioFormat は音声フォーマット検証のテスト
func TestValidateAudioFormat(t *testing.T) {
	tests := []struct {
		name        string
		audioData   []byte
		expectValid bool
	}{
		{
			name:        "有効な音声フォーマット",
			audioData:   []byte("valid audio format"),
			expectValid: true,
		},
		{
			name:        "無効な音声フォーマット",
			audioData:   []byte("invalid"),
			expectValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			processor := NewAudioProcessor()
			isValid, err := processor.ValidateFormat(tt.audioData)

			require.NoError(t, err)
			assert.Equal(t, tt.expectValid, isValid)
		})
	}
}
