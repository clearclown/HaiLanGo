package handler

import (
	"net/http"

	"github.com/clearclown/HaiLanGo/backend/internal/models"
	"github.com/clearclown/HaiLanGo/backend/internal/repository"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// STTHandler はSTT APIのハンドラー
type STTHandler struct {
	repo repository.STTRepositoryInterface
}

// NewSTTHandler はSTTハンドラーを作成
func NewSTTHandler(repo repository.STTRepositoryInterface) *STTHandler {
	return &STTHandler{
		repo: repo,
	}
}

// RegisterRoutes はSTT APIのルートを登録
func (h *STTHandler) RegisterRoutes(rg *gin.RouterGroup) {
	stt := rg.Group("/stt")
	{
		// 音声認識・発音評価
		stt.POST("/recognize", h.Recognize)

		// サポート言語一覧
		stt.GET("/languages", h.GetLanguages)

		// STTジョブのステータス取得
		stt.GET("/jobs/:jobId", h.GetJobStatus)

		// 書籍のSTTジョブ一覧
		stt.GET("/books/:bookId/jobs", h.GetBookJobs)

		// ユーザーのSTT統計
		stt.GET("/statistics", h.GetStatistics)
	}
}

// Recognize は音声認識・発音評価を実行
// POST /api/v1/stt/recognize
func (h *STTHandler) Recognize(c *gin.Context) {
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

	var req models.STTRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// デフォルト値の設定
	if req.Options.Format == "" {
		req.Options.Format = "wav"
	}
	if !req.Options.EnablePunctuation {
		req.Options.EnablePunctuation = true
	}
	if !req.Options.EnableWordTiming {
		req.Options.EnableWordTiming = true
	}
	if !req.Options.Evaluate && req.ReferenceText != "" {
		req.Options.Evaluate = true
	}

	// ダミーのbookIDとpageNumberを使用（単一音声認識の場合）
	bookID := uuid.New()
	pageNumber := 0

	// AudioURLの生成（実際にはストレージにアップロードされたURLを使用）
	audioURL := "/storage/audio/temp_" + uuid.New().String() + "." + req.Options.Format

	// STTジョブを作成
	job, err := h.repo.CreateJob(c.Request.Context(), userID, bookID, pageNumber, audioURL, req.Language, req.ReferenceText, req.Options)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create STT job"})
		return
	}

	// バックグラウンドでSTT処理を開始
	go func() {
		if inMemRepo, ok := h.repo.(*repository.InMemorySTTRepository); ok {
			inMemRepo.SimulateSTTProcessing(c.Request.Context(), job.ID)
		}
	}()

	response := &models.STTJobResponse{
		JobID:     job.ID,
		Status:    job.Status,
		Progress:  job.Progress,
		CreatedAt: job.CreatedAt,
		UpdatedAt: job.UpdatedAt,
	}

	c.JSON(http.StatusAccepted, response)
}

// GetLanguages はサポート言語一覧を取得
// GET /api/v1/stt/languages
func (h *STTHandler) GetLanguages(c *gin.Context) {
	_, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	languages, err := h.repo.GetSupportedLanguages(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get languages"})
		return
	}

	c.JSON(http.StatusOK, languages)
}

// GetJobStatus はSTTジョブのステータスを取得
// GET /api/v1/stt/jobs/:jobId
func (h *STTHandler) GetJobStatus(c *gin.Context) {
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

	response := &models.STTJobResponse{
		JobID:      job.ID,
		BookID:     job.BookID,
		PageNumber: job.PageNumber,
		Status:     job.Status,
		Progress:   job.Progress,
		Result:     job.Result,
		Score:      job.Score,
		CreatedAt:  job.CreatedAt,
		UpdatedAt:  job.UpdatedAt,
	}

	c.JSON(http.StatusOK, response)
}

// GetBookJobs は書籍のSTTジョブ一覧を取得
// GET /api/v1/stt/books/:bookId/jobs
func (h *STTHandler) GetBookJobs(c *gin.Context) {
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get STT jobs"})
		return
	}

	responses := make([]*models.STTJobResponse, 0, len(jobs))
	for _, job := range jobs {
		responses = append(responses, &models.STTJobResponse{
			JobID:      job.ID,
			BookID:     job.BookID,
			PageNumber: job.PageNumber,
			Status:     job.Status,
			Progress:   job.Progress,
			Result:     job.Result,
			Score:      job.Score,
			CreatedAt:  job.CreatedAt,
			UpdatedAt:  job.UpdatedAt,
		})
	}

	c.JSON(http.StatusOK, responses)
}

// GetStatistics はユーザーのSTT統計情報を取得
// GET /api/v1/stt/statistics
func (h *STTHandler) GetStatistics(c *gin.Context) {
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get statistics"})
		return
	}

	c.JSON(http.StatusOK, stats)
}
