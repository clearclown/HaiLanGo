package learning

import (
	"context"
	"testing"
	"time"

	"github.com/clearclown/HaiLanGo/backend/internal/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockRepository はリポジトリのモック
type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) GetPage(ctx context.Context, bookID uuid.UUID, pageNumber int) (*models.PageWithProgress, error) {
	args := m.Called(ctx, bookID, pageNumber)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.PageWithProgress), args.Error(1)
}

func (m *MockRepository) MarkPageCompleted(ctx context.Context, userID, bookID uuid.UUID, pageNumber int, studyTime int) error {
	args := m.Called(ctx, userID, bookID, pageNumber, studyTime)
	return args.Error(0)
}

func (m *MockRepository) GetProgress(ctx context.Context, userID, bookID uuid.UUID) (*models.LearningProgress, error) {
	args := m.Called(ctx, userID, bookID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.LearningProgress), args.Error(1)
}

func TestGetPage(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockRepository)
	service := NewService(mockRepo)

	bookID := uuid.New()
	pageNumber := 1

	expectedPage := &models.PageWithProgress{
		Page: models.Page{
			ID:          uuid.New(),
			BookID:      bookID,
			PageNumber:  pageNumber,
			ImageURL:    "https://example.com/page1.png",
			OCRText:     "Здравствуйте!",
			Translation: "こんにちは！",
			AudioURL:    "https://example.com/audio1.mp3",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		IsCompleted: false,
		CompletedAt: nil,
	}

	mockRepo.On("GetPage", ctx, bookID, pageNumber).Return(expectedPage, nil)

	// テスト実行
	page, err := service.GetPage(ctx, bookID, pageNumber)

	// 検証
	assert.NoError(t, err)
	assert.NotNil(t, page)
	assert.Equal(t, bookID, page.BookID)
	assert.Equal(t, pageNumber, page.PageNumber)
	assert.Equal(t, "Здравствуйте!", page.OCRText)
	mockRepo.AssertExpectations(t)
}

func TestMarkPageCompleted(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockRepository)
	service := NewService(mockRepo)

	userID := uuid.New()
	bookID := uuid.New()
	pageNumber := 1
	studyTime := 300 // 5分

	mockRepo.On("MarkPageCompleted", ctx, userID, bookID, pageNumber, studyTime).Return(nil)

	// テスト実行
	err := service.MarkPageCompleted(ctx, userID, bookID, pageNumber, studyTime)

	// 検証
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestGetProgress(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockRepository)
	service := NewService(mockRepo)

	userID := uuid.New()
	bookID := uuid.New()

	expectedProgress := &models.LearningProgress{
		BookID:         bookID,
		TotalPages:     150,
		CompletedPages: 45,
		Progress:       30.0,
		TotalStudyTime: 6750, // 1時間52分30秒
	}

	mockRepo.On("GetProgress", ctx, userID, bookID).Return(expectedProgress, nil)

	// テスト実行
	progress, err := service.GetProgress(ctx, userID, bookID)

	// 検証
	assert.NoError(t, err)
	assert.NotNil(t, progress)
	assert.Equal(t, 150, progress.TotalPages)
	assert.Equal(t, 45, progress.CompletedPages)
	assert.Equal(t, 30.0, progress.Progress)
	mockRepo.AssertExpectations(t)
}

func TestGetPage_NotFound(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockRepository)
	service := NewService(mockRepo)

	bookID := uuid.New()
	pageNumber := 999

	mockRepo.On("GetPage", ctx, bookID, pageNumber).Return(nil, ErrPageNotFound)

	// テスト実行
	page, err := service.GetPage(ctx, bookID, pageNumber)

	// 検証
	assert.Error(t, err)
	assert.Nil(t, page)
	assert.Equal(t, ErrPageNotFound, err)
	mockRepo.AssertExpectations(t)
}
