package storage

import (
	"github.com/Mudicat-pr/firstTgBot/pkg/e"
)

type ContractHandle struct {
	S *Storage
}

type ContractManager interface {
	Adder
	Deleter
	Editor
	DetailViewer
}

func (a *ContractHandle) Add(tariffName string, userID int64, contract int, fullname, address, email, phone, status string) (err error) {
	defer func() { err = e.WrapIfErr("Can't insert data to contracts", err) }()

	ctx, cancel := e.Ctx()
	defer cancel()
	_, err = a.S.db.ExecContext(ctx, `INSERT INTO 
	contracts(tariff_name, id_user, contract, fullname, address, email, phone, status)
	VALUES (?,?,?,?,?,?,?,?)`,
		tariffName, userID, contract, fullname, address, email, phone, status)
	if err != nil {
		return err
	}

	return nil
}

func (a *ContractHandle) Contract(userID int64) (contract int, err error) {
	defer func() { err = e.WrapIfErr("Cannot get actual contract", err) }()
	ctx, cancel := e.Ctx()
	defer cancel()

	var ap Contract
	q := `SELECT contract FROM contracts WHERE id_user = ?`
	err = a.S.db.QueryRowContext(ctx, q, userID).Scan(&ap.ContractID)
	if err != nil {
		return 0, nil
	}
	return ap.ContractID, nil
}

func (a *ContractHandle) UserIDContract(contractID int) (userID int64, err error) {
	defer func() { err = e.WrapIfErr("Can't get userID for contract ID", err) }()
	ctx, cancel := e.Ctx()
	defer cancel()

	var ap Contract
	q := "SELECT id_user FROM contracts WHERE contract = ?"
	err = a.S.db.QueryRowContext(ctx, q, contractID).Scan(&ap.UserID)
	if err != nil {
		return 0, nil
	}
	return ap.UserID, nil
}

func (a *ContractHandle) Detail(userID int64, contractID int) (d *Contract, err error) {
	defer func() { err = e.WrapIfErr("Failed to select contracts data", err) }()
	ctx, cancel := e.Ctx()
	defer cancel()
	var ap Contract

	q := `SELECT tariff_name, contract, fullname, address, email, phone, status FROM contracts WHERE id_user = ? AND contract = ?`
	err = a.S.db.QueryRowContext(ctx, q, userID, contractID).
		Scan(&ap.TariffName,
			&ap.ContractID,
			&ap.ContractData.FullName,
			&ap.ContractData.Address,
			&ap.ContractData.Email,
			&ap.ContractData.Phone,
			&ap.ContractData.Status)
	if err != nil {
		return &ap, err
	}
	return &ap, nil
}

// Берет строку из FSM и в зависимости от нее дает функционал к изменению записи
func (a *ContractHandle) Edit(userID int64, contractID int, tariffName, fullname, address, email, phone string) (err error) {
	ctx, cancel := e.Ctx()
	q := `UPDATE contracts SET 
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

func (a *ContractHandle) Del(contractID int) (err error) {
	ctx, cancel := e.Ctx()
	q := "DELETE FROM contracts WHERE contract = ?"
	defer cancel()
	return a.S.ExecQuery(ctx, q, contractID)
}

func (a *ContractHandle) SwitchStatus(contractID int, status string) error {
	ctx, cancel := e.Ctx()
	q := "UPDATE contracts SET status = ? WHERE contract = ?"
	defer cancel()
	_, err := a.S.db.ExecContext(ctx, q, status, contractID)
	return err
}
