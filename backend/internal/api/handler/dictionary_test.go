package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/clearclown/HaiLanGo/backend/internal/models"
	"github.com/clearclown/HaiLanGo/backend/internal/service/dictionary"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMain(m *testing.M) {
	// Enable mock mode for all tests
	os.Setenv("TEST_USE_MOCKS", "true")
	code := m.Run()
	os.Exit(code)
}

func setupTestHandler(t *testing.T) *DictionaryHandler {
	service, err := dictionary.NewService()
	require.NoError(t, err)
	return NewDictionaryHandler(service)
}

func TestDictionaryHandler_LookupWord(t *testing.T) {
	handler := setupTestHandler(t)

	t.Run("Success", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/v1/dictionary/words/hello?language=en", nil)
		w := httptest.NewRecorder()

		handler.LookupWord(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, "application/json", w.Header().Get("Content-Type"))

		var entry models.WordEntry
		err := json.NewDecoder(w.Body).Decode(&entry)
		require.NoError(t, err)
		assert.Equal(t, "hello", entry.Word)
		assert.NotEmpty(t, entry.Meanings)
	})

	t.Run("DefaultLanguage", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/v1/dictionary/words/hello", nil)
		w := httptest.NewRecorder()

		handler.LookupWord(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var entry models.WordEntry
		err := json.NewDecoder(w.Body).Decode(&entry)
		require.NoError(t, err)
		assert.Equal(t, "en", entry.Language)
	})

	t.Run("WordNotFound", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/v1/dictionary/words/xyzabc123notaword?language=en", nil)
		w := httptest.NewRecorder()

		handler.LookupWord(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("InvalidURL", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/v1/dictionary", nil)
		w := httptest.NewRecorder()

		handler.LookupWord(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestDictionaryHandler_LookupWordDetails(t *testing.T) {
	handler := setupTestHandler(t)

	t.Run("Success", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/v1/dictionary/words/hello/details?language=en", nil)
		w := httptest.NewRecorder()

		handler.LookupWordDetails(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, "application/json", w.Header().Get("Content-Type"))

		var entry models.WordEntry
		err := json.NewDecoder(w.Body).Decode(&entry)
		require.NoError(t, err)
		assert.Equal(t, "hello", entry.Word)
		assert.NotEmpty(t, entry.Meanings)
	})

	t.Run("DefaultLanguage", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/v1/dictionary/words/hello/details", nil)
		w := httptest.NewRecorder()

		handler.LookupWordDetails(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var entry models.WordEntry
		err := json.NewDecoder(w.Body).Decode(&entry)
		require.NoError(t, err)
		assert.Equal(t, "en", entry.Language)
	})

	t.Run("WordNotFound", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/v1/dictionary/words/xyzabc123notaword/details?language=en", nil)
		w := httptest.NewRecorder()

		handler.LookupWordDetails(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("InvalidURL", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/v1/dictionary", nil)
		w := httptest.NewRecorder()

		handler.LookupWordDetails(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestDictionaryHandler_Integration(t *testing.T) {
	handler := setupTestHandler(t)

	t.Run("MultipleLanguages", func(t *testing.T) {
		languages := []string{"en", "es", "fr", "de"}

		for _, lang := range languages {
			req := httptest.NewRequest("GET", "/api/v1/dictionary/words/hello?language="+lang, nil)
			w := httptest.NewRecorder()

			handler.LookupWord(w, req)

			// Should succeed or return not found (depending on mock data)
			assert.True(t, w.Code == http.StatusOK || w.Code == http.StatusNotFound)
		}
	})

	t.Run("CacheEffectiveness", func(t *testing.T) {
		// First request
		req1 := httptest.NewRequest("GET", "/api/v1/dictionary/words/hello?language=en", nil)
		w1 := httptest.NewRecorder()
		handler.LookupWord(w1, req1)
		require.Equal(t, http.StatusOK, w1.Code)

		// Second request (should hit cache)
		req2 := httptest.NewRequest("GET", "/api/v1/dictionary/words/hello?language=en", nil)
		w2 := httptest.NewRecorder()
		handler.LookupWord(w2, req2)
		require.Equal(t, http.StatusOK, w2.Code)

		// Verify both responses are identical
		var entry1, entry2 models.WordEntry
		require.NoError(t, json.NewDecoder(w1.Body).Decode(&entry1))
		require.NoError(t, json.NewDecoder(w2.Body).Decode(&entry2))
		assert.Equal(t, entry1.Word, entry2.Word)
	})
}

// Benchmark tests
func BenchmarkDictionaryHandler_LookupWord(b *testing.B) {
	os.Setenv("TEST_USE_MOCKS", "true")
	service, _ := dictionary.NewService()
	handler := NewDictionaryHandler(service)

	req := httptest.NewRequest("GET", "/api/v1/dictionary/words/hello?language=en", nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		handler.LookupWord(w, req)
	}
}

func BenchmarkDictionaryHandler_LookupWordDetails(b *testing.B) {
	os.Setenv("TEST_USE_MOCKS", "true")
	service, _ := dictionary.NewService()
	handler := NewDictionaryHandler(service)

	req := httptest.NewRequest("GET", "/api/v1/dictionary/words/hello/details?language=en", nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		handler.LookupWordDetails(w, req)
	}
}
