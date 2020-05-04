package repository

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
	"github.com/pkg/errors"

	"tourtoster/hotel"
)

type (
	sqlite struct {
		db *sql.DB
	}
)

const (
	insertHotel = "INSERT INTO hotel (name) VALUES ($1);"
	updateHotel = "UPDATE hotel SET name=$1 WHERE id=$2;"
	deleteHotel = "DELETE FROM hotel WHERE id=$1;"
	updateUsers = "UPDATE users SET hotel_id=0, hotel_name=(SELECT name FROM hotel WHERE id=$1) WHERE hotel_id=$1;"

	selectHotelByName = "SELECT id FROM hotel WHERE name=$1;"
	selectHotelByID   = "SELECT name FROM hotel WHERE id=$1;"
	selectHotels      = "SELECT id,name FROM hotel ORDER BY name COLLATE NOCASE;"
)

func NewSQLite(db *sql.DB) *sqlite {
	return &sqlite{
		db: db,
	}
}

func (p *sqlite) HotelByName(name string) (*hotel.Hotel, error) {
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

func (p *sqlite) Hotel(ID int64) (*hotel.Hotel, error) {
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

func (p *sqlite) List() ([]hotel.Hotel, error) {
	rows, err := p.db.Query(selectHotels)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = rows.Close()
	}()

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

func (p *sqlite) Save(h *hotel.Hotel) error {
	if h.ID == 0 {
		return p.insert(h)
	}

	return p.update(h)
}

func (p *sqlite) Delete(ID int64) error {
	tx, txErr := p.db.Begin()
	if txErr != nil {
		return errors.Wrap(txErr, "error create transaction")
	}
	defer func() {
		_ = tx.Rollback()
	}()

	if _, err := tx.Exec(updateUsers, ID); err != nil {
		return errors.Wrap(err, "error update users")
	}

	if _, err := tx.Exec(deleteHotel, ID); err != nil {
		return errors.Wrap(err, "error delete hotel")
	}

	if err := tx.Commit(); err != nil {
		return errors.Wrap(err, "error commit transaction")
	}

	return nil
}

func (p *sqlite) insert(h *hotel.Hotel) error {
	tx, txErr := p.db.Begin()
	if txErr != nil {
		return errors.Wrap(txErr, "error create transaction")
	}
	defer func() {
		_ = tx.Rollback()
	}()

	r, execErr := tx.Exec(insertHotel, h.Name)
	if execErr != nil {
		return errors.Wrap(execErr, "error insert hotel")
	}

	var err error
	h.ID, err = r.LastInsertId()
	if err != nil {
		return errors.Wrap(err, "error get last insert hotel ID")
	}

	if err := tx.Commit(); err != nil {
		return errors.Wrap(err, "error commit transaction")
	}

	return nil
}

func (p *sqlite) update(h *hotel.Hotel) error {
	if _, err := p.db.Exec(updateHotel, h.Name, h.ID); err != nil {
		return errors.Wrap(err, "error update hotel")
	}

	return nil
}
