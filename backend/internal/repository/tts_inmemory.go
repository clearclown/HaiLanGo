package repository

import (
	"context"
	"database/sql"
	"fmt"
	"sync"
	"time"

	"github.com/clearclown/HaiLanGo/backend/internal/models"
	"github.com/google/uuid"
)

// TTSRepositoryInterface はTTSリポジトリのインターフェース
type TTSRepositoryInterface interface {
	// CreateJob はTTSジョブを作成
	CreateJob(ctx context.Context, userID, bookID uuid.UUID, pageNumber int, text string, language string, options models.TTSSynthesizeOptions) (*models.TTSJobDetail, error)

	// GetJob はTTSジョブを取得
	GetJob(ctx context.Context, jobID string) (*models.TTSJobDetail, error)

	// UpdateJobStatus はTTSジョブのステータスを更新
	UpdateJobStatus(ctx context.Context, jobID string, status models.TTSStatus, progress int) error

	// UpdateJobResult はTTSジョブの結果を更新
	UpdateJobResult(ctx context.Context, jobID string, audioID string, audioURL string) error

	// UpdateJobError はTTSジョブのエラーを更新
	UpdateJobError(ctx context.Context, jobID string, errorMsg string) error

	// GetJobsByBookID は書籍IDでTTSジョブを取得
	GetJobsByBookID(ctx context.Context, bookID uuid.UUID) ([]*models.TTSJobDetail, error)

	// GetJobsByUserID はユーザーIDでTTSジョブを取得
	GetJobsByUserID(ctx context.Context, userID uuid.UUID) ([]*models.TTSJobDetail, error)

	// GetCacheStats はキャッシュ統計情報を取得
	GetCacheStats(ctx context.Context, userID uuid.UUID) (*models.TTSCacheStats, error)

	// DeleteJob はTTSジョブを削除
	DeleteJob(ctx context.Context, jobID string) error

	// GetAudioURL は音声URLを取得（AudioIDから）
	GetAudioURL(ctx context.Context, audioID string) (string, error)

	// GetSupportedLanguages はサポート言語一覧を取得
	GetSupportedLanguages(ctx context.Context) ([]*models.TTSLanguage, error)
}

// InMemoryTTSRepository はインメモリTTSリポジトリ
type InMemoryTTSRepository struct {
	mu         sync.RWMutex
	jobs       map[string]*models.TTSJobDetail // jobID -> TTSJobDetail
	userJobs   map[string][]string             // userID -> []jobID
	bookJobs   map[string][]string             // bookID -> []jobID
	audioCache map[string]string               // audioID -> audioURL
	languages  []*models.TTSLanguage            // サポート言語
}

// NewInMemoryTTSRepository はインメモリTTSリポジトリを作成
func NewInMemoryTTSRepository() *InMemoryTTSRepository {
	repo := &InMemoryTTSRepository{
		jobs:       make(map[string]*models.TTSJobDetail),
		userJobs:   make(map[string][]string),
		bookJobs:   make(map[string][]string),
		audioCache: make(map[string]string),
		languages:  make([]*models.TTSLanguage, 0),
	}

	// サポート言語を初期化
	repo.initSupportedLanguages()

	// サンプルデータを初期化
	repo.initSampleData()

	return repo
}

