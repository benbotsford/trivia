-- +goose Up
-- Empty-for-now subscriptions table. Billing implementation is deferred (see
-- PROJECT.md "Subscription-Readiness"). Populating happens later via Stripe
-- webhooks; until then the EntitlementChecker implementation should always
-- grant access.
CREATE TABLE subscriptions (
    id                     uuid PRIMARY KEY,
    user_id                uuid        NOT NULL UNIQUE REFERENCES users(id) ON DELETE CASCADE,
    plan                   text        NOT NULL,
    status                 text        NOT NULL,
    current_period_end     timestamptz,
    stripe_customer_id     text        UNIQUE,
    stripe_subscription_id text        UNIQUE,
    created_at             timestamptz NOT NULL DEFAULT now(),
    updated_at             timestamptz NOT NULL DEFAULT now(),
    CONSTRAINT subscriptions_plan_values   CHECK (plan   IN ('monthly', 'yearly')),
    CONSTRAINT subscriptions_status_values CHECK (status IN ('active', 'trialing', 'past_due', 'canceled', 'incomplete', 'incomplete_expired', 'unpaid'))
);

CREATE INDEX subscriptions_status_idx ON subscriptions (status);

CREATE TRIGGER subscriptions_set_updated_at
    BEFORE UPDATE ON subscriptions
    FOR EACH ROW EXECUTE FUNCTION trigger_set_updated_at();

-- +goose Down
DROP TABLE IF EXISTS subscriptions;
