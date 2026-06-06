package repository

import (
	"context"
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
RETURNING telegram_id, first_name, last_name, username, created_at, updated_at;
`
	var u model.User
	err := r.db.QueryRow(ctx, q, p.TelegramID, p.FirstName, p.LastName, p.Username).Scan(
		&u.TelegramID, &u.FirstName, &u.LastName, &u.Username, &u.CreatedAt, &u.UpdatedAt,
	)
	return u, err
}

func (r *UsersRepo) GetByTelegramID(ctx context.Context, telegramID int64) (model.User, error) {
	const q = `
SELECT telegram_id, first_name, last_name, username, created_at, updated_at
FROM users WHERE telegram_id = $1;
`
	var u model.User
	err := r.db.QueryRow(ctx, q, telegramID).Scan(
		&u.TelegramID, &u.FirstName, &u.LastName, &u.Username, &u.CreatedAt, &u.UpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return model.User{}, usecase.ErrNotFound
	}
	return u, err
}
