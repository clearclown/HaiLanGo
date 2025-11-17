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

// OCRRepositoryInterface はOCRリポジトリのインターフェース
type OCRRepositoryInterface interface {
	// CreateJob はOCRジョブを作成
	CreateJob(ctx context.Context, userID, bookID uuid.UUID, pageNumber int, imageURL string, language string) (*models.OCRJobDetail, error)

	// GetJob はOCRジョブを取得
	GetJob(ctx context.Context, jobID string) (*models.OCRJobDetail, error)

	// UpdateJobStatus はOCRジョブのステータスを更新
	UpdateJobStatus(ctx context.Context, jobID string, status models.OCRStatus, progress int) error

	// UpdateJobResult はOCRジョブの結果を更新
	UpdateJobResult(ctx context.Context, jobID string, result *models.OCRResult) error

	// UpdateJobError はOCRジョブのエラーを更新
	UpdateJobError(ctx context.Context, jobID string, errorMsg string) error

	// GetJobsByBookID は書籍IDでOCRジョブを取得
	GetJobsByBookID(ctx context.Context, bookID uuid.UUID) ([]*models.OCRJobDetail, error)

	// GetJobsByUserID はユーザーIDでOCRジョブを取得
	GetJobsByUserID(ctx context.Context, userID uuid.UUID) ([]*models.OCRJobDetail, error)

	// GetStatistics はOCR統計情報を取得
	GetStatistics(ctx context.Context, userID uuid.UUID) (*models.OCRStatistics, error)

	// DeleteJob はOCRジョブを削除
	DeleteJob(ctx context.Context, jobID string) error
}

// InMemoryOCRRepository はインメモリOCRリポジトリ
type InMemoryOCRRepository struct {
	mu       sync.RWMutex
	jobs     map[string]*models.OCRJobDetail // jobID -> OCRJobDetail
	userJobs map[string][]string             // userID -> []jobID
	bookJobs map[string][]string             // bookID -> []jobID
}

// NewInMemoryOCRRepository はインメモリOCRリポジトリを作成
func NewInMemoryOCRRepository() *InMemoryOCRRepository {
	repo := &InMemoryOCRRepository{
		jobs:     make(map[string]*models.OCRJobDetail),
		userJobs: make(map[string][]string),
		bookJobs: make(map[string][]string),
	}

	// サンプルデータを初期化
	repo.initSampleData()

	return repo
}

