package middleware

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

// InternalUserAuth requires X-Internal-Token and Telegram identity headers (for bot/admin).
func InternalUserAuth(internalToken string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if internalToken == "" {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "internal token is not configured"})
			return
		}
		if c.GetHeader(headerInternalToken) != internalToken {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}
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
	}
}
