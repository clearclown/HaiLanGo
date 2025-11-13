// Package teachermode provides teacher mode API handlers
package teachermode

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGeneratePlaylistHandler(t *testing.T) {
	t.Run("正常にプレイリストが生成される", func(t *testing.T) {
		handler := NewHandler()

		reqBody := map[string]interface{}{
			"settings": map[string]interface{}{
				"speed":        1.0,
				"pageInterval": 5,
				"repeatCount":  1,
				"audioQuality": "standard",
				"content": map[string]bool{
					"includeTranslation":           true,
					"includeWordExplanation":       true,
					"includeGrammarExplanation":    false,
					"includePronunciationPractice": false,
					"includeExampleSentences":      false,
				},
			},
			"pageRange": map[string]int{
				"start": 1,
				"end":   150,
			},
		}

		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest("POST", "/api/v1/books/test-book/teacher-mode/generate", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		handler.GeneratePlaylist(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		var response map[string]interface{}
		err := json.Unmarshal(rec.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.NotEmpty(t, response["playlistId"])
		assert.NotEmpty(t, response["pages"])
	})

	t.Run("無効なリクエストボディでエラーが返る", func(t *testing.T) {
		handler := NewHandler()

		req := httptest.NewRequest("POST", "/api/v1/books/test-book/teacher-mode/generate", bytes.NewReader([]byte("invalid json")))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		handler.GeneratePlaylist(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("認証なしでエラーが返る", func(t *testing.T) {
		handler := NewHandler()

		reqBody := map[string]interface{}{
			"settings": map[string]interface{}{
				"speed":        1.0,
				"pageInterval": 5,
				"repeatCount":  1,
				"audioQuality": "standard",
			},
		}

		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest("POST", "/api/v1/books/test-book/teacher-mode/generate", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		// 認証ヘッダーなし
		rec := httptest.NewRecorder()

		handler.GeneratePlaylist(rec, req)

		assert.Equal(t, http.StatusUnauthorized, rec.Code)
	})
}

func TestGenerateDownloadPackageHandler(t *testing.T) {
	t.Run("正常にダウンロードパッケージが生成される", func(t *testing.T) {
		handler := NewHandler()

		reqBody := map[string]interface{}{
			"settings": map[string]interface{}{
				"speed":        1.0,
				"pageInterval": 5,
				"repeatCount":  1,
				"audioQuality": "standard",
			},
		}

		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest("POST", "/api/v1/books/test-book/teacher-mode/download-package", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer test-token")
		rec := httptest.NewRecorder()

		handler.GenerateDownloadPackage(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		var response map[string]interface{}
		err := json.Unmarshal(rec.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.NotEmpty(t, response["packageId"])
		assert.NotEmpty(t, response["downloadUrl"])
		assert.NotZero(t, response["totalSize"])
	})
}

func TestUpdatePlaybackStateHandler(t *testing.T) {
	t.Run("正常に再生状態が保存される", func(t *testing.T) {
		handler := NewHandler()

		reqBody := map[string]interface{}{
			"currentPage":         12,
			"currentSegmentIndex": 2,
			"elapsedTime":         3600,
		}

		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest("PUT", "/api/v1/books/test-book/teacher-mode/playback-state", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer test-token")
		rec := httptest.NewRecorder()

		handler.UpdatePlaybackState(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		var response map[string]interface{}
		err := json.Unmarshal(rec.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.True(t, response["success"].(bool))
	})

	t.Run("無効なデータでエラーが返る", func(t *testing.T) {
		handler := NewHandler()

		reqBody := map[string]interface{}{
			"currentPage": -1, // 無効なページ番号
		}

		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest("PUT", "/api/v1/books/test-book/teacher-mode/playback-state", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer test-token")
		rec := httptest.NewRecorder()

		handler.UpdatePlaybackState(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})
}

func TestGetPlaylistHandler(t *testing.T) {
	t.Run("正常にプレイリストが取得される", func(t *testing.T) {
		handler := NewHandler()

		req := httptest.NewRequest("GET", "/api/v1/books/test-book/teacher-mode/playlist", nil)
		req.Header.Set("Authorization", "Bearer test-token")
		rec := httptest.NewRecorder()

		handler.GetPlaylist(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		var response map[string]interface{}
		err := json.Unmarshal(rec.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.NotEmpty(t, response["playlistId"])
		assert.NotEmpty(t, response["pages"])
	})

	t.Run("プレイリストが存在しない場合は404が返る", func(t *testing.T) {
		handler := NewHandler()

		req := httptest.NewRequest("GET", "/api/v1/books/nonexistent-book/teacher-mode/playlist", nil)
		req.Header.Set("Authorization", "Bearer test-token")
		rec := httptest.NewRecorder()

		handler.GetPlaylist(rec, req)

		assert.Equal(t, http.StatusNotFound, rec.Code)
	})
}
