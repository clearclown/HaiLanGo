package handler

import (
	"net/http"
	"strconv"

	"github.com/clearclown/HaiLanGo/backend/internal/models"
	"github.com/clearclown/HaiLanGo/backend/internal/repository"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// LearningHandler handles learning-related HTTP requests
type LearningHandler struct {
	repo repository.LearningRepositoryInterface
}

// NewLearningHandler creates a new learning handler
func NewLearningHandler(repo repository.LearningRepositoryInterface) *LearningHandler {
	return &LearningHandler{
		repo: repo,
	}
}

// GetPageLearning handles GET /api/v1/learning/books/:bookId/pages/:pageNumber
// @Summary Get learning page data
// @Description Get page data for learning including OCR, phrases, vocabulary
// @Tags learning
// @Accept json
// @Produce json
// @Param bookId path string true "Book ID"
// @Param pageNumber path int true "Page Number"
// @Security BearerAuth
// @Success 200 {object} models.PageLearning
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/learning/books/{bookId}/pages/{pageNumber} [get]
func (h *LearningHandler) GetPageLearning(c *gin.Context) {
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

	bookIDStr := c.Param("bookId")
	bookID, err := uuid.Parse(bookIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid book ID"})
		return
	}

	pageNumberStr := c.Param("pageNumber")
	pageNumber, err := strconv.Atoi(pageNumberStr)
	if err != nil || pageNumber < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid page number"})
		return
	}

	pageLearning, err := h.repo.GetPageLearning(c.Request.Context(), userID, bookID, pageNumber)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Page not found"})
		return
	}

	c.JSON(http.StatusOK, pageLearning)
}

// CompletePage handles POST /api/v1/learning/books/:bookId/pages/:pageNumber/complete
// @Summary Mark page as completed
// @Description Mark a learning page as completed
// @Tags learning
// @Accept json
// @Produce json
// @Param bookId path string true "Book ID"
// @Param pageNumber path int true "Page Number"
// @Param request body models.CompletePageRequest true "Complete page request"
// @Security BearerAuth
// @Success 200 {object} models.CompletePageResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/learning/books/{bookId}/pages/{pageNumber}/complete [post]
func (h *LearningHandler) CompletePage(c *gin.Context) {
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

	bookIDStr := c.Param("bookId")
	bookID, err := uuid.Parse(bookIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid book ID"})
		return
	}

	pageNumberStr := c.Param("pageNumber")
	pageNumber, err := strconv.Atoi(pageNumberStr)
	if err != nil || pageNumber < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid page number"})
		return
	}

	var req models.CompletePageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	progress, err := h.repo.CompletePage(c.Request.Context(), userID, bookID, pageNumber, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to complete page"})
		return
	}

	response := models.CompletePageResponse{
		Message:  "Page marked as completed",
		Progress: *progress,
		NextPage: pageNumber + 1,
	}

	c.JSON(http.StatusOK, response)
}

// RecordSession handles POST /api/v1/learning/books/:bookId/pages/:pageNumber/session
// @Summary Record learning session
// @Description Record a learning session (start/end)
// @Tags learning
// @Accept json
// @Produce json
// @Param bookId path string true "Book ID"
// @Param pageNumber path int true "Page Number"
// @Param request body models.SessionRequest true "Session request"
// @Security BearerAuth
// @Success 200 {object} models.SessionResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/learning/books/{bookId}/pages/{pageNumber}/session [post]
func (h *LearningHandler) RecordSession(c *gin.Context) {
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

	bookIDStr := c.Param("bookId")
	bookID, err := uuid.Parse(bookIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid book ID"})
		return
	}

	pageNumberStr := c.Param("pageNumber")
	pageNumber, err := strconv.Atoi(pageNumberStr)
	if err != nil || pageNumber < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid page number"})
		return
	}

	var req models.SessionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	session, err := h.repo.RecordSession(c.Request.Context(), userID, bookID, pageNumber, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to record session"})
		return
	}

	c.JSON(http.StatusOK, session)
}

// GetBookProgress handles GET /api/v1/learning/books/:bookId/progress
// @Summary Get book progress
// @Description Get overall progress for a book
// @Tags learning
// @Accept json
// @Produce json
// @Param bookId path string true "Book ID"
// @Security BearerAuth
// @Success 200 {object} models.BookProgressSummary
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/learning/books/{bookId}/progress [get]
func (h *LearningHandler) GetBookProgress(c *gin.Context) {
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

	bookIDStr := c.Param("bookId")
	bookID, err := uuid.Parse(bookIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid book ID"})
		return
	}

	progress, err := h.repo.GetBookProgress(c.Request.Context(), userID, bookID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get book progress"})
		return
	}

	c.JSON(http.StatusOK, progress)
}

// RegisterRoutes registers learning routes
func (h *LearningHandler) RegisterRoutes(rg *gin.RouterGroup) {
	learning := rg.Group("/learning")
	{
		learning.GET("/books/:bookId/pages/:pageNumber", h.GetPageLearning)
		learning.POST("/books/:bookId/pages/:pageNumber/complete", h.CompletePage)
		learning.POST("/books/:bookId/pages/:pageNumber/session", h.RecordSession)
		learning.GET("/books/:bookId/progress", h.GetBookProgress)
	}
}
