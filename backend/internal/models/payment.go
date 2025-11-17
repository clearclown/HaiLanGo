package models

import (
	"time"

	"github.com/google/uuid"
)

// SubscriptionPlan はサブスクリプションプラン
type SubscriptionPlan string

const (
	PlanFree    SubscriptionPlan = "free"     // 無料プラン
	PlanPremium SubscriptionPlan = "premium"  // プレミアムプラン（月額）
	PlanYearly  SubscriptionPlan = "yearly"   // 年間プラン
)

// SubscriptionStatus はサブスクリプションステータス
type SubscriptionStatus string

const (
	SubscriptionStatusActive    SubscriptionStatus = "active"     // アクティブ
	SubscriptionStatusCanceled  SubscriptionStatus = "canceled"   // キャンセル済み
	SubscriptionStatusPastDue   SubscriptionStatus = "past_due"   // 支払い遅延
	SubscriptionStatusTrialing  SubscriptionStatus = "trialing"   // トライアル中
	SubscriptionStatusIncomplete SubscriptionStatus = "incomplete" // 不完全
)

// PaymentStatus は支払いステータス
type PaymentStatus string

const (
	PaymentStatusPending   PaymentStatus = "pending"    // 保留中
	PaymentStatusSucceeded PaymentStatus = "succeeded"  // 成功
	PaymentStatusFailed    PaymentStatus = "failed"     // 失敗
	PaymentStatusRefunded  PaymentStatus = "refunded"   // 返金済み
)

// CreateSubscriptionRequest はサブスクリプション作成リクエスト
type CreateSubscriptionRequest struct {
	Plan          SubscriptionPlan `json:"plan" binding:"required"`            // プラン（premium/yearly）
	PaymentMethod string           `json:"payment_method" binding:"required"`  // Stripe Payment Method ID
}

// SubscriptionResponse はサブスクリプションレスポンス
type SubscriptionResponse struct {
	ID                   string             `json:"id"`
	UserID               string             `json:"user_id"`
	Plan                 SubscriptionPlan   `json:"plan"`
	Status               SubscriptionStatus `json:"status"`
	StripeSubscriptionID string             `json:"stripe_subscription_id,omitempty"`
	StripeCustomerID     string             `json:"stripe_customer_id,omitempty"`
	CurrentPeriodStart   time.Time          `json:"current_period_start"`
	CurrentPeriodEnd     time.Time          `json:"current_period_end"`
	CancelAtPeriodEnd    bool               `json:"cancel_at_period_end"`
	CreatedAt            time.Time          `json:"created_at"`
	UpdatedAt            time.Time          `json:"updated_at"`
}

// Subscription はサブスクリプション情報
type Subscription struct {
	ID                   uuid.UUID
	UserID               uuid.UUID
	Plan                 SubscriptionPlan
	Status               SubscriptionStatus
	StripeSubscriptionID string
	StripeCustomerID     string
	CurrentPeriodStart   time.Time
	CurrentPeriodEnd     time.Time
	CancelAtPeriodEnd    bool
	CreatedAt            time.Time
	UpdatedAt            time.Time
}

// PaymentHistoryItem は支払い履歴アイテム
type PaymentHistoryItem struct {
	ID              string        `json:"id"`
	Amount          int64         `json:"amount"`           // 金額（セント単位）
	Currency        string        `json:"currency"`         // 通貨コード（例: "usd", "jpy"）
	Status          PaymentStatus `json:"status"`
	Description     string        `json:"description"`
	InvoiceURL      string        `json:"invoice_url,omitempty"`
	ReceiptURL      string        `json:"receipt_url,omitempty"`
	StripePaymentID string        `json:"stripe_payment_id"`
	CreatedAt       time.Time     `json:"created_at"`
}

// Payment は支払い情報
type Payment struct {
	ID              uuid.UUID
	UserID          uuid.UUID
	SubscriptionID  uuid.UUID
	Amount          int64
	Currency        string
	Status          PaymentStatus
	Description     string
	InvoiceURL      string
	ReceiptURL      string
	StripePaymentID string
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

// CancelSubscriptionRequest はサブスクリプションキャンセルリクエスト
type CancelSubscriptionRequest struct {
	CancelAtPeriodEnd bool `json:"cancel_at_period_end"` // 期間終了時にキャンセルするか
}

// UpdatePaymentMethodRequest は支払い方法更新リクエスト
type UpdatePaymentMethodRequest struct {
	PaymentMethod string `json:"payment_method" binding:"required"` // Stripe Payment Method ID
}

// StripeWebhookEvent はStripe Webhookイベント
type StripeWebhookEvent struct {
	Type string      `json:"type"`
	Data interface{} `json:"data"`
}

// PlanPricing はプラン料金情報
type PlanPricing struct {
	Plan        SubscriptionPlan `json:"plan"`
	Name        string           `json:"name"`
	Description string           `json:"description"`
	Price       int64            `json:"price"`        // 金額（セント単位）
	Currency    string           `json:"currency"`     // 通貨コード
	Interval    string           `json:"interval"`     // 請求間隔（month/year）
	Features    []string         `json:"features"`     // 機能一覧
}

// SubscriptionUsage はサブスクリプション使用状況
type SubscriptionUsage struct {
	Plan               SubscriptionPlan `json:"plan"`
	DailyPagesLimit    int              `json:"daily_pages_limit"`     // 1日あたりのページ制限
	DailyMinutesLimit  int              `json:"daily_minutes_limit"`   // 1日あたりの分数制限
	PagesUsedToday     int              `json:"pages_used_today"`      // 今日使用したページ数
	MinutesUsedToday   int              `json:"minutes_used_today"`    // 今日使用した分数
	TTSQualityLevel    string           `json:"tts_quality_level"`     // TTS音質レベル
	OfflineDownload    bool             `json:"offline_download"`      // オフラインダウンロード可能か
	ResetAt            time.Time        `json:"reset_at"`              // 制限リセット時刻
}
