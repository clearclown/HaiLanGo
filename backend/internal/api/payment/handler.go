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
	r.HandleFunc("/api/v1/payment/webhook", h.HandleWebhook).Methods("POST")
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
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	userID, err := uuid.Parse(req.UserID)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	planID, err := uuid.Parse(req.PlanID)
	if err != nil {
		http.Error(w, "Invalid plan ID", http.StatusBadRequest)
		return
	}

	subscription, err := h.service.CreateSubscription(r.Context(), userID, planID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(subscription)
}

// ListPlans handles listing all subscription plans
func (h *Handler) ListPlans(w http.ResponseWriter, r *http.Request) {
	plans, err := h.service.ListPlans(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(plans)
}

// GetSubscription handles retrieving a subscription by ID
func (h *Handler) GetSubscription(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "Invalid subscription ID", http.StatusBadRequest)
		return
	}

	subscription, err := h.service.GetSubscription(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(subscription)
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
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	subscriptionID, err := uuid.Parse(req.SubscriptionID)
	if err != nil {
		http.Error(w, "Invalid subscription ID", http.StatusBadRequest)
		return
	}

	if err := h.service.CancelSubscription(r.Context(), subscriptionID, req.CancelAtPeriodEnd); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "success"})
}

// HandleWebhook handles Stripe webhook events
func (h *Handler) HandleWebhook(w http.ResponseWriter, r *http.Request) {
	// Get event type from header
	eventType := r.Header.Get("Stripe-Event-Type")
	if eventType == "" {
		eventType = r.Header.Get("X-Stripe-Event-Type")
	}

	// Read payload
	var payload []byte
	if _, err := r.Body.Read(payload); err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}

	// Handle event
	if err := h.service.HandleWebhookEvent(r.Context(), eventType, payload); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
