package repository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TrialActivationsRepo struct {
	db *pgxpool.Pool
}

func NewTrialActivationsRepo(db *pgxpool.Pool) *TrialActivationsRepo {
	return &TrialActivationsRepo{db: db}
}

func (r *TrialActivationsRepo) HasUsed(ctx context.Context, telegramID int64) (bool, error) {
	const q = `SELECT 1 FROM trial_activations WHERE telegram_id = $1 LIMIT 1;`
	var one int
	err := r.db.QueryRow(ctx, q, telegramID).Scan(&one)
	if errors.Is(err, pgx.ErrNoRows) {
		return false, nil
	}
	return err == nil, err
}

func (r *TrialActivationsRepo) Record(ctx context.Context, telegramID int64) error {
	const q = `INSERT INTO trial_activations (telegram_id) VALUES ($1);`
	_, err := r.db.Exec(ctx, q, telegramID)
	return err
}
