package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/clearclown/HaiLanGo/backend/internal/models"
	"github.com/google/uuid"
)

// StatsRepository handles statistics data access
type StatsRepository struct {
	db *sql.DB
}

// NewStatsRepository creates a new stats repository
func NewStatsRepository(db *sql.DB) *StatsRepository {
	return &StatsRepository{
		db: db,
	}
}

// GetLearningSessions retrieves learning sessions within a date range
func (r *StatsRepository) GetLearningSessions(ctx context.Context, userID uuid.UUID, startDate, endDate time.Time) ([]models.LearningSession, error) {
	query := `
		SELECT id, user_id, book_id, page_id, start_time, end_time, duration, created_at
		FROM learning_sessions
		WHERE user_id = $1 AND start_time >= $2 AND start_time <= $3
		ORDER BY start_time ASC
	`

	rows, err := r.db.QueryContext(ctx, query, userID, startDate, endDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sessions []models.LearningSession
	for rows.Next() {
		var session models.LearningSession
		err := rows.Scan(
			&session.ID,
			&session.UserID,
			&session.BookID,
			&session.PageID,
			&session.StartTime,
			&session.EndTime,
			&session.Duration,
			&session.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		sessions = append(sessions, session)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return sessions, nil
}

// GetVocabularyProgress retrieves all vocabulary progress for a user
func (r *StatsRepository) GetVocabularyProgress(ctx context.Context, userID uuid.UUID) ([]models.VocabularyProgress, error) {
	query := `
		SELECT id, user_id, word, language, mastery_level, last_reviewed, review_count, created_at, updated_at
		FROM vocabulary_progress
		WHERE user_id = $1
		ORDER BY mastery_level ASC, last_reviewed DESC
	`

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var vocabList []models.VocabularyProgress
	for rows.Next() {
		var vocab models.VocabularyProgress
		err := rows.Scan(
			&vocab.ID,
			&vocab.UserID,
			&vocab.Word,
			&vocab.Language,
			&vocab.MasteryLevel,
			&vocab.LastReviewed,
			&vocab.ReviewCount,
			&vocab.CreatedAt,
			&vocab.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		vocabList = append(vocabList, vocab)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return vocabList, nil
}

// GetPhraseProgress retrieves all phrase progress for a user
func (r *StatsRepository) GetPhraseProgress(ctx context.Context, userID uuid.UUID) ([]models.PhraseProgress, error) {
	query := `
		SELECT id, user_id, phrase, language, mastery_level, last_reviewed, review_count, created_at, updated_at
		FROM phrase_progress
		WHERE user_id = $1
		ORDER BY mastery_level ASC, last_reviewed DESC
	`

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var phraseList []models.PhraseProgress
	for rows.Next() {
		var phrase models.PhraseProgress
		err := rows.Scan(
			&phrase.ID,
			&phrase.UserID,
			&phrase.Phrase,
			&phrase.Language,
			&phrase.MasteryLevel,
			&phrase.LastReviewed,
			&phrase.ReviewCount,
			&phrase.CreatedAt,
			&phrase.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		phraseList = append(phraseList, phrase)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return phraseList, nil
}

// GetPronunciationScores retrieves pronunciation scores within a date range
func (r *StatsRepository) GetPronunciationScores(ctx context.Context, userID uuid.UUID, startDate, endDate time.Time) ([]models.PronunciationScoreRecord, error) {
	query := `
		SELECT id, user_id, text, language, score, accuracy, fluency, created_at
		FROM pronunciation_scores
		WHERE user_id = $1 AND created_at >= $2 AND created_at <= $3
		ORDER BY created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, userID, startDate, endDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var scores []models.PronunciationScoreRecord
	for rows.Next() {
		var score models.PronunciationScoreRecord
		err := rows.Scan(
			&score.ID,
			&score.UserID,
			&score.Text,
			&score.Language,
			&score.Score,
			&score.Accuracy,
			&score.Fluency,
			&score.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		scores = append(scores, score)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return scores, nil
}

// GetCurrentStreak retrieves the current study streak
func (r *StatsRepository) GetCurrentStreak(ctx context.Context, userID uuid.UUID) (int, error) {
	query := `
		SELECT current_streak
		FROM user_streaks
		WHERE user_id = $1
	`

	var currentStreak int
	err := r.db.QueryRowContext(ctx, query, userID).Scan(&currentStreak)
	if err == sql.ErrNoRows {
		return 0, nil
	}
	if err != nil {
		return 0, err
	}

	return currentStreak, nil
}

// GetLongestStreak retrieves the longest study streak
func (r *StatsRepository) GetLongestStreak(ctx context.Context, userID uuid.UUID) (int, error) {
	query := `
		SELECT longest_streak
		FROM user_streaks
		WHERE user_id = $1
	`

	var longestStreak int
	err := r.db.QueryRowContext(ctx, query, userID).Scan(&longestStreak)
	if err == sql.ErrNoRows {
		return 0, nil
	}
	if err != nil {
		return 0, err
	}

	return longestStreak, nil
}

// GetCompletedPagesCount retrieves the count of completed pages
func (r *StatsRepository) GetCompletedPagesCount(ctx context.Context, userID uuid.UUID) (int, error) {
	query := `
		SELECT COUNT(DISTINCT page_id)
		FROM learning_sessions
		WHERE user_id = $1
	`

	var count int
	err := r.db.QueryRowContext(ctx, query, userID).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

// GetCompletedBooksCount retrieves the count of completed books
func (r *StatsRepository) GetCompletedBooksCount(ctx context.Context, userID uuid.UUID) (int, error) {
	query := `
		SELECT COUNT(DISTINCT book_id)
		FROM learning_sessions
		WHERE user_id = $1
	`

	var count int
	err := r.db.QueryRowContext(ctx, query, userID).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

// GetDashboardStats retrieves dashboard statistics for a user
func (r *StatsRepository) GetDashboardStats(ctx context.Context, userID uuid.UUID) (*models.DashboardStatsFlat, error) {
	stats := &models.DashboardStatsFlat{}

	// Get current streak
	currentStreak, _ := r.GetCurrentStreak(ctx, userID)
	stats.CurrentStreak = currentStreak

	// Get longest streak
	longestStreak, _ := r.GetLongestStreak(ctx, userID)
	stats.LongestStreak = longestStreak

	// Get completed pages count
	completedPages, _ := r.GetCompletedPagesCount(ctx, userID)
	stats.CompletedPages = completedPages

	// Get completed books count
	completedBooks, _ := r.GetCompletedBooksCount(ctx, userID)
	stats.CompletedBooks = completedBooks

	// Get total study time (sum of all session durations)
	query := `
		SELECT COALESCE(SUM(EXTRACT(EPOCH FROM (end_time - start_time))), 0)
		FROM learning_sessions
		WHERE user_id = $1
	`
	var totalSeconds float64
	_ = r.db.QueryRowContext(ctx, query, userID).Scan(&totalSeconds)
	stats.TotalLearningTime = int(totalSeconds / 60)

	return stats, nil
}

// GetLearningTimeData retrieves learning time data for a given period
func (r *StatsRepository) GetLearningTimeData(ctx context.Context, userID uuid.UUID, period string) (*models.LearningTimeData, error) {
	// Implementation simplified - return empty data for now
	return &models.LearningTimeData{
		Period:         period,
		Data:           []models.DailyLearningTime{},
		TotalMinutes:   0,
		AverageMinutes: 0,
	}, nil
}

// GetProgressData retrieves progress data for a given period
func (r *StatsRepository) GetProgressData(ctx context.Context, userID uuid.UUID, period string) (*models.ProgressData, error) {
	// Implementation simplified - return empty data for now
	return &models.ProgressData{
		Period:  period,
		Words:   []models.TimeSeriesData{},
		Phrases: []models.TimeSeriesData{},
		Pages:   []models.TimeSeriesData{},
	}, nil
}

// GetWeakPoints retrieves weak points (words/phrases with low scores)
func (r *StatsRepository) GetWeakPoints(ctx context.Context, userID uuid.UUID, limit int) (*models.WeakPointsData, error) {
	// Implementation simplified - return empty data for now
	return &models.WeakPointsData{
		WeakWords:   []models.WeakItem{},
		WeakPhrases: []models.WeakItem{},
	}, nil
}

// RecordLearningSession records a learning session
func (r *StatsRepository) RecordLearningSession(ctx context.Context, session *models.LearningSession) error {
	query := `
		INSERT INTO learning_sessions (id, user_id, book_id, page_id, start_time, end_time, duration, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`

	_, err := r.db.ExecContext(
		ctx,
		query,
		session.ID,
		session.UserID,
		session.BookID,
		session.PageID,
		session.StartTime,
		session.EndTime,
		session.Duration,
		session.CreatedAt,
	)

	return err
}

// UpdateUserProgress updates user progress for a specific day
func (r *StatsRepository) UpdateUserProgress(ctx context.Context, progress *models.UserProgressDaily) error {
	// Implementation simplified - no daily progress tracking table in current schema
	return nil
}

// UpdateStreak updates the user's study streak
func (r *StatsRepository) UpdateStreak(ctx context.Context, userID uuid.UUID, activityDate time.Time) error {
	// Get the last study date
	query := `
		SELECT last_study_date, current_streak, longest_streak
		FROM user_streaks
		WHERE user_id = $1
	`

	var lastStudyDate time.Time
	var currentStreak, longestStreak int
	err := r.db.QueryRowContext(ctx, query, userID).Scan(&lastStudyDate, &currentStreak, &longestStreak)

	today := activityDate.Truncate(24 * time.Hour)

	if err == sql.ErrNoRows {
		// Create new streak record
		insertQuery := `
			INSERT INTO user_streaks (user_id, current_streak, longest_streak, last_study_date)
			VALUES ($1, 1, 1, $2)
		`
		_, err = r.db.ExecContext(ctx, insertQuery, userID, today)
		return err
	}

	if err != nil {
		return err
	}

	// Check if the user studied today
	if lastStudyDate.Equal(today) {
		// Already counted today
		return nil
	}

	yesterday := today.AddDate(0, 0, -1)
	if lastStudyDate.Equal(yesterday) {
		// Consecutive day
		currentStreak++
		if currentStreak > longestStreak {
			longestStreak = currentStreak
		}
	} else {
		// Streak broken
		currentStreak = 1
	}

	// Update streak
	updateQuery := `
		UPDATE user_streaks
		SET current_streak = $1, longest_streak = $2, last_study_date = $3, updated_at = NOW()
		WHERE user_id = $4
	`
	_, err = r.db.ExecContext(ctx, updateQuery, currentStreak, longestStreak, today, userID)
	return err
}
