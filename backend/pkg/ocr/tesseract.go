package ocr

import (
	"context"
	"fmt"
)

// TesseractClient はTesseract OCRクライアント
type TesseractClient struct{}

// NewTesseractClient は新しいTesseract OCRクライアントを作成する
func NewTesseractClient() *TesseractClient {
	return &TesseractClient{}
}

// ProcessImage は画像データをOCR処理する
func (t *TesseractClient) ProcessImage(ctx context.Context, imageData []byte, languages []string) (*OCRResult, error) {
	// TODO: Tesseract OCRの実装
	// 現在はスタブ実装
	return nil, fmt.Errorf("Tesseract OCR not yet implemented")
}
