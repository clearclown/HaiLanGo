package handler

import (
	"net/http"
	"strconv"

	"github.com/clearclown/HaiLanGo/backend/internal/models"
	"github.com/clearclown/HaiLanGo/backend/internal/service/vocabulary"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// VocabularyHandler は単語帳APIのハンドラー
type VocabularyHandler struct {
	service vocabulary.VocabularyService
}

// NewVocabularyHandler は新しいVocabularyHandlerを作成する
func NewVocabularyHandler(service vocabulary.VocabularyService) *VocabularyHandler {
	return &VocabularyHandler{
		service: service,
	}
}

// RegisterRoutes は単語帳APIのルートを登録する
func (h *VocabularyHandler) RegisterRoutes(rg *gin.RouterGroup) {
	vocab := rg.Group("/vocabulary")
	{
		vocab.GET("", h.ListWords)
		vocab.POST("", h.AddWord)
		vocab.GET("/stats", h.GetStats)
		vocab.GET("/export", h.ExportCSV)
		vocab.GET("/:id", h.GetWord)
		vocab.PUT("/:id", h.UpdateWord)
		vocab.DELETE("/:id", h.DeleteWord)
		vocab.POST("/:id/review", h.RecordReview)
		vocab.POST("/:id/tags", h.AddTags)
		vocab.POST("/auto-collect", h.AutoCollect)
	}
}

// AddWordRequest は単語追加リクエスト
type AddWordRequest struct {
	BookID        uuid.UUID `json:"book_id"`
	PageNumber    int       `json:"page_number"`
	Text          string    `json:"text" binding:"required"`
	Meaning       string    `json:"meaning"`
	Pronunciation string    `json:"pronunciation"`
	PartOfSpeech  string    `json:"part_of_speech"`
	Example       string    `json:"example"`
	Language      string    `json:"language" binding:"required"`
	Tags          []string  `json:"tags"`
}

// UpdateWordRequest は単語更新リクエスト
type UpdateWordRequest struct {
	Text          string   `json:"text"`
	Meaning       string   `json:"meaning"`
	Pronunciation string   `json:"pronunciation"`
	PartOfSpeech  string   `json:"part_of_speech"`
	Example       string   `json:"example"`
	Tags          []string `json:"tags"`
}

// RecordReviewRequest は学習記録リクエスト
type RecordReviewRequest struct {
	Score float64 `json:"score" binding:"required,min=0,max=100"`
}

// AddTagsRequest はタグ追加リクエスト
type AddTagsRequest struct {
	Tags []string `json:"tags" binding:"required"`
}

// AutoCollectRequest は自動収集リクエスト
type AutoCollectRequest struct {
	BookID     uuid.UUID `json:"book_id" binding:"required"`
	PageNumber int       `json:"page_number" binding:"required"`
	Text       string    `json:"text" binding:"required"`
	Language   string    `json:"language" binding:"required"`
}

// ListWords は単語一覧を取得する
// GET /api/v1/vocabulary?book_id=xxx&language=ru&query=hello&limit=20&offset=0
func (h *VocabularyHandler) ListWords(c *gin.Context) {
	userIDVal, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	userID := userIDVal.(uuid.UUID)

	// フィルタパラメータを解析
	filter := &models.WordFilter{
		UserID:    userID,
		Language:  c.Query("language"),
		Query:     c.Query("query"),
		SortBy:    c.DefaultQuery("sort_by", "created_at"),
		SortOrder: c.DefaultQuery("sort_order", "desc"),
	}

	// book_id
	if bookIDStr := c.Query("book_id"); bookIDStr != "" {
		bookID, err := uuid.Parse(bookIDStr)
		if err == nil {
			filter.BookID = bookID
		}
	}

	// tags
	if tagsStr := c.Query("tags"); tagsStr != "" {
		filter.Tags = []string{tagsStr}
	}

	// min_mastery
	if minMasteryStr := c.Query("min_mastery"); minMasteryStr != "" {
		if minMastery, err := strconv.ParseFloat(minMasteryStr, 64); err == nil {
			filter.MinMastery = minMastery
		}
	}

	// max_mastery
	if maxMasteryStr := c.Query("max_mastery"); maxMasteryStr != "" {
		if maxMastery, err := strconv.ParseFloat(maxMasteryStr, 64); err == nil {
			filter.MaxMastery = maxMastery
		}
	}

	// limit
	if limitStr := c.Query("limit"); limitStr != "" {
		if limit, err := strconv.Atoi(limitStr); err == nil && limit > 0 {
			filter.Limit = limit
		}
	}
	if filter.Limit == 0 {
		filter.Limit = 50 // default
	}

	// offset
	if offsetStr := c.Query("offset"); offsetStr != "" {
		if offset, err := strconv.Atoi(offsetStr); err == nil && offset >= 0 {
			filter.Offset = offset
		}
	}

	words, err := h.service.GetWords(c.Request.Context(), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get words"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"words": words,
		"count": len(words),
	})
}

