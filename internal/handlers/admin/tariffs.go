package admin

import (
	"strconv"

	h "github.com/Mudicat-pr/firstTgBot/internal/handlers"
	"github.com/Mudicat-pr/firstTgBot/internal/storage"
	"github.com/Mudicat-pr/firstTgBot/pkg/e"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const (
	InvalidInput = "Вы ввели данные для добавления нового тарифа неверно. Нужно Имя; Описание; Цена"
)

type AdminHandle struct {
	*h.BaseVar
}

func (a *AdminHandle) Add(msg *tgbotapi.Message) (err error) {
	defer func() {
		err = e.WrapIfErr("Failed to insert data for table tariffs from DB. METHOD: Add from admin/tariffs.go", err)
	}()
	data, err := h.SplitInput(msg.Text, "addTariff", h.TariffParser)
	userID := msg.From.ID
	if err != nil {
		h.MsgForUser(*a.Bot, msg.Chat.ID, InvalidInput)
		return err
	}
	err = a.createAppeal(data, userID)
	if err != nil {
		h.MsgForUser(*a.Bot, msg.Chat.ID, InvalidInput)
		return err
	}
	h.MsgForUser(*a.Bot, msg.Chat.ID, "Тариф успешно добавлен!")
	return nil
}

func (a AdminHandle) createAppeal(data *storage.Tariff, userID int64) error {
	err := a.TariffDB.Add(data.Title, data.Body, data.Price)
	a.F.ClearState(userID)
	return err
}

func (t *AdminHandle) Del(msg *tgbotapi.Message) (err error) {
	defer func() { err = e.WrapIfErr("Can't delete data. METHOD: DEL from admin/tariffs.go", err) }()
	tariff, err := t.getTariff(msg.Text)
	if err != nil {
		h.MsgForUser(*t.Bot, msg.From.ID, "Невозможно удалить тариф. Возможно введен неверный тариф или такого ID не существует")
		return err
	}

	if err = t.TariffDB.Del(tariff.ID); err != nil {
		h.MsgForUser(*t.Bot, msg.From.ID, "Невозможно удалить тариф. Возможно введен неверный тариф или такого ID не существует")
		return err
	}
	h.MsgForUser(*t.Bot, msg.From.ID, "Тариф успешно удален из БД. Вернуть его будет невозможно")
	return nil
}

func (t *AdminHandle) getTariff(msg string) (data storage.Tariff, err error) {
	id, _ := strconv.Atoi(msg)
	data, err = t.TariffDB.Details(id)
	if data.ID == 0 || err != nil {
		return data, err
	}
	return data, nil
}

func (t *AdminHandle) HideByID(msg *tgbotapi.Message) (err error) {
	defer func() { err = e.WrapIfErr("Can't delete data. METHOD: HideByID from admin/tariffs.go", err) }()
	data, err := h.SplitInput(msg.Text, "hideTariff", h.TariffParser)
	if err != nil {
		h.MsgForUser(*t.Bot, msg.Chat.ID, "Неверные данные или ввод")
		return err
	}
	err = t.TariffDB.Hide(data.ID, data.IsHide)
	if err != nil {
		h.MsgForUser(*t.Bot, msg.Chat.ID, "Неверные данные или ввод")
		return err
	}
	h.MsgForUser(*t.Bot, msg.Chat.ID, "Успех!")
	t.F.ClearState(msg.From.ID)
	return err
}

func (t *AdminHandle) Edit(msg *tgbotapi.Message) (err error) {
	defer func() { err = e.WrapIfErr("Failed to edit data. METHOD Edit from admin/tariffs.go", err) }()
	data, err := h.SplitInput(msg.Text, "editTariff", h.TariffParser)
	if err != nil {
		h.MsgForUser(*t.Bot, msg.Chat.ID, "Неверные данные или порядок ввода")
		return err
	}
	err = t.TariffDB.Edit(data.ID, data.Title, data.Body, data.Price)
	if err != nil {
		h.MsgForUser(*t.Bot, msg.Chat.ID, "Неверные данные или порядок ввода")
		return err
	}
	h.MsgForUser(*t.Bot, msg.Chat.ID, "Успех!")
	t.F.ClearState(msg.From.ID)
	return err
}
