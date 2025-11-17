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

func setupTTSTestRouter() (*gin.Engine, repository.TTSRepositoryInterface) {
	gin.SetMode(gin.TestMode)

	ttsRepo := repository.NewInMemoryTTSRepository()
	ttsHandler := NewTTSHandler(ttsRepo)

	r := gin.New()

	// テスト用の認証ミドルウェア
	r.Use(func(c *gin.Context) {
		c.Set("user_id", "550e8400-e29b-41d4-a716-446655440000")
		c.Next()
	})

	// ルート登録
	v1 := r.Group("/api/v1")
	ttsHandler.RegisterRoutes(v1)

	return r, ttsRepo
}

// TestSynthesize は音声合成のテスト
func TestTTSSynthesize(t *testing.T) {
	router, _ := setupTTSTestRouter()

	requestBody := models.TTSRequest{
		Text:     "こんにちは",
		Language: "ja",
		Options: models.TTSSynthesizeOptions{
			Speed:   1.0,
			Quality: models.TTSQualityStandard,
			Format:  models.TTSAudioFormatMP3,
		},
	}

	body, _ := json.Marshal(requestBody)
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/tts/synthesize", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusAccepted, w.Code)

	var response models.TTSJobResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.NotEmpty(t, response.JobID)
	assert.Equal(t, models.TTSStatusPending, response.Status)
}

// TestGetAudio は音声取得のテスト
func TestTTSGetAudio(t *testing.T) {
	router, repo := setupTTSTestRouter()

	// 完了済みジョブを取得
	testBookID := "550e8400-e29b-41d4-a716-446655440000"

	ctx := context.Background()
	bookID, _ := uuid.Parse(testBookID)
	jobs, _ := repo.GetJobsByBookID(ctx, bookID)

	var audioID string
	for _, job := range jobs {
		if job.Status == models.TTSStatusCompleted && job.AudioID != "" {
			audioID = job.AudioID
			break
		}
	}

	assert.NotEmpty(t, audioID)

	req, _ := http.NewRequest(http.MethodGet, "/api/v1/tts/audio/"+audioID, nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// リダイレクトが返されるはず
	assert.Equal(t, http.StatusFound, w.Code)
}

// TestGetLanguages はサポート言語一覧取得のテスト
func TestTTSGetLanguages(t *testing.T) {
	router, _ := setupTTSTestRouter()

	req, _ := http.NewRequest(http.MethodGet, "/api/v1/tts/languages", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response []*models.TTSLanguage
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
			break
		}
	}
	assert.True(t, foundJapanese)
}

// TestBatchSynthesize はバッチ音声生成のテスト
func TestTTSBatchSynthesize(t *testing.T) {
	router, _ := setupTTSTestRouter()

	requestBody := models.TTSBatchRequest{
		BookID:   "550e8400-e29b-41d4-a716-446655440000",
		Language: "ja",
		Options: models.TTSSynthesizeOptions{
			Speed:   1.0,
			Quality: models.TTSQualityStandard,
		},
	}

	body, _ := json.Marshal(requestBody)
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/tts/books/550e8400-e29b-41d4-a716-446655440000/batch", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusAccepted, w.Code)

	var response models.TTSBatchResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Equal(t, "550e8400-e29b-41d4-a716-446655440000", response.BookID)
	assert.Equal(t, 100, response.TotalPages)
	assert.Equal(t, 100, len(response.JobIDs))
}

// TestGetJobStatus はジョブステータス取得のテスト
func TestTTSGetJobStatus(t *testing.T) {
	router, repo := setupTTSTestRouter()

	// 完了済みジョブを取得
	testBookID := "550e8400-e29b-41d4-a716-446655440000"

	ctx := context.Background()
	bookID, _ := uuid.Parse(testBookID)
	jobs, _ := repo.GetJobsByBookID(ctx, bookID)

	assert.NotEmpty(t, jobs)

	// 最初のジョブIDを取得
	var completedJobID string
	for _, job := range jobs {
		if job.Status == models.TTSStatusCompleted {
			completedJobID = job.ID
			break
		}
	}

	assert.NotEmpty(t, completedJobID)

	req, _ := http.NewRequest(http.MethodGet, "/api/v1/tts/jobs/"+completedJobID, nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response models.TTSJobResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Equal(t, completedJobID, response.JobID)
	assert.Equal(t, models.TTSStatusCompleted, response.Status)
	assert.NotEmpty(t, response.AudioID)
	assert.NotEmpty(t, response.AudioURL)
}

// TestGetBookJobs は書籍のTTSジョブ一覧取得のテスト
func TestTTSGetBookJobs(t *testing.T) {
	router, _ := setupTTSTestRouter()

	testBookID := "550e8400-e29b-41d4-a716-446655440000"

	req, _ := http.NewRequest(http.MethodGet, "/api/v1/tts/books/"+testBookID+"/jobs", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response []*models.TTSJobResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	// サンプルデータには35ジョブ存在
	assert.Equal(t, 35, len(response))

	// 最初のジョブはページ1で完了済み
	assert.Equal(t, testBookID, response[0].BookID)
	assert.Equal(t, models.TTSStatusCompleted, response[0].Status)
}

// TestGetCacheStats はキャッシュ統計取得のテスト
func TestTTSGetCacheStats(t *testing.T) {
	router, _ := setupTTSTestRouter()

	req, _ := http.NewRequest(http.MethodGet, "/api/v1/tts/cache/stats", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response models.TTSCacheStats
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	// サンプルデータには30音声がキャッシュされている
	assert.Equal(t, 30, response.TotalCached)
	assert.Greater(t, response.CacheHitRate, 0.0)
	assert.Greater(t, response.TotalSize, int64(0))
}

// TestTTSUnauthorized は認証なしのテスト
func TestTTSUnauthorized(t *testing.T) {
	gin.SetMode(gin.TestMode)

	ttsRepo := repository.NewInMemoryTTSRepository()
	ttsHandler := NewTTSHandler(ttsRepo)

	r := gin.New()
	v1 := r.Group("/api/v1")
	ttsHandler.RegisterRoutes(v1)

	tests := []struct {
		name   string
		method string
		url    string
	}{
		{"Synthesize", http.MethodPost, "/api/v1/tts/synthesize"},
		{"Get Audio", http.MethodGet, "/api/v1/tts/audio/test-audio-id"},
		{"Get Languages", http.MethodGet, "/api/v1/tts/languages"},
		{"Batch Synthesize", http.MethodPost, "/api/v1/tts/books/550e8400-e29b-41d4-a716-446655440000/batch"},
		{"Get Job Status", http.MethodGet, "/api/v1/tts/jobs/test-job-id"},
		{"Get Book Jobs", http.MethodGet, "/api/v1/tts/books/550e8400-e29b-41d4-a716-446655440000/jobs"},
		{"Get Cache Stats", http.MethodGet, "/api/v1/tts/cache/stats"},
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

// TestTTSInvalidBookID は無効な書籍IDのテスト
func TestTTSInvalidBookID(t *testing.T) {
	router, _ := setupTTSTestRouter()

	req, _ := http.NewRequest(http.MethodGet, "/api/v1/tts/books/invalid-uuid/jobs", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// TestTTSInvalidAudioID は無効な音声IDのテスト
func TestTTSInvalidAudioID(t *testing.T) {
	router, _ := setupTTSTestRouter()

	req, _ := http.NewRequest(http.MethodGet, "/api/v1/tts/audio/invalid-audio-id", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}
