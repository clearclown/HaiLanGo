-- Drop tables in reverse order (due to foreign key constraints)
DROP TABLE IF EXISTS payments;
DROP TABLE IF EXISTS subscriptions;
DROP TABLE IF EXISTS subscription_plans;
