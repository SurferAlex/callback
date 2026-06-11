CREATE TABLE subscription_notifications (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users (telegram_id) ON DELETE CASCADE,
    notification_type TEXT NOT NULL,
    sent_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    CONSTRAINT uq_subscription_notifications_user_type UNIQUE (user_id, notification_type)
);

CREATE INDEX idx_subscription_notifications_user_id ON subscription_notifications (user_id);