func (r *InMemoryOCRRepository) initSampleData() {
	testUserID := "550e8400-e29b-41d4-a716-446655440000"
	testBookID := "550e8400-e29b-41d4-a716-446655440000"

	// 完了済みのOCRジョブ（最初の45ページ）
	for i := 1; i <= 45; i++ {
		jobID := uuid.New().String()
		now := time.Now().Add(-time.Duration(46-i) * time.Hour)
		completedAt := now.Add(5 * time.Minute)

		result := &models.OCRResult{
			Text:             fmt.Sprintf("Здравствуйте! Как дела? (ページ %d)", i),
			DetectedLanguage: "ru",
			Confidence:       0.95 + float64(i%5)*0.01,
			Words: []models.OCRWord{
				{
					Text:       "Здравствуйте",
					Confidence: 0.98,
					BoundingBox: models.BoundingBox{X: 10, Y: 10, Width: 150, Height: 30},
				},
				{
					Text:       "Как",
					Confidence: 0.97,
					BoundingBox: models.BoundingBox{X: 170, Y: 10, Width: 50, Height: 30},
				},
				{
					Text:       "дела",
					Confidence: 0.96,
					BoundingBox: models.BoundingBox{X: 230, Y: 10, Width: 70, Height: 30},
				},
			},
			Lines: []models.OCRLine{
				{
					Text:       "Здравствуйте! Как дела?",
					Confidence: 0.97,
					BoundingBox: models.BoundingBox{X: 10, Y: 10, Width: 300, Height: 30},
					Words: []models.OCRWord{
						{Text: "Здравствуйте", Confidence: 0.98},
						{Text: "Как", Confidence: 0.97},
						{Text: "дела", Confidence: 0.96},
					},
				},
			},
			Blocks: []models.OCRBlock{
				{
					Type:       "text",
					Text:       "Здравствуйте! Как дела?",
					Confidence: 0.97,
					BoundingBox: models.BoundingBox{X: 10, Y: 10, Width: 300, Height: 30},
				},
			},
			HasRuby:        false,
			Orientation:    0,
			ProcessingTime: 3000 + i*10,
		}

		job := &models.OCRJobDetail{
			ID:          jobID,
			BookID:      testBookID,
			PageNumber:  i,
			ImageURL:    fmt.Sprintf("/storage/books/%s/pages/%d.jpg", testBookID, i),
			Status:      models.OCRStatusCompleted,
			Progress:    100,
			Result:      result,
			CreatedAt:   now,
			UpdatedAt:   completedAt,
			CompletedAt: &completedAt,
		}

		r.jobs[jobID] = job
		r.userJobs[testUserID] = append(r.userJobs[testUserID], jobID)
		r.bookJobs[testBookID] = append(r.bookJobs[testBookID], jobID)
	}

	// 処理中のOCRジョブ（ページ46-48）
	for i := 46; i <= 48; i++ {
		jobID := uuid.New().String()
		now := time.Now().Add(-time.Duration(5) * time.Minute)

		job := &models.OCRJobDetail{
			ID:         jobID,
			BookID:     testBookID,
			PageNumber: i,
			ImageURL:   fmt.Sprintf("/storage/books/%s/pages/%d.jpg", testBookID, i),
			Status:     models.OCRStatusProcessing,
			Progress:   50 + (i-46)*10,
			CreatedAt:  now,
			UpdatedAt:  time.Now(),
		}

		r.jobs[jobID] = job
		r.userJobs[testUserID] = append(r.userJobs[testUserID], jobID)
		r.bookJobs[testBookID] = append(r.bookJobs[testBookID], jobID)
	}

	// ペンディングのOCRジョブ（ページ49-50）
	for i := 49; i <= 50; i++ {
		jobID := uuid.New().String()
		now := time.Now().Add(-time.Duration(1) * time.Minute)

		job := &models.OCRJobDetail{
			ID:         jobID,
			BookID:     testBookID,
			PageNumber: i,
			ImageURL:   fmt.Sprintf("/storage/books/%s/pages/%d.jpg", testBookID, i),
			Status:     models.OCRStatusPending,
			Progress:   0,
			CreatedAt:  now,
			UpdatedAt:  now,
		}

		r.jobs[jobID] = job
		r.userJobs[testUserID] = append(r.userJobs[testUserID], jobID)
		r.bookJobs[testBookID] = append(r.bookJobs[testBookID], jobID)
	}
}

func (r *InMemoryOCRRepository) CreateJob(ctx context.Context, userID, bookID uuid.UUID, pageNumber int, imageURL string, language string) (*models.OCRJobDetail, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	jobID := uuid.New().String()
	now := time.Now()

	job := &models.OCRJobDetail{
		ID:         jobID,
		BookID:     bookID.String(),
		PageNumber: pageNumber,
		ImageURL:   imageURL,
		Status:     models.OCRStatusPending,
		Progress:   0,
		CreatedAt:  now,
		UpdatedAt:  now,
	}

	r.jobs[jobID] = job
	r.userJobs[userID.String()] = append(r.userJobs[userID.String()], jobID)
	r.bookJobs[bookID.String()] = append(r.bookJobs[bookID.String()], jobID)

	return job, nil
}

func (r *InMemoryOCRRepository) GetJob(ctx context.Context, jobID string) (*models.OCRJobDetail, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	job, exists := r.jobs[jobID]
	if !exists {
		return nil, fmt.Errorf("job not found: %s", jobID)
	}

	return job, nil
}

func (r *InMemoryOCRRepository) UpdateJobStatus(ctx context.Context, jobID string, status models.OCRStatus, progress int) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	job, exists := r.jobs[jobID]
	if !exists {
		return fmt.Errorf("job not found: %s", jobID)
	}

	job.Status = status
	job.Progress = progress
	job.UpdatedAt = time.Now()

	if status == models.OCRStatusCompleted || status == models.OCRStatusFailed {
		now := time.Now()
		job.CompletedAt = &now
	}

	return nil
}

func (r *InMemoryOCRRepository) UpdateJobResult(ctx context.Context, jobID string, result *models.OCRResult) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	job, exists := r.jobs[jobID]
	if !exists {
		return fmt.Errorf("job not found: %s", jobID)
	}

	job.Result = result
	job.Status = models.OCRStatusCompleted
	job.Progress = 100
	job.UpdatedAt = time.Now()
	now := time.Now()
	job.CompletedAt = &now

	return nil
}

