package vpnapi

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type Client struct {
	base  string
	token string
	http  *http.Client
}

func New(baseURL, token string) *Client {
	return &Client{
		base:  strings.TrimRight(baseURL, "/"),
		token: token,
		http:  &http.Client{Timeout: 45 * time.Second},
	}
}

type UserMe struct {
	TelegramID   int64            `json:"telegramId"`
	FirstName    string           `json:"firstName"`
	VpnKey       string           `json:"vpnKey"`
	Subscription SubscriptionView `json:"subscription"`
}

type SubscriptionView struct {
	Status    string `json:"status"`
	Plan      string `json:"plan"`
	ExpiresAt string `json:"expiresAt"`
	DaysLeft  int    `json:"daysLeft"`
}

type ConfigResponse struct {
	VlessURI string `json:"vlessUri"`
}

func (c *Client) setUserHeaders(req *http.Request, userID int64, first, last, username string) {
	req.Header.Set("X-Internal-Token", c.token)
	req.Header.Set("X-Telegram-User-Id", strconv.FormatInt(userID, 10))
	if first != "" {
		req.Header.Set("X-Telegram-First-Name", first)
	}
	if last != "" {
		req.Header.Set("X-Telegram-Last-Name", last)
	}
	if username != "" {
		req.Header.Set("X-Telegram-Username", username)
	}
}

func (c *Client) GetMe(ctx context.Context, userID int64, first, last, username string) (UserMe, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.base+"/api/v1/user/me", nil)
	if err != nil {
		return UserMe{}, err
	}
	c.setUserHeaders(req, userID, first, last, username)
	resp, err := c.http.Do(req)
	if err != nil {
		return UserMe{}, err
	}
	defer resp.Body.Close()
	return decodeJSON[UserMe](resp)
}

func (c *Client) MockActivate(ctx context.Context, userID int64, first, last, username, plan string) (UserMe, error) {
	body, _ := json.Marshal(map[string]string{"plan": plan})
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.base+"/api/v1/user/subscription/mock-activate", bytes.NewReader(body))
	if err != nil {
		return UserMe{}, err
	}
	req.Header.Set("Content-Type", "application/json")
	c.setUserHeaders(req, userID, first, last, username)
	resp, err := c.http.Do(req)
	if err != nil {
		return UserMe{}, err
	}
	defer resp.Body.Close()
	return decodeJSON[UserMe](resp)
}

func (c *Client) GetConfig(ctx context.Context, userID int64, first, last, username string) (ConfigResponse, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.base+"/api/v1/user/config", nil)
	if err != nil {
		return ConfigResponse{}, err
	}
	c.setUserHeaders(req, userID, first, last, username)
	resp, err := c.http.Do(req)
	if err != nil {
		return ConfigResponse{}, err
	}
	defer resp.Body.Close()
	return decodeJSON[ConfigResponse](resp)
}

func (c *Client) RefreshConfig(ctx context.Context, userID int64, first, last, username string) (ConfigResponse, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.base+"/api/v1/user/config/refresh", nil)
	if err != nil {
		return ConfigResponse{}, err
	}
	c.setUserHeaders(req, userID, first, last, username)
	resp, err := c.http.Do(req)
	if err != nil {
		return ConfigResponse{}, err
	}
	defer resp.Body.Close()
	return decodeJSON[ConfigResponse](resp)
}

func decodeJSON[T any](resp *http.Response) (T, error) {
	var zero T
	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		var errBody struct {
			Error string `json:"error"`
		}
		_ = json.Unmarshal(body, &errBody)
		msg := errBody.Error
		if msg == "" {
			msg = fmt.Sprintf("status %d", resp.StatusCode)
		}
		return zero, fmt.Errorf("vpnapi: %s", msg)
	}
	var out T
	if err := json.Unmarshal(body, &out); err != nil {
		return zero, err
	}
	return out, nil
}

func IsNoSubscription(err error) bool {
	return err != nil && strings.Contains(err.Error(), "no active subscription")
}

func IsTrialAlreadyUsed(err error) bool {
	return err != nil && strings.Contains(err.Error(), "trial already used")
}

func IsTrialActiveSubscription(err error) bool {
	return err != nil && strings.Contains(err.Error(), "active subscription exists")
}

func (c *Client) ActivateTrial(ctx context.Context, userID int64, first, last, username string) (UserMe, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.base+"/api/v1/user/trial/activate", nil)
	if err != nil {
		return UserMe{}, err
	}
	c.setUserHeaders(req, userID, first, last, username)
	resp, err := c.http.Do(req)
	if err != nil {
		return UserMe{}, err
	}
	defer resp.Body.Close()
	return decodeJSON[UserMe](resp)
}
