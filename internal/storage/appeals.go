package storage

import "github.com/Mudicat-pr/firstTgBot/pkg/e"

type AppealHandle struct {
	S *Storage
}

type AppealManager interface {
	Adder
	Deleter
	Editor
	DetailViewer
}

func (a *AppealHandle) Details(userID int64) (bool, error) {
	data, err := a.S.db.Query("SELECT id_user FROM appeals WHERE id_user = ?")
	if err != nil || data == nil {
		return false, err
	}
	return true, nil
}

func (a *AppealHandle) Add(tariffName string, userID int64, contract int, fullname, address, email, phone string) (err error) {
	defer func() { err = e.WrapIfErr("Can't insert data to appeals", err) }()

	ctx, cancel := e.Ctx()
	defer cancel()
	_, err = a.S.db.ExecContext(ctx, `INSERT INTO 
	appeals(tariff_name, id_user, contract, fullname, address, email, phone)
	VALUES (?,?,?,?,?,?,?)`,
		tariffName, userID, contract, fullname, address, email, phone)
	if err != nil {
		return err
	}

	return nil
}

// Берет строку из FSM и в зависимости от нее дает функционал к изменению записи
func (a *AppealHandle) Edit(userID int64, contractID int, tariffName, fullname, address, email, phone string) (err error) {
	ctx, cancel := e.Ctx()
	q := `UPDATE appeals SET 
	tariff_name = ?,
	fullname = ?,
	address = ?,
	email = ?,
	phone = ?,
	WHERE contract = ? AND id_user = ?`
	defer cancel()
	_, err = a.S.db.ExecContext(ctx, q, tariffName, fullname, address, email, phone, contractID, userID)
	return err
}

func (a *AppealHandle) DataForEdit(contractID int) (data Appeal, err error) {
	defer func() { err = e.WrapIfErr("Can't select data from table appeals", err) }()
	var ap Appeal
	ctx, cancel := e.Ctx()
	defer cancel()
	err = a.S.db.QueryRowContext(ctx, "SELECT tariff_name, fullname, address, email, phone FROM appeals WHERE contract = ?", contractID).
		Scan(&ap.TariffName, &ap.AppealData.FullName,
			&ap.AppealData.Address, &ap.AppealData.Email, &ap.AppealData.Phone)
	if err != nil {
		return ap, err
	}
	return ap, nil
}

func (a *AppealHandle) Del(contractID int) (err error) {
	return nil // Заглушка
}
