package handlers

import (
	"errors"
	"log"
	"net/http"

	"api-vpn/internal/usecase"

	"github.com/gin-gonic/gin"
)

// SubscriptionFeed serves Happ-compatible subscription content for a personal token.
func (h *Handlers) SubscriptionFeed(c *gin.Context) {
	token := c.Param("token")
	feed, err := h.Users.RenderSubscriptionFeed(c.Request.Context(), token)
	if err != nil {
		switch {
		case errors.Is(err, usecase.ErrNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		case errors.Is(err, usecase.ErrExpired):
			c.JSON(http.StatusForbidden, gin.H{"error": "subscription expired"})
		default:
			log.Printf("[sub/feed] token=%s err=%v", token, err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		}
		return
	}
	for k, v := range feed.Headers {
		c.Header(k, v)
	}
	c.Data(http.StatusOK, "text/plain; charset=utf-8", feed.Body)
}
