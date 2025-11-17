package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/clearclown/HaiLanGo/backend/internal/models"
	"github.com/google/uuid"
)

// STTRepositoryInterface はSTTリポジトリのインターフェース
type STTRepositoryInterface interface {
	// CreateJob はSTTジョブを作成
	CreateJob(ctx context.Context, userID uuid.UUID, bookID uuid.UUID, pageNumber int, audioURL string, language string, referenceText string, options models.STTRecognizeOptions) (*models.STTJobDetail, error)

	// GetJob はSTTジョブを取得
	GetJob(ctx context.Context, jobID string) (*models.STTJobDetail, error)

	// UpdateJobStatus はSTTジョブのステータスを更新
	UpdateJobStatus(ctx context.Context, jobID string, status models.STTStatus, progress int) error

	// UpdateJobResult はSTTジョブの結果を更新
	UpdateJobResult(ctx context.Context, jobID string, result *models.STTResult, score *models.PronunciationScore) error

	// UpdateJobError はSTTジョブのエラーを更新
	UpdateJobError(ctx context.Context, jobID string, errorMsg string) error

	// GetJobsByBookID は書籍IDでSTTジョブを取得
	GetJobsByBookID(ctx context.Context, bookID uuid.UUID) ([]*models.STTJobDetail, error)

	// GetJobsByUserID はユーザーIDでSTTジョブを取得
	GetJobsByUserID(ctx context.Context, userID uuid.UUID) ([]*models.STTJobDetail, error)

	// GetStatistics は統計情報を取得
	GetStatistics(ctx context.Context, userID uuid.UUID) (*models.STTStatistics, error)

	// DeleteJob はSTTジョブを削除
	DeleteJob(ctx context.Context, jobID string) error

	// GetSupportedLanguages はサポート言語一覧を取得
	GetSupportedLanguages(ctx context.Context) ([]*models.STTLanguage, error)
}

// InMemorySTTRepository はインメモリSTTリポジトリ
type InMemorySTTRepository struct {
	mu        sync.RWMutex
	jobs      map[string]*models.STTJobDetail // jobID -> STTJobDetail
	userJobs  map[string][]string             // userID -> []jobID
	bookJobs  map[string][]string             // bookID -> []jobID
	languages []*models.STTLanguage            // サポート言語
}

// NewInMemorySTTRepository はインメモリSTTリポジトリを作成
func NewInMemorySTTRepository() *InMemorySTTRepository {
	repo := &InMemorySTTRepository{
		jobs:      make(map[string]*models.STTJobDetail),
		userJobs:  make(map[string][]string),
		bookJobs:  make(map[string][]string),
		languages: make([]*models.STTLanguage, 0),
	}

	// サポート言語を初期化
	repo.initSupportedLanguages()

	// サンプルデータを初期化
	repo.initSampleData()

	return repo
}

func (r *InMemorySTTRepository) initSupportedLanguages() {
	r.languages = []*models.STTLanguage{
		{Code: "ja", Name: "Japanese", NativeName: "日本語", IsSupported: true, SupportsPronunciation: true},
		{Code: "en", Name: "English", NativeName: "English", IsSupported: true, SupportsPronunciation: true},
		{Code: "zh", Name: "Chinese", NativeName: "中文", IsSupported: true, SupportsPronunciation: true},
		{Code: "ru", Name: "Russian", NativeName: "Русский", IsSupported: true, SupportsPronunciation: true},
		{Code: "fa", Name: "Persian", NativeName: "فارسی", IsSupported: true, SupportsPronunciation: true},
		{Code: "he", Name: "Hebrew", NativeName: "עברית", IsSupported: true, SupportsPronunciation: true},
		{Code: "es", Name: "Spanish", NativeName: "Español", IsSupported: true, SupportsPronunciation: true},
		{Code: "fr", Name: "French", NativeName: "Français", IsSupported: true, SupportsPronunciation: true},
		{Code: "pt", Name: "Portuguese", NativeName: "Português", IsSupported: true, SupportsPronunciation: true},
		{Code: "de", Name: "German", NativeName: "Deutsch", IsSupported: true, SupportsPronunciation: true},
		{Code: "it", Name: "Italian", NativeName: "Italiano", IsSupported: true, SupportsPronunciation: true},
		{Code: "tr", Name: "Turkish", NativeName: "Türkçe", IsSupported: true, SupportsPronunciation: true},
	}
}

