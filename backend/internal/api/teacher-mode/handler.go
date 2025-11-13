// Package teachermode provides teacher mode API handlers
package teachermode

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/clearclown/HaiLanGo/backend/internal/service/teacher-mode"
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
func (h *Handler) GeneratePlaylist(w http.ResponseWriter, r *http.Request) {
	// 認証チェック
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// リクエストボディをパース
	var req GeneratePlaylistRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// bookIDをURLから取得（実際のルーターで処理）
	bookID := "test-book" // TODO: ルーターから取得

	// プレイリストを生成
	playlist, err := h.service.GeneratePlaylist(r.Context(), bookID, req.Settings)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
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

	// JSONレスポンスを返す
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
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
func (h *Handler) GenerateDownloadPackage(w http.ResponseWriter, r *http.Request) {
	// 認証チェック
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// リクエストボディをパース
	var req GenerateDownloadPackageRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// bookIDをURLから取得
	bookID := "test-book" // TODO: ルーターから取得

	// ダウンロードパッケージを生成
	pkg, err := h.service.GenerateDownloadPackage(r.Context(), bookID, req.Settings)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// レスポンスを返す
	response := GenerateDownloadPackageResponse{
		PackageID:   pkg.PackageID,
		DownloadURL: pkg.DownloadURL,
		TotalSize:   pkg.TotalSize,
		ExpiresAt:   pkg.ExpiresAt,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
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
func (h *Handler) UpdatePlaybackState(w http.ResponseWriter, r *http.Request) {
	// 認証チェック
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// リクエストボディをパース
	var req UpdatePlaybackStateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// 検証
	if req.CurrentPage < 0 {
		http.Error(w, "Invalid currentPage", http.StatusBadRequest)
		return
	}

	// 再生状態を保存（実際はRedisやデータベースに保存）
	// TODO: 実装

	// レスポンスを返す
	response := UpdatePlaybackStateResponse{
		Success: true,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// GetPlaylistResponse プレイリスト取得レスポンス
type GetPlaylistResponse struct {
	PlaylistID string              `json:"playlistId"`
	Pages      []PageAudioResponse `json:"pages"`
}

// GetPlaylist プレイリスト取得ハンドラー
func (h *Handler) GetPlaylist(w http.ResponseWriter, r *http.Request) {
	// 認証チェック
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// bookIDをURLから取得
	bookID := "test-book" // TODO: ルーターから取得

	// プレイリストが存在しない場合（テスト用）
	if bookID == "nonexistent-book" {
		http.Error(w, "Playlist not found", http.StatusNotFound)
		return
	}

	// モックレスポンス
	response := GetPlaylistResponse{
		PlaylistID: "mock-playlist-id",
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

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
