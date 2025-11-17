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
	ocrservice "github.com/clearclown/HaiLanGo/backend/internal/service/ocr"
	"github.com/clearclown/HaiLanGo/backend/internal/websocket"
	"github.com/clearclown/HaiLanGo/backend/pkg/cache"
	"github.com/clearclown/HaiLanGo/backend/pkg/ocr"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func setupOCRTestRouter() (*gin.Engine, repository.OCRRepositoryInterface) {
	gin.SetMode(gin.TestMode)

	ocrRepo := repository.NewInMemoryOCRRepository()

	// OCRサービスのセットアップ
	ocrClient, _ := ocr.NewOCRClient()
	mockCache := cache.NewMockCache()
	ocrSvc := ocrservice.NewOCRService(ocrClient, mockCache)

	// WebSocketハブのセットアップ
	wsHub := websocket.NewHub()

	ocrHandler := NewOCRHandler(ocrRepo, ocrSvc, wsHub)

	r := gin.New()

	// テスト用の認証ミドルウェア
	r.Use(func(c *gin.Context) {
		c.Set("user_id", "550e8400-e29b-41d4-a716-446655440000")
		c.Next()
	})

	// ルート登録
	v1 := r.Group("/api/v1")
	ocrHandler.RegisterRoutes(v1)

	return r, ocrRepo
}

// TestProcessPage はページOCR処理のテスト
func TestProcessPage(t *testing.T) {
	router, _ := setupOCRTestRouter()

	requestBody := models.ProcessPageOCRRequest{
		PageNumber: 51,
		Language:   "ru",
		Options: models.OCROptions{
			DetectOrientation: true,
			DetectLanguage:    false,
		},
	}

	body, _ := json.Marshal(requestBody)
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/ocr/books/550e8400-e29b-41d4-a716-446655440000/pages/51", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusAccepted, w.Code)

	var response models.OCRJobResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.NotEmpty(t, response.JobID)
	assert.Equal(t, "550e8400-e29b-41d4-a716-446655440000", response.BookID)
	assert.Equal(t, 51, response.PageNumber)
	assert.Equal(t, models.OCRStatusPending, response.Status)
}

