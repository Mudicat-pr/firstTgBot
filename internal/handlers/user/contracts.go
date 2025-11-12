package user

import (
	"errors"
	"fmt"
	"log"
	"strconv"

	h "github.com/Mudicat-pr/firstTgBot/internal/handlers"
	"github.com/Mudicat-pr/firstTgBot/internal/storage"
	"github.com/Mudicat-pr/firstTgBot/internal/tools"
	"github.com/Mudicat-pr/firstTgBot/pkg/e"
	"github.com/Mudicat-pr/firstTgBot/pkg/idgen"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const (
	fieldTariffName = 1 << iota
	fieldFullname
	fieldAddress
	fieldEmail
	fieldPhone

	editMode
	addMode
)

type StepConfig struct {
	NextState int
	Prompt    string
	Field     int
}

func (u *UserHandle) handle(msg *tgbotapi.Message, data interface{}, nextState, field int) (state int, newData interface{}, err error) {
	ErrMessage := fmt.Sprintf("Failed to handle with msg: %v and data %v", msg.Text, data)
	defer func() { err = e.WrapIfErr(ErrMessage, err) }()
	trf := setStruct(data)
	userID := msg.From.ID
	currentState := u.F.UserState(userID)

	if msg.Text == h.Skip {
		actualData, err := u.actualField(trf, msg, field)
		if err != nil {
			h.MsgForUser(*u.Bot, userID, "Невозможно пропустить шаг")
			return currentState, trf, err
		}
		return nextState, actualData, nil
	}

	newData, err = u.validateField(trf, msg, field)
	if err != nil {
		h.MsgForUser(*u.Bot, userID, "Невозможно добавить данные")
		return currentState, trf, err
	}

	return nextState, newData, nil
}

func (u *UserHandle) actualField(data *storage.Contract, msg *tgbotapi.Message, field int) (*storage.Contract, error) {
	actualData, err := u.ContractDB.Detail(msg.From.ID, data.ContractID)
	if err != nil {
		h.MsgForUser(*u.Bot, msg.From.ID, h.SkipText)
		return data, e.Wrap("Failed select actual data", err)
	}
	switch field {
	case fieldTariffName:
		data.TariffName = actualData.TariffName
	case fieldFullname:
		data.ContractData.FullName = actualData.ContractData.FullName
	case fieldAddress:
		data.ContractData.Address = actualData.ContractData.Address
	case fieldEmail:
		data.ContractData.Email = actualData.ContractData.Email
	case fieldPhone:
		data.ContractData.Phone = actualData.ContractData.Phone
	}
	return data, nil
}

func (u *UserHandle) validateField(data *storage.Contract, msg *tgbotapi.Message, field int) (*storage.Contract, error) {
	newData := msg.Text
	switch field {
	case fieldTariffName:
		data.TariffName = newData
	case fieldFullname:
		data.ContractData.FullName = newData
	case fieldAddress:
		data.ContractData.Address = newData
	case fieldEmail:
		data.ContractData.Email = newData
	case fieldPhone:
		data.ContractData.Phone = newData
	}
	return data, nil
}

func (u *UserHandle) chainHelper(msg *tgbotapi.Message, data interface{}, steps map[int]StepConfig, mode int) (state int, newData interface{}, err error) {
	userID := msg.From.ID
	step, exists := steps[u.F.UserState(userID)]
	if !exists {
		return 0, nil, errors.New("State is not defined")
	}
	h.MsgForUser(*u.Bot, userID, step.Prompt)
	state, newData, err = u.handle(msg, data, step.NextState, step.Field)

	if step.NextState == 0 && err == nil {
		ap := setStruct(data)
		switch mode {
		case addMode:
			newContract := idgen.IDgenerator()
			err = u.ContractDB.Add(
				ap.TariffName,
				msg.From.ID,
				newContract,
				ap.ContractData.FullName,
				ap.ContractData.Address,
				ap.ContractData.Email,
				ap.ContractData.Phone,
				h.ContractOpened)
			if err != nil {
				h.MsgForUser(*u.Bot, userID, "Произошла ошибка при оформлении")
				return 0, nil, err
			}
			go func(contractCopy *storage.Contract) {
				if err := h.EmailNotification(contractCopy); err != nil {
					log.Printf("Ошибка отправки email для контракта %d: %v", contractCopy.ContractID, err)
				}
			}(ap)
			text := fmt.Sprintf("Ваша заявки отправлена. Ее статус и номер: %s, №%d", h.ContractOpened, newContract)
			h.MsgForUser(*u.Bot, userID, text)
		case editMode:
			text := fmt.Sprintf(`Ваша заявка/договор выглядят теперь так:
Тариф: %s
ФИО: %s
Адрес проживания: %s
Эл. почта: %s
Номер телефона: %s
Ваш номер заявки/договора: %d
Текущий  статус: %s
`,
				ap.TariffName,
				ap.ContractData.FullName,
				ap.ContractData.Address,
				ap.ContractData.Email,
				ap.ContractData.Phone,
				ap.ContractID,
				ap.ContractData.Status)

			u.ContractDB.Edit(msg.From.ID,
				ap.ContractID,
				ap.TariffName,
				ap.ContractData.FullName,
				ap.ContractData.Address,
				ap.ContractData.Email,
				ap.ContractData.Phone)

			h.MsgForUser(*u.Bot, msg.From.ID, text)
		}
	}
	return state, newData, err
}

func setStruct(data interface{}) *storage.Contract {
	var appeal *storage.Contract
	if data == nil {
		appeal = &storage.Contract{}
	} else {
		appeal = data.(*storage.Contract)
	}
	return appeal
}

func (u *UserHandle) StartAdd(msg *tgbotapi.Message, data interface{}) (state int, newData interface{}, err error) {
	if cont, _ := u.ContractDB.Contract(msg.From.ID); cont != 0 {
		h.MsgForUser(*u.Bot, msg.From.ID, "У вас уже есть заявка/договор. Вы не можете создать новый")
		return 0, nil, nil
	}
	h.MsgForUser(*u.Bot, msg.From.ID, "Введите ваше ФИО полностью")
	return u.handle(msg, data, tools.ContractFullname, fieldTariffName)
}

func (u *UserHandle) AddChain(msg *tgbotapi.Message, data interface{}) (state int, newData interface{}, err error) {
	defer func() { err = e.WrapIfErr("Failed add chain handle", err) }()
	steps := map[int]StepConfig{
		tools.ContractFullname: {tools.ContractAddress, "Укажите адрес проживания", fieldFullname},
		tools.ContractAddress:  {tools.ContractEmail, "Введите адрес электронной почты", fieldAddress},
		tools.ContractEmail:    {tools.ContractPhone, "Укажите номер телефона", fieldEmail},
		tools.ContractPhone:    {0, "", fieldPhone},
	}
	return u.chainHelper(msg, data, steps, addMode)
}

func (u *UserHandle) EditChain(msg *tgbotapi.Message, data interface{}) (state int, newData interface{}, err error) {
	defer func() { err = e.WrapIfErr("Failed to edit appeal/contract", err) }()
	steps := map[int]StepConfig{
		tools.ContractEdit:         {tools.ContractEditFullname, "Введите новое ФИО" + h.SkipHint, fieldTariffName},
		tools.ContractEditFullname: {tools.ContractEditAddress, "Введите новый адрес проживания" + h.SkipHint, fieldFullname},
		tools.ContractEditAddress:  {tools.ContractEditEmail, "Введите новый адрес электронной почты" + h.SkipHint, fieldAddress},
		tools.ContractEditEmail:    {tools.ContractEditPhone, "Введите новый номер телефона" + h.SkipHint, fieldEmail},
		tools.ContractEditPhone:    {0, "", fieldPhone},
	}
	return u.chainHelper(msg, data, steps, editMode)
}

func (u *UserHandle) DetailContract(msg *tgbotapi.Message, data interface{}) (state int, newData interface{}, err error) {
	defer func() { err = e.WrapIfErr("Failed to view details", err) }()
	ap := setStruct(data)
	var userID int64
	var contract int
	if h.IsAdmin(msg) {
		contract, err = strconv.Atoi(msg.Text)
		if err != nil {
			h.MsgForUser(*u.Bot, msg.From.ID, "Введен неверный номер заявки/договора")
			return 0, nil, nil
		}
		userID, err = u.ContractDB.UserIDContract(contract)
		if err != nil {
			h.MsgForUser(*u.Bot, msg.From.ID, "У пользователя нет заяки или она не найдена")
		}

	} else {
		userID = msg.From.ID
		contract = ap.ContractID
	}
	dtls, err := u.ContractDB.Detail(userID, contract)
	if err != nil {
		h.MsgForUser(*u.Bot, msg.From.ID, "Заявки не существует или она не найдена")
		return 0, nil, nil
	}
	text := fmt.Sprintf(`Заявка с номером контракта: %d
	
Выбранный тарифный план: %s
ФИО: %s
Адрес проживания: %s
Эл. почта: %s
Номер телефона: %s

Текущий статус заявки/договора: %s
`, dtls.ContractID,
		dtls.TariffName,
		dtls.ContractData.FullName,
		dtls.ContractData.Address,
		dtls.ContractData.Email,
		dtls.ContractData.Phone,
		dtls.ContractData.Status)
	h.MsgForUser(*u.Bot, msg.From.ID, text)
	return 0, nil, nil
}
