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
		Bot:        bot,
		F:          fsm,
		ContractDB: &storage.ContractHandle{S: store},
		TariffDB:   &storage.TariffHandle{S: store},
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

	f.Register(tools.TariffTitle, tools.BindState(a, (*admin.AdminHandle).AddChain))
	f.Register(tools.TariffBody, tools.BindState(a, (*admin.AdminHandle).AddChain))
	f.Register(tools.TariffPrice, tools.BindState(a, (*admin.AdminHandle).AddChain))

	f.Register(tools.DelTariff, tools.BindState(a, (*admin.AdminHandle).Del))

	f.Register(tools.TariffEdit, tools.BindState(a, (*admin.AdminHandle).StartEdit))
	f.Register(tools.TariffTitleEdit, tools.BindState(a, (*admin.AdminHandle).EditChain))
	f.Register(tools.TariffBodyEdit, tools.BindState(a, (*admin.AdminHandle).EditChain))
	f.Register(tools.TariffPriceEdit, tools.BindState(a, (*admin.AdminHandle).EditChain))

	f.Register(tools.Hide, tools.BindState(a, (*admin.AdminHandle).Hide))

	f.Register(tools.DeleteContract, tools.BindState(a, (*admin.AdminHandle).DelContract))

	f.Register(tools.SwitchStart, tools.BindState(a, (*admin.AdminHandle).ContractStatus))
	f.Register(tools.SwitchEnd, tools.BindState(a, (*admin.AdminHandle).SwitchStatus))

	// User states
	f.Register(tools.TariffDetails, tools.BindState(u, (*user.UserHandle).DetailTariff))

	f.Register(tools.ContractSubmit, tools.BindState(u, (*user.UserHandle).StartAdd))
	f.Register(tools.ContractFullname, tools.BindState(u, (*user.UserHandle).AddChain))
	f.Register(tools.ContractAddress, tools.BindState(u, (*user.UserHandle).AddChain))
	f.Register(tools.ContractEmail, tools.BindState(u, (*user.UserHandle).AddChain))
	f.Register(tools.ContractPhone, tools.BindState(u, (*user.UserHandle).AddChain))
	f.Register(tools.ContractEdit, tools.BindState(u, (*user.UserHandle).EditChain))
	f.Register(tools.ContractEditFullname, tools.BindState(u, (*user.UserHandle).EditChain))
	f.Register(tools.ContractEditAddress, tools.BindState(u, (*user.UserHandle).EditChain))
	f.Register(tools.ContractEditEmail, tools.BindState(u, (*user.UserHandle).EditChain))
	f.Register(tools.ContractEditPhone, tools.BindState(u, (*user.UserHandle).EditChain))

	f.Register(tools.ContractDetails, tools.BindState(u, (*user.UserHandle).DetailContract))

}
