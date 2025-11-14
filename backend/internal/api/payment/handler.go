package payment

import (
	"net/http"

	"github.com/clearclown/HaiLanGo/backend/internal/service/payment"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Handler handles payment-related HTTP requests
type Handler struct {
	service payment.Service
}

// NewHandler creates a new payment handler
func NewHandler(service payment.Service) *Handler {
	return &Handler{
		service: service,
	}
}

// CreateSubscriptionRequest represents the request body for creating a subscription
type CreateSubscriptionRequest struct {
	UserID string `json:"user_id"`
	PlanID string `json:"plan_id"`
}

// CreateSubscription handles subscription creation requests
func (h *Handler) CreateSubscription(c *gin.Context) {
	var req CreateSubscriptionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	userID, err := uuid.Parse(req.UserID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	planID, err := uuid.Parse(req.PlanID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid plan ID"})
		return
	}

	subscription, err := h.service.CreateSubscription(c.Request.Context(), userID, planID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, subscription)
}

// ListPlans handles listing all subscription plans
func (h *Handler) ListPlans(c *gin.Context) {
	plans, err := h.service.ListPlans(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, plans)
}

// GetSubscription handles retrieving a subscription by ID
func (h *Handler) GetSubscription(c *gin.Context) {
	idStr := c.Param("id")

	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid subscription ID"})
		return
	}

	subscription, err := h.service.GetSubscription(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, subscription)
}

// CancelSubscriptionRequest represents the request body for canceling a subscription
type CancelSubscriptionRequest struct {
	SubscriptionID    string `json:"subscription_id"`
	CancelAtPeriodEnd bool   `json:"cancel_at_period_end"`
}

// CancelSubscription handles subscription cancellation requests
func (h *Handler) CancelSubscription(c *gin.Context) {
	var req CancelSubscriptionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	subscriptionID, err := uuid.Parse(req.SubscriptionID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid subscription ID"})
		return
	}

	if err := h.service.CancelSubscription(c.Request.Context(), subscriptionID, req.CancelAtPeriodEnd); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Subscription canceled successfully"})
}

// RegisterRoutes registers payment routes
func (h *Handler) RegisterRoutes(rg *gin.RouterGroup) {
	payment := rg.Group("/payment")
	{
		payment.POST("/subscribe", h.CreateSubscription)
		payment.GET("/plans", h.ListPlans)
		payment.GET("/subscription/:id", h.GetSubscription)
		payment.POST("/cancel", h.CancelSubscription)
		// Note: Webhook handler should be added separately without auth middleware
	}
}

