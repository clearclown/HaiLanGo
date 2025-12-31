package payment

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/clearclown/HaiLanGo/backend/internal/models"
	"github.com/clearclown/HaiLanGo/backend/internal/repository"
	"github.com/google/uuid"
	"github.com/stripe/stripe-go/v76"
	"github.com/stripe/stripe-go/v76/customer"
	"github.com/stripe/stripe-go/v76/price"
	"github.com/stripe/stripe-go/v76/subscription"
)

// StripeService implements the Service interface using Stripe API
type StripeService struct {
	secretKey string
	repo      repository.PaymentRepositoryInterface
}

// NewStripeService creates a new Stripe payment service
func NewStripeService() *StripeService {
	return NewStripeServiceWithRepo(nil)
}

// NewStripeServiceWithRepo creates a new Stripe payment service with a repository
func NewStripeServiceWithRepo(repo repository.PaymentRepositoryInterface) *StripeService {
	secretKey := os.Getenv("STRIPE_SECRET_KEY")
	stripe.Key = secretKey

	// Use InMemory repository if none provided
	if repo == nil {
		repo = repository.NewInMemoryPaymentRepository()
	}

	return &StripeService{
		secretKey: secretKey,
		repo:      repo,
	}
}

// CreateSubscription creates a new subscription for a user
func (s *StripeService) CreateSubscription(ctx context.Context, userID, planID uuid.UUID) (*models.Subscription, error) {
	// Get plan details
	plan, err := s.GetPlan(ctx, planID)
	if err != nil {
		return nil, fmt.Errorf("failed to get plan: %w", err)
	}

	// Determine plan type from plan
	planType := models.PlanTypePremium
	if plan.Interval == "year" {
		planType = models.PlanTypeYearly
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

	// Save to repository
	repoSub, err := s.repo.CreateSubscription(ctx, userID, planType, sub.ID, cust.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to save subscription: %w", err)
	}

	// Update with Stripe data
	repoSub.PlanID = planID
	repoSub.Status = models.SubscriptionStatus(sub.Status)
	repoSub.CurrentPeriodStart = time.Unix(sub.CurrentPeriodStart, 0)
	repoSub.CurrentPeriodEnd = time.Unix(sub.CurrentPeriodEnd, 0)
	repoSub.CancelAtPeriodEnd = sub.CancelAtPeriodEnd

	return repoSub, nil
}

// GetSubscription retrieves a subscription by ID
func (s *StripeService) GetSubscription(ctx context.Context, subscriptionID uuid.UUID) (*models.Subscription, error) {
	return s.repo.GetSubscription(ctx, subscriptionID)
}

// GetUserSubscription retrieves a user's active subscription
func (s *StripeService) GetUserSubscription(ctx context.Context, userID uuid.UUID) (*models.Subscription, error) {
	return s.repo.GetSubscriptionByUserID(ctx, userID)
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
	// Get plan pricing from repository
	pricing, err := s.repo.GetPlanPricing(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get plan pricing: %w", err)
	}

	// Find matching plan by ID or return default premium plan
	for _, p := range pricing {
		plan := &models.SubscriptionPlan{
			ID:            planID,
			Name:          p.Name,
			Description:   p.Description,
			PlanType:      p.Plan,
			Price:         p.Price,
			Currency:      p.Currency,
			Interval:      p.Interval,
			StripePriceID: fmt.Sprintf("price_%s", p.Plan),
			Active:        true,
			Features:      p.Features,
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
		}

		// Match by plan type (simplified matching)
		if p.Plan == models.PlanTypePremium || p.Plan == models.PlanTypeYearly {
			return plan, nil
		}
	}

	// Return default premium plan if not found
	return &models.SubscriptionPlan{
		ID:            planID,
		Name:          "Premium",
		Description:   "Premium subscription plan",
		PlanType:      models.PlanTypePremium,
		Price:         999,
		Currency:      "usd",
		Interval:      "month",
		StripePriceID: "price_premium_monthly",
		Active:        true,
		Features:      []string{"Unlimited learning", "Premium TTS", "Offline download"},
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}, nil
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

// StripeSubscriptionData represents Stripe subscription webhook data
type StripeSubscriptionData struct {
	Object struct {
		ID                string `json:"id"`
		Status            string `json:"status"`
		CustomerID        string `json:"customer"`
		CancelAtPeriodEnd bool   `json:"cancel_at_period_end"`
		Metadata          struct {
			UserID string `json:"user_id"`
			PlanID string `json:"plan_id"`
		} `json:"metadata"`
		CurrentPeriodStart int64 `json:"current_period_start"`
		CurrentPeriodEnd   int64 `json:"current_period_end"`
	} `json:"object"`
}

// StripePaymentData represents Stripe payment intent webhook data
type StripePaymentData struct {
	Object struct {
		ID       string `json:"id"`
		Amount   int64  `json:"amount"`
		Currency string `json:"currency"`
		Status   string `json:"status"`
		Metadata struct {
			UserID         string `json:"user_id"`
			SubscriptionID string `json:"subscription_id"`
		} `json:"metadata"`
	} `json:"object"`
}

func (s *StripeService) handleSubscriptionEvent(ctx context.Context, eventType string, payload []byte) error {
	var eventData struct {
		Data StripeSubscriptionData `json:"data"`
	}
	if err := json.Unmarshal(payload, &eventData); err != nil {
		return fmt.Errorf("failed to parse subscription event: %w", err)
	}

	subData := eventData.Data.Object

	// Parse user ID from metadata
	if subData.Metadata.UserID == "" {
		// If no user ID in metadata, skip (might be a test event)
		return nil
	}

	userID, err := uuid.Parse(subData.Metadata.UserID)
	if err != nil {
		return fmt.Errorf("invalid user_id in metadata: %w", err)
	}

	// Get existing subscription by user ID
	existingSub, err := s.repo.GetSubscriptionByUserID(ctx, userID)
	if err != nil {
		// Subscription might not exist yet for "created" events
		if eventType == "customer.subscription.created" {
			// Create subscription in repository
			planType := models.PlanTypePremium
			_, err = s.repo.CreateSubscription(ctx, userID, planType, subData.ID, subData.CustomerID)
			if err != nil {
				return fmt.Errorf("failed to create subscription: %w", err)
			}
			return nil
		}
		return fmt.Errorf("subscription not found for user: %w", err)
	}

	// Update subscription status based on event type
	switch eventType {
	case "customer.subscription.updated":
		status := models.SubscriptionStatus(subData.Status)
		if err := s.repo.UpdateSubscriptionStatus(ctx, existingSub.ID, status); err != nil {
			return fmt.Errorf("failed to update subscription status: %w", err)
		}

		// Update cancel at period end flag
		if subData.CancelAtPeriodEnd != existingSub.CancelAtPeriodEnd {
			if err := s.repo.CancelSubscription(ctx, existingSub.ID, subData.CancelAtPeriodEnd); err != nil {
				return fmt.Errorf("failed to update cancel_at_period_end: %w", err)
			}
		}

	case "customer.subscription.deleted":
		if err := s.repo.UpdateSubscriptionStatus(ctx, existingSub.ID, models.SubscriptionStatusCanceled); err != nil {
			return fmt.Errorf("failed to cancel subscription: %w", err)
		}
	}

	return nil
}

func (s *StripeService) handlePaymentEvent(ctx context.Context, eventType string, payload []byte) error {
	var eventData struct {
		Data StripePaymentData `json:"data"`
	}
	if err := json.Unmarshal(payload, &eventData); err != nil {
		return fmt.Errorf("failed to parse payment event: %w", err)
	}

	paymentData := eventData.Data.Object

	// Parse user ID from metadata
	if paymentData.Metadata.UserID == "" {
		// If no user ID in metadata, skip
		return nil
	}

	userID, err := uuid.Parse(paymentData.Metadata.UserID)
	if err != nil {
		return fmt.Errorf("invalid user_id in metadata: %w", err)
	}

	// Get subscription for this user
	sub, err := s.repo.GetSubscriptionByUserID(ctx, userID)
	if err != nil {
		return fmt.Errorf("subscription not found for user: %w", err)
	}

	switch eventType {
	case "payment_intent.succeeded":
		// Create payment record
		_, err = s.repo.CreatePayment(ctx, userID, sub.ID, paymentData.Amount, paymentData.Currency, paymentData.ID)
		if err != nil {
			return fmt.Errorf("failed to create payment record: %w", err)
		}

	case "payment_intent.payment_failed":
		// Update subscription status to past_due
		if err := s.repo.UpdateSubscriptionStatus(ctx, sub.ID, models.SubscriptionStatusPastDue); err != nil {
			return fmt.Errorf("failed to update subscription status: %w", err)
		}
	}

	return nil
}
