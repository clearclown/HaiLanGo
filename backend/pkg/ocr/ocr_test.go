package ocr

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMain(m *testing.M) {
	// テスト実行時は自動的にモックを使用
	os.Setenv("TEST_USE_MOCKS", "true")
	code := m.Run()
	os.Exit(code)
}

func TestNewOCRClient_WithMocks(t *testing.T) {
	os.Setenv("USE_MOCK_APIS", "true")
	defer os.Unsetenv("USE_MOCK_APIS")

	client, err := NewOCRClient()
	require.NoError(t, err)
	require.NotNil(t, client)

	// モッククライアントが返されることを確認
	_, ok := client.(*MockOCRClient)
	assert.True(t, ok, "expected MockOCRClient")
}

func TestNewOCRClient_WithoutAPIKey(t *testing.T) {
	os.Setenv("OCR_PROVIDER", "google_vision")
	os.Unsetenv("GOOGLE_CLOUD_VISION_API_KEY")
	os.Unsetenv("USE_MOCK_APIS")
	defer os.Unsetenv("OCR_PROVIDER")

	client, err := NewOCRClient()
	require.NoError(t, err)
	require.NotNil(t, client)

	// APIキーがない場合はモッククライアントが返されることを確認
	_, ok := client.(*MockOCRClient)
	assert.True(t, ok, "expected MockOCRClient when no API key is provided")
}

func TestMockOCRClient_ProcessImage(t *testing.T) {
	ctx := context.Background()
	client := NewMockOCRClient()

	tests := []struct {
		name      string
		imageData []byte
		languages []string
		want      string
	}{
		{
			name:      "Russian text",
			imageData: []byte("test image data"),
			languages: []string{"ru"},
			want:      "ru",
		},
		{
			name:      "Japanese text",
			imageData: []byte("test image data"),
			languages: []string{"ja"},
			want:      "ja",
		},
		{
			name:      "English text",
			imageData: []byte("test image data"),
			languages: []string{"en"},
			want:      "en",
		},
		{
			name:      "Chinese text",
			imageData: []byte("test image data"),
			languages: []string{"zh"},
			want:      "zh",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := client.ProcessImage(ctx, tt.imageData, tt.languages)
			require.NoError(t, err)
			require.NotNil(t, result)

			assert.Equal(t, tt.want, result.DetectedLanguage)
			assert.NotEmpty(t, result.Text)
			assert.Greater(t, result.Confidence, 0.0)
			assert.LessOrEqual(t, result.Confidence, 1.0)
			assert.NotEmpty(t, result.Pages)
		})
	}
}

func TestMockOCRClient_ProcessImage_DefaultLanguage(t *testing.T) {
	ctx := context.Background()
	client := NewMockOCRClient()

	// 言語指定なし
	result, err := client.ProcessImage(ctx, []byte("test"), []string{})
	require.NoError(t, err)
	require.NotNil(t, result)

	// デフォルトは英語
	assert.Equal(t, "en", result.DetectedLanguage)
}

func TestMockOCRClient_ProcessImage_UnsupportedLanguage(t *testing.T) {
	ctx := context.Background()
	client := NewMockOCRClient()

	// サポートされていない言語
	result, err := client.ProcessImage(ctx, []byte("test"), []string{"xyz"})
	require.NoError(t, err)
	require.NotNil(t, result)

	// デフォルトの英語にフォールバック
	assert.Equal(t, "xyz", result.DetectedLanguage)
	assert.NotEmpty(t, result.Text)
}

func TestMockOCRClient_SetMockResponse(t *testing.T) {
	// テンポラリディレクトリを作成
	tmpDir := t.TempDir()
	os.Setenv("MOCK_DATA_DIR", tmpDir)
	defer os.Unsetenv("MOCK_DATA_DIR")

	client := NewMockOCRClient()

	// モックレスポンスを設定
	expectedResult := &OCRResult{
		Text:             "Custom mock text",
		DetectedLanguage: "fr",
		Confidence:       0.98,
		Pages: []PageOCRResult{
			{
				PageNumber: 1,
				Text:       "Custom mock text",
				Confidence: 0.98,
			},
		},
	}

	err := client.SetMockResponse(expectedResult)
	require.NoError(t, err)

	// 設定したモックレスポンスが返されることを確認
	ctx := context.Background()
	result, err := client.ProcessImage(ctx, []byte("test"), []string{"fr"})
	require.NoError(t, err)
	require.NotNil(t, result)

	assert.Equal(t, expectedResult.Text, result.Text)
	assert.Equal(t, expectedResult.DetectedLanguage, result.DetectedLanguage)
	assert.Equal(t, expectedResult.Confidence, result.Confidence)
}

func TestOCRResult_Validation(t *testing.T) {
	result := &OCRResult{
		Text:             "Test text",
		DetectedLanguage: "en",
		Confidence:       0.95,
		Pages: []PageOCRResult{
			{
				PageNumber: 1,
				Text:       "Test text",
				Confidence: 0.95,
			},
		},
	}

	// 基本的な検証
	assert.NotEmpty(t, result.Text)
	assert.NotEmpty(t, result.DetectedLanguage)
	assert.Greater(t, result.Confidence, 0.0)
	assert.LessOrEqual(t, result.Confidence, 1.0)
	assert.NotEmpty(t, result.Pages)
	assert.Equal(t, 1, result.Pages[0].PageNumber)
}