// AddWord は新しい単語を追加する
// POST /api/v1/vocabulary
func (h *VocabularyHandler) AddWord(c *gin.Context) {
	userIDVal, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	userID := userIDVal.(uuid.UUID)

	var req AddWordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	word := &models.Word{
		UserID:        userID,
		BookID:        req.BookID,
		PageNumber:    req.PageNumber,
		Text:          req.Text,
		Meaning:       req.Meaning,
		Pronunciation: req.Pronunciation,
		PartOfSpeech:  req.PartOfSpeech,
		Example:       req.Example,
		Language:      req.Language,
		Tags:          req.Tags,
		Mastery:       0,
		ReviewCount:   0,
		AverageScore:  0,
	}

	if word.Tags == nil {
		word.Tags = []string{}
	}

	if err := h.service.AddWord(c.Request.Context(), word); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add word"})
		return
	}

	c.JSON(http.StatusCreated, word)
}

// GetWord は単語を取得する
// GET /api/v1/vocabulary/:id
func (h *VocabularyHandler) GetWord(c *gin.Context) {
	_, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid word ID"})
		return
	}

	word, err := h.service.GetWordByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get word"})
		return
	}

	if word == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Word not found"})
		return
	}

	c.JSON(http.StatusOK, word)
}

// UpdateWord は単語を更新する
// PUT /api/v1/vocabulary/:id
func (h *VocabularyHandler) UpdateWord(c *gin.Context) {
	_, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid word ID"})
		return
	}

	// 既存の単語を取得
	word, err := h.service.GetWordByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get word"})
		return
	}
	if word == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Word not found"})
		return
	}

	var req UpdateWordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// 更新可能なフィールドのみ更新
	if req.Text != "" {
		word.Text = req.Text
	}
	if req.Meaning != "" {
		word.Meaning = req.Meaning
	}
	if req.Pronunciation != "" {
		word.Pronunciation = req.Pronunciation
	}
	if req.PartOfSpeech != "" {
		word.PartOfSpeech = req.PartOfSpeech
	}
	if req.Example != "" {
		word.Example = req.Example
	}
	if req.Tags != nil {
		word.Tags = req.Tags
	}

	if err := h.service.UpdateWord(c.Request.Context(), word); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update word"})
		return
	}

	c.JSON(http.StatusOK, word)
}

// DeleteWord は単語を削除する
// DELETE /api/v1/vocabulary/:id
func (h *VocabularyHandler) DeleteWord(c *gin.Context) {
	_, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid word ID"})
		return
	}

	if err := h.service.DeleteWord(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete word"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Word deleted successfully"})
}

// RecordReview は学習記録を保存する
// POST /api/v1/vocabulary/:id/review
func (h *VocabularyHandler) RecordReview(c *gin.Context) {
	_, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid word ID"})
		return
	}

	var req RecordReviewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	if err := h.service.RecordReview(c.Request.Context(), id, req.Score); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to record review"})
		return
	}

	// 更新された単語を取得して返す
	word, err := h.service.GetWordByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get updated word"})
		return
	}

	c.JSON(http.StatusOK, word)
}

// AddTags は単語にタグを追加する
// POST /api/v1/vocabulary/:id/tags
func (h *VocabularyHandler) AddTags(c *gin.Context) {
	_, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid word ID"})
		return
	}

	var req AddTagsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	if err := h.service.AddTags(c.Request.Context(), id, req.Tags); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add tags"})
		return
	}

	// 更新された単語を取得して返す
	word, err := h.service.GetWordByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get updated word"})
		return
	}

	c.JSON(http.StatusOK, word)
}

// GetStats は単語統計を取得する
// GET /api/v1/vocabulary/stats?book_id=xxx
func (h *VocabularyHandler) GetStats(c *gin.Context) {
	userIDVal, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	userID := userIDVal.(uuid.UUID)

	var bookID uuid.UUID
	if bookIDStr := c.Query("book_id"); bookIDStr != "" {
		var err error
		bookID, err = uuid.Parse(bookIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid book_id"})
			return
		}
	}

	stats, err := h.service.GetStats(c.Request.Context(), userID, bookID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get stats"})
		return
	}

	c.JSON(http.StatusOK, stats)
}

// ExportCSV は単語をCSV形式でエクスポートする
// GET /api/v1/vocabulary/export?book_id=xxx
func (h *VocabularyHandler) ExportCSV(c *gin.Context) {
	userIDVal, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	userID := userIDVal.(uuid.UUID)

	filter := &models.WordFilter{
		UserID: userID,
	}

	if bookIDStr := c.Query("book_id"); bookIDStr != "" {
		bookID, err := uuid.Parse(bookIDStr)
		if err == nil {
			filter.BookID = bookID
		}
	}

	csvData, err := h.service.ExportWordsToCSV(c.Request.Context(), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to export words"})
		return
	}

	c.Header("Content-Type", "text/csv; charset=utf-8")
	c.Header("Content-Disposition", "attachment; filename=vocabulary.csv")
	c.Data(http.StatusOK, "text/csv", csvData)
}

// AutoCollect はテキストから単語を自動収集する
// POST /api/v1/vocabulary/auto-collect
func (h *VocabularyHandler) AutoCollect(c *gin.Context) {
	userIDVal, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	userID := userIDVal.(uuid.UUID)

	var req AutoCollectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	if err := h.service.AutoCollectWords(c.Request.Context(), userID, req.BookID, req.PageNumber, req.Text, req.Language); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to auto-collect words"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Words collected successfully"})
}
