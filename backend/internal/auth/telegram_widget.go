package auth

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"
)

var ErrInvalidWidget = errors.New("invalid telegram login")

type WidgetUser struct {
	ID        int64
	FirstName string
	LastName  string
	Username  string
	PhotoURL  string
}

// VerifyWidget checks Telegram Login Widget payload.
// https://core.telegram.org/widgets/login#checking-authorization
func VerifyWidget(botToken string, fields map[string]string, maxAge time.Duration) (WidgetUser, error) {
	receivedHash := strings.TrimSpace(fields["hash"])
	if receivedHash == "" {
		return WidgetUser{}, ErrInvalidWidget
	}

	var pairs []string
	for k, v := range fields {
		if k == "hash" || strings.TrimSpace(v) == "" {
			continue
		}
		pairs = append(pairs, k+"="+v)
	}
	sort.Strings(pairs)
	dataCheck := strings.Join(pairs, "\n")

	secret := sha256.Sum256([]byte(botToken))
	mac := hmac.New(sha256.New, secret[:])
	mac.Write([]byte(dataCheck))
	expected := hex.EncodeToString(mac.Sum(nil))
	if !hmac.Equal([]byte(expected), []byte(receivedHash)) {
		return WidgetUser{}, ErrInvalidWidget
	}

	if authDate := strings.TrimSpace(fields["auth_date"]); authDate != "" && maxAge > 0 {
		sec, err := strconv.ParseInt(authDate, 10, 64)
		if err != nil || time.Since(time.Unix(sec, 0)) > maxAge {
			return WidgetUser{}, ErrInvalidWidget
		}
	}

	id, err := strconv.ParseInt(strings.TrimSpace(fields["id"]), 10, 64)
	if err != nil || id <= 0 {
		return WidgetUser{}, ErrInvalidWidget
	}

	return WidgetUser{
		ID:        id,
		FirstName: strings.TrimSpace(fields["first_name"]),
		LastName:  strings.TrimSpace(fields["last_name"]),
		Username:  strings.TrimSpace(fields["username"]),
		PhotoURL:  strings.TrimSpace(fields["photo_url"]),
	}, nil
}

func WidgetFieldsFromMap(m map[string]interface{}) map[string]string {
	out := make(map[string]string)
	for k, v := range m {
		out[k] = fmt.Sprint(v)
	}
	return out
}
