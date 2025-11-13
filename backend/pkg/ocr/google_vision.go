package ocr

import (
	"context"
	"fmt"
)

// GoogleVisionClient はGoogle Vision APIクライアント
type GoogleVisionClient struct {
	apiKey string
}

// NewGoogleVisionClient は新しいGoogle Vision APIクライアントを作成する
func NewGoogleVisionClient(apiKey string) *GoogleVisionClient {
	return &GoogleVisionClient{
		apiKey: apiKey,
	}
}

// ProcessImage は画像データをOCR処理する
func (g *GoogleVisionClient) ProcessImage(ctx context.Context, imageData []byte, languages []string) (*OCRResult, error) {
	// TODO: Google Vision APIの実装
	// 現在はスタブ実装
	return nil, fmt.Errorf("Google Vision API not yet implemented")
}