func (r *InMemoryTTSRepository) initSupportedLanguages() {
	r.languages = []*models.TTSLanguage{
		{Code: "ja", Name: "Japanese", NativeName: "日本語", Voices: []string{"ja-JP-Neural2-B", "ja-JP-Neural2-C"}, IsSupported: true},
		{Code: "en", Name: "English", NativeName: "English", Voices: []string{"en-US-Neural2-A", "en-US-Neural2-C"}, IsSupported: true},
		{Code: "zh", Name: "Chinese", NativeName: "中文", Voices: []string{"zh-CN-Neural2-A", "zh-CN-Neural2-B"}, IsSupported: true},
		{Code: "ru", Name: "Russian", NativeName: "Русский", Voices: []string{"ru-RU-Wavenet-A", "ru-RU-Wavenet-B"}, IsSupported: true},
		{Code: "fa", Name: "Persian", NativeName: "فارسی", Voices: []string{"fa-IR-Wavenet-A"}, IsSupported: true},
		{Code: "he", Name: "Hebrew", NativeName: "עברית", Voices: []string{"he-IL-Wavenet-A"}, IsSupported: true},
		{Code: "es", Name: "Spanish", NativeName: "Español", Voices: []string{"es-ES-Neural2-A", "es-ES-Neural2-B"}, IsSupported: true},
		{Code: "fr", Name: "French", NativeName: "Français", Voices: []string{"fr-FR-Neural2-A", "fr-FR-Neural2-B"}, IsSupported: true},
		{Code: "pt", Name: "Portuguese", NativeName: "Português", Voices: []string{"pt-BR-Neural2-A", "pt-BR-Neural2-B"}, IsSupported: true},
		{Code: "de", Name: "German", NativeName: "Deutsch", Voices: []string{"de-DE-Neural2-A", "de-DE-Neural2-B"}, IsSupported: true},
		{Code: "it", Name: "Italian", NativeName: "Italiano", Voices: []string{"it-IT-Neural2-A", "it-IT-Neural2-B"}, IsSupported: true},
		{Code: "tr", Name: "Turkish", NativeName: "Türkçe", Voices: []string{"tr-TR-Wavenet-A"}, IsSupported: true},
	}
}

func (r *InMemoryTTSRepository) initSampleData() {
	testUserID := "550e8400-e29b-41d4-a716-446655440000"
	testBookID := "550e8400-e29b-41d4-a716-446655440000"

	// 完了済みのTTSジョブ（最初の30ページ）
	for i := 1; i <= 30; i++ {
		jobID := uuid.New().String()
		audioID := uuid.New().String()
		now := time.Now().Add(-time.Duration(31-i) * time.Hour)
		completedAt := now.Add(3 * time.Second)

		audioURL := fmt.Sprintf("/storage/audio/%s.mp3", audioID)
		r.audioCache[audioID] = audioURL

		job := &models.TTSJobDetail{
			ID:          jobID,
			BookID:      testBookID,
			PageNumber:  i,
			Text:        fmt.Sprintf("Здравствуйте! Как дела? (ページ %d)", i),
			Language:    "ru",
			Status:      models.TTSStatusCompleted,
			Progress:    100,
			AudioID:     audioID,
			AudioURL:    audioURL,
			CreatedAt:   now,
			UpdatedAt:   completedAt,
			CompletedAt: &completedAt,
		}

		r.jobs[jobID] = job
		r.userJobs[testUserID] = append(r.userJobs[testUserID], jobID)
		r.bookJobs[testBookID] = append(r.bookJobs[testBookID], jobID)
	}

	// 処理中のTTSジョブ（ページ31-33）
	for i := 31; i <= 33; i++ {
		jobID := uuid.New().String()
		now := time.Now().Add(-time.Duration(2) * time.Minute)

		job := &models.TTSJobDetail{
			ID:         jobID,
			BookID:     testBookID,
			PageNumber: i,
			Text:       fmt.Sprintf("こんにちは！元気ですか？ (ページ %d)", i),
			Language:   "ja",
			Status:     models.TTSStatusProcessing,
			Progress:   50 + (i-31)*10,
			CreatedAt:  now,
			UpdatedAt:  time.Now(),
		}

		r.jobs[jobID] = job
		r.userJobs[testUserID] = append(r.userJobs[testUserID], jobID)
		r.bookJobs[testBookID] = append(r.bookJobs[testBookID], jobID)
	}

	// ペンディングのTTSジョブ（ページ34-35）
	for i := 34; i <= 35; i++ {
		jobID := uuid.New().String()
		now := time.Now().Add(-time.Duration(30) * time.Second)

		job := &models.TTSJobDetail{
			ID:         jobID,
			BookID:     testBookID,
			PageNumber: i,
			Text:       fmt.Sprintf("Hello! How are you? (Page %d)", i),
			Language:   "en",
			Status:     models.TTSStatusPending,
			Progress:   0,
			CreatedAt:  now,
			UpdatedAt:  now,
		}

		r.jobs[jobID] = job
		r.userJobs[testUserID] = append(r.userJobs[testUserID], jobID)
		r.bookJobs[testBookID] = append(r.bookJobs[testBookID], jobID)
	}
}

