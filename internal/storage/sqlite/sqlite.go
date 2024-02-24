package sqlite

import (
	"database/sql"
	"fmt"

	"github.com/eckeriaue/golang-url-shortener/internal/storage"
	"github.com/mattn/go-sqlite3"
)

type Storage struct {
	db *sql.DB
}

func New(storagePath string) (*Storage, error) {
	const op = "storage.sqlite.New"

	db, err := sql.Open("sqlite3", storagePath)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	stmt, err := db.Prepare(`
		CREATE TABLE iF NOT EXISTS url(
			id INTEGER PRIMARY KEY,
			alias TEXT NOT NULL UNIQUE,
			url TEXT NOT NULL
		);
		CREATE INDEX IF NOT EXISTS idx_ailas ON url(alias);
	`)

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	_, err = stmt.Exec()

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{db}, nil
}

func (s *Storage) SaveUrl(urlToSave, alias string) (int64, error) {
	const op = "storage.sqlite.SaveUrl"

	stmt, err := s.db.Prepare(`insert into url("url", "alias") values(?, ?) `)
	if err != nil {
		return -1, fmt.Errorf("%s: %w", op, err)
	}
	res, err := stmt.Exec(urlToSave, alias)
	if err != nil {

		if sqliteErr, ok := err.(sqlite3.Error); ok && sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {
			return -1, fmt.Errorf("%s: %w", op, storage.ErrUrlExists)
		}

		return -1, fmt.Errorf("%s: %w", op, err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return -1, fmt.Errorf("%s: failed to get last insert id %w", op, err)
	}
	return id, nil
}
