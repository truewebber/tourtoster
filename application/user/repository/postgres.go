package repository

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"

	"tourtoster/user"
)

type (
	Postgres struct {
		db *sql.DB
	}
)

const (
	selectUserByID    = "SELECT email, status, role FROM users WHERE id=$1;"
	selectUserByEmail = "SELECT id, password_hash, status, role FROM users WHERE email=$1;"
)

func NewPostgres(db *sql.DB) *Postgres {
	return &Postgres{
		db: db,
	}
}

func (p *Postgres) User(ID int64) (*user.User, error) {
	u := new(user.User)
	if err := p.db.QueryRow(selectUserByID, ID).Scan(&u.Email, &u.Status, &u.Role); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, err
	}
	u.ID = ID

	return u, nil
}

func (p *Postgres) UserWithEmail(email string) (*user.User, error) {
	u := new(user.User)
	if err := p.db.QueryRow(selectUserByEmail, email).Scan(&u.ID, &u.PasswordHash, &u.Status, &u.Role); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, err
	}
	u.Email = email

	return u, nil
}

func (p *Postgres) Save(_ *user.User) error {
	return nil
}

func (p *Postgres) Delete(_ int64) error {
	return nil
}
