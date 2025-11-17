package repository

import (
	"context"
	"database/sql"
	"fmt"
	"sync"
	"time"

	"github.com/clearclown/HaiLanGo/backend/internal/models"
	"github.com/google/uuid"
)

// PaymentRepositoryInterface は決済リポジトリのインターフェース
type PaymentRepositoryInterface interface {
	// CreateSubscription はサブスクリプションを作成
	CreateSubscription(ctx context.Context, userID uuid.UUID, plan models.SubscriptionPlan, stripeSubscriptionID, stripeCustomerID string) (*models.Subscription, error)

	// GetSubscription はサブスクリプションを取得
	GetSubscription(ctx context.Context, subscriptionID uuid.UUID) (*models.Subscription, error)

	// GetSubscriptionByUserID はユーザーIDでサブスクリプションを取得
	GetSubscriptionByUserID(ctx context.Context, userID uuid.UUID) (*models.Subscription, error)

	// UpdateSubscriptionStatus はサブスクリプションステータスを更新
	UpdateSubscriptionStatus(ctx context.Context, subscriptionID uuid.UUID, status models.SubscriptionStatus) error

	// CancelSubscription はサブスクリプションをキャンセル
	CancelSubscription(ctx context.Context, subscriptionID uuid.UUID, cancelAtPeriodEnd bool) error

	// CreatePayment は支払いを作成
	CreatePayment(ctx context.Context, userID, subscriptionID uuid.UUID, amount int64, currency string, stripePaymentID string) (*models.Payment, error)

	// GetPaymentHistory は支払い履歴を取得
	GetPaymentHistory(ctx context.Context, userID uuid.UUID, limit int) ([]*models.Payment, error)

	// GetPlanPricing はプラン料金情報を取得
	GetPlanPricing(ctx context.Context) ([]*models.PlanPricing, error)

	// GetSubscriptionUsage はサブスクリプション使用状況を取得
	GetSubscriptionUsage(ctx context.Context, userID uuid.UUID) (*models.SubscriptionUsage, error)

	// UpdateUsage は使用状況を更新
	UpdateUsage(ctx context.Context, userID uuid.UUID, pagesUsed, minutesUsed int) error
}

// InMemoryPaymentRepository はインメモリ決済リポジトリ
type InMemoryPaymentRepository struct {
	mu             sync.RWMutex
	subscriptions  map[string]*models.Subscription       // subscriptionID -> Subscription
	userSubscriptions map[string]string                  // userID -> subscriptionID
	payments       map[string]*models.Payment            // paymentID -> Payment
	userPayments   map[string][]string                   // userID -> []paymentID
	usages         map[string]*models.SubscriptionUsage  // userID -> Usage
	planPricing    []*models.PlanPricing                 // プラン料金情報
}

// NewInMemoryPaymentRepository はインメモリ決済リポジトリを作成
func NewInMemoryPaymentRepository() *InMemoryPaymentRepository {
	repo := &InMemoryPaymentRepository{
		subscriptions:     make(map[string]*models.Subscription),
		userSubscriptions: make(map[string]string),
		payments:          make(map[string]*models.Payment),
		userPayments:      make(map[string][]string),
		usages:            make(map[string]*models.SubscriptionUsage),
		planPricing:       make([]*models.PlanPricing, 0),
	}

	// プラン料金情報を初期化
	repo.initPlanPricing()

	// サンプルデータを初期化
	repo.initSampleData()

	return repo
}

func (r *InMemoryPaymentRepository) initPlanPricing() {
	r.planPricing = []*models.PlanPricing{
		{
			Plan:        models.PlanFree,
			Name:        "Free",
			Description: "基本的な学習機能",
			Price:       0,
			Currency:    "usd",
			Interval:    "month",
			Features: []string{
				"1日1ページまで学習",
				"1日30分まで使用",
				"標準品質のTTS",
				"基本的な学習統計",
			},
		},
		{
			Plan:        models.PlanPremium,
			Name:        "Premium",
			Description: "すべての機能が使い放題",
			Price:       999, // $9.99
			Currency:    "usd",
			Interval:    "month",
			Features: []string{
				"無制限の学習",
				"高品質TTS（ElevenLabs）",
				"オフライン音声ダウンロード",
				"詳細な学習分析",
				"優先サポート",
			},
		},
		{
			Plan:        models.PlanYearly,
			Name:        "Yearly Premium",
			Description: "年間プラン（2ヶ月分お得）",
			Price:       9999, // $99.99
			Currency:    "usd",
			Interval:    "year",
			Features: []string{
				"無制限の学習",
				"高品質TTS（ElevenLabs）",
				"オフライン音声ダウンロード",
				"詳細な学習分析",
				"優先サポート",
				"年間プラン特典",
			},
		},
	}
}

