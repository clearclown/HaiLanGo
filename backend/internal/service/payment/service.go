package payment

import (
	"context"

	"github.com/clearclown/HaiLanGo/backend/internal/models"
	"github.com/google/uuid"
)

// Service defines the payment service interface
type Service interface {
	// CreateSubscription creates a new subscription for a user
	CreateSubscription(ctx context.Context, userID, planID uuid.UUID) (*models.Subscription, error)

	// GetSubscription retrieves a subscription by ID
	GetSubscription(ctx context.Context, subscriptionID uuid.UUID) (*models.Subscription, error)

	// GetUserSubscription retrieves a user's active subscription
	GetUserSubscription(ctx context.Context, userID uuid.UUID) (*models.Subscription, error)

	// UpdateSubscription updates a subscription
	UpdateSubscription(ctx context.Context, subscription *models.Subscription) (*models.Subscription, error)

	// CancelSubscription cancels a subscription
	CancelSubscription(ctx context.Context, subscriptionID uuid.UUID, cancelAtPeriodEnd bool) error

	// ListPlans lists all available subscription plans
	ListPlans(ctx context.Context) ([]*models.SubscriptionPlan, error)

	// GetPlan retrieves a plan by ID
	GetPlan(ctx context.Context, planID uuid.UUID) (*models.SubscriptionPlan, error)

	// HandleWebhookEvent handles Stripe webhook events
	HandleWebhookEvent(ctx context.Context, eventType string, payload []byte) error
}

// NewPaymentService creates a new payment service
func NewPaymentService() Service {
	// Check if we should use mocks
	if shouldUseMocks() {
		return NewMockPaymentService()
	}

	// Return real Stripe service
	return NewStripeService()
}
