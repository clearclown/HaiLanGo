package handler

import (
	"net/http"

	"github.com/clearclown/HaiLanGo/backend/internal/models"
	"github.com/clearclown/HaiLanGo/backend/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// TeacherModeHandler は教師モードAPIのハンドラー
type TeacherModeHandler struct {
	service *service.TeacherModeService
}

// NewTeacherModeHandler は新しいTeacherModeHandlerを作成
func NewTeacherModeHandler(service *service.TeacherModeService) *TeacherModeHandler {
	return &TeacherModeHandler{
		service: service,
	}
}

// GeneratePlaylistRequest はプレイリスト生成リクエスト
type GeneratePlaylistRequest struct {
	Settings  models.TeacherModeSettings `json:"settings" binding:"required"`
	PageRange *models.PageRange          `json:"page_range,omitempty"`
}

// GeneratePlaylistResponse はプレイリスト生成レスポンス
type GeneratePlaylistResponse struct {
	PlaylistID        string                  `json:"playlist_id"`
	TotalPages        int                     `json:"total_pages"`
	EstimatedDuration int                     `json:"estimated_duration"` // 秒
	Pages             []models.PageAudio      `json:"pages"`
}

// GenerateDownloadPackageRequest はダウンロードパッケージ生成リクエスト
type GenerateDownloadPackageRequest struct {
	Settings models.TeacherModeSettings `json:"settings" binding:"required"`
}

// GenerateDownloadPackageResponse はダウンロードパッケージ生成レスポンス
type GenerateDownloadPackageResponse struct {
	PackageID   string `json:"package_id"`
	DownloadURL string `json:"download_url"`
	TotalSize   int64  `json:"total_size"`
	ExpiresAt   string `json:"expires_at"`
}

// UpdatePlaybackStateRequest は再生状態更新リクエスト
type UpdatePlaybackStateRequest struct {
	CurrentPage         int `json:"current_page" binding:"required"`
	CurrentSegmentIndex int `json:"current_segment_index" binding:"required"`
	ElapsedTime         int `json:"elapsed_time" binding:"required"` // ミリ秒
}

// RegisterRoutes はルートを登録する
func (h *TeacherModeHandler) RegisterRoutes(r *gin.RouterGroup) {
	teacherMode := r.Group("/books/:id/teacher-mode")
	{
		teacherMode.POST("/generate", h.GeneratePlaylist)
		teacherMode.POST("/download-package", h.GenerateDownloadPackage)
		teacherMode.PUT("/playback-state", h.UpdatePlaybackState)
		teacherMode.GET("/playback-state", h.GetPlaybackState)
	}
}

// GeneratePlaylist godoc
// @Summary Generate teacher mode playlist
// @Tags teacher-mode
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Book ID"
// @Param request body GeneratePlaylistRequest true "Generate playlist request"
// @Success 200 {object} GeneratePlaylistResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/books/{id}/teacher-mode/generate [post]
func (h *TeacherModeHandler) GeneratePlaylist(c *gin.Context) {
	// ユーザーIDを取得
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

	// 書籍IDを取得
	bookIDStr := c.Param("id")
	bookID, err := uuid.Parse(bookIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid book ID"})
		return
	}

	// リクエストをパース
	var req GeneratePlaylistRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// プレイリストを生成
	playlist, err := h.service.GeneratePlaylist(
		c.Request.Context(),
		userID,
		bookID,
		&req.Settings,
		req.PageRange,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// レスポンスを作成
	response := GeneratePlaylistResponse{
		PlaylistID:        playlist.ID,
		TotalPages:        len(playlist.Pages),
		EstimatedDuration: playlist.TotalDuration / 1000, // ミリ秒を秒に変換
		Pages:             playlist.Pages,
	}

	c.JSON(http.StatusOK, response)
}

// GenerateDownloadPackage godoc
// @Summary Generate teacher mode download package
// @Tags teacher-mode
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Book ID"
// @Param request body GenerateDownloadPackageRequest true "Generate download package request"
// @Success 200 {object} GenerateDownloadPackageResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/books/{id}/teacher-mode/download-package [post]
func (h *TeacherModeHandler) GenerateDownloadPackage(c *gin.Context) {
	// ユーザーIDを取得
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

	// 書籍IDを取得
	bookIDStr := c.Param("id")
	bookID, err := uuid.Parse(bookIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid book ID"})
		return
	}

	// リクエストをパース
	var req GenerateDownloadPackageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// ダウンロードパッケージを生成
	packageID, downloadURL, totalSize, expiresAt, err := h.service.GenerateDownloadPackage(
		c.Request.Context(),
		userID,
		bookID,
		&req.Settings,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// レスポンスを作成
	response := GenerateDownloadPackageResponse{
		PackageID:   packageID.String(),
		DownloadURL: downloadURL,
		TotalSize:   totalSize,
		ExpiresAt:   expiresAt.Format("2006-01-02T15:04:05Z"),
	}

	c.JSON(http.StatusOK, response)
}

// UpdatePlaybackState godoc
// @Summary Update teacher mode playback state
// @Tags teacher-mode
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Book ID"
// @Param request body UpdatePlaybackStateRequest true "Update playback state request"
// @Success 200 {object} map[string]bool
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/books/{id}/teacher-mode/playback-state [put]
func (h *TeacherModeHandler) UpdatePlaybackState(c *gin.Context) {
	// ユーザーIDを取得
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

	// 書籍IDを取得
	bookIDStr := c.Param("id")
	bookID, err := uuid.Parse(bookIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid book ID"})
		return
	}

	// リクエストをパース
	var req UpdatePlaybackStateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 再生状態を更新
	state := &models.PlaybackState{
		Status:              models.PlaybackStatusPlaying,
		CurrentPage:         req.CurrentPage,
		CurrentSegmentIndex: req.CurrentSegmentIndex,
		ElapsedTime:         req.ElapsedTime,
		TotalDuration:       0,
	}

	if err := h.service.UpdatePlaybackState(c.Request.Context(), userID, bookID, state); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true})
}

// GetPlaybackState godoc
// @Summary Get teacher mode playback state
// @Tags teacher-mode
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Book ID"
// @Success 200 {object} models.PlaybackState
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/books/{id}/teacher-mode/playback-state [get]
func (h *TeacherModeHandler) GetPlaybackState(c *gin.Context) {
	// ユーザーIDを取得
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

	// 書籍IDを取得
	bookIDStr := c.Param("id")
	bookID, err := uuid.Parse(bookIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid book ID"})
		return
	}

	// 再生状態を取得
	state, err := h.service.GetPlaybackState(c.Request.Context(), userID, bookID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, state)
}
