package repository

import (
	"database/sql"

	"tourtoster/group_tour"
)

type (
	sqlite struct {
		db *sql.DB
	}
)

const (
	selectTour = ``
)

func NewSQLite(db *sql.DB) *sqlite {
	return &sqlite{
		db: db,
	}
}

func (p *sqlite) Tour(ID int64) (*group_tour.Tour, error) {
	return nil, nil
}

func (p *sqlite) Save(t *group_tour.Tour) error {
	return nil
}

func (p *sqlite) Delete(ID int64) error {
	return nil
}

func (p *sqlite) List() ([]group_tour.Tour, error) {
	return nil, nil
}
