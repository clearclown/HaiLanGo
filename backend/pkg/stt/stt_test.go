package stt

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestGoogleSTTClient はGoogle Cloud STTクライアントのテスト
func TestGoogleSTTClient(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name        string
		audioData   []byte
		lang        string
		expectError bool
	}{
		{
			name:        "正常な認識",
			audioData:   []byte("test audio"),
			lang:        "en-US",
			expectError: false,
		},
		{
			name:        "空の音声データ",
			audioData:   []byte{},
			lang:        "en-US",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := NewGoogleSTTClient()
			result, err := client.Recognize(ctx, tt.audioData, tt.lang)

			if tt.expectError {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.NotNil(t, result)
			assert.NotEmpty(t, result.Text)
		})
	}
}

// TestWhisperClient はWhisper APIクライアントのテスト
func TestWhisperClient(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name        string
		audioData   []byte
		lang        string
		expectError bool
	}{
		{
			name:        "正常な認識",
			audioData:   []byte("test audio"),
			lang:        "en",
			expectError: false,
		},
		{
			name:        "空の音声データ",
			audioData:   []byte{},
			lang:        "en",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := NewWhisperClient()
			result, err := client.Recognize(ctx, tt.audioData, tt.lang)

			if tt.expectError {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.NotNil(t, result)
			assert.NotEmpty(t, result.Text)
		})
	}
}

// TestMockSTTClient はモックSTTクライアントのテスト
func TestMockSTTClient(t *testing.T) {
	ctx := context.Background()

	client := NewMockSTTClient()
	audioData := []byte("test audio")
	lang := "en"

	result, err := client.Recognize(ctx, audioData, lang)

	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.NotEmpty(t, result.Text)
	assert.Equal(t, lang, result.Language)
}

// TestSTTClientFactory はSTTクライアントファクトリーのテスト
func TestSTTClientFactory(t *testing.T) {
	tests := []struct {
		name         string
		useMock      bool
		apiKey       string
		expectType   string
	}{
		{
			name:       "モックを使用",
			useMock:    true,
			apiKey:     "",
			expectType: "mock",
		},
		{
			name:       "Google STTを使用",
			useMock:    false,
			apiKey:     "test-api-key",
			expectType: "google",
		},
		{
			name:       "APIキーなしでモックを使用",
			useMock:    false,
			apiKey:     "",
			expectType: "mock",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := NewSTTClient(tt.useMock, tt.apiKey)
			assert.NotNil(t, client)
		})
	}
}
