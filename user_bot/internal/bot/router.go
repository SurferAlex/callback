package bot

import (
	"context"
	"log"
	"strings"

	"user-bot/internal/botapp"
	"user-bot/internal/config"
	"user-bot/vpnapi"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Router struct {
	Cfg    config.Config
	API    *vpnapi.Client
	Send   *Sender
	Plans  map[string]botapp.PlanOffer
}

func NewRouter(cfg config.Config, api *vpnapi.Client, send *Sender) *Router {
	plans := make(map[string]botapp.PlanOffer)
	for _, p := range botapp.LoadPlans() {
		plans[p.Code] = p
	}
	return &Router{Cfg: cfg, API: api, Send: send, Plans: plans}
}

func (r *Router) HandleUpdate(ctx context.Context, update tgbotapi.Update) {
	if update.CallbackQuery != nil {
		r.handleCallback(ctx, update.CallbackQuery)
		return
	}
	if update.Message == nil || update.Message.From == nil {
		return
	}
	if update.Message.IsCommand() && update.Message.Command() == "start" {
		r.handleStart(update.Message.Chat.ID, update.Message.From)
	}
}

func (r *Router) handleStart(chatID int64, from *tgbotapi.User) {
	r.Send.SendSticker(chatID, r.Cfg.WelcomeStickerID)
	r.Send.SendHTML(chatID, botapp.MainMenuText(), botapp.MainMenuKeyboard(r.Cfg.MiniAppURL))
}

func (r *Router) handleCallback(ctx context.Context, cq *tgbotapi.CallbackQuery) {
	if cq.From == nil || cq.Message == nil {
		return
	}
	data := strings.TrimSpace(cq.Data)
	chatID := cq.Message.Chat.ID
	msgID := cq.Message.MessageID
	from := cq.From

	r.Send.AnswerCallback(cq.ID, "")

	showMain := func() {
		r.Send.EditOrSend(chatID, msgID, botapp.MainMenuText(), botapp.MainMenuKeyboard(r.Cfg.MiniAppURL))
	}

	switch {
	case data == botapp.CBBackMain:
		showMain()

	case data == botapp.CBTrialMenu:
		r.Send.EditOrSend(chatID, msgID, botapp.TrialIntroText(), botapp.TrialIntroKeyboard())

	case data == botapp.CBTrialActivate:
		_, err := r.API.ActivateTrial(ctx, from.ID, from.FirstName, from.LastName, from.UserName)
		if err != nil {
			if vpnapi.IsTrialAlreadyUsed(err) {
				r.Send.SendSticker(chatID, r.Cfg.ErrorStickerID)
				r.Send.EditOrSend(chatID, msgID, botapp.TrialUsedText(), botapp.TrialUsedKeyboard())
				return
			}
			if vpnapi.IsTrialActiveSubscription(err) {
				r.Send.EditOrSend(chatID, msgID, "⚡ У вас уже есть активная подписка.\n\nПолучите конфиг в <b>🌊 Личном кабинете</b> или <b>🌐 Веб-кабинете</b>.", botapp.MainMenuKeyboard(r.Cfg.MiniAppURL))
				return
			}
			r.Send.SendSticker(chatID, r.Cfg.ErrorStickerID)
			r.Send.EditOrSend(chatID, msgID, botapp.GenericErrorText(), botapp.MainMenuKeyboard(r.Cfg.MiniAppURL))
			log.Printf("trial activate: %v", err)
			return
		}
		r.Send.SendSticker(chatID, r.Cfg.SuccessStickerID)
		r.Send.EditOrSend(chatID, msgID, botapp.TrialSuccessText(), botapp.AfterSuccessKeyboard(r.Cfg.MiniAppURL))

	case data == botapp.CBBuyVPN || data == botapp.CBChangePlan:
		r.Send.EditOrSend(chatID, msgID, botapp.BuyPlansIntroText(), botapp.PlansKeyboard())

	case data == botapp.CBPlan1m, data == botapp.CBPlan2m, data == botapp.CBPlan3m, data == botapp.CBPlan6m, data == botapp.CBPlan12m:
		code := botapp.PlanFromCallback(data)
		p, ok := r.Plans[code]
		if !ok {
			return
		}
		r.Send.EditOrSend(chatID, msgID, botapp.CheckoutText(p), botapp.CheckoutKeyboard())

	case data == botapp.CBPaymentSoon:
		r.Send.EditOrSend(chatID, msgID, botapp.PaymentSoonText(), botapp.TrialIntroKeyboard())

	case data == botapp.CBGetConfig:
		cfgResp, err := r.API.GetConfig(ctx, from.ID, from.FirstName, from.LastName, from.UserName)
		if err != nil {
			if vpnapi.IsNoSubscription(err) {
				r.Send.EditOrSend(chatID, msgID, botapp.NoSubscriptionText(), botapp.TrialIntroKeyboard())
				return
			}
			r.Send.SendSticker(chatID, r.Cfg.ErrorStickerID)
			r.Send.EditOrSend(chatID, msgID, botapp.GenericErrorText(), botapp.MainMenuKeyboard(r.Cfg.MiniAppURL))
			log.Printf("[get_config] telegram_id=%d err=%v", from.ID, err)
			return
		}
		r.Send.EditOrSend(chatID, msgID, botapp.ConfigText("🔑 <b>Ваш VPN-конфиг</b>", cfgResp.VlessURI), botapp.MainMenuKeyboard(r.Cfg.MiniAppURL))

	case data == botapp.CBRefreshConfig:
		cfgResp, err := r.API.RefreshConfig(ctx, from.ID, from.FirstName, from.LastName, from.UserName)
		if err != nil {
			if vpnapi.IsNoSubscription(err) {
				r.Send.EditOrSend(chatID, msgID, botapp.NoSubscriptionText(), botapp.TrialIntroKeyboard())
				return
			}
			r.Send.SendSticker(chatID, r.Cfg.ErrorStickerID)
			r.Send.EditOrSend(chatID, msgID, botapp.GenericErrorText(), botapp.MainMenuKeyboard(r.Cfg.MiniAppURL))
			log.Printf("[refresh_config] telegram_id=%d err=%v", from.ID, err)
			return
		}
		r.Send.EditOrSend(chatID, msgID, botapp.ConfigText("🔄 <b>Новый конфиг</b>", cfgResp.VlessURI), botapp.MainMenuKeyboard(r.Cfg.MiniAppURL))

	case data == botapp.CBInstructions:
		r.Send.EditOrSend(chatID, msgID, botapp.GuideMenuText(), botapp.InstructionsKeyboard(r.Cfg.MiniAppURL))

	case data == botapp.CBSupport:
		r.Send.EditOrSend(chatID, msgID, botapp.SupportText(r.Cfg.SupportContact, r.Cfg.SupportTelegramURL), botapp.MainMenuKeyboard(r.Cfg.MiniAppURL))

	case data == botapp.CBCabinet:
		if r.Cfg.MiniAppURL == "" {
			r.Send.EditOrSend(chatID, msgID, botapp.CabinetUnavailableText(), botapp.MainMenuKeyboard(""))
		}

	case data == botapp.CBWebCabinet:
		webURL := r.Cfg.MiniAppURL
		if webURL == "" {
			webURL = "https://app.surfwave.space"
		}
		r.Send.EditOrSend(chatID, msgID, botapp.WebCabinetText(webURL), botapp.WebCabinetKeyboard(webURL))
	}
}
