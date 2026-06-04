package auth

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var ErrInvalidToken = errors.New("invalid token")

type AccessClaims struct {
	TelegramID int64 `json:"tid"`
	jwt.RegisteredClaims
}

func IssueAccess(secret string, telegramID int64, ttl time.Duration) (string, time.Time, error) {
	exp := time.Now().UTC().Add(ttl)
	claims := AccessClaims{
		TelegramID: telegramID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(exp),
			IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
			Subject:   fmt.Sprintf("%d", telegramID),
		},
	}
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	s, err := t.SignedString([]byte(secret))
	return s, exp, err
}

func ParseAccess(secret, token string) (int64, error) {
	parsed, err := jwt.ParseWithClaims(token, &AccessClaims{}, func(t *jwt.Token) (interface{}, error) {
		if t.Method != jwt.SigningMethodHS256 {
			return nil, ErrInvalidToken
		}
		return []byte(secret), nil
	})
	if err != nil {
		return 0, ErrInvalidToken
	}
	claims, ok := parsed.Claims.(*AccessClaims)
	if !ok || !parsed.Valid || claims.TelegramID <= 0 {
		return 0, ErrInvalidToken
	}
	return claims.TelegramID, nil
}
