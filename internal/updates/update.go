package upd

import (
	"fmt"
	"log/slog"

	h "github.com/Mudicat-pr/firstTgBot/internal/handlers"
	"github.com/Mudicat-pr/firstTgBot/internal/handlers/admin"
	"github.com/Mudicat-pr/firstTgBot/internal/handlers/user"
	"github.com/Mudicat-pr/firstTgBot/internal/storage"
	"github.com/Mudicat-pr/firstTgBot/internal/tools"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const UserHelper = `
Для общения с ботом используйте следующие ключевые слова:

<b>Все тарифы</b> - <i>Список всех доступных тарифных планов</i>

<b>О тарифе</b> - <i>Просмотр деталей тарифного плана</i>

<b>Создать заявку</b> - <i>Заполнение формы для отправки вашей заявки на обработку</i>

<b>Моя заявка или Мой договор</b> - <i>Просмотр вашей заявки/договора</i>

<b>Изменить заявку</b> - <i>Перезаполнить форму для изменения заявки</i>
`

var AdminHelper = fmt.Sprintf(`
%s

Команды для управления (доступны только администраторам):
<b>Все спрятанные</b> - <i>Список всех спрятанных тарифов</i>

<b>Новый тариф</b> - <i>Создание нового тарифного плата и его добавление в список Все</i>

<b>Удалить тариф</b> - <i>Удаленный тариф вернуть будет невозможно</i>

<b>Изменить тариф</b> - <i>Перезаполнение тарифа с пропуском полей</i>

<b>Спрятать тариф, Открыть тариф</b> - <i>Чтоб тариф не было видно в общем списке. Альтернатива удалению</i>

<b>Изменить статус</b> - <i>При изменении статуса пользователь получит уведомление</i>

<b>Удалить заявку или Удалить договор</b> - <i>Рекомендую перед этим изменить статус, чтоб пользователь получил уведомление</i>
`, UserHelper)

func UpdateTg(bot *tgbotapi.BotAPI,
	updates tgbotapi.UpdatesChannel,
	s *storage.Storage,
	f *tools.FSM,
	a *admin.AdminHandle,
	u *user.UserHandle,
	contractHandle *storage.ContractHandle,
	tariffHandle *storage.TariffHandle) {

	r := New(bot, f, h.IsAdmin, contractHandle, tariffHandle) // Роутер

	r.Register(Command{
		Name: cmdAll,
		Handle: func(msg *tgbotapi.Message) {
			u.All(msg, h.FlagTrue)
		},
	})
	r.Register(Command{
		Name:   cmdDetailsT,
		State:  tools.TariffDetails,
		Prompt: PromptDetailsTariff,
	})

	r.Register(Command{
		Name:   cmdSubmit,
		State:  tools.ContractSubmit,
		Prompt: PromptContractCreate,
	})
	r.Register(Command{
		Name:      cmdAddT,
		AdminOnly: h.FlagTrue,
		State:     tools.TariffTitle,
		Prompt:    PromptAddTariff,
	})
	r.Register(Command{
		Name:      cmdDelT,
		AdminOnly: h.FlagTrue,
		State:     tools.DelTariff,
		Prompt:    PromptDelTariff,
	})
	r.Register(Command{
		Name:      cmdHide,
		AdminOnly: h.FlagTrue,
		State:     tools.Hide,
		Prompt:    PromptHideTariff,
	})
	r.Register(Command{
		Name:      cmdAllHidden,
		AdminOnly: h.FlagTrue,
		Handle: func(msg *tgbotapi.Message) {
			u.All(msg, h.FlagFalse)
		},
	})
	r.Register(Command{
		Name:   cmdEditC,
		State:  tools.ContractEdit,
		Prompt: PromptEditContract,
	})
	r.Register(Command{
		Name:      cmdEditT,
		State:     tools.TariffEdit,
		AdminOnly: h.FlagTrue,
		Prompt:    PromptEditTariff,
	})
	r.Register(Command{
		Name:   cmdDetailsC,
		State:  tools.ContractDetails,
		Prompt: PromptDetailsContract,
	})
	r.Register(Command{
		Name:      cmdDelC,
		State:     tools.DeleteContract,
		AdminOnly: h.FlagTrue,
		Prompt:    PromptDeleteContract,
	})
	r.Register(Command{
		Name:      cmdSwitch,
		State:     tools.SwitchStart,
		AdminOnly: h.FlagTrue,
		Prompt:    PromptSwitchStatus,
	})
	for update := range updates {
		if update.Message == nil {
			continue
		}
		msg := update.Message

		if msg.Text == "/cancel" {
			if f.UserState(msg.From.ID) != 0 {
				f.ClearState(msg.From.ID)
				h.MsgForUser(*bot, msg.From.ID, "Операция успешно отменена!")
			} else {
				h.MsgForUser(*bot, msg.From.ID, "Нет операций для отмены")
			}
		}

		if state := f.UserState(msg.From.ID); state != 0 {
			f.HandleState(msg)
			continue
		}
		slog.Info("BOT REQUEST", "Username", msg.From.UserName, "User input", msg.Text)
		r.Handle(msg)
	}
}
