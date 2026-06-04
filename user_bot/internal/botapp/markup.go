package botapp

// JSONMarkup supports web_app buttons not present in older tgbotapi structs.
type JSONMarkup struct {
	InlineKeyboard [][]inlineBtn `json:"inline_keyboard"`
}

type inlineBtn struct {
	Text         string      `json:"text"`
	CallbackData *string     `json:"callback_data,omitempty"`
	URL          *string     `json:"url,omitempty"`
	WebApp       *webAppInfo `json:"web_app,omitempty"`
}

type webAppInfo struct {
	URL string `json:"url"`
}

func (m JSONMarkup) Markup() interface{} { return m }

func cbBtn(text, data string) inlineBtn {
	return inlineBtn{Text: text, CallbackData: &data}
}

func webAppBtn(text, url string) inlineBtn {
	return inlineBtn{Text: text, WebApp: &webAppInfo{URL: url}}
}

func MainMenuKeyboard(miniAppURL string) JSONMarkup {
	rows := [][]inlineBtn{
		{cbBtn("🎁 Бесплатно на 24 часа", CBTrialMenu)},
		{cbBtn("💳 Купить VPN", CBBuyVPN)},
		{cbBtn("🔑 Получить конфиг", CBGetConfig), cbBtn("🔄 Обновить конфиг", CBRefreshConfig)},
		{cbBtn("📱 Инструкции", CBInstructions), cbBtn("🛟 Поддержка", CBSupport)},
	}
	if miniAppURL != "" {
		rows = append(rows, []inlineBtn{webAppBtn("🌊 Личный кабинет", miniAppURL)})
	} else {
		rows = append(rows, []inlineBtn{cbBtn("🌊 Личный кабинет", CBCabinet)})
	}
	return JSONMarkup{InlineKeyboard: rows}
}

func BackMainRow() []inlineBtn {
	return []inlineBtn{cbBtn("⬅️ Назад", CBBackMain)}
}

func TrialIntroKeyboard() JSONMarkup {
	return JSONMarkup{InlineKeyboard: [][]inlineBtn{
		{cbBtn("🚀 Активировать тест", CBTrialActivate)},
		BackMainRow(),
	}}
}

func TrialUsedKeyboard() JSONMarkup {
	return JSONMarkup{InlineKeyboard: [][]inlineBtn{
		{cbBtn("💳 Купить VPN", CBBuyVPN)},
		{cbBtn("🏠 Главное меню", CBBackMain)},
	}}
}

func PlansKeyboard() JSONMarkup {
	return JSONMarkup{InlineKeyboard: [][]inlineBtn{
		{cbBtn("1 месяц", CBPlan1m)},
		{cbBtn("3 месяца", CBPlan3m), cbBtn("6 месяцев", CBPlan6m)},
		{cbBtn("12 месяцев", CBPlan12m)},
		BackMainRow(),
	}}
}

func CheckoutKeyboard() JSONMarkup {
	return JSONMarkup{InlineKeyboard: [][]inlineBtn{
		{cbBtn("🔜 Оплата скоро", CBPaymentSoon)},
		{cbBtn("🎁 Бесплатно 24 часа", CBTrialMenu)},
		{cbBtn("⬅️ Изменить тариф", CBChangePlan)},
		{cbBtn("🏠 Главное меню", CBBackMain)},
	}}
}

func InstructionsKeyboard() JSONMarkup {
	return JSONMarkup{InlineKeyboard: [][]inlineBtn{
		{cbBtn("iPhone", CBGuideIOS), cbBtn("Android", CBGuideAndroid)},
		{cbBtn("Windows", CBGuideWin), cbBtn("macOS", CBGuideMac)},
		BackMainRow(),
	}}
}

func AfterSuccessKeyboard(miniAppURL string) JSONMarkup {
	rows := [][]inlineBtn{
		{cbBtn("🔑 Получить конфиг", CBGetConfig)},
	}
	if miniAppURL != "" {
		rows = append(rows, []inlineBtn{webAppBtn("🌊 Личный кабинет", miniAppURL)})
	}
	rows = append(rows, BackMainRow())
	return JSONMarkup{InlineKeyboard: rows}
}

// ParseMockPay returns plan code from "mock_pay:1m".
func ParseMockPay(data string) (string, bool) {
	const p = CBMockPay + ":"
	if len(data) > len(p) && data[:len(p)] == p {
		return data[len(p):], true
	}
	return "", false
}