func (r *InMemorySTTRepository) initSampleData() {
	testUserID := "550e8400-e29b-41d4-a716-446655440000"
	testBookID := "550e8400-e29b-41d4-a716-446655440000"

	// 完了済みのSTTジョブ（最初の25個）
	for i := 1; i <= 25; i++ {
		jobID := uuid.New().String()
		now := time.Now().Add(-time.Duration(26-i) * time.Hour)
		completedAt := now.Add(5 * time.Second)

		score := 70 + i // スコア70-95
		job := &models.STTJobDetail{
			ID:            jobID,
			BookID:        testBookID,
			PageNumber:    i,
			AudioURL:      fmt.Sprintf("/storage/audio/recording_%d.wav", i),
			ReferenceText: fmt.Sprintf("Здравствуйте! Как дела? (ページ %d)", i),
			Language:      "ru",
			Status:        models.STTStatusCompleted,
			Progress:      100,
			Result: &models.STTResult{
				Text:       fmt.Sprintf("Здравствуйте! Как дела? (ページ %d)", i),
				Language:   "ru",
				Confidence: 0.90 + float64(i)*0.001,
				Duration:   3.5,
				CreatedAt:  completedAt,
			},
			Score: &models.PronunciationScore{
				TotalScore:     score,
				AccuracyScore:  score + 2,
				FluencyScore:   score - 3,
				PronuncScore:   score + 1,
				ExpectedText:   fmt.Sprintf("Здравствуйте! Как дела? (ページ %d)", i),
				RecognizedText: fmt.Sprintf("Здравствуйте! Как дела? (ページ %d)", i),
				EvaluationID:   jobID,
				UserID:         testUserID,
				Feedback: &models.Feedback{
					Level:          "good",
					Message:        "良い発音です！",
					PositivePoints: []string{"イントネーションが自然", "発音が明瞭"},
					Improvements:   []string{"もう少し流暢に話しましょう"},
					SpecificAdvice: []string{"「как дела」の部分をもっとはっきりと"},
				},
				CreatedAt: completedAt,
			},
			CreatedAt:   now,
			UpdatedAt:   completedAt,
			CompletedAt: &completedAt,
		}

		r.jobs[jobID] = job
		r.userJobs[testUserID] = append(r.userJobs[testUserID], jobID)
		r.bookJobs[testBookID] = append(r.bookJobs[testBookID], jobID)
	}

	// 処理中のSTTジョブ（ページ26-28）
	for i := 26; i <= 28; i++ {
		jobID := uuid.New().String()
		now := time.Now().Add(-time.Duration(2) * time.Minute)

		job := &models.STTJobDetail{
			ID:            jobID,
			BookID:        testBookID,
			PageNumber:    i,
			AudioURL:      fmt.Sprintf("/storage/audio/recording_%d.wav", i),
			ReferenceText: fmt.Sprintf("こんにちは！元気ですか？ (ページ %d)", i),
			Language:      "ja",
			Status:        models.STTStatusProcessing,
			Progress:      50 + (i-26)*15,
			CreatedAt:     now,
			UpdatedAt:     time.Now(),
		}

		r.jobs[jobID] = job
		r.userJobs[testUserID] = append(r.userJobs[testUserID], jobID)
		r.bookJobs[testBookID] = append(r.bookJobs[testBookID], jobID)
	}

	// ペンディングのSTTジョブ（ページ29-30）
	for i := 29; i <= 30; i++ {
		jobID := uuid.New().String()
		now := time.Now().Add(-time.Duration(30) * time.Second)

		job := &models.STTJobDetail{
			ID:            jobID,
			BookID:        testBookID,
			PageNumber:    i,
			AudioURL:      fmt.Sprintf("/storage/audio/recording_%d.wav", i),
			ReferenceText: fmt.Sprintf("Hello! How are you? (Page %d)", i),
			Language:      "en",
			Status:        models.STTStatusPending,
			Progress:      0,
			CreatedAt:     now,
			UpdatedAt:     now,
		}

		r.jobs[jobID] = job
		r.userJobs[testUserID] = append(r.userJobs[testUserID], jobID)
		r.bookJobs[testBookID] = append(r.bookJobs[testBookID], jobID)
	}
}

