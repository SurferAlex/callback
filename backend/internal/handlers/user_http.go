package handlers

import (
	"errors"
	"log"
	"net/http"
	"strings"
	"time"

	"api-vpn/internal/middleware"
	"api-vpn/internal/model"
	"api-vpn/internal/usecase"

	"github.com/gin-gonic/gin"
)

type MockActivateRequest struct {
	Plan string `json:"plan" binding:"required"`
}

type UserMeResponse struct {
	TelegramID   int64              `json:"telegramId"`
	FirstName    string             `json:"firstName"`
	LastName     *string            `json:"lastName,omitempty"`
	Username     *string            `json:"username,omitempty"`
	VpnKey       string             `json:"vpnKey"`
	Subscription SubscriptionView   `json:"subscription"`
	Server       ServerView         `json:"server"`
}

type SubscriptionView struct {
	Status    string `json:"status"`
	Plan      string `json:"plan"`
	ExpiresAt string `json:"expiresAt"`
	DaysLeft  int    `json:"daysLeft"`
	AutoRenew bool   `json:"autoRenew"`
}

type ServerView struct {
	ID          string `json:"id"`
	City        string `json:"city"`
	CountryCode string `json:"countryCode"`
	Country     string `json:"country"`
}

type ConfigResponse struct {
	VlessURI   string `json:"vlessUri"`
	ClientUUID string `json:"clientUuid,omitempty"`
}

func profileFromTelegram(tg middleware.TelegramIdentity) model.UpsertUserParams {
	var last, user *string
	if tg.LastName != "" {
		last = &tg.LastName
	}
	if tg.Username != "" {
		user = &tg.Username
	}
	return model.UpsertUserParams{
		TelegramID: tg.ID,
		FirstName:  tg.FirstName,
		LastName:   last,
		Username:   user,
	}
}

func toUserMeResponse(p usecase.UserProfile) UserMeResponse {
	resp := UserMeResponse{
		TelegramID: p.User.TelegramID,
		FirstName:  p.User.FirstName,
		LastName:   p.User.LastName,
		Username:   p.User.Username,
		Subscription: SubscriptionView{
			Status:    "none",
			Plan:      "",
			ExpiresAt: "",
			DaysLeft:  0,
			AutoRenew: false,
		},
		Server: ServerView{
			ID:          "default",
			City:        "VPN",
			CountryCode: "XX",
			Country:     "Сервер",
		},
	}
	if p.Client != nil {
		resp.Server.ID = p.Client.ServerID
	}
	if p.Subscription != nil {
		resp.Subscription.Status = "active"
		resp.Subscription.Plan = p.Subscription.PlanLabel
		resp.Subscription.ExpiresAt = p.Subscription.EndsAt.UTC().Format(time.RFC3339)
		resp.Subscription.DaysLeft = daysLeft(p.Subscription.EndsAt)
	}
	if p.Client != nil && p.Client.KeyExpiresAt.After(time.Now()) {
		if p.Subscription == nil {
			resp.Subscription.Status = "active"
			resp.Subscription.ExpiresAt = p.Client.KeyExpiresAt.UTC().Format(time.RFC3339)
			resp.Subscription.Plan = "VPN"
			resp.Subscription.DaysLeft = daysLeft(p.Client.KeyExpiresAt)
		}
	} else if p.Subscription == nil {
		resp.Subscription.Status = "none"
	}
	if p.Access != nil {
		resp.VpnKey = p.Access.VLESSURI
	}
	if resp.FirstName == "" {
		resp.FirstName = "Пользователь"
	}
	return resp
}

func daysLeft(t time.Time) int {
	ms := t.Sub(time.Now())
	if ms <= 0 {
		return 0
	}
	return int((ms + 24*time.Hour - 1) / (24 * time.Hour))
}

func (h *Handlers) UserMe(c *gin.Context) {
	tg, ok := middleware.GetTelegram(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	_, _ = h.Users.UpsertUser(c.Request.Context(), profileFromTelegram(tg))
	prof, err := h.Users.GetProfile(c.Request.Context(), tg.ID)
	if err != nil {
		log.Printf("user me failed: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}
	if prof.User.TelegramID == 0 {
		prof.User.TelegramID = tg.ID
		prof.User.FirstName = tg.FirstName
	}
	c.JSON(http.StatusOK, toUserMeResponse(prof))
}

func (h *Handlers) UserMockActivate(c *gin.Context) {
	tg, ok := middleware.GetTelegram(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	var req MockActivateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	prof, err := h.Users.MockActivate(c.Request.Context(), tg.ID, strings.TrimSpace(req.Plan), profileFromTelegram(tg))
	if err != nil {
		switch {
		case errors.Is(err, usecase.ErrInvalidPlan):
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid plan"})
		case errors.Is(err, usecase.ErrInvalidServer):
			c.JSON(http.StatusBadRequest, gin.H{"error": "vpn server not configured"})
		default:
			log.Printf("mock activate failed: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		}
		return
	}
	c.JSON(http.StatusOK, toUserMeResponse(prof))
}

func (h *Handlers) UserGetConfig(c *gin.Context) {
	tg, ok := middleware.GetTelegram(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	acc, err := h.Users.GetConfig(c.Request.Context(), tg.ID)
	if err != nil {
		if errors.Is(err, usecase.ErrNotFound) || errors.Is(err, usecase.ErrInactive) || errors.Is(err, usecase.ErrExpired) {
			c.JSON(http.StatusNotFound, gin.H{"error": "no active subscription"})
			return
		}
		log.Printf("get config failed: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}
	c.JSON(http.StatusOK, ConfigResponse{
		VlessURI:   acc.VLESSURI,
		ClientUUID: acc.ClientUUID.String(),
	})
}

func (h *Handlers) UserRefreshConfig(c *gin.Context) {
	tg, ok := middleware.GetTelegram(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	acc, err := h.Users.RefreshConfig(c.Request.Context(), tg.ID)
	if err != nil {
		if errors.Is(err, usecase.ErrNotFound) || errors.Is(err, usecase.ErrInactive) || errors.Is(err, usecase.ErrExpired) {
			c.JSON(http.StatusNotFound, gin.H{"error": "no active subscription"})
			return
		}
		log.Printf("refresh config failed: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}
	c.JSON(http.StatusOK, ConfigResponse{
		VlessURI:   acc.VLESSURI,
		ClientUUID: acc.ClientUUID.String(),
	})
}
