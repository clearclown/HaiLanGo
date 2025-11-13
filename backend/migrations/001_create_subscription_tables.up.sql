-- Create subscription_plans table
CREATE TABLE IF NOT EXISTS subscription_plans (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    description TEXT,
    price BIGINT NOT NULL,  -- Price in cents
    currency VARCHAR(3) NOT NULL DEFAULT 'usd',
    interval VARCHAR(20) NOT NULL,  -- 'month' or 'year'
    stripe_price_id VARCHAR(255) NOT NULL UNIQUE,
    active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Create index on active plans
CREATE INDEX idx_subscription_plans_active ON subscription_plans(active);

-- Create subscriptions table
CREATE TABLE IF NOT EXISTS subscriptions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL,
    plan_id UUID NOT NULL REFERENCES subscription_plans(id),
    stripe_customer_id VARCHAR(255) NOT NULL,
    stripe_subscription_id VARCHAR(255) NOT NULL UNIQUE,
    status VARCHAR(20) NOT NULL,  -- 'active', 'canceled', 'past_due', 'trialing', 'unpaid'
    current_period_start TIMESTAMP NOT NULL,
    current_period_end TIMESTAMP NOT NULL,
    cancel_at_period_end BOOLEAN NOT NULL DEFAULT false,
    canceled_at TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Create indexes for subscriptions
CREATE INDEX idx_subscriptions_user_id ON subscriptions(user_id);
CREATE INDEX idx_subscriptions_status ON subscriptions(status);
CREATE INDEX idx_subscriptions_stripe_subscription_id ON subscriptions(stripe_subscription_id);

-- Create payments table for payment history
CREATE TABLE IF NOT EXISTS payments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL,
    subscription_id UUID NOT NULL REFERENCES subscriptions(id),
    stripe_payment_intent VARCHAR(255) NOT NULL,
    amount BIGINT NOT NULL,  -- Amount in cents
    currency VARCHAR(3) NOT NULL DEFAULT 'usd',
    status VARCHAR(20) NOT NULL,  -- 'succeeded', 'failed', 'pending'
    description TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Create indexes for payments
CREATE INDEX idx_payments_user_id ON payments(user_id);
CREATE INDEX idx_payments_subscription_id ON payments(subscription_id);
CREATE INDEX idx_payments_status ON payments(status);

-- Insert default subscription plans
INSERT INTO subscription_plans (id, name, description, price, currency, interval, stripe_price_id, active)
VALUES
    (gen_random_uuid(), 'Premium Monthly', 'Monthly premium subscription', 999, 'usd', 'month', 'price_premium_monthly', true),
    (gen_random_uuid(), 'Premium Yearly', 'Yearly premium subscription with 20% savings', 9999, 'usd', 'year', 'price_premium_yearly', true)
ON CONFLICT (stripe_price_id) DO NOTHING;
