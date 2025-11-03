package user

import (
	"fmt"
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
		h.MsgForUser(*t.Bot, msg.Chat.ID, "На данный момент тарифов нет или они находятся в разработке")
		return err
	}
	h.MsgForUser(*t.Bot, msg.Chat.ID, buildRows(rows, flag))
	return nil
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
		header = "<b>Все текущие тарифы:</b>\n\n"
	case h.FlagFalse:
		header = "<b>Все скрытые тарифы:</b> \n\n"
	}
	footer := "<i>Если желаете узнать подробнее - пишите в чат /details</i>"
	return header + message.String() + "\n" + footer
}

// Return details by one of selected tariffs
func (t *UserHandle) Detail(msg *tgbotapi.Message) (err error) {
	defer func() { err = e.WrapIfErr("Can't select tariff for view details", err) }()

	tariffID, err := h.MsgToInt(msg.Text)
	if err != nil {
		h.MsgForUser(*t.Bot, msg.Chat.ID, "Тарифа с таким ID не существует или он не доступен в данный момент")
		return err
	}
	row, err := t.TariffDB.Details(tariffID)
	if row.IsHide == h.FlagTrue {
		t.F.ClearState(msg.From.ID)
		h.MsgForUser(*t.Bot, msg.Chat.ID, "Тарифа с таким ID не существует или он не доступен в данный момент")
		return err
	}
	header := fmt.Sprintf("Вы выбрали тариф %d - %s: \n", row.ID, row.Title)
	body := fmt.Sprintf("Описание: %s\nЦена: %d рублей", row.Body, row.Price)

	h.MsgForUser(*t.Bot, msg.Chat.ID, header+body)
	defer t.F.ClearState(msg.From.ID)
	return err
}
