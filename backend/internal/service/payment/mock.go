package payment

import (
	"context"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/clearclown/HaiLanGo/backend/internal/models"
	"github.com/google/uuid"
)

// MockPaymentService is a mock implementation of the Service interface
type MockPaymentService struct {
	subscriptions map[uuid.UUID]*models.Subscription
	plans         map[uuid.UUID]*models.SubscriptionPlan
	mu            sync.RWMutex
}

// NewMockPaymentService creates a new mock payment service
func NewMockPaymentService() *MockPaymentService {
	service := &MockPaymentService{
		subscriptions: make(map[uuid.UUID]*models.Subscription),
		plans:         make(map[uuid.UUID]*models.SubscriptionPlan),
	}

	// Initialize default plans
	service.initializePlans()

	return service
}

// initializePlans initializes the default subscription plans
func (m *MockPaymentService) initializePlans() {
	monthlyPlan := &models.SubscriptionPlan{
		ID:            uuid.New(),
		Name:          "Premium Monthly",
		Description:   "Monthly premium subscription",
		Price:         999,  // $9.99
		Currency:      "usd",
		Interval:      "month",
		StripePriceID: "price_mock_monthly",
		Active:        true,
		Features:      []string{"Unlimited learning", "High-quality TTS", "Offline downloads", "Priority support"},
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	yearlyPlan := &models.SubscriptionPlan{
		ID:            uuid.New(),
		Name:          "Premium Yearly",
		Description:   "Yearly premium subscription",
		Price:         9999, // $99.99
		Currency:      "usd",
		Interval:      "year",
		StripePriceID: "price_mock_yearly",
		Active:        true,
		Features:      []string{"Unlimited learning", "High-quality TTS", "Offline downloads", "Priority support", "20% savings"},
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	m.plans[monthlyPlan.ID] = monthlyPlan
	m.plans[yearlyPlan.ID] = yearlyPlan
}

// CreateSubscription creates a new mock subscription
func (m *MockPaymentService) CreateSubscription(ctx context.Context, userID, planID uuid.UUID) (*models.Subscription, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Verify plan exists
	plan, ok := m.plans[planID]
	if !ok {
		return nil, fmt.Errorf("plan not found")
	}

	// Create subscription
	subscription := &models.Subscription{
		ID:                   uuid.New(),
		UserID:               userID,
		PlanID:               planID,
		StripeCustomerID:     fmt.Sprintf("cus_mock_%s", uuid.New().String()[:8]),
		StripeSubscriptionID: fmt.Sprintf("sub_mock_%s", uuid.New().String()[:8]),
		Status:               models.SubscriptionStatusActive,
		CurrentPeriodStart:   time.Now(),
		CurrentPeriodEnd:     calculatePeriodEnd(time.Now(), plan.Interval),
		CancelAtPeriodEnd:    false,
		CreatedAt:            time.Now(),
		UpdatedAt:            time.Now(),
	}

	m.subscriptions[subscription.ID] = subscription

	return subscription, nil
}

// GetSubscription retrieves a mock subscription by ID
func (m *MockPaymentService) GetSubscription(ctx context.Context, subscriptionID uuid.UUID) (*models.Subscription, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	subscription, ok := m.subscriptions[subscriptionID]
	if !ok {
		return nil, fmt.Errorf("subscription not found")
	}

	return subscription, nil
}

// GetUserSubscription retrieves a user's active mock subscription
func (m *MockPaymentService) GetUserSubscription(ctx context.Context, userID uuid.UUID) (*models.Subscription, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	for _, sub := range m.subscriptions {
		if sub.UserID == userID && sub.Status == models.SubscriptionStatusActive {
			return sub, nil
		}
	}

	return nil, fmt.Errorf("no active subscription found")
}

// UpdateSubscription updates a mock subscription
func (m *MockPaymentService) UpdateSubscription(ctx context.Context, subscription *models.Subscription) (*models.Subscription, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, ok := m.subscriptions[subscription.ID]; !ok {
		return nil, fmt.Errorf("subscription not found")
	}

	subscription.UpdatedAt = time.Now()
	m.subscriptions[subscription.ID] = subscription

	return subscription, nil
}

// CancelSubscription cancels a mock subscription
func (m *MockPaymentService) CancelSubscription(ctx context.Context, subscriptionID uuid.UUID, cancelAtPeriodEnd bool) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	subscription, ok := m.subscriptions[subscriptionID]
	if !ok {
		return fmt.Errorf("subscription not found")
	}

	if cancelAtPeriodEnd {
		subscription.CancelAtPeriodEnd = true
	} else {
		subscription.Status = models.SubscriptionStatusCanceled
		now := time.Now()
		subscription.CanceledAt = &now
	}

	subscription.UpdatedAt = time.Now()
	m.subscriptions[subscriptionID] = subscription

	return nil
}

// ListPlans lists all mock subscription plans
func (m *MockPaymentService) ListPlans(ctx context.Context) ([]*models.SubscriptionPlan, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	plans := make([]*models.SubscriptionPlan, 0, len(m.plans))
	for _, plan := range m.plans {
		if plan.Active {
			plans = append(plans, plan)
		}
	}

	return plans, nil
}

// GetPlan retrieves a mock plan by ID
func (m *MockPaymentService) GetPlan(ctx context.Context, planID uuid.UUID) (*models.SubscriptionPlan, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	plan, ok := m.plans[planID]
	if !ok {
		return nil, fmt.Errorf("plan not found")
	}

	return plan, nil
}

// HandleWebhookEvent handles mock webhook events
func (m *MockPaymentService) HandleWebhookEvent(ctx context.Context, eventType string, payload []byte) error {
	// Mock webhook handling - always succeeds
	return nil
}

// shouldUseMocks checks if mocks should be used
func shouldUseMocks() bool {
	return os.Getenv("USE_MOCK_APIS") == "true" ||
		os.Getenv("TEST_USE_MOCKS") == "true" ||
		os.Getenv("STRIPE_SECRET_KEY") == ""
}

// calculatePeriodEnd calculates the end of the billing period
func calculatePeriodEnd(start time.Time, interval string) time.Time {
	switch interval {
	case "month":
		return start.AddDate(0, 1, 0)
	case "year":
		return start.AddDate(1, 0, 0)
	default:
		return start.AddDate(0, 1, 0)
	}
}
