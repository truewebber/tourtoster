package repository

import (
	"database/sql"
	"strings"

	"github.com/pkg/errors"

	"github.com/truewebber/tourtoster/log"
	"github.com/truewebber/tourtoster/tour"
	"github.com/truewebber/tourtoster/user"
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

	selectTours = `SELECT id,tour_type_id,creator_id,status,recurrence_rule,title,image,description,map,max_persons,
					price_per_children_3_6,price_per_children_0_6,price_per_children_7_17,price_per_adults,
					updated_at,created_at
				FROM tours` + whereField + orderField + `;`

	insertTour = `INSERT INTO tours (tour_type_id,creator_id,status,recurrence_rule,title,image,description,
									map,max_persons,price_per_children_3_6,price_per_children_0_6,
									price_per_children_7_17,price_per_adults,updated_at,created_at)
				VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,current_timestamp,current_timestamp);`

	updateTour = `UPDATE tours SET status=$1,recurrence_rule=$2,title=$3,image=$4,description=$5,map=$6,max_persons=$7,
								price_per_children_3_6=$8,price_per_children_0_6=$9,price_per_children_7_17=$10,
								price_per_adults=$11,updated_at=current_timestamp
					WHERE id=$12;`

	deleteTour = `DELETE FROM tours WHERE id=$1;`
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

		if err := rows.Scan(&t.ID, &t.Type, &t.Creator.ID, &t.Status, &t.Recurrence, &t.Title, &t.Image, &t.Description,
			&t.Map, &t.MaxPersons, &t.PricePerChildrenThreeSix, &t.PricePerChildrenZeroSix,
			&t.PricePerChildrenSevenSeventeen, &t.PricePerAdults, &t.UpdatedAt, &t.CreatedAt); err != nil {
			return nil, err
		}

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

	if err := s.db.QueryRow(query, args...).Scan(&t.ID, &t.Type, &t.Creator.ID, &t.Status, &t.Recurrence, &t.Title,
		&t.Image, &t.Description, &t.Map, &t.MaxPersons, &t.PricePerChildrenThreeSix, &t.PricePerChildrenZeroSix,
		&t.PricePerChildrenSevenSeventeen, &t.PricePerAdults, &t.UpdatedAt, &t.CreatedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, err
	}

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
	if t.ID == 0 {
		return s.insert(t)
	}

	return s.update(t)
}

func (s *sqlite) insert(t *tour.Tour) error {
	tx, txErr := s.db.Begin()
	if txErr != nil {
		return errors.Wrap(txErr, "error create `insert tour` transaction")
	}
	defer func() {
		_ = tx.Rollback()
	}()

	if _, err := tx.Exec(insertTour, t.Type, t.Creator.ID, t.Status, t.Recurrence.String(), t.Title, t.Image,
		t.Description, t.Map, t.MaxPersons, t.PricePerChildrenThreeSix, t.PricePerChildrenZeroSix,
		t.PricePerChildrenSevenSeventeen, t.PricePerAdults); err != nil {
		return errors.Wrapf(err, "error insert tour: %#v\n", t)
	}

	if err := tx.QueryRow("SELECT last_insert_rowid();").Scan(&t.ID); err != nil {
		return errors.Wrap(err, "error scan tour last_insert_rowid()")
	}

	err := tx.Commit()

	return errors.Wrap(err, "error commit tour insert transaction")
}

func (s *sqlite) update(t *tour.Tour) error {
	_, err := s.db.Exec(updateTour, t.Status, t.Recurrence.String(), t.Title, t.Image,
		t.Description, t.Map, t.MaxPersons, t.PricePerChildrenThreeSix, t.PricePerChildrenZeroSix,
		t.PricePerChildrenSevenSeventeen, t.PricePerAdults, t.ID)

	return errors.Wrapf(err, "error update tour: %#v\n", t)
}

func (s *sqlite) Delete(t *tour.Tour) error {
	_, err := s.db.Exec(deleteTour, t.ID)

	return errors.Wrapf(err, "error delete tour: %#v\n", t)
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
