package learning

import (
	"context"
	"net/http"
	"strconv"

	"github.com/clearclown/HaiLanGo/backend/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Service は学習サービスのインターフェース
type Service interface {
	GetPage(ctx context.Context, bookID uuid.UUID, pageNumber int) (*models.PageWithProgress, error)
	MarkPageCompleted(ctx context.Context, userID, bookID uuid.UUID, pageNumber int, studyTime int) error
	GetProgress(ctx context.Context, userID, bookID uuid.UUID) (*models.LearningProgress, error)
}

// Handler は学習APIのハンドラー
type Handler struct {
	service Service
}

// NewHandler は新しいHandlerを作成する
func NewHandler(service Service) *Handler {
	return &Handler{
		service: service,
	}
}

// GetPage はページを取得するハンドラー
// GET /api/v1/books/:bookId/pages/:pageNumber
func (h *Handler) GetPage(c *gin.Context) {
	// パラメータ取得
	bookIDStr := c.Param("bookId")
	pageNumberStr := c.Param("pageNumber")

	// bookIDの検証
	bookID, err := uuid.Parse(bookIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid book ID"})
		return
	}

	// pageNumberの検証
	pageNumber, err := strconv.Atoi(pageNumberStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid page number"})
		return
	}

	// ページ取得
	page, err := h.service.GetPage(c.Request.Context(), bookID, pageNumber)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, page)
}

// MarkPageCompletedRequest はページ完了リクエストの構造体
type MarkPageCompletedRequest struct {
	UserID    string `json:"userId" binding:"required"`
	StudyTime int    `json:"studyTime"` // 秒単位
}

// MarkPageCompleted はページを完了としてマークするハンドラー
// POST /api/v1/books/:bookId/pages/:pageNumber/complete
func (h *Handler) MarkPageCompleted(c *gin.Context) {
	// パラメータ取得
	bookIDStr := c.Param("bookId")
	pageNumberStr := c.Param("pageNumber")

	// bookIDの検証
	bookID, err := uuid.Parse(bookIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid book ID"})
		return
	}

	// pageNumberの検証
	pageNumber, err := strconv.Atoi(pageNumberStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid page number"})
		return
	}

	// リクエストボディの解析
	var req MarkPageCompletedRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// userIDの検証
	userID, err := uuid.Parse(req.UserID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// ページ完了マーク
	if err := h.service.MarkPageCompleted(c.Request.Context(), userID, bookID, pageNumber, req.StudyTime); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Page completed successfully"})
}

// GetProgress は学習進捗を取得するハンドラー
// GET /api/v1/books/:bookId/progress?userId=xxx
func (h *Handler) GetProgress(c *gin.Context) {
	// パラメータ取得
	bookIDStr := c.Param("bookId")
	userIDStr := c.Query("userId")

	// bookIDの検証
	bookID, err := uuid.Parse(bookIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid book ID"})
		return
	}

	// userIDの検証
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// 進捗取得
	progress, err := h.service.GetProgress(c.Request.Context(), userID, bookID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, progress)
}

// RegisterRoutes registers learning routes
func (h *Handler) RegisterRoutes(rg *gin.RouterGroup) {
	learning := rg.Group("/learning")
	{
		learning.GET("/books/:bookId/pages/:pageNumber", h.GetPage)
		learning.POST("/books/:bookId/pages/:pageNumber/complete", h.MarkPageCompleted)
		learning.GET("/books/:bookId/progress", h.GetProgress)
	}
}
