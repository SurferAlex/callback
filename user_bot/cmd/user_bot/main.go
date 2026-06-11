package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"user-bot/internal/bot"
	"user-bot/internal/config"
	"user-bot/internal/notify"
	"user-bot/vpnapi"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Printf(".env not loaded: %v", err)
	}
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	tg, err := tgbotapi.NewBotAPI(cfg.TelegramBotToken)
	if err != nil {
		log.Fatalf("telegram: %v", err)
	}
	log.Printf("user_bot: @%s", tg.Self.UserName)

	api := vpnapi.New(cfg.VPNAPIBaseURL, cfg.VPNAPIInternalToken)
	sender := &bot.Sender{Bot: tg}
	router := bot.NewRouter(cfg, api, sender)

	if cfg.NotifyEnabled && cfg.DatabaseURL != "" {
		pool, err := notify.OpenPool(ctx, cfg.DatabaseURL)
		if err != nil {
			log.Printf("[notify] postgres connect failed, scheduler disabled: %v", err)
		} else {
			defer pool.Close()
			worker := notify.NewWorker(notify.NewStore(pool), sender)
			go notify.StartScheduler(ctx, worker, cfg.NotifyInterval)
		}
	} else if cfg.NotifyEnabled && cfg.DatabaseURL == "" {
		log.Printf("[notify] DATABASE_URL not set, scheduler disabled")
	}

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 30
	updates := tg.GetUpdatesChan(u)

	for {
		select {
		case <-ctx.Done():
			tg.StopReceivingUpdates()
			return
		case upd, ok := <-updates:
			if !ok {
				return
			}
			router.HandleUpdate(ctx, upd)
		}
	}
}
