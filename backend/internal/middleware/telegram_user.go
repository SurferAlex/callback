package middleware

import (
	"net/http"
	"strconv"
	"strings"

	"api-vpn/internal/auth"
	"api-vpn/internal/usecase"

	"github.com/gin-gonic/gin"
)

const (
	headerInternalToken  = "X-Internal-Token"
	headerTelegramUserID = "X-Telegram-User-Id"
	headerTelegramFirst  = "X-Telegram-First-Name"
	headerTelegramLast   = "X-Telegram-Last-Name"
	headerTelegramUser   = "X-Telegram-Username"
)

type TelegramIdentity struct {
	ID        int64
	FirstName string
	LastName  string
	Username  string
}

// UserAuth: Bearer JWT, internal token, or Mini App tma initData.
func UserAuth(internalToken, botToken string, sessions *usecase.AuthSession) gin.HandlerFunc {
	return func(c *gin.Context) {
		if sessions != nil {
			if authz := strings.TrimSpace(c.GetHeader("Authorization")); strings.HasPrefix(authz, "Bearer ") {
				token := strings.TrimSpace(strings.TrimPrefix(authz, "Bearer "))
				id, err := sessions.ParseAccess(token)
				if err == nil && id > 0 {
					c.Set("telegram", TelegramIdentity{ID: id})
					c.Next()
					return
				}
			}
		}

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

		authHeader := strings.TrimSpace(c.GetHeader("Authorization"))
		if !strings.HasPrefix(authHeader, "tma ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}
		initData := strings.TrimSpace(strings.TrimPrefix(authHeader, "tma "))
		if initData == "" || strings.TrimSpace(botToken) == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}
		id, first, last, user, err := auth.ValidateInitData(initData, botToken)
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
