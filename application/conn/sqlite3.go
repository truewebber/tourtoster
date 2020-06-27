package conn

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
	"github.com/pkg/errors"
)

func NewConn(file string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", file)
	if err != nil {
		return nil, errors.Wrap(err, "error connect to db")
	}
	if err := db.Ping(); err != nil {
		return nil, errors.Wrap(err, "ping error")
	}

	return db, nil
}
