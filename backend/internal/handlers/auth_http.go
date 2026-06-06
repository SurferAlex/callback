package handlers

import (
	"errors"
	"log"
	"net/http"
	"strings"
	"time"

	"api-vpn/internal/auth"
	"api-vpn/internal/model"
	"api-vpn/internal/usecase"

	"github.com/gin-gonic/gin"
)

const refreshCookieName = "surf_refresh"

type AuthHandlers struct {
	Auth      *usecase.AuthSession
	Users     *usecase.UserService
	BotToken  string
	CookieDomain string
	CookieSecure bool
	RefreshTTL   time.Duration
}

type tokenResponse struct {
	AccessToken string  `json:"accessToken"`
	ExpiresIn   int64   `json:"expiresIn"`
	FirstName   string  `json:"firstName,omitempty"`
	Username    *string `json:"username,omitempty"`
}

func (h *AuthHandlers) setRefreshCookie(c *gin.Context, raw string, exp time.Time) {
	c.SetSameSite(http.SameSiteLaxMode)
	secure := h.CookieSecure
	c.SetCookie(refreshCookieName, raw, int(time.Until(exp).Seconds()), "/", h.CookieDomain, secure, true)
}

func (h *AuthHandlers) clearRefreshCookie(c *gin.Context) {
	c.SetCookie(refreshCookieName, "", -1, "/", h.CookieDomain, h.CookieSecure, true)
}

func (h *AuthHandlers) writeTokens(c *gin.Context, pair usecase.TokenPair) {
	h.setRefreshCookie(c, pair.RefreshRaw, pair.RefreshExp)
	resp := tokenResponse{
		AccessToken: pair.AccessToken,
		ExpiresIn:   int64(time.Until(pair.AccessExp).Seconds()),
	}
	if prof, err := h.Users.GetProfile(c.Request.Context(), pair.TelegramID); err == nil {
		if strings.TrimSpace(prof.User.FirstName) != "" {
			resp.FirstName = strings.TrimSpace(prof.User.FirstName)
		}
		resp.Username = prof.User.Username
	}
	c.JSON(http.StatusOK, resp)
}

// SessionTelegramWebApp exchanges Mini App initData for JWT + refresh cookie.
func (h *AuthHandlers) SessionTelegramWebApp(c *gin.Context) {
	authHeader := strings.TrimSpace(c.GetHeader("Authorization"))
	if !strings.HasPrefix(authHeader, "tma ") {
		c.JSON(http.StatusBadRequest, gin.H{"error": "tma authorization required"})
		return
	}
	initData := strings.TrimSpace(strings.TrimPrefix(authHeader, "tma "))
	id, first, last, user, err := auth.ValidateInitData(initData, h.BotToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid init data"})
		return
	}
	var lastP, userP *string
	if last != "" {
		lastP = &last
	}
	if user != "" {
		userP = &user
	}
	pair, err := h.Auth.IssueForTelegram(c.Request.Context(), model.UpsertUserParams{
		TelegramID: id,
		FirstName:  first,
		LastName:   lastP,
		Username:   userP,
	})
	if err != nil {
		log.Printf("session webapp: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}
	h.writeTokens(c, pair)
}

// SessionTelegramWidget exchanges Telegram Login Widget for JWT + refresh cookie.
func (h *AuthHandlers) SessionTelegramWidget(c *gin.Context) {
	var body map[string]interface{}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	fields := auth.WidgetFieldsFromMap(body)
	wu, err := auth.VerifyWidget(h.BotToken, fields, 24*time.Hour)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid telegram login"})
		return
	}
	var lastP, userP *string
	if wu.LastName != "" {
		lastP = &wu.LastName
	}
	if wu.Username != "" {
		userP = &wu.Username
	}
	pair, err := h.Auth.IssueForTelegram(c.Request.Context(), model.UpsertUserParams{
		TelegramID: wu.ID,
		FirstName:  wu.FirstName,
		LastName:   lastP,
		Username:   userP,
	})
	if err != nil {
		log.Printf("session widget: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}
	h.writeTokens(c, pair)
}

// Refresh rotates refresh cookie and returns new access token.
func (h *AuthHandlers) Refresh(c *gin.Context) {
	raw, err := c.Cookie(refreshCookieName)
	if err != nil || strings.TrimSpace(raw) == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "refresh required"})
		return
	}
	pair, err := h.Auth.Refresh(c.Request.Context(), raw)
	if err != nil {
		if errors.Is(err, usecase.ErrInvalidRefresh) {
			h.clearRefreshCookie(c)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid refresh"})
			return
		}
		log.Printf("refresh: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}
	h.writeTokens(c, pair)
}

// Logout revokes refresh token and clears cookie.
func (h *AuthHandlers) Logout(c *gin.Context) {
	raw, _ := c.Cookie(refreshCookieName)
	_ = h.Auth.Logout(c.Request.Context(), raw)
	h.clearRefreshCookie(c)
	c.Status(http.StatusNoContent)
}

