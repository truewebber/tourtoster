package repository

import (
	"database/sql"
	"strings"

	"github.com/pkg/errors"

	"tourtoster/log"
	"tourtoster/tour"
	"tourtoster/user"
)

type (
	sqlite struct {
		db       *sql.DB
		userRepo user.Repository
		logger   log.Logger
	}
)

const (
	whereField = "{WHERE}"
	orderField = "{ORDER}"
)

const (
	selectFeatures = `SELECT id, tour_type_id, icon, title FROM features;`

	selectTours = `SELECT id,tour_type_id,creator_id,status,recurrence_rule,title,image,description,map,max_persons,price_per_children_3_6,price_per_children_0_6,price_per_children_7_17,price_per_adults,updated_at,created_at FROM tours` + whereField + orderField + `;`
)

func NewSQLite(db *sql.DB, userRepo user.Repository, logger log.Logger) *sqlite {
	return &sqlite{
		db:       db,
		userRepo: userRepo,
		logger:   logger,
	}
}

func (s *sqlite) buildListQuery(o *tour.Order, ff []tour.Filter) (string, []interface{}) {
	replaceWhere, args := s.buildWhere(ff)
	replaceOrder := s.buildOrder(o)

	query := strings.ReplaceAll(selectTours, whereField, replaceWhere)
	query = strings.ReplaceAll(query, orderField, replaceOrder)

	return query, args
}

func (s *sqlite) buildWhere(ff []tour.Filter) (string, []interface{}) {
	out := ""
	outArgs := make([]interface{}, 0, len(ff))

	ss := make([]string, 0)
	for i := 0; i < len(ff); i++ {
		s, sArgs := ff[i].Build(len(outArgs) + 1)
		ss = append(ss, s)
		outArgs = append(outArgs, sArgs...)
	}

	out = strings.Join(ss, " AND ")
	if out != "" {
		out = " WHERE " + out
	}

	return out, outArgs
}

func (s *sqlite) buildOrder(o *tour.Order) string {
	if o == nil {
		return ""
	}

	return " ORDER BY " + o.Build()
}

func (s *sqlite) List(o *tour.Order, ff ...tour.Filter) ([]tour.Tour, error) {
	query, args := s.buildListQuery(o, ff)

	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			s.logger.Error("error close db rows", "error", err.Error())
		}
	}()

	out := make([]tour.Tour, 0)

	for rows.Next() {
		t := tour.Tour{
			Creator: &user.User{},
		}

		if err := rows.Scan(&t.ID, &t.Creator.ID, &t.Title, &t.Image, &t.Description,
			&t.Map, &t.MaxPersons, &t.PricePerAdults, &t.PricePerChildrenSevenSeventeen); err != nil {
			return nil, err
		}

		t.Status = tour.Enabled
		t.Type = tour.GroupType

		var err error
		t.Creator, err = s.userRepo.User(t.Creator.ID)
		if err != nil {
			return nil, err
		}

		out = append(out, t)
	}

	return out, nil
}

func (s *sqlite) Tour(ID int64) (*tour.Tour, error) {
	t := tour.Tour{
		Creator: &user.User{},
	}

	query, args := s.buildListQuery(nil, []tour.Filter{tour.FilterTourID(ID)})

	if err := s.db.QueryRow(query, args...).Scan(&t.Creator.ID, &t.Title, &t.Image, &t.Description,
		&t.Map, &t.MaxPersons, &t.PricePerAdults, &t.PricePerChildrenSevenSeventeen, &t.Status); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, err
	}

	t.ID = ID
	t.Type = tour.GroupType

	var err error
	t.Creator, err = s.userRepo.User(t.Creator.ID)
	if err != nil {
		return nil, err
	}

	if t.Creator == nil {
		return nil, errors.New("creator can't be nil")
	}

	return &t, nil
}

func (s *sqlite) Save(t *tour.Tour) error {
	return nil
}

func (s *sqlite) Delete(ID int64) error {
	return nil
}

// ---------------------------------------------------------------------------------------------------------------------

func (s *sqlite) Features() ([]tour.Feature, error) {
	rows, err := s.db.Query(selectFeatures)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			s.logger.Error("error close db rows", "error", err.Error())
		}
	}()

	out := make([]tour.Feature, 0)

	for rows.Next() {
		f := tour.Feature{}

		if err := rows.Scan(&f.ID, &f.TourType, &f.Icon, &f.Title); err != nil {
			return nil, err
		}

		out = append(out, f)
	}

	return out, nil
}
