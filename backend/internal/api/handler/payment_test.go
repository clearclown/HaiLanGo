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

func setupPaymentTestRouter() (*gin.Engine, repository.PaymentRepositoryInterface) {
	gin.SetMode(gin.TestMode)

	paymentRepo := repository.NewInMemoryPaymentRepository()
	paymentHandler := NewPaymentHandler(paymentRepo)

	r := gin.New()

	// テスト用の認証ミドルウェア
	r.Use(func(c *gin.Context) {
		c.Set("user_id", "550e8400-e29b-41d4-a716-446655440000")
		c.Next()
	})

	// ルート登録
	v1 := r.Group("/api/v1")
	paymentHandler.RegisterRoutes(v1)

	return r, paymentRepo
}

// TestCreateSubscription はサブスクリプション作成のテスト
func TestCreateSubscription(t *testing.T) {
	router, _ := setupPaymentTestRouter()

	requestBody := models.CreateSubscriptionRequest{
		Plan:          models.PlanPremium,
		PaymentMethod: "pm_test_123456",
	}

	body, _ := json.Marshal(requestBody)
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/payment/subscription/create", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var response models.SubscriptionResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.NotEmpty(t, response.ID)
	assert.Equal(t, models.PlanPremium, response.Plan)
	assert.Equal(t, models.SubscriptionStatusActive, response.Status)
	assert.NotEmpty(t, response.StripeSubscriptionID)
	assert.NotEmpty(t, response.StripeCustomerID)
}

// TestCreateSubscriptionInvalidPlan は無効なプランのテスト
func TestCreateSubscriptionInvalidPlan(t *testing.T) {
	router, _ := setupPaymentTestRouter()

	requestBody := map[string]interface{}{
		"plan":           "invalid_plan",
		"payment_method": "pm_test_123456",
	}

	body, _ := json.Marshal(requestBody)
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/payment/subscription/create", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// TestGetSubscription はサブスクリプション取得のテスト
func TestGetSubscription(t *testing.T) {
	router, _ := setupPaymentTestRouter()

	req, _ := http.NewRequest(http.MethodGet, "/api/v1/payment/subscription", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response models.SubscriptionResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.NotEmpty(t, response.ID)
	assert.Equal(t, models.PlanFree, response.Plan)
	assert.Equal(t, models.SubscriptionStatusActive, response.Status)
}

// TestCancelSubscription はサブスクリプションキャンセルのテスト
func TestCancelSubscription(t *testing.T) {
	tests := []struct {
		name              string
		cancelAtPeriodEnd bool
		expectedStatus    models.SubscriptionStatus
	}{
		{
			name:              "即時キャンセル",
			cancelAtPeriodEnd: false,
			expectedStatus:    models.SubscriptionStatusCanceled,
		},
		{
			name:              "期間終了時にキャンセル",
			cancelAtPeriodEnd: true,
			expectedStatus:    models.SubscriptionStatusActive,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 各テストで新しいルーターを作成
			router, _ := setupPaymentTestRouter()

			requestBody := models.CancelSubscriptionRequest{
				CancelAtPeriodEnd: tt.cancelAtPeriodEnd,
			}

			body, _ := json.Marshal(requestBody)
			req, _ := http.NewRequest(http.MethodPost, "/api/v1/payment/subscription/cancel", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, http.StatusOK, w.Code)

			var response models.SubscriptionResponse
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)

			assert.Equal(t, tt.expectedStatus, response.Status)
			assert.Equal(t, tt.cancelAtPeriodEnd, response.CancelAtPeriodEnd)
		})
	}
}

// TestUpdatePaymentMethod は支払い方法更新のテスト
func TestUpdatePaymentMethod(t *testing.T) {
	router, _ := setupPaymentTestRouter()

	requestBody := models.UpdatePaymentMethodRequest{
		PaymentMethod: "pm_new_123456",
	}

	body, _ := json.Marshal(requestBody)
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/payment/subscription/payment-method", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Equal(t, "Payment method updated successfully", response["message"])
	assert.NotEmpty(t, response["subscription_id"])
}

// TestGetPaymentHistory は支払い履歴取得のテスト
func TestGetPaymentHistory(t *testing.T) {
	router, _ := setupPaymentTestRouter()

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
		{
			name:         "limit=20",
			limitParam:   "20",
			expectedCode: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			url := "/api/v1/payment/history"
			if tt.limitParam != "" {
				url += "?limit=" + tt.limitParam
			}

			req, _ := http.NewRequest(http.MethodGet, url, nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedCode, w.Code)

			var response []*models.PaymentHistoryItem
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)

			// サンプルデータには支払い履歴がないため空配列
			assert.NotNil(t, response)
		})
	}
}

// TestGetPlans はプラン情報取得のテスト
func TestGetPlans(t *testing.T) {
	router, _ := setupPaymentTestRouter()

	req, _ := http.NewRequest(http.MethodGet, "/api/v1/payment/plans", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response []*models.PlanPricing
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	// 3つのプラン（Free, Premium, Yearly）
	assert.Equal(t, 3, len(response))

	// Freeプランの確認
	freeFound := false
	for _, plan := range response {
		if plan.Plan == models.PlanFree {
			freeFound = true
			assert.Equal(t, "Free", plan.Name)
			assert.Equal(t, int64(0), plan.Price)
			assert.Equal(t, "usd", plan.Currency)
			assert.Equal(t, "month", plan.Interval)
			assert.Greater(t, len(plan.Features), 0)
		}
	}
	assert.True(t, freeFound)
}

// TestGetUsage は使用状況取得のテスト
func TestGetUsage(t *testing.T) {
	router, _ := setupPaymentTestRouter()

	req, _ := http.NewRequest(http.MethodGet, "/api/v1/payment/usage", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response models.SubscriptionUsage
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Equal(t, models.PlanFree, response.Plan)
	assert.Equal(t, 1, response.DailyPagesLimit)
	assert.Equal(t, 30, response.DailyMinutesLimit)
	assert.False(t, response.OfflineDownload)
	assert.Equal(t, "standard", response.TTSQualityLevel)
}

// TestPaymentUnauthorized は認証なしのテスト
func TestPaymentUnauthorized(t *testing.T) {
	gin.SetMode(gin.TestMode)

	paymentRepo := repository.NewInMemoryPaymentRepository()
	paymentHandler := NewPaymentHandler(paymentRepo)

	r := gin.New()
	v1 := r.Group("/api/v1")
	paymentHandler.RegisterRoutes(v1)

	tests := []struct {
		name   string
		method string
		url    string
	}{
		{"Create Subscription", http.MethodPost, "/api/v1/payment/subscription/create"},
		{"Get Subscription", http.MethodGet, "/api/v1/payment/subscription"},
		{"Cancel Subscription", http.MethodPost, "/api/v1/payment/subscription/cancel"},
		{"Update Payment Method", http.MethodPost, "/api/v1/payment/subscription/payment-method"},
		{"Get Payment History", http.MethodGet, "/api/v1/payment/history"},
		{"Get Plans", http.MethodGet, "/api/v1/payment/plans"},
		{"Get Usage", http.MethodGet, "/api/v1/payment/usage"},
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

// TestCreateSubscriptionInvalidRequest は無効なリクエストボディのテスト
func TestCreateSubscriptionInvalidRequest(t *testing.T) {
	router, _ := setupPaymentTestRouter()

	req, _ := http.NewRequest(http.MethodPost, "/api/v1/payment/subscription/create", bytes.NewBuffer([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}
