package stt

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/clearclown/HaiLanGo/backend/internal/models"
)

// GoogleSTTClient はGoogle Cloud STTクライアント
type GoogleSTTClient struct {
	apiKey string
}

// NewGoogleSTTClient は新しいGoogle STTクライアントを作成する
func NewGoogleSTTClient() *GoogleSTTClient {
	apiKey := os.Getenv("GOOGLE_CLOUD_STT_API_KEY")
	return &GoogleSTTClient{
		apiKey: apiKey,
	}
}

// Recognize はGoogle Cloud STTを使用して音声認識を実行する
func (g *GoogleSTTClient) Recognize(ctx context.Context, audioData []byte, language string) (*models.STTResult, error) {
	if len(audioData) == 0 {
		return nil, fmt.Errorf("音声データが空です")
	}

	// 実際のAPIキーがない場合はモックを返す
	if g.apiKey == "" || os.Getenv("USE_MOCK_APIS") == "true" {
		mockClient := NewMockSTTClient()
		return mockClient.Recognize(ctx, audioData, language)
	}

	// TODO: 実際のGoogle Cloud STT API呼び出しを実装
	// 現時点ではモックレスポンスを返す
	result := &models.STTResult{
		Text:       "Hello, world!",
		Language:   language,
		Confidence: 0.92,
		Duration:   1.5,
		Words: []models.WordInfo{
			{
				Word:       "Hello",
				StartTime:  0.0,
				EndTime:    0.5,
				Confidence: 0.93,
			},
			{
				Word:       "world",
				StartTime:  0.6,
				EndTime:    1.1,
				Confidence: 0.91,
			},
		},
		CreatedAt: time.Now(),
	}

	return result, nil
}
