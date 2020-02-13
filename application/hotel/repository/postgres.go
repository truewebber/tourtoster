package repository

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
	"github.com/pkg/errors"

	"tourtoster/hotel"
)

type (
	Postgres struct {
		db *sql.DB
	}
)

const (
	insertHotel     = "INSERT INTO hotel (id, name) VALUES ($1, $2);"
	updateHotel     = "UPDATE hotel SET name=$1 WHERE id=$2;"
	deleteHotelByID = "DELETE FROM hotel WHERE id=$1;"

	selectHotelByName = "SELECT id FROM hotel WHERE name=$1;"
	selectHotelByID   = "SELECT name FROM hotel WHERE id=$1;"
	selectHotels      = "SELECT id,name FROM hotel ORDER BY name;"
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

func (p *Postgres) List() ([]hotel.Hotel, error) {
	rows, err := p.db.Query(selectHotels)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	hh := make([]hotel.Hotel, 0)
	for rows.Next() {
		h := hotel.Hotel{}
		if err := rows.Scan(&h.ID, &h.Name); err != nil {
			return nil, errors.Wrap(err, "error scan hotel")
		}

		hh = append(hh, h)
	}

	return hh, nil
}

func (p *Postgres) Save(h *hotel.Hotel) error {
	query := insertHotel
	args := []interface{}{h.Name}

	if h.ID != 0 {
		query = updateHotel
		args = []interface{}{h.Name, h.ID}
	}

	_, err := p.db.Exec(query, args)
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
