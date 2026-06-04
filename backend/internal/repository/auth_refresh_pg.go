package repository

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type AuthRefreshRepo struct {
	db *pgxpool.Pool
}

func NewAuthRefreshRepo(db *pgxpool.Pool) *AuthRefreshRepo {
	return &AuthRefreshRepo{db: db}
}

func (r *AuthRefreshRepo) Store(ctx context.Context, telegramUserID int64, tokenHash string, expiresAt time.Time) error {
	const q = `
INSERT INTO auth_refresh_tokens (telegram_user_id, token_hash, expires_at)
VALUES ($1, $2, $3);
`
	_, err := r.db.Exec(ctx, q, telegramUserID, tokenHash, expiresAt.UTC())
	return err
}

func (r *AuthRefreshRepo) FindTelegramID(ctx context.Context, tokenHash string, now time.Time) (int64, error) {
	const q = `
SELECT telegram_user_id FROM auth_refresh_tokens
WHERE token_hash = $1 AND expires_at > $2
LIMIT 1;
`
	var id int64
	err := r.db.QueryRow(ctx, q, tokenHash, now.UTC()).Scan(&id)
	return id, err
}

func (r *AuthRefreshRepo) Revoke(ctx context.Context, tokenHash string) error {
	const q = `DELETE FROM auth_refresh_tokens WHERE token_hash = $1;`
	_, err := r.db.Exec(ctx, q, tokenHash)
	return err
}

func (r *AuthRefreshRepo) RevokeAllForUser(ctx context.Context, telegramUserID int64) error {
	const q = `DELETE FROM auth_refresh_tokens WHERE telegram_user_id = $1;`
	_, err := r.db.Exec(ctx, q, telegramUserID)
	return err
}
