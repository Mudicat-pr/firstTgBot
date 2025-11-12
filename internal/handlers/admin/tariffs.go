package admin

import (
	"errors"
	"fmt"
	"strconv"

	h "github.com/Mudicat-pr/firstTgBot/internal/handlers"
	"github.com/Mudicat-pr/firstTgBot/internal/storage"
	"github.com/Mudicat-pr/firstTgBot/internal/tools"
	"github.com/Mudicat-pr/firstTgBot/pkg/e"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type AdminHandle struct {
	*h.BaseVar
}

type StepConfig struct {
	NextState int
	Prompt    string
	Field     int
}

const (
	fieldTitle = iota
	fieldBody
	fieldPrice

	editMode
	addMode
)

func (a *AdminHandle) handle(msg *tgbotapi.Message, data interface{}, nextState, field int) (state int, newData interface{}, err error) {
	ErrMessage := fmt.Sprintf("Failed to handle with msg: %v and data %v", msg.Text, data)
	defer func() { err = e.WrapIfErr(ErrMessage, err) }()
	trf := setStruct(data)
	userID := msg.From.ID
	currentState := a.F.UserState(userID)

	if msg.Text == h.Skip {
		actualData, err := a.actualField(trf, msg, field)
		if err != nil {
			h.MsgForUser(*a.Bot, userID, "Невозможно пропустить шаг")
			return currentState, trf, err
		}
		return nextState, actualData, nil
	}

	newData, err = a.validateField(trf, msg, field)
	if err != nil {
		h.MsgForUser(*a.Bot, userID, "Невозможно добавить данные")
		return currentState, trf, err
	}

	return nextState, newData, nil
}

func (a *AdminHandle) actualField(data *storage.Tariff, msg *tgbotapi.Message, field int) (*storage.Tariff, error) {
	actualData, err := a.TariffDB.Details(data.ID)
	if err != nil {
		h.MsgForUser(*a.Bot, msg.From.ID, h.SkipText)
		return data, e.Wrap("Failed select actual data", err)
	}
	switch field {
	case fieldTitle:
		data.Title = actualData.Title
	case fieldBody:
		data.Body = actualData.Body
	case fieldPrice:
		data.Price = actualData.Price
	}
	return data, nil
}

func (a *AdminHandle) validateField(data *storage.Tariff, msg *tgbotapi.Message, field int) (*storage.Tariff, error) {
	newData := msg.Text
	switch field {
	case fieldTitle:
		if len(newData) < 3 {
			h.MsgForUser(*a.Bot, msg.From.ID, "Слишком короткое имя")
			return data, errors.New("Too short name. Failed validate")
		}
		data.Title = newData
	case fieldBody:
		if len(newData) > 500 {
			h.MsgForUser(*a.Bot, msg.From.ID, "Слишком длинное описание (более 500 символов)")
			return data, errors.New("Too long description, Failed validate")
		}
		data.Body = newData
	case fieldPrice:
		price, err := strconv.Atoi(newData)
		if err != nil {
			h.MsgForUser(*a.Bot, msg.From.ID, "Неверный ввод. Должно быть только целое число")
			return data, err
		}
		data.Price = price
	}
	return data, nil
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

func (a *AdminHandle) chainHelper(msg *tgbotapi.Message, data interface{}, steps map[int]StepConfig, mode int) (state int, newData interface{}, err error) {
	userID := msg.From.ID
	step, exists := steps[a.F.UserState(userID)]
	if !exists {
		return 0, nil, errors.New("State is not defined")
	}
	h.MsgForUser(*a.Bot, userID, step.Prompt)
	state, newData, err = a.handle(msg, data, step.NextState, step.Field)

	if step.NextState == 0 && err == nil {
		trf := setStruct(data)
		switch mode {
		case addMode:
			if err = a.TariffDB.Add(trf.Title, trf.Body, trf.Price); err != nil {
				h.MsgForUser(*a.Bot, userID, "Ошибка добавления")
				return 0, nil, err
			}
			h.MsgForUser(*a.Bot, userID, "Тариф добавлен")
		case editMode:
			if err = a.TariffDB.Edit(trf.ID, trf.Title, trf.Body, trf.Price); err != nil {
				h.MsgForUser(*a.Bot, userID, "Не удалось изменить тариф")
				return 0, nil, err
			}
			msgBot := fmt.Sprintf("Раньше тариф выглядил так:\nИМЯ - %s\nОПИСАНИЕ - %s\nЦЕНА - %d", trf.Title, trf.Body, trf.Price)
			h.MsgForUser(*a.Bot, userID, "Тариф успешно изменен\n"+msgBot)
		}
	}
	return state, newData, err
}

func (a *AdminHandle) AddChain(msg *tgbotapi.Message, data interface{}) (state int, newData interface{}, err error) {
	defer func() { err = e.WrapIfErr("Failed handle add chain in admin manage", err) }()
	steps := map[int]StepConfig{
		tools.TariffTitle: {tools.TariffBody, "Описание тарифа. Не более 500 символов", fieldTitle},
		tools.TariffBody:  {tools.TariffPrice, "Укажите ценник целым числом", fieldBody},
		tools.TariffPrice: {0, "", fieldPrice},
	}
	return a.chainHelper(msg, data, steps, addMode)
}

func (a *AdminHandle) EditChain(msg *tgbotapi.Message, data interface{}) (state int, newData interface{}, err error) {
	defer func() { err = e.WrapIfErr("Failed in edit handle", err) }()
	steps := map[int]StepConfig{
		tools.TariffTitleEdit: {tools.TariffBodyEdit, "Введите новое описание тарифа, не более 500 символов", fieldTitle},
		tools.TariffBodyEdit:  {tools.TariffPriceEdit, "Введите новую стоимость тарифного плана. Целым числом", fieldBody},
		tools.TariffPriceEdit: {0, "", fieldPrice},
	}
	return a.chainHelper(msg, data, steps, editMode)
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
