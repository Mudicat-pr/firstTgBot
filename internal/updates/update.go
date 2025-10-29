package upd

import (
	"fmt"
	"log"

	h "github.com/Mudicat-pr/firstTgBot/internal/handlers"
	"github.com/Mudicat-pr/firstTgBot/internal/handlers/admin"
	"github.com/Mudicat-pr/firstTgBot/internal/handlers/user"
	"github.com/Mudicat-pr/firstTgBot/internal/storage"
	"github.com/Mudicat-pr/firstTgBot/internal/tools"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func UpdateTg(bot *tgbotapi.BotAPI,
	updates tgbotapi.UpdatesChannel,
	s *storage.Storage,
	f *tools.FSM,
	a *admin.AdminHandle,
	u *user.UserHandle) {

	for update := range updates {
		if update.Message == nil {
			continue
		}
		administrator := h.IsAdmin(update.Message)
		if f.State(update.Message.From.ID) != "" {
			f.Handle(update.Message)
			continue
		}

		if update.Message != nil {
			log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
			text := update.Message.Text

			switch text {
			case "/all":
				u.All(update.Message, h.FlagTrue)
			case "/details":
				bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Выберите ID тарифа для просмотра полного описания"))
				f.SetState(update.Message.From.ID, tools.DetailsTariff)
			case "/add":
				if !administrator {
					bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, h.UnknownCommand))
					continue
				}
				bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Введите значение в следующем виде: ИМЯ; ОПИСАНИЕ; ЦЕНА (цена только целым числом)"))
				f.SetState(update.Message.From.ID, tools.AddTariff)
			case "/del":
				if !administrator {
					bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, h.UnknownCommand))
					continue
				}
				bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Введите ID удаляемого тарифа"))
				f.SetState(update.Message.From.ID, tools.DelTariff)
			case "/hide":
				if !administrator {
					bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, h.UnknownCommand))
					continue
				}
				msg := fmt.Sprintf(`Выберите тариф по его ID и далее %s или %s. Например: 1; %s`,
					h.IsHidden, h.IsOpened, h.IsHidden)
				bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, msg))
				f.SetState(update.Message.From.ID, tools.HideByTariffID)
			case "/all_hidden":
				if !administrator {
					bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, h.UnknownCommand))
					continue
				}
				u.All(update.Message, h.FlagFalse)
			case "/edit":
				if !administrator {
					bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, h.UnknownCommand))
					continue
				}
				bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Введите текст следующим образом: ID; ИМЯ; ОПИСАНИЕ; ЦЕНА"))
				f.SetState(update.Message.From.ID, tools.EditTariff)
			case "/sub":
				bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, `
Оставьте заявку в следующем виде, где данные вводятся черезз точку с запятой:
Имя тарифа; ФИО; Адрес проживания; Электронная почта; Номер телефона`))
				f.SetState(update.Message.From.ID, tools.SubmitAppeal)
			default:
				bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, h.UnknownCommand))
			}
		}
	}
}

//Введите команду "%s" или "%s"`, admin.IsHidden, admin.IsOpened
