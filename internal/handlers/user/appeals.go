package user

import (
	"fmt"

	"github.com/Mudicat-pr/firstTgBot/pkg/e"

	h "github.com/Mudicat-pr/firstTgBot/internal/handlers"
	"github.com/Mudicat-pr/firstTgBot/pkg/idgen"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (a *UserHandle) Add(msg *tgbotapi.Message) (err error) {
	defer func() { err = e.WrapIfErr("Can't send appeal. METHOD Add from user/appeals.go", err) }()

	contract := idgen.IDgenerator(msg.From.ID)
	data, err := h.SplitInput(msg.Text, "addAppeal", h.AppealParcer)
	if err != nil {
		h.MsgForUser(*a.Bot, msg.Chat.ID, "Данные не прошли валидацию. Попробуйте снова")
		return err
	}
	err = a.AppealDB.Add(data.TariffName,
		msg.From.ID,
		contract,
		data.AppealData.FullName,
		data.AppealData.Address,
		data.AppealData.Email,
		data.AppealData.Phone)
	if err != nil {
		h.MsgForUser(*a.Bot, msg.Chat.ID, "Данные не прошли валидацию. Попробуйте снова")
		return err
	}
	textMsg := fmt.Sprintf("Ваше заявление успешно отправлено в нашу базу данных!\nВыдан статус <b>'%v'</b>", h.AppealOpened)
	h.MsgForUser(*a.Bot, msg.Chat.ID, textMsg)
	a.F.ClearState(msg.From.ID)
	return err
}

func (a *UserHandle) Details(msg *tgbotapi.Message) (err error) {
	defer func() { err = e.WrapIfErr("No found appeal for view details", err) }()
	// Сделать функцию просмотра деталей заявления
	return nil
}

func (a *UserHandle) Edit(msg *tgbotapi.Message) (err error) {
	// Изменение тарифа пользователем
	return nil
}
