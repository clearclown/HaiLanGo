package stt

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/clearclown/HaiLanGo/backend/internal/models"
)

// WhisperClient はWhisper APIクライアント
type WhisperClient struct {
	apiKey string
}

// NewWhisperClient は新しいWhisper APIクライアントを作成する
func NewWhisperClient() *WhisperClient {
	apiKey := os.Getenv("OPENAI_API_KEY")
	return &WhisperClient{
		apiKey: apiKey,
	}
}

// Recognize はWhisper APIを使用して音声認識を実行する
func (w *WhisperClient) Recognize(ctx context.Context, audioData []byte, language string) (*models.STTResult, error) {
	if len(audioData) == 0 {
		return nil, fmt.Errorf("音声データが空です")
	}

	// 実際のAPIキーがない場合はモックを返す
	if w.apiKey == "" || os.Getenv("USE_MOCK_APIS") == "true" {
		mockClient := NewMockSTTClient()
		return mockClient.Recognize(ctx, audioData, language)
	}

	// TODO: 実際のWhisper API呼び出しを実装
	// 現時点ではモックレスポンスを返す
	result := &models.STTResult{
		Text:       "Hello, world!",
		Language:   language,
		Confidence: 0.90,
		Duration:   1.5,
		Words: []models.WordInfo{
			{
				Word:       "Hello",
				StartTime:  0.0,
				EndTime:    0.5,
				Confidence: 0.91,
			},
			{
				Word:       "world",
				StartTime:  0.6,
				EndTime:    1.1,
				Confidence: 0.89,
			},
		},
		CreatedAt: time.Now(),
	}

	return result, nil
}