func (r *InMemoryPaymentRepository) initSampleData() {
	testUserID := uuid.MustParse("550e8400-e29b-41d4-a716-446655440000")

	// 無料プランのサブスクリプション
	freeSubID := uuid.New()
	freeSub := &models.Subscription{
		ID:                 freeSubID,
		UserID:             testUserID,
		Plan:               models.PlanFree,
		Status:             models.SubscriptionStatusActive,
		CurrentPeriodStart: time.Now().AddDate(0, -1, 0),
		CurrentPeriodEnd:   time.Now().AddDate(0, 1, 0),
		CancelAtPeriodEnd:  false,
		CreatedAt:          time.Now().AddDate(0, -1, 0),
		UpdatedAt:          time.Now(),
	}

	r.subscriptions[freeSubID.String()] = freeSub
	r.userSubscriptions[testUserID.String()] = freeSubID.String()

	// 無料プランの使用状況
	r.usages[testUserID.String()] = &models.SubscriptionUsage{
		Plan:              models.PlanFree,
		DailyPagesLimit:   1,
		DailyMinutesLimit: 30,
		PagesUsedToday:    0,
		MinutesUsedToday:  15,
		TTSQualityLevel:   "standard",
		OfflineDownload:   false,
		ResetAt:           time.Now().Add(24 * time.Hour),
	}

	// 別のテストユーザー（プレミアムプラン）
	premiumUserID := uuid.New()
	premiumSubID := uuid.New()
	premiumSub := &models.Subscription{
		ID:                   premiumSubID,
		UserID:               premiumUserID,
		Plan:                 models.PlanPremium,
		Status:               models.SubscriptionStatusActive,
		StripeSubscriptionID: "sub_" + uuid.New().String(),
		StripeCustomerID:     "cus_" + uuid.New().String(),
		CurrentPeriodStart:   time.Now().AddDate(0, 0, -15),
		CurrentPeriodEnd:     time.Now().AddDate(0, 0, 15),
		CancelAtPeriodEnd:    false,
		CreatedAt:            time.Now().AddDate(0, -2, 0),
		UpdatedAt:            time.Now(),
	}

	r.subscriptions[premiumSubID.String()] = premiumSub
	r.userSubscriptions[premiumUserID.String()] = premiumSubID.String()

	// プレミアムプランの使用状況
	r.usages[premiumUserID.String()] = &models.SubscriptionUsage{
		Plan:              models.PlanPremium,
		DailyPagesLimit:   -1, // 無制限
		DailyMinutesLimit: -1, // 無制限
		PagesUsedToday:    25,
		MinutesUsedToday:  180,
		TTSQualityLevel:   "premium",
		OfflineDownload:   true,
		ResetAt:           time.Now().Add(24 * time.Hour),
	}

	// 支払い履歴（プレミアムユーザー）
	for i := 1; i <= 3; i++ {
		paymentID := uuid.New()
		payment := &models.Payment{
			ID:              paymentID,
			UserID:          premiumUserID,
			SubscriptionID:  premiumSubID,
			Amount:          999,
			Currency:        "usd",
			Status:          models.PaymentStatusSucceeded,
			Description:     fmt.Sprintf("Premium subscription - Month %d", i),
			InvoiceURL:      fmt.Sprintf("https://invoice.stripe.com/i/%s", uuid.New().String()),
			ReceiptURL:      fmt.Sprintf("https://pay.stripe.com/receipts/%s", uuid.New().String()),
			StripePaymentID: "pi_" + uuid.New().String(),
			CreatedAt:       time.Now().AddDate(0, -i, 0),
			UpdatedAt:       time.Now().AddDate(0, -i, 0),
		}

		r.payments[paymentID.String()] = payment
		r.userPayments[premiumUserID.String()] = append(r.userPayments[premiumUserID.String()], paymentID.String())
	}
}

