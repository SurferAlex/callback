package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"user-bot/internal/botapp"
	"user-bot/internal/config"
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

	bot, err := tgbotapi.NewBotAPI(cfg.TelegramBotToken)
	if err != nil {
		log.Fatalf("telegram: %v", err)
	}
	log.Printf("user_bot: @%s", bot.Self.UserName)

	api := vpnapi.New(cfg.VPNAPIBaseURL, cfg.VPNAPIInternalToken)
	menu := botapp.MainMenuMarkup()

	type pendingPlan struct {
		plan   string
		chatID int64
		at     time.Time
	}
	pending := map[int64]pendingPlan{}

	send := func(chatID int64, text string, markup interface{}) {
		msg := tgbotapi.NewMessage(chatID, text)
		if markup != nil {
			switch m := markup.(type) {
			case tgbotapi.ReplyKeyboardMarkup:
				msg.ReplyMarkup = m
			case tgbotapi.InlineKeyboardMarkup:
				msg.ReplyMarkup = m
			}
		}
		if _, err := bot.Send(msg); err != nil {
			log.Printf("send failed chat=%d: %v", chatID, err)
		}
	}

	sendConfig := func(chatID int64, from *tgbotapi.User, prefix string) {
		if from == nil {
			return
		}
		cfgResp, err := api.GetConfig(ctx, from.ID, from.FirstName, from.LastName, from.UserName)
		if err != nil {
			if vpnapi.IsNoSubscription(err) {
				send(chatID, "Активной подписки нет.\nОформите тестовую подписку: «💳 Купить VPN» → тариф → «✅ Тестовая активация».", botapp.MockPayMarkup())
				return
			}
			log.Printf("get config: %v", err)
			send(chatID, "Не удалось получить конфиг. Попробуйте позже.", menu)
			return
		}
		text := prefix + "\n\n<code>" + escapeHTML(cfgResp.VlessURI) + "</code>"
		msg := tgbotapi.NewMessage(chatID, text)
		msg.ParseMode = tgbotapi.ModeHTML
		msg.ReplyMarkup = menu
		_, _ = bot.Send(msg)
	}

	cabinetMarkup := func() tgbotapi.InlineKeyboardMarkup {
		if cfg.MiniAppURL == "" {
			return tgbotapi.InlineKeyboardMarkup{}
		}
		return tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonURL("Открыть личный кабинет", cfg.MiniAppURL),
			),
		)
	}

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 30
	updates := bot.GetUpdatesChan(u)

	for {
		select {
		case <-ctx.Done():
			bot.StopReceivingUpdates()
			return
		case update, ok := <-updates:
			if !ok {
				return
			}
			if update.Message == nil || update.Message.From == nil {
				continue
			}
			from := update.Message.From
			chatID := update.Message.Chat.ID
			text := strings.TrimSpace(update.Message.Text)

			if update.Message.IsCommand() {
				switch update.Message.Command() {
				case "start":
					delete(pending, from.ID)
					send(chatID, botapp.WelcomeText(from.FirstName), menu)
					continue
				}
			}

			if text == botapp.BtnBack {
				delete(pending, from.ID)
				send(chatID, "Главное меню:", menu)
				continue
			}

			if p, ok := pending[from.ID]; ok && text == botapp.BtnMockPay {
				if time.Since(p.at) > 10*time.Minute {
					delete(pending, from.ID)
					send(chatID, "Сессия истекла. Выберите тариф снова.", botapp.PlansMarkup())
					continue
				}
				delete(pending, from.ID)
				me, err := api.MockActivate(ctx, from.ID, from.FirstName, from.LastName, from.UserName, p.plan)
				if err != nil {
					log.Printf("mock activate: %v", err)
					send(chatID, "Не удалось активировать подписку. Попробуйте позже.", menu)
					continue
				}
				info := fmt.Sprintf(
					"✅ Тестовая подписка активирована!\n\nТариф: %s\nДействует до: %s\nОсталось дней: %d",
					me.Subscription.Plan,
					formatDate(me.Subscription.ExpiresAt),
					me.Subscription.DaysLeft,
				)
				send(chatID, info, menu)
				if me.VpnKey != "" {
					msg := tgbotapi.NewMessage(chatID, "Ваш конфиг:\n\n<code>"+escapeHTML(me.VpnKey)+"</code>")
					msg.ParseMode = tgbotapi.ModeHTML
					msg.ReplyMarkup = menu
					_, _ = bot.Send(msg)
				}
				continue
			}

			if code := botapp.PlanCode(text); code != "" {
				pending[from.ID] = pendingPlan{plan: code, chatID: chatID, at: time.Now()}
				send(chatID, "Тариф: "+text+"\n\nНажмите «✅ Тестовая активация» для mock-оплаты.", botapp.MockPayMarkup())
				continue
			}

			switch text {
			case botapp.BtnBuy:
				delete(pending, from.ID)
				send(chatID, "Выберите срок подписки:", botapp.PlansMarkup())
			case botapp.BtnGetConfig:
				sendConfig(chatID, from, "🔑 Ваш VPN-конфиг:")
			case botapp.BtnRefresh:
				cfgResp, err := api.RefreshConfig(ctx, from.ID, from.FirstName, from.LastName, from.UserName)
				if err != nil {
					if vpnapi.IsNoSubscription(err) {
						send(chatID, "Нет активной подписки. Оформите тестовую через «💳 Купить VPN».", botapp.PlansMarkup())
						break
					}
					log.Printf("refresh: %v", err)
					send(chatID, "Не удалось обновить конфиг.", menu)
					break
				}
				msg := tgbotapi.NewMessage(chatID, "🔄 Новый конфиг:\n\n<code>"+escapeHTML(cfgResp.VlessURI)+"</code>")
				msg.ParseMode = tgbotapi.ModeHTML
				msg.ReplyMarkup = menu
				_, _ = bot.Send(msg)
			case botapp.BtnGuide:
				send(chatID, "Выберите платформу:", botapp.GuideMarkup())
			case botapp.BtnGuideIOS:
				send(chatID, guideIOS(), menu)
			case botapp.BtnGuideAndroid:
				send(chatID, guideAndroid(), menu)
			case botapp.BtnGuideWin:
				send(chatID, guideWindows(), menu)
			case botapp.BtnGuideMac:
				send(chatID, guideMacOS(), menu)
			case botapp.BtnSupport:
				send(chatID, fmt.Sprintf(
					"🛟 Поддержка\n\nКонтакт: %s\nTelegram: %s",
					cfg.SupportContact,
					cfg.SupportTelegramURL,
				), menu)
			case botapp.BtnCabinet:
				if cfg.MiniAppURL == "" {
					send(chatID, "Личный кабинет скоро будет доступен. Укажите MINI_APP_URL в настройках бота.", menu)
				} else {
					send(chatID, "Откройте личный кабинет:", cabinetMarkup())
				}
			}
		}
	}
}

func escapeHTML(s string) string {
	s = strings.ReplaceAll(s, "&", "&amp;")
	s = strings.ReplaceAll(s, "<", "&lt;")
	s = strings.ReplaceAll(s, ">", "&gt;")
	return s
}

func formatDate(iso string) string {
	if iso == "" {
		return "—"
	}
	t, err := time.Parse(time.RFC3339, iso)
	if err != nil {
		return iso
	}
	return t.Format("02.01.2006")
}

func guideIOS() string {
	return "📱 iPhone (заглушка)\n\n1. Установите Happ из App Store\n2. Скопируйте конфиг из бота\n3. Откройте Happ → добавить подписку"
}

func guideAndroid() string {
	return "📱 Android (заглушка)\n\n1. Установите Happ\n2. Получите конфиг в боте\n3. Импортируйте ключ в приложение"
}

func guideWindows() string {
	return "🖥 Windows (заглушка)\n\n1. Установите совместимый VPN-клиент\n2. Вставьте VLESS-конфиг из бота"
}

func guideMacOS() string {
	return "🍎 macOS (заглушка)\n\n1. Установите Happ\n2. Импортируйте конфиг из раздела «Получить конфиг»"
}