func (r *InMemoryTTSRepository) CreateJob(ctx context.Context, userID, bookID uuid.UUID, pageNumber int, text string, language string, options models.TTSSynthesizeOptions) (*models.TTSJobDetail, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	jobID := uuid.New().String()
	now := time.Now()

	job := &models.TTSJobDetail{
		ID:         jobID,
		BookID:     bookID.String(),
		PageNumber: pageNumber,
		Text:       text,
		Language:   language,
		Status:     models.TTSStatusPending,
		Progress:   0,
		CreatedAt:  now,
		UpdatedAt:  now,
	}

	r.jobs[jobID] = job
	r.userJobs[userID.String()] = append(r.userJobs[userID.String()], jobID)
	r.bookJobs[bookID.String()] = append(r.bookJobs[bookID.String()], jobID)

	return job, nil
}

func (r *InMemoryTTSRepository) GetJob(ctx context.Context, jobID string) (*models.TTSJobDetail, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	job, exists := r.jobs[jobID]
	if !exists {
		return nil, fmt.Errorf("job not found: %s", jobID)
	}

	return job, nil
}

func (r *InMemoryTTSRepository) UpdateJobStatus(ctx context.Context, jobID string, status models.TTSStatus, progress int) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	job, exists := r.jobs[jobID]
	if !exists {
		return fmt.Errorf("job not found: %s", jobID)
	}

	job.Status = status
	job.Progress = progress
	job.UpdatedAt = time.Now()

	if status == models.TTSStatusCompleted || status == models.TTSStatusFailed {
		now := time.Now()
		job.CompletedAt = &now
	}

	return nil
}

func (r *InMemoryTTSRepository) UpdateJobResult(ctx context.Context, jobID string, audioID string, audioURL string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	job, exists := r.jobs[jobID]
	if !exists {
		return fmt.Errorf("job not found: %s", jobID)
	}

	job.AudioID = audioID
	job.AudioURL = audioURL
	job.Status = models.TTSStatusCompleted
	job.Progress = 100
	job.UpdatedAt = time.Now()
	now := time.Now()
	job.CompletedAt = &now

	// キャッシュに保存
	r.audioCache[audioID] = audioURL

	return nil
}

func (r *InMemoryTTSRepository) UpdateJobError(ctx context.Context, jobID string, errorMsg string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	job, exists := r.jobs[jobID]
	if !exists {
		return fmt.Errorf("job not found: %s", jobID)
	}

	job.Error = errorMsg
	job.Status = models.TTSStatusFailed
	job.UpdatedAt = time.Now()
	now := time.Now()
	job.CompletedAt = &now

	return nil
}

func (r *InMemoryTTSRepository) GetJobsByBookID(ctx context.Context, bookID uuid.UUID) ([]*models.TTSJobDetail, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	jobIDs, exists := r.bookJobs[bookID.String()]
	if !exists {
		return []*models.TTSJobDetail{}, nil
	}

	jobs := make([]*models.TTSJobDetail, 0, len(jobIDs))
	for _, jobID := range jobIDs {
		if job, exists := r.jobs[jobID]; exists {
			jobs = append(jobs, job)
		}
	}

	return jobs, nil
}

func (r *InMemoryTTSRepository) GetJobsByUserID(ctx context.Context, userID uuid.UUID) ([]*models.TTSJobDetail, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	jobIDs, exists := r.userJobs[userID.String()]
	if !exists {
		return []*models.TTSJobDetail{}, nil
	}

	jobs := make([]*models.TTSJobDetail, 0, len(jobIDs))
	for _, jobID := range jobIDs {
		if job, exists := r.jobs[jobID]; exists {
			jobs = append(jobs, job)
		}
	}

	return jobs, nil
}

