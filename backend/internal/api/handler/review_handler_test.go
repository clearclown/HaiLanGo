package handler

import (
	"bytes"
	"context"
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

func setupReviewTestRouter() (*gin.Engine, *repository.InMemoryReviewRepository) {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	repo := repository.NewInMemoryReviewRepository()
	handler := NewReviewHandler(repo)

	// テスト用のミドルウェア：user_idを設定
	r.Use(func(c *gin.Context) {
		c.Set("user_id", "550e8400-e29b-41d4-a716-446655440001")
		c.Next()
	})

	handler.RegisterRoutes(r.Group("/"))

	return r, repo
}

func TestGetStats(t *testing.T) {
	router, _ := setupReviewTestRouter()

	// リクエストを作成
	req, _ := http.NewRequest("GET", "/review/stats", nil)
	w := httptest.NewRecorder()

	// リクエストを実行
	router.ServeHTTP(w, req)

	// レスポンスを確認
	assert.Equal(t, http.StatusOK, w.Code)

	var stats models.ReviewStats
	err := json.Unmarshal(w.Body.Bytes(), &stats)
	assert.NoError(t, err)

	// サンプルデータの確認
	assert.Equal(t, 3, stats.UrgentCount)
	assert.Equal(t, 5, stats.RecommendedCount)
	assert.Equal(t, 4, stats.OptionalCount)

	t.Logf("Stats: %+v", stats)
}

func TestGetItems(t *testing.T) {
	router, _ := setupReviewTestRouter()

	// リクエストを作成
	req, _ := http.NewRequest("GET", "/review/items", nil)
	w := httptest.NewRecorder()

	// リクエストを実行
	router.ServeHTTP(w, req)

	// レスポンスを確認
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string][]*models.ReviewItem
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	items := response["items"]
	assert.NotNil(t, items)
	assert.Equal(t, 12, len(items)) // 3緊急 + 5推奨 + 4余裕 = 12

	t.Logf("Items count: %d", len(items))

	// 優先度が設定されているか確認
	for _, item := range items {
		assert.NotEmpty(t, item.Priority)
		t.Logf("Item: %s - %s (priority: %s)", item.Text, item.Translation, item.Priority)
	}
}

func TestGetItemsWithPriorityFilter(t *testing.T) {
	router, _ := setupReviewTestRouter()

	tests := []struct {
		priority      string
		expectedCount int
	}{
		{"urgent", 3},
		{"recommended", 5},
		{"optional", 4},
	}

	for _, tt := range tests {
		t.Run(tt.priority, func(t *testing.T) {
			req, _ := http.NewRequest("GET", "/review/items?priority="+tt.priority, nil)
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, http.StatusOK, w.Code)

			var response map[string][]*models.ReviewItem
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)

			items := response["items"]
			assert.Equal(t, tt.expectedCount, len(items))

			// すべてのアイテムの優先度が一致するか確認
			for _, item := range items {
				assert.Equal(t, tt.priority, item.Priority)
			}

			t.Logf("Priority %s: %d items", tt.priority, len(items))
		})
	}
}

func TestSubmitReview(t *testing.T) {
	router, repo := setupReviewTestRouter()

	// 最初のアイテムを取得
	ctx := context.Background()
	items, err := repo.FindByUserID(ctx, "550e8400-e29b-41d4-a716-446655440001")
	assert.NoError(t, err)
	assert.NotEmpty(t, items)

	testItem := items[0]
	originalMastery := testItem.MasteryLevel

	// 復習結果を送信
	result := models.ReviewResult{
		ItemID:      testItem.ID,
		Score:       90, // 高スコア
		CompletedAt: time.Now(),
	}

	body, _ := json.Marshal(result)
	req, _ := http.NewRequest("POST", "/review/submit", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// リクエストを実行
	router.ServeHTTP(w, req)

	// レスポンスを確認
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.True(t, response["success"].(bool))
	assert.NotEmpty(t, response["next_review"])

	// アイテムが更新されているか確認
	updatedItem, err := repo.FindByID(ctx, testItem.ID)
	assert.NoError(t, err)
	assert.Greater(t, updatedItem.MasteryLevel, originalMastery) // 習熟度が上がっているはず

	t.Logf("Original mastery: %d, Updated mastery: %d", originalMastery, updatedItem.MasteryLevel)
	t.Logf("Next review: %s", response["next_review"])
}

func TestSubmitReview_LowScore(t *testing.T) {
	router, repo := setupReviewTestRouter()

	// 最初のアイテムを取得
	ctx := context.Background()
	items, err := repo.FindByUserID(ctx, "550e8400-e29b-41d4-a716-446655440001")
	assert.NoError(t, err)
	assert.NotEmpty(t, items)

	testItem := items[1]
	originalMastery := testItem.MasteryLevel

	// 低スコアの復習結果を送信
	result := models.ReviewResult{
		ItemID:      testItem.ID,
		Score:       30, // 低スコア
		CompletedAt: time.Now(),
	}

	body, _ := json.Marshal(result)
	req, _ := http.NewRequest("POST", "/review/submit", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// リクエストを実行
	router.ServeHTTP(w, req)

	// レスポンスを確認
	assert.Equal(t, http.StatusOK, w.Code)

	// アイテムが更新されているか確認
	updatedItem, err := repo.FindByID(ctx, testItem.ID)
	assert.NoError(t, err)
	assert.Less(t, updatedItem.MasteryLevel, originalMastery) // 習熟度が下がっているはず
	assert.Equal(t, 1, updatedItem.IntervalDays) // 失敗したので1日後に再復習

	t.Logf("Original mastery: %d, Updated mastery: %d", originalMastery, updatedItem.MasteryLevel)
	t.Logf("Interval reset to: %d days", updatedItem.IntervalDays)
}

func TestSubmitReview_InvalidItemID(t *testing.T) {
	router, _ := setupReviewTestRouter()

	// 存在しないアイテムID
	result := models.ReviewResult{
		ItemID:      "non-existent-id",
		Score:       90,
		CompletedAt: time.Now(),
	}

	body, _ := json.Marshal(result)
	req, _ := http.NewRequest("POST", "/review/submit", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// リクエストを実行
	router.ServeHTTP(w, req)

	// レスポンスを確認
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestSubmitReview_Unauthorized(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	repo := repository.NewInMemoryReviewRepository()
	handler := NewReviewHandler(repo)

	// 異なるユーザーIDを設定
	r.Use(func(c *gin.Context) {
		c.Set("user_id", "different-user-id")
		c.Next()
	})

	handler.RegisterRoutes(r.Group("/"))

	// サンプルユーザーのアイテムを取得しようとする
	ctx := context.Background()
	items, _ := repo.FindByUserID(ctx, "550e8400-e29b-41d4-a716-446655440001")
	testItem := items[0]

	result := models.ReviewResult{
		ItemID:      testItem.ID,
		Score:       90,
		CompletedAt: time.Now(),
	}

	body, _ := json.Marshal(result)
	req, _ := http.NewRequest("POST", "/review/submit", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// リクエストを実行
	r.ServeHTTP(w, req)

	// レスポンスを確認（403 Forbidden）
	assert.Equal(t, http.StatusForbidden, w.Code)
}
