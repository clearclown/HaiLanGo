package handler

import (
	"net/http"
	"strconv"

	"github.com/clearclown/HaiLanGo/backend/internal/models"
	"github.com/clearclown/HaiLanGo/backend/internal/repository"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// PaymentHandler は決済APIのハンドラー
type PaymentHandler struct {
	repo repository.PaymentRepositoryInterface
}

// NewPaymentHandler は決済ハンドラーを作成
func NewPaymentHandler(repo repository.PaymentRepositoryInterface) *PaymentHandler {
	return &PaymentHandler{
		repo: repo,
	}
}

// RegisterRoutes は決済APIのルートを登録
func (h *PaymentHandler) RegisterRoutes(rg *gin.RouterGroup) {
	payment := rg.Group("/payment")
	{
		// サブスクリプション管理
		payment.POST("/subscription/create", h.CreateSubscription)
		payment.GET("/subscription", h.GetSubscription)
		payment.POST("/subscription/cancel", h.CancelSubscription)
		payment.POST("/subscription/payment-method", h.UpdatePaymentMethod)

		// 支払い履歴
		payment.GET("/history", h.GetPaymentHistory)

		// プラン情報
		payment.GET("/plans", h.GetPlans)

		// 使用状況
		payment.GET("/usage", h.GetUsage)

		// Webhook（認証不要）
		// Note: 本番環境ではWebhook署名検証が必要
	}
}

// CreateSubscription はサブスクリプションを作成
// POST /api/v1/payment/subscription/create
func (h *PaymentHandler) CreateSubscription(c *gin.Context) {
	userIDStr, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	var req models.CreateSubscriptionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// プランの検証
	if req.Plan != models.PlanPremium && req.Plan != models.PlanYearly {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid plan. Must be 'premium' or 'yearly'"})
		return
	}

	// 実際の実装では、ここでStripe APIを呼び出してサブスクリプションを作成
	// 今はダミーのStripe IDを使用
	stripeSubscriptionID := "sub_" + uuid.New().String()
	stripeCustomerID := "cus_" + uuid.New().String()

	subscription, err := h.repo.CreateSubscription(
		c.Request.Context(),
		userID,
		req.Plan,
		stripeSubscriptionID,
		stripeCustomerID,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create subscription"})
		return
	}

	response := &models.SubscriptionResponse{
		ID:                   subscription.ID.String(),
		UserID:               subscription.UserID.String(),
		Plan:                 subscription.Plan,
		Status:               subscription.Status,
		StripeSubscriptionID: subscription.StripeSubscriptionID,
		StripeCustomerID:     subscription.StripeCustomerID,
		CurrentPeriodStart:   subscription.CurrentPeriodStart,
		CurrentPeriodEnd:     subscription.CurrentPeriodEnd,
		CancelAtPeriodEnd:    subscription.CancelAtPeriodEnd,
		CreatedAt:            subscription.CreatedAt,
		UpdatedAt:            subscription.UpdatedAt,
	}

	c.JSON(http.StatusCreated, response)
}

// GetSubscription はサブスクリプション情報を取得
// GET /api/v1/payment/subscription
func (h *PaymentHandler) GetSubscription(c *gin.Context) {
	userIDStr, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	subscription, err := h.repo.GetSubscriptionByUserID(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Subscription not found"})
		return
	}

	response := &models.SubscriptionResponse{
		ID:                   subscription.ID.String(),
		UserID:               subscription.UserID.String(),
		Plan:                 subscription.Plan,
		Status:               subscription.Status,
		StripeSubscriptionID: subscription.StripeSubscriptionID,
		StripeCustomerID:     subscription.StripeCustomerID,
		CurrentPeriodStart:   subscription.CurrentPeriodStart,
		CurrentPeriodEnd:     subscription.CurrentPeriodEnd,
		CancelAtPeriodEnd:    subscription.CancelAtPeriodEnd,
		CreatedAt:            subscription.CreatedAt,
		UpdatedAt:            subscription.UpdatedAt,
	}

	c.JSON(http.StatusOK, response)
}

// CancelSubscription はサブスクリプションをキャンセル
// POST /api/v1/payment/subscription/cancel
func (h *PaymentHandler) CancelSubscription(c *gin.Context) {
	userIDStr, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	var req models.CancelSubscriptionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	subscription, err := h.repo.GetSubscriptionByUserID(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Subscription not found"})
		return
	}

	// 実際の実装では、ここでStripe APIを呼び出してサブスクリプションをキャンセル

	err = h.repo.CancelSubscription(c.Request.Context(), subscription.ID, req.CancelAtPeriodEnd)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to cancel subscription"})
		return
	}

	// 更新されたサブスクリプションを取得
	subscription, _ = h.repo.GetSubscription(c.Request.Context(), subscription.ID)

	response := &models.SubscriptionResponse{
		ID:                   subscription.ID.String(),
		UserID:               subscription.UserID.String(),
		Plan:                 subscription.Plan,
		Status:               subscription.Status,
		StripeSubscriptionID: subscription.StripeSubscriptionID,
		StripeCustomerID:     subscription.StripeCustomerID,
		CurrentPeriodStart:   subscription.CurrentPeriodStart,
		CurrentPeriodEnd:     subscription.CurrentPeriodEnd,
		CancelAtPeriodEnd:    subscription.CancelAtPeriodEnd,
		CreatedAt:            subscription.CreatedAt,
		UpdatedAt:            subscription.UpdatedAt,
	}

	c.JSON(http.StatusOK, response)
}

