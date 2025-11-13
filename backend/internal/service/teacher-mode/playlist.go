// Package teachermode provides teacher mode functionality
package teachermode

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
)

// GeneratePlaylist プレイリストを生成
func (s *Service) GeneratePlaylist(ctx context.Context, bookID string, settings *TeacherModeSettings) (*TeacherModePlaylist, error) {
	// 入力検証
	if bookID == "" {
		return nil, errors.New("invalid book ID")
	}

	if settings == nil {
		return nil, errors.New("settings is required")
	}

	// 設定の検証
	if err := s.ValidateSettings(settings); err != nil {
		return nil, err
	}

	// プレイリストID生成
	playlistID := uuid.New().String()

	// ページデータを取得（実際はデータベースから取得）
	pages, err := s.fetchBookPages(ctx, bookID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch book pages: %w", err)
	}

	// 各ページの音声セグメントを生成
	var pageAudios []PageAudio
	var totalDuration int64

	for _, page := range pages {
		segments, err := s.GenerateAudioSegments(ctx, &page, settings)
		if err != nil {
			return nil, fmt.Errorf("failed to generate audio segments for page %d: %w", page.PageNumber, err)
		}

		pageDuration := s.CalculateTotalDuration(segments)

		pageAudio := PageAudio{
			PageNumber:    page.PageNumber,
			Segments:      segments,
			TotalDuration: pageDuration,
		}

		pageAudios = append(pageAudios, pageAudio)
		totalDuration += pageDuration

		// ページ間隔を追加
		if page.PageNumber < len(pages) {
			totalDuration += int64(settings.PageInterval) * 1000
		}
	}

	playlist := &TeacherModePlaylist{
		ID:            playlistID,
		BookID:        bookID,
		Pages:         pageAudios,
		Settings:      settings,
		TotalDuration: totalDuration,
	}

	return playlist, nil
}

// fetchBookPages 書籍のページデータを取得（モック実装）
func (s *Service) fetchBookPages(ctx context.Context, bookID string) ([]PageData, error) {
	// 実際はデータベースから取得
	// 現在はテスト用のモックデータを返す
	mockPages := []PageData{
		{
			PageNumber:      1,
			PhraseText:      "Hello",
			TranslationText: "こんにちは",
			Language:        "en",
		},
		{
			PageNumber:      2,
			PhraseText:      "World",
			TranslationText: "世界",
			Language:        "en",
		},
	}

	return mockPages, nil
}

// GenerateDownloadPackage オフラインダウンロードパッケージを生成
func (s *Service) GenerateDownloadPackage(ctx context.Context, bookID string, settings *TeacherModeSettings) (*DownloadPackage, error) {
	// プレイリストを生成
	playlist, err := s.GeneratePlaylist(ctx, bookID, settings)
	if err != nil {
		return nil, err
	}

	// パッケージID生成
	packageID := uuid.New().String()

	// 音声ファイルをZIP化（実際の実装）
	// 現在はモックデータを返す
	downloadPackage := &DownloadPackage{
		PackageID:   packageID,
		DownloadURL: fmt.Sprintf("http://example.com/packages/%s.zip", packageID),
		TotalSize:   248000000, // 約248MB
		ExpiresAt:   "2025-12-01T00:00:00Z",
	}

	return downloadPackage, nil
}

// DownloadPackage ダウンロードパッケージ
type DownloadPackage struct {
	PackageID   string `json:"packageId"`
	DownloadURL string `json:"downloadUrl"`
	TotalSize   int64  `json:"totalSize"`
	ExpiresAt   string `json:"expiresAt"`
}
