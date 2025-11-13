package models

import (
	"time"

	"github.com/google/uuid"
)

// SubscriptionPlan represents a subscription plan
type SubscriptionPlan struct {
	ID          uuid.UUID `json:"id" db:"id"`
	Name        string    `json:"name" db:"name"`                 // e.g., "Premium Monthly", "Premium Yearly"
	Description string    `json:"description" db:"description"`   // Plan description
	Price       int64     `json:"price" db:"price"`               // Price in cents
	Currency    string    `json:"currency" db:"currency"`         // e.g., "usd"
	Interval    string    `json:"interval" db:"interval"`         // "month" or "year"
	StripePriceID string  `json:"stripe_price_id" db:"stripe_price_id"` // Stripe Price ID
	Active      bool      `json:"active" db:"active"`             // Is this plan active?
	Features    []string  `json:"features" db:"features"`         // List of features
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

// SubscriptionStatus represents the status of a subscription
type SubscriptionStatus string

const (
	SubscriptionStatusActive   SubscriptionStatus = "active"
	SubscriptionStatusCanceled SubscriptionStatus = "canceled"
	SubscriptionStatusPastDue  SubscriptionStatus = "past_due"
	SubscriptionStatusTrialing SubscriptionStatus = "trialing"
	SubscriptionStatusUnpaid   SubscriptionStatus = "unpaid"
)

// Subscription represents a user's subscription
type Subscription struct {
	ID                   uuid.UUID          `json:"id" db:"id"`
	UserID               uuid.UUID          `json:"user_id" db:"user_id"`
	PlanID               uuid.UUID          `json:"plan_id" db:"plan_id"`
	StripeCustomerID     string             `json:"stripe_customer_id" db:"stripe_customer_id"`
	StripeSubscriptionID string             `json:"stripe_subscription_id" db:"stripe_subscription_id"`
	Status               SubscriptionStatus `json:"status" db:"status"`
	CurrentPeriodStart   time.Time          `json:"current_period_start" db:"current_period_start"`
	CurrentPeriodEnd     time.Time          `json:"current_period_end" db:"current_period_end"`
	CancelAtPeriodEnd    bool               `json:"cancel_at_period_end" db:"cancel_at_period_end"`
	CanceledAt           *time.Time         `json:"canceled_at,omitempty" db:"canceled_at"`
	CreatedAt            time.Time          `json:"created_at" db:"created_at"`
	UpdatedAt            time.Time          `json:"updated_at" db:"updated_at"`
}

// Payment represents a payment transaction
type Payment struct {
	ID                  uuid.UUID `json:"id" db:"id"`
	UserID              uuid.UUID `json:"user_id" db:"user_id"`
	SubscriptionID      uuid.UUID `json:"subscription_id" db:"subscription_id"`
	StripePaymentIntent string    `json:"stripe_payment_intent" db:"stripe_payment_intent"`
	Amount              int64     `json:"amount" db:"amount"`                   // Amount in cents
	Currency            string    `json:"currency" db:"currency"`               // e.g., "usd"
	Status              string    `json:"status" db:"status"`                   // "succeeded", "failed", "pending"
	Description         string    `json:"description" db:"description"`
	CreatedAt           time.Time `json:"created_at" db:"created_at"`
}
