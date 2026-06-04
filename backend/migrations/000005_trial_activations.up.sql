CREATE TABLE trial_activations (
    telegram_id BIGINT PRIMARY KEY,
    activated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
