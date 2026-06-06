package botapp

import (
	"fmt"
	"strings"
)

func MainMenuText(name string) string {
	if name == "" {
		name = "друг"
	}
	return fmt.Sprintf(`🏄 <b>Surf VPN</b>

Привет, %s!

Выберите действие:`, escapeHTML(name))
}

func TrialIntroText() string {
	return `🏄 <b>Добро пожаловать в Surf VPN</b>

Попробуйте сервис бесплатно в течение <b>24 часов</b>.

✅ Полный доступ
✅ Все функции VPN
✅ До 2 устройств

После активации тест нельзя получить повторно.`
}

func TrialUsedText() string {
	return `⚠️ <b>Бесплатный тест уже был активирован</b>

Платная подписка скоро будет доступна в боте.`
}

func TrialSuccessText() string {
	return `🎉 <b>Бесплатный доступ активирован</b>

⏳ Срок действия: <b>24 часа</b>
📱 Подключений: до <b>2 устройств</b>
⚡ Скорость: без ограничений

Ваш VPN-конфиг готов к использованию.`
}

func MockSuccessText(plan PlanOffer, expires string, daysLeft int) string {
	return fmt.Sprintf(`🎉 <b>Подписка активирована</b>

💎 Тариф: %s
⏳ Действует до: %s
📱 Устройств: %d

Осталось дней: <b>%d</b>`, plan.PlanName, expires, plan.Devices, daysLeft)
}

func BuyPlansIntroText() string {
	return `💳 <b>Premium подписка</b>

Выберите срок. Оплата подключается — пока доступен только бесплатный тест на 24 часа.`
}

func CheckoutText(p PlanOffer) string {
	return fmt.Sprintf(`━━━━━━━━━━━━━━

🏄 <b>Surf VPN Premium</b>

📱 Устройств: <b>%d</b>
⏳ Срок: <b>%s</b>
💎 Тариф: <b>%s</b>
💰 К оплате: <b>%d ₽</b>

🔜 Оплата временно недоступна

VPN-ключ сейчас выдаётся только в разделе
<b>🎁 Бесплатно на 24 часа</b>.

━━━━━━━━━━━━━━`, p.Devices, p.Duration, p.PlanName, p.PriceRub)
}

func PaymentSoonText() string {
	return `🔜 <b>Оплата скоро</b>

Мы подключаем платёжную систему. Пока воспользуйтесь бесплатным доступом на 24 часа — с полным VPN-конфигом.`
}

func ConfigText(prefix, vless string) string {
	return prefix + "\n\n<code>" + escapeHTML(vless) + "</code>"
}

func SupportText(contact, url string) string {
	return fmt.Sprintf(`🛟 <b>Поддержка Surf VPN</b>

Напишите нам — поможем с подключением.

Контакт: %s
Telegram: %s`, escapeHTML(contact), escapeHTML(url))
}

func CabinetUnavailableText() string {
	return `⚠️ <b>Личный кабинет временно недоступен</b>

Попробуйте позже или воспользуйтесь меню бота.`
}

func NoSubscriptionText() string {
	return `⚠️ <b>VPN-конфиг недоступен</b>

Сейчас ключ выдаётся только после активации
<b>🎁 Бесплатно на 24 часа</b> (один раз).`
}

func GenericErrorText() string {
	return `⚠️ <b>Что-то пошло не так</b>

Попробуйте ещё раз через минуту.`
}

func WebCabinetText(webURL string) string {
	if webURL == "" {
		webURL = "https://app.surfwave.space"
	}
	return fmt.Sprintf(`🌐 <b>Веб-кабинет Surf VPN</b>

Если у вас нет доступа к Telegram или необходимо получить новый VPN-ключ, используйте веб-версию кабинета.

Возможности:

• Получение конфигов
• Обновление VPN-ключей
• Управление подпиской
• Просмотр данных аккаунта

Ссылка:
%s`, escapeHTML(webURL))
}

func GuideMenuText() string {
	return `📱 <b>Инструкция по использованию Surf VPN</b>

Для подключения VPN и управления подпиской используйте:

🌊 <b>Личный кабинет</b> — внутри Telegram

🌐 <b>Веб-кабинет</b> — через браузер

В кабинетах доступны:

• Получение конфигов
• Обновление VPN-ключей
• Информация о подписке
• Управление VPN

⚠️ Рекомендуем один раз открыть веб-кабинет и авторизоваться через Telegram.

Если доступ к Telegram будет недоступен, вы сможете войти в веб-кабинет и самостоятельно получить новый конфиг или обновить существующий ключ.`
}

func GuideIOSText() string {
	return `📱 <b>iPhone</b>

1. Установите Happ из App Store
2. Получите конфиг в боте
3. Импортируйте ключ в Happ`
}

func GuideAndroidText() string {
	return `📱 <b>Android</b>

1. Установите Happ
2. Нажмите «Получить конфиг»
3. Добавьте ключ в приложение`
}

func GuideWindowsText() string {
	return `🖥 <b>Windows</b>

1. Установите VPN-клиент с поддержкой VLESS
2. Скопируйте конфиг из бота
3. Вставьте в клиент`
}

func GuideMacText() string {
	return `🍎 <b>macOS</b>

1. Установите Happ
2. Получите конфиг в боте
3. Подключитесь в один тап`
}

func escapeHTML(s string) string {
	s = strings.ReplaceAll(s, "&", "&amp;")
	s = strings.ReplaceAll(s, "<", "&lt;")
	s = strings.ReplaceAll(s, ">", "&gt;")
	return s
}
