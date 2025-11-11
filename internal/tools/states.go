package tools

const (
	AllTarrifs = 1 << iota

	TariffTitle
	TariffBody
	TariffPrice
	TariffAdd
	DelTariff

	// Скрытие тарифа (админка)
	Hide

	// Константы для редактирования договора. Сигнатура a_ от сокращения appeal
	EditTariff
	EdtiAppeal

	// Изменение тарифов
	TariffEdit
	TariffTitleEdit
	TariffBodyEdit
	TariffPriceEdit
	TariffEditConfirm
)
