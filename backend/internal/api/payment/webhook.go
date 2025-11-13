package payment

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/stripe/stripe-go/v76/webhook"
)

// VerifyWebhookSignature verifies the Stripe webhook signature
func VerifyWebhookSignature(payload []byte, signature string) error {
	webhookSecret := os.Getenv("STRIPE_WEBHOOK_SECRET")

	// Skip verification if using mocks or no webhook secret
	if os.Getenv("USE_MOCK_APIS") == "true" || webhookSecret == "" {
		return nil
	}

	// Verify the signature
	_, err := webhook.ConstructEvent(payload, signature, webhookSecret)
	if err != nil {
		return fmt.Errorf("webhook signature verification failed: %w", err)
	}

	return nil
}

// HandleWebhookWithVerification handles webhook with signature verification
func (h *Handler) HandleWebhookWithVerification(w http.ResponseWriter, r *http.Request) {
	// Read the request body
	payload, err := io.ReadAll(r.Body)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Failed to read request body")
		return
	}

	// Get the signature from header
	signature := r.Header.Get("Stripe-Signature")
	if signature == "" {
		respondWithError(w, http.StatusBadRequest, "Missing Stripe signature")
		return
	}

	// Verify the signature
	if err := VerifyWebhookSignature(payload, signature); err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid webhook signature")
		return
	}

	// Parse the event
	var event map[string]interface{}
	if err := json.Unmarshal(payload, &event); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid JSON payload")
		return
	}

	// Get event type
	eventType, ok := event["type"].(string)
	if !ok {
		respondWithError(w, http.StatusBadRequest, "Missing event type")
		return
	}

	// Handle the event
	if err := h.service.HandleWebhookEvent(r.Context(), eventType, payload); err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to handle webhook: %v", err))
		return
	}

	// Respond with success
	respondWithJSON(w, http.StatusOK, map[string]string{"status": "success"})
}
