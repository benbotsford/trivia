-- Subscriptions queries are defined but unused until billing ships. Keeping them
-- here so the EntitlementChecker implementation has a ready-to-go data access
-- layer when Stripe is wired up.

-- name: UpsertSubscription :one
INSERT INTO subscriptions (
    id,
    user_id,
    plan,
    status,
    current_period_end,
    stripe_customer_id,
    stripe_subscription_id
) VALUES ($1, $2, $3, $4, $5, $6, $7)
ON CONFLICT (user_id) DO UPDATE SET
    plan                   = EXCLUDED.plan,
    status                 = EXCLUDED.status,
    current_period_end     = EXCLUDED.current_period_end,
    stripe_customer_id     = EXCLUDED.stripe_customer_id,
    stripe_subscription_id = EXCLUDED.stripe_subscription_id
RETURNING *;

-- name: GetSubscriptionByUser :one
SELECT * FROM subscriptions WHERE user_id = $1;

-- name: GetSubscriptionByStripeSubID :one
SELECT * FROM subscriptions WHERE stripe_subscription_id = $1;
