package upd

import (
	"log"

	h "github.com/Mudicat-pr/firstTgBot/internal/handlers"
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
	Name      string
	AdminOnly bool
	State     int
	Prompt    string
	Handle    TypeHandle
}

// Структура для работы роутеров
// bot - сам API бота с поинтером
// fsm - Finite State Machine (его код в директории internal/tools в одноименом .go файле)
// isAdmin - Получает булево значение для проверки юзера на права администратора из мидлвара
// cmds - Карта всех команд в виде ключ:значение - строка и структура Command
type CommandRouter struct {
	bot     *tgbotapi.BotAPI
	fsm     *tools.FSM
	isAdmin AdminHandle
	cmds    map[string]Command
}

// За подробностями возвращаемого типа смотреть структуру CommandRouter
func New(bot *tgbotapi.BotAPI, fsm *tools.FSM, isAdmin AdminHandle) *CommandRouter {
	return &CommandRouter{
		bot:     bot,
		fsm:     fsm,
		isAdmin: isAdmin,
		cmds:    make(map[string]Command),
	}
}

// Регистрация новых команд, т.е. добавление их в карту cmds структуры CommandRouter
func (r *CommandRouter) Register(cmd Command) {
	r.cmds[cmd.Name] = cmd
}

// Стандартный набор роутов для апдейтов и обработка зарегистрированных команд
func (r *CommandRouter) Handle(msg *tgbotapi.Message) {
	cmd, ok := r.cmds[msg.Text]
	log.Printf("Command found: %v, state: %d", ok, cmd.State)
	isAdmin := h.IsAdmin(msg)

	response := tgbotapi.NewMessage(msg.Chat.ID, cmd.Prompt)
	response.ParseMode = tgbotapi.ModeHTML
	// Отправка help-сообщения в ответ на неизвестную команду или сообщение
	if !ok {
		text := UserHelper
		if isAdmin {
			text = AdminHelper
		}
		r.bot.Send(tgbotapi.NewMessage(msg.Chat.ID, text))
		return
	}

	// Если не совпадают поля - вывод пользовательского help-текста
	if cmd.AdminOnly && !isAdmin {
		r.bot.Send(tgbotapi.NewMessage(msg.Chat.ID, UserHelper))
		return
	}
	if cmd.Handle != nil {
		cmd.Handle(msg)
	}

	if cmd.Prompt != "" {
		r.bot.Send(response)
	}

	if cmd.State != 0 {
		r.fsm.SetState(msg.From.ID, cmd.State)
	}

}