func (r *InMemoryOCRRepository) UpdateJobError(ctx context.Context, jobID string, errorMsg string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	job, exists := r.jobs[jobID]
	if !exists {
		return fmt.Errorf("job not found: %s", jobID)
	}

	job.Error = errorMsg
	job.Status = models.OCRStatusFailed
	job.UpdatedAt = time.Now()
	now := time.Now()
	job.CompletedAt = &now

	return nil
}

func (r *InMemoryOCRRepository) GetJobsByBookID(ctx context.Context, bookID uuid.UUID) ([]*models.OCRJobDetail, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	jobIDs, exists := r.bookJobs[bookID.String()]
	if !exists {
		return []*models.OCRJobDetail{}, nil
	}

	jobs := make([]*models.OCRJobDetail, 0, len(jobIDs))
	for _, jobID := range jobIDs {
		if job, exists := r.jobs[jobID]; exists {
			jobs = append(jobs, job)
		}
	}

	return jobs, nil
}

func (r *InMemoryOCRRepository) GetJobsByUserID(ctx context.Context, userID uuid.UUID) ([]*models.OCRJobDetail, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	jobIDs, exists := r.userJobs[userID.String()]
	if !exists {
		return []*models.OCRJobDetail{}, nil
	}

	jobs := make([]*models.OCRJobDetail, 0, len(jobIDs))
	for _, jobID := range jobIDs {
		if job, exists := r.jobs[jobID]; exists {
			jobs = append(jobs, job)
		}
	}

	return jobs, nil
}

func (r *InMemoryOCRRepository) GetStatistics(ctx context.Context, userID uuid.UUID) (*models.OCRStatistics, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	jobIDs, exists := r.userJobs[userID.String()]
	if !exists {
		return &models.OCRStatistics{}, nil
	}

	stats := &models.OCRStatistics{}
	totalConfidence := 0.0
	confidenceCount := 0

	for _, jobID := range jobIDs {
		job, exists := r.jobs[jobID]
		if !exists {
			continue
		}

		stats.TotalJobs++

		switch job.Status {
		case models.OCRStatusCompleted:
			stats.CompletedJobs++
			if job.Result != nil {
				totalConfidence += job.Result.Confidence
				confidenceCount++
				stats.TotalProcessingTime += job.Result.ProcessingTime
			}
		case models.OCRStatusFailed:
			stats.FailedJobs++
		case models.OCRStatusPending:
			stats.PendingJobs++
		case models.OCRStatusProcessing:
			stats.ProcessingJobs++
		}
	}

	if confidenceCount > 0 {
		stats.AverageConfidence = totalConfidence / float64(confidenceCount)
	}

	return stats, nil
}

func (r *InMemoryOCRRepository) DeleteJob(ctx context.Context, jobID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	_, exists := r.jobs[jobID]
	if !exists {
		return fmt.Errorf("job not found: %s", jobID)
	}

	// ユーザーとブックのジョブリストから削除
	// これは簡略化した実装（実際にはより効率的な方法が必要）
	delete(r.jobs, jobID)

	return nil
}

// SimulateOCRProcessing はOCR処理をシミュレート（テスト用）
func (r *InMemoryOCRRepository) SimulateOCRProcessing(ctx context.Context, jobID string) error {
	// ステータスを処理中に更新
	if err := r.UpdateJobStatus(ctx, jobID, models.OCRStatusProcessing, 10); err != nil {
		return err
	}

	// 処理進行をシミュレート
	for i := 20; i <= 90; i += 20 {
		time.Sleep(500 * time.Millisecond)
		if err := r.UpdateJobStatus(ctx, jobID, models.OCRStatusProcessing, i); err != nil {
			return err
		}
	}

	// 結果を生成
	result := &models.OCRResult{
		Text:             "Здравствуйте! Как дела?",
		DetectedLanguage: "ru",
		Confidence:       0.95,
		Words: []models.OCRWord{
			{Text: "Здравствуйте", Confidence: 0.98, BoundingBox: models.BoundingBox{X: 10, Y: 10, Width: 150, Height: 30}},
			{Text: "Как", Confidence: 0.97, BoundingBox: models.BoundingBox{X: 170, Y: 10, Width: 50, Height: 30}},
			{Text: "дела", Confidence: 0.96, BoundingBox: models.BoundingBox{X: 230, Y: 10, Width: 70, Height: 30}},
		},
		Lines: []models.OCRLine{
			{
				Text:       "Здравствуйте! Как дела?",
				Confidence: 0.97,
				BoundingBox: models.BoundingBox{X: 10, Y: 10, Width: 300, Height: 30},
			},
		},
		Blocks: []models.OCRBlock{
			{
				Type:       "text",
				Text:       "Здравствуйте! Как дела?",
				Confidence: 0.97,
				BoundingBox: models.BoundingBox{X: 10, Y: 10, Width: 300, Height: 30},
			},
		},
		HasRuby:        false,
		Orientation:    0,
		ProcessingTime: 3500,
	}

	return r.UpdateJobResult(ctx, jobID, result)
}

