package auth

import (
	"errors"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"
)

var ErrInvalidInit = errors.New("invalid init data")

type initDataUser struct {
	ID        int64  `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Username  string `json:"username"`
}

// ValidateInitData validates Telegram Mini App initData.
func ValidateInitData(initData, botToken string) (id int64, first, last, username string, err error) {
	vals, err := url.ParseQuery(initData)
	if err != nil {
		return 0, "", "", "", ErrInvalidInit
	}
	hash := vals.Get("hash")
	if hash == "" {
		return 0, "", "", "", ErrInvalidInit
	}
	vals.Del("hash")

	var pairs []string
	for k := range vals {
		pairs = append(pairs, k+"="+vals.Get(k))
	}
	sort.Strings(pairs)
	dataCheck := strings.Join(pairs, "\n")

	secret := hmac.New(sha256.New, []byte("WebAppData"))
	secret.Write([]byte(botToken))
	key := secret.Sum(nil)

	mac := hmac.New(sha256.New, key)
	mac.Write([]byte(dataCheck))
	expected := hex.EncodeToString(mac.Sum(nil))
	if !hmac.Equal([]byte(expected), []byte(hash)) {
		return 0, "", "", "", ErrInvalidInit
	}

	if authDate := vals.Get("auth_date"); authDate != "" {
		sec, parseErr := strconv.ParseInt(authDate, 10, 64)
		if parseErr == nil && time.Since(time.Unix(sec, 0)) > 24*time.Hour {
			return 0, "", "", "", ErrInvalidInit
		}
	}

	var u initDataUser
	if err := json.Unmarshal([]byte(vals.Get("user")), &u); err != nil || u.ID <= 0 {
		return 0, "", "", "", ErrInvalidInit
	}
	return u.ID, u.FirstName, u.LastName, u.Username, nil
}
