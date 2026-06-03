package middleware

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

const (
	headerInternalToken  = "X-Internal-Token"
	headerTelegramUserID = "X-Telegram-User-Id"
	headerTelegramFirst  = "X-Telegram-First-Name"
	headerTelegramLast   = "X-Telegram-Last-Name"
	headerTelegramUser   = "X-Telegram-Username"
)

var errInvalidInit = errors.New("invalid init data")

type TelegramIdentity struct {
	ID        int64
	FirstName string
	LastName  string
	Username  string
}

type initDataUser struct {
	ID        int64  `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Username  string `json:"username"`
}

// UserAuth accepts bot calls (internal token + telegram headers) or Mini App (Authorization: tma <initData>).
func UserAuth(internalToken, botToken string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if tok := c.GetHeader(headerInternalToken); tok != "" && tok == internalToken {
			id, err := strconv.ParseInt(strings.TrimSpace(c.GetHeader(headerTelegramUserID)), 10, 64)
			if err != nil || id <= 0 {
				c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid telegram user id"})
				return
			}
			c.Set("telegram", TelegramIdentity{
				ID:        id,
				FirstName: strings.TrimSpace(c.GetHeader(headerTelegramFirst)),
				LastName:  strings.TrimSpace(c.GetHeader(headerTelegramLast)),
				Username:  strings.TrimSpace(strings.TrimPrefix(c.GetHeader(headerTelegramUser), "@")),
			})
			c.Next()
			return
		}

		auth := strings.TrimSpace(c.GetHeader("Authorization"))
		if !strings.HasPrefix(auth, "tma ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}
		initData := strings.TrimSpace(strings.TrimPrefix(auth, "tma "))
		if initData == "" || strings.TrimSpace(botToken) == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}
		id, first, last, user, err := validateInitData(initData, botToken)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid init data"})
			return
		}
		c.Set("telegram", TelegramIdentity{ID: id, FirstName: first, LastName: last, Username: user})
		c.Next()
	}
}

func GetTelegram(c *gin.Context) (TelegramIdentity, bool) {
	v, ok := c.Get("telegram")
	if !ok {
		return TelegramIdentity{}, false
	}
	tg, ok := v.(TelegramIdentity)
	return tg, ok
}

func validateInitData(initData, botToken string) (id int64, first, last, username string, err error) {
	vals, err := url.ParseQuery(initData)
	if err != nil {
		return 0, "", "", "", errInvalidInit
	}
	hash := vals.Get("hash")
	if hash == "" {
		return 0, "", "", "", errInvalidInit
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
		return 0, "", "", "", errInvalidInit
	}

	if authDate := vals.Get("auth_date"); authDate != "" {
		sec, parseErr := strconv.ParseInt(authDate, 10, 64)
		if parseErr == nil && time.Since(time.Unix(sec, 0)) > 24*time.Hour {
			return 0, "", "", "", errInvalidInit
		}
	}

	var u initDataUser
	if err := json.Unmarshal([]byte(vals.Get("user")), &u); err != nil || u.ID <= 0 {
		return 0, "", "", "", errInvalidInit
	}
	return u.ID, u.FirstName, u.LastName, u.Username, nil
}
