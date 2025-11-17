package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/clearclown/HaiLanGo/backend/internal/models"
	"github.com/clearclown/HaiLanGo/backend/internal/repository"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// TTSHandler はTTS APIのハンドラー
type TTSHandler struct {
	repo repository.TTSRepositoryInterface
}

// NewTTSHandler はTTSハンドラーを作成
func NewTTSHandler(repo repository.TTSRepositoryInterface) *TTSHandler {
	return &TTSHandler{
		repo: repo,
	}
}

// RegisterRoutes はTTS APIのルートを登録
func (h *TTSHandler) RegisterRoutes(rg *gin.RouterGroup) {
	tts := rg.Group("/tts")
	{
		// 音声合成
		tts.POST("/synthesize", h.Synthesize)

		// 音声ファイル取得
		tts.GET("/audio/:audioId", h.GetAudio)

		// サポート言語一覧
		tts.GET("/languages", h.GetLanguages)

		// バッチ音声生成
		tts.POST("/books/:bookId/batch", h.BatchSynthesize)

		// TTSジョブのステータス取得
		tts.GET("/jobs/:jobId", h.GetJobStatus)

		// 書籍のTTSジョブ一覧
		tts.GET("/books/:bookId/jobs", h.GetBookJobs)

		// キャッシュ統計
		tts.GET("/cache/stats", h.GetCacheStats)
	}
}

// Synthesize はテキストから音声を生成
// POST /api/v1/tts/synthesize
func (h *TTSHandler) Synthesize(c *gin.Context) {
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

	var req models.TTSRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// デフォルト値の設定
	if req.Options.Speed == 0 {
		req.Options.Speed = 1.0
	}
	if req.Options.Quality == "" {
		req.Options.Quality = models.TTSQualityStandard
	}
	if req.Options.Format == "" {
		req.Options.Format = models.TTSAudioFormatMP3
	}

	// ダミーのbookIDとpageNumberを使用（単一テキスト合成の場合）
	bookID := uuid.New()
	pageNumber := 0

	// TTSジョブを作成
	job, err := h.repo.CreateJob(c.Request.Context(), userID, bookID, pageNumber, req.Text, req.Language, req.Options)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create TTS job"})
		return
	}

	// バックグラウンドでTTS処理を開始
	go func() {
		if inMemRepo, ok := h.repo.(*repository.InMemoryTTSRepository); ok {
			inMemRepo.SimulateTTSProcessing(c.Request.Context(), job.ID)
		}
	}()

	response := &models.TTSJobResponse{
		JobID:     job.ID,
		Status:    job.Status,
		Progress:  job.Progress,
		CreatedAt: job.CreatedAt,
		UpdatedAt: job.UpdatedAt,
	}

	c.JSON(http.StatusAccepted, response)
}

// GetAudio は音声ファイルを取得（リダイレクト）
// GET /api/v1/tts/audio/:audioId
func (h *TTSHandler) GetAudio(c *gin.Context) {
	_, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	audioID := c.Param("audioId")

	audioURL, err := h.repo.GetAudioURL(c.Request.Context(), audioID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Audio not found"})
		return
	}

	// 実際の実装では、ここで音声ファイルにリダイレクトまたはストリーミング
	c.Redirect(http.StatusFound, audioURL)
}

// GetLanguages はサポート言語一覧を取得
// GET /api/v1/tts/languages
func (h *TTSHandler) GetLanguages(c *gin.Context) {
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

// BatchSynthesize はバッチ音声生成を開始
// POST /api/v1/tts/books/:bookId/batch
func (h *TTSHandler) BatchSynthesize(c *gin.Context) {
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

	var req models.TTSBatchRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// 実際の実装では書籍のページ数とテキストを取得
	totalPages := 100 // サンプル値

	jobIDs := make([]string, 0, totalPages)

	// 各ページのTTSジョブを作成
	for i := 1; i <= totalPages; i++ {
		text := "Sample text for page " + strconv.Itoa(i) // 実際にはOCR結果から取得

		job, err := h.repo.CreateJob(c.Request.Context(), userID, bookID, i, text, req.Language, req.Options)
		if err != nil {
			continue
		}

		jobIDs = append(jobIDs, job.ID)

		// バックグラウンドでTTS処理を開始
		go func(jobID string) {
			if inMemRepo, ok := h.repo.(*repository.InMemoryTTSRepository); ok {
				inMemRepo.SimulateTTSProcessing(c.Request.Context(), jobID)
			}
		}(job.ID)
	}

	response := &models.TTSBatchResponse{
		BatchID:    uuid.New().String(),
		BookID:     bookID.String(),
		TotalPages: totalPages,
		JobIDs:     jobIDs,
		CreatedAt:  time.Now(),
	}

	c.JSON(http.StatusAccepted, response)
}

// GetJobStatus はTTSジョブのステータスを取得
// GET /api/v1/tts/jobs/:jobId
func (h *TTSHandler) GetJobStatus(c *gin.Context) {
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

	response := &models.TTSJobResponse{
		JobID:      job.ID,
		BookID:     job.BookID,
		PageNumber: job.PageNumber,
		Status:     job.Status,
		Progress:   job.Progress,
		AudioID:    job.AudioID,
		AudioURL:   job.AudioURL,
		CreatedAt:  job.CreatedAt,
		UpdatedAt:  job.UpdatedAt,
	}

	c.JSON(http.StatusOK, response)
}

// GetBookJobs は書籍のTTSジョブ一覧を取得
// GET /api/v1/tts/books/:bookId/jobs
func (h *TTSHandler) GetBookJobs(c *gin.Context) {
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get TTS jobs"})
		return
	}

	responses := make([]*models.TTSJobResponse, 0, len(jobs))
	for _, job := range jobs {
		responses = append(responses, &models.TTSJobResponse{
			JobID:      job.ID,
			BookID:     job.BookID,
			PageNumber: job.PageNumber,
			Status:     job.Status,
			Progress:   job.Progress,
			AudioID:    job.AudioID,
			AudioURL:   job.AudioURL,
			CreatedAt:  job.CreatedAt,
			UpdatedAt:  job.UpdatedAt,
		})
	}

	c.JSON(http.StatusOK, responses)
}

// GetCacheStats はキャッシュ統計情報を取得
// GET /api/v1/tts/cache/stats
func (h *TTSHandler) GetCacheStats(c *gin.Context) {
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

	stats, err := h.repo.GetCacheStats(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get cache stats"})
		return
	}

	c.JSON(http.StatusOK, stats)
}