// Helper function to convert result to JSON string (for database storage)
func resultToJSON(result *models.OCRResult) string {
	if result == nil {
		return ""
	}
	data, _ := json.Marshal(result)
	return string(data)
}

// Helper function to parse JSON string to result (from database)
func jsonToResult(jsonStr string) *models.OCRResult {
	if jsonStr == "" {
		return nil
	}
	var result models.OCRResult
	json.Unmarshal([]byte(jsonStr), &result)
	return &result
}

// PostgreSQL Implementation

type OCRRepositoryPostgres struct {
	db *sql.DB
}

func NewOCRRepositoryPostgres(db *sql.DB) OCRRepositoryInterface {
	return &OCRRepositoryPostgres{db: db}
}

func (r *OCRRepositoryPostgres) CreateJob(ctx context.Context, userID, bookID uuid.UUID, pageNumber int, imageURL string, language string) (*models.OCRJobDetail, error) {
	jobID := uuid.New()
	now := time.Now()
	
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO ocr_jobs (id, user_id, book_id, page_number, status, progress, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`, jobID, userID, bookID, pageNumber, "pending", 0, now, now)
	
	if err != nil {
		return nil, err
	}
	
	return &models.OCRJobDetail{
		ID:         jobID.String(),
		BookID:     bookID.String(),
		PageNumber: pageNumber,
		ImageURL:   imageURL,
		Status:     models.OCRStatusPending,
		Progress:   0,
		CreatedAt:  now,
		UpdatedAt:  now,
	}, nil
}

func (r *OCRRepositoryPostgres) GetJob(ctx context.Context, jobID string) (*models.OCRJobDetail, error) {
	var job models.OCRJobDetail
	var status string
	var errorMsg sql.NullString
	var completedAt sql.NullTime
	
	err := r.db.QueryRowContext(ctx, `
		SELECT id, book_id, page_number, status, progress, error_message, created_at, updated_at, completed_at
		FROM ocr_jobs
		WHERE id = $1
	`, jobID).Scan(&job.ID, &job.BookID, &job.PageNumber, &status, &job.Progress, &errorMsg, &job.CreatedAt, &job.UpdatedAt, &completedAt)
	
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("job not found: %s", jobID)
	}
	if err != nil {
		return nil, err
	}
	
	job.Status = models.OCRStatus(status)
	if errorMsg.Valid {
		job.Error = errorMsg.String
	}
	if completedAt.Valid {
		job.CompletedAt = &completedAt.Time
	}
	
	return &job, nil
}

func (r *OCRRepositoryPostgres) UpdateJobStatus(ctx context.Context, jobID string, status models.OCRStatus, progress int) error {
	query := `UPDATE ocr_jobs SET status = $1, progress = $2, updated_at = $3 WHERE id = $4`
	args := []interface{}{string(status), progress, time.Now(), jobID}
	
	if status == models.OCRStatusCompleted || status == models.OCRStatusFailed {
		query = `UPDATE ocr_jobs SET status = $1, progress = $2, updated_at = $3, completed_at = $3 WHERE id = $4`
	}
	
	_, err := r.db.ExecContext(ctx, query, args...)
	return err
}

func (r *OCRRepositoryPostgres) UpdateJobResult(ctx context.Context, jobID string, result *models.OCRResult) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	
	_, err = tx.ExecContext(ctx, `
		UPDATE ocr_jobs SET status = $1, progress = $2, updated_at = $3, completed_at = $3
		WHERE id = $4
	`, string(models.OCRStatusCompleted), 100, time.Now(), jobID)
	if err != nil {
		return err
	}
	
	job, err := r.GetJob(ctx, jobID)
	if err != nil {
		return err
	}
	
	_, err = tx.ExecContext(ctx, `
		INSERT INTO ocr_results (job_id, page_number, original_text, translated_text, confidence, processing_time_ms)
		VALUES ($1, $2, $3, $4, $5, $6)
	`, jobID, job.PageNumber, result.Text, "", result.Confidence, result.ProcessingTime)
	if err != nil {
		return err
	}
	
	return tx.Commit()
}

func (r *OCRRepositoryPostgres) UpdateJobError(ctx context.Context, jobID string, errorMsg string) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE ocr_jobs SET status = $1, error_message = $2, updated_at = $3, completed_at = $3
		WHERE id = $4
	`, string(models.OCRStatusFailed), errorMsg, time.Now(), jobID)
	return err
}

