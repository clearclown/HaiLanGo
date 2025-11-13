package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/clearclown/HaiLanGo/backend/internal/models"
	"github.com/google/uuid"
)

// SubscriptionRepository defines the interface for subscription data access
type SubscriptionRepository interface {
	// Create creates a new subscription
	Create(ctx context.Context, subscription *models.Subscription) error

	// GetByID retrieves a subscription by ID
	GetByID(ctx context.Context, id uuid.UUID) (*models.Subscription, error)

	// GetByUserID retrieves a user's active subscription
	GetByUserID(ctx context.Context, userID uuid.UUID) (*models.Subscription, error)

	// Update updates a subscription
	Update(ctx context.Context, subscription *models.Subscription) error

	// Delete deletes a subscription
	Delete(ctx context.Context, id uuid.UUID) error

	// ListByUserID lists all subscriptions for a user
	ListByUserID(ctx context.Context, userID uuid.UUID) ([]*models.Subscription, error)
}

// PostgresSubscriptionRepository implements SubscriptionRepository using PostgreSQL
type PostgresSubscriptionRepository struct {
	db *sql.DB
}

// NewPostgresSubscriptionRepository creates a new PostgreSQL subscription repository
func NewPostgresSubscriptionRepository(db *sql.DB) *PostgresSubscriptionRepository {
	return &PostgresSubscriptionRepository{db: db}
}

// Create creates a new subscription
func (r *PostgresSubscriptionRepository) Create(ctx context.Context, subscription *models.Subscription) error {
	query := `
		INSERT INTO subscriptions (
			id, user_id, plan_id, stripe_customer_id, stripe_subscription_id,
			status, current_period_start, current_period_end, cancel_at_period_end,
			created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	`

	_, err := r.db.ExecContext(
		ctx,
		query,
		subscription.ID,
		subscription.UserID,
		subscription.PlanID,
		subscription.StripeCustomerID,
		subscription.StripeSubscriptionID,
		subscription.Status,
		subscription.CurrentPeriodStart,
		subscription.CurrentPeriodEnd,
		subscription.CancelAtPeriodEnd,
		subscription.CreatedAt,
		subscription.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to create subscription: %w", err)
	}

	return nil
}

// GetByID retrieves a subscription by ID
func (r *PostgresSubscriptionRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Subscription, error) {
	query := `
		SELECT
			id, user_id, plan_id, stripe_customer_id, stripe_subscription_id,
			status, current_period_start, current_period_end, cancel_at_period_end,
			canceled_at, created_at, updated_at
		FROM subscriptions
		WHERE id = $1
	`

	subscription := &models.Subscription{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&subscription.ID,
		&subscription.UserID,
		&subscription.PlanID,
		&subscription.StripeCustomerID,
		&subscription.StripeSubscriptionID,
		&subscription.Status,
		&subscription.CurrentPeriodStart,
		&subscription.CurrentPeriodEnd,
		&subscription.CancelAtPeriodEnd,
		&subscription.CanceledAt,
		&subscription.CreatedAt,
		&subscription.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("subscription not found")
	}

	if err != nil {
		return nil, fmt.Errorf("failed to get subscription: %w", err)
	}

	return subscription, nil
}

// GetByUserID retrieves a user's active subscription
func (r *PostgresSubscriptionRepository) GetByUserID(ctx context.Context, userID uuid.UUID) (*models.Subscription, error) {
	query := `
		SELECT
			id, user_id, plan_id, stripe_customer_id, stripe_subscription_id,
			status, current_period_start, current_period_end, cancel_at_period_end,
			canceled_at, created_at, updated_at
		FROM subscriptions
		WHERE user_id = $1 AND status = 'active'
		ORDER BY created_at DESC
		LIMIT 1
	`

	subscription := &models.Subscription{}
	err := r.db.QueryRowContext(ctx, query, userID).Scan(
		&subscription.ID,
		&subscription.UserID,
		&subscription.PlanID,
		&subscription.StripeCustomerID,
		&subscription.StripeSubscriptionID,
		&subscription.Status,
		&subscription.CurrentPeriodStart,
		&subscription.CurrentPeriodEnd,
		&subscription.CancelAtPeriodEnd,
		&subscription.CanceledAt,
		&subscription.CreatedAt,
		&subscription.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("no active subscription found")
	}

	if err != nil {
		return nil, fmt.Errorf("failed to get subscription: %w", err)
	}

	return subscription, nil
}

// Update updates a subscription
func (r *PostgresSubscriptionRepository) Update(ctx context.Context, subscription *models.Subscription) error {
	query := `
		UPDATE subscriptions
		SET
			status = $2,
			current_period_start = $3,
			current_period_end = $4,
			cancel_at_period_end = $5,
			canceled_at = $6,
			updated_at = $7
		WHERE id = $1
	`

	_, err := r.db.ExecContext(
		ctx,
		query,
		subscription.ID,
		subscription.Status,
		subscription.CurrentPeriodStart,
		subscription.CurrentPeriodEnd,
		subscription.CancelAtPeriodEnd,
		subscription.CanceledAt,
		subscription.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to update subscription: %w", err)
	}

	return nil
}

