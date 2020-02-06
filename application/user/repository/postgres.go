package repository

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
	"github.com/pkg/errors"

	"tourtoster/hotel"
	"tourtoster/user"
)

type (
	Postgres struct {
		db *sql.DB
	}
)

const (
	selectUserByID = `SELECT first_name, second_name, last_name, hotel_name, hotel_id,
       						note, email, phone, password_hash, status, role
					FROM users WHERE id=$1;`
	selectUserByEmail = `SELECT id, first_name, second_name, last_name, hotel_name, hotel_id,
       						note, phone, password_hash, status, role
						FROM users WHERE email=$1;`
	selectUsers = `SELECT id, first_name, second_name, last_name, hotel_name, hotel_id,
       						note, email, phone, password_hash, status, role
					FROM users;`
)

func NewPostgres(db *sql.DB) *Postgres {
	return &Postgres{
		db: db,
	}
}

func (p *Postgres) User(ID int64) (*user.User, error) {
	u := new(user.User)
	u.Hotel = new(hotel.Hotel)

	if err := p.db.QueryRow(selectUserByID, ID).Scan(
		&u.FirstName, &u.SecondName, &u.LastName, &u.Hotel.Name, &u.Hotel.ID,
		&u.Note, &u.Email, &u.Phone, &u.PasswordHash, &u.Status, &u.Permissions,
	); err != nil {
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
	u.Hotel = new(hotel.Hotel)
	if err := p.db.QueryRow(selectUserByEmail, email).Scan(
		&u.ID, &u.FirstName, &u.SecondName, &u.LastName, &u.Hotel.Name, &u.Hotel.ID,
		&u.Note, &u.Phone, &u.PasswordHash, &u.Status, &u.Permissions,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, err
	}
	u.Email = email

	return u, nil
}

func (p *Postgres) List() ([]user.User, error) {
	rows, err := p.db.Query(selectUsers)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	uu := make([]user.User, 0)
	for rows.Next() {
		var u user.User
		u.Hotel = new(hotel.Hotel)

		if err := rows.Scan(
			&u.ID, &u.FirstName, &u.SecondName, &u.LastName, &u.Hotel.Name, &u.Hotel.ID,
			&u.Note, &u.Email, &u.Phone, &u.PasswordHash, &u.Status, &u.Permissions,
		); err != nil {
			return nil, errors.Wrap(err, "error scan user")
		}

		uu = append(uu, u)
	}

	return uu, nil
}

func (p *Postgres) Save(_ *user.User) error {
	return nil
}

func (p *Postgres) Delete(_ int64) error {
	return nil
}
