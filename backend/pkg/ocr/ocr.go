package ocr

import (
	"context"
)

// OCRClient はOCR APIのインターフェース
type OCRClient interface {
	// ProcessImage は画像データをOCR処理する
	ProcessImage(ctx context.Context, imageData []byte, languages []string) (*OCRResult, error)
}

// OCRResult はOCR処理の結果を表す
type OCRResult struct {
	Text             string   `json:"text"`               // 抽出されたテキスト
	DetectedLanguage string   `json:"detected_language"`  // 検出された言語
	Confidence       float64  `json:"confidence"`         // 信頼度（0.0-1.0）
	Pages            []PageOCRResult `json:"pages"`       // ページごとの結果
}

// PageOCRResult はページごとのOCR結果を表す
type PageOCRResult struct {
	PageNumber int     `json:"page_number"`
	Text       string  `json:"text"`
	Confidence float64 `json:"confidence"`
}

// OCRProvider はOCRプロバイダーの種類
type OCRProvider string

const (
	ProviderGoogleVision OCRProvider = "google_vision"
	ProviderAzureVision  OCRProvider = "azure_vision"
	ProviderTesseract    OCRProvider = "tesseract"
)