// TestGetJobStatus はジョブステータス取得のテスト
func TestGetJobStatus(t *testing.T) {
	router, repo := setupOCRTestRouter()

	// サンプルデータには50ページ分のジョブが存在
	// ページ1-45: 完了済み
	// ページ46-48: 処理中
	// ページ49-50: ペンディング

	// 完了済みジョブを取得
	testBookID := "550e8400-e29b-41d4-a716-446655440000"

	ctx := context.Background()
	bookID, _ := uuid.Parse(testBookID)
	jobs, _ := repo.GetJobsByBookID(ctx, bookID)

	assert.NotEmpty(t, jobs)

	// 最初のジョブIDを取得
	var completedJobID string
	for _, job := range jobs {
		if job.Status == models.OCRStatusCompleted {
			completedJobID = job.ID
			break
		}
	}

	assert.NotEmpty(t, completedJobID)

	req, _ := http.NewRequest(http.MethodGet, "/api/v1/ocr/jobs/"+completedJobID, nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response models.OCRJobResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Equal(t, completedJobID, response.JobID)
	assert.Equal(t, models.OCRStatusCompleted, response.Status)
}

// TestGetJobResult は処理結果取得のテスト
func TestGetJobResult(t *testing.T) {
	router, repo := setupOCRTestRouter()

	// 完了済みジョブを取得
	testBookID := "550e8400-e29b-41d4-a716-446655440000"

	ctx := context.Background()
	bookID, _ := uuid.Parse(testBookID)
	jobs, _ := repo.GetJobsByBookID(ctx, bookID)

	var completedJobID string
	for _, job := range jobs {
		if job.Status == models.OCRStatusCompleted {
			completedJobID = job.ID
			break
		}
	}

	req, _ := http.NewRequest(http.MethodGet, "/api/v1/ocr/jobs/"+completedJobID+"/result", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response models.OCRResultResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Equal(t, completedJobID, response.JobID)
	assert.Equal(t, models.OCRStatusCompleted, response.Status)
	assert.NotNil(t, response.Result)
	assert.NotEmpty(t, response.Result.Text)
	assert.Equal(t, "ru", response.Result.DetectedLanguage)
	assert.Greater(t, response.Result.Confidence, 0.9)
}

// TestGetJobResultStillProcessing は処理中のジョブ結果取得のテスト
func TestGetJobResultStillProcessing(t *testing.T) {
	router, repo := setupOCRTestRouter()

	// 処理中のジョブを取得
	testBookID := "550e8400-e29b-41d4-a716-446655440000"

	ctx := context.Background()
	bookID, _ := uuid.Parse(testBookID)
	jobs, _ := repo.GetJobsByBookID(ctx, bookID)

	var processingJobID string
	for _, job := range jobs {
		if job.Status == models.OCRStatusProcessing {
			processingJobID = job.ID
			break
		}
	}

	assert.NotEmpty(t, processingJobID)

	req, _ := http.NewRequest(http.MethodGet, "/api/v1/ocr/jobs/"+processingJobID+"/result", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusAccepted, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Equal(t, "Job is still processing", response["message"])
	assert.Equal(t, "processing", response["status"])
}

// TestGetBookJobs は書籍のOCRジョブ一覧取得のテスト
func TestGetBookJobs(t *testing.T) {
	router, _ := setupOCRTestRouter()

	testBookID := "550e8400-e29b-41d4-a716-446655440000"

	req, _ := http.NewRequest(http.MethodGet, "/api/v1/ocr/books/"+testBookID+"/jobs", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response []*models.OCRJobResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	// サンプルデータには50ジョブ存在
	assert.Equal(t, 50, len(response))

	// 最初のジョブはページ1で完了済み
	assert.Equal(t, testBookID, response[0].BookID)
	assert.Equal(t, models.OCRStatusCompleted, response[0].Status)
}

// TestGetStatistics はOCR統計情報取得のテスト
func TestGetStatistics(t *testing.T) {
	router, _ := setupOCRTestRouter()

	req, _ := http.NewRequest(http.MethodGet, "/api/v1/ocr/statistics", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response models.OCRStatistics
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	// サンプルデータには50ジョブ存在
	assert.Equal(t, 50, response.TotalJobs)

	// ページ1-45: 完了済み
	assert.Equal(t, 45, response.CompletedJobs)

	// ページ46-48: 処理中
	assert.Equal(t, 3, response.ProcessingJobs)

	// ページ49-50: ペンディング
	assert.Equal(t, 2, response.PendingJobs)

	// 平均信頼度
	assert.Greater(t, response.AverageConfidence, 0.9)
}

// TestBatchProcess はバッチOCR処理のテスト
func TestBatchProcess(t *testing.T) {
	router, _ := setupOCRTestRouter()

	requestBody := models.BatchOCRRequest{
		BookID:   "550e8400-e29b-41d4-a716-446655440000",
		Language: "ru",
		Options: models.OCROptions{
			DetectOrientation: true,
		},
	}

	body, _ := json.Marshal(requestBody)
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/ocr/books/550e8400-e29b-41d4-a716-446655440000/batch", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusAccepted, w.Code)

	var response models.BatchOCRResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Equal(t, "550e8400-e29b-41d4-a716-446655440000", response.BookID)
	assert.Equal(t, 150, response.TotalPages)
	assert.Equal(t, 150, len(response.JobIDs))
}

// TestOCRUnauthorized は認証なしのテスト
func TestOCRUnauthorized(t *testing.T) {
	gin.SetMode(gin.TestMode)

	ocrRepo := repository.NewInMemoryOCRRepository()

	// OCRサービスのセットアップ
	ocrClient, _ := ocr.NewOCRClient()
	mockCache := cache.NewMockCache()
	ocrSvc := ocrservice.NewOCRService(ocrClient, mockCache)

	// WebSocketハブのセットアップ
	wsHub := websocket.NewHub()

	ocrHandler := NewOCRHandler(ocrRepo, ocrSvc, wsHub)

	r := gin.New()
	v1 := r.Group("/api/v1")
	ocrHandler.RegisterRoutes(v1)

	tests := []struct {
		name   string
		method string
		url    string
	}{
		{"Process Page", http.MethodPost, "/api/v1/ocr/books/550e8400-e29b-41d4-a716-446655440000/pages/1"},
		{"Get Job Status", http.MethodGet, "/api/v1/ocr/jobs/test-job-id"},
		{"Get Job Result", http.MethodGet, "/api/v1/ocr/jobs/test-job-id/result"},
		{"Get Book Jobs", http.MethodGet, "/api/v1/ocr/books/550e8400-e29b-41d4-a716-446655440000/jobs"},
		{"Get Statistics", http.MethodGet, "/api/v1/ocr/statistics"},
		{"Batch Process", http.MethodPost, "/api/v1/ocr/books/550e8400-e29b-41d4-a716-446655440000/batch"},
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

// TestOCRInvalidBookID は無効な書籍IDのテスト
func TestOCRInvalidBookID(t *testing.T) {
	router, _ := setupOCRTestRouter()

	req, _ := http.NewRequest(http.MethodGet, "/api/v1/ocr/books/invalid-uuid/jobs", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// TestOCRInvalidPageNumber は無効なページ番号のテスト
func TestOCRInvalidPageNumber(t *testing.T) {
	router, _ := setupOCRTestRouter()

	requestBody := models.ProcessPageOCRRequest{
		PageNumber: 0, // 無効（1以上が必要）
		Language:   "ru",
	}

	body, _ := json.Marshal(requestBody)
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/ocr/books/550e8400-e29b-41d4-a716-446655440000/pages/0", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}