func (r *InMemoryTTSRepository) GetCacheStats(ctx context.Context, userID uuid.UUID) (*models.TTSCacheStats, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	stats := &models.TTSCacheStats{
		TotalCached:  len(r.audioCache),
		CacheHitRate: 0.85, // サンプル値
		TotalSize:    int64(len(r.audioCache) * 1024 * 500), // 推定: 1音声 = 500KB
		AvgDuration:  15,                                     // 平均15秒
		Languages:    make(map[string]int),
	}

	// 言語ごとのカウント
	jobIDs, exists := r.userJobs[userID.String()]
	if exists {
		for _, jobID := range jobIDs {
			if job, exists := r.jobs[jobID]; exists && job.Status == models.TTSStatusCompleted {
				stats.Languages[job.Language]++
			}
		}
	}

	return stats, nil
}

func (r *InMemoryTTSRepository) DeleteJob(ctx context.Context, jobID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	_, exists := r.jobs[jobID]
	if !exists {
		return fmt.Errorf("job not found: %s", jobID)
	}

	delete(r.jobs, jobID)

	return nil
}

func (r *InMemoryTTSRepository) GetAudioURL(ctx context.Context, audioID string) (string, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	audioURL, exists := r.audioCache[audioID]
	if !exists {
		return "", fmt.Errorf("audio not found: %s", audioID)
	}

	return audioURL, nil
}

func (r *InMemoryTTSRepository) GetSupportedLanguages(ctx context.Context) ([]*models.TTSLanguage, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	return r.languages, nil
}

// SimulateTTSProcessing はTTS処理をシミュレート（テスト用）
func (r *InMemoryTTSRepository) SimulateTTSProcessing(ctx context.Context, jobID string) error {
	// ステータスを処理中に更新
	if err := r.UpdateJobStatus(ctx, jobID, models.TTSStatusProcessing, 10); err != nil {
		return err
	}

	// 処理進行をシミュレート
	for i := 20; i <= 90; i += 20 {
		time.Sleep(300 * time.Millisecond)
		if err := r.UpdateJobStatus(ctx, jobID, models.TTSStatusProcessing, i); err != nil {
			return err
		}
	}

	// 音声を生成
	audioID := uuid.New().String()
	audioURL := fmt.Sprintf("/storage/audio/%s.mp3", audioID)

	return r.UpdateJobResult(ctx, jobID, audioID, audioURL)
}

// PostgreSQL Implementation

type TTSRepositoryPostgres struct {
	db *sql.DB
}

func NewTTSRepositoryPostgres(db *sql.DB) TTSRepositoryInterface {
	return &TTSRepositoryPostgres{db: db}
}

func (r *TTSRepositoryPostgres) CreateJob(ctx context.Context, userID, bookID uuid.UUID, pageNumber int, text string, language string, options models.TTSSynthesizeOptions) (*models.TTSJobDetail, error) {
	jobID := uuid.New()
	now := time.Now()

	_, err := r.db.ExecContext(ctx, `
		INSERT INTO tts_jobs (id, user_id, book_id, status, progress, total_items, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`, jobID, userID, bookID, "pending", 0, 1, now, now)

	if err != nil {
		return nil, err
	}

	return &models.TTSJobDetail{
		ID:         jobID.String(),
		BookID:     bookID.String(),
		PageNumber: pageNumber,
		Text:       text,
		Language:   language,
		Status:     models.TTSStatusPending,
		Progress:   0,
		CreatedAt:  now,
		UpdatedAt:  now,
	}, nil
}

