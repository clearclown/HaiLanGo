package stt

import (
	"context"

	"github.com/clearclown/HaiLanGo/backend/internal/models"
)

// STTClient は音声認識クライアントのインターフェース
type STTClient interface {
	// Recognize は音声データをテキストに変換する
	Recognize(ctx context.Context, audioData []byte, language string) (*models.STTResult, error)
}

// NewSTTClient は環境変数とAPIキーに基づいて適切なSTTクライアントを返す
func NewSTTClient(useMock bool, apiKey string) STTClient {
	if useMock || apiKey == "" {
		return NewMockSTTClient()
	}
	return NewGoogleSTTClient()
}