func (r *InMemoryPaymentRepository) CreateSubscription(ctx context.Context, userID uuid.UUID, plan models.SubscriptionPlan, stripeSubscriptionID, stripeCustomerID string) (*models.Subscription, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	subscriptionID := uuid.New()
	now := time.Now()

	var periodEnd time.Time
	if plan == models.PlanYearly {
		periodEnd = now.AddDate(1, 0, 0) // 1年後
	} else {
		periodEnd = now.AddDate(0, 1, 0) // 1ヶ月後
	}

	subscription := &models.Subscription{
		ID:                   subscriptionID,
		UserID:               userID,
		Plan:                 plan,
		Status:               models.SubscriptionStatusActive,
		StripeSubscriptionID: stripeSubscriptionID,
		StripeCustomerID:     stripeCustomerID,
		CurrentPeriodStart:   now,
		CurrentPeriodEnd:     periodEnd,
		CancelAtPeriodEnd:    false,
		CreatedAt:            now,
		UpdatedAt:            now,
	}

	r.subscriptions[subscriptionID.String()] = subscription
	r.userSubscriptions[userID.String()] = subscriptionID.String()

	// 使用状況を初期化
	r.initUsageForPlan(userID, plan)

	return subscription, nil
}

func (r *InMemoryPaymentRepository) initUsageForPlan(userID uuid.UUID, plan models.SubscriptionPlan) {
	usage := &models.SubscriptionUsage{
		Plan:             plan,
		PagesUsedToday:   0,
		MinutesUsedToday: 0,
		ResetAt:          time.Now().Add(24 * time.Hour),
	}

	switch plan {
	case models.PlanFree:
		usage.DailyPagesLimit = 1
		usage.DailyMinutesLimit = 30
		usage.TTSQualityLevel = "standard"
		usage.OfflineDownload = false
	case models.PlanPremium, models.PlanYearly:
		usage.DailyPagesLimit = -1  // 無制限
		usage.DailyMinutesLimit = -1 // 無制限
		usage.TTSQualityLevel = "premium"
		usage.OfflineDownload = true
	}

	r.usages[userID.String()] = usage
}

func (r *InMemoryPaymentRepository) GetSubscription(ctx context.Context, subscriptionID uuid.UUID) (*models.Subscription, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	sub, exists := r.subscriptions[subscriptionID.String()]
	if !exists {
		return nil, fmt.Errorf("subscription not found: %s", subscriptionID)
	}

	return sub, nil
}

func (r *InMemoryPaymentRepository) GetSubscriptionByUserID(ctx context.Context, userID uuid.UUID) (*models.Subscription, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	subID, exists := r.userSubscriptions[userID.String()]
	if !exists {
		return nil, fmt.Errorf("subscription not found for user: %s", userID)
	}

	sub, exists := r.subscriptions[subID]
	if !exists {
		return nil, fmt.Errorf("subscription not found: %s", subID)
	}

	return sub, nil
}

func (r *InMemoryPaymentRepository) UpdateSubscriptionStatus(ctx context.Context, subscriptionID uuid.UUID, status models.SubscriptionStatus) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	sub, exists := r.subscriptions[subscriptionID.String()]
	if !exists {
		return fmt.Errorf("subscription not found: %s", subscriptionID)
	}

	sub.Status = status
	sub.UpdatedAt = time.Now()

	return nil
}

func (r *InMemoryPaymentRepository) CancelSubscription(ctx context.Context, subscriptionID uuid.UUID, cancelAtPeriodEnd bool) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	sub, exists := r.subscriptions[subscriptionID.String()]
	if !exists {
		return fmt.Errorf("subscription not found: %s", subscriptionID)
	}

	sub.CancelAtPeriodEnd = cancelAtPeriodEnd
	if !cancelAtPeriodEnd {
		sub.Status = models.SubscriptionStatusCanceled
	}
	sub.UpdatedAt = time.Now()

	return nil
}

