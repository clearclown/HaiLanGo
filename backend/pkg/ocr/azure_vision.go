package ocr

import (
	"context"
	"fmt"
)

// AzureVisionClient はAzure Computer Vision APIクライアント
type AzureVisionClient struct {
	endpoint string
	apiKey   string
}

// NewAzureVisionClient は新しいAzure Computer Vision APIクライアントを作成する
func NewAzureVisionClient(endpoint, apiKey string) *AzureVisionClient {
	return &AzureVisionClient{
		endpoint: endpoint,
		apiKey:   apiKey,
	}
}

// ProcessImage は画像データをOCR処理する
func (a *AzureVisionClient) ProcessImage(ctx context.Context, imageData []byte, languages []string) (*OCRResult, error) {
	// TODO: Azure Computer Vision APIの実装
	// 現在はスタブ実装
	return nil, fmt.Errorf("Azure Computer Vision API not yet implemented")
}
