package service

import (
	"context"
	"fmt"
	"time"

	"github.com/clearclown/HaiLanGo/backend/internal/models"
	"github.com/clearclown/HaiLanGo/backend/internal/repository"
	"github.com/google/uuid"
)

// TeacherModeService は教師モードのサービス
type TeacherModeService struct {
	teacherModeRepo repository.TeacherModeRepository
	pageRepo        repository.PageRepository
	bookRepo        repository.BookRepository
	ttsRepo         repository.TTSRepositoryInterface
}

// NewTeacherModeService は新しいTeacherModeServiceを作成する
func NewTeacherModeService(
	teacherModeRepo repository.TeacherModeRepository,
	pageRepo repository.PageRepository,
	bookRepo repository.BookRepository,
	ttsRepo repository.TTSRepositoryInterface,
) *TeacherModeService {
	return &TeacherModeService{
		teacherModeRepo: teacherModeRepo,
		pageRepo:        pageRepo,
		bookRepo:        bookRepo,
		ttsRepo:         ttsRepo,
	}
}

// GeneratePlaylist は教師モードのプレイリストを生成する
func (s *TeacherModeService) GeneratePlaylist(
	ctx context.Context,
	userID uuid.UUID,
	bookID uuid.UUID,
	settings *models.TeacherModeSettings,
	pageRange *models.PageRange,
) (*models.TeacherModePlaylist, error) {
	// 書籍情報を取得
	book, err := s.bookRepo.GetByID(ctx, bookID)
	if err != nil {
		return nil, fmt.Errorf("failed to get book: %w", err)
	}
	if book == nil {
		return nil, fmt.Errorf("book not found")
	}

	// ページ情報を取得
	pages, err := s.pageRepo.FindByBookID(ctx, bookID)
	if err != nil {
		return nil, fmt.Errorf("failed to get pages: %w", err)
	}

	// ページ範囲をフィルタリング
	var filteredPages []*models.Page
	for _, page := range pages {
		if pageRange != nil {
			if page.PageNumber >= pageRange.Start && page.PageNumber <= pageRange.End {
				filteredPages = append(filteredPages, page)
			}
		} else {
			filteredPages = append(filteredPages, page)
		}
	}

	// プレイリストを作成
	playlist := &models.TeacherModePlaylist{
		ID:       uuid.New().String(),
		BookID:   bookID,
		Pages:    make([]models.PageAudio, 0),
		Settings: *settings,
	}

	// TTS options for speed
	ttsOptions := models.TTSSynthesizeOptions{
		Speed: settings.Speed,
	}
	if settings.AudioQuality == "premium" {
		ttsOptions.Voice = "premium" // Use premium voice
	}

	totalDuration := 0

	// 各ページの音声セグメントを生成
	for _, page := range filteredPages {
		pageAudio := models.PageAudio{
			PageNumber: page.PageNumber,
			Segments:   make([]models.AudioSegment, 0),
		}

		segmentID := 0

		// 1. 学習先言語のフレーズ（必須）
		if page.OCRText != "" {
			phraseSegment, duration, err := s.createAudioSegment(
				ctx,
				userID,
				bookID,
				page.PageNumber,
				segmentID,
				models.AudioSegmentTypePhrase,
				page.OCRText,
				book.TargetLanguage,
				ttsOptions,
			)
			if err != nil {
				return nil, fmt.Errorf("failed to create phrase segment for page %d: %w", page.PageNumber, err)
			}
			pageAudio.Segments = append(pageAudio.Segments, *phraseSegment)
			pageAudio.TotalDuration += duration
			totalDuration += duration
			segmentID++
		}

		// 2. 母国語訳（オプション）
		if settings.Content.IncludeTranslation && page.OCRText != "" {
			// TODO: 実際には翻訳APIを使用する
			translationText := fmt.Sprintf("Translation of: %s", page.OCRText)
			translationSegment, duration, err := s.createAudioSegment(
				ctx,
				userID,
				bookID,
				page.PageNumber,
				segmentID,
				models.AudioSegmentTypeTranslation,
				translationText,
				book.NativeLanguage,
				ttsOptions,
			)
			if err != nil {
				return nil, fmt.Errorf("failed to create translation segment for page %d: %w", page.PageNumber, err)
			}
			pageAudio.Segments = append(pageAudio.Segments, *translationSegment)
			pageAudio.TotalDuration += duration
			totalDuration += duration
			segmentID++
		}

		// 3. 単語解説（オプション）
		if settings.Content.IncludeWordExplanation && page.OCRText != "" {
			// TODO: 実際には辞書APIを使用する
			explanationText := fmt.Sprintf("Word explanation for: %s", page.OCRText)
			explanationSegment, duration, err := s.createAudioSegment(
				ctx,
				userID,
				bookID,
				page.PageNumber,
				segmentID,
				models.AudioSegmentTypeExplanation,
				explanationText,
				book.NativeLanguage,
				ttsOptions,
			)
			if err != nil {
				return nil, fmt.Errorf("failed to create explanation segment for page %d: %w", page.PageNumber, err)
			}
			pageAudio.Segments = append(pageAudio.Segments, *explanationSegment)
			pageAudio.TotalDuration += duration
			totalDuration += duration
			segmentID++
		}

		// 4. ページ間隔（一時停止）
		if settings.PageInterval > 0 {
			pauseSegment := &models.AudioSegment{
				ID:       fmt.Sprintf("page-%d-segment-%d", page.PageNumber, segmentID),
				Type:     models.AudioSegmentTypePause,
				AudioURL: "",
				Duration: settings.PageInterval * 1000, // 秒をミリ秒に変換
				Text:     "",
				Language: "",
			}
			pageAudio.Segments = append(pageAudio.Segments, *pauseSegment)
			pageAudio.TotalDuration += pauseSegment.Duration
			totalDuration += pauseSegment.Duration
		}

		playlist.Pages = append(playlist.Pages, pageAudio)
	}

	playlist.TotalDuration = totalDuration

	return playlist, nil
}

