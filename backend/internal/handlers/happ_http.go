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

	lower := strings.ToLower(key)
	target := key
	if strings.HasPrefix(lower, "happ://") {
		target = key
	} else {
		// VLESS key, subscription https:// URL, or plain config string.
		target = "happ://add/" + key
	}

	c.Redirect(http.StatusFound, target)
}
