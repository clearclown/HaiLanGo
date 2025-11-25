package ocr

import (
	"fmt"
	"os"
)

// NewOCRClient は環境変数に基づいて適切なOCRクライアントを返す
func NewOCRClient() (OCRClient, error) {
	// モック使用の判定
	useMocks := os.Getenv("USE_MOCK_APIS") == "true" ||
		os.Getenv("TEST_USE_MOCKS") == "true"

	if useMocks {
		return NewMockOCRClient(), nil
	}

	// プロバイダーの選択
	provider := OCRProvider(os.Getenv("OCR_PROVIDER"))
	if provider == "" {
		provider = ProviderGoogleVision // デフォルト
	}

	switch provider {
	case ProviderGoogleVision:
		// サービスアカウント認証を優先
		credentialsFile := os.Getenv("GOOGLE_APPLICATION_CREDENTIALS")
		if credentialsFile != "" {
			// ファイルが存在するか確認
			if _, err := os.Stat(credentialsFile); err == nil {
				return NewGoogleVisionClientWithCredentials(credentialsFile), nil
			}
		}

		// APIキー認証にフォールバック
		apiKey := os.Getenv("GOOGLE_CLOUD_VISION_API_KEY")
		if apiKey == "" {
			// 認証情報がない場合は自動的にモックを使用
			return NewMockOCRClient(), nil
		}
		return NewGoogleVisionClient(apiKey), nil

	case ProviderAzureVision:
		endpoint := os.Getenv("AZURE_COMPUTER_VISION_ENDPOINT")
		apiKey := os.Getenv("AZURE_COMPUTER_VISION_API_KEY")
		if endpoint == "" || apiKey == "" {
			// APIキーがない場合は自動的にモックを使用
			return NewMockOCRClient(), nil
		}
		return NewAzureVisionClient(endpoint, apiKey), nil

	case ProviderTesseract:
		return NewTesseractClient(), nil

	default:
		return nil, fmt.Errorf("unsupported OCR provider: %s", provider)
	}
}
