package main

import (
	"log"

	"github.com/Mudicat-pr/firstTgBot/config"
	"github.com/Mudicat-pr/firstTgBot/internal/handlers"
	"github.com/Mudicat-pr/firstTgBot/internal/handlers/admin"
	"github.com/Mudicat-pr/firstTgBot/internal/handlers/user"
	"github.com/Mudicat-pr/firstTgBot/internal/storage"
	"github.com/Mudicat-pr/firstTgBot/internal/tools"
	upd "github.com/Mudicat-pr/firstTgBot/internal/updates"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func main() {
	store, err := storage.New("./internal/storage/storage.db")
	if err != nil {
		log.Fatal(err)
	}

	token := config.MustToken()

	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true
	log.Printf("Authorized on account %s", bot.Self.UserName)

	fsm := tools.New()

	baseVariables := handlers.BaseVar{
		Bot:      bot,
		F:        fsm,
		AppealDB: &storage.AppealHandle{S: store},
		TariffDB: &storage.TariffHandle{S: store},
	}
	adminVar := admin.AdminHandle{BaseVar: &baseVariables}
	userVar := user.UserHandle{BaseVar: &baseVariables}
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := bot.GetUpdatesChan(u)
	regState(fsm, &adminVar, &userVar)
	upd.UpdateTg(bot, updates, store, fsm, &adminVar, &userVar)
}

func regState(f *tools.FSM, a *admin.AdminHandle, u *user.UserHandle) {
	// Admin states

	f.Register(tools.TariffTitle, tools.BindState(a, (*admin.AdminHandle).SetTitle))
	f.Register(tools.TariffBody, tools.BindState(a, (*admin.AdminHandle).SetBody))
	f.Register(tools.TariffPrice, tools.BindState(a, (*admin.AdminHandle).SetPrice))
	f.Register(tools.DelTariff, tools.BindState(a, (*admin.AdminHandle).Del))

	f.Register(tools.TariffEdit, tools.BindState(a, (*admin.AdminHandle).StartEdit))
	f.Register(tools.TariffTitleEdit, tools.BindState(a, (*admin.AdminHandle).EditTitle))
	f.Register(tools.TariffBodyEdit, tools.BindState(a, (*admin.AdminHandle).EditBody))
	f.Register(tools.TariffPriceEdit, tools.BindState(a, (*admin.AdminHandle).EditPrice))
	f.Register(tools.TariffEditConfirm, tools.BindState(a, (*admin.AdminHandle).EditConfirm))

	f.Register(tools.Hide, tools.BindState(a, (*admin.AdminHandle).Hide))
	/*
		f.Register(tools.HideByTariffID, tools.BindState(a, (*admin.AdminHandle).HideByID))
		f.Register(tools.EditTariff, tools.BindState(a, (*admin.AdminHandle).Edit))

		// User states
		f.Register(tools.DetailsTariff, tools.BindState(u, (*user.UserHandle).Detail))
		f.Register(tools.SubmitAppeal, tools.BindState(u, (*user.UserHandle).Add))

	*/
}