func (r *TTSRepositoryPostgres) GetJob(ctx context.Context, jobID string) (*models.TTSJobDetail, error) {
	var job models.TTSJobDetail
	var status string
	var errorMsg sql.NullString
	var completedAt sql.NullTime

	err := r.db.QueryRowContext(ctx, `
		SELECT id, book_id, status, progress, error_message, created_at, updated_at, completed_at
		FROM tts_jobs
		WHERE id = $1
	`, jobID).Scan(&job.ID, &job.BookID, &status, &job.Progress, &errorMsg, &job.CreatedAt, &job.UpdatedAt, &completedAt)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("job not found: %s", jobID)
	}
	if err != nil {
		return nil, err
	}

	job.Status = models.TTSStatus(status)
	if errorMsg.Valid {
		job.Error = errorMsg.String
	}
	if completedAt.Valid {
		job.CompletedAt = &completedAt.Time
	}

	return &job, nil
}

func (r *TTSRepositoryPostgres) UpdateJobStatus(ctx context.Context, jobID string, status models.TTSStatus, progress int) error {
	query := `UPDATE tts_jobs SET status = $1, progress = $2, updated_at = $3 WHERE id = $4`
	args := []interface{}{string(status), progress, time.Now(), jobID}

	if status == models.TTSStatusCompleted || status == models.TTSStatusFailed {
		query = `UPDATE tts_jobs SET status = $1, progress = $2, updated_at = $3, completed_at = $3 WHERE id = $4`
	}

	_, err := r.db.ExecContext(ctx, query, args...)
	return err
}

func (r *TTSRepositoryPostgres) UpdateJobResult(ctx context.Context, jobID string, audioID string, audioURL string) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Update job status
	_, err = tx.ExecContext(ctx, `
		UPDATE tts_jobs SET status = $1, progress = $2, updated_at = $3, completed_at = $3
		WHERE id = $4
	`, string(models.TTSStatusCompleted), 100, time.Now(), jobID)
	if err != nil {
		return err
	}

	// Insert audio cache
	_, err = tx.ExecContext(ctx, `
		INSERT INTO tts_audio (id, text, language, audio_url, created_at, last_accessed_at, access_count)
		VALUES ($1, $2, $3, $4, NOW(), NOW(), 1)
		ON CONFLICT (id) DO UPDATE SET
			last_accessed_at = NOW(),
			access_count = tts_audio.access_count + 1
	`, audioID, "", "", audioURL)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (r *TTSRepositoryPostgres) UpdateJobError(ctx context.Context, jobID string, errorMsg string) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE tts_jobs SET status = $1, error_message = $2, updated_at = $3, completed_at = $3
		WHERE id = $4
	`, string(models.TTSStatusFailed), errorMsg, time.Now(), jobID)
	return err
}

func (r *TTSRepositoryPostgres) GetJobsByBookID(ctx context.Context, bookID uuid.UUID) ([]*models.TTSJobDetail, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, book_id, status, progress, error_message, created_at, updated_at, completed_at
		FROM tts_jobs
		WHERE book_id = $1
		ORDER BY created_at DESC
	`, bookID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	jobs := []*models.TTSJobDetail{}
	for rows.Next() {
		var job models.TTSJobDetail
		var status string
		var errorMsg sql.NullString
		var completedAt sql.NullTime

		err := rows.Scan(&job.ID, &job.BookID, &status, &job.Progress, &errorMsg, &job.CreatedAt, &job.UpdatedAt, &completedAt)
		if err != nil {
			return nil, err
		}

		job.Status = models.TTSStatus(status)
		if errorMsg.Valid {
			job.Error = errorMsg.String
		}
		if completedAt.Valid {
			job.CompletedAt = &completedAt.Time
		}

		jobs = append(jobs, &job)
	}

	return jobs, nil
}

func (r *TTSRepositoryPostgres) GetJobsByUserID(ctx context.Context, userID uuid.UUID) ([]*models.TTSJobDetail, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, book_id, status, progress, error_message, created_at, updated_at, completed_at
		FROM tts_jobs
		WHERE user_id = $1
		ORDER BY created_at DESC
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	jobs := []*models.TTSJobDetail{}
	for rows.Next() {
		var job models.TTSJobDetail
		var status string
		var errorMsg sql.NullString
		var completedAt sql.NullTime

		err := rows.Scan(&job.ID, &job.BookID, &status, &job.Progress, &errorMsg, &job.CreatedAt, &job.UpdatedAt, &completedAt)
		if err != nil {
			return nil, err
		}

		job.Status = models.TTSStatus(status)
		if errorMsg.Valid {
			job.Error = errorMsg.String
		}
		if completedAt.Valid {
			job.CompletedAt = &completedAt.Time
		}

		jobs = append(jobs, &job)
	}

	return jobs, nil
}

