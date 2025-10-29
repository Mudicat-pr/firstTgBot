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

const (
	AppealOpened  = "Открыта"
	AppealProcess = "В процессе"
	AppealClosed  = "Закрыта/Решена"
	AppealBan     = "Отклонена"
)

const (
	FlagTrue  = true
	FlagFalse = false
)

const (
	EmptyTariffList = "На данный момент тарифов нет, они в процессе добавления"
	NoFoundTariff   = "Тариф по вашему запросу был не найден"
	UnknownCommand  = "Неизвестная команда. Введите /help для получения списка доступных функций"
)

const (
	InputSeparator = "; "
	IsHidden       = "Скрыть"
	IsOpened       = "Открыть"
)

type Parcer[T any] func([]string) (*T, error)

var TariffParser = map[string]Parcer[storage.Tariff]{
	"addTariff": func(data []string) (*storage.Tariff, error) {
		if len(data) < 3 {
			return nil, e.ErrInvalidInputFormat
		}
		price, err := strconv.Atoi(data[2])
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
		id, err := strconv.Atoi(data[0])
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
		id, err := strconv.Atoi(data[0])
		if err != nil {
			return nil, err
		}
		price, err := strconv.Atoi(data[3])
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

func SplitInput[T any](input, mode string, parserType map[string]Parcer[T]) (*T, error) {
	data := strings.Split(input, InputSeparator)
	parser, ok := parserType[mode]
	if !ok {
		return nil, e.ErrInvalidInputFormat
	}
	return parser(data)
}

func MsgForUser(bot tgbotapi.BotAPI, userChatID int64, text string) {
	bot.Send(tgbotapi.NewMessage(
		userChatID,
		text,
	))
}
