CREATE TABLE auth_refresh_tokens (
    id BIGSERIAL PRIMARY KEY,
    telegram_user_id BIGINT NOT NULL REFERENCES users (telegram_id) ON DELETE CASCADE,
    token_hash TEXT NOT NULL,
    expires_at TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    CONSTRAINT uq_auth_refresh_token_hash UNIQUE (token_hash)
);

CREATE INDEX idx_auth_refresh_tokens_user ON auth_refresh_tokens (telegram_user_id);
CREATE INDEX idx_auth_refresh_tokens_expires ON auth_refresh_tokens (expires_at);
