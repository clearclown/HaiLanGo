package handler

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/clearclown/HaiLanGo/backend/internal/models"
	"github.com/clearclown/HaiLanGo/backend/internal/repository"
	ocrservice "github.com/clearclown/HaiLanGo/backend/internal/service/ocr"
	"github.com/clearclown/HaiLanGo/backend/internal/websocket"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// OCRHandler はOCR APIのハンドラー
type OCRHandler struct {
	repo       repository.OCRRepositoryInterface
	ocrService *ocrservice.OCRService
	wsHub      *websocket.Hub
}

// NewOCRHandler はOCRハンドラーを作成
func NewOCRHandler(repo repository.OCRRepositoryInterface, ocrService *ocrservice.OCRService, wsHub *websocket.Hub) *OCRHandler {
	return &OCRHandler{
		repo:       repo,
		ocrService: ocrService,
		wsHub:      wsHub,
	}
}

// RegisterRoutes はOCR APIのルートを登録
func (h *OCRHandler) RegisterRoutes(rg *gin.RouterGroup) {
	ocr := rg.Group("/ocr")
	{
		// 特定ページのOCR処理
		ocr.POST("/books/:bookId/pages/:pageNumber", h.ProcessPage)

		// バッチOCR処理
		ocr.POST("/books/:bookId/batch", h.BatchProcess)

		// OCRジョブのステータス取得
		ocr.GET("/jobs/:jobId", h.GetJobStatus)

		// OCR処理結果取得
		ocr.GET("/jobs/:jobId/result", h.GetJobResult)

		// 書籍のOCRジョブ一覧
		ocr.GET("/books/:bookId/jobs", h.GetBookJobs)

		// OCR統計情報
		ocr.GET("/statistics", h.GetStatistics)
	}
}

// ProcessPage は特定ページのOCR処理を開始
// POST /api/v1/ocr/books/:bookId/pages/:pageNumber
func (h *OCRHandler) ProcessPage(c *gin.Context) {
	userIDStr, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	bookID, err := uuid.Parse(c.Param("bookId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid book ID"})
		return
	}

	pageNumber, err := strconv.Atoi(c.Param("pageNumber"))
	if err != nil || pageNumber < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid page number"})
		return
	}

	var req models.ProcessPageOCRRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// 画像URLを構築（実際の実装ではストレージから取得）
	imageURL := "/storage/books/" + bookID.String() + "/pages/" + strconv.Itoa(pageNumber) + ".jpg"

	// OCRジョブを作成
	job, err := h.repo.CreateJob(c.Request.Context(), userID, bookID, pageNumber, imageURL, req.Language)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create OCR job"})
		return
	}

	// バックグラウンドでOCR処理を開始
	go func() {
		ctx := c.Request.Context()
		jobIDStr := job.ID.String()

		// ジョブステータスを処理中に更新
		h.repo.UpdateJobStatus(ctx, jobIDStr, models.OCRStatusProcessing, 10)

		// 画像データを取得（実際の実装ではストレージから取得）
		imageData := h.getImageData(ctx, imageURL)
		if imageData == nil {
			h.repo.UpdateJobStatus(ctx, jobIDStr, models.OCRStatusFailed, 0)
			h.repo.UpdateJobError(ctx, jobIDStr, "Failed to load image data")
			return
		}

		h.repo.UpdateJobStatus(ctx, jobIDStr, models.OCRStatusProcessing, 30)

		// OCRサービスを使用してOCR処理を実行
		if h.ocrService != nil {
			languages := []string{req.Language}
			page, err := h.ocrService.ProcessPage(ctx, job.ID, imageData, languages)
			if err != nil {
				h.repo.UpdateJobStatus(ctx, jobIDStr, models.OCRStatusFailed, 0)
				h.repo.UpdateJobError(ctx, jobIDStr, err.Error())
				return
			}

			// 結果を設定
			result := &models.OCRResult{
				Text:             page.OCRText,
				Confidence:       page.OCRConfidence,
				DetectedLanguage: page.DetectedLang,
			}
			h.repo.UpdateJobResult(ctx, jobIDStr, result)
			h.repo.UpdateJobStatus(ctx, jobIDStr, models.OCRStatusCompleted, 100)

			// WebSocket経由で完了を通知
			if h.wsHub != nil {
				message, err := websocket.NewOCRProgressMessage(
					bookID,
					1,
					pageNumber,
					"completed",
					fmt.Sprintf("OCR completed for page %d", pageNumber),
				)
				if err == nil {
					h.wsHub.SendToUser(userID, message)
				}
			}
		} else {
			// OCRサービスがない場合はInMemoryシミュレーション
			if inMemRepo, ok := h.repo.(*repository.InMemoryOCRRepository); ok {
				inMemRepo.SimulateOCRProcessing(ctx, jobIDStr)
			}
		}
	}()

	response := &models.OCRJobResponse{
		JobID:      job.ID,
		BookID:     job.BookID,
		PageNumber: job.PageNumber,
		Status:     job.Status,
		Progress:   job.Progress,
		CreatedAt:  job.CreatedAt,
		UpdatedAt:  job.UpdatedAt,
	}

	c.JSON(http.StatusAccepted, response)
}

