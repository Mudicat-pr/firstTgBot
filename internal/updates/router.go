package upd

import (
	h "github.com/Mudicat-pr/firstTgBot/internal/handlers"
	"github.com/Mudicat-pr/firstTgBot/internal/storage"
	"github.com/Mudicat-pr/firstTgBot/internal/tools"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// Стандартный вид хендлера для telegram-бота
type TypeHandle func(msg *tgbotapi.Message)

// Мидлвар для проверки user id на его наличие в списке администраторов (временно внутри конфига .yaml)
type AdminHandle func(msg *tgbotapi.Message) bool

// Структура команды для роутинга и дальнейшего использования этих роутов внутри апдейтов
// Name - имя команды
// AdminOnly - булево значение, если true - команда для администрирования ботом
// State - Текущий шаг (стейт) для FSM
// Prompt - Сообщение от бота пользователю в ответ на команду
// Handle - Хендлер, закрепленный за командой
type Command struct {
	Name      int
	AdminOnly bool
	State     int
	Prompt    func(msg *tgbotapi.Message, f *tools.FSM, un *UnitedStruct) (string, tgbotapi.ReplyKeyboardMarkup)
	Handle    TypeHandle
}

// Структура для работы роутеров
// bot - сам API бота с поинтером
// fsm - Finite State Machine (его код в директории internal/tools в одноименом .go файле)
// isAdmin - Получает булево значение для проверки юзера на права администратора из мидлвара
// cmds - Карта всех команд в виде ключ:значение - строка и структура Command
type CommandRouter struct {
	bot            *tgbotapi.BotAPI
	fsm            *tools.FSM
	isAdmin        AdminHandle
	cmds           map[int]Command
	contractHandle *storage.ContractHandle
	tariffHandle   *storage.TariffHandle
}

// За подробностями возвращаемого типа смотреть структуру CommandRouter
func New(bot *tgbotapi.BotAPI, fsm *tools.FSM, isAdmin AdminHandle, contractHandle *storage.ContractHandle, tariffHandle *storage.TariffHandle) *CommandRouter {
	return &CommandRouter{
		bot:            bot,
		fsm:            fsm,
		isAdmin:        isAdmin,
		cmds:           make(map[int]Command),
		contractHandle: contractHandle,
		tariffHandle:   tariffHandle,
	}
}

type UnitedStruct struct {
	ap  *storage.ContractHandle
	trf *storage.TariffHandle
}

// Регистрация новых команд, т.е. добавление их в карту cmds структуры CommandRouter
func (r *CommandRouter) Register(cmd Command) {
	r.cmds[cmd.Name] = cmd
}

// Стандартный набор роутов для апдейтов и обработка зарегистрированных команд
func (r *CommandRouter) Handle(msg *tgbotapi.Message) {

	currentCmd := textToCmd[msg.Text]

	cmd, ok := r.cmds[currentCmd]

	isAdmin := h.IsAdmin(msg)
	// Отправка help-сообщения в ответ на неизвестную команду или сообщение
	if !ok {
		msgText := tgbotapi.NewMessage(msg.From.ID, UserHelper)
		if isAdmin {
			msgText = tgbotapi.NewMessage(msg.From.ID, AdminHelper)
		}
		msgText.ParseMode = tgbotapi.ModeHTML
		msgText.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
		r.bot.Send(msgText)
		return
	}

	un := &UnitedStruct{
		ap:  r.contractHandle,
		trf: r.tariffHandle,
	}

	// Если не совпадают поля - вывод пользовательского help-текста
	if cmd.AdminOnly && !isAdmin {
		r.bot.Send(tgbotapi.NewMessage(msg.Chat.ID, UserHelper))
		return
	}
	if cmd.Handle != nil {
		cmd.Handle(msg)
	}

	if cmd.Prompt != nil {
		text, keyboard := cmd.Prompt(msg, r.fsm, un)
		response := tgbotapi.NewMessage(msg.Chat.ID, text)
		response.ParseMode = tgbotapi.ModeHTML

		if len(keyboard.Keyboard) > 0 {
			response.ReplyMarkup = keyboard
		}
		r.bot.Send(response)
	}

	if cmd.State != 0 {
		r.fsm.SetState(msg.From.ID, cmd.State)
	}

}
