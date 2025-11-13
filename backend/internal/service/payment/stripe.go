package payment

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/clearclown/HaiLanGo/backend/internal/models"
	"github.com/google/uuid"
	"github.com/stripe/stripe-go/v76"
	"github.com/stripe/stripe-go/v76/customer"
	"github.com/stripe/stripe-go/v76/price"
	"github.com/stripe/stripe-go/v76/subscription"
)

// StripeService implements the Service interface using Stripe API
type StripeService struct {
	secretKey string
}

// NewStripeService creates a new Stripe payment service
func NewStripeService() *StripeService {
	secretKey := os.Getenv("STRIPE_SECRET_KEY")
	stripe.Key = secretKey

	return &StripeService{
		secretKey: secretKey,
	}
}

// CreateSubscription creates a new subscription for a user
func (s *StripeService) CreateSubscription(ctx context.Context, userID, planID uuid.UUID) (*models.Subscription, error) {
	// Get plan details
	plan, err := s.GetPlan(ctx, planID)
	if err != nil {
		return nil, fmt.Errorf("failed to get plan: %w", err)
	}

	// Create or get Stripe customer
	customerParams := &stripe.CustomerParams{
		Metadata: map[string]string{
			"user_id": userID.String(),
		},
	}
	cust, err := customer.New(customerParams)
	if err != nil {
		return nil, fmt.Errorf("failed to create customer: %w", err)
	}

	// Create subscription
	subParams := &stripe.SubscriptionParams{
		Customer: stripe.String(cust.ID),
		Items: []*stripe.SubscriptionItemsParams{
			{
				Price: stripe.String(plan.StripePriceID),
			},
		},
		Metadata: map[string]string{
			"user_id": userID.String(),
			"plan_id": planID.String(),
		},
	}
	sub, err := subscription.New(subParams)
	if err != nil {
		return nil, fmt.Errorf("failed to create subscription: %w", err)
	}

	// Convert to internal model
	return &models.Subscription{
		ID:                   uuid.New(),
		UserID:               userID,
		PlanID:               planID,
		StripeCustomerID:     cust.ID,
		StripeSubscriptionID: sub.ID,
		Status:               models.SubscriptionStatus(sub.Status),
		CurrentPeriodStart:   time.Unix(sub.CurrentPeriodStart, 0),
		CurrentPeriodEnd:     time.Unix(sub.CurrentPeriodEnd, 0),
		CancelAtPeriodEnd:    sub.CancelAtPeriodEnd,
		CreatedAt:            time.Now(),
		UpdatedAt:            time.Now(),
	}, nil
}

// GetSubscription retrieves a subscription by ID
func (s *StripeService) GetSubscription(ctx context.Context, subscriptionID uuid.UUID) (*models.Subscription, error) {
	// This would typically query the database
	// For now, return a placeholder
	return nil, fmt.Errorf("not implemented")
}

// GetUserSubscription retrieves a user's active subscription
func (s *StripeService) GetUserSubscription(ctx context.Context, userID uuid.UUID) (*models.Subscription, error) {
	// This would typically query the database
	// For now, return a placeholder
	return nil, fmt.Errorf("not implemented")
}

// UpdateSubscription updates a subscription
func (s *StripeService) UpdateSubscription(ctx context.Context, sub *models.Subscription) (*models.Subscription, error) {
	// Update Stripe subscription
	params := &stripe.SubscriptionParams{
		CancelAtPeriodEnd: stripe.Bool(sub.CancelAtPeriodEnd),
	}
	_, err := subscription.Update(sub.StripeSubscriptionID, params)
	if err != nil {
		return nil, fmt.Errorf("failed to update subscription: %w", err)
	}

	sub.UpdatedAt = time.Now()
	return sub, nil
}

// CancelSubscription cancels a subscription
func (s *StripeService) CancelSubscription(ctx context.Context, subscriptionID uuid.UUID, cancelAtPeriodEnd bool) error {
	// Get subscription from database
	sub, err := s.GetSubscription(ctx, subscriptionID)
	if err != nil {
		return err
	}

	if cancelAtPeriodEnd {
		// Cancel at period end
		params := &stripe.SubscriptionParams{
			CancelAtPeriodEnd: stripe.Bool(true),
		}
		_, err = subscription.Update(sub.StripeSubscriptionID, params)
	} else {
		// Cancel immediately
		_, err = subscription.Cancel(sub.StripeSubscriptionID, nil)
	}

	if err != nil {
		return fmt.Errorf("failed to cancel subscription: %w", err)
	}

	return nil
}

// ListPlans lists all available subscription plans
func (s *StripeService) ListPlans(ctx context.Context) ([]*models.SubscriptionPlan, error) {
	// List prices from Stripe
	params := &stripe.PriceListParams{}
	params.Filters.AddFilter("active", "", "true")

	i := price.List(params)
	plans := []*models.SubscriptionPlan{}

	for i.Next() {
		p := i.Price()
		plan := &models.SubscriptionPlan{
			ID:            uuid.New(),
			Name:          fmt.Sprintf("Premium %s", p.Recurring.Interval),
			Price:         p.UnitAmount,
			Currency:      string(p.Currency),
			Interval:      string(p.Recurring.Interval),
			StripePriceID: p.ID,
			Active:        p.Active,
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
		}
		plans = append(plans, plan)
	}

	return plans, nil
}

// GetPlan retrieves a plan by ID
func (s *StripeService) GetPlan(ctx context.Context, planID uuid.UUID) (*models.SubscriptionPlan, error) {
	// This would typically query the database
	// For now, return a placeholder
	return nil, fmt.Errorf("not implemented")
}

// HandleWebhookEvent handles Stripe webhook events
func (s *StripeService) HandleWebhookEvent(ctx context.Context, eventType string, payload []byte) error {
	// Handle different webhook events
	switch eventType {
	case "customer.subscription.created",
		"customer.subscription.updated",
		"customer.subscription.deleted":
		// Handle subscription events
		return s.handleSubscriptionEvent(ctx, eventType, payload)

	case "payment_intent.succeeded",
		"payment_intent.payment_failed":
		// Handle payment events
		return s.handlePaymentEvent(ctx, eventType, payload)

	default:
		// Log unhandled event
		return nil
	}
}

func (s *StripeService) handleSubscriptionEvent(ctx context.Context, eventType string, payload []byte) error {
	// Parse and handle subscription events
	return nil
}

func (s *StripeService) handlePaymentEvent(ctx context.Context, eventType string, payload []byte) error {
	// Parse and handle payment events
	return nil
}