// BatchProcess はバッチOCR処理を開始
// POST /api/v1/ocr/books/:bookId/batch
func (h *OCRHandler) BatchProcess(c *gin.Context) {
	userIDStr, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	bookID, err := uuid.Parse(c.Param("bookId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid book ID"})
		return
	}

	var req models.BatchOCRRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// 実際の実装では書籍のページ数を取得
	totalPages := 150 // サンプルデータ

	jobIDs := make([]uuid.UUID, 0, totalPages)

	// 各ページのOCRジョブを作成
	for i := 1; i <= totalPages; i++ {
		imageURL := "/storage/books/" + bookID.String() + "/pages/" + strconv.Itoa(i) + ".jpg"

		job, err := h.repo.CreateJob(c.Request.Context(), userID, bookID, i, imageURL, req.Language)
		if err != nil {
			continue
		}

		jobIDs = append(jobIDs, job.ID)

		// バックグラウンドでOCR処理を開始
		go func(jobID uuid.UUID, pageNum int, imgURL string) {
			ctx := context.Background()
			jobIDStr := jobID.String()

			// OCRサービスを使用して処理
			if h.ocrService != nil {
				h.repo.UpdateJobStatus(ctx, jobIDStr, models.OCRStatusProcessing, 10)

				imageData := h.getImageData(ctx, imgURL)
				if imageData == nil {
					h.repo.UpdateJobStatus(ctx, jobIDStr, models.OCRStatusFailed, 0)
					h.repo.UpdateJobError(ctx, jobIDStr, "Failed to load image data")
					return
				}

				languages := []string{req.Language}
				page, err := h.ocrService.ProcessPage(ctx, jobID, imageData, languages)
				if err != nil {
					h.repo.UpdateJobStatus(ctx, jobIDStr, models.OCRStatusFailed, 0)
					h.repo.UpdateJobError(ctx, jobIDStr, err.Error())
					return
				}

				result := &models.OCRResult{
					Text:             page.OCRText,
					Confidence:       page.OCRConfidence,
					DetectedLanguage: page.DetectedLang,
				}
				h.repo.UpdateJobResult(ctx, jobIDStr, result)
				h.repo.UpdateJobStatus(ctx, jobIDStr, models.OCRStatusCompleted, 100)
			} else {
				if inMemRepo, ok := h.repo.(*repository.InMemoryOCRRepository); ok {
					inMemRepo.SimulateOCRProcessing(ctx, jobIDStr)
				}
			}

			// WebSocket経由で進捗を通知
			if h.wsHub != nil {
				message, err := websocket.NewOCRProgressMessage(
					bookID,
					totalPages,
					pageNum,
					"processing",
					fmt.Sprintf("Processing page %d of %d", pageNum, totalPages),
				)
				if err == nil {
					h.wsHub.SendToUser(userID, message)
				}
			}
		}(job.ID, i, imageURL)
	}

	response := &models.BatchOCRResponse{
		BookID:     bookID,
		TotalPages: totalPages,
		JobIDs:     jobIDs,
		CreatedAt:  time.Now(),
	}

	c.JSON(http.StatusAccepted, response)
}

// GetJobStatus はOCRジョブのステータスを取得
// GET /api/v1/ocr/jobs/:jobId
func (h *OCRHandler) GetJobStatus(c *gin.Context) {
	_, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	jobID := c.Param("jobId")

	job, err := h.repo.GetJob(c.Request.Context(), jobID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Job not found"})
		return
	}

	response := &models.OCRJobResponse{
		JobID:      job.ID,
		BookID:     job.BookID,
		PageNumber: job.PageNumber,
		Status:     job.Status,
		Progress:   job.Progress,
		CreatedAt:  job.CreatedAt,
		UpdatedAt:  job.UpdatedAt,
	}

	c.JSON(http.StatusOK, response)
}

// GetJobResult はOCR処理結果を取得
// GET /api/v1/ocr/jobs/:jobId/result
func (h *OCRHandler) GetJobResult(c *gin.Context) {
	_, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	jobID := c.Param("jobId")

	job, err := h.repo.GetJob(c.Request.Context(), jobID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Job not found"})
		return
	}

	if job.Status != models.OCRStatusCompleted && job.Status != models.OCRStatusFailed {
		c.JSON(http.StatusAccepted, gin.H{
			"message":  "Job is still processing",
			"status":   job.Status,
			"progress": job.Progress,
		})
		return
	}

	response := &models.OCRResultResponse{
		JobID:       job.ID,
		BookID:      job.BookID,
		PageNumber:  job.PageNumber,
		Status:      job.Status,
		Result:      job.Result,
		Error:       job.Error,
		CreatedAt:   job.CreatedAt,
		CompletedAt: job.CompletedAt,
	}

	c.JSON(http.StatusOK, response)
}

// GetBookJobs は書籍のOCRジョブ一覧を取得
// GET /api/v1/ocr/books/:bookId/jobs
func (h *OCRHandler) GetBookJobs(c *gin.Context) {
	_, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	bookID, err := uuid.Parse(c.Param("bookId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid book ID"})
		return
	}

	jobs, err := h.repo.GetJobsByBookID(c.Request.Context(), bookID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get OCR jobs"})
		return
	}

	responses := make([]*models.OCRJobResponse, 0, len(jobs))
	for _, job := range jobs {
		responses = append(responses, &models.OCRJobResponse{
			JobID:      job.ID,
			BookID:     job.BookID,
			PageNumber: job.PageNumber,
			Status:     job.Status,
			Progress:   job.Progress,
			CreatedAt:  job.CreatedAt,
			UpdatedAt:  job.UpdatedAt,
		})
	}

	c.JSON(http.StatusOK, responses)
}

// getImageData は画像データを取得するヘルパー関数
func (h *OCRHandler) getImageData(ctx context.Context, imageURL string) []byte {
	// 実際の実装ではストレージ（S3、R2など）から画像を取得
	// 現在はモックデータを返す
	return []byte("mock-image-data-for-ocr-processing")
}

// GetStatistics はOCR統計情報を取得
// GET /api/v1/ocr/statistics
func (h *OCRHandler) GetStatistics(c *gin.Context) {
	userIDStr, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	stats, err := h.repo.GetStatistics(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get OCR statistics"})
		return
	}

	c.JSON(http.StatusOK, stats)
}
