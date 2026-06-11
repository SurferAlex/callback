package notify

import "time"

// Type identifies a subscription expiry reminder (stored in subscription_notifications).
type Type string

const (
	Type7Days  Type = "7_days"
	Type3Days  Type = "3_days"
	Type1Day   Type = "1_day"
	Type2Hours Type = "2_hours"
)

// Rule defines when to send a reminder relative to subscriptions.ends_at.
// Add new entries here to support extra intervals without changing worker logic.
type Rule struct {
	Type    Type
	Before  time.Duration
	Message string
}

// Rules returns reminder rules ordered from earliest threshold to latest (largest Before first).
func Rules() []Rule {
	return []Rule{
		{
			Type:   Type7Days,
			Before: 7 * 24 * time.Hour,
			Message: "🌊 Ваша подписка SurfWave VPN закончится через 7 дней.\n\n" +
				"Продлите подписку заранее, чтобы сохранить доступ к VPN.",
		},
		{
			Type:   Type3Days,
			Before: 3 * 24 * time.Hour,
			Message: "🌊 До окончания подписки осталось 3 дня.\n\n" +
				"Не забудьте продлить доступ.",
		},
		{
			Type:   Type1Day,
			Before: 24 * time.Hour,
			Message: "⚠️ До окончания подписки остался 1 день.\n\n" +
				"Чтобы не потерять доступ к VPN, продлите подписку заранее.",
		},
		{
			Type:   Type2Hours,
			Before: 2 * time.Hour,
			Message: "🚨 Подписка закончится через 2 часа.\n\n" +
				"Продлите её сейчас, чтобы избежать отключения.",
		},
	}
}