func (r *InMemorySTTRepository) CreateJob(ctx context.Context, userID, bookID uuid.UUID, pageNumber int, audioURL string, language string, referenceText string, options models.STTRecognizeOptions) (*models.STTJobDetail, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	jobID := uuid.New().String()
	now := time.Now()

	job := &models.STTJobDetail{
		ID:            jobID,
		BookID:        bookID.String(),
		PageNumber:    pageNumber,
		AudioURL:      audioURL,
		ReferenceText: referenceText,
		Language:      language,
		Status:        models.STTStatusPending,
		Progress:      0,
		CreatedAt:     now,
		UpdatedAt:     now,
	}

	r.jobs[jobID] = job
	r.userJobs[userID.String()] = append(r.userJobs[userID.String()], jobID)
	r.bookJobs[bookID.String()] = append(r.bookJobs[bookID.String()], jobID)

	return job, nil
}

func (r *InMemorySTTRepository) GetJob(ctx context.Context, jobID string) (*models.STTJobDetail, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	job, exists := r.jobs[jobID]
	if !exists {
		return nil, fmt.Errorf("job not found: %s", jobID)
	}

	return job, nil
}

func (r *InMemorySTTRepository) UpdateJobStatus(ctx context.Context, jobID string, status models.STTStatus, progress int) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	job, exists := r.jobs[jobID]
	if !exists {
		return fmt.Errorf("job not found: %s", jobID)
	}

	job.Status = status
	job.Progress = progress
	job.UpdatedAt = time.Now()

	if status == models.STTStatusCompleted || status == models.STTStatusFailed {
		now := time.Now()
		job.CompletedAt = &now
	}

	return nil
}

func (r *InMemorySTTRepository) UpdateJobResult(ctx context.Context, jobID string, result *models.STTResult, score *models.PronunciationScore) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	job, exists := r.jobs[jobID]
	if !exists {
		return fmt.Errorf("job not found: %s", jobID)
	}

	job.Result = result
	job.Score = score
	job.Status = models.STTStatusCompleted
	job.Progress = 100
	job.UpdatedAt = time.Now()
	now := time.Now()
	job.CompletedAt = &now

	return nil
}

func (r *InMemorySTTRepository) UpdateJobError(ctx context.Context, jobID string, errorMsg string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	job, exists := r.jobs[jobID]
	if !exists {
		return fmt.Errorf("job not found: %s", jobID)
	}

	job.Error = errorMsg
	job.Status = models.STTStatusFailed
	job.UpdatedAt = time.Now()
	now := time.Now()
	job.CompletedAt = &now

	return nil
}

func (r *InMemorySTTRepository) GetJobsByBookID(ctx context.Context, bookID uuid.UUID) ([]*models.STTJobDetail, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	jobIDs, exists := r.bookJobs[bookID.String()]
	if !exists {
		return []*models.STTJobDetail{}, nil
	}

	jobs := make([]*models.STTJobDetail, 0, len(jobIDs))
	for _, jobID := range jobIDs {
		if job, exists := r.jobs[jobID]; exists {
			jobs = append(jobs, job)
		}
	}

	return jobs, nil
}

func (r *InMemorySTTRepository) GetJobsByUserID(ctx context.Context, userID uuid.UUID) ([]*models.STTJobDetail, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	jobIDs, exists := r.userJobs[userID.String()]
	if !exists {
		return []*models.STTJobDetail{}, nil
	}

	jobs := make([]*models.STTJobDetail, 0, len(jobIDs))
	for _, jobID := range jobIDs {
		if job, exists := r.jobs[jobID]; exists {
			jobs = append(jobs, job)
		}
	}

	return jobs, nil
}

