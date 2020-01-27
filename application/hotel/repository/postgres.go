package repository

import (
	"database/sql"
	"github.com/pkg/errors"

	_ "github.com/mattn/go-sqlite3"

	"tourtoster/hotel"
)

type (
	Postgres struct {
		db *sql.DB
	}
)

const (
	selectHotelByName = `SELECT id FROM hotel WHERE name=$1;`
	selectHotelByID   = `SELECT name FROM hotel WHERE id=$1;`
	insertHotel       = `INSERT INTO hotel (id, name) VALUES ($1, $2);`
	deleteHotelByID   = `DELETE FROM hotel WHERE id=$1;`
)

func NewPostgres(db *sql.DB) *Postgres {
	return &Postgres{
		db: db,
	}
}

func (p *Postgres) HotelByName(name string) (*hotel.Hotel, error) {
	h := new(hotel.Hotel)

	if err := p.db.QueryRow(selectHotelByName, name).Scan(&h.ID); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, err
	}
	h.Name = name

	return h, nil
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

func (p *Postgres) Save(h *hotel.Hotel) error {
	_, err := p.db.Exec(insertHotel, h.ID, h.Name)
	if err != nil {
		return errors.Wrap(err, "error insert hotel")
	}

	return nil
}

func (p *Postgres) Delete(ID int64) error {
	_, err := p.db.Exec(deleteHotelByID, ID)
	if err != nil {
		return errors.Wrap(err, "error delete hotel")
	}

	return nil
}
