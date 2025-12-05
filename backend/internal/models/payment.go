package models

import (
	"time"

	"github.com/google/uuid"
)

// PlanType はサブスクリプションプランの種類
type PlanType string

const (
	PlanTypeFree    PlanType = "free"    // 無料プラン
	PlanTypePremium PlanType = "premium" // プレミアムプラン（月額）
	PlanTypeYearly  PlanType = "yearly"  // 年間プラン
)

// SubscriptionPlan はサブスクリプションプラン詳細
type SubscriptionPlan struct {
	ID            uuid.UUID `json:"id"`
	Name          string    `json:"name"`
	Description   string    `json:"description"`
	PlanType      PlanType  `json:"plan_type"`
	Price         int64     `json:"price"`           // 金額（セント単位）
	Currency      string    `json:"currency"`        // 通貨コード
	Interval      string    `json:"interval"`        // 請求間隔（month/year）
	StripePriceID string    `json:"stripe_price_id"` // Stripe Price ID
	Active        bool      `json:"active"`
	Features      []string  `json:"features"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// SubscriptionStatus はサブスクリプションステータス
type SubscriptionStatus string

const (
	SubscriptionStatusActive     SubscriptionStatus = "active"     // アクティブ
	SubscriptionStatusCanceled   SubscriptionStatus = "canceled"   // キャンセル済み
	SubscriptionStatusPastDue    SubscriptionStatus = "past_due"   // 支払い遅延
	SubscriptionStatusTrialing   SubscriptionStatus = "trialing"   // トライアル中
	SubscriptionStatusIncomplete SubscriptionStatus = "incomplete" // 不完全
)

// PaymentStatus は支払いステータス
type PaymentStatus string

const (
	PaymentStatusPending   PaymentStatus = "pending"   // 保留中
	PaymentStatusSucceeded PaymentStatus = "succeeded" // 成功
	PaymentStatusFailed    PaymentStatus = "failed"    // 失敗
	PaymentStatusRefunded  PaymentStatus = "refunded"  // 返金済み
)

// CreateSubscriptionRequest はサブスクリプション作成リクエスト
type CreateSubscriptionRequest struct {
	PlanID        uuid.UUID `json:"plan_id"`                           // プランID（新形式）
	Plan          PlanType  `json:"plan"`                              // プラン種類（後方互換性）
	PaymentMethod string    `json:"payment_method" binding:"required"` // Stripe Payment Method ID
}

// SubscriptionResponse はサブスクリプションレスポンス
type SubscriptionResponse struct {
	ID                   uuid.UUID          `json:"id"`
	UserID               uuid.UUID          `json:"user_id"`
	PlanID               uuid.UUID          `json:"plan_id,omitempty"`
	Plan                 PlanType           `json:"plan"`
	Status               SubscriptionStatus `json:"status"`
	StripeSubscriptionID string             `json:"stripe_subscription_id,omitempty"`
	StripeCustomerID     string             `json:"stripe_customer_id,omitempty"`
	CurrentPeriodStart   time.Time          `json:"current_period_start"`
	CurrentPeriodEnd     time.Time          `json:"current_period_end"`
	CancelAtPeriodEnd    bool               `json:"cancel_at_period_end"`
	CanceledAt           *time.Time         `json:"canceled_at,omitempty"`
	CreatedAt            time.Time          `json:"created_at"`
	UpdatedAt            time.Time          `json:"updated_at"`
}

// Subscription はサブスクリプション情報
type Subscription struct {
	ID                   uuid.UUID
	UserID               uuid.UUID
	PlanID               uuid.UUID          // 新形式: プランIDへの参照
	Plan                 PlanType           // 後方互換性: プラン種類を直接保持
	Status               SubscriptionStatus
	StripeSubscriptionID string
	StripeCustomerID     string
	CurrentPeriodStart   time.Time
	CurrentPeriodEnd     time.Time
	CancelAtPeriodEnd    bool
	CanceledAt           *time.Time
	CreatedAt            time.Time
	UpdatedAt            time.Time
}

// PaymentHistoryItem は支払い履歴アイテム
type PaymentHistoryItem struct {
	ID              uuid.UUID     `json:"id"`
	Amount          int64         `json:"amount"`   // 金額（セント単位）
	Currency        string        `json:"currency"` // 通貨コード（例: "usd", "jpy"）
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

// SubscriptionUsage はサブスクリプション使用状況
type SubscriptionUsage struct {
	Plan              PlanType  `json:"plan_type"`           // 後方互換性のためPlanとしても使用
	PlanType          PlanType  `json:"-"`                   // deprecated: Planを使用
	DailyPagesLimit   int       `json:"daily_pages_limit"`   // 1日あたりのページ制限
	DailyMinutesLimit int       `json:"daily_minutes_limit"` // 1日あたりの分数制限
	PagesUsedToday    int       `json:"pages_used_today"`    // 今日使用したページ数
	MinutesUsedToday  int       `json:"minutes_used_today"`  // 今日使用した分数
	TTSQualityLevel   string    `json:"tts_quality_level"`   // TTS音質レベル
	OfflineDownload   bool      `json:"offline_download"`    // オフラインダウンロード可能か
	ResetAt           time.Time `json:"reset_at"`            // 制限リセット時刻
}

// PlanPricing はプラン料金情報（InMemoryリポジトリとの互換性のため）
type PlanPricing struct {
	Plan        PlanType `json:"plan"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Price       int64    `json:"price"`    // セント単位
	Currency    string   `json:"currency"` // 通貨コード
	Interval    string   `json:"interval"` // month/year
	Features    []string `json:"features"`
}

// 後方互換性のためのプラン定数エイリアス
const (
	PlanFree    = PlanTypeFree
	PlanPremium = PlanTypePremium
	PlanYearly  = PlanTypeYearly
)
