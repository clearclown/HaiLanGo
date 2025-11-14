package stats

import (
	"context"
	"testing"
	"time"

	"github.com/clearclown/HaiLanGo/backend/internal/models"
	"github.com/google/uuid"
)

// MockStatsRepository is a mock implementation of StatsRepository
type MockStatsRepository struct {
	learningSessions   []models.LearningSession
	vocabularyProgress []models.VocabularyProgress
	phraseProgress     []models.PhraseProgress
	pronunciationScores []models.PronunciationScore
}

func NewMockStatsRepository() *MockStatsRepository {
	return &MockStatsRepository{
		learningSessions:   []models.LearningSession{},
		vocabularyProgress: []models.VocabularyProgress{},
		phraseProgress:     []models.PhraseProgress{},
		pronunciationScores: []models.PronunciationScore{},
	}
}

func (m *MockStatsRepository) GetLearningSessions(ctx context.Context, userID uuid.UUID, startDate, endDate time.Time) ([]models.LearningSession, error) {
	var result []models.LearningSession
	for _, session := range m.learningSessions {
		if session.UserID == userID &&
		   session.StartTime.After(startDate) &&
		   session.StartTime.Before(endDate) {
			result = append(result, session)
		}
	}
	return result, nil
}

func (m *MockStatsRepository) GetVocabularyProgress(ctx context.Context, userID uuid.UUID) ([]models.VocabularyProgress, error) {
	var result []models.VocabularyProgress
	for _, vocab := range m.vocabularyProgress {
		if vocab.UserID == userID {
			result = append(result, vocab)
		}
	}
	return result, nil
}

func (m *MockStatsRepository) GetPhraseProgress(ctx context.Context, userID uuid.UUID) ([]models.PhraseProgress, error) {
	var result []models.PhraseProgress
	for _, phrase := range m.phraseProgress {
		if phrase.UserID == userID {
			result = append(result, phrase)
		}
	}
	return result, nil
}

func (m *MockStatsRepository) GetPronunciationScores(ctx context.Context, userID uuid.UUID, startDate, endDate time.Time) ([]models.PronunciationScore, error) {
	var result []models.PronunciationScore
	for _, score := range m.pronunciationScores {
		if score.UserID == userID &&
		   score.CreatedAt.After(startDate) &&
		   score.CreatedAt.Before(endDate) {
			result = append(result, score)
		}
	}
	return result, nil
}

func (m *MockStatsRepository) GetCurrentStreak(ctx context.Context, userID uuid.UUID) (int, error) {
	return 7, nil // Mock: 7 days streak
}

func (m *MockStatsRepository) GetLongestStreak(ctx context.Context, userID uuid.UUID) (int, error) {
	return 15, nil // Mock: 15 days longest streak
}

func (m *MockStatsRepository) GetCompletedPagesCount(ctx context.Context, userID uuid.UUID) (int, error) {
	return len(m.learningSessions), nil
}

func (m *MockStatsRepository) GetCompletedBooksCount(ctx context.Context, userID uuid.UUID) (int, error) {
	// Count unique book IDs
	bookMap := make(map[uuid.UUID]bool)
	for _, session := range m.learningSessions {
		if session.UserID == userID {
			bookMap[session.BookID] = true
		}
	}
	return len(bookMap), nil
}

// Test cases

func TestGetDashboardStats_NoData(t *testing.T) {
	// Arrange
	repo := NewMockStatsRepository()
	service := NewService(repo)
	userID := uuid.New()
	ctx := context.Background()

	// Act
	stats, err := service.GetDashboardStats(ctx, userID)

	// Assert
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if stats.LearningTime.TotalSeconds != 0 {
		t.Errorf("Expected TotalSeconds to be 0, got %d", stats.LearningTime.TotalSeconds)
	}

	if stats.Progress.CompletedPages != 0 {
		t.Errorf("Expected CompletedPages to be 0, got %d", stats.Progress.CompletedPages)
	}

	if stats.Progress.MasteredWords != 0 {
		t.Errorf("Expected MasteredWords to be 0, got %d", stats.Progress.MasteredWords)
	}
}

