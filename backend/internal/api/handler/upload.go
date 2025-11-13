package handler

import (
	"net/http"

	"github.com/clearclown/HaiLanGo/internal/models"
	"github.com/clearclown/HaiLanGo/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// UploadHandler はファイルアップロードのHTTPハンドラー
type UploadHandler struct {
	uploadService *service.UploadService
}

// NewUploadHandler はUploadHandlerの新しいインスタンスを作成する
func NewUploadHandler(uploadService *service.UploadService) *UploadHandler {
	return &UploadHandler{
		uploadService: uploadService,
	}
}

// CreateBook は新しい書籍を作成するハンドラー
// POST /api/v1/books
func (h *UploadHandler) CreateBook(c *gin.Context) {
	// TODO: 実際の実装では認証ミドルウェアからユーザーIDを取得
	userID := uuid.New()

	var metadata models.BookMetadata
	if err := c.ShouldBindJSON(&metadata); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request body",
			"details": err.Error(),
		})
		return
	}

	book, err := h.uploadService.CreateBook(c.Request.Context(), userID, metadata)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to create book",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, book)
}

// UploadFiles はファイルをアップロードするハンドラー
// POST /api/v1/books/:book_id/upload
func (h *UploadHandler) UploadFiles(c *gin.Context) {
	// パラメータからbook_idを取得
	bookIDStr := c.Param("book_id")
	bookID, err := uuid.Parse(bookIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid book_id",
		})
		return
	}

	// TODO: 実際の実装では認証ミドルウェアからユーザーIDを取得
	userID := uuid.New()

	// マルチパートフォームを解析
	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "failed to parse multipart form",
			"details": err.Error(),
		})
		return
	}

	files := form.File["files"]
	if len(files) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "no files provided",
		})
		return
	}

	// ファイルをアップロード
	bookFiles, err := h.uploadService.UploadMultipleFiles(c.Request.Context(), userID, bookID, files)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to upload files",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "files uploaded successfully",
		"files": bookFiles,
		"count": len(bookFiles),
	})
}

// GetUploadProgress はアップロード進捗を取得するハンドラー
// GET /api/v1/books/:book_id/upload-status
func (h *UploadHandler) GetUploadProgress(c *gin.Context) {
	bookIDStr := c.Param("book_id")
	bookID, err := uuid.Parse(bookIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid book_id",
		})
		return
	}

	progress, err := h.uploadService.GetUploadProgress(c.Request.Context(), bookID)
	if err != nil {
		if err == service.ErrBookNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "upload progress not found",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to get upload progress",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, progress)
}

// RegisterRoutes はルートを登録する
func (h *UploadHandler) RegisterRoutes(router *gin.RouterGroup) {
	books := router.Group("/books")
	{
		books.POST("", h.CreateBook)
		books.POST("/:book_id/upload", h.UploadFiles)
		books.GET("/:book_id/upload-status", h.GetUploadProgress)
	}
}
