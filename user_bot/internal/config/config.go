package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
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
	DatabaseURL         string
	NotifyEnabled       bool
	NotifyInterval      time.Duration
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
	cfg.DatabaseURL = strings.TrimSpace(os.Getenv("DATABASE_URL"))
	cfg.NotifyEnabled = envBoolDefault("NOTIFY_SCHEDULER_ENABLED", true)
	cfg.NotifyInterval = envDurationDefault("NOTIFY_SCHEDULER_INTERVAL", 30*time.Minute)
	return cfg, nil
}

func envBoolDefault(key string, def bool) bool {
	v := strings.TrimSpace(os.Getenv(key))
	if v == "" {
		return def
	}
	b, err := strconv.ParseBool(v)
	if err != nil {
		return def
	}
	return b
}

func envDurationDefault(key string, def time.Duration) time.Duration {
	v := strings.TrimSpace(os.Getenv(key))
	if v == "" {
		return def
	}
	d, err := time.ParseDuration(v)
	if err != nil {
		return def
	}
	return d
}
