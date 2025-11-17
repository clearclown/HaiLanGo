package ocr

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/clearclown/HaiLanGo/backend/internal/models"
	"github.com/clearclown/HaiLanGo/backend/internal/repository"
	"github.com/clearclown/HaiLanGo/backend/internal/websocket"
	"github.com/clearclown/HaiLanGo/backend/pkg/cache"
	"github.com/clearclown/HaiLanGo/backend/pkg/ocr"
	"github.com/google/uuid"
)

// OCRService はOCR処理サービス
type OCRService struct {
	ocrClient ocr.OCRClient
	cache     cache.Cache
	cacheTTL  time.Duration
	pageRepo  repository.PageRepository
	wsHub     *websocket.Hub
}

// NewOCRService は新しいOCRサービスを作成する
func NewOCRService(ocrClient ocr.OCRClient, cache cache.Cache) *OCRService {
	return &OCRService{
		ocrClient: ocrClient,
		cache:     cache,
		cacheTTL:  7 * 24 * time.Hour, // 7日間
		pageRepo:  repository.NewMockPageRepository(),
	}
}

// SetWebSocketHub はWebSocketハブを設定する
func (s *OCRService) SetWebSocketHub(hub *websocket.Hub) {
	s.wsHub = hub
}

// SetPageRepository はページリポジトリを設定する
func (s *OCRService) SetPageRepository(repo repository.PageRepository) {
	s.pageRepo = repo
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

// PageData はページデータ
type PageData struct {
	PageID    uuid.UUID
	ImageData []byte
}

// ProcessBookResult は書籍処理結果
type ProcessBookResult struct {
	TotalPages      int
	ProcessedPages  int
	FailedPages     int
	ProcessingTime  time.Duration
	Errors          []error
}

// ProcessBook は書籍の全ページをOCR処理する（並列処理）
func (s *OCRService) ProcessBook(ctx context.Context, bookID uuid.UUID, userID string, pages []PageData, languages []string, maxConcurrency int) (*ProcessBookResult, error) {
	startTime := time.Now()

	if maxConcurrency <= 0 {
		maxConcurrency = 5 // デフォルトの並列数
	}

	result := &ProcessBookResult{
		TotalPages: len(pages),
		Errors:     []error{},
	}

	// 処理済みページ数をアトミックにカウント
	var processedCount int32

	// ワーカープールを作成
	jobs := make(chan PageData, len(pages))
	results := make(chan error, len(pages))

	var wg sync.WaitGroup

	// ワーカーを起動
	for i := 0; i < maxConcurrency; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for pageData := range jobs {
				// コンテキストがキャンセルされていないかチェック
				select {
				case <-ctx.Done():
					results <- ctx.Err()
					return
				default:
				}

				// ページをOCR処理
				page, err := s.ProcessPage(ctx, pageData.PageID, pageData.ImageData, languages)
				if err != nil {
					results <- fmt.Errorf("page %s failed: %w", pageData.PageID, err)
					continue
				}

				// データベースに保存
				if s.pageRepo != nil {
					page.BookID = bookID
					if err := s.pageRepo.Create(ctx, page); err != nil {
						results <- fmt.Errorf("failed to save page %s: %w", pageData.PageID, err)
						continue
					}
				}

				results <- nil // 成功

				// 処理済みページ数を更新してWebSocket通知を送信
				processed := atomic.AddInt32(&processedCount, 1)
				s.sendOCRProgress(userID, bookID.String(), len(pages), int(processed))
			}
		}()
	}

	// ジョブをワーカーに送信
	go func() {
		for _, page := range pages {
			jobs <- page
		}
		close(jobs)
	}()

	// ワーカーの完了を待つ
	go func() {
		wg.Wait()
		close(results)
	}()

	// 結果を集計
	for err := range results {
		if err != nil {
			result.FailedPages++
			result.Errors = append(result.Errors, err)
		} else {
			result.ProcessedPages++
		}
	}

	result.ProcessingTime = time.Since(startTime)

	// 完了通知を送信
	s.sendBookReady(userID, bookID.String(), result.ProcessedPages, result.TotalPages)

	return result, nil
}

// sendOCRProgress はOCR処理の進捗をWebSocket経由で送信する
func (s *OCRService) sendOCRProgress(userID, bookID string, totalPages, processedPages int) {
	if s.wsHub == nil {
		return
	}

	// UUIDに変換
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return
	}
	bookUUID, err := uuid.Parse(bookID)
	if err != nil {
		return
	}

	message, err := websocket.NewOCRProgressMessage(
		bookUUID,
		totalPages,
		processedPages,
		"processing",
		fmt.Sprintf("Processing page %d of %d", processedPages, totalPages),
	)
	if err != nil {
		return
	}

	s.wsHub.SendToUser(userUUID, message)
}

// sendBookReady は書籍処理完了をWebSocket経由で送信する
func (s *OCRService) sendBookReady(userID, bookID string, processedPages, totalPages int) {
	if s.wsHub == nil {
		return
	}

	// UUIDに変換
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return
	}
	bookUUID, err := uuid.Parse(bookID)
	if err != nil {
		return
	}

	message, err := websocket.NewBookReadyMessage(
		bookUUID,
		"", // title - 実際の実装ではbookRepoから取得する
		totalPages,
	)
	if err != nil {
		return
	}

	s.wsHub.SendToUser(userUUID, message)
}
