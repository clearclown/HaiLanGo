package handler

import (
	"net/http"
	"strconv"

	"github.com/clearclown/HaiLanGo/backend/internal/models"
	"github.com/clearclown/HaiLanGo/backend/internal/service/pattern"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// PatternHandler handles pattern-related HTTP requests
type PatternHandler struct {
	extractor *pattern.Extractor
}

// NewPatternHandler creates a new pattern handler
func NewPatternHandler() *PatternHandler {
	return &PatternHandler{
		extractor: pattern.NewExtractor(),
	}
}

// ExtractPatterns handles POST /api/v1/books/:bookId/patterns/extract
func (h *PatternHandler) ExtractPatterns(c *gin.Context) {
	bookIDStr := c.Param("bookId")

	// Parse request
	var req models.PatternExtractionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Parse bookID
	bookID, err := uuid.Parse(bookIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid book ID"})
		return
	}
	req.BookID = bookID

	// TODO: Fetch book pages from database
	// For now, this is a placeholder
	pages := []pattern.PageText{} // Would fetch from database

	// Extract patterns
	patterns, err := h.extractor.ExtractPatterns(c.Request.Context(), req.BookID, pages, req.MinFrequency)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to extract patterns"})
		return
	}

	// Create response
	resp := models.PatternExtractionResponse{
		Patterns:       patterns,
		TotalFound:     len(patterns),
		ProcessedPages: req.PageEnd - req.PageStart + 1,
	}

	c.JSON(http.StatusOK, resp)
}

// GetPatterns handles GET /api/v1/books/:bookId/patterns
func (h *PatternHandler) GetPatterns(c *gin.Context) {
	// Extract book_id from URL
	bookIDStr := c.Param("bookId")
	bookID, err := uuid.Parse(bookIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid book ID"})
		return
	}

	// TODO: Fetch patterns from database
	// For now, return empty list
	patterns := []models.Pattern{}

	c.JSON(http.StatusOK, gin.H{
		"patterns": patterns,
		"book_id":  bookID,
	})
}

// GetPatternPractice handles GET /api/v1/patterns/:patternId/practice
func (h *PatternHandler) GetPatternPractice(c *gin.Context) {
	// Extract pattern_id from URL
	patternIDStr := c.Param("patternId")
	patternID, err := uuid.Parse(patternIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid pattern ID"})
		return
	}

	// Extract count parameter (default: 10)
	count := 10
	if countStr := c.Query("count"); countStr != "" {
		if c, err := strconv.Atoi(countStr); err == nil {
			count = c
		}
	}

	// TODO: Fetch pattern and generate practice exercises from database
	// For now, return placeholder
	practices := []models.PatternPractice{}

	c.JSON(http.StatusOK, gin.H{
		"pattern_id": patternID,
		"practices":  practices,
		"count":      count,
	})
}

// RegisterRoutes registers pattern routes
func (h *PatternHandler) RegisterRoutes(rg *gin.RouterGroup) {
	patterns := rg.Group("/patterns")
	{
		patterns.GET("/:patternId/practice", h.GetPatternPractice)
	}

	// Book-specific pattern routes
	books := rg.Group("/books")
	{
		books.POST("/:bookId/patterns/extract", h.ExtractPatterns)
		books.GET("/:bookId/patterns", h.GetPatterns)
	}
}
