package config

import (
	"fmt"
	"os"
	"strings"
)

type Config struct {
	TelegramBotToken    string
	VPNAPIBaseURL       string
	VPNAPIInternalToken string
	MiniAppURL          string
	SupportContact      string
	SupportTelegramURL  string
	WelcomeStickerID    string
	SuccessStickerID    string
	ErrorStickerID      string
}

func Load() (Config, error) {
	cfg := Config{
		TelegramBotToken:    strings.TrimSpace(os.Getenv("TELEGRAM_BOT_TOKEN")),
		VPNAPIBaseURL:       strings.TrimRight(strings.TrimSpace(os.Getenv("VPNAPI_BASE_URL")), "/"),
		VPNAPIInternalToken: strings.TrimSpace(os.Getenv("VPNAPI_INTERNAL_TOKEN")),
		MiniAppURL:          strings.TrimSpace(os.Getenv("MINI_APP_URL")),
		SupportContact:      strings.TrimSpace(os.Getenv("SUPPORT_CONTACT")),
		SupportTelegramURL:  strings.TrimSpace(os.Getenv("SUPPORT_TELEGRAM_URL")),
		WelcomeStickerID:    strings.TrimSpace(os.Getenv("WELCOME_STICKER_ID")),
		SuccessStickerID:    strings.TrimSpace(os.Getenv("SUCCESS_STICKER_ID")),
		ErrorStickerID:      strings.TrimSpace(os.Getenv("ERROR_STICKER_ID")),
	}
	if cfg.TelegramBotToken == "" {
		return Config{}, fmt.Errorf("TELEGRAM_BOT_TOKEN is required")
	}
	if cfg.VPNAPIBaseURL == "" {
		return Config{}, fmt.Errorf("VPNAPI_BASE_URL is required")
	}
	if cfg.VPNAPIInternalToken == "" {
		return Config{}, fmt.Errorf("VPNAPI_INTERNAL_TOKEN is required")
	}
	if cfg.SupportContact == "" {
		cfg.SupportContact = "@surfervpn_support"
	}
	if cfg.SupportTelegramURL == "" {
		cfg.SupportTelegramURL = "https://t.me/surfervpn_support"
	}
	return cfg, nil
}