func (r *InMemoryPaymentRepository) CreatePayment(ctx context.Context, userID, subscriptionID uuid.UUID, amount int64, currency string, stripePaymentID string) (*models.Payment, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	paymentID := uuid.New()
	now := time.Now()

	payment := &models.Payment{
		ID:              paymentID,
		UserID:          userID,
		SubscriptionID:  subscriptionID,
		Amount:          amount,
		Currency:        currency,
		Status:          models.PaymentStatusSucceeded,
		Description:     "Subscription payment",
		InvoiceURL:      fmt.Sprintf("https://invoice.stripe.com/i/%s", uuid.New().String()),
		ReceiptURL:      fmt.Sprintf("https://pay.stripe.com/receipts/%s", uuid.New().String()),
		StripePaymentID: stripePaymentID,
		CreatedAt:       now,
		UpdatedAt:       now,
	}

	r.payments[paymentID.String()] = payment
	r.userPayments[userID.String()] = append(r.userPayments[userID.String()], paymentID.String())

	return payment, nil
}

func (r *InMemoryPaymentRepository) GetPaymentHistory(ctx context.Context, userID uuid.UUID, limit int) ([]*models.Payment, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	paymentIDs, exists := r.userPayments[userID.String()]
	if !exists {
		return []*models.Payment{}, nil
	}

	payments := make([]*models.Payment, 0)
	for i := len(paymentIDs) - 1; i >= 0 && len(payments) < limit; i-- {
		if payment, exists := r.payments[paymentIDs[i]]; exists {
			payments = append(payments, payment)
		}
	}

	return payments, nil
}

func (r *InMemoryPaymentRepository) GetPlanPricing(ctx context.Context) ([]*models.PlanPricing, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	return r.planPricing, nil
}

func (r *InMemoryPaymentRepository) GetSubscriptionUsage(ctx context.Context, userID uuid.UUID) (*models.SubscriptionUsage, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	usage, exists := r.usages[userID.String()]
	if !exists {
		return nil, fmt.Errorf("usage not found for user: %s", userID)
	}

	return usage, nil
}

func (r *InMemoryPaymentRepository) UpdateUsage(ctx context.Context, userID uuid.UUID, pagesUsed, minutesUsed int) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	usage, exists := r.usages[userID.String()]
	if !exists {
		return fmt.Errorf("usage not found for user: %s", userID)
	}

	usage.PagesUsedToday += pagesUsed
	usage.MinutesUsedToday += minutesUsed

	return nil
}

// PostgreSQL Implementation

type PaymentRepositoryPostgres struct {
	db *sql.DB
}

func NewPaymentRepositoryPostgres(db *sql.DB) PaymentRepositoryInterface {
	return &PaymentRepositoryPostgres{db: db}
}

func (r *PaymentRepositoryPostgres) CreateSubscription(ctx context.Context, userID uuid.UUID, plan models.SubscriptionPlan, stripeSubscriptionID, stripeCustomerID string) (*models.Subscription, error) {
	subscriptionID := uuid.New()
	now := time.Now()

	var periodEnd time.Time
	if plan == models.PlanYearly {
		periodEnd = now.AddDate(1, 0, 0)
	} else {
		periodEnd = now.AddDate(0, 1, 0)
	}

	// Get plan_id from subscription_plans table
	var planID uuid.UUID
	err := r.db.QueryRowContext(ctx, `
		SELECT id FROM subscription_plans WHERE stripe_price_id = $1
	`, "price_"+string(plan)).Scan(&planID)
	if err != nil {
		return nil, fmt.Errorf("plan not found: %w", err)
	}

	_, err = r.db.ExecContext(ctx, `
		INSERT INTO subscriptions (id, user_id, plan_id, stripe_customer_id, stripe_subscription_id, status, current_period_start, current_period_end, cancel_at_period_end, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	`, subscriptionID, userID, planID, stripeCustomerID, stripeSubscriptionID, "active", now, periodEnd, false, now, now)

	if err != nil {
		return nil, err
	}

	return &models.Subscription{
		ID:                   subscriptionID,
		UserID:               userID,
		Plan:                 plan,
		Status:               models.SubscriptionStatusActive,
		StripeSubscriptionID: stripeSubscriptionID,
		StripeCustomerID:     stripeCustomerID,
		CurrentPeriodStart:   now,
		CurrentPeriodEnd:     periodEnd,
		CancelAtPeriodEnd:    false,
		CreatedAt:            now,
		UpdatedAt:            now,
	}, nil
}

