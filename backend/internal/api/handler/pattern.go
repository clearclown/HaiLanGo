package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/clearclown/HaiLanGo/backend/internal/repository"
	"github.com/clearclown/HaiLanGo/backend/internal/service/pattern"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// PatternHandler はパターンAPIのハンドラー
type PatternHandler struct {
	repo repository.PatternRepositoryInterface
}

// NewPatternHandler はパターンハンドラーを作成
func NewPatternHandler(repo repository.PatternRepositoryInterface) *PatternHandler {
	return &PatternHandler{
		repo: repo,
	}
}

// RegisterRoutes はパターンAPIのルートを登録
func (h *PatternHandler) RegisterRoutes(rg *gin.RouterGroup) {
	patterns := rg.Group("/patterns")
	{
		// パターン抽出
		patterns.POST("/extract", h.ExtractPatterns)

		// 書籍のパターン一覧取得
		patterns.GET("/books/:book_id", h.GetPatternsByBook)

		// パターン詳細取得
		patterns.GET("/:pattern_id", h.GetPatternByID)

		// パターンの使用例取得
		patterns.GET("/:pattern_id/examples", h.GetPatternExamples)

		// パターンの練習問題取得
		patterns.GET("/:pattern_id/practice", h.GetPatternPractice)

		// パターン学習進捗更新
		patterns.POST("/:pattern_id/progress", h.UpdatePatternProgress)
	}
}

// ExtractPatternsRequest はパターン抽出リクエスト
type ExtractPatternsRequest struct {
	BookID       uuid.UUID `json:"book_id" binding:"required"`
	MinFrequency int       `json:"min_frequency"`
	// 実際のページデータは省略（モック実装のため）
}

// ExtractPatterns はパターンを抽出
// POST /api/v1/patterns/extract
func (h *PatternHandler) ExtractPatterns(c *gin.Context) {
	_, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var req ExtractPatternsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// デフォルト最小頻度
	if req.MinFrequency == 0 {
		req.MinFrequency = 2
	}

	// サンプルページデータ（実際はDBから取得）
	samplePages := []pattern.PageText{
		{
			PageNumber:  1,
			Text:        "Здравствуйте! Как дела? Здравствуйте!",
			Translation: "こんにちは！調子はどう？こんにちは！",
		},
		{
			PageNumber:  2,
			Text:        "Спасибо, хорошо. А у вас? Здравствуйте!",
			Translation: "ありがとう、元気です。あなたは？こんにちは！",
		},
	}

	startTime := time.Now()

	patterns, err := h.repo.ExtractPatterns(c.Request.Context(), req.BookID, samplePages, req.MinFrequency)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to extract patterns"})
		return
	}

	duration := time.Since(startTime)

	c.JSON(http.StatusOK, gin.H{
		"patterns":        patterns,
		"total_found":     len(patterns),
		"processed_pages": len(samplePages),
		"duration_ms":     duration.Milliseconds(),
	})
}

// GetPatternsByBook は書籍のパターン一覧を取得
// GET /api/v1/patterns/books/:book_id
func (h *PatternHandler) GetPatternsByBook(c *gin.Context) {
	_, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	bookIDStr := c.Param("book_id")
	bookID, err := uuid.Parse(bookIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid book ID"})
		return
	}

	patterns, err := h.repo.GetPatternsByBookID(c.Request.Context(), bookID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get patterns"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"patterns": patterns,
		"book_id":  bookID,
		"count":    len(patterns),
	})
}

// GetPatternByID はパターン詳細を取得
// GET /api/v1/patterns/:pattern_id
func (h *PatternHandler) GetPatternByID(c *gin.Context) {
	_, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	patternIDStr := c.Param("pattern_id")
	patternID, err := uuid.Parse(patternIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid pattern ID"})
		return
	}

	pattern, err := h.repo.GetPatternByID(c.Request.Context(), patternID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get pattern"})
		return
	}

	if pattern == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Pattern not found"})
		return
	}

	c.JSON(http.StatusOK, pattern)
}

// GetPatternExamples はパターンの使用例を取得
// GET /api/v1/patterns/:pattern_id/examples?limit=10
func (h *PatternHandler) GetPatternExamples(c *gin.Context) {
	_, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	patternIDStr := c.Param("pattern_id")
	patternID, err := uuid.Parse(patternIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid pattern ID"})
		return
	}

	// クエリパラメータからlimitを取得（デフォルト: 10）
	limit := 10
	if limitStr := c.Query("limit"); limitStr != "" {
		if parsedLimit, err := strconv.Atoi(limitStr); err == nil && parsedLimit > 0 {
			limit = parsedLimit
		}
	}

	examples, err := h.repo.GetPatternExamples(c.Request.Context(), patternID, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get examples"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"examples":   examples,
		"pattern_id": patternID,
		"count":      len(examples),
	})
}

// GetPatternPractice はパターンの練習問題を取得
// GET /api/v1/patterns/:pattern_id/practice?count=10
func (h *PatternHandler) GetPatternPractice(c *gin.Context) {
	_, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	patternIDStr := c.Param("pattern_id")
	patternID, err := uuid.Parse(patternIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid pattern ID"})
		return
	}

	// クエリパラメータからcountを取得（デフォルト: 10）
	count := 10
	if countStr := c.Query("count"); countStr != "" {
		if parsedCount, err := strconv.Atoi(countStr); err == nil && parsedCount > 0 {
			count = parsedCount
		}
	}

	practices, err := h.repo.GetPatternPractice(c.Request.Context(), patternID, count)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get practice"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"practices":  practices,
		"pattern_id": patternID,
		"count":      len(practices),
	})
}

// UpdatePatternProgressRequest はパターン学習進捗更新リクエスト
type UpdatePatternProgressRequest struct {
	Correct bool `json:"correct"`
}

// UpdatePatternProgress はパターン学習進捗を更新
// POST /api/v1/patterns/:pattern_id/progress
func (h *PatternHandler) UpdatePatternProgress(c *gin.Context) {
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

	patternIDStr := c.Param("pattern_id")
	patternID, err := uuid.Parse(patternIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid pattern ID"})
		return
	}

	var req UpdatePatternProgressRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	progress, err := h.repo.UpdatePatternProgress(c.Request.Context(), userID, patternID, req.Correct)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update progress"})
		return
	}

	c.JSON(http.StatusOK, progress)
}
