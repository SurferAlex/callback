package bot

import (
	"log"

	"user-bot/internal/botapp"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Sender struct {
	Bot *tgbotapi.BotAPI
}

func (s *Sender) SendSticker(chatID int64, stickerID string) {
	if stickerID == "" {
		return
	}
	msg := tgbotapi.NewSticker(chatID, tgbotapi.FileID(stickerID))
	if _, err := s.Bot.Send(msg); err != nil {
		log.Printf("sticker send failed chat=%d: %v", chatID, err)
	}
}

func (s *Sender) SendHTML(chatID int64, text string, markup interface{}) {
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = tgbotapi.ModeHTML
	if markup != nil {
		msg.ReplyMarkup = markup
	}
	if _, err := s.Bot.Send(msg); err != nil {
		log.Printf("send html failed chat=%d: %v", chatID, err)
	}
}

func (s *Sender) AnswerCallback(id, text string) {
	cb := tgbotapi.NewCallback(id, text)
	if _, err := s.Bot.Request(cb); err != nil {
		log.Printf("answer callback: %v", err)
	}
}

func (s *Sender) EditOrSend(chatID int64, messageID int, text string, markup botapp.JSONMarkup) {
	del := tgbotapi.NewDeleteMessage(chatID, messageID)
	if _, err := s.Bot.Request(del); err != nil {
		log.Printf("delete message chat=%d msg=%d: %v", chatID, messageID, err)
	}
	s.SendHTML(chatID, text, markup)
}