// createAudioSegment は音声セグメントを作成する
func (s *TeacherModeService) createAudioSegment(
	ctx context.Context,
	userID uuid.UUID,
	bookID uuid.UUID,
	pageNumber int,
	segmentID int,
	segmentType models.AudioSegmentType,
	text string,
	language string,
	options models.TTSSynthesizeOptions,
) (*models.AudioSegment, int, error) {
	// TTS APIを使用して音声を生成
	job, err := s.ttsRepo.CreateJob(ctx, userID, bookID, pageNumber, text, language, options)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to create TTS job: %w", err)
	}

	// 仮の音声URLを生成（実際にはTTS APIから取得）
	audioURL := job.AudioURL
	if audioURL == "" {
		audioURL = fmt.Sprintf("/api/v1/tts/audio/%s", job.AudioID)
	}

	// 仮の長さを計算（実際にはTTS APIから取得）
	// 1文字あたり約100ミリ秒として計算
	duration := len(text) * 100

	segment := &models.AudioSegment{
		ID:       fmt.Sprintf("page-%d-segment-%d", pageNumber, segmentID),
		Type:     segmentType,
		AudioURL: audioURL,
		Duration: duration,
		Text:     text,
		Language: language,
	}

	return segment, duration, nil
}

// GenerateDownloadPackage は教師モードのダウンロードパッケージを生成する
func (s *TeacherModeService) GenerateDownloadPackage(
	ctx context.Context,
	userID uuid.UUID,
	bookID uuid.UUID,
	settings *models.TeacherModeSettings,
) (packageID uuid.UUID, downloadURL string, totalSize int64, expiresAt time.Time, error error) {
	// プレイリストを生成
	playlist, err := s.GeneratePlaylist(ctx, userID, bookID, settings, nil)
	if err != nil {
		return uuid.Nil, "", 0, time.Time{}, fmt.Errorf("failed to generate playlist: %w", err)
	}

	// 総サイズを計算（仮）
	// 実際にはすべての音声ファイルのサイズを合計する
	totalSize = int64(len(playlist.Pages) * 1024 * 1024) // 1ページあたり1MB

	// 有効期限を設定（7日後）
	expiresAt = time.Now().Add(7 * 24 * time.Hour)

	// ダウンロード履歴を保存
	download := &models.TeacherModeDownload{
		ID:             uuid.New(),
		UserID:         userID,
		BookID:         bookID,
		Settings:       *settings,
		TotalSizeBytes: totalSize,
		DownloadedAt:   time.Now(),
		ExpiresAt:      &expiresAt,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	if err := s.teacherModeRepo.SaveDownload(ctx, download); err != nil {
		return uuid.Nil, "", 0, time.Time{}, fmt.Errorf("failed to save download: %w", err)
	}

	// ダウンロードURLを生成（仮）
	// 実際にはZIPファイルを生成してストレージにアップロードする
	downloadURL = fmt.Sprintf("/api/v1/teacher-mode/download/%s", download.ID.String())

	return download.ID, downloadURL, totalSize, expiresAt, nil
}

// UpdatePlaybackState は再生状態を更新する
func (s *TeacherModeService) UpdatePlaybackState(
	ctx context.Context,
	userID uuid.UUID,
	bookID uuid.UUID,
	state *models.PlaybackState,
) error {
	// 経過時間をミリ秒から秒に変換
	elapsedTimeSeconds := state.ElapsedTime / 1000

	return s.teacherModeRepo.UpdatePlaybackState(
		ctx,
		userID,
		bookID,
		state.CurrentPage,
		state.CurrentSegmentIndex,
		elapsedTimeSeconds,
	)
}

// GetPlaybackState は再生状態を取得する
func (s *TeacherModeService) GetPlaybackState(
	ctx context.Context,
	userID uuid.UUID,
	bookID uuid.UUID,
) (*models.PlaybackState, error) {
	history, err := s.teacherModeRepo.GetPlaybackState(ctx, userID, bookID)
	if err != nil {
		return nil, fmt.Errorf("failed to get playback state: %w", err)
	}

	// 履歴が存在しない場合はデフォルト状態を返す
	if history == nil {
		return &models.PlaybackState{
			Status:              models.PlaybackStatusStopped,
			CurrentPage:         1,
			CurrentSegmentIndex: 0,
			ElapsedTime:         0,
			TotalDuration:       0,
		}, nil
	}

	// 履歴を PlaybackState に変換
	state := &models.PlaybackState{
		Status:              models.PlaybackStatusStopped,
		CurrentPage:         history.CurrentPage,
		CurrentSegmentIndex: history.CurrentSegmentIndex,
		ElapsedTime:         history.ElapsedTime * 1000, // 秒をミリ秒に変換
		TotalDuration:       0,                          // プレイリストから計算する必要がある
	}

	return state, nil
}
