package notify

import (
	"context"
	"log"
	"time"
)

// Messenger sends a plain-text Telegram message to a user (chat_id = telegram user id).
type Messenger interface {
	SendNotification(chatID int64, text string) error
}

type Worker struct {
	store  *Store
	sender Messenger
	rules  []Rule
	now    func() time.Time
}

func NewWorker(store *Store, sender Messenger) *Worker {
	return &Worker{
		store:  store,
		sender: sender,
		rules:  Rules(),
		now:    time.Now,
	}
}

// RunOnce scans active subscriptions and sends due reminders.
func (w *Worker) RunOnce(ctx context.Context) {
	now := w.now().UTC()
	subs, err := w.store.ListActiveSubscriptions(ctx, now)
	if err != nil {
		log.Printf("[notify] list subscriptions: %v", err)
		return
	}

	checked := len(subs)
	var sent, errorsN int

	for _, sub := range subs {
		remaining := sub.EndsAt.Sub(now)
		if remaining <= 0 {
			continue
		}
		for _, rule := range w.rules {
			if remaining > rule.Before {
				continue
			}
			already, err := w.store.WasSent(ctx, sub.UserID, rule.Type)
			if err != nil {
				log.Printf("[notify] was_sent user=%d type=%s: %v", sub.UserID, rule.Type, err)
				errorsN++
				continue
			}
			if already {
				continue
			}
			if err := w.sender.SendNotification(sub.UserID, rule.Message); err != nil {
				log.Printf("[notify] telegram send user=%d type=%s: %v", sub.UserID, rule.Type, err)
				errorsN++
				continue
			}
			if err := w.store.MarkSent(ctx, sub.UserID, rule.Type, now); err != nil {
				log.Printf("[notify] mark_sent user=%d type=%s: %v", sub.UserID, rule.Type, err)
				errorsN++
				continue
			}
			sent++
			log.Printf("[notify] sent user=%d type=%s ends_at=%s", sub.UserID, rule.Type, sub.EndsAt.UTC().Format(time.RFC3339))
		}
	}

	log.Printf("[notify] tick complete: checked=%d sent=%d errors=%d", checked, sent, errorsN)
}

// StartScheduler runs RunOnce immediately, then every interval until ctx is cancelled.
func StartScheduler(ctx context.Context, w *Worker, interval time.Duration) {
	if interval <= 0 {
		interval = 30 * time.Minute
	}
	log.Printf("[notify] scheduler started interval=%s", interval)

	w.RunOnce(ctx)

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Printf("[notify] scheduler stopped")
			return
		case <-ticker.C:
			w.RunOnce(ctx)
		}
	}
}
