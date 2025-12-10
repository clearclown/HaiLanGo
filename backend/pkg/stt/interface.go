package stt

import (
	"context"

	"github.com/clearclown/HaiLanGo/backend/internal/models"
)

// LegacySTTClient は旧STT APIクライアントのインターフェース（後方互換性のため保持）
// 新規実装はSTTClient (stt.go) を使用してください
type LegacySTTClient interface {
	// Recognize は音声データをテキストに変換する
	Recognize(ctx context.Context, audioData []byte, language string) (*models.STTResult, error)
}

// NewLegacySTTClient は環境変数とAPIキーに基づいて適切なレガシーSTTクライアントを返す
// 新規実装はNewSTTClient() (factory.go) を使用してください
func NewLegacySTTClient(useMock bool, apiKey string) LegacySTTClient {
	if useMock || apiKey == "" {
		return NewMockSTTClient()
	}
	return NewGoogleSTTClient()
}
