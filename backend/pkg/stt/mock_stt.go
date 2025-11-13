package stt

import (
	"context"
	"fmt"
	"time"

	"github.com/clearclown/HaiLanGo/backend/internal/models"
)

// MockSTTClient はモックSTTクライアント
type MockSTTClient struct{}

// NewMockSTTClient は新しいモックSTTクライアントを作成する
func NewMockSTTClient() *MockSTTClient {
	return &MockSTTClient{}
}

// Recognize はモック音声認識を実行する
func (m *MockSTTClient) Recognize(ctx context.Context, audioData []byte, language string) (*models.STTResult, error) {
	if len(audioData) == 0 {
		return nil, fmt.Errorf("音声データが空です")
	}

	// モックレスポンスを生成
	var text string
	switch language {
	case "en", "en-US":
		text = "Hello, world!"
	case "ru", "ru-RU":
		text = "Здравствуйте"
	case "ja", "ja-JP":
		text = "こんにちは"
	case "zh", "zh-CN":
		text = "你好"
	default:
		text = "Hello, world!"
	}

	result := &models.STTResult{
		Text:       text,
		Language:   language,
		Confidence: 0.95,
		Duration:   1.5,
		Words: []models.WordInfo{
			{
				Word:       "Hello",
				StartTime:  0.0,
				EndTime:    0.5,
				Confidence: 0.96,
			},
			{
				Word:       "world",
				StartTime:  0.6,
				EndTime:    1.1,
				Confidence: 0.94,
			},
		},
		CreatedAt: time.Now(),
	}

	return result, nil
}
