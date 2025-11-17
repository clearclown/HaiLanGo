package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/clearclown/HaiLanGo/backend/internal/models"
	"github.com/clearclown/HaiLanGo/backend/internal/repository"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func setupSTTTestRouter() (*gin.Engine, repository.STTRepositoryInterface) {
	gin.SetMode(gin.TestMode)

	sttRepo := repository.NewInMemorySTTRepository()
	sttHandler := NewSTTHandler(sttRepo)

	r := gin.New()

	// テスト用の認証ミドルウェア
	r.Use(func(c *gin.Context) {
		c.Set("user_id", "550e8400-e29b-41d4-a716-446655440000")
		c.Next()
	})

	// ルート登録
	v1 := r.Group("/api/v1")
	sttHandler.RegisterRoutes(v1)

	return r, sttRepo
}

// TestSTTRecognize は音声認識のテスト
func TestSTTRecognize(t *testing.T) {
	router, _ := setupSTTTestRouter()

	requestBody := models.STTRequest{
		AudioData:     "base64_encoded_audio_data",
		Language:      "ru",
		ReferenceText: "Здравствуйте!",
		Options: models.STTRecognizeOptions{
			Format:            "wav",
			EnablePunctuation: true,
			EnableWordTiming:  true,
			Evaluate:          true,
		},
	}

	body, _ := json.Marshal(requestBody)
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/stt/recognize", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusAccepted, w.Code)

	var response models.STTJobResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.NotEmpty(t, response.JobID)
	assert.Equal(t, models.STTStatusPending, response.Status)
}

// TestSTTGetLanguages はサポート言語一覧取得のテスト
func TestSTTGetLanguages(t *testing.T) {
	router, _ := setupSTTTestRouter()

	req, _ := http.NewRequest(http.MethodGet, "/api/v1/stt/languages", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response []*models.STTLanguage
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	// 12言語がサポートされているはず
	assert.Equal(t, 12, len(response))

	// 日本語が含まれているか確認
	foundJapanese := false
	for _, lang := range response {
		if lang.Code == "ja" {
			foundJapanese = true
			assert.Equal(t, "Japanese", lang.Name)
			assert.Equal(t, "日本語", lang.NativeName)
			assert.True(t, lang.IsSupported)
			assert.True(t, lang.SupportsPronunciation)
			break
		}
	}
	assert.True(t, foundJapanese)
}

// TestSTTGetJobStatus はジョブステータス取得のテスト
func TestSTTGetJobStatus(t *testing.T) {
	router, repo := setupSTTTestRouter()

	// 完了済みジョブを取得
	testBookID := "550e8400-e29b-41d4-a716-446655440000"

	ctx := context.Background()
	bookID, _ := uuid.Parse(testBookID)
	jobs, _ := repo.GetJobsByBookID(ctx, bookID)

	var completedJobID string
	for _, job := range jobs {
		if job.Status == models.STTStatusCompleted {
			completedJobID = job.ID
			break
		}
	}

	assert.NotEmpty(t, completedJobID)

	req, _ := http.NewRequest(http.MethodGet, "/api/v1/stt/jobs/"+completedJobID, nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response models.STTJobResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Equal(t, completedJobID, response.JobID)
	assert.Equal(t, models.STTStatusCompleted, response.Status)
	assert.NotNil(t, response.Result)
	assert.NotNil(t, response.Score)
}

// TestSTTGetBookJobs は書籍のSTTジョブ一覧取得のテスト
func TestSTTGetBookJobs(t *testing.T) {
	router, _ := setupSTTTestRouter()

	testBookID := "550e8400-e29b-41d4-a716-446655440000"

	req, _ := http.NewRequest(http.MethodGet, "/api/v1/stt/books/"+testBookID+"/jobs", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response []*models.STTJobResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	// サンプルデータには30ジョブ存在
	assert.Equal(t, 30, len(response))

	// 最初のジョブはページ1で完了済み
	assert.Equal(t, testBookID, response[0].BookID)
	assert.Equal(t, models.STTStatusCompleted, response[0].Status)
}

// TestSTTGetStatistics は統計情報取得のテスト
func TestSTTGetStatistics(t *testing.T) {
	router, _ := setupSTTTestRouter()

	req, _ := http.NewRequest(http.MethodGet, "/api/v1/stt/statistics", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response models.STTStatistics
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	// サンプルデータには25件の完了ジョブ
	assert.Equal(t, 25, response.TotalRecognitions)
	assert.Equal(t, 25, response.TotalEvaluations)
	assert.Greater(t, response.AverageScore, 0.0)
	assert.Greater(t, response.BestScore, 0)
}

// TestSTTUnauthorized は認証なしのテスト
func TestSTTUnauthorized(t *testing.T) {
	gin.SetMode(gin.TestMode)

	sttRepo := repository.NewInMemorySTTRepository()
	sttHandler := NewSTTHandler(sttRepo)

	r := gin.New()
	v1 := r.Group("/api/v1")
	sttHandler.RegisterRoutes(v1)

	tests := []struct {
		name   string
		method string
		url    string
	}{
		{"Recognize", http.MethodPost, "/api/v1/stt/recognize"},
		{"Get Languages", http.MethodGet, "/api/v1/stt/languages"},
		{"Get Job Status", http.MethodGet, "/api/v1/stt/jobs/test-job-id"},
		{"Get Book Jobs", http.MethodGet, "/api/v1/stt/books/550e8400-e29b-41d4-a716-446655440000/jobs"},
		{"Get Statistics", http.MethodGet, "/api/v1/stt/statistics"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest(tt.method, tt.url, nil)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			assert.Equal(t, http.StatusUnauthorized, w.Code)
		})
	}
}

// TestSTTInvalidBookID は無効な書籍IDのテスト
func TestSTTInvalidBookID(t *testing.T) {
	router, _ := setupSTTTestRouter()

	req, _ := http.NewRequest(http.MethodGet, "/api/v1/stt/books/invalid-uuid/jobs", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// TestSTTInvalidJobID は無効なジョブIDのテスト
func TestSTTInvalidJobID(t *testing.T) {
	router, _ := setupSTTTestRouter()

	req, _ := http.NewRequest(http.MethodGet, "/api/v1/stt/jobs/invalid-job-id", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}
