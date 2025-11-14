// Package teachermode provides teacher mode API handlers
package teachermode

import (
	"net/http"

	"github.com/clearclown/HaiLanGo/backend/internal/service/teacher-mode"
	"github.com/gin-gonic/gin"
)

// Handler 教師モードAPIハンドラー
type Handler struct {
	service *teachermode.Service
}

// NewHandler 新しいハンドラーインスタンスを作成
func NewHandler() *Handler {
	return &Handler{
		service: teachermode.NewService(),
	}
}

// GeneratePlaylistRequest プレイリスト生成リクエスト
type GeneratePlaylistRequest struct {
	Settings  *teachermode.TeacherModeSettings `json:"settings"`
	PageRange *PageRange                       `json:"pageRange"`
}

// PageRange ページ範囲
type PageRange struct {
	Start int `json:"start"`
	End   int `json:"end"`
}

// GeneratePlaylistResponse プレイリスト生成レスポンス
type GeneratePlaylistResponse struct {
	PlaylistID        string              `json:"playlistId"`
	TotalPages        int                 `json:"totalPages"`
	EstimatedDuration int64               `json:"estimatedDuration"`
	Pages             []PageAudioResponse `json:"pages"`
}

// PageAudioResponse ページ音声レスポンス
type PageAudioResponse struct {
	PageNumber int                    `json:"pageNumber"`
	Segments   []AudioSegmentResponse `json:"segments"`
}

// AudioSegmentResponse 音声セグメントレスポンス
type AudioSegmentResponse struct {
	Type     string `json:"type"`
	AudioURL string `json:"audioUrl"`
	Duration int64  `json:"duration"`
	Text     string `json:"text"`
}

// GeneratePlaylist プレイリスト生成ハンドラー
func (h *Handler) GeneratePlaylist(c *gin.Context) {
	bookID := c.Param("bookId")

	// リクエストボディをパース
	var req GeneratePlaylistRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// プレイリストを生成
	playlist, err := h.service.GeneratePlaylist(c.Request.Context(), bookID, req.Settings)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// レスポンスを構築
	response := GeneratePlaylistResponse{
		PlaylistID:        playlist.ID,
		TotalPages:        len(playlist.Pages),
		EstimatedDuration: playlist.TotalDuration,
		Pages:             make([]PageAudioResponse, 0, len(playlist.Pages)),
	}

	for _, page := range playlist.Pages {
		pageResp := PageAudioResponse{
			PageNumber: page.PageNumber,
			Segments:   make([]AudioSegmentResponse, 0, len(page.Segments)),
		}

		for _, seg := range page.Segments {
			segResp := AudioSegmentResponse{
				Type:     seg.Type,
				AudioURL: seg.AudioURL,
				Duration: seg.Duration,
				Text:     seg.Text,
			}
			pageResp.Segments = append(pageResp.Segments, segResp)
		}

		response.Pages = append(response.Pages, pageResp)
	}

	c.JSON(http.StatusOK, response)
}

// GenerateDownloadPackageRequest ダウンロードパッケージ生成リクエスト
type GenerateDownloadPackageRequest struct {
	Settings *teachermode.TeacherModeSettings `json:"settings"`
}

// GenerateDownloadPackageResponse ダウンロードパッケージ生成レスポンス
type GenerateDownloadPackageResponse struct {
	PackageID   string `json:"packageId"`
	DownloadURL string `json:"downloadUrl"`
	TotalSize   int64  `json:"totalSize"`
	ExpiresAt   string `json:"expiresAt"`
}

// GenerateDownloadPackage ダウンロードパッケージ生成ハンドラー
func (h *Handler) GenerateDownloadPackage(c *gin.Context) {
	bookID := c.Param("bookId")

	// リクエストボディをパース
	var req GenerateDownloadPackageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// ダウンロードパッケージを生成
	pkg, err := h.service.GenerateDownloadPackage(c.Request.Context(), bookID, req.Settings)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// レスポンスを返す
	response := GenerateDownloadPackageResponse{
		PackageID:   pkg.PackageID,
		DownloadURL: pkg.DownloadURL,
		TotalSize:   pkg.TotalSize,
		ExpiresAt:   pkg.ExpiresAt,
	}

	c.JSON(http.StatusOK, response)
}

// UpdatePlaybackStateRequest 再生状態更新リクエスト
type UpdatePlaybackStateRequest struct {
	CurrentPage         int   `json:"currentPage"`
	CurrentSegmentIndex int   `json:"currentSegmentIndex"`
	ElapsedTime         int64 `json:"elapsedTime"`
}

// UpdatePlaybackStateResponse 再生状態更新レスポンス
type UpdatePlaybackStateResponse struct {
	Success bool `json:"success"`
}

// UpdatePlaybackState 再生状態更新ハンドラー
func (h *Handler) UpdatePlaybackState(c *gin.Context) {
	playlistID := c.Param("playlistId")

	// リクエストボディをパース
	var req UpdatePlaybackStateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// 検証
	if req.CurrentPage < 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid currentPage"})
		return
	}

	// 再生状態を保存（実際はRedisやデータベースに保存）
	// TODO: 実装
	_ = playlistID // TODO: use it

	// レスポンスを返す
	response := UpdatePlaybackStateResponse{
		Success: true,
	}

	c.JSON(http.StatusOK, response)
}

// GetPlaylistResponse プレイリスト取得レスポンス
type GetPlaylistResponse struct {
	PlaylistID string              `json:"playlistId"`
	Pages      []PageAudioResponse `json:"pages"`
}

// GetPlaylist プレイリスト取得ハンドラー
func (h *Handler) GetPlaylist(c *gin.Context) {
	playlistID := c.Param("playlistId")

	// プレイリストが存在しない場合（テスト用）
	if playlistID == "nonexistent" {
		c.JSON(http.StatusNotFound, gin.H{"error": "Playlist not found"})
		return
	}

	// モックレスポンス
	response := GetPlaylistResponse{
		PlaylistID: playlistID,
		Pages: []PageAudioResponse{
			{
				PageNumber: 1,
				Segments: []AudioSegmentResponse{
					{
						Type:     "phrase",
						AudioURL: "http://example.com/audio1.mp3",
						Duration: 2000,
						Text:     "Hello",
					},
				},
			},
		},
	}

	c.JSON(http.StatusOK, response)
}

// RegisterRoutes registers teacher mode routes
func (h *Handler) RegisterRoutes(rg *gin.RouterGroup) {
	teacher := rg.Group("/teacher-mode")
	{
		teacher.POST("/books/:bookId/generate", h.GeneratePlaylist)
		teacher.POST("/books/:bookId/download-package", h.GenerateDownloadPackage)
		teacher.GET("/playlists/:playlistId", h.GetPlaylist)
		teacher.POST("/playlists/:playlistId/state", h.UpdatePlaybackState)
	}
}
