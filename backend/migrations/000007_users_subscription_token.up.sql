ALTER TABLE users
    ADD COLUMN IF NOT EXISTS subscription_token TEXT;

CREATE UNIQUE INDEX IF NOT EXISTS uq_users_subscription_token
    ON users (subscription_token)
    WHERE subscription_token IS NOT NULL;
