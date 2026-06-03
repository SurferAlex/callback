package botapp

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

const (
	BtnBuy       = "💳 Купить VPN"
	BtnGetConfig = "🔑 Получить конфиг"
	BtnRefresh   = "🔄 Обновить конфиг"
	BtnGuide     = "📱 Инструкции"
	BtnSupport   = "🛟 Поддержка"
	BtnCabinet   = "🌊 Личный кабинет"

	BtnPlan1m  = "1 месяц"
	BtnPlan3m  = "3 месяца"
	BtnPlan6m  = "6 месяцев"
	BtnPlan12m = "12 месяцев"
	BtnMockPay = "✅ Тестовая активация"
	BtnBack    = "⬅️ В меню"

	BtnGuideIOS     = "iPhone"
	BtnGuideAndroid = "Android"
	BtnGuideWin     = "Windows"
	BtnGuideMac     = "macOS"
)

func MainMenuMarkup() tgbotapi.ReplyKeyboardMarkup {
	return tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(BtnBuy),
			tgbotapi.NewKeyboardButton(BtnGetConfig),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(BtnRefresh),
			tgbotapi.NewKeyboardButton(BtnGuide),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(BtnSupport),
			tgbotapi.NewKeyboardButton(BtnCabinet),
		),
	)
}

func PlansMarkup() tgbotapi.ReplyKeyboardMarkup {
	return tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(tgbotapi.NewKeyboardButton(BtnPlan1m)),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(BtnPlan3m),
			tgbotapi.NewKeyboardButton(BtnPlan6m),
		),
		tgbotapi.NewKeyboardButtonRow(tgbotapi.NewKeyboardButton(BtnPlan12m)),
		tgbotapi.NewKeyboardButtonRow(tgbotapi.NewKeyboardButton(BtnBack)),
	)
}

func MockPayMarkup() tgbotapi.ReplyKeyboardMarkup {
	return tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(tgbotapi.NewKeyboardButton(BtnMockPay)),
		tgbotapi.NewKeyboardButtonRow(tgbotapi.NewKeyboardButton(BtnBack)),
	)
}

func GuideMarkup() tgbotapi.ReplyKeyboardMarkup {
	return tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(BtnGuideIOS),
			tgbotapi.NewKeyboardButton(BtnGuideAndroid),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(BtnGuideWin),
			tgbotapi.NewKeyboardButton(BtnGuideMac),
		),
		tgbotapi.NewKeyboardButtonRow(tgbotapi.NewKeyboardButton(BtnBack)),
	)
}

func PlanCode(button string) string {
	switch button {
	case BtnPlan1m:
		return "1m"
	case BtnPlan3m:
		return "3m"
	case BtnPlan6m:
		return "6m"
	case BtnPlan12m:
		return "12m"
	default:
		return ""
	}
}

func WelcomeText(name string) string {
	if name == "" {
		name = "друг"
	}
	return "Привет, " + name + "!\n\nSurfer VPN — быстрый доступ в один тап.\nВыберите действие в меню ниже."
}
