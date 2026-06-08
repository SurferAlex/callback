package usecase

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"api-vpn/internal/brand"
	"api-vpn/internal/model"
)

// SubscriptionFeed is the Happ / V2Ray-compatible subscription payload.
type SubscriptionFeed struct {
	Body    []byte
	Headers map[string]string
}

// SubscriptionURL builds the public HTTPS subscription link for a token.
func (s *UserService) SubscriptionURL(token string) string {
	token = strings.TrimSpace(token)
	if token == "" {
		return ""
	}
	base := strings.TrimRight(s.subscriptionBaseURL, "/")
	if base == "" {
		base = "https://sub.surfwave.space"
	}
	return base + "/sub/" + token
}

// EnsureSubscriptionToken assigns a unique token to the user when missing.
func (s *UserService) EnsureSubscriptionToken(ctx context.Context, telegramID int64) (string, error) {
	return s.users.EnsureSubscriptionToken(ctx, telegramID)
}

// RenderSubscriptionFeed returns base64-encoded proxy list for Happ (PostgreSQL source of truth).
func (s *UserService) RenderSubscriptionFeed(ctx context.Context, token string) (SubscriptionFeed, error) {
	token = strings.TrimSpace(token)
	if token == "" {
		return SubscriptionFeed{}, ErrNotFound
	}

	user, err := s.users.GetBySubscriptionToken(ctx, token)
	if err != nil {
		return SubscriptionFeed{}, err
	}

	sub, err := s.subs.GetActiveForUser(ctx, user.TelegramID, s.now())
	if err != nil {
		return SubscriptionFeed{}, ErrExpired
	}
	if !sub.EndsAt.After(s.now()) {
		return SubscriptionFeed{}, ErrExpired
	}

	uris, err := s.collectSubscriptionURIs(ctx, user.TelegramID)
	if err != nil {
		return SubscriptionFeed{}, err
	}
	if len(uris) == 0 {
		return SubscriptionFeed{}, ErrNotFound
	}

	payload := strings.Join(uris, "\n")
	encoded := base64.StdEncoding.EncodeToString([]byte(payload))

	headers := map[string]string{
		"Content-Disposition":       `attachment; filename="surf-vpn.txt"`,
		"Profile-Title":             brand.SubscriptionProfileTitle,
		"Profile-Update-Interval":   "12",
		"Subscription-Userinfo":     subscriptionUserinfo(sub.EndsAt),
		"Cache-Control":             "no-store",
	}

	return SubscriptionFeed{
		Body:    []byte(encoded),
		Headers: headers,
	}, nil
}

// collectSubscriptionURIs returns fresh VLESS URIs for the active client.
// Each subscription fetch rebuilds the URI from 3x-ui inbound + vpn_servers settings
// so Happ "update subscription" picks up fingerprint and panel changes without a new UUID.
func (s *UserService) collectSubscriptionURIs(ctx context.Context, telegramID int64) ([]string, error) {
	client, err := s.activeClient(ctx, telegramID)
	if err != nil {
		return nil, err
	}
	if s.xui == nil {
		return nil, fmt.Errorf("xui access not configured")
	}

	var acc model.XUIAccess
	if _, err := s.xui.Get(ctx, client.ClientUUID); errors.Is(err, ErrNotFound) {
		acc, err = s.xui.Provision(ctx, client.ClientUUID)
	} else if err != nil {
		return nil, err
	} else {
		acc, err = s.xui.RebuildVLESSURI(ctx, client)
	}
	if err != nil {
		return nil, err
	}
	if strings.TrimSpace(acc.VLESSURI) == "" {
		return nil, ErrNotFound
	}
	return []string{strings.TrimSpace(acc.VLESSURI)}, nil
}

func subscriptionUserinfo(expiresAt time.Time) string {
	expire := expiresAt.UTC().Unix()
	return fmt.Sprintf("upload=0; download=0; total=0; expire=%s", strconv.FormatInt(expire, 10))
}
