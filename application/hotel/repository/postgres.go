package repository

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"

	"tourtoster/hotel"
)

type (
	Postgres struct {
		db *sql.DB
	}
)

const (
	selectHotelByID = `SELECT name FROM hotel WHERE id=$1;`
)

func NewPostgres(db *sql.DB) *Postgres {
	return &Postgres{
		db: db,
	}
}

func (p *Postgres) Hotel(ID int64) (*hotel.Hotel, error) {
	h := new(hotel.Hotel)

	if err := p.db.QueryRow(selectHotelByID, ID).Scan(&h.Name); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, err
	}
	h.ID = ID

	return h, nil
}

func (p *Postgres) Save(_ *hotel.Hotel) error {
	return nil
}

func (p *Postgres) Delete(_ int64) error {
	return nil
}
