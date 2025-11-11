package upd

import (
	"log"

	h "github.com/Mudicat-pr/firstTgBot/internal/handlers"
	"github.com/Mudicat-pr/firstTgBot/internal/handlers/admin"
	"github.com/Mudicat-pr/firstTgBot/internal/handlers/user"
	"github.com/Mudicat-pr/firstTgBot/internal/storage"
	"github.com/Mudicat-pr/firstTgBot/internal/tools"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const UserHelper = `Я могу помочь вам  выбрать интересующий вас тариф нашего сотового оператора, и оставить заявку по подключению!

/all - Список всех доступных тарифных планов.
/details - Открыть описание интересующего вас тарифа.
/submit - Оставить заявку на подключение тарифного плана.
/cancel - Отменить действие.

Я ещё молодой бот, возможно список доступных команд будет расширяться. Если я вам буду нужен - просто напишите любое сообщение в чат☺️`

const AdminHelper = `Дорогой администратор, вот все доступные команды:

Пользовательские (общие):
/all - Список всех доступных тарифных планов.
/details - Открыть описание интересующего вас тарифа.
/submit - Оставить заявку на подключение тарифного плана.
/cancel - Отменить действие.

Команды для управления (доступны только администраторам):
/add - Добавить новый тарифный план.
/del - Удалить тарифный план. ‼️ВНИМАНИЕ‼: Удаленный тарифный план не подлежит восстановлению.
/hide - Спрятать тарифный план. Спрятанный тарифный план невидим для списка /all у пользователя.
/all_hidden - Просмотреть все скрытые тарифные планы.
/edit - Изменить тарифный план.
`

func UpdateTg(bot *tgbotapi.BotAPI,
	updates tgbotapi.UpdatesChannel,
	s *storage.Storage,
	f *tools.FSM,
	a *admin.AdminHandle,
	u *user.UserHandle) {

	r := New(bot, f, h.IsAdmin) // Роутер

	r.Register(Command{
		Name: "/all",
		Handle: func(msg *tgbotapi.Message) {
			u.All(msg, h.FlagTrue)
		},
	})
	/*
		r.Register(Command{
			Name:   "/details",
			State:  tools.DetailsTariff,
			Prompt: PromptDetailsTariff + CancelMessage,
		})
		r.Register(Command{
			Name:   "/submit",
			State:  tools.SubmitAppeal,
			Prompt: PromptSubmitTariff + CancelMessage,
		}) */
	r.Register(Command{
		Name:      "/add",
		AdminOnly: h.FlagTrue,
		State:     tools.TariffTitle,
		Prompt:    PromptAddTariff + CancelMessage,
	})
	r.Register(Command{
		Name:      "/del",
		AdminOnly: h.FlagTrue,
		State:     tools.DelTariff,
		Prompt:    PromptDelTariff + CancelMessage,
	})
	r.Register(Command{
		Name:      "/hide",
		AdminOnly: h.FlagTrue,
		State:     tools.Hide,
		Prompt:    PromptHideTariff + CancelMessage,
	})
	r.Register(Command{
		Name:      "/all_hidden",
		AdminOnly: h.FlagTrue,
		Handle: func(msg *tgbotapi.Message) {
			u.All(msg, h.FlagFalse)
		},
	})
	r.Register(Command{
		Name:   "/edit_appeal",
		State:  tools.EdtiAppeal,
		Prompt: PromptEditAppeal + CancelMessage,
	})
	r.Register(Command{
		Name:      "/edit_t",
		State:     tools.TariffEdit,
		AdminOnly: h.FlagTrue,
		Prompt:    PromptEditTariff + CancelMessage,
	})
	for update := range updates {
		if update.Message == nil {
			continue
		}
		msg := update.Message

		if msg.Text == "/cancel" {
			f.ClearState(msg.From.ID)
			bot.Send(tgbotapi.NewMessage(msg.Chat.ID, "Операция успешно отменена!"))
		}

		if state := f.UserState(msg.From.ID); state != 0 {
			log.Printf("Текущий стейт: %v", state)
			f.HandleState(msg)
			continue
		}
		log.Printf("[%s] %s", msg.From.UserName, msg.Text)
		log.Printf("Before router: user state = %d", f.UserState(msg.From.ID))
		r.Handle(msg)
		log.Printf("After router: user state = %d", f.UserState(msg.From.ID))
	}
}

//Введите команду "%s" или "%s"`, admin.IsHidden, admin.IsOpened