// UpdatePaymentMethod は支払い方法を更新
// POST /api/v1/payment/subscription/payment-method
func (h *PaymentHandler) UpdatePaymentMethod(c *gin.Context) {
	userIDStr, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	var req models.UpdatePaymentMethodRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	subscription, err := h.repo.GetSubscriptionByUserID(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Subscription not found"})
		return
	}

	// 実際の実装では、ここでStripe APIを呼び出して支払い方法を更新

	c.JSON(http.StatusOK, gin.H{
		"message":        "Payment method updated successfully",
		"subscription_id": subscription.ID.String(),
	})
}

// GetPaymentHistory は支払い履歴を取得
// GET /api/v1/payment/history
func (h *PaymentHandler) GetPaymentHistory(c *gin.Context) {
	userIDStr, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// クエリパラメータからlimitを取得（デフォルト: 10）
	limit := 10
	if limitStr := c.Query("limit"); limitStr != "" {
		if parsedLimit, err := strconv.Atoi(limitStr); err == nil && parsedLimit > 0 {
			limit = parsedLimit
		}
	}

	payments, err := h.repo.GetPaymentHistory(c.Request.Context(), userID, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get payment history"})
		return
	}

	// レスポンスを整形
	items := make([]*models.PaymentHistoryItem, 0, len(payments))
	for _, payment := range payments {
		items = append(items, &models.PaymentHistoryItem{
			ID:              payment.ID.String(),
			Amount:          payment.Amount,
			Currency:        payment.Currency,
			Status:          payment.Status,
			Description:     payment.Description,
			InvoiceURL:      payment.InvoiceURL,
			ReceiptURL:      payment.ReceiptURL,
			StripePaymentID: payment.StripePaymentID,
			CreatedAt:       payment.CreatedAt,
		})
	}

	c.JSON(http.StatusOK, items)
}

// GetPlans はプラン情報を取得
// GET /api/v1/payment/plans
func (h *PaymentHandler) GetPlans(c *gin.Context) {
	_, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	plans, err := h.repo.GetPlanPricing(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get plans"})
		return
	}

	c.JSON(http.StatusOK, plans)
}

// GetUsage は使用状況を取得
// GET /api/v1/payment/usage
func (h *PaymentHandler) GetUsage(c *gin.Context) {
	userIDStr, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	usage, err := h.repo.GetSubscriptionUsage(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Usage not found"})
		return
	}

	c.JSON(http.StatusOK, usage)
}
