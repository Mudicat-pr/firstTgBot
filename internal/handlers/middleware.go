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

/*func IsHasAppeal(msg tgbotapi.Message, s *storage.AppealHandle) bool {
	userID := msg.From.ID
	flag, err := s.Details(userID)
	if err != nil {
		fmt.Printf("Can't insert data: %v", err)
		return flag
	}
	return flag
}*/

/*
func CheckTariff(tariffID int, s *storage.TariffHandle) (bool, error) {
	trf, err := s.Details(tariffID)
	if err != nil {
		return false, e.Wrap("Failed check data", err)
	}
	return trf.ID != 0, nil
} */
