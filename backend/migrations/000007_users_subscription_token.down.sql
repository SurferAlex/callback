DROP INDEX IF EXISTS uq_users_subscription_token;
ALTER TABLE users DROP COLUMN IF EXISTS subscription_token;