func (r *OCRRepositoryPostgres) GetJobsByBookID(ctx context.Context, bookID uuid.UUID) ([]*models.OCRJobDetail, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, book_id, page_number, status, progress, error_message, created_at, updated_at, completed_at
		FROM ocr_jobs
		WHERE book_id = $1
		ORDER BY created_at DESC
	`, bookID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	jobs := []*models.OCRJobDetail{}
	for rows.Next() {
		var job models.OCRJobDetail
		var status string
		var errorMsg sql.NullString
		var completedAt sql.NullTime
		
		err := rows.Scan(&job.ID, &job.BookID, &job.PageNumber, &status, &job.Progress, &errorMsg, &job.CreatedAt, &job.UpdatedAt, &completedAt)
		if err != nil {
			return nil, err
		}
		
		job.Status = models.OCRStatus(status)
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

func (r *OCRRepositoryPostgres) GetJobsByUserID(ctx context.Context, userID uuid.UUID) ([]*models.OCRJobDetail, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, book_id, page_number, status, progress, error_message, created_at, updated_at, completed_at
		FROM ocr_jobs
		WHERE user_id = $1
		ORDER BY created_at DESC
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	jobs := []*models.OCRJobDetail{}
	for rows.Next() {
		var job models.OCRJobDetail
		var status string
		var errorMsg sql.NullString
		var completedAt sql.NullTime
		
		err := rows.Scan(&job.ID, &job.BookID, &job.PageNumber, &status, &job.Progress, &errorMsg, &job.CreatedAt, &job.UpdatedAt, &completedAt)
		if err != nil {
			return nil, err
		}
		
		job.Status = models.OCRStatus(status)
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

func (r *OCRRepositoryPostgres) GetStatistics(ctx context.Context, userID uuid.UUID) (*models.OCRStatistics, error) {
	stats := &models.OCRStatistics{}
	
	err := r.db.QueryRowContext(ctx, `
		SELECT
			COUNT(*) as total_jobs,
			COUNT(CASE WHEN status = 'completed' THEN 1 END) as completed_jobs,
			COUNT(CASE WHEN status = 'failed' THEN 1 END) as failed_jobs,
			COUNT(CASE WHEN status = 'pending' THEN 1 END) as pending_jobs,
			COUNT(CASE WHEN status = 'processing' THEN 1 END) as processing_jobs
		FROM ocr_jobs
		WHERE user_id = $1
	`, userID).Scan(&stats.TotalJobs, &stats.CompletedJobs, &stats.FailedJobs, &stats.PendingJobs, &stats.ProcessingJobs)
	
	if err != nil {
		return nil, err
	}
	
	// Get average confidence and total processing time
	r.db.QueryRowContext(ctx, `
		SELECT
			COALESCE(AVG(confidence), 0) as avg_confidence,
			COALESCE(SUM(processing_time_ms), 0) as total_processing_time
		FROM ocr_results
		WHERE job_id IN (SELECT id FROM ocr_jobs WHERE user_id = $1)
	`, userID).Scan(&stats.AverageConfidence, &stats.TotalProcessingTime)
	
	return stats, nil
}

func (r *OCRRepositoryPostgres) DeleteJob(ctx context.Context, jobID string) error {
	_, err := r.db.ExecContext(ctx, "DELETE FROM ocr_jobs WHERE id = $1", jobID)
	return err
}
