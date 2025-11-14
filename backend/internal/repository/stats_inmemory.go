package repository

import (
	"context"
	"sync"
	"time"

	"github.com/clearclown/HaiLanGo/backend/internal/models"
	"github.com/google/uuid"
)

// StatsRepositoryInterface はStats APIのインターフェース
type StatsRepositoryInterface interface {
	GetDashboardStats(ctx context.Context, userID uuid.UUID) (*models.DashboardStatsFlat, error)
	GetLearningTimeData(ctx context.Context, userID uuid.UUID, period string) (*models.LearningTimeData, error)
	GetProgressData(ctx context.Context, userID uuid.UUID, period string) (*models.ProgressData, error)
	GetWeakPoints(ctx context.Context, userID uuid.UUID, limit int) (*models.WeakPointsData, error)
	RecordLearningSession(ctx context.Context, session *models.LearningSession) error
	UpdateUserProgress(ctx context.Context, progress *models.UserProgressDaily) error
	UpdateStreak(ctx context.Context, userID uuid.UUID, activityDate time.Time) error
}

// InMemoryStatsRepository はInMemory実装
type InMemoryStatsRepository struct {
	sessions map[string]*models.LearningSession
	progress map[string]map[string]*models.UserProgressDaily // userID -> date -> progress
	streaks  map[string]*models.LearningStreakRecord
	mu       sync.RWMutex
}

// NewInMemoryStatsRepository は新しいInMemoryStatsRepositoryを作成
func NewInMemoryStatsRepository() *InMemoryStatsRepository {
	repo := &InMemoryStatsRepository{
		sessions: make(map[string]*models.LearningSession),
		progress: make(map[string]map[string]*models.UserProgressDaily),
		streaks:  make(map[string]*models.LearningStreakRecord),
	}

	// サンプルデータ初期化
	repo.initSampleData()

	return repo
}

func (r *InMemoryStatsRepository) initSampleData() {
	// テストユーザーID
	testUserID := uuid.MustParse("550e8400-e29b-41d4-a716-446655440000")

	// ストリークデータ
	r.streaks[testUserID.String()] = &models.LearningStreakRecord{
		ID:               uuid.New(),
		UserID:           testUserID,
		CurrentStreak:    7,
		LongestStreak:    15,
		LastActivityDate: time.Now(),
		UpdatedAt:        time.Now(),
	}

	// 過去7日間の進捗データ
	r.progress[testUserID.String()] = make(map[string]*models.UserProgressDaily)
	for i := 0; i < 7; i++ {
		date := time.Now().AddDate(0, 0, -i)
		dateStr := date.Format("2006-01-02")
		r.progress[testUserID.String()][dateStr] = &models.UserProgressDaily{
			ID:                      uuid.New(),
			UserID:                  testUserID,
			Date:                    date,
			CompletedPages:          i * 2,
			MasteredWords:           i * 10,
			MasteredPhrases:         i * 2,
			LearningMinutes:         25 + i*5,
			PronunciationAttempts:   i * 3,
			PronunciationTotalScore: i * 250,
			UpdatedAt:               time.Now(),
		}
	}
}

func (r *InMemoryStatsRepository) GetDashboardStats(ctx context.Context, userID uuid.UUID) (*models.DashboardStatsFlat, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	streak, _ := r.streaks[userID.String()]

	// 進捗データを集計
	totalPages := 0
	totalWords := 0
	totalPhrases := 0
	todayMinutes := 0
	weekMinutes := 0
	totalMinutes := 0
	totalPronunciationScore := 0
	totalPronunciationAttempts := 0

	userProgress, exists := r.progress[userID.String()]
	if exists {
		for dateStr, prog := range userProgress {
			date, _ := time.Parse("2006-01-02", dateStr)

			totalPages += prog.CompletedPages
			totalWords += prog.MasteredWords
			totalPhrases += prog.MasteredPhrases
			totalMinutes += prog.LearningMinutes
			totalPronunciationScore += prog.PronunciationTotalScore
			totalPronunciationAttempts += prog.PronunciationAttempts

			if isToday(date) {
				todayMinutes = prog.LearningMinutes
			}
			if isThisWeek(date) {
				weekMinutes += prog.LearningMinutes
			}
		}
	}

	avgPronunciation := 0.0
	if totalPronunciationAttempts > 0 {
		avgPronunciation = float64(totalPronunciationScore) / float64(totalPronunciationAttempts)
	}

	stats := &models.DashboardStatsFlat{
		LearningTimeToday:         todayMinutes,
		LearningTimeThisWeek:      weekMinutes,
		TotalLearningTime:         totalMinutes,
		CurrentStreak:             0,
		LongestStreak:             0,
		CompletedPages:            totalPages,
		TotalPages:                150, // 仮の値
		MasteredWords:             totalWords,
		MasteredPhrases:           totalPhrases,
		CompletedBooks:            1,
		TotalBooks:                3,
		AveragePronunciationScore: avgPronunciation,
	}

	if streak != nil {
		stats.CurrentStreak = streak.CurrentStreak
		stats.LongestStreak = streak.LongestStreak
	}

	return stats, nil
}