// Delete deletes a subscription
func (r *PostgresSubscriptionRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM subscriptions WHERE id = $1`

	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete subscription: %w", err)
	}

	return nil
}

// ListByUserID lists all subscriptions for a user
func (r *PostgresSubscriptionRepository) ListByUserID(ctx context.Context, userID uuid.UUID) ([]*models.Subscription, error) {
	query := `
		SELECT
			id, user_id, plan_id, stripe_customer_id, stripe_subscription_id,
			status, current_period_start, current_period_end, cancel_at_period_end,
			canceled_at, created_at, updated_at
		FROM subscriptions
		WHERE user_id = $1
		ORDER BY created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to list subscriptions: %w", err)
	}
	defer rows.Close()

	subscriptions := []*models.Subscription{}
	for rows.Next() {
		subscription := &models.Subscription{}
		err := rows.Scan(
			&subscription.ID,
			&subscription.UserID,
			&subscription.PlanID,
			&subscription.StripeCustomerID,
			&subscription.StripeSubscriptionID,
			&subscription.Status,
			&subscription.CurrentPeriodStart,
			&subscription.CurrentPeriodEnd,
			&subscription.CancelAtPeriodEnd,
			&subscription.CanceledAt,
			&subscription.CreatedAt,
			&subscription.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan subscription: %w", err)
		}
		subscriptions = append(subscriptions, subscription)
	}

	return subscriptions, nil
}

// PlanRepository defines the interface for plan data access
type PlanRepository interface {
	// GetByID retrieves a plan by ID
	GetByID(ctx context.Context, id uuid.UUID) (*models.SubscriptionPlan, error)

	// List lists all active plans
	List(ctx context.Context) ([]*models.SubscriptionPlan, error)

	// Create creates a new plan
	Create(ctx context.Context, plan *models.SubscriptionPlan) error

	// Update updates a plan
	Update(ctx context.Context, plan *models.SubscriptionPlan) error
}

// PostgresPlanRepository implements PlanRepository using PostgreSQL
type PostgresPlanRepository struct {
	db *sql.DB
}

// NewPostgresPlanRepository creates a new PostgreSQL plan repository
func NewPostgresPlanRepository(db *sql.DB) *PostgresPlanRepository {
	return &PostgresPlanRepository{db: db}
}

// GetByID retrieves a plan by ID
func (r *PostgresPlanRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.SubscriptionPlan, error) {
	query := `
		SELECT
			id, name, description, price, currency, interval,
			stripe_price_id, active, created_at, updated_at
		FROM subscription_plans
		WHERE id = $1
	`

	plan := &models.SubscriptionPlan{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&plan.ID,
		&plan.Name,
		&plan.Description,
		&plan.Price,
		&plan.Currency,
		&plan.Interval,
		&plan.StripePriceID,
		&plan.Active,
		&plan.CreatedAt,
		&plan.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("plan not found")
	}

	if err != nil {
		return nil, fmt.Errorf("failed to get plan: %w", err)
	}

	return plan, nil
}

// List lists all active plans
func (r *PostgresPlanRepository) List(ctx context.Context) ([]*models.SubscriptionPlan, error) {
	query := `
		SELECT
			id, name, description, price, currency, interval,
			stripe_price_id, active, created_at, updated_at
		FROM subscription_plans
		WHERE active = true
		ORDER BY price ASC
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to list plans: %w", err)
	}
	defer rows.Close()

	plans := []*models.SubscriptionPlan{}
	for rows.Next() {
		plan := &models.SubscriptionPlan{}
		err := rows.Scan(
			&plan.ID,
			&plan.Name,
			&plan.Description,
			&plan.Price,
			&plan.Currency,
			&plan.Interval,
			&plan.StripePriceID,
			&plan.Active,
			&plan.CreatedAt,
			&plan.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan plan: %w", err)
		}
		plans = append(plans, plan)
	}

	return plans, nil
}

// Create creates a new plan
func (r *PostgresPlanRepository) Create(ctx context.Context, plan *models.SubscriptionPlan) error {
	query := `
		INSERT INTO subscription_plans (
			id, name, description, price, currency, interval,
			stripe_price_id, active, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`

	_, err := r.db.ExecContext(
		ctx,
		query,
		plan.ID,
		plan.Name,
		plan.Description,
		plan.Price,
		plan.Currency,
		plan.Interval,
		plan.StripePriceID,
		plan.Active,
		plan.CreatedAt,
		plan.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to create plan: %w", err)
	}

	return nil
}

// Update updates a plan
func (r *PostgresPlanRepository) Update(ctx context.Context, plan *models.SubscriptionPlan) error {
	query := `
		UPDATE subscription_plans
		SET
			name = $2,
			description = $3,
			price = $4,
			active = $5,
			updated_at = $6
		WHERE id = $1
	`

	_, err := r.db.ExecContext(
		ctx,
		query,
		plan.ID,
		plan.Name,
		plan.Description,
		plan.Price,
		plan.Active,
		plan.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to update plan: %w", err)
	}

	return nil
}
