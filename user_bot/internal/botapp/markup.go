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

func urlBtn(text, url string) inlineBtn {
	u := url
	return inlineBtn{Text: text, URL: &u}
}

func MainMenuKeyboard(miniAppURL string) JSONMarkup {
	rows := [][]inlineBtn{
		{cbBtn("🎁 Бесплатно на 24 часа", CBTrialMenu)},
		{cbBtn("💳 Купить VPN", CBBuyVPN)},
		{cbBtn("📖 Как подключиться", CBInstructions), cbBtn("🛟 Поддержка", CBSupport)},
	}
	if miniAppURL != "" {
		rows = append(rows, []inlineBtn{webAppBtn("🌊 Личный кабинет", miniAppURL)})
	} else {
		rows = append(rows, []inlineBtn{cbBtn("🌊 Личный кабинет", CBCabinet)})
	}
	rows = append(rows, []inlineBtn{cbBtn("🌐 Веб-кабинет", CBWebCabinet)})
	return JSONMarkup{InlineKeyboard: rows}
}

func WebCabinetKeyboard(webURL string) JSONMarkup {
	return JSONMarkup{InlineKeyboard: [][]inlineBtn{
		{urlBtn("🌐 Открыть веб-кабинет", webURL)},
		BackMainRow(),
	}}
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
		BackMainRow(),
	}}
}

func AfterSuccessKeyboard(miniAppURL string) JSONMarkup {
	var rows [][]inlineBtn
	if miniAppURL != "" {
		rows = append(rows, []inlineBtn{webAppBtn("🌊 Личный кабинет", miniAppURL)})
	}
	rows = append(rows, []inlineBtn{cbBtn("🌐 Веб-кабинет", CBWebCabinet)})
	rows = append(rows, BackMainRow())
	return JSONMarkup{InlineKeyboard: rows}
}

