package tools

const (
	AllTarrifs = 1 << iota

	TariffTitle
	TariffBody
	TariffPrice
	TariffAdd
	TariffAddConfirm

	DelTariff

	// Скрытие тарифа (админка)
	Hide

	//
	ContractDetails

	// Изменение тарифов
	TariffEdit
	TariffTitleEdit
	TariffBodyEdit
	TariffPriceEdit
	TariffEditConfirm

	TariffDetails

	// Создание заявки
	ContractSubmit
	ContractFullname
	ContractAddress
	ContractEmail
	ContractPhone

	// Редактирование заявки
	ContractEdit
	ContractEditFullname
	ContractEditAddress
	ContractEditEmail
	ContractEditPhone
	ContractEditConfirm

	DeleteContract

	SwitchStart
	SwitchEnd
)
