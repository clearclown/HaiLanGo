package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/clearclown/HaiLanGo/backend/internal/models"
	"github.com/clearclown/HaiLanGo/backend/internal/repository"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setupLearningTestRouter() (*gin.Engine, *repository.InMemoryLearningRepository) {
	gin.SetMode(gin.TestMode)

	repo := repository.NewInMemoryLearningRepository()
	handler := NewLearningHandler(repo)

	r := gin.New()
	// テスト用にuser_idをセット
	r.Use(func(c *gin.Context) {
		c.Set("user_id", "550e8400-e29b-41d4-a716-446655440000")
		c.Next()
	})

	handler.RegisterRoutes(r.Group("/api/v1"))

	return r, repo
}

func TestGetPageLearning(t *testing.T) {
	router, _ := setupLearningTestRouter()

	testBookID := "550e8400-e29b-41d4-a716-446655440000"
	testPageNumber := 12

	req, _ := http.NewRequest(http.MethodGet, "/api/v1/learning/books/"+testBookID+"/pages/12", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var pageLearning models.PageLearning
	err := json.Unmarshal(w.Body.Bytes(), &pageLearning)
	assert.NoError(t, err)

	// ページデータの確認
	assert.Equal(t, testBookID, pageLearning.Page.BookID)
	assert.Equal(t, testPageNumber, pageLearning.Page.PageNumber)
	assert.NotEmpty(t, pageLearning.Page.OCRText)
	assert.NotEmpty(t, pageLearning.Page.Translation)

	// フレーズデータの確認
	assert.GreaterOrEqual(t, len(pageLearning.Phrases), 0)

	// ナビゲーション情報の確認
	assert.Equal(t, testPageNumber, pageLearning.Navigation.CurrentPage)
	assert.True(t, pageLearning.Navigation.HasPrevious)
	assert.True(t, pageLearning.Navigation.HasNext)
	assert.Equal(t, 150, pageLearning.Navigation.TotalPages)
}

func TestGetPageLearningFirstPage(t *testing.T) {
	router, _ := setupLearningTestRouter()

	testBookID := "550e8400-e29b-41d4-a716-446655440000"

	req, _ := http.NewRequest(http.MethodGet, "/api/v1/learning/books/"+testBookID+"/pages/1", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var pageLearning models.PageLearning
	err := json.Unmarshal(w.Body.Bytes(), &pageLearning)
	assert.NoError(t, err)

	// 最初のページなのでHasPreviousはfalse
	assert.False(t, pageLearning.Navigation.HasPrevious)
	assert.True(t, pageLearning.Navigation.HasNext)
}

func TestGetPageLearningLastPage(t *testing.T) {
	router, _ := setupLearningTestRouter()

	testBookID := "550e8400-e29b-41d4-a716-446655440000"

	req, _ := http.NewRequest(http.MethodGet, "/api/v1/learning/books/"+testBookID+"/pages/150", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var pageLearning models.PageLearning
	err := json.Unmarshal(w.Body.Bytes(), &pageLearning)
	assert.NoError(t, err)

	// 最後のページなのでHasNextはfalse
	assert.True(t, pageLearning.Navigation.HasPrevious)
	assert.False(t, pageLearning.Navigation.HasNext)
}

func TestCompletePage(t *testing.T) {
	router, _ := setupLearningTestRouter()

	testBookID := "550e8400-e29b-41d4-a716-446655440000"

	reqBody := models.CompletePageRequest{
		StudyTime: 180,
		Notes:     "難しかったが理解できた",
	}
	bodyBytes, _ := json.Marshal(reqBody)

	req, _ := http.NewRequest(http.MethodPost, "/api/v1/learning/books/"+testBookID+"/pages/50/complete", bytes.NewBuffer(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response models.CompletePageResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Equal(t, "Page marked as completed", response.Message)
	assert.True(t, response.Progress.IsCompleted)
	assert.NotNil(t, response.Progress.CompletedAt)
	assert.Equal(t, 51, response.NextPage)
}

func TestCompletePageInvalidStudyTime(t *testing.T) {
	router, _ := setupLearningTestRouter()

	testBookID := "550e8400-e29b-41d4-a716-446655440000"

	reqBody := models.CompletePageRequest{
		StudyTime: 0, // Invalid
		Notes:     "",
	}
	bodyBytes, _ := json.Marshal(reqBody)

	req, _ := http.NewRequest(http.MethodPost, "/api/v1/learning/books/"+testBookID+"/pages/50/complete", bytes.NewBuffer(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestRecordSessionStart(t *testing.T) {
	router, _ := setupLearningTestRouter()

	testBookID := "550e8400-e29b-41d4-a716-446655440000"

	reqBody := models.SessionRequest{
		Action:    "start",
		Timestamp: time.Now(),
	}
	bodyBytes, _ := json.Marshal(reqBody)

	req, _ := http.NewRequest(http.MethodPost, "/api/v1/learning/books/"+testBookID+"/pages/12/session", bytes.NewBuffer(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response models.SessionResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.NotEmpty(t, response.SessionID)
	assert.NotZero(t, response.StartedAt)
	assert.Nil(t, response.EndedAt)
}

func TestRecordSessionEnd(t *testing.T) {
	router, _ := setupLearningTestRouter()

	testBookID := "550e8400-e29b-41d4-a716-446655440000"

	reqBody := models.SessionRequest{
		Action:    "end",
		Timestamp: time.Now(),
	}
	bodyBytes, _ := json.Marshal(reqBody)

	req, _ := http.NewRequest(http.MethodPost, "/api/v1/learning/books/"+testBookID+"/pages/12/session", bytes.NewBuffer(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response models.SessionResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.NotEmpty(t, response.SessionID)
	assert.NotNil(t, response.EndedAt)
}

func TestGetBookProgress(t *testing.T) {
	router, _ := setupLearningTestRouter()

	testBookID := "550e8400-e29b-41d4-a716-446655440000"

	req, _ := http.NewRequest(http.MethodGet, "/api/v1/learning/books/"+testBookID+"/progress", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var progress models.BookProgressSummary
	err := json.Unmarshal(w.Body.Bytes(), &progress)
	assert.NoError(t, err)

	assert.Equal(t, testBookID, progress.BookID)
	assert.Equal(t, 150, progress.TotalPages)
	assert.Equal(t, 45, progress.CompletedPages) // サンプルデータ
	assert.Greater(t, progress.CompletionPercentage, 0.0)
	assert.Greater(t, progress.TotalStudyTime, 0)
	assert.Len(t, progress.Pages, 150)
}

func TestLearningUnauthorized(t *testing.T) {
	gin.SetMode(gin.TestMode)

	repo := repository.NewInMemoryLearningRepository()
	handler := NewLearningHandler(repo)

	r := gin.New()
	// user_idをセットしない
	handler.RegisterRoutes(r.Group("/api/v1"))

	testBookID := "550e8400-e29b-41d4-a716-446655440000"

	tests := []struct {
		name   string
		method string
		path   string
	}{
		{"Get Page", "GET", "/api/v1/learning/books/" + testBookID + "/pages/12"},
		{"Complete Page", "POST", "/api/v1/learning/books/" + testBookID + "/pages/12/complete"},
		{"Record Session", "POST", "/api/v1/learning/books/" + testBookID + "/pages/12/session"},
		{"Get Progress", "GET", "/api/v1/learning/books/" + testBookID + "/progress"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest(tt.method, tt.path, nil)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			assert.Equal(t, http.StatusUnauthorized, w.Code)
		})
	}
}

func TestInvalidBookID(t *testing.T) {
	router, _ := setupLearningTestRouter()

	req, _ := http.NewRequest(http.MethodGet, "/api/v1/learning/books/invalid-uuid/pages/12", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestInvalidPageNumber(t *testing.T) {
	router, _ := setupLearningTestRouter()

	testBookID := "550e8400-e29b-41d4-a716-446655440000"

	req, _ := http.NewRequest(http.MethodGet, "/api/v1/learning/books/"+testBookID+"/pages/0", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}