func TestGetDashboardStats_WithData(t *testing.T) {
	// Arrange
	repo := NewMockStatsRepository()
	userID := uuid.New()
	bookID := uuid.New()
	pageID := uuid.New()

	// Add mock data
	repo.learningSessions = []models.LearningSession{
		{
			ID:        uuid.New(),
			UserID:    userID,
			BookID:    bookID,
			PageID:    pageID,
			StartTime: time.Now().Add(-2 * time.Hour),
			EndTime:   time.Now().Add(-1 * time.Hour),
			Duration:  3600, // 1 hour
			CreatedAt: time.Now(),
		},
		{
			ID:        uuid.New(),
			UserID:    userID,
			BookID:    bookID,
			PageID:    uuid.New(),
			StartTime: time.Now().Add(-1 * time.Hour),
			EndTime:   time.Now(),
			Duration:  1800, // 30 minutes
			CreatedAt: time.Now(),
		},
	}

	repo.vocabularyProgress = []models.VocabularyProgress{
		{
			ID:           uuid.New(),
			UserID:       userID,
			Word:         "hello",
			Language:     "en",
			MasteryLevel: 85,
			ReviewCount:  5,
			CreatedAt:    time.Now(),
		},
		{
			ID:           uuid.New(),
			UserID:       userID,
			Word:         "world",
			Language:     "en",
			MasteryLevel: 90,
			ReviewCount:  7,
			CreatedAt:    time.Now(),
		},
	}

	repo.phraseProgress = []models.PhraseProgress{
		{
			ID:           uuid.New(),
			UserID:       userID,
			Phrase:       "Good morning",
			Language:     "en",
			MasteryLevel: 80,
			ReviewCount:  3,
			CreatedAt:    time.Now(),
		},
	}

	repo.pronunciationScores = []models.PronunciationScore{
		{
			ID:        uuid.New(),
			UserID:    userID,
			Text:      "hello",
			Language:  "en",
			Score:     85.5,
			Accuracy:  88.0,
			Fluency:   83.0,
			CreatedAt: time.Now(),
		},
		{
			ID:        uuid.New(),
			UserID:    userID,
			Text:      "world",
			Language:  "en",
			Score:     92.0,
			Accuracy:  94.0,
			Fluency:   90.0,
			CreatedAt: time.Now(),
		},
	}

	service := NewService(repo)
	ctx := context.Background()

	// Act
	stats, err := service.GetDashboardStats(ctx, userID)

	// Assert
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Check learning time (3600 + 1800 = 5400 seconds)
	if stats.LearningTime.TotalSeconds != 5400 {
		t.Errorf("Expected TotalSeconds to be 5400, got %d", stats.LearningTime.TotalSeconds)
	}

	// Check total hours (5400 / 3600 = 1.5 hours)
	expectedHours := 1.5
	if stats.LearningTime.TotalHours != expectedHours {
		t.Errorf("Expected TotalHours to be %.2f, got %.2f", expectedHours, stats.LearningTime.TotalHours)
	}

	// Check completed pages (2 sessions)
	if stats.Progress.CompletedPages != 2 {
		t.Errorf("Expected CompletedPages to be 2, got %d", stats.Progress.CompletedPages)
	}

	// Check mastered words (mastery level >= 80)
	if stats.Progress.MasteredWords != 2 {
		t.Errorf("Expected MasteredWords to be 2, got %d", stats.Progress.MasteredWords)
	}

	// Check mastered phrases
	if stats.Progress.MasteredPhrases != 1 {
		t.Errorf("Expected MasteredPhrases to be 1, got %d", stats.Progress.MasteredPhrases)
	}

	// Check streak
	if stats.Streak.CurrentStreak != 7 {
		t.Errorf("Expected CurrentStreak to be 7, got %d", stats.Streak.CurrentStreak)
	}

	if stats.Streak.LongestStreak != 15 {
		t.Errorf("Expected LongestStreak to be 15, got %d", stats.Streak.LongestStreak)
	}

	// Check pronunciation average ((85.5 + 92.0) / 2 = 88.75)
	expectedPronAvg := 88.75
	if stats.PronunciationAvg != expectedPronAvg {
		t.Errorf("Expected PronunciationAvg to be %.2f, got %.2f", expectedPronAvg, stats.PronunciationAvg)
	}
}

