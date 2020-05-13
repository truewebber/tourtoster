package repository

import (
	"database/sql"

	"github.com/mgutz/logxi/v1"
	"github.com/pkg/errors"

	"tourtoster/group_tour"
	"tourtoster/tour"
	"tourtoster/user"
)

type (
	sqlite struct {
		db       *sql.DB
		userRepo user.Repository
	}
)

const (
	selectTours = `SELECT id, creator_id, title, image, description, map,
       				max_persons, price_per_adult, price_per_child
					FROM tours WHERE tour_type_id = $1 AND status = $2;`

	selectTourWithID = `SELECT creator_id, title, image, description, map,
       					max_persons, price_per_adult, price_per_child, status
						FROM tours WHERE id = $1;`
)

func NewSQLite(db *sql.DB, userRepo user.Repository) *sqlite {
	return &sqlite{
		db:       db,
		userRepo: userRepo,
	}
}

func (p *sqlite) List() ([]group_tour.Tour, error) {
	rows, err := p.db.Query(selectTours, tour.GroupType, tour.Enabled)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			log.Error("error close db rows", "error", err.Error())
		}
	}()

	out := make([]group_tour.Tour, 0)

	for rows.Next() {
		t := group_tour.Tour{
			Creator: &user.User{},
		}

		if err := rows.Scan(&t.ID, &t.Creator.ID, &t.Title, &t.Image, &t.Description,
			&t.Map, &t.MaxPersons, &t.PricePerAdult, &t.PricePerChild); err != nil {
			return nil, err
		}

		t.Status = tour.Enabled
		t.Type = tour.GroupType

		var err error
		t.Creator, err = p.userRepo.User(t.Creator.ID)
		if err != nil {
			return nil, err
		}

		out = append(out, t)
	}

	return out, nil
}

func (p *sqlite) Tour(ID int64) (*group_tour.Tour, error) {
	t := group_tour.Tour{
		Creator: &user.User{},
	}

	if err := p.db.QueryRow(selectTourWithID, ID).Scan(&t.Creator.ID, &t.Title, &t.Image, &t.Description,
		&t.Map, &t.MaxPersons, &t.PricePerAdult, &t.PricePerChild, &t.Status); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, err
	}

	t.ID = ID
	t.Type = tour.GroupType

	var err error
	t.Creator, err = p.userRepo.User(t.Creator.ID)
	if err != nil {
		return nil, err
	}

	if t.Creator == nil {
		return nil, errors.New("creator can't be nil")
	}

	return &t, nil
}

func (p *sqlite) Save(t *group_tour.Tour) error {
	return nil
}

func (p *sqlite) Delete(ID int64) error {
	return nil
}
