package handlers

import (
	"strconv"
	"strings"

	"github.com/Mudicat-pr/firstTgBot/internal/storage"
	"github.com/Mudicat-pr/firstTgBot/internal/tools"
	"github.com/Mudicat-pr/firstTgBot/pkg/e"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type BaseVar struct {
	Bot      *tgbotapi.BotAPI
	AppealDB *storage.AppealHandle
	TariffDB *storage.TariffHandle
	F        *tools.FSM
}

// Состояния для заявления
const (
	AppealOpened  = "Открыта"
	AppealProcess = "В процессе"
	AppealClosed  = "Закрыта/Решена"
	AppealBan     = "Отклонена"
)

// Флаги и иные константы для работы с булевой датой
const (
	FlagTrue  = true
	FlagFalse = false

	IsOpened = "Открыть"
	IsHidden = "Скрыть"
)

// Настройки парсера мапы при сплите сообщения
const (
	InputSeparator = "; "
)

// Функция парсера как объект для валидации мапы
type Parcer[T any] func([]string) (*T, error)

// Парсер сплита текста юзера в мапу (Для тарифов)
var TariffParser = map[string]Parcer[storage.Tariff]{
	"addTariff": func(data []string) (*storage.Tariff, error) {
		if len(data) < 3 {
			return nil, e.ErrInvalidInputFormat
		}
		price, err := MsgToInt(data[2])
		if err != nil {
			return nil, err
		}
		return &storage.Tariff{
			Title: data[0],
			Body:  data[1],
			Price: price,
		}, nil
	},
	"hideTariff": func(data []string) (*storage.Tariff, error) {
		if len(data) < 2 {
			return nil, e.ErrInvalidInputFormat
		}
		id, err := MsgToInt(data[0])
		if err != nil {
			return nil, err
		}
		var isHide bool
		switch data[1] {
		case IsHidden:
			isHide = true
		case IsOpened:
			isHide = false
		default:
			return nil, e.ErrInvalidInputFormat
		}
		return &storage.Tariff{
			ID:     id,
			IsHide: isHide,
		}, nil
	},
	"editTariff": func(data []string) (*storage.Tariff, error) {
		if len(data) < 4 {
			return nil, e.ErrInvalidInputFormat
		}
		id, err := MsgToInt(data[0])
		if err != nil {
			return nil, err
		}
		price, err := MsgToInt(data[3])
		if err != nil {
			return nil, err
		}
		return &storage.Tariff{
			ID:    id,
			Title: data[1],
			Body:  data[2],
			Price: price,
		}, nil
	},
}

// Парсер сплита текста юзера в мапу (Дла заявок)
var AppealParcer = map[string]Parcer[storage.Appeal]{
	"addAppeal": func(data []string) (*storage.Appeal, error) {
		if len(data) < 5 {
			return nil, e.ErrInvalidInputFormat
		}
		return &storage.Appeal{
			TariffName: data[0],
			AppealData: storage.AppealData{
				FullName: data[1],
				Email:    data[2],
				Address:  data[3],
				Phone:    data[4],
			},
		}, nil
	},
}

// Сам сплит. Данные самостоятельно не на валидирует, только отдает результат.
// На вход идет сообщение пользователя input,
// мод (режим) из парсеров выше, вся логика описана там вместе в валидацией
// parcerType - тип парсера: тарифный или для заявлений
func SplitInput[T any](input, mode string, parserType map[string]Parcer[T]) (*T, error) {
	data := strings.Split(input, InputSeparator)
	parser, ok := parserType[mode]
	if !ok {
		return nil, e.ErrInvalidInputFormat
	}
	return parser(data)
}

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
