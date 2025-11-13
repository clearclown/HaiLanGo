package ocr

import (
	"context"
	"testing"

	"github.com/clearclown/HaiLanGo/backend/internal/models"
	"github.com/clearclown/HaiLanGo/backend/pkg/cache"
	"github.com/clearclown/HaiLanGo/backend/pkg/ocr"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestProcessPage(t *testing.T) {
	ctx := context.Background()
	ocrClient := ocr.NewMockOCRClient()
	mockCache := cache.NewMockCache()
	service := NewOCRService(ocrClient, mockCache)

	pageID := uuid.New()
	imageData := []byte("test image data")
	languages := []string{"ru", "en"}

	// 初回処理
	page, err := service.ProcessPage(ctx, pageID, imageData, languages)
	require.NoError(t, err)
	require.NotNil(t, page)

	assert.Equal(t, pageID, page.ID)
	assert.NotEmpty(t, page.OCRText)
	assert.Equal(t, "ru", page.DetectedLang)
	assert.Greater(t, page.OCRConfidence, 0.0)
	assert.LessOrEqual(t, page.OCRConfidence, 1.0)
	assert.Equal(t, models.OCRStatusCompleted, page.OCRStatus)
}

func TestProcessPage_WithCache(t *testing.T) {
	ctx := context.Background()
	ocrClient := ocr.NewMockOCRClient()
	mockCache := cache.NewMockCache()
	service := NewOCRService(ocrClient, mockCache)

	pageID := uuid.New()
	imageData := []byte("test image data")
	languages := []string{"ja"}

	// 1回目の処理
	page1, err := service.ProcessPage(ctx, pageID, imageData, languages)
	require.NoError(t, err)

	// 2回目の処理（キャッシュから取得）
	page2, err := service.ProcessPage(ctx, pageID, imageData, languages)
	require.NoError(t, err)

	// 同じ結果が返されることを確認
	assert.Equal(t, page1.OCRText, page2.OCRText)
	assert.Equal(t, page1.DetectedLang, page2.DetectedLang)
	assert.Equal(t, page1.OCRConfidence, page2.OCRConfidence)
}

func TestProcessPage_DifferentLanguages(t *testing.T) {
	ctx := context.Background()
	ocrClient := ocr.NewMockOCRClient()
	mockCache := cache.NewMockCache()
	service := NewOCRService(ocrClient, mockCache)

	pageID := uuid.New()
	imageData := []byte("test image data")

	tests := []struct {
		name      string
		languages []string
		wantLang  string
	}{
		{
			name:      "Russian",
			languages: []string{"ru"},
			wantLang:  "ru",
		},
		{
			name:      "Japanese",
			languages: []string{"ja"},
			wantLang:  "ja",
		},
		{
			name:      "Chinese",
			languages: []string{"zh"},
			wantLang:  "zh",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			page, err := service.ProcessPage(ctx, pageID, imageData, tt.languages)
			require.NoError(t, err)
			assert.Equal(t, tt.wantLang, page.DetectedLang)
		})
	}
}

func TestGenerateCacheKey(t *testing.T) {
	service := NewOCRService(ocr.NewMockOCRClient(), cache.NewMockCache())

	imageData1 := []byte("test image 1")
	imageData2 := []byte("test image 2")
	languages1 := []string{"ru", "en"}
	languages2 := []string{"ja", "zh"}

	// 同じデータと言語で同じキーが生成される
	key1a := service.generateCacheKey(imageData1, languages1)
	key1b := service.generateCacheKey(imageData1, languages1)
	assert.Equal(t, key1a, key1b)

	// 異なるデータで異なるキーが生成される
	key2 := service.generateCacheKey(imageData2, languages1)
	assert.NotEqual(t, key1a, key2)

	// 異なる言語で異なるキーが生成される
	key3 := service.generateCacheKey(imageData1, languages2)
	assert.NotEqual(t, key1a, key3)

	// キーのフォーマット確認
	assert.Contains(t, key1a, "ocr:")
}

func TestBuildPageFromOCRResult(t *testing.T) {
	service := NewOCRService(ocr.NewMockOCRClient(), cache.NewMockCache())

	pageID := uuid.New()
	result := &ocr.OCRResult{
		Text:             "Test OCR text",
		DetectedLanguage: "en",
		Confidence:       0.95,
	}

	page := service.buildPageFromOCRResult(pageID, result)

	assert.Equal(t, pageID, page.ID)
	assert.Equal(t, result.Text, page.OCRText)
	assert.Equal(t, result.DetectedLanguage, page.DetectedLang)
	assert.Equal(t, result.Confidence, page.OCRConfidence)
	assert.Equal(t, models.OCRStatusCompleted, page.OCRStatus)
	assert.False(t, page.UpdatedAt.IsZero())
}
