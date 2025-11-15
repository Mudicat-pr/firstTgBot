package upd

import (
	h "github.com/Mudicat-pr/firstTgBot/internal/handlers"
	"github.com/Mudicat-pr/firstTgBot/internal/tools"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const CancelMessage = "\nЧтобы отменить действие, введите /cancel"

// Команды для бота. E: ADD: "Добавить"
const (
	cmdAll = 1 << iota
	cmdDetailsC
	cmdDetailsT
	cmdSubmit
	cmdContract
	cmdEditC
	cmdEditT
	cmdAllHidden
	cmdAddT
	cmdDelT
	cmdDelC
	cmdHide
	cmdSwitch
)

var (
	Cmds = map[int][]string{
		// Пользовательские
		cmdAll:      {"Все", "Все тарифы", "Тарифы все"},
		cmdDetailsT: {"Подробнее", "Детали", "О тарифе"},
		cmdSubmit:   {"Отправить заявку", "Создать заявку", "Сформировать заявку"},
		cmdDetailsC: {"Заявка", "Моя заявка", "Договор", "Мой договор"},
		cmdEditC:    {"Изменить заявку"},

		// Команды администратора
		cmdAllHidden: {"Спрятанные", "Все спрятанные"},
		cmdAddT:      {"Новый тариф", "Создать тариф", "Добавить тариф"},
		cmdDelT:      {"Удалить тариф", "Тариф удалить"},
		cmdEditT:     {"Изменить тариф"},
		cmdHide:      {"Спрятать тариф", "Открыть тариф"},
		cmdSwitch:    {"Изменить статус заявки", "Изменить статус договора", "Изменить статус"},
		cmdDelC:      {"Удалить заявку", "Удалить договор"},
	}
	textToCmd = reverseCmdMap(Cmds)
)

func reverseCmdMap(Cmds map[int][]string) map[string]int {
	rev := make(map[string]int)
	for bit, aliases := range Cmds {
		for _, alias := range aliases {
			rev[alias] = bit
		}

	}
	return rev
}

func PromptAddTariff(msg *tgbotapi.Message, f *tools.FSM, un *UnitedStruct) (string, tgbotapi.ReplyKeyboardMarkup) {
	return "Введите имя для нового тарифа" + CancelMessage, tgbotapi.ReplyKeyboardMarkup{}
}

func PromptDelTariff(msg *tgbotapi.Message, f *tools.FSM, un *UnitedStruct) (string, tgbotapi.ReplyKeyboardMarkup) {
	return `
Для удаления тарифа введите его ID.

<b>ВЕРНУТЬ УДАЛЕННЫЙ ТАРИФ БУДЕТ НЕВОЗМОЖНО</b>` + CancelMessage, tgbotapi.ReplyKeyboardMarkup{}
}

func PromptEditTariff(msg *tgbotapi.Message, f *tools.FSM, un *UnitedStruct) (string, tgbotapi.ReplyKeyboardMarkup) {
	return "Введите ID тарифа для его изменения" + CancelMessage, tgbotapi.ReplyKeyboardMarkup{}
}

func PromptDetailsTariff(msg *tgbotapi.Message, f *tools.FSM, un *UnitedStruct) (string, tgbotapi.ReplyKeyboardMarkup) {
	_, err := un.trf.AllTariffs()
	if err != nil {
		f.ClearState(msg.From.ID)
		return "Нет тарифов для просмотра подробностей", tgbotapi.ReplyKeyboardMarkup{}
	}
	return "Введите ID интересующего тарифа" + CancelMessage, tgbotapi.ReplyKeyboardMarkup{}
}

func PromptHideTariff(msg *tgbotapi.Message, f *tools.FSM, un *UnitedStruct) (string, tgbotapi.ReplyKeyboardMarkup) {
	return "Введите ID тарифа для его скрытия", tgbotapi.ReplyKeyboardMarkup{}
}

func PromptContractCreate(msg *tgbotapi.Message, f *tools.FSM, un *UnitedStruct) (string, tgbotapi.ReplyKeyboardMarkup) {
	contract, err := un.ap.Contract(msg.From.ID)
	if err != nil {
		f.ClearState(msg.From.ID)
		return "Произошла неизвестная ошибка", tgbotapi.ReplyKeyboardMarkup{}
	}
	if contract > 0 {
		f.ClearState(msg.From.ID)
		return "У вас уже есть действующий тариф, создать новый невозможно", tgbotapi.ReplyKeyboardMarkup{}
	}
	return "Введите имя тарифного плана" + CancelMessage, tgbotapi.ReplyKeyboardMarkup{}
}

func PromptEditContract(msg *tgbotapi.Message, f *tools.FSM, un *UnitedStruct) (string, tgbotapi.ReplyKeyboardMarkup) {
	contract, err := un.ap.Contract(msg.From.ID)
	if err != nil {
		f.ClearState(msg.From.ID)
		return "Произошла неизвестная ошибка", tgbotapi.ReplyKeyboardMarkup{}
	}
	if contract == 0 {
		f.ClearState(msg.From.ID)
		return "У вас нет существующей заявки на подключение", tgbotapi.ReplyKeyboardMarkup{}
	}

	return "Введите имя интересующего тарифа", h.CreateSkipKey()
}

func PromptDeleteContract(msg *tgbotapi.Message, f *tools.FSM, un *UnitedStruct) (string, tgbotapi.ReplyKeyboardMarkup) {
	return `Введите номер контракта.
	
Внимание: после удаления вернуть данные о договоре будет невозможно` + CancelMessage, tgbotapi.ReplyKeyboardMarkup{}
}

func PromptSwitchStatus(msg *tgbotapi.Message, f *tools.FSM, un *UnitedStruct) (string, tgbotapi.ReplyKeyboardMarkup) {
	return "Введите номер контракта/заявки для изменения статуса", tgbotapi.ReplyKeyboardMarkup{}
}

func PromptDetailsContract(msg *tgbotapi.Message, f *tools.FSM, un *UnitedStruct) (string, tgbotapi.ReplyKeyboardMarkup) {
	return "Введите номер договора/заявки", tgbotapi.ReplyKeyboardMarkup{}
}
