CREATE TABLE users (
    telegram_id BIGINT PRIMARY KEY,
    first_name TEXT NOT NULL DEFAULT '',
    last_name TEXT,
    username TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE subscriptions (
    id BIGSERIAL PRIMARY KEY,
    telegram_user_id BIGINT NOT NULL REFERENCES users (telegram_id) ON DELETE CASCADE,
    plan_code TEXT NOT NULL,
    plan_label TEXT NOT NULL,
    status TEXT NOT NULL DEFAULT 'active',
    starts_at TIMESTAMPTZ NOT NULL,
    ends_at TIMESTAMPTZ NOT NULL,
    client_uuid UUID REFERENCES vpn_clients (client_uuid) ON DELETE SET NULL,
    is_mock BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_subscriptions_telegram_user_id ON subscriptions (telegram_user_id);
CREATE INDEX idx_subscriptions_active ON subscriptions (telegram_user_id, ends_at DESC)
    WHERE status = 'active';

CREATE UNIQUE INDEX uq_vpn_clients_one_active_per_telegram
    ON vpn_clients (telegram_user_id)
    WHERE is_active = true AND telegram_user_id IS NOT NULL;
