package handlers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// HappOpen redirects to happ://add/{vless-key}. Used by Telegram Mini App (HTTPS bridge).
func (h *Handlers) HappOpen(c *gin.Context) {
	key := strings.TrimSpace(c.Query("key"))
	if key == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "key required"})
		return
	}

	target := key
	if !strings.HasPrefix(strings.ToLower(key), "happ://") {
		target = "happ://add/" + key
	}

	c.Redirect(http.StatusFound, target)
}
