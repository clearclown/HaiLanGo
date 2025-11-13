package stats

import (
	"context"
	"sort"
	"time"

	"github.com/clearclown/HaiLanGo/internal/models"
	"github.com/google/uuid"
)

// StatsRepository defines the interface for statistics data access
type StatsRepository interface {
	GetLearningSessions(ctx context.Context, userID uuid.UUID, startDate, endDate time.Time) ([]models.LearningSession, error)
	GetVocabularyProgress(ctx context.Context, userID uuid.UUID) ([]models.VocabularyProgress, error)
	GetPhraseProgress(ctx context.Context, userID uuid.UUID) ([]models.PhraseProgress, error)
	GetPronunciationScores(ctx context.Context, userID uuid.UUID, startDate, endDate time.Time) ([]models.PronunciationScore, error)
	GetCurrentStreak(ctx context.Context, userID uuid.UUID) (int, error)
	GetLongestStreak(ctx context.Context, userID uuid.UUID) (int, error)
	GetCompletedPagesCount(ctx context.Context, userID uuid.UUID) (int, error)
	GetCompletedBooksCount(ctx context.Context, userID uuid.UUID) (int, error)
}

// Service provides statistics-related business logic
type Service struct {
	repo StatsRepository
}

// NewService creates a new stats service
func NewService(repo StatsRepository) *Service {
	return &Service{
		repo: repo,
	}
}