func (r *InMemorySTTRepository) GetStatistics(ctx context.Context, userID uuid.UUID) (*models.STTStatistics, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	stats := &models.STTStatistics{
		TotalRecognitions: 0,
		TotalEvaluations:  0,
		AverageScore:      0.0,
		BestScore:         0,
		LanguageStats:     make(map[string]int),
		RecentJobs:        make([]*models.STTJobDetail, 0),
		TotalDuration:     0,
	}

	jobIDs, exists := r.userJobs[userID.String()]
	if !exists {
		return stats, nil
	}

	totalScore := 0
	scoreCount := 0

	for _, jobID := range jobIDs {
		if job, exists := r.jobs[jobID]; exists {
			if job.Status == models.STTStatusCompleted {
				stats.TotalRecognitions++
				stats.LanguageStats[job.Language]++

				if job.Result != nil {
					stats.TotalDuration += int(job.Result.Duration)
				}

				if job.Score != nil {
					stats.TotalEvaluations++
					totalScore += job.Score.TotalScore
					scoreCount++

					if job.Score.TotalScore > stats.BestScore {
						stats.BestScore = job.Score.TotalScore
					}
				}

				// 最近のジョブ（最大10件）
				if len(stats.RecentJobs) < 10 {
					stats.RecentJobs = append(stats.RecentJobs, job)
				}
			}
		}
	}

	if scoreCount > 0 {
		stats.AverageScore = float64(totalScore) / float64(scoreCount)
	}

	return stats, nil
}

func (r *InMemorySTTRepository) DeleteJob(ctx context.Context, jobID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	_, exists := r.jobs[jobID]
	if !exists {
		return fmt.Errorf("job not found: %s", jobID)
	}

	delete(r.jobs, jobID)

	return nil
}

func (r *InMemorySTTRepository) GetSupportedLanguages(ctx context.Context) ([]*models.STTLanguage, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	return r.languages, nil
}

// SimulateSTTProcessing はSTT処理をシミュレート（テスト用）
func (r *InMemorySTTRepository) SimulateSTTProcessing(ctx context.Context, jobID string) error {
	// ステータスを処理中に更新
	if err := r.UpdateJobStatus(ctx, jobID, models.STTStatusProcessing, 10); err != nil {
		return err
	}

	// 処理進行をシミュレート
	for i := 20; i <= 90; i += 20 {
		time.Sleep(200 * time.Millisecond)
		if err := r.UpdateJobStatus(ctx, jobID, models.STTStatusProcessing, i); err != nil {
			return err
		}
	}

	// 結果を生成
	job, err := r.GetJob(ctx, jobID)
	if err != nil {
		return err
	}

	result := &models.STTResult{
		Text:       job.ReferenceText,
		Language:   job.Language,
		Confidence: 0.92,
		Duration:   3.5,
		Words: []models.WordInfo{
			{Word: "Здравствуйте", StartTime: 0.0, EndTime: 0.8, Confidence: 0.95},
			{Word: "Как", StartTime: 0.9, EndTime: 1.1, Confidence: 0.90},
			{Word: "дела", StartTime: 1.2, EndTime: 1.5, Confidence: 0.91},
		},
		CreatedAt: time.Now(),
	}

	score := &models.PronunciationScore{
		TotalScore:     85,
		AccuracyScore:  87,
		FluencyScore:   82,
		PronuncScore:   86,
		ExpectedText:   job.ReferenceText,
		RecognizedText: job.ReferenceText,
		EvaluationID:   jobID,
		UserID:         "550e8400-e29b-41d4-a716-446655440000",
		Feedback: &models.Feedback{
			Level:          "good",
			Message:        "素晴らしい発音です！",
			PositivePoints: []string{"イントネーションが自然", "発音の明瞭さ"},
			Improvements:   []string{"もう少し流暢に話しましょう"},
			SpecificAdvice: []string{"「дела」の部分をもっとはっきりと発音しましょう"},
		},
		CreatedAt: time.Now(),
	}

	return r.UpdateJobResult(ctx, jobID, result, score)
}

// PostgreSQL Implementation

type STTRepositoryPostgres struct {
	db *sql.DB
}

func NewSTTRepositoryPostgres(db *sql.DB) STTRepositoryInterface {
	return &STTRepositoryPostgres{db: db}
}

func (r *STTRepositoryPostgres) CreateJob(ctx context.Context, userID, bookID uuid.UUID, pageNumber int, audioURL string, language string, referenceText string, options models.STTRecognizeOptions) (*models.STTJobDetail, error) {
	jobID := uuid.New()
	now := time.Now()

	_, err := r.db.ExecContext(ctx, `
		INSERT INTO stt_jobs (id, user_id, book_id, audio_url, language, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`, jobID, userID, bookID, audioURL, language, "pending", now, now)

	if err != nil {
		return nil, err
	}

	return &models.STTJobDetail{
		ID:            jobID.String(),
		BookID:        bookID.String(),
		PageNumber:    pageNumber,
		AudioURL:      audioURL,
		ReferenceText: referenceText,
		Language:      language,
		Status:        models.STTStatusPending,
		Progress:      0,
		CreatedAt:     now,
		UpdatedAt:     now,
	}, nil
}

