package payment

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/clearclown/HaiLanGo/backend/internal/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMain(m *testing.M) {
	// Ensure tests use mocks
	os.Setenv("TEST_USE_MOCKS", "true")
	os.Setenv("USE_MOCK_APIS", "true")
	code := m.Run()
	os.Exit(code)
}

func TestNewPaymentService(t *testing.T) {
	service := NewPaymentService()
	assert.NotNil(t, service)
}

func TestCreateSubscription(t *testing.T) {
	service := NewPaymentService()
	ctx := context.Background()

	// Get available plans first
	plans, err := service.ListPlans(ctx)
	require.NoError(t, err)
	require.NotEmpty(t, plans)

	userID := uuid.New()
	planID := plans[0].ID

	tests := []struct {
		name    string
		userID  uuid.UUID
		planID  uuid.UUID
		wantErr bool
	}{
		{
			name:    "successful subscription creation",
			userID:  userID,
			planID:  planID,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			subscription, err := service.CreateSubscription(ctx, tt.userID, tt.planID)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.NotNil(t, subscription)
			assert.Equal(t, tt.userID, subscription.UserID)
			assert.Equal(t, tt.planID, subscription.PlanID)
			assert.NotEmpty(t, subscription.StripeCustomerID)
			assert.NotEmpty(t, subscription.StripeSubscriptionID)
			assert.Equal(t, models.SubscriptionStatusActive, subscription.Status)
		})
	}
}

func TestGetSubscription(t *testing.T) {
	service := NewPaymentService()
	ctx := context.Background()

	// Get available plans first
	plans, err := service.ListPlans(ctx)
	require.NoError(t, err)
	require.NotEmpty(t, plans)

	userID := uuid.New()
	planID := plans[0].ID

	// Create a subscription first
	subscription, err := service.CreateSubscription(ctx, userID, planID)
	require.NoError(t, err)

	// Test getting the subscription
	retrieved, err := service.GetSubscription(ctx, subscription.ID)
	require.NoError(t, err)
	assert.NotNil(t, retrieved)
	assert.Equal(t, subscription.ID, retrieved.ID)
	assert.Equal(t, subscription.UserID, retrieved.UserID)
}

func TestGetUserSubscription(t *testing.T) {
	service := NewPaymentService()
	ctx := context.Background()

	// Get available plans first
	plans, err := service.ListPlans(ctx)
	require.NoError(t, err)
	require.NotEmpty(t, plans)

	userID := uuid.New()
	planID := plans[0].ID

	// Create a subscription first
	subscription, err := service.CreateSubscription(ctx, userID, planID)
	require.NoError(t, err)

	// Test getting user's subscription
	retrieved, err := service.GetUserSubscription(ctx, userID)
	require.NoError(t, err)
	assert.NotNil(t, retrieved)
	assert.Equal(t, subscription.UserID, retrieved.UserID)
}

func TestCancelSubscription(t *testing.T) {
	service := NewPaymentService()
	ctx := context.Background()

	// Get available plans first
	plans, err := service.ListPlans(ctx)
	require.NoError(t, err)
	require.NotEmpty(t, plans)

	userID := uuid.New()
	planID := plans[0].ID

	// Create a subscription first
	subscription, err := service.CreateSubscription(ctx, userID, planID)
	require.NoError(t, err)

	// Test canceling the subscription
	err = service.CancelSubscription(ctx, subscription.ID, true)
	require.NoError(t, err)

	// Verify cancellation
	canceled, err := service.GetSubscription(ctx, subscription.ID)
	require.NoError(t, err)
	assert.True(t, canceled.CancelAtPeriodEnd)
}

func TestListPlans(t *testing.T) {
	service := NewPaymentService()
	ctx := context.Background()

	plans, err := service.ListPlans(ctx)
	require.NoError(t, err)
	assert.NotNil(t, plans)
	assert.GreaterOrEqual(t, len(plans), 2) // At least monthly and yearly plans
}

func TestGetPlan(t *testing.T) {
	service := NewPaymentService()
	ctx := context.Background()

	// Get all plans
	plans, err := service.ListPlans(ctx)
	require.NoError(t, err)
	require.NotEmpty(t, plans)

	// Get a specific plan
	plan, err := service.GetPlan(ctx, plans[0].ID)
	require.NoError(t, err)
	assert.NotNil(t, plan)
	assert.Equal(t, plans[0].ID, plan.ID)
}

func TestUpdateSubscription(t *testing.T) {
	service := NewPaymentService()
	ctx := context.Background()

	// Get available plans first
	plans, err := service.ListPlans(ctx)
	require.NoError(t, err)
	require.NotEmpty(t, plans)

	userID := uuid.New()
	planID := plans[0].ID

	// Create a subscription first
	subscription, err := service.CreateSubscription(ctx, userID, planID)
	require.NoError(t, err)

	// Update the subscription
	newEndDate := time.Now().Add(30 * 24 * time.Hour)
	subscription.CurrentPeriodEnd = newEndDate

	updated, err := service.UpdateSubscription(ctx, subscription)
	require.NoError(t, err)
	assert.NotNil(t, updated)
	assert.Equal(t, newEndDate.Unix(), updated.CurrentPeriodEnd.Unix())
}

func TestHandleWebhookEvent(t *testing.T) {
	service := NewPaymentService()
	ctx := context.Background()

	tests := []struct {
		name      string
		eventType string
		payload   []byte
		wantErr   bool
	}{
		{
			name:      "successful webhook handling - subscription created",
			eventType: "customer.subscription.created",
			payload:   []byte(`{"id":"sub_test","customer":"cus_test"}`),
			wantErr:   false,
		},
		{
			name:      "successful webhook handling - subscription updated",
			eventType: "customer.subscription.updated",
			payload:   []byte(`{"id":"sub_test","customer":"cus_test"}`),
			wantErr:   false,
		},
		{
			name:      "successful webhook handling - payment succeeded",
			eventType: "payment_intent.succeeded",
			payload:   []byte(`{"id":"pi_test","amount":999}`),
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := service.HandleWebhookEvent(ctx, tt.eventType, tt.payload)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
