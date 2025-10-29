package storage

import "github.com/Mudicat-pr/firstTgBot/pkg/e"

type TariffHandle struct {
	S *Storage
}

type TariffManager interface {
	Adder
	Deleter
	Editor
	DetailViewer
}

func (t *TariffHandle) AllTariffs() (tariff []Tariff, err error) {
	defer func() { err = e.WrapIfErr("can't query tariffs", err) }()
	rows, err := t.S.db.Query("SELECT id, title, price, is_hide FROM tariffs")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tariffs []Tariff

	for rows.Next() {
		var trf Tariff
		if err = rows.Scan(&trf.ID, &trf.Title, &trf.Price, &trf.IsHide); err != nil {
			return tariffs, err
		}
		/*
			if trf.IsHide {
				continue
			}*/
		tariffs = append(tariffs, trf)
	}
	if err = rows.Err(); err != nil {
		return tariffs, err
	}
	return tariffs, nil
}

func (t *TariffHandle) Details(tariffID int) (tariff Tariff, err error) {
	defer func() { err = e.WrapIfErr("Failed to select data", err) }()
	ctx, cancel := e.Ctx()
	defer cancel()
	var trf Tariff
	err = t.S.db.QueryRowContext(ctx, "SELECT id, title, body, price, is_hide FROM tariffs WHERE id = ?", tariffID).
		Scan(&trf.ID, &trf.Title, &trf.Body, &trf.Price, &trf.IsHide)
	if err != nil {
		return trf, err
	}
	return trf, nil
}

func (t *TariffHandle) Add(title, description string, price int) error {
	ctx, cancel := e.Ctx()
	q := "INSERT INTO tariffs(title, body, price) VALUES (?,?,?)"
	defer cancel()

	return t.S.ExecQuery(ctx, q, title, description, price)
}

func (t *TariffHandle) Del(tariffID int) error {
	ctx, cancel := e.Ctx()
	q := "DELETE FROM tariffs WHERE id = ?"
	defer cancel()

	return t.S.ExecQuery(ctx, q, tariffID)
}

func (t *TariffHandle) Hide(tariffID int, isHide bool) error {
	ctx, cancel := e.Ctx()
	q := "UPDATE tariffs SET is_hide = ? WHERE id = ?"
	defer cancel()

	return t.S.ExecQuery(ctx, q, isHide, tariffID)
}

func (t *TariffHandle) Edit(tariffID int, title, body string, price int) error {
	ctx, cancel := e.Ctx()
	q := "UPDATE tariffs SET title = ?, body = ?, price = ? WHERE id = ?"
	defer cancel()

	return t.S.ExecQuery(ctx, q, title, body, price, tariffID)
}