func (r *STTRepositoryPostgres) GetJob(ctx context.Context, jobID string) (*models.STTJobDetail, error) {
	var job models.STTJobDetail
	var status string
	var errorMsg sql.NullString
	var completedAt sql.NullTime

	err := r.db.QueryRowContext(ctx, `
		SELECT id, book_id, audio_url, language, status, error_message, created_at, updated_at, completed_at
		FROM stt_jobs
		WHERE id = $1
	`, jobID).Scan(&job.ID, &job.BookID, &job.AudioURL, &job.Language, &status, &errorMsg, &job.CreatedAt, &job.UpdatedAt, &completedAt)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("job not found: %s", jobID)
	}
	if err != nil {
		return nil, err
	}

	job.Status = models.STTStatus(status)
	if errorMsg.Valid {
		job.Error = errorMsg.String
	}
	if completedAt.Valid {
		job.CompletedAt = &completedAt.Time
	}

	// Get result if completed
	if job.Status == models.STTStatusCompleted {
		var transcript string
		var confidence sql.NullFloat64
		var score sql.NullInt64
		var feedbackJSON []byte

		err = r.db.QueryRowContext(ctx, `
			SELECT transcript, confidence, pronunciation_score, pronunciation_feedback
			FROM stt_results
			WHERE job_id = $1
		`, jobID).Scan(&transcript, &confidence, &score, &feedbackJSON)

		if err == nil {
			job.Result = &models.STTResult{
				Text:       transcript,
				Language:   job.Language,
				Confidence: confidence.Float64,
			}

			if score.Valid {
				var feedback models.Feedback
				if len(feedbackJSON) > 0 {
					json.Unmarshal(feedbackJSON, &feedback)
				}

				job.Score = &models.PronunciationScore{
					TotalScore: int(score.Int64),
					Feedback:   &feedback,
				}
			}
		}
	}

	return &job, nil
}

func (r *STTRepositoryPostgres) UpdateJobStatus(ctx context.Context, jobID string, status models.STTStatus, progress int) error {
	query := `UPDATE stt_jobs SET status = $1, updated_at = $2 WHERE id = $3`
	args := []interface{}{string(status), time.Now(), jobID}

	if status == models.STTStatusCompleted || status == models.STTStatusFailed {
		query = `UPDATE stt_jobs SET status = $1, updated_at = $2, completed_at = $2 WHERE id = $3`
	}

	_, err := r.db.ExecContext(ctx, query, args...)
	return err
}