func (r *InMemoryStatsRepository) GetLearningTimeData(ctx context.Context, userID uuid.UUID, period string) (*models.LearningTimeData, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	data := &models.LearningTimeData{
		Period: period,
		Data:   []models.DailyLearningTime{},
	}

	userProgress, exists := r.progress[userID.String()]
	if !exists {
		return data, nil
	}

	days := getDaysForPeriod(period)
	totalMinutes := 0

	for i := days - 1; i >= 0; i-- {
		date := time.Now().AddDate(0, 0, -i)
		dateStr := date.Format("2006-01-02")

		minutes := 0
		if prog, ok := userProgress[dateStr]; ok {
			minutes = prog.LearningMinutes
			totalMinutes += minutes
		}

		data.Data = append(data.Data, models.DailyLearningTime{
			Date:    dateStr,
			Minutes: minutes,
		})
	}

	data.TotalMinutes = totalMinutes
	if len(data.Data) > 0 {
		data.AverageMinutes = float64(totalMinutes) / float64(len(data.Data))
	}

	return data, nil
}

func (r *InMemoryStatsRepository) GetProgressData(ctx context.Context, userID uuid.UUID, period string) (*models.ProgressData, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	data := &models.ProgressData{
		Period:  period,
		Words:   []models.TimeSeriesData{},
		Phrases: []models.TimeSeriesData{},
		Pages:   []models.TimeSeriesData{},
	}

	userProgress, exists := r.progress[userID.String()]
	if !exists {
		return data, nil
	}

	days := getDaysForPeriod(period)
	cumulativeWords := 0
	cumulativePhrases := 0
	cumulativePages := 0

	for i := days - 1; i >= 0; i-- {
		date := time.Now().AddDate(0, 0, -i)
		dateStr := date.Format("2006-01-02")

		if prog, ok := userProgress[dateStr]; ok {
			cumulativeWords += prog.MasteredWords
			cumulativePhrases += prog.MasteredPhrases
			cumulativePages += prog.CompletedPages
		}

		data.Words = append(data.Words, models.TimeSeriesData{
			Date:  dateStr,
			Count: cumulativeWords,
		})
		data.Phrases = append(data.Phrases, models.TimeSeriesData{
			Date:  dateStr,
			Count: cumulativePhrases,
		})
		data.Pages = append(data.Pages, models.TimeSeriesData{
			Date:  dateStr,
			Count: cumulativePages,
		})
	}

	return data, nil
}

func (r *InMemoryStatsRepository) GetWeakPoints(ctx context.Context, userID uuid.UUID, limit int) (*models.WeakPointsData, error) {
	// TODO: 実装（STT/発音データが必要）
	return &models.WeakPointsData{
		WeakWords:   []models.WeakItem{},
		WeakPhrases: []models.WeakItem{},
	}, nil
}

func (r *InMemoryStatsRepository) RecordLearningSession(ctx context.Context, session *models.LearningSession) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.sessions[session.ID.String()] = session
	return nil
}

func (r *InMemoryStatsRepository) UpdateUserProgress(ctx context.Context, progress *models.UserProgressDaily) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	userIDStr := progress.UserID.String()
	if _, exists := r.progress[userIDStr]; !exists {
		r.progress[userIDStr] = make(map[string]*models.UserProgressDaily)
	}

	dateStr := progress.Date.Format("2006-01-02")
	r.progress[userIDStr][dateStr] = progress

	return nil
}

func (r *InMemoryStatsRepository) UpdateStreak(ctx context.Context, userID uuid.UUID, activityDate time.Time) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	userIDStr := userID.String()
	streak, exists := r.streaks[userIDStr]

	if !exists {
		streak = &models.LearningStreakRecord{
			ID:               uuid.New(),
			UserID:           userID,
			CurrentStreak:    1,
			LongestStreak:    1,
			LastActivityDate: activityDate,
			UpdatedAt:        time.Now(),
		}
		r.streaks[userIDStr] = streak
		return nil
	}

	// ストリーク計算ロジック
	daysDiff := int(activityDate.Sub(streak.LastActivityDate).Hours() / 24)

	if daysDiff == 1 {
		// 連続
		streak.CurrentStreak++
		if streak.CurrentStreak > streak.LongestStreak {
			streak.LongestStreak = streak.CurrentStreak
		}
	} else if daysDiff > 1 {
		// 途切れた
		streak.CurrentStreak = 1
	}
	// daysDiff == 0 なら同じ日なので何もしない

	streak.LastActivityDate = activityDate
	streak.UpdatedAt = time.Now()

	return nil
}

// ヘルパー関数
func isToday(date time.Time) bool {
	now := time.Now()
	return date.Year() == now.Year() && date.YearDay() == now.YearDay()
}

func isThisWeek(date time.Time) bool {
	now := time.Now()
	_, week := now.ISOWeek()
	_, dateWeek := date.ISOWeek()
	return week == dateWeek && now.Year() == date.Year()
}

func getDaysForPeriod(period string) int {
	switch period {
	case "day":
		return 1
	case "week":
		return 7
	case "month":
		return 30
	case "year":
		return 365
	default:
		return 7
	}
}
