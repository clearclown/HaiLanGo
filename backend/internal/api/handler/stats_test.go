package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/clearclown/HaiLanGo/backend/internal/models"
	"github.com/clearclown/HaiLanGo/backend/internal/repository"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setupStatsTestRouter() (*gin.Engine, *repository.InMemoryStatsRepository) {
	gin.SetMode(gin.TestMode)

	repo := repository.NewInMemoryStatsRepository()
	handler := NewStatsHandler(repo)

	r := gin.New()
	// テスト用にuser_idをセット
	r.Use(func(c *gin.Context) {
		c.Set("user_id", "550e8400-e29b-41d4-a716-446655440000")
		c.Next()
	})

	handler.RegisterRoutes(r.Group("/api/v1"))

	return r, repo
}

func TestGetDashboard(t *testing.T) {
	router, _ := setupStatsTestRouter()

	req, _ := http.NewRequest(http.MethodGet, "/api/v1/stats/dashboard", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var stats models.DashboardStatsFlat
	err := json.Unmarshal(w.Body.Bytes(), &stats)
	assert.NoError(t, err)

	// サンプルデータが正しく返されているか確認
	assert.GreaterOrEqual(t, stats.CurrentStreak, 0)
	assert.GreaterOrEqual(t, stats.LongestStreak, 0)
	assert.GreaterOrEqual(t, stats.MasteredWords, 0)
	assert.Equal(t, 7, stats.CurrentStreak) // サンプルデータの期待値
	assert.Equal(t, 15, stats.LongestStreak) // サンプルデータの期待値
}

func TestGetLearningTime(t *testing.T) {
	router, _ := setupStatsTestRouter()

	req, _ := http.NewRequest(http.MethodGet, "/api/v1/stats/learning-time?period=week", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var data models.LearningTimeData
	err := json.Unmarshal(w.Body.Bytes(), &data)
	assert.NoError(t, err)

	assert.Equal(t, "week", data.Period)
	assert.Len(t, data.Data, 7)
	assert.GreaterOrEqual(t, data.TotalMinutes, 0)
}

func TestGetLearningTimeDay(t *testing.T) {
	router, _ := setupStatsTestRouter()

	req, _ := http.NewRequest(http.MethodGet, "/api/v1/stats/learning-time?period=day", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var data models.LearningTimeData
	err := json.Unmarshal(w.Body.Bytes(), &data)
	assert.NoError(t, err)

	assert.Equal(t, "day", data.Period)
	assert.Len(t, data.Data, 1)
}

func TestGetLearningTimeMonth(t *testing.T) {
	router, _ := setupStatsTestRouter()

	req, _ := http.NewRequest(http.MethodGet, "/api/v1/stats/learning-time?period=month", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var data models.LearningTimeData
	err := json.Unmarshal(w.Body.Bytes(), &data)
	assert.NoError(t, err)

	assert.Equal(t, "month", data.Period)
	assert.Len(t, data.Data, 30)
}

func TestGetProgress(t *testing.T) {
	router, _ := setupStatsTestRouter()

	req, _ := http.NewRequest(http.MethodGet, "/api/v1/stats/progress?period=month", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var data models.ProgressData
	err := json.Unmarshal(w.Body.Bytes(), &data)
	assert.NoError(t, err)

	assert.Equal(t, "month", data.Period)
	assert.Len(t, data.Words, 30)
	assert.Len(t, data.Phrases, 30)
	assert.Len(t, data.Pages, 30)
}

func TestGetProgressWeek(t *testing.T) {
	router, _ := setupStatsTestRouter()

	req, _ := http.NewRequest(http.MethodGet, "/api/v1/stats/progress?period=week", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var data models.ProgressData
	err := json.Unmarshal(w.Body.Bytes(), &data)
	assert.NoError(t, err)

	assert.Equal(t, "week", data.Period)
	assert.Len(t, data.Words, 7)
	assert.Len(t, data.Phrases, 7)
	assert.Len(t, data.Pages, 7)
}

func TestGetWeakPoints(t *testing.T) {
	router, _ := setupStatsTestRouter()

	req, _ := http.NewRequest(http.MethodGet, "/api/v1/stats/weak-points?limit=5", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var data models.WeakPointsData
	err := json.Unmarshal(w.Body.Bytes(), &data)
	assert.NoError(t, err)

	// InMemory実装では空のデータが返される
	assert.NotNil(t, data.WeakWords)
	assert.NotNil(t, data.WeakPhrases)
}

func TestGetWeakPointsDefaultLimit(t *testing.T) {
	router, _ := setupStatsTestRouter()

	req, _ := http.NewRequest(http.MethodGet, "/api/v1/stats/weak-points", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var data models.WeakPointsData
	err := json.Unmarshal(w.Body.Bytes(), &data)
	assert.NoError(t, err)

	assert.NotNil(t, data.WeakWords)
	assert.NotNil(t, data.WeakPhrases)
}

func TestStatsUnauthorized(t *testing.T) {
	gin.SetMode(gin.TestMode)

	repo := repository.NewInMemoryStatsRepository()
	handler := NewStatsHandler(repo)

	r := gin.New()
	// user_idをセットしない
	handler.RegisterRoutes(r.Group("/api/v1"))

	tests := []struct {
		name string
		path string
	}{
		{"Dashboard", "/api/v1/stats/dashboard"},
		{"Learning Time", "/api/v1/stats/learning-time"},
		{"Progress", "/api/v1/stats/progress"},
		{"Weak Points", "/api/v1/stats/weak-points"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest(http.MethodGet, tt.path, nil)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			assert.Equal(t, http.StatusUnauthorized, w.Code)
		})
	}
}
