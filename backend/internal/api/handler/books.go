package handler

import (
	"fmt"
	"net/http"
	"time"

	"github.com/clearclown/HaiLanGo/backend/internal/models"
	"github.com/clearclown/HaiLanGo/backend/internal/repository"
	"github.com/clearclown/HaiLanGo/backend/internal/websocket"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// BooksHandler は書籍APIのハンドラー
type BooksHandler struct {
	repo  repository.BookRepository
	wsHub *websocket.Hub
}

// NewBooksHandler は新しいBooksHandlerを作成
func NewBooksHandler(repo repository.BookRepository, wsHub *websocket.Hub) *BooksHandler {
	return &BooksHandler{
		repo:  repo,
		wsHub: wsHub,
	}
}

// CreateBookRequest は本の作成リクエスト
type CreateBookRequest struct {
	Title             string `json:"title" binding:"required"`
	TargetLanguage    string `json:"target_language" binding:"required"`
	NativeLanguage    string `json:"native_language" binding:"required"`
	ReferenceLanguage string `json:"reference_language,omitempty"`
}

// GetBooks godoc
// @Summary Get all books for user
// @Tags books
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string][]models.Book
// @Failure 401 {object} map[string]string
// @Router /api/v1/books [get]
func (h *BooksHandler) GetBooks(c *gin.Context) {
	// ミドルウェアからユーザーIDを取得
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

	books, err := h.repo.GetByUserID(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch books"})
		return
	}

	// booksがnilの場合は空配列を返す
	if books == nil {
		books = []*models.Book{}
	}

	c.JSON(http.StatusOK, gin.H{"books": books})
}

// GetBook godoc
// @Summary Get book by ID
// @Tags books
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Book ID"
// @Success 200 {object} models.Book
// @Failure 404 {object} map[string]string
// @Router /api/v1/books/{id} [get]
func (h *BooksHandler) GetBook(c *gin.Context) {
	bookIDStr := c.Param("id")
	bookID, err := uuid.Parse(bookIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid book ID"})
		return
	}

	// ミドルウェアからユーザーIDを取得
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

	book, err := h.repo.GetByID(c.Request.Context(), bookID)
	if err != nil {
		if err == repository.ErrBookNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Book not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch book"})
		return
	}

	// ユーザー所有権チェック
	if book.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden"})
		return
	}

	c.JSON(http.StatusOK, book)
}

// CreateBook godoc
// @Summary Create new book
// @Tags books
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param book body CreateBookRequest true "Book data"
// @Success 201 {object} map[string]models.Book
// @Failure 400 {object} map[string]string
// @Router /api/v1/books [post]
func (h *BooksHandler) CreateBook(c *gin.Context) {
	var req CreateBookRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body", "details": err.Error()})
		return
	}

	// ミドルウェアからユーザーIDを取得
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

	now := time.Now()
	book := &models.Book{
		ID:                uuid.New(),
		UserID:            userID,
		Title:             req.Title,
		TargetLanguage:    req.TargetLanguage,
		NativeLanguage:    req.NativeLanguage,
		ReferenceLanguage: req.ReferenceLanguage,
		TotalPages:        0,
		ProcessedPages:    0,
		Status:            models.BookStatusUploading,
		OCRStatus:         models.OCRStatusPending,
		CreatedAt:         now,
		UpdatedAt:         now,
	}

	if err := h.repo.Create(c.Request.Context(), book); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create book", "details": err.Error()})
		return
	}

	// WebSocket通知: 書籍が作成されたことを通知
	if h.wsHub != nil {
		message, err := websocket.NewNotificationMessage(
			"書籍を作成しました",
			fmt.Sprintf("「%s」の登録が完了しました", book.Title),
			websocket.NotificationLevel("success"),
		)
		if err == nil {
			h.wsHub.SendToUser(userID, message)
		}
	}

	c.JSON(http.StatusCreated, gin.H{"book": book})
}

// DeleteBook godoc
// @Summary Delete book
// @Tags books
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Book ID"
// @Success 200 {object} map[string]bool
// @Failure 404 {object} map[string]string
// @Router /api/v1/books/{id} [delete]
func (h *BooksHandler) DeleteBook(c *gin.Context) {
	bookIDStr := c.Param("id")
	bookID, err := uuid.Parse(bookIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid book ID"})
		return
	}

	// ミドルウェアからユーザーIDを取得
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

	// 所有権チェック
	book, err := h.repo.GetByID(c.Request.Context(), bookID)
	if err != nil {
		if err == repository.ErrBookNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Book not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch book"})
		return
	}

	if book.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden"})
		return
	}

	if err := h.repo.Delete(c.Request.Context(), bookID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete book"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true})
}

// RegisterRoutes registers book routes
func (h *BooksHandler) RegisterRoutes(rg *gin.RouterGroup) {
	books := rg.Group("/books")
	{
		books.GET("", h.GetBooks)
		books.POST("", h.CreateBook)
		books.GET("/:id", h.GetBook)
		books.DELETE("/:id", h.DeleteBook)
	}
}
