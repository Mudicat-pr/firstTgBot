package handlers

import (
	"strconv"

	"github.com/Mudicat-pr/firstTgBot/internal/storage"
	"github.com/Mudicat-pr/firstTgBot/internal/tools"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type BaseVar struct {
	Bot        *tgbotapi.BotAPI
	ContractDB *storage.ContractHandle
	TariffDB   *storage.TariffHandle
	F          *tools.FSM
}

// Состояния для заявления
const (
	ContractOpened  = "Открыта"
	ContractProcess = "В процессе"
	ContractClosed  = "Закрыта/Решена"
	ContractBan     = "Отклонена и удалена"

	SkipText = "Нет данных для пропуска"
	Skip     = "Пропустить"
	SkipHint = "\n\n<i>Вы можете проспутить шаг, введя ключевое слово </i>" + Skip
)

// Флаги и иные константы для работы с булевой датой
const (
	FlagTrue  = true
	FlagFalse = false

	IsOpened = "Открыть"
	IsHidden = "Скрыть"
)

const (
	tariffType = 1 << iota
	appealType
)

// Просто создано для удобства, чтоб не писать из раза в раз
// одну и ту же сигнатуру
// С недавних пор теперь может парсить строку в HTML для стилизации
func MsgForUser(bot tgbotapi.BotAPI, userChatID int64, text string) {
	response := tgbotapi.NewMessage(userChatID, text)
	response.ParseMode = tgbotapi.ModeHTML
	bot.Send(response)
}

// Перевод строки в чисто. Возможно изменю эту функцию для вывода
// более оптимальных типов данных числа (например uint8, int32 и т.д.)
func MsgToInt(msg string) (res int, err error) {
	res, err = strconv.Atoi(msg)
	return res, err
}
