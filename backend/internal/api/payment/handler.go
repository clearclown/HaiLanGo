package payment

import (
	"encoding/json"
	"net/http"

	"github.com/clearclown/HaiLanGo/backend/internal/service/payment"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
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

// RegisterRoutes registers payment routes
func (h *Handler) RegisterRoutes(r *mux.Router) {
	r.HandleFunc("/api/v1/payment/subscribe", h.CreateSubscription).Methods("POST")
	r.HandleFunc("/api/v1/payment/plans", h.ListPlans).Methods("GET")
	r.HandleFunc("/api/v1/payment/subscription/{id}", h.GetSubscription).Methods("GET")
	r.HandleFunc("/api/v1/payment/cancel", h.CancelSubscription).Methods("POST")
	r.HandleFunc("/api/v1/payment/webhook", h.HandleWebhookWithVerification).Methods("POST")
}

// CreateSubscriptionRequest represents the request body for creating a subscription
type CreateSubscriptionRequest struct {
	UserID string `json:"user_id"`
	PlanID string `json:"plan_id"`
}

// CreateSubscription handles subscription creation requests
func (h *Handler) CreateSubscription(w http.ResponseWriter, r *http.Request) {
	var req CreateSubscriptionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	userID, err := uuid.Parse(req.UserID)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	planID, err := uuid.Parse(req.PlanID)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid plan ID")
		return
	}

	subscription, err := h.service.CreateSubscription(r.Context(), userID, planID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, subscription)
}

// ListPlans handles listing all subscription plans
func (h *Handler) ListPlans(w http.ResponseWriter, r *http.Request) {
	plans, err := h.service.ListPlans(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, plans)
}

// GetSubscription handles retrieving a subscription by ID
func (h *Handler) GetSubscription(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	id, err := uuid.Parse(idStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid subscription ID")
		return
	}

	subscription, err := h.service.GetSubscription(r.Context(), id)
	if err != nil {
		respondWithError(w, http.StatusNotFound, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, subscription)
}

// CancelSubscriptionRequest represents the request body for canceling a subscription
type CancelSubscriptionRequest struct {
	SubscriptionID    string `json:"subscription_id"`
	CancelAtPeriodEnd bool   `json:"cancel_at_period_end"`
}

// CancelSubscription handles subscription cancellation requests
func (h *Handler) CancelSubscription(w http.ResponseWriter, r *http.Request) {
	var req CancelSubscriptionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	subscriptionID, err := uuid.Parse(req.SubscriptionID)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid subscription ID")
		return
	}

	if err := h.service.CancelSubscription(r.Context(), subscriptionID, req.CancelAtPeriodEnd); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithSuccess(w, nil, "Subscription canceled successfully")
}

