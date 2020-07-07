package repository

import (
	"database/sql"
	"strconv"
	"strings"

	"tourtoster/currency"
	"tourtoster/log"
)

type (
	sqlite struct {
		db     *sql.DB
		logger log.Logger
	}
)

const (
	inValues         = "inValues"
	selectCurrencies = `SELECT name, value FROM currencies WHERE name IN (` + inValues + `);`
)

func NewSQLite(db *sql.DB, logger log.Logger) *sqlite {
	return &sqlite{
		db:     db,
		logger: logger,
	}
}

func (s *sqlite) List(currencies ...currency.Currency) (map[currency.Currency]float64, error) {
	rows, err := s.db.Query(s.buildQuery(currencies), s.buildArgs(currencies)...)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			s.logger.Error("error close db rows", "error", err.Error())
		}
	}()

	out := make(map[currency.Currency]float64)

	for rows.Next() {
		var (
			name  currency.Currency
			value float64
		)

		if err := rows.Scan(&name, &value); err != nil {
			return nil, err
		}

		out[name] = value
	}

	return out, nil
}

func (s *sqlite) buildQuery(currencies []currency.Currency) string {
	values := make([]string, 0, len(currencies))
	for i := 0; i < len(currencies); i++ {
		values = append(values, "$"+strconv.Itoa(i+1))
	}

	return strings.ReplaceAll(selectCurrencies, inValues, strings.Join(values, ","))
}

func (s *sqlite) buildArgs(currencies []currency.Currency) []interface{} {
	args := make([]interface{}, 0, len(currencies))
	for i := 0; i < len(currencies); i++ {
		args = append(args, currencies[i])
	}

	return args
}