func (r *PaymentRepositoryPostgres) GetSubscription(ctx context.Context, subscriptionID uuid.UUID) (*models.Subscription, error) {
	var sub models.Subscription
	var status string
	var canceledAt sql.NullTime

	err := r.db.QueryRowContext(ctx, `
		SELECT id, user_id, stripe_customer_id, stripe_subscription_id, status, current_period_start, current_period_end, cancel_at_period_end, canceled_at, created_at, updated_at
		FROM subscriptions
		WHERE id = $1
	`, subscriptionID).Scan(&sub.ID, &sub.UserID, &sub.StripeCustomerID, &sub.StripeSubscriptionID, &status, &sub.CurrentPeriodStart, &sub.CurrentPeriodEnd, &sub.CancelAtPeriodEnd, &canceledAt, &sub.CreatedAt, &sub.UpdatedAt)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("subscription not found: %s", subscriptionID)
	}
	if err != nil {
		return nil, err
	}

	sub.Status = models.SubscriptionStatus(status)
	// Note: canceled_at exists in database but not in model, so we ignore it

	return &sub, nil
}

func (r *PaymentRepositoryPostgres) GetSubscriptionByUserID(ctx context.Context, userID uuid.UUID) (*models.Subscription, error) {
	var sub models.Subscription
	var status string
	var canceledAt sql.NullTime

	err := r.db.QueryRowContext(ctx, `
		SELECT id, user_id, stripe_customer_id, stripe_subscription_id, status, current_period_start, current_period_end, cancel_at_period_end, canceled_at, created_at, updated_at
		FROM subscriptions
		WHERE user_id = $1
		ORDER BY created_at DESC
		LIMIT 1
	`, userID).Scan(&sub.ID, &sub.UserID, &sub.StripeCustomerID, &sub.StripeSubscriptionID, &status, &sub.CurrentPeriodStart, &sub.CurrentPeriodEnd, &sub.CancelAtPeriodEnd, &canceledAt, &sub.CreatedAt, &sub.UpdatedAt)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("subscription not found for user: %s", userID)
	}
	if err != nil {
		return nil, err
	}

	sub.Status = models.SubscriptionStatus(status)
	// Note: canceled_at exists in database but not in model, so we ignore it

	return &sub, nil
}

