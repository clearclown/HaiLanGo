package learning

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/clearclown/HaiLanGo/backend/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockService はサービスのモック
type MockService struct {
	mock.Mock
}

func (m *MockService) GetPage(ctx context.Context, bookID uuid.UUID, pageNumber int) (*models.PageWithProgress, error) {
	args := m.Called(ctx, bookID, pageNumber)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.PageWithProgress), args.Error(1)
}

func (m *MockService) MarkPageCompleted(ctx context.Context, userID, bookID uuid.UUID, pageNumber int, studyTime int) error {
	args := m.Called(ctx, userID, bookID, pageNumber, studyTime)
	return args.Error(0)
}

func (m *MockService) GetProgress(ctx context.Context, userID, bookID uuid.UUID) (*models.LearningProgress, error) {
	args := m.Called(ctx, userID, bookID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.LearningProgress), args.Error(1)
}

func setupRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	return gin.Default()
}

func TestGetPageHandler(t *testing.T) {
	mockService := new(MockService)
	handler := NewHandler(mockService)

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

	mockService.On("GetPage", mock.Anything, bookID, pageNumber).Return(expectedPage, nil)

	// テスト用ルーター設定
	router := setupRouter()
	router.GET("/api/v1/books/:bookId/pages/:pageNumber", handler.GetPage)

	// リクエスト作成
	req := httptest.NewRequest("GET", "/api/v1/books/"+bookID.String()+"/pages/1", nil)
	w := httptest.NewRecorder()

	// リクエスト実行
	router.ServeHTTP(w, req)

	// 検証
	assert.Equal(t, http.StatusOK, w.Code)

	var response models.PageWithProgress
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, bookID, response.BookID)
	assert.Equal(t, pageNumber, response.PageNumber)
	assert.Equal(t, "Здравствуйте!", response.OCRText)
	mockService.AssertExpectations(t)
}

func TestMarkPageCompletedHandler(t *testing.T) {
	mockService := new(MockService)
	handler := NewHandler(mockService)

	userID := uuid.New()
	bookID := uuid.New()
	pageNumber := 1
	studyTime := 300

	mockService.On("MarkPageCompleted", mock.Anything, userID, bookID, pageNumber, studyTime).Return(nil)

	// テスト用ルーター設定
	router := setupRouter()
	router.POST("/api/v1/books/:bookId/pages/:pageNumber/complete", handler.MarkPageCompleted)

	// リクエストボディ作成
	requestBody := map[string]interface{}{
		"userId":    userID.String(),
		"studyTime": studyTime,
	}
	jsonBody, _ := json.Marshal(requestBody)

	// リクエスト作成
	req := httptest.NewRequest("POST", "/api/v1/books/"+bookID.String()+"/pages/1/complete", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// リクエスト実行
	router.ServeHTTP(w, req)

	// 検証
	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}

func TestGetProgressHandler(t *testing.T) {
	mockService := new(MockService)
	handler := NewHandler(mockService)

	userID := uuid.New()
	bookID := uuid.New()

	expectedProgress := &models.LearningProgress{
		BookID:         bookID,
		TotalPages:     150,
		CompletedPages: 45,
		Progress:       30.0,
		TotalStudyTime: 6750,
	}

	mockService.On("GetProgress", mock.Anything, userID, bookID).Return(expectedProgress, nil)

	// テスト用ルーター設定
	router := setupRouter()
	router.GET("/api/v1/books/:bookId/progress", handler.GetProgress)

	// リクエスト作成
	req := httptest.NewRequest("GET", "/api/v1/books/"+bookID.String()+"/progress?userId="+userID.String(), nil)
	w := httptest.NewRecorder()

	// リクエスト実行
	router.ServeHTTP(w, req)

	// 検証
	assert.Equal(t, http.StatusOK, w.Code)

	var response models.LearningProgress
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, bookID, response.BookID)
	assert.Equal(t, 150, response.TotalPages)
	assert.Equal(t, 45, response.CompletedPages)
	assert.Equal(t, 30.0, response.Progress)
	mockService.AssertExpectations(t)
}

func TestGetPageHandler_InvalidBookID(t *testing.T) {
	mockService := new(MockService)
	handler := NewHandler(mockService)

	router := setupRouter()
	router.GET("/api/v1/books/:bookId/pages/:pageNumber", handler.GetPage)

	// 不正なbookIDでリクエスト
	req := httptest.NewRequest("GET", "/api/v1/books/invalid-uuid/pages/1", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	// 検証
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestGetPageHandler_InvalidPageNumber(t *testing.T) {
	mockService := new(MockService)
	handler := NewHandler(mockService)

	bookID := uuid.New()

	router := setupRouter()
	router.GET("/api/v1/books/:bookId/pages/:pageNumber", handler.GetPage)

	// 不正なpageNumberでリクエスト
	req := httptest.NewRequest("GET", "/api/v1/books/"+bookID.String()+"/pages/abc", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	// 検証
	assert.Equal(t, http.StatusBadRequest, w.Code)
}
