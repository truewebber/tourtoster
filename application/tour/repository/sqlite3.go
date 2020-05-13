package repository

import (
	"database/sql"
	"strconv"
	"strings"

	"github.com/mgutz/logxi/v1"

	"tourtoster/tour"
)

type (
	sqlite struct {
		db *sql.DB
	}
)

const (
	inValues         = "inValues"
	selectCurrencies = `SELECT name, value FROM currencies WHERE name IN (` + inValues + `);`
)

func NewSQLite(db *sql.DB) *sqlite {
	return &sqlite{
		db: db,
	}
}

func (s *sqlite) List(currencies ...tour.Currency) (map[tour.Currency]float64, error) {
	rows, err := s.db.Query(s.buildQuery(currencies), s.buildArgs(currencies)...)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			log.Error("error close db rows", "error", err.Error())
		}
	}()

	out := make(map[tour.Currency]float64)

	for rows.Next() {
		var (
			name  tour.Currency
			value float64
		)

		if err := rows.Scan(&name, &value); err != nil {
			return nil, err
		}

		out[name] = value
	}

	return out, nil
}

func (s *sqlite) buildQuery(currencies []tour.Currency) string {
	values := make([]string, 0, len(currencies))
	for i := 0; i < len(currencies); i++ {
		values = append(values, "$"+strconv.Itoa(i+1))
	}

	return strings.ReplaceAll(selectCurrencies, inValues, strings.Join(values, ","))
}

func (s *sqlite) buildArgs(currencies []tour.Currency) []interface{} {
	args := make([]interface{}, 0, len(currencies))
	for i := 0; i < len(currencies); i++ {
		args = append(args, currencies[i])
	}

	return args
}
