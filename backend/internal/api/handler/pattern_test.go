package handler

import (
	"bytes"
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

func setupPatternTestRouter() (*gin.Engine, repository.PatternRepositoryInterface) {
	gin.SetMode(gin.TestMode)

	patternRepo := repository.NewInMemoryPatternRepository()
	patternHandler := NewPatternHandler(patternRepo)

	r := gin.New()

	// テスト用の認証ミドルウェア
	r.Use(func(c *gin.Context) {
		c.Set("user_id", "550e8400-e29b-41d4-a716-446655440000")
		c.Next()
	})

	// ルート登録
	v1 := r.Group("/api/v1")
	patternHandler.RegisterRoutes(v1)

	return r, patternRepo
}

// TestExtractPatterns はパターン抽出のテスト
func TestExtractPatterns(t *testing.T) {
	router, _ := setupPatternTestRouter()

	requestBody := ExtractPatternsRequest{
		BookID:       uuid.MustParse("550e8400-e29b-41d4-a716-446655440001"),
		MinFrequency: 2,
	}

	body, _ := json.Marshal(requestBody)
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/patterns/extract", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.NotNil(t, response["patterns"])
	assert.NotNil(t, response["total_found"])
	assert.NotNil(t, response["processed_pages"])
	assert.NotNil(t, response["duration_ms"])
}

// TestExtractPatternsInvalidRequest は無効なリクエストのテスト
func TestExtractPatternsInvalidRequest(t *testing.T) {
	router, _ := setupPatternTestRouter()

	req, _ := http.NewRequest(http.MethodPost, "/api/v1/patterns/extract", bytes.NewBuffer([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// TestGetPatternsByBook は書籍のパターン一覧取得のテスト
func TestGetPatternsByBook(t *testing.T) {
	router, _ := setupPatternTestRouter()

	// サンプルデータの書籍ID
	bookID := "550e8400-e29b-41d4-a716-446655440001"

	req, _ := http.NewRequest(http.MethodGet, "/api/v1/patterns/books/"+bookID, nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.NotNil(t, response["patterns"])
	assert.NotNil(t, response["book_id"])
	assert.NotNil(t, response["count"])

	// サンプルデータには2つのパターンがある
	count := response["count"].(float64)
	assert.Equal(t, float64(2), count)
}

// TestGetPatternsByBookInvalidID は無効な書籍IDのテスト
func TestGetPatternsByBookInvalidID(t *testing.T) {
	router, _ := setupPatternTestRouter()

	req, _ := http.NewRequest(http.MethodGet, "/api/v1/patterns/books/invalid-id", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// TestGetPatternByID はパターン詳細取得のテスト
func TestGetPatternByID(t *testing.T) {
	router, patternRepo := setupPatternTestRouter()

	// サンプルデータから最初のパターンを取得
	sampleBookID := uuid.MustParse("550e8400-e29b-41d4-a716-446655440001")
	patterns, _ := patternRepo.GetPatternsByBookID(nil, sampleBookID)
	assert.Greater(t, len(patterns), 0)

	patternID := patterns[0].ID

	req, _ := http.NewRequest(http.MethodGet, "/api/v1/patterns/"+patternID.String(), nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response models.Pattern
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Equal(t, patternID, response.ID)
	assert.NotEmpty(t, response.Pattern)
}

// TestGetPatternByIDNotFound は存在しないパターンのテスト
func TestGetPatternByIDNotFound(t *testing.T) {
	router, _ := setupPatternTestRouter()

	nonExistentID := uuid.New()

	req, _ := http.NewRequest(http.MethodGet, "/api/v1/patterns/"+nonExistentID.String(), nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

// TestGetPatternExamples はパターン使用例取得のテスト
func TestGetPatternExamples(t *testing.T) {
	router, patternRepo := setupPatternTestRouter()

	// サンプルデータから最初のパターンを取得
	sampleBookID := uuid.MustParse("550e8400-e29b-41d4-a716-446655440001")
	patterns, _ := patternRepo.GetPatternsByBookID(nil, sampleBookID)
	assert.Greater(t, len(patterns), 0)

	patternID := patterns[0].ID

	tests := []struct {
		name         string
		limitParam   string
		expectedCode int
	}{
		{
			name:         "デフォルトlimit",
			limitParam:   "",
			expectedCode: http.StatusOK,
		},
		{
			name:         "limit=5",
			limitParam:   "5",
			expectedCode: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			url := "/api/v1/patterns/" + patternID.String() + "/examples"
			if tt.limitParam != "" {
				url += "?limit=" + tt.limitParam
			}

			req, _ := http.NewRequest(http.MethodGet, url, nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedCode, w.Code)

			var response map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)

			assert.NotNil(t, response["examples"])
			assert.NotNil(t, response["pattern_id"])
			assert.NotNil(t, response["count"])
		})
	}
}

// TestGetPatternPractice はパターン練習問題取得のテスト
func TestGetPatternPractice(t *testing.T) {
	router, patternRepo := setupPatternTestRouter()

	// サンプルデータから最初のパターンを取得
	sampleBookID := uuid.MustParse("550e8400-e29b-41d4-a716-446655440001")
	patterns, _ := patternRepo.GetPatternsByBookID(nil, sampleBookID)
	assert.Greater(t, len(patterns), 0)

	patternID := patterns[0].ID

	tests := []struct {
		name         string
		countParam   string
		expectedCode int
	}{
		{
			name:         "デフォルトcount",
			countParam:   "",
			expectedCode: http.StatusOK,
		},
		{
			name:         "count=5",
			countParam:   "5",
			expectedCode: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			url := "/api/v1/patterns/" + patternID.String() + "/practice"
			if tt.countParam != "" {
				url += "?count=" + tt.countParam
			}

			req, _ := http.NewRequest(http.MethodGet, url, nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedCode, w.Code)

			var response map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)

			assert.NotNil(t, response["practices"])
			assert.NotNil(t, response["pattern_id"])
			assert.NotNil(t, response["count"])
		})
	}
}

// TestUpdatePatternProgress はパターン学習進捗更新のテスト
func TestUpdatePatternProgress(t *testing.T) {
	tests := []struct {
		name         string
		correct      bool
		expectedCode int
	}{
		{
			name:         "正解の場合",
			correct:      true,
			expectedCode: http.StatusOK,
		},
		{
			name:         "不正解の場合",
			correct:      false,
			expectedCode: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 各テストで新しいルーターを作成
			router, patternRepo := setupPatternTestRouter()

			// サンプルデータから最初のパターンを取得
			sampleBookID := uuid.MustParse("550e8400-e29b-41d4-a716-446655440001")
			patterns, _ := patternRepo.GetPatternsByBookID(nil, sampleBookID)
			assert.Greater(t, len(patterns), 0)

			patternID := patterns[0].ID

			requestBody := UpdatePatternProgressRequest{
				Correct: tt.correct,
			}

			body, _ := json.Marshal(requestBody)
			req, _ := http.NewRequest(http.MethodPost, "/api/v1/patterns/"+patternID.String()+"/progress", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedCode, w.Code)

			var response models.PatternProgress
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)

			assert.Equal(t, patternID, response.PatternID)
			assert.NotNil(t, response.MasteryLevel)
			assert.Greater(t, response.PracticeCount, 0)
		})
	}
}

// TestPatternUnauthorized は認証なしのテスト
func TestPatternUnauthorized(t *testing.T) {
	gin.SetMode(gin.TestMode)

	patternRepo := repository.NewInMemoryPatternRepository()
	patternHandler := NewPatternHandler(patternRepo)

	r := gin.New()
	v1 := r.Group("/api/v1")
	patternHandler.RegisterRoutes(v1)

	tests := []struct {
		name   string
		method string
		url    string
		body   string
	}{
		{"Extract Patterns", http.MethodPost, "/api/v1/patterns/extract", `{"book_id":"550e8400-e29b-41d4-a716-446655440001"}`},
		{"Get Patterns By Book", http.MethodGet, "/api/v1/patterns/books/550e8400-e29b-41d4-a716-446655440001", ""},
		{"Get Pattern By ID", http.MethodGet, "/api/v1/patterns/" + uuid.New().String(), ""},
		{"Get Pattern Examples", http.MethodGet, "/api/v1/patterns/" + uuid.New().String() + "/examples", ""},
		{"Get Pattern Practice", http.MethodGet, "/api/v1/patterns/" + uuid.New().String() + "/practice", ""},
		{"Update Pattern Progress", http.MethodPost, "/api/v1/patterns/" + uuid.New().String() + "/progress", `{"correct":true}`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var req *http.Request
			if tt.body != "" {
				req, _ = http.NewRequest(tt.method, tt.url, bytes.NewBufferString(tt.body))
				req.Header.Set("Content-Type", "application/json")
			} else {
				req, _ = http.NewRequest(tt.method, tt.url, nil)
			}
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			assert.Equal(t, http.StatusUnauthorized, w.Code)
		})
	}
}