func (r *STTRepositoryPostgres) UpdateJobResult(ctx context.Context, jobID string, result *models.STTResult, score *models.PronunciationScore) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Update job status
	_, err = tx.ExecContext(ctx, `
		UPDATE stt_jobs SET status = $1, updated_at = $2, completed_at = $2
		WHERE id = $3
	`, string(models.STTStatusCompleted), time.Now(), jobID)
	if err != nil {
		return err
	}

	// Prepare feedback JSON
	var feedbackJSON []byte
	if score != nil && score.Feedback != nil {
		feedbackJSON, _ = json.Marshal(score.Feedback)
	}

	// Insert result
	pronunciationScore := sql.NullInt64{}
	if score != nil {
		pronunciationScore.Int64 = int64(score.TotalScore)
		pronunciationScore.Valid = true
	}

	_, err = tx.ExecContext(ctx, `
		INSERT INTO stt_results (job_id, transcript, confidence, pronunciation_score, pronunciation_feedback)
		VALUES ($1, $2, $3, $4, $5)
	`, jobID, result.Text, result.Confidence, pronunciationScore, feedbackJSON)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (r *STTRepositoryPostgres) UpdateJobError(ctx context.Context, jobID string, errorMsg string) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE stt_jobs SET status = $1, error_message = $2, updated_at = $3, completed_at = $3
		WHERE id = $4
	`, string(models.STTStatusFailed), errorMsg, time.Now(), jobID)
	return err
}

func (r *STTRepositoryPostgres) GetJobsByBookID(ctx context.Context, bookID uuid.UUID) ([]*models.STTJobDetail, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, book_id, audio_url, language, status, error_message, created_at, updated_at, completed_at
		FROM stt_jobs
		WHERE book_id = $1
		ORDER BY created_at DESC
	`, bookID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	jobs := []*models.STTJobDetail{}
	for rows.Next() {
		var job models.STTJobDetail
		var status string
		var errorMsg sql.NullString
		var completedAt sql.NullTime

		err := rows.Scan(&job.ID, &job.BookID, &job.AudioURL, &job.Language, &status, &errorMsg, &job.CreatedAt, &job.UpdatedAt, &completedAt)
		if err != nil {
			return nil, err
		}

		job.Status = models.STTStatus(status)
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

func (r *STTRepositoryPostgres) GetJobsByUserID(ctx context.Context, userID uuid.UUID) ([]*models.STTJobDetail, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, book_id, audio_url, language, status, error_message, created_at, updated_at, completed_at
		FROM stt_jobs
		WHERE user_id = $1
		ORDER BY created_at DESC
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	jobs := []*models.STTJobDetail{}
	for rows.Next() {
		var job models.STTJobDetail
		var status string
		var errorMsg sql.NullString
		var completedAt sql.NullTime

		err := rows.Scan(&job.ID, &job.BookID, &job.AudioURL, &job.Language, &status, &errorMsg, &job.CreatedAt, &job.UpdatedAt, &completedAt)
		if err != nil {
			return nil, err
		}

		job.Status = models.STTStatus(status)
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

func (r *STTRepositoryPostgres) GetStatistics(ctx context.Context, userID uuid.UUID) (*models.STTStatistics, error) {
	stats := &models.STTStatistics{
		LanguageStats: make(map[string]int),
		RecentJobs:    make([]*models.STTJobDetail, 0),
	}

	// Get total recognitions and evaluations
	err := r.db.QueryRowContext(ctx, `
		SELECT
			COUNT(*) as total_recognitions,
			COUNT(sr.id) as total_evaluations,
			COALESCE(AVG(sr.pronunciation_score), 0) as avg_score,
			COALESCE(MAX(sr.pronunciation_score), 0) as best_score
		FROM stt_jobs sj
		LEFT JOIN stt_results sr ON sj.id = sr.job_id
		WHERE sj.user_id = $1 AND sj.status = 'completed'
	`, userID).Scan(&stats.TotalRecognitions, &stats.TotalEvaluations, &stats.AverageScore, &stats.BestScore)

	if err != nil {
		return nil, err
	}

	return stats, nil
}

func (r *STTRepositoryPostgres) DeleteJob(ctx context.Context, jobID string) error {
	_, err := r.db.ExecContext(ctx, "DELETE FROM stt_jobs WHERE id = $1", jobID)
	return err
}

func (r *STTRepositoryPostgres) GetSupportedLanguages(ctx context.Context) ([]*models.STTLanguage, error) {
	// Return hardcoded supported languages
	languages := []*models.STTLanguage{
		{Code: "ja", Name: "Japanese", NativeName: "日本語", IsSupported: true, SupportsPronunciation: true},
		{Code: "en", Name: "English", NativeName: "English", IsSupported: true, SupportsPronunciation: true},
		{Code: "zh", Name: "Chinese", NativeName: "中文", IsSupported: true, SupportsPronunciation: true},
		{Code: "ru", Name: "Russian", NativeName: "Русский", IsSupported: true, SupportsPronunciation: true},
		{Code: "fa", Name: "Persian", NativeName: "فارسی", IsSupported: true, SupportsPronunciation: true},
		{Code: "he", Name: "Hebrew", NativeName: "עברית", IsSupported: true, SupportsPronunciation: true},
		{Code: "es", Name: "Spanish", NativeName: "Español", IsSupported: true, SupportsPronunciation: true},
		{Code: "fr", Name: "French", NativeName: "Français", IsSupported: true, SupportsPronunciation: true},
		{Code: "pt", Name: "Portuguese", NativeName: "Português", IsSupported: true, SupportsPronunciation: true},
		{Code: "de", Name: "German", NativeName: "Deutsch", IsSupported: true, SupportsPronunciation: true},
		{Code: "it", Name: "Italian", NativeName: "Italiano", IsSupported: true, SupportsPronunciation: true},
		{Code: "tr", Name: "Turkish", NativeName: "Türkçe", IsSupported: true, SupportsPronunciation: true},
	}
	return languages, nil
}
