package handlers

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

func CreateSkipKey() tgbotapi.ReplyKeyboardMarkup {
	return tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(Skip),
			tgbotapi.NewKeyboardButton(Cancel),
		),
	)
}

func CancelKey() tgbotapi.ReplyKeyboardMarkup {
	return tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(Cancel),
		),
	)
}

func SwitchStatusKey() tgbotapi.ReplyKeyboardMarkup {
	return tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(ContractOpened),
			tgbotapi.NewKeyboardButton(ContractProcess),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(ContractProcess),
			tgbotapi.NewKeyboardButton(ContractBan),
		),
	)
}
