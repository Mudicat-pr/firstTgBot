package admin

import (
	"fmt"
	"strconv"
	"strings"

	h "github.com/Mudicat-pr/firstTgBot/internal/handlers"
	"github.com/Mudicat-pr/firstTgBot/internal/storage"
	"github.com/Mudicat-pr/firstTgBot/internal/tools"
	"github.com/Mudicat-pr/firstTgBot/pkg/e"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type AdminHandle struct {
	*h.BaseVar
}

const (
	fieldTitle = iota
	fieldBody
	fieldPrice
)

func (a *AdminHandle) handle(msg *tgbotapi.Message, data interface{}, nextState int, field int) (state int, newData interface{}, err error) {
	ErrMessage := fmt.Sprintf("Failed to handle with msg: %v and data %v", msg.Text, data)
	defer func() { err = e.WrapIfErr(ErrMessage, err) }()
	trf := setStruct(data)
	text := msg.Text
	userID := msg.From.ID
	currentState := a.F.UserState(userID)

	switch field {
	case fieldBody:
		if len(text) > 500 {
			h.MsgForUser(*a.Bot, userID, "Текст слишком большой! Нужно менее 500 символов")
			return currentState, trf, nil
		}
		if text == h.Skip {
			actualData, err := a.TariffDB.Details(trf.ID)
			if err != nil {
				return currentState, trf, err
			}
			trf.Body = actualData.Body
		} else {
			trf.Body = text
		}
	case fieldTitle:
		if len(text) < 3 {
			h.MsgForUser(*a.Bot, userID, "Имя слишком короткое")
			return currentState, trf, nil
		}
		if text == h.Skip {
			actualData, err := a.TariffDB.Details(trf.ID)
			if err != nil {
				return currentState, trf, err
			}
			trf.Title = actualData.Title
		} else {
			trf.Title = text
		}
	case fieldPrice:
		if text == h.Skip {
			actualData, err := a.TariffDB.Details(trf.ID)
			if err != nil {
				return currentState, trf, err
			}
			trf.Price = actualData.Price
		} else {
			price, err := strconv.Atoi(text)
			if err != nil {
				h.MsgForUser(*a.Bot, userID, "Неверный ввод. Должно быть только целое число")
				return currentState, trf, err
			}
			trf.Price = price
		}
	}

	if nextState == 0 {
		h.MsgForUser(*a.Bot, userID, "Успех!")
	}
	return nextState, trf, nil
}

func setStruct(data interface{}) *storage.Tariff {
	var trf *storage.Tariff
	if data == nil {
		trf = &storage.Tariff{}
	} else {
		trf = data.(*storage.Tariff)
	}
	return trf
}

func (a *AdminHandle) SetTitle(msg *tgbotapi.Message, data interface{}) (state int, newData interface{}, err error) {
	h.MsgForUser(*a.Bot, msg.From.ID, "Введите описание тарифа, но не более 500 символов")
	return a.handle(msg, data, tools.TariffBody, fieldTitle)
}

func (a *AdminHandle) SetBody(msg *tgbotapi.Message, data interface{}) (state int, newData interface{}, err error) {
	h.MsgForUser(*a.Bot, msg.From.ID, "Введите ценник для нового тарифа, только целое число")
	return a.handle(msg, data, tools.TariffPrice, fieldBody)
}

func (a *AdminHandle) SetPrice(msg *tgbotapi.Message, data interface{}) (state int, newData interface{}, err error) {
	return a.handle(msg, data, 0, fieldPrice)
}

func (a *AdminHandle) Del(msg *tgbotapi.Message, data interface{}) (state int, newData interface{}, err error) {
	defer func() { err = e.WrapIfErr("Can't delete this tariff from DB", err) }()
	currentState := a.F.UserState(msg.From.ID)
	tariffID, err := strconv.Atoi(msg.Text)

	if err != nil {
		h.MsgForUser(*a.Bot, msg.From.ID, "Введено не целое число")
		return currentState, nil, err
	}
	if err = a.TariffDB.Del(tariffID); err != nil {
		h.MsgForUser(*a.Bot, msg.From.ID, "Произошла ошибка! Попробуйте снова или отмените действие")
		return currentState, nil, err
	}
	return 0, nil, err

}

func (a *AdminHandle) getTariffID(msg *tgbotapi.Message) (*storage.Tariff, error) {
	tariffID, err := strconv.Atoi(msg.Text)
	if err != nil {
		h.MsgForUser(*a.Bot, msg.From.ID, "Введен неверный ID")
		return nil, err
	}
	trf, err := a.TariffDB.Details(tariffID)
	if err != nil || trf.ID == 0 {
		h.MsgForUser(*a.Bot, msg.From.ID, "Введен неверный ID")
		return nil, err
	}

	return trf, err
}

func (a *AdminHandle) StartEdit(msg *tgbotapi.Message, data interface{}) (state int, newData interface{}, err error) {
	currentState := a.F.UserState(msg.From.ID)
	trf, err := a.getTariffID(msg)
	if err != nil {
		h.MsgForUser(*a.Bot, msg.From.ID, "Неверный ID")
		return currentState, trf, err
	}

	h.MsgForUser(*a.Bot, msg.From.ID, "Введите новое имя тарифа")
	return tools.TariffTitleEdit, trf, nil
}

func (a *AdminHandle) EditTitle(msg *tgbotapi.Message, data interface{}) (state int, newData interface{}, err error) {
	h.MsgForUser(*a.Bot, msg.From.ID, "Введите новое описание тарифа, не более 500 символов")
	return a.handle(msg, data, tools.TariffBodyEdit, fieldTitle)
}

func (a *AdminHandle) EditBody(msg *tgbotapi.Message, data interface{}) (state int, newData interface{}, err error) {
	h.MsgForUser(*a.Bot, msg.From.ID, "Введите новый ценник. Только целое числоа")
	return a.handle(msg, data, tools.TariffPriceEdit, fieldBody)
}

func (a *AdminHandle) EditPrice(msg *tgbotapi.Message, data interface{}) (state int, newData interface{}, err error) {
	trf := setStruct(data)
	t, _ := a.TariffDB.Details(trf.ID)
	msgBot := fmt.Sprintf("Раньше тариф выглядил так: ИМЯ - %s\nОПИСАНИЕ - %s\nЦЕНА - %d.\n\nУверены? [ДА/НЕТ]", t.Title, t.Body, t.Price)
	h.MsgForUser(*a.Bot, msg.From.ID, msgBot)
	return a.handle(msg, data, tools.TariffEditConfirm, fieldPrice)
}

func (a *AdminHandle) EditConfirm(msg *tgbotapi.Message, data interface{}) (state int, newData interface{}, err error) {
	defer func() { err = e.WrapIfErr("Can't edit tariff", err) }()
	if strings.ToUpper(msg.Text) == "НЕТ" {
		h.MsgForUser(*a.Bot, msg.From.ID, "Изменение отменено")
		a.F.ClearState(msg.From.ID)
		return 0, nil, nil
	}
	trf := setStruct(data)

	err = a.TariffDB.Edit(trf.ID, trf.Title, trf.Body, trf.Price)
	if err != nil {
		h.MsgForUser(*a.Bot, msg.From.ID, "Произошла ошибка! Попробуйте снова")
		return tools.TariffEdit, nil, nil
	}
	h.MsgForUser(*a.Bot, msg.From.ID, "Тариф успешно изменен!")
	return 0, nil, nil
}

func (a *AdminHandle) Hide(msg *tgbotapi.Message, data interface{}) (state int, newData interface{}, err error) {
	trf, err := a.getTariffID(msg)
	if err != nil {
		h.MsgForUser(*a.Bot, msg.From.ID, "Неверный ID тарифа")
		return a.F.UserState(msg.From.ID), trf, err
	}
	var flag string
	if trf.IsHide == h.FlagFalse {
		flag = h.IsHidden
		trf.IsHide = h.FlagTrue
	} else {
		flag = h.IsOpened
		trf.IsHide = h.FlagFalse
	}
	a.TariffDB.Hide(trf.ID, trf.IsHide)
	msgBot := fmt.Sprintf("Успех! Тариф сменил статус на: %v", flag)
	h.MsgForUser(*a.Bot, msg.From.ID, msgBot)
	return 0, nil, nil
}
