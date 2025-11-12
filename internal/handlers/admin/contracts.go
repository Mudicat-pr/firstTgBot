package admin

import (
	"fmt"
	"log"
	"strconv"

	h "github.com/Mudicat-pr/firstTgBot/internal/handlers"
	"github.com/Mudicat-pr/firstTgBot/internal/storage"
	"github.com/Mudicat-pr/firstTgBot/internal/tools"
	"github.com/Mudicat-pr/firstTgBot/pkg/e"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (a *AdminHandle) DelContract(msg *tgbotapi.Message, data interface{}) (state int, newData interface{}, err error) {
	defer func() { err = e.WrapIfErr("Can't delete contract", err) }()

	contractID, err := strconv.Atoi(msg.Text)
	if err != nil {
		h.MsgForUser(*a.Bot, msg.From.ID, "Номер контракта для удаления не найден")
		return 0, nil, err
	}
	if err = a.ContractDB.Del(contractID); err != nil {
		h.MsgForUser(*a.Bot, msg.From.ID, "Произошла непредвиденная ошибка при попытке удалить контракт")
		return 0, nil, err
	}

	h.MsgForUser(*a.Bot, msg.From.ID, "Данные о заявке/контракте безвозвратно удалены")
	return 0, nil, nil

}

func (a *AdminHandle) switchHelper(msg *tgbotapi.Message, data *storage.Contract, status string) (state int, newData interface{}, err error) {
	defer func() { err = e.WrapIfErr("Failed to switching status of contract", err) }()
	currentState := a.F.UserState(msg.From.ID)
	var newStatus string
	switch status {
	case h.ContractOpened:
		newStatus = h.ContractOpened
	case h.ContractClosed:
		newStatus = h.ContractClosed
	case h.ContractProcess:
		newStatus = h.ContractProcess
	case h.ContractBan:
		newStatus = h.ContractBan
	default:
		h.MsgForUser(*a.Bot, msg.From.ID, "Неизвестное значение для статуса")
		return currentState, data, nil
	}
	if err = a.ContractDB.SwitchStatus(data.ContractID, newStatus); err != nil {
		h.MsgForUser(*a.Bot, msg.From.ID, "Не удалось изменить статус заявки/договора")
		return 0, nil, err
	}
	notificationText := fmt.Sprintf("Ваше заявление изменило статус на: %s", newStatus)
	botCopy := *a.Bot
	go func(bot tgbotapi.BotAPI, userID int64, text string) {
		msg := tgbotapi.NewMessage(userID, text)
		if _, err := bot.Send(msg); err != nil {
			log.Printf("Произошла ошибка при попытке отправить уведомление юзеру: %v", userID)
		}
	}(botCopy, data.UserID, notificationText)

	h.MsgForUser(*a.Bot, msg.From.ID, "Статус успешно изменен!")
	return 0, nil, nil
}

func setStructContract(data interface{}) *storage.Contract {
	var ap *storage.Contract
	if data == nil {
		ap = &storage.Contract{}
	} else {
		ap = data.(*storage.Contract)
	}
	return ap
}

func (a *AdminHandle) ContractStatus(msg *tgbotapi.Message, data interface{}) (state int, newData interface{}, err error) {
	currentState := a.F.UserState(msg.From.ID)
	contract, err := strconv.Atoi(msg.Text)
	if err != nil {
		h.MsgForUser(*a.Bot, msg.From.ID, "Договора/заявки с текущим номером не найдено")
		return currentState, data, err
	}
	ap := setStructContract(data)
	ap.ContractID = contract
	ap.UserID, err = a.ContractDB.UserIDContract(contract)
	if err != nil {
		h.MsgForUser(*a.Bot, msg.From.ID, "Невозможно получить пользователя для работы со статусом")
		return 0, nil, err
	}

	selectStatus := fmt.Sprintf("Выберите один из возможных статусов: \n%v\r\n%v\r\n%v\r\n%v\r",
		h.ContractProcess,
		h.ContractOpened,
		h.ContractClosed,
		h.ContractBan)
	h.MsgForUser(*a.Bot, msg.From.ID, selectStatus)
	return tools.SwitchEnd, ap, nil
}

func (a *AdminHandle) SwitchStatus(msg *tgbotapi.Message, data interface{}) (state int, newData interface{}, err error) {
	ap := setStructContract(data)
	status := msg.Text
	return a.switchHelper(msg, ap, status)
}
