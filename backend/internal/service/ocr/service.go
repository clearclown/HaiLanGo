package ocr

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"time"

	"github.com/clearclown/HaiLanGo/backend/internal/models"
	"github.com/clearclown/HaiLanGo/backend/pkg/cache"
	"github.com/clearclown/HaiLanGo/backend/pkg/ocr"
	"github.com/google/uuid"
)

// OCRService はOCR処理サービス
type OCRService struct {
	ocrClient ocr.OCRClient
	cache     cache.Cache
	cacheTTL  time.Duration
}

// NewOCRService は新しいOCRサービスを作成する
func NewOCRService(ocrClient ocr.OCRClient, cache cache.Cache) *OCRService {
	return &OCRService{
		ocrClient: ocrClient,
		cache:     cache,
		cacheTTL:  7 * 24 * time.Hour, // 7日間
	}
}

// ProcessPage はページのOCR処理を行う
func (s *OCRService) ProcessPage(ctx context.Context, pageID uuid.UUID, imageData []byte, languages []string) (*models.Page, error) {
	// キャッシュキーを生成
	cacheKey := s.generateCacheKey(imageData, languages)

	// キャッシュから取得を試みる
	if cachedData, err := s.cache.Get(ctx, cacheKey); err == nil {
		var result ocr.OCRResult
		if err := json.Unmarshal(cachedData, &result); err == nil {
			return s.buildPageFromOCRResult(pageID, &result), nil
		}
	}

	// OCR処理を実行
	result, err := s.ocrClient.ProcessImage(ctx, imageData, languages)
	if err != nil {
		return nil, fmt.Errorf("OCR processing failed: %w", err)
	}

	// キャッシュに保存
	if data, err := json.Marshal(result); err == nil {
		s.cache.Set(ctx, cacheKey, data, s.cacheTTL)
	}

	return s.buildPageFromOCRResult(pageID, result), nil
}

// generateCacheKey は画像データと言語からキャッシュキーを生成する
func (s *OCRService) generateCacheKey(imageData []byte, languages []string) string {
	hash := sha256.New()
	hash.Write(imageData)
	for _, lang := range languages {
		hash.Write([]byte(lang))
	}
	return "ocr:" + hex.EncodeToString(hash.Sum(nil))
}

// buildPageFromOCRResult はOCR結果からPageモデルを構築する
func (s *OCRService) buildPageFromOCRResult(pageID uuid.UUID, result *ocr.OCRResult) *models.Page {
	now := time.Now()
	return &models.Page{
		ID:            pageID,
		OCRText:       result.Text,
		OCRConfidence: result.Confidence,
		DetectedLang:  result.DetectedLanguage,
		OCRStatus:     models.OCRStatusCompleted,
		UpdatedAt:     now,
	}
}
