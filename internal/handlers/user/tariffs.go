package user

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/Mudicat-pr/firstTgBot/pkg/e"

	h "github.com/Mudicat-pr/firstTgBot/internal/handlers"

	"github.com/Mudicat-pr/firstTgBot/internal/storage"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type UserHandle struct {
	*h.BaseVar
}

// Return in tg chat all tariffs
func (t *UserHandle) All(msg *tgbotapi.Message, flag bool) (err error) {
	defer func() { err = e.WrapIfErr("Failed to select tariffs from DB", err) }()
	rows, err := t.TariffDB.AllTariffs()

	if err != nil {
		h.MsgForUser(*t.Bot, msg.Chat.ID, h.EmptyTariffList)
		return err
	}
	_, err = t.Bot.Send(tgbotapi.NewMessage(
		msg.Chat.ID,
		buildRows(rows, flag),
	))
	return err
}

func buildRows(rows []storage.Tariff, flag bool) string {
	var message strings.Builder
	for _, t := range rows {
		if t.IsHide == flag {
			continue
		}
		fmt.Fprintf(
			&message,
			"ID %d: %s - %d рублей\n",
			t.ID, t.Title, t.Price)
	}
	var header string
	switch flag {
	case h.FlagTrue:
		header = "Все текущие тарифы: \n\n"
	case h.FlagFalse:
		header = "Все скрытые тарифы: \n\n"
	}
	footer := "Если желаете узнать подробнее - пишите в чат /details"
	return header + message.String() + "\n" + footer
}

// Return details by one of selected tariffs
func (t *UserHandle) Detail(msg *tgbotapi.Message) (err error) {
	defer func() { err = e.WrapIfErr("Can't select tariff for view details", err) }()

	tariffID, err := msgToInt(msg.Text)
	if err != nil {
		h.MsgForUser(*t.Bot, msg.Chat.ID, h.NoFoundTariff)
		return err
	}
	row, err := t.TariffDB.Details(tariffID)
	if row.IsHide == h.FlagTrue {
		t.F.ClearState(msg.From.ID)
		h.MsgForUser(*t.Bot, msg.Chat.ID, h.NoFoundTariff)
		return err
	}
	header := fmt.Sprintf("Вы выбрали тариф %d - %s: \n", row.ID, row.Title)
	body := fmt.Sprintf("Описание: %s\nЦена: %d рублей", row.Body, row.Price)

	h.MsgForUser(*t.Bot, msg.Chat.ID, header+body)
	defer t.F.ClearState(msg.From.ID)
	return err
}

func msgToInt(msg string) (res int, err error) {
	res, err = strconv.Atoi(msg)
	return res, err
}
