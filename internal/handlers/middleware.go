package handlers

import (
	"github.com/Mudicat-pr/firstTgBot/config"
	"github.com/Mudicat-pr/firstTgBot/pkg/e"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func IsAdmin(msg *tgbotapi.Message) bool {
	if msg.From == nil || msg == nil {
		return false
	}
	user := msg.From.ID

	cfg, err := config.ReadConfig()
	if err != nil {
		e.Wrap("failed to read config", err)
		return false
	}

	if user != cfg.AdminID {
		return false
	}
	return true
}