func (r *PaymentRepositoryPostgres) UpdateSubscriptionStatus(ctx context.Context, subscriptionID uuid.UUID, status models.SubscriptionStatus) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE subscriptions SET status = $1, updated_at = $2
		WHERE id = $3
	`, string(status), time.Now(), subscriptionID)
	return err
}

func (r *PaymentRepositoryPostgres) CancelSubscription(ctx context.Context, subscriptionID uuid.UUID, cancelAtPeriodEnd bool) error {
	query := `UPDATE subscriptions SET cancel_at_period_end = $1, updated_at = $2 WHERE id = $3`
	args := []interface{}{cancelAtPeriodEnd, time.Now(), subscriptionID}

	if !cancelAtPeriodEnd {
		query = `UPDATE subscriptions SET cancel_at_period_end = $1, status = $2, canceled_at = $3, updated_at = $3 WHERE id = $4`
		args = []interface{}{cancelAtPeriodEnd, string(models.SubscriptionStatusCanceled), time.Now(), subscriptionID}
	}

	_, err := r.db.ExecContext(ctx, query, args...)
	return err
}

func (r *PaymentRepositoryPostgres) CreatePayment(ctx context.Context, userID, subscriptionID uuid.UUID, amount int64, currency string, stripePaymentID string) (*models.Payment, error) {
	paymentID := uuid.New()
	now := time.Now()

	_, err := r.db.ExecContext(ctx, `
		INSERT INTO payments (id, user_id, subscription_id, stripe_payment_intent, amount, currency, status, description, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`, paymentID, userID, subscriptionID, stripePaymentID, amount, currency, "succeeded", "Subscription payment", now)

	if err != nil {
		return nil, err
	}

	return &models.Payment{
		ID:              paymentID,
		UserID:          userID,
		SubscriptionID:  subscriptionID,
		Amount:          amount,
		Currency:        currency,
		Status:          models.PaymentStatusSucceeded,
		Description:     "Subscription payment",
		StripePaymentID: stripePaymentID,
		CreatedAt:       now,
		UpdatedAt:       now,
	}, nil
}

func (r *PaymentRepositoryPostgres) GetPaymentHistory(ctx context.Context, userID uuid.UUID, limit int) ([]*models.Payment, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, user_id, subscription_id, stripe_payment_intent, amount, currency, status, description, created_at
		FROM payments
		WHERE user_id = $1
		ORDER BY created_at DESC
		LIMIT $2
	`, userID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	payments := []*models.Payment{}
	for rows.Next() {
		var payment models.Payment
		var status string

		err := rows.Scan(&payment.ID, &payment.UserID, &payment.SubscriptionID, &payment.StripePaymentID, &payment.Amount, &payment.Currency, &status, &payment.Description, &payment.CreatedAt)
		if err != nil {
			return nil, err
		}

		payment.Status = models.PaymentStatus(status)
		payments = append(payments, &payment)
	}

	return payments, nil
}

func (r *PaymentRepositoryPostgres) GetPlanPricing(ctx context.Context) ([]*models.PlanPricing, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT name, description, price, currency, interval
		FROM subscription_plans
		WHERE active = true
		ORDER BY price ASC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	planPricing := []*models.PlanPricing{}
	for rows.Next() {
		var pricing models.PlanPricing
		err := rows.Scan(&pricing.Name, &pricing.Description, &pricing.Price, &pricing.Currency, &pricing.Interval)
		if err != nil {
			return nil, err
		}

		// Map name to plan enum
		switch pricing.Name {
		case "Premium Monthly":
			pricing.Plan = models.PlanPremium
		case "Premium Yearly":
			pricing.Plan = models.PlanYearly
		default:
			pricing.Plan = models.PlanFree
		}

		planPricing = append(planPricing, &pricing)
	}

	return planPricing, nil
}

func (r *PaymentRepositoryPostgres) GetSubscriptionUsage(ctx context.Context, userID uuid.UUID) (*models.SubscriptionUsage, error) {
	// Get subscription plan
	sub, err := r.GetSubscriptionByUserID(ctx, userID)
	if err != nil {
		// Return default free plan usage if no subscription
		return &models.SubscriptionUsage{
			Plan:              models.PlanFree,
			DailyPagesLimit:   1,
			DailyMinutesLimit: 30,
			PagesUsedToday:    0,
			MinutesUsedToday:  0,
			TTSQualityLevel:   "standard",
			OfflineDownload:   false,
			ResetAt:           time.Now().Add(24 * time.Hour),
		}, nil
	}

	usage := &models.SubscriptionUsage{
		Plan:             sub.Plan,
		PagesUsedToday:   0,
		MinutesUsedToday: 0,
		ResetAt:          time.Now().Add(24 * time.Hour),
	}

	switch sub.Plan {
	case models.PlanFree:
		usage.DailyPagesLimit = 1
		usage.DailyMinutesLimit = 30
		usage.TTSQualityLevel = "standard"
		usage.OfflineDownload = false
	case models.PlanPremium, models.PlanYearly:
		usage.DailyPagesLimit = -1
		usage.DailyMinutesLimit = -1
		usage.TTSQualityLevel = "premium"
		usage.OfflineDownload = true
	}

	return usage, nil
}

func (r *PaymentRepositoryPostgres) UpdateUsage(ctx context.Context, userID uuid.UUID, pagesUsed, minutesUsed int) error {
	// For now, just return nil as we don't store usage in DB
	// In production, you would track daily usage in a separate table
	return nil
}
