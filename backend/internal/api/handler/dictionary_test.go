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
	"github.com/stretchr/testify/assert"
)

func setupDictionaryTestRouter() (*gin.Engine, repository.DictionaryRepositoryInterface) {
	gin.SetMode(gin.TestMode)

	dictionaryRepo := repository.NewInMemoryDictionaryRepository()
	dictionaryHandler := NewDictionaryHandler(dictionaryRepo)

	r := gin.New()

	// テスト用の認証ミドルウェア
	r.Use(func(c *gin.Context) {
		c.Set("user_id", "550e8400-e29b-41d4-a716-446655440000")
		c.Next()
	})

	// ルート登録
	v1 := r.Group("/api/v1")
	dictionaryHandler.RegisterRoutes(v1)

	return r, dictionaryRepo
}

// TestLookupWord は単語検索のテスト
func TestLookupWord(t *testing.T) {
	router, _ := setupDictionaryTestRouter()

	tests := []struct {
		name         string
		word         string
		language     string
		expectedCode int
		checkWord    string
	}{
		{
			name:         "英語単語検索",
			word:         "hello",
			language:     "en",
			expectedCode: http.StatusOK,
			checkWord:    "hello",
		},
		{
			name:         "デフォルト言語（英語）",
			word:         "book",
			language:     "",
			expectedCode: http.StatusOK,
			checkWord:    "book",
		},
		{
			name:         "ロシア語単語検索",
			word:         "здравствуйте",
			language:     "ru",
			expectedCode: http.StatusOK,
			checkWord:    "здравствуйте",
		},
		{
			name:         "日本語単語検索",
			word:         "こんにちは",
			language:     "ja",
			expectedCode: http.StatusOK,
			checkWord:    "こんにちは",
		},
		{
			name:         "存在しない単語（ダミーエントリ生成）",
			word:         "nonexistent",
			language:     "en",
			expectedCode: http.StatusOK,
			checkWord:    "nonexistent",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			url := "/api/v1/dictionary/words/" + tt.word
			if tt.language != "" {
				url += "?language=" + tt.language
			}

			req, _ := http.NewRequest(http.MethodGet, url, nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedCode, w.Code)

			if tt.expectedCode == http.StatusOK {
				var entry models.WordEntry
				err := json.Unmarshal(w.Body.Bytes(), &entry)
				assert.NoError(t, err)
				assert.Equal(t, tt.checkWord, entry.Word)
				assert.NotEmpty(t, entry.Meanings)
			}
		})
	}
}

// TestBatchLookup は複数単語検索のテスト
func TestBatchLookup(t *testing.T) {
	router, _ := setupDictionaryTestRouter()

	tests := []struct {
		name         string
		requestBody  BatchLookupRequest
		expectedCode int
		expectedCount int
	}{
		{
			name: "複数単語検索（英語）",
			requestBody: BatchLookupRequest{
				Words:    []string{"hello", "book", "test"},
				Language: "en",
			},
			expectedCode: http.StatusOK,
			expectedCount: 3,
		},
		{
			name: "1単語のみ",
			requestBody: BatchLookupRequest{
				Words:    []string{"hello"},
				Language: "en",
			},
			expectedCode: http.StatusOK,
			expectedCount: 1,
		},
		{
			name: "空の配列",
			requestBody: BatchLookupRequest{
				Words:    []string{},
				Language: "en",
			},
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "50単語制限を超える",
			requestBody: BatchLookupRequest{
				Words:    make([]string, 51), // 51単語
				Language: "en",
			},
			expectedCode: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 51単語の配列を生成
			if len(tt.requestBody.Words) == 51 {
				for i := 0; i < 51; i++ {
					tt.requestBody.Words[i] = "word"
				}
			}

			body, _ := json.Marshal(tt.requestBody)
			req, _ := http.NewRequest(http.MethodPost, "/api/v1/dictionary/batch", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedCode, w.Code)

			if tt.expectedCode == http.StatusOK {
				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, float64(tt.expectedCount), response["count"])

				entries, ok := response["entries"].([]interface{})
				assert.True(t, ok)
				assert.Equal(t, tt.expectedCount, len(entries))
			}
		})
	}
}

// TestGetSupportedLanguages はサポート言語一覧取得のテスト
func TestGetSupportedLanguages(t *testing.T) {
	router, _ := setupDictionaryTestRouter()

	req, _ := http.NewRequest(http.MethodGet, "/api/v1/dictionary/languages", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	languages, ok := response["languages"].([]interface{})
	assert.True(t, ok)
	assert.Equal(t, 12, len(languages)) // 12言語サポート

	// 主要言語が含まれているか確認
	languageStrs := make([]string, len(languages))
	for i, lang := range languages {
		languageStrs[i] = lang.(string)
	}
	assert.Contains(t, languageStrs, "en")
	assert.Contains(t, languageStrs, "ja")
	assert.Contains(t, languageStrs, "ru")
	assert.Contains(t, languageStrs, "zh")
}

// TestDictionaryUnauthorized は認証なしのテスト
func TestDictionaryUnauthorized(t *testing.T) {
	gin.SetMode(gin.TestMode)

	dictionaryRepo := repository.NewInMemoryDictionaryRepository()
	dictionaryHandler := NewDictionaryHandler(dictionaryRepo)

	r := gin.New()
	v1 := r.Group("/api/v1")
	dictionaryHandler.RegisterRoutes(v1)

	tests := []struct {
		name   string
		method string
		url    string
		body   string
	}{
		{"Lookup Word", http.MethodGet, "/api/v1/dictionary/words/hello", ""},
		{"Batch Lookup", http.MethodPost, "/api/v1/dictionary/batch", `{"words":["hello"],"language":"en"}`},
		{"Get Languages", http.MethodGet, "/api/v1/dictionary/languages", ""},
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

// TestBatchLookupInvalidRequest は無効なリクエストボディのテスト
func TestBatchLookupInvalidRequest(t *testing.T) {
	router, _ := setupDictionaryTestRouter()

	req, _ := http.NewRequest(http.MethodPost, "/api/v1/dictionary/batch", bytes.NewBuffer([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// TestLookupWordCaseInsensitive は言語コードの大文字小文字非依存テスト
func TestLookupWordCaseInsensitive(t *testing.T) {
	router, _ := setupDictionaryTestRouter()

	tests := []struct {
		name     string
		language string
	}{
		{"小文字", "en"},
		{"大文字", "EN"},
		{"混在", "En"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			url := "/api/v1/dictionary/words/hello?language=" + tt.language

			req, _ := http.NewRequest(http.MethodGet, url, nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, http.StatusOK, w.Code)

			var entry models.WordEntry
			err := json.Unmarshal(w.Body.Bytes(), &entry)
			assert.NoError(t, err)
			assert.Equal(t, "hello", entry.Word)
		})
	}
}
