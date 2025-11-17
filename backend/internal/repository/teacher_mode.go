package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"time"

	"github.com/clearclown/HaiLanGo/backend/internal/models"
	"github.com/google/uuid"
)

// TeacherModeRepository は教師モードのリポジトリインターフェース
type TeacherModeRepository interface {
	// SaveDownload はダウンロード履歴を保存する
	SaveDownload(ctx context.Context, download *models.TeacherModeDownload) error

	// GetDownloadByID はIDでダウンロード履歴を取得する
	GetDownloadByID(ctx context.Context, id uuid.UUID) (*models.TeacherModeDownload, error)

	// GetDownloadsByUserID はユーザーIDでダウンロード履歴を取得する
	GetDownloadsByUserID(ctx context.Context, userID uuid.UUID) ([]*models.TeacherModeDownload, error)

	// SavePlaybackState は再生状態を保存する
	SavePlaybackState(ctx context.Context, state *models.TeacherModePlaybackHistory) error

	// GetPlaybackState は再生状態を取得する
	GetPlaybackState(ctx context.Context, userID uuid.UUID, bookID uuid.UUID) (*models.TeacherModePlaybackHistory, error)

	// UpdatePlaybackState は再生状態を更新する
	UpdatePlaybackState(ctx context.Context, userID uuid.UUID, bookID uuid.UUID, currentPage int, currentSegmentIndex int, elapsedTime int) error
}

// teacherModeRepositoryPostgres はPostgreSQLベースの教師モードリポジトリ実装
type teacherModeRepositoryPostgres struct {
	db *sql.DB
}

// NewTeacherModeRepositoryPostgres は新しいPostgreSQL実装のTeacherModeRepositoryを作成する
func NewTeacherModeRepositoryPostgres(db *sql.DB) TeacherModeRepository {
	return &teacherModeRepositoryPostgres{db: db}
}

// SaveDownload はダウンロード履歴を保存する
func (r *teacherModeRepositoryPostgres) SaveDownload(ctx context.Context, download *models.TeacherModeDownload) error {
	// SettingsをJSONBに変換
	settingsJSON, err := json.Marshal(download.Settings)
	if err != nil {
		return err
	}

	query := `
		INSERT INTO teacher_mode_downloads (
			id, user_id, book_id, settings, total_size_bytes,
			downloaded_at, expires_at, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`

	_, err = r.db.ExecContext(
		ctx,
		query,
		download.ID,
		download.UserID,
		download.BookID,
		settingsJSON,
		download.TotalSizeBytes,
		download.DownloadedAt,
		download.ExpiresAt,
		download.CreatedAt,
		download.UpdatedAt,
	)

	return err
}

// GetDownloadByID はIDでダウンロード履歴を取得する
func (r *teacherModeRepositoryPostgres) GetDownloadByID(ctx context.Context, id uuid.UUID) (*models.TeacherModeDownload, error) {
	query := `
		SELECT id, user_id, book_id, settings, total_size_bytes,
		       downloaded_at, expires_at, created_at, updated_at
		FROM teacher_mode_downloads
		WHERE id = $1
	`

	download := &models.TeacherModeDownload{}
	var settingsJSON []byte

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&download.ID,
		&download.UserID,
		&download.BookID,
		&settingsJSON,
		&download.TotalSizeBytes,
		&download.DownloadedAt,
		&download.ExpiresAt,
		&download.CreatedAt,
		&download.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	// JSONBをSettingsに変換
	if err := json.Unmarshal(settingsJSON, &download.Settings); err != nil {
		return nil, err
	}

	return download, nil
}

