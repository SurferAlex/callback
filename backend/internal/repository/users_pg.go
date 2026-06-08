package repository

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"

	"api-vpn/internal/model"
	"api-vpn/internal/usecase"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UsersRepo struct {
	db *pgxpool.Pool
}

func NewUsersRepo(db *pgxpool.Pool) *UsersRepo {
	return &UsersRepo{db: db}
}

func (r *UsersRepo) Upsert(ctx context.Context, p model.UpsertUserParams) (model.User, error) {
	const q = `
INSERT INTO users (telegram_id, first_name, last_name, username)
VALUES ($1, $2, $3, $4)
ON CONFLICT (telegram_id) DO UPDATE SET
  first_name = COALESCE(NULLIF(EXCLUDED.first_name, ''), users.first_name),
  last_name = COALESCE(EXCLUDED.last_name, users.last_name),
  username = COALESCE(EXCLUDED.username, users.username),
  updated_at = now()
RETURNING telegram_id, first_name, last_name, username, subscription_token, created_at, updated_at;
`
	var u model.User
	err := r.db.QueryRow(ctx, q, p.TelegramID, p.FirstName, p.LastName, p.Username).Scan(
		&u.TelegramID, &u.FirstName, &u.LastName, &u.Username, &u.SubscriptionToken, &u.CreatedAt, &u.UpdatedAt,
	)
	return u, err
}

func (r *UsersRepo) GetByTelegramID(ctx context.Context, telegramID int64) (model.User, error) {
	const q = `
SELECT telegram_id, first_name, last_name, username, subscription_token, created_at, updated_at
FROM users WHERE telegram_id = $1;
`
	var u model.User
	err := r.db.QueryRow(ctx, q, telegramID).Scan(
		&u.TelegramID, &u.FirstName, &u.LastName, &u.Username, &u.SubscriptionToken, &u.CreatedAt, &u.UpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return model.User{}, usecase.ErrNotFound
	}
	return u, err
}

func (r *UsersRepo) GetBySubscriptionToken(ctx context.Context, token string) (model.User, error) {
	const q = `
SELECT telegram_id, first_name, last_name, username, subscription_token, created_at, updated_at
FROM users WHERE subscription_token = $1;
`
	var u model.User
	err := r.db.QueryRow(ctx, q, token).Scan(
		&u.TelegramID, &u.FirstName, &u.LastName, &u.Username, &u.SubscriptionToken, &u.CreatedAt, &u.UpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return model.User{}, usecase.ErrNotFound
	}
	return u, err
}

func (r *UsersRepo) EnsureSubscriptionToken(ctx context.Context, telegramID int64) (string, error) {
	u, err := r.GetByTelegramID(ctx, telegramID)
	if err != nil {
		return "", err
	}
	if u.SubscriptionToken != nil && *u.SubscriptionToken != "" {
		return *u.SubscriptionToken, nil
	}
	token, err := newSubscriptionToken()
	if err != nil {
		return "", err
	}
	const q = `
UPDATE users SET subscription_token = $2, updated_at = now()
WHERE telegram_id = $1 AND (subscription_token IS NULL OR subscription_token = '');
`
	tag, err := r.db.Exec(ctx, q, telegramID, token)
	if err != nil {
		return "", err
	}
	if tag.RowsAffected() == 0 {
		u, err = r.GetByTelegramID(ctx, telegramID)
		if err != nil {
			return "", err
		}
		if u.SubscriptionToken != nil {
			return *u.SubscriptionToken, nil
		}
		return "", errors.New("subscription token not set")
	}
	return token, nil
}

func newSubscriptionToken() (string, error) {
	b := make([]byte, 24)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}