func (r *TTSRepositoryPostgres) GetCacheStats(ctx context.Context, userID uuid.UUID) (*models.TTSCacheStats, error) {
	stats := &models.TTSCacheStats{
		Languages: make(map[string]int),
	}

	// Get total cached audio count
	r.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM tts_audio").Scan(&stats.TotalCached)

	// Get total size (estimated)
	stats.TotalSize = int64(stats.TotalCached * 1024 * 500) // 500KB per audio
	stats.CacheHitRate = 0.85                                // Default
	stats.AvgDuration = 15                                   // Default

	return stats, nil
}

func (r *TTSRepositoryPostgres) DeleteJob(ctx context.Context, jobID string) error {
	_, err := r.db.ExecContext(ctx, "DELETE FROM tts_jobs WHERE id = $1", jobID)
	return err
}

func (r *TTSRepositoryPostgres) GetAudioURL(ctx context.Context, audioID string) (string, error) {
	var audioURL string
	err := r.db.QueryRowContext(ctx, `
		SELECT audio_url FROM tts_audio WHERE id = $1
	`, audioID).Scan(&audioURL)

	if err == sql.ErrNoRows {
		return "", fmt.Errorf("audio not found: %s", audioID)
	}
	if err != nil {
		return "", err
	}

	// Update access tracking
	r.db.ExecContext(ctx, `
		UPDATE tts_audio SET last_accessed_at = NOW(), access_count = access_count + 1
		WHERE id = $1
	`, audioID)

	return audioURL, nil
}

func (r *TTSRepositoryPostgres) GetSupportedLanguages(ctx context.Context) ([]*models.TTSLanguage, error) {
	// Return hardcoded supported languages (could be from DB in future)
	languages := []*models.TTSLanguage{
		{Code: "ja", Name: "Japanese", NativeName: "日本語", Voices: []string{"ja-JP-Neural2-B"}, IsSupported: true},
		{Code: "en", Name: "English", NativeName: "English", Voices: []string{"en-US-Neural2-A"}, IsSupported: true},
		{Code: "zh", Name: "Chinese", NativeName: "中文", Voices: []string{"zh-CN-Neural2-A"}, IsSupported: true},
		{Code: "ru", Name: "Russian", NativeName: "Русский", Voices: []string{"ru-RU-Wavenet-A"}, IsSupported: true},
		{Code: "fa", Name: "Persian", NativeName: "فارسی", Voices: []string{"fa-IR-Wavenet-A"}, IsSupported: true},
		{Code: "he", Name: "Hebrew", NativeName: "עברית", Voices: []string{"he-IL-Wavenet-A"}, IsSupported: true},
		{Code: "es", Name: "Spanish", NativeName: "Español", Voices: []string{"es-ES-Neural2-A"}, IsSupported: true},
		{Code: "fr", Name: "French", NativeName: "Français", Voices: []string{"fr-FR-Neural2-A"}, IsSupported: true},
		{Code: "pt", Name: "Portuguese", NativeName: "Português", Voices: []string{"pt-BR-Neural2-A"}, IsSupported: true},
		{Code: "de", Name: "German", NativeName: "Deutsch", Voices: []string{"de-DE-Neural2-A"}, IsSupported: true},
		{Code: "it", Name: "Italian", NativeName: "Italiano", Voices: []string{"it-IT-Neural2-A"}, IsSupported: true},
		{Code: "tr", Name: "Turkish", NativeName: "Türkçe", Voices: []string{"tr-TR-Wavenet-A"}, IsSupported: true},
	}
	return languages, nil
}
