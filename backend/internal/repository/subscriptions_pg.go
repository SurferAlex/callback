package repository

import (
	"context"
	"errors"
	"time"

	"api-vpn/internal/model"
	"api-vpn/internal/usecase"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type SubscriptionsRepo struct {
	db *pgxpool.Pool
}

func NewSubscriptionsRepo(db *pgxpool.Pool) *SubscriptionsRepo {
	return &SubscriptionsRepo{db: db}
}

func (r *SubscriptionsRepo) Create(ctx context.Context, p model.CreateSubscriptionParams) (model.Subscription, error) {
	const q = `
INSERT INTO subscriptions (
  telegram_user_id, plan_code, plan_label, status, starts_at, ends_at, client_uuid, is_mock
) VALUES ($1, $2, $3, 'active', $4, $5, $6, $7)
RETURNING id, telegram_user_id, plan_code, plan_label, status, starts_at, ends_at, client_uuid, is_mock, created_at, updated_at;
`
	var s model.Subscription
	var clientUUID *string
	if p.ClientUUID != nil {
		clientUUID = p.ClientUUID
	}
	err := r.db.QueryRow(ctx, q,
		p.TelegramUserID, p.PlanCode, p.PlanLabel, p.StartsAt.UTC(), p.EndsAt.UTC(), clientUUID, p.IsMock,
	).Scan(
		&s.ID, &s.TelegramUserID, &s.PlanCode, &s.PlanLabel, &s.Status,
		&s.StartsAt, &s.EndsAt, &s.ClientUUID, &s.IsMock, &s.CreatedAt, &s.UpdatedAt,
	)
	return s, err
}

func (r *SubscriptionsRepo) DeactivateActiveForUser(ctx context.Context, telegramUserID int64) error {
	const q = `
UPDATE subscriptions SET status = 'replaced', updated_at = now()
WHERE telegram_user_id = $1 AND status = 'active';
`
	_, err := r.db.Exec(ctx, q, telegramUserID)
	return err
}

func (r *SubscriptionsRepo) UpdateActiveClientUUID(ctx context.Context, telegramUserID int64, clientUUID string, now time.Time) error {
	const q = `
UPDATE subscriptions
SET client_uuid = $3, updated_at = now()
WHERE telegram_user_id = $1 AND status = 'active' AND ends_at > $2;
`
	_, err := r.db.Exec(ctx, q, telegramUserID, now.UTC(), clientUUID)
	return err
}

func (r *SubscriptionsRepo) GetActiveForUser(ctx context.Context, telegramUserID int64, now time.Time) (model.Subscription, error) {
	const q = `
SELECT id, telegram_user_id, plan_code, plan_label, status, starts_at, ends_at, client_uuid, is_mock, created_at, updated_at
FROM subscriptions
WHERE telegram_user_id = $1 AND status = 'active' AND ends_at > $2
ORDER BY ends_at DESC
LIMIT 1;
`
	var s model.Subscription
	err := r.db.QueryRow(ctx, q, telegramUserID, now.UTC()).Scan(
		&s.ID, &s.TelegramUserID, &s.PlanCode, &s.PlanLabel, &s.Status,
		&s.StartsAt, &s.EndsAt, &s.ClientUUID, &s.IsMock, &s.CreatedAt, &s.UpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return model.Subscription{}, usecase.ErrNotFound
	}
	return s, err
}
