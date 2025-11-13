package ocr

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// MockOCRClient はモックOCRクライアント
type MockOCRClient struct {
	dataDir string
}

// NewMockOCRClient は新しいモックOCRクライアントを作成する
func NewMockOCRClient() *MockOCRClient {
	dataDir := os.Getenv("MOCK_DATA_DIR")
	if dataDir == "" {
		dataDir = "./mocks/data"
	}

	return &MockOCRClient{
		dataDir: dataDir,
	}
}

// ProcessImage は画像データをOCR処理する（モック）
func (m *MockOCRClient) ProcessImage(ctx context.Context, imageData []byte, languages []string) (*OCRResult, error) {
	// モックデータファイルから読み込み
	mockFile := filepath.Join(m.dataDir, "ocr", "sample_response.json")

	data, err := os.ReadFile(mockFile)
	if err != nil {
		// ファイルがない場合はデフォルトのモックレスポンスを返す
		return m.generateDefaultResponse(imageData, languages), nil
	}

	var result OCRResult
	if err := json.Unmarshal(data, &result); err != nil {
		return m.generateDefaultResponse(imageData, languages), nil
	}

	return &result, nil
}

// generateDefaultResponse はデフォルトのモックレスポンスを生成する
func (m *MockOCRClient) generateDefaultResponse(imageData []byte, languages []string) *OCRResult {
	detectedLang := "en"
	if len(languages) > 0 {
		detectedLang = languages[0]
	}

	// 言語ごとのサンプルテキスト
	sampleTexts := map[string]string{
		"ru": "Здравствуйте! Это пример текста из OCR.",
		"ja": "こんにちは！これはOCRからのサンプルテキストです。",
		"zh": "你好！这是来自OCR的示例文本。",
		"en": "Hello! This is sample text from OCR.",
		"es": "¡Hola! Este es un texto de muestra de OCR.",
		"fr": "Bonjour! Ceci est un exemple de texte OCR.",
		"de": "Hallo! Dies ist ein Beispieltext aus OCR.",
		"ar": "مرحبا! هذا نص نموذجي من OCR.",
		"he": "שלום! זהו טקסט לדוגמה מ-OCR.",
		"fa": "سلام! این یک متن نمونه از OCR است.",
	}

	text, ok := sampleTexts[detectedLang]
	if !ok {
		text = sampleTexts["en"]
	}

	return &OCRResult{
		Text:             text,
		DetectedLanguage: detectedLang,
		Confidence:       0.95,
		Pages: []PageOCRResult{
			{
				PageNumber: 1,
				Text:       text,
				Confidence: 0.95,
			},
		},
	}
}

// SetMockResponse はモックレスポンスを設定する（テスト用）
func (m *MockOCRClient) SetMockResponse(result *OCRResult) error {
	mockFile := filepath.Join(m.dataDir, "ocr", "sample_response.json")

	// ディレクトリが存在しない場合は作成
	if err := os.MkdirAll(filepath.Dir(mockFile), 0755); err != nil {
		return fmt.Errorf("failed to create mock data directory: %w", err)
	}

	data, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal mock response: %w", err)
	}

	if err := os.WriteFile(mockFile, data, 0644); err != nil {
		return fmt.Errorf("failed to write mock response: %w", err)
	}

	return nil
}