// GetDownloadsByUserID はユーザーIDでダウンロード履歴を取得する
func (r *teacherModeRepositoryPostgres) GetDownloadsByUserID(ctx context.Context, userID uuid.UUID) ([]*models.TeacherModeDownload, error) {
	query := `
		SELECT id, user_id, book_id, settings, total_size_bytes,
		       downloaded_at, expires_at, created_at, updated_at
		FROM teacher_mode_downloads
		WHERE user_id = $1
		ORDER BY downloaded_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var downloads []*models.TeacherModeDownload
	for rows.Next() {
		download := &models.TeacherModeDownload{}
		var settingsJSON []byte

		err := rows.Scan(
			&download.ID,
			&download.UserID,
			&download.BookID,
			&settingsJSON,
			&download.TotalSizeBytes,
			&download.DownloadedAt,
			&download.ExpiresAt,
			&download.CreatedAt,
			&download.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		// JSONBをSettingsに変換
		if err := json.Unmarshal(settingsJSON, &download.Settings); err != nil {
			return nil, err
		}

		downloads = append(downloads, download)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return downloads, nil
}

// SavePlaybackState は再生状態を保存する
func (r *teacherModeRepositoryPostgres) SavePlaybackState(ctx context.Context, state *models.TeacherModePlaybackHistory) error {
	query := `
		INSERT INTO teacher_mode_playback_history (
			id, user_id, book_id, current_page, current_segment_index,
			elapsed_time, total_play_time_seconds, last_played_at, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		ON CONFLICT (user_id, book_id) DO UPDATE SET
			current_page = EXCLUDED.current_page,
			current_segment_index = EXCLUDED.current_segment_index,
			elapsed_time = EXCLUDED.elapsed_time,
			total_play_time_seconds = EXCLUDED.total_play_time_seconds,
			last_played_at = EXCLUDED.last_played_at,
			updated_at = EXCLUDED.updated_at
	`

	_, err := r.db.ExecContext(
		ctx,
		query,
		state.ID,
		state.UserID,
		state.BookID,
		state.CurrentPage,
		state.CurrentSegmentIndex,
		state.ElapsedTime,
		state.TotalPlayTimeSeconds,
		state.LastPlayedAt,
		state.CreatedAt,
		state.UpdatedAt,
	)

	return err
}

// GetPlaybackState は再生状態を取得する
func (r *teacherModeRepositoryPostgres) GetPlaybackState(ctx context.Context, userID uuid.UUID, bookID uuid.UUID) (*models.TeacherModePlaybackHistory, error) {
	query := `
		SELECT id, user_id, book_id, current_page, current_segment_index,
		       elapsed_time, total_play_time_seconds, last_played_at, created_at, updated_at
		FROM teacher_mode_playback_history
		WHERE user_id = $1 AND book_id = $2
	`

	state := &models.TeacherModePlaybackHistory{}
	err := r.db.QueryRowContext(ctx, query, userID, bookID).Scan(
		&state.ID,
		&state.UserID,
		&state.BookID,
		&state.CurrentPage,
		&state.CurrentSegmentIndex,
		&state.ElapsedTime,
		&state.TotalPlayTimeSeconds,
		&state.LastPlayedAt,
		&state.CreatedAt,
		&state.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return state, nil
}

// UpdatePlaybackState は再生状態を更新する
func (r *teacherModeRepositoryPostgres) UpdatePlaybackState(ctx context.Context, userID uuid.UUID, bookID uuid.UUID, currentPage int, currentSegmentIndex int, elapsedTime int) error {
	query := `
		UPDATE teacher_mode_playback_history
		SET current_page = $1,
		    current_segment_index = $2,
		    elapsed_time = $3,
		    last_played_at = NOW(),
		    updated_at = NOW()
		WHERE user_id = $4 AND book_id = $5
	`

	result, err := r.db.ExecContext(
		ctx,
		query,
		currentPage,
		currentSegmentIndex,
		elapsedTime,
		userID,
		bookID,
	)

	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	// レコードが存在しない場合は新規作成
	if rowsAffected == 0 {
		state := &models.TeacherModePlaybackHistory{
			ID:                   uuid.New(),
			UserID:               userID,
			BookID:               bookID,
			CurrentPage:          currentPage,
			CurrentSegmentIndex:  currentSegmentIndex,
			ElapsedTime:          elapsedTime,
			TotalPlayTimeSeconds: 0,
			LastPlayedAt:         time.Now(),
			CreatedAt:            time.Now(),
			UpdatedAt:            time.Now(),
		}
		return r.SavePlaybackState(ctx, state)
	}

	return nil
}