// GetDashboardStats retrieves all statistics for the dashboard
func (s *Service) GetDashboardStats(ctx context.Context, userID uuid.UUID) (*models.DashboardStats, error) {
	// Get learning time stats
	learningTime, err := s.GetLearningTimeStats(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Get progress stats
	progress, err := s.GetProgressStats(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Get streak stats
	streak, err := s.GetStreakStats(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Get pronunciation average
	pronunciationAvg, err := s.GetPronunciationAverage(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Get weak words
	weakWords, err := s.GetWeakWords(ctx, userID, 10)
	if err != nil {
		return nil, err
	}

	// Get learning time chart (last 7 days)
	learningTimeChart, err := s.GetLearningTimeChart(ctx, userID, 7)
	if err != nil {
		return nil, err
	}

	// Get progress chart (last 30 days)
	progressChart, err := s.GetProgressChart(ctx, userID, 30)
	if err != nil {
		return nil, err
	}

	return &models.DashboardStats{
		LearningTime:      *learningTime,
		Progress:          *progress,
		Streak:            *streak,
		PronunciationAvg:  pronunciationAvg,
		WeakWords:         weakWords,
		LearningTimeChart: learningTimeChart,
		ProgressChart:     progressChart,
	}, nil
}

// GetLearningTimeStats calculates learning time statistics
func (s *Service) GetLearningTimeStats(ctx context.Context, userID uuid.UUID) (*models.LearningTimeStats, error) {
	// Get all learning sessions (from the beginning of time to now)
	sessions, err := s.repo.GetLearningSessions(ctx, userID, time.Time{}, time.Now())
	if err != nil {
		return nil, err
	}

	totalSeconds := 0
	for _, session := range sessions {
		totalSeconds += session.Duration
	}

	totalHours := float64(totalSeconds) / 3600.0

	// Calculate daily average (last 30 days)
	thirtyDaysAgo := time.Now().AddDate(0, 0, -30)
	recentSessions, err := s.repo.GetLearningSessions(ctx, userID, thirtyDaysAgo, time.Now())
	if err != nil {
		return nil, err
	}

	recentSeconds := 0
	for _, session := range recentSessions {
		recentSeconds += session.Duration
	}

	dailyAverage := 0
	if len(recentSessions) > 0 {
		// Calculate unique days
		dayMap := make(map[string]bool)
		for _, session := range recentSessions {
			dateKey := session.StartTime.Format("2006-01-02")
			dayMap[dateKey] = true
		}
		daysCount := len(dayMap)
		if daysCount > 0 {
			dailyAverage = recentSeconds / daysCount
		}
	}

	// Calculate weekly average (last 12 weeks)
	twelveWeeksAgo := time.Now().AddDate(0, 0, -84) // 12 * 7 days
	weeklySessions, err := s.repo.GetLearningSessions(ctx, userID, twelveWeeksAgo, time.Now())
	if err != nil {
		return nil, err
	}

	weeklySeconds := 0
	for _, session := range weeklySessions {
		weeklySeconds += session.Duration
	}

	weeklyAverage := 0
	if len(weeklySessions) > 0 {
		// Calculate unique weeks
		weekMap := make(map[string]bool)
		for _, session := range weeklySessions {
			year, week := session.StartTime.ISOWeek()
			weekKey := string(rune(year)) + "-" + string(rune(week))
			weekMap[weekKey] = true
		}
		weeksCount := len(weekMap)
		if weeksCount > 0 {
			weeklyAverage = weeklySeconds / weeksCount
		}
	}

	// Calculate monthly average (last 12 months)
	twelveMonthsAgo := time.Now().AddDate(0, -12, 0)
	monthlySessions, err := s.repo.GetLearningSessions(ctx, userID, twelveMonthsAgo, time.Now())
	if err != nil {
		return nil, err
	}

	monthlySeconds := 0
	for _, session := range monthlySessions {
		monthlySeconds += session.Duration
	}

	monthlyAverage := 0
	if len(monthlySessions) > 0 {
		// Calculate unique months
		monthMap := make(map[string]bool)
		for _, session := range monthlySessions {
			monthKey := session.StartTime.Format("2006-01")
			monthMap[monthKey] = true
		}
		monthsCount := len(monthMap)
		if monthsCount > 0 {
			monthlyAverage = monthlySeconds / monthsCount
		}
	}

	return &models.LearningTimeStats{
		TotalSeconds:   totalSeconds,
		TotalHours:     totalHours,
		DailyAverage:   dailyAverage,
		WeeklyAverage:  weeklyAverage,
		MonthlyAverage: monthlyAverage,
	}, nil
}

// GetProgressStats retrieves progress statistics
func (s *Service) GetProgressStats(ctx context.Context, userID uuid.UUID) (*models.ProgressStats, error) {
	// Get completed pages count
	completedPages, err := s.repo.GetCompletedPagesCount(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Get vocabulary progress
	vocabProgress, err := s.repo.GetVocabularyProgress(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Count mastered words (mastery level >= 80)
	masteredWords := 0
	for _, vocab := range vocabProgress {
		if vocab.MasteryLevel >= 80 {
			masteredWords++
		}
	}

	// Get phrase progress
	phraseProgress, err := s.repo.GetPhraseProgress(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Count mastered phrases (mastery level >= 80)
	masteredPhrases := 0
	for _, phrase := range phraseProgress {
		if phrase.MasteryLevel >= 80 {
			masteredPhrases++
		}
	}

	// Get completed books count
	completedBooks, err := s.repo.GetCompletedBooksCount(ctx, userID)
	if err != nil {
		return nil, err
	}

	return &models.ProgressStats{
		CompletedPages:  completedPages,
		MasteredWords:   masteredWords,
		MasteredPhrases: masteredPhrases,
		CompletedBooks:  completedBooks,
	}, nil
}

// GetStreakStats retrieves streak statistics
func (s *Service) GetStreakStats(ctx context.Context, userID uuid.UUID) (*models.StreakStats, error) {
	currentStreak, err := s.repo.GetCurrentStreak(ctx, userID)
	if err != nil {
		return nil, err
	}

	longestStreak, err := s.repo.GetLongestStreak(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Get the last study date
	sessions, err := s.repo.GetLearningSessions(ctx, userID, time.Now().AddDate(0, 0, -7), time.Now())
	if err != nil {
		return nil, err
	}

	lastStudyDate := time.Time{}
	if len(sessions) > 0 {
		// Find the most recent session
		for _, session := range sessions {
			if session.StartTime.After(lastStudyDate) {
				lastStudyDate = session.StartTime
			}
		}
	}

	return &models.StreakStats{
		CurrentStreak: currentStreak,
		LongestStreak: longestStreak,
		LastStudyDate: lastStudyDate,
	}, nil
}

// GetPronunciationAverage calculates the average pronunciation score
func (s *Service) GetPronunciationAverage(ctx context.Context, userID uuid.UUID) (float64, error) {
	// Get pronunciation scores (last 30 days)
	thirtyDaysAgo := time.Now().AddDate(0, 0, -30)
	scores, err := s.repo.GetPronunciationScores(ctx, userID, thirtyDaysAgo, time.Now())
	if err != nil {
		return 0, err
	}

	if len(scores) == 0 {
		return 0, nil
	}

	totalScore := 0.0
	for _, score := range scores {
		totalScore += score.Score
	}

	return totalScore / float64(len(scores)), nil
}

// GetWeakWords retrieves weak words (mastery level < 50)
func (s *Service) GetWeakWords(ctx context.Context, userID uuid.UUID, limit int) ([]string, error) {
	vocabProgress, err := s.repo.GetVocabularyProgress(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Filter weak words (mastery level < 50)
	var weakWords []models.VocabularyProgress
	for _, vocab := range vocabProgress {
		if vocab.MasteryLevel < 50 {
			weakWords = append(weakWords, vocab)
		}
	}

	// Sort by mastery level (ascending)
	sort.Slice(weakWords, func(i, j int) bool {
		return weakWords[i].MasteryLevel < weakWords[j].MasteryLevel
	})

	// Limit the results
	if len(weakWords) > limit {
		weakWords = weakWords[:limit]
	}

	// Extract word strings
	result := make([]string, len(weakWords))
	for i, vocab := range weakWords {
		result[i] = vocab.Word
	}

	return result, nil
}

// GetLearningTimeChart retrieves learning time chart data
func (s *Service) GetLearningTimeChart(ctx context.Context, userID uuid.UUID, days int) ([]models.LearningTimeDataPoint, error) {
	startDate := time.Now().AddDate(0, 0, -days)
	sessions, err := s.repo.GetLearningSessions(ctx, userID, startDate, time.Now())
	if err != nil {
		return nil, err
	}

	// Group sessions by date
	dateMap := make(map[string]int)
	for _, session := range sessions {
		dateKey := session.StartTime.Format("2006-01-02")
		dateMap[dateKey] += session.Duration
	}

	// Create data points for all days in the range
	var dataPoints []models.LearningTimeDataPoint
	currentDate := startDate
	for i := 0; i < days; i++ {
		dateKey := currentDate.Format("2006-01-02")
		seconds := dateMap[dateKey]
		dataPoints = append(dataPoints, models.LearningTimeDataPoint{
			Date:    currentDate,
			Seconds: seconds,
		})
		currentDate = currentDate.AddDate(0, 0, 1)
	}

	return dataPoints, nil
}

// GetProgressChart retrieves progress chart data
func (s *Service) GetProgressChart(ctx context.Context, userID uuid.UUID, days int) ([]models.ProgressDataPoint, error) {
	startDate := time.Now().AddDate(0, 0, -days)

	// Get all vocabulary progress
	vocabProgress, err := s.repo.GetVocabularyProgress(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Get all phrase progress
	phraseProgress, err := s.repo.GetPhraseProgress(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Get all sessions
	sessions, err := s.repo.GetLearningSessions(ctx, userID, startDate, time.Now())
	if err != nil {
		return nil, err
	}

	// Group by date
	vocabMap := make(map[string]int)
	for _, vocab := range vocabProgress {
		if vocab.CreatedAt.After(startDate) {
			dateKey := vocab.CreatedAt.Format("2006-01-02")
			vocabMap[dateKey]++
		}
	}

	phraseMap := make(map[string]int)
	for _, phrase := range phraseProgress {
		if phrase.CreatedAt.After(startDate) {
			dateKey := phrase.CreatedAt.Format("2006-01-02")
			phraseMap[dateKey]++
		}
	}

	pageMap := make(map[string]int)
	for _, session := range sessions {
		dateKey := session.StartTime.Format("2006-01-02")
		pageMap[dateKey]++
	}

	// Create cumulative data points
	var dataPoints []models.ProgressDataPoint
	currentDate := startDate
	cumulativeWords := 0
	cumulativePhrases := 0
	cumulativePages := 0

	for i := 0; i < days; i++ {
		dateKey := currentDate.Format("2006-01-02")
		cumulativeWords += vocabMap[dateKey]
		cumulativePhrases += phraseMap[dateKey]
		cumulativePages += pageMap[dateKey]

		dataPoints = append(dataPoints, models.ProgressDataPoint{
			Date:    currentDate,
			Words:   cumulativeWords,
			Phrases: cumulativePhrases,
			Pages:   cumulativePages,
		})
		currentDate = currentDate.AddDate(0, 0, 1)
	}

	return dataPoints, nil
}