func TestGetLearningTimeStats(t *testing.T) {
	// Arrange
	repo := NewMockStatsRepository()
	userID := uuid.New()
	bookID := uuid.New()

	// Add sessions over 7 days
	now := time.Now()
	for i := 0; i < 7; i++ {
		repo.learningSessions = append(repo.learningSessions, models.LearningSession{
			ID:        uuid.New(),
			UserID:    userID,
			BookID:    bookID,
			PageID:    uuid.New(),
			StartTime: now.AddDate(0, 0, -i),
			EndTime:   now.AddDate(0, 0, -i).Add(1 * time.Hour),
			Duration:  3600, // 1 hour each day
			CreatedAt: now.AddDate(0, 0, -i),
		})
	}

	service := NewService(repo)
	ctx := context.Background()

	// Act
	stats, err := service.GetLearningTimeStats(ctx, userID)

	// Assert
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Total: 7 days * 3600 seconds = 25200 seconds
	if stats.TotalSeconds != 25200 {
		t.Errorf("Expected TotalSeconds to be 25200, got %d", stats.TotalSeconds)
	}

	// Total hours: 25200 / 3600 = 7 hours
	expectedHours := 7.0
	if stats.TotalHours != expectedHours {
		t.Errorf("Expected TotalHours to be %.2f, got %.2f", expectedHours, stats.TotalHours)
	}

	// Daily average: 25200 / 7 = 3600 seconds (1 hour)
	if stats.DailyAverage != 3600 {
		t.Errorf("Expected DailyAverage to be 3600, got %d", stats.DailyAverage)
	}
}

func TestGetProgressStats(t *testing.T) {
	// Arrange
	repo := NewMockStatsRepository()
	userID := uuid.New()

	// Add vocabulary progress
	for i := 0; i < 10; i++ {
		masteryLevel := 50 + i*5 // 50, 55, 60, ..., 95
		repo.vocabularyProgress = append(repo.vocabularyProgress, models.VocabularyProgress{
			ID:           uuid.New(),
			UserID:       userID,
			Word:         "word" + string(rune(i)),
			Language:     "en",
			MasteryLevel: masteryLevel,
			CreatedAt:    time.Now(),
		})
	}

	service := NewService(repo)
	ctx := context.Background()

	// Act
	stats, err := service.GetProgressStats(ctx, userID)

	// Assert
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Mastered words: mastery level >= 80 (80, 85, 90, 95) = 4 words
	if stats.MasteredWords != 4 {
		t.Errorf("Expected MasteredWords to be 4, got %d", stats.MasteredWords)
	}
}

func TestGetLearningTimeChart_Last7Days(t *testing.T) {
	// Arrange
	repo := NewMockStatsRepository()
	userID := uuid.New()
	bookID := uuid.New()

	// Add sessions for the last 7 days
	now := time.Now()
	for i := 0; i < 7; i++ {
		date := now.AddDate(0, 0, -i)
		repo.learningSessions = append(repo.learningSessions, models.LearningSession{
			ID:        uuid.New(),
			UserID:    userID,
			BookID:    bookID,
			PageID:    uuid.New(),
			StartTime: date,
			EndTime:   date.Add(time.Duration(i+1) * time.Hour),
			Duration:  (i + 1) * 3600, // 1, 2, 3, ..., 7 hours
			CreatedAt: date,
		})
	}

	service := NewService(repo)
	ctx := context.Background()

	// Act
	chart, err := service.GetLearningTimeChart(ctx, userID, 7)

	// Assert
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Should have 7 data points
	if len(chart) != 7 {
		t.Errorf("Expected 7 data points, got %d", len(chart))
	}

	// Verify data points are in chronological order
	for i := 1; i < len(chart); i++ {
		if chart[i].Date.Before(chart[i-1].Date) {
			t.Errorf("Data points are not in chronological order")
		}
	}
}

func TestGetWeakWords(t *testing.T) {
	// Arrange
	repo := NewMockStatsRepository()
	userID := uuid.New()

	// Add weak words (mastery level < 50)
	weakWords := []string{"difficult", "complex", "challenging"}
	for _, word := range weakWords {
		repo.vocabularyProgress = append(repo.vocabularyProgress, models.VocabularyProgress{
			ID:           uuid.New(),
			UserID:       userID,
			Word:         word,
			Language:     "en",
			MasteryLevel: 30,
			CreatedAt:    time.Now(),
		})
	}

	// Add strong words (mastery level >= 80)
	strongWords := []string{"easy", "simple"}
	for _, word := range strongWords {
		repo.vocabularyProgress = append(repo.vocabularyProgress, models.VocabularyProgress{
			ID:           uuid.New(),
			UserID:       userID,
			Word:         word,
			Language:     "en",
			MasteryLevel: 90,
			CreatedAt:    time.Now(),
		})
	}

	service := NewService(repo)
	ctx := context.Background()

	// Act
	weakWords, err := service.GetWeakWords(ctx, userID, 10)

	// Assert
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Should have 3 weak words
	if len(weakWords) != 3 {
		t.Errorf("Expected 3 weak words, got %d", len(weakWords))
	}
}
