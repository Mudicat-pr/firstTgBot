package storage

import (
	"context"
	"database/sql"

	"github.com/Mudicat-pr/firstTgBot/pkg/e"
	_ "github.com/mattn/go-sqlite3"
)

type Tariff struct {
	ID     int
	Title  string
	Body   string
	Price  int
	IsHide bool
}

type Contract struct {
	ID           int
	UserID       int64
	ContractID   int
	TariffName   string
	ContractData ContractData
}

type ContractData struct {
	FullName string
	Address  string
	Email    string
	Phone    string
	Status   string
}

type Adder interface {
	Add(args ...any) (err error)
}

type Deleter interface {
	Del(args ...any) (err error)
}

type Editor interface {
	Edit(args ...any) (err error)
}

type DetailViewer interface {
	Details(id int) (err error) // Любое id - юзера, договора
}

type Storage struct {
	db *sql.DB
}

func New(storagePath string) (s *Storage, err error) {
	defer func() { err = e.WrapIfErr("can't create tables or database", err) }()
	db, err := sql.Open("sqlite3", storagePath)
	if err != nil {
		return nil, err
	}

	// Можно дополнить количество таблиц в этот срез
	tables := []string{

		`CREATE TABLE IF NOT EXISTS tariffs(
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT NOT NULL,
		body TEXT NOT NULL,
		price INTEGER,
		is_hide INTEGER NOT NULL DEFAULT 0
		)`,

		`CREATE TABLE IF NOT EXISTS contracts(
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		id_user INTEGER NOT NULL,
		tariff_name TEXT NOT NULL,
		contract INTEGER NOT NULL,
		fullname TEXT NOT NULL,
		address TEXT NOT NULL,
		email TEXT NOT NULL,
		phone TEXT NOT NULL,
		status TEXT NOT NULL
		)`,
	}

	for _, table := range tables {
		_, err := db.Exec(table)
		if err != nil {
			return nil, err
		}
	}
	return &Storage{db: db}, nil
}

func (s *Storage) ExecQuery(ctx context.Context, query string, args ...any) error {
	_, err := s.db.ExecContext(ctx, query, args...)
	if err != nil {
		return e.Wrap("can't execute task", err)
	}
	return nil
}
