package notify

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ActiveSubscription struct {
	UserID int64
	EndsAt time.Time
}

type Store struct {
	db *pgxpool.Pool
}

func NewStore(db *pgxpool.Pool) *Store {
	return &Store{db: db}
}

func OpenPool(ctx context.Context, databaseURL string) (*pgxpool.Pool, error) {
	cfg, err := pgxpool.ParseConfig(databaseURL)
	if err != nil {
		return nil, err
	}
	return pgxpool.NewWithConfig(ctx, cfg)
}

// ListActiveSubscriptions returns users with status=active and ends_at in the future.
func (s *Store) ListActiveSubscriptions(ctx context.Context, now time.Time) ([]ActiveSubscription, error) {
	const q = `
SELECT telegram_user_id, ends_at
FROM subscriptions
WHERE status = 'active' AND ends_at > $1
ORDER BY ends_at ASC;
`
	rows, err := s.db.Query(ctx, q, now.UTC())
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []ActiveSubscription
	for rows.Next() {
		var item ActiveSubscription
		if err := rows.Scan(&item.UserID, &item.EndsAt); err != nil {
			return nil, err
		}
		out = append(out, item)
	}
	return out, rows.Err()
}

func (s *Store) WasSent(ctx context.Context, userID int64, notificationType Type) (bool, error) {
	const q = `
SELECT 1
FROM subscription_notifications
WHERE user_id = $1 AND notification_type = $2
LIMIT 1;
`
	var one int
	err := s.db.QueryRow(ctx, q, userID, string(notificationType)).Scan(&one)
	if errors.Is(err, pgx.ErrNoRows) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

func (s *Store) MarkSent(ctx context.Context, userID int64, notificationType Type, sentAt time.Time) error {
	const q = `
INSERT INTO subscription_notifications (user_id, notification_type, sent_at)
VALUES ($1, $2, $3)
ON CONFLICT (user_id, notification_type) DO NOTHING;
`
	_, err := s.db.Exec(ctx, q, userID, string(notificationType), sentAt.UTC())
	return err
}
