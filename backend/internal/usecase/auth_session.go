package usecase

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"time"

	"api-vpn/internal/auth"
	"api-vpn/internal/model"

	"github.com/jackc/pgx/v5"
)

var ErrInvalidRefresh = errors.New("invalid refresh token")

type AuthRefreshRepo interface {
	Store(ctx context.Context, telegramUserID int64, tokenHash string, expiresAt time.Time) error
	FindTelegramID(ctx context.Context, tokenHash string, now time.Time) (int64, error)
	Revoke(ctx context.Context, tokenHash string) error
	RevokeAllForUser(ctx context.Context, telegramUserID int64) error
}

type AuthSession struct {
	users    UsersRepo
	refresh  AuthRefreshRepo
	jwtSecret string
	accessTTL time.Duration
	refreshTTL time.Duration
	now       func() time.Time
}

func NewAuthSession(users UsersRepo, refresh AuthRefreshRepo, jwtSecret string, accessTTL, refreshTTL time.Duration) *AuthSession {
	if accessTTL <= 0 {
		accessTTL = 15 * time.Minute
	}
	if refreshTTL <= 0 {
		refreshTTL = 30 * 24 * time.Hour
	}
	return &AuthSession{
		users:      users,
		refresh:    refresh,
		jwtSecret:  jwtSecret,
		accessTTL:  accessTTL,
		refreshTTL: refreshTTL,
		now:        time.Now,
	}
}

type TokenPair struct {
	AccessToken string
	AccessExp   time.Time
	RefreshRaw  string
	RefreshExp  time.Time
	TelegramID  int64
}

func (s *AuthSession) IssueForTelegram(ctx context.Context, p model.UpsertUserParams) (TokenPair, error) {
	if _, err := s.users.Upsert(ctx, p); err != nil {
		return TokenPair{}, err
	}
	return s.issue(ctx, p.TelegramID)
}

func (s *AuthSession) Refresh(ctx context.Context, refreshRaw string) (TokenPair, error) {
	hash := hashToken(refreshRaw)
	id, err := s.refresh.FindTelegramID(ctx, hash, s.now())
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return TokenPair{}, ErrInvalidRefresh
		}
		return TokenPair{}, err
	}
	_ = s.refresh.Revoke(ctx, hash)
	return s.issue(ctx, id)
}

func (s *AuthSession) Logout(ctx context.Context, refreshRaw string) error {
	if refreshRaw == "" {
		return nil
	}
	return s.refresh.Revoke(ctx, hashToken(refreshRaw))
}

func (s *AuthSession) ParseAccess(token string) (int64, error) {
	return auth.ParseAccess(s.jwtSecret, token)
}

func (s *AuthSession) issue(ctx context.Context, telegramID int64) (TokenPair, error) {
	access, exp, err := auth.IssueAccess(s.jwtSecret, telegramID, s.accessTTL)
	if err != nil {
		return TokenPair{}, err
	}
	raw, rhash, rexp, err := newRefreshToken(s.refreshTTL)
	if err != nil {
		return TokenPair{}, err
	}
	if err := s.refresh.Store(ctx, telegramID, rhash, rexp); err != nil {
		return TokenPair{}, err
	}
	return TokenPair{
		AccessToken: access,
		AccessExp:   exp,
		RefreshRaw:  raw,
		RefreshExp:  rexp,
		TelegramID:  telegramID,
	}, nil
}

func newRefreshToken(ttl time.Duration) (raw, hash string, exp time.Time, err error) {
	b := make([]byte, 32)
	if _, err = rand.Read(b); err != nil {
		return "", "", time.Time{}, err
	}
	raw = base64.RawURLEncoding.EncodeToString(b)
	exp = time.Now().UTC().Add(ttl)
	return raw, hashToken(raw), exp, nil
}

func hashToken(raw string) string {
	sum := sha256.Sum256([]byte(raw))
	return base64.RawURLEncoding.EncodeToString(sum[:])
}
