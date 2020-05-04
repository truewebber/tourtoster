package repository

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"

	"github.com/mattn/go-sqlite3"
	"github.com/pkg/errors"

	"tourtoster/hotel"
	"tourtoster/user"
)

type (
	sqlite struct {
		db *sql.DB
	}
)

const (
	selectUsers = `SELECT id, first_name, second_name, last_name, hotel_name, hotel_id,
       						note, email, phone, password_hash, status, role
					FROM users ` + field + ` ORDER BY last_name COLLATE NOCASE, first_name COLLATE NOCASE;`
	field = `{FIELD}`

	selectUserByID = `SELECT first_name, second_name, last_name, hotel_name, hotel_id,
       						note, email, phone, password_hash, status, role
					FROM users WHERE id=$1;`
	selectUserByEmail = `SELECT id, first_name, second_name, last_name, hotel_name, hotel_id,
       						note, phone, password_hash, status, role
						FROM users WHERE email=$1;`
	insertUser = `INSERT INTO users (first_name, second_name, last_name, hotel_name, 
									hotel_id, note, email, phone, password_hash, status, role)
					VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11);`
	updatePasswordUser = `UPDATE users SET password_hash=$1 WHERE id=$2;`
	deleteUser         = `DELETE FROM users WHERE id=$1;`
)

func NewSQLite(db *sql.DB) *sqlite {
	return &sqlite{
		db: db,
	}
}

func (p *sqlite) User(ID int64) (*user.User, error) {
	u := new(user.User)
	u.Hotel = new(hotel.Hotel)

	if err := p.db.QueryRow(selectUserByID, ID).Scan(
		&u.FirstName, &u.SecondName, &u.LastName, &u.Hotel.Name, &u.Hotel.ID,
		&u.Note, &u.Email, &u.Phone, &u.PasswordHash, &u.Status, &u.Permissions,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, err
	}
	u.ID = ID

	return u, nil
}

func (p *sqlite) UserWithEmail(email string) (*user.User, error) {
	u := new(user.User)
	u.Hotel = new(hotel.Hotel)
	if err := p.db.QueryRow(selectUserByEmail, email).Scan(
		&u.ID, &u.FirstName, &u.SecondName, &u.LastName, &u.Hotel.Name, &u.Hotel.ID,
		&u.Note, &u.Phone, &u.PasswordHash, &u.Status, &u.Permissions,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, err
	}
	u.Email = email

	return u, nil
}

func (p *sqlite) List(filters map[string]interface{}) ([]user.User, error) {
	query, params := buildListQueryParams(selectUsers, filters)

	rows, err := p.db.Query(query, params...)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	uu := make([]user.User, 0)
	for rows.Next() {
		var u user.User
		u.Hotel = new(hotel.Hotel)

		if err := rows.Scan(
			&u.ID, &u.FirstName, &u.SecondName, &u.LastName, &u.Hotel.Name, &u.Hotel.ID,
			&u.Note, &u.Email, &u.Phone, &u.PasswordHash, &u.Status, &u.Permissions,
		); err != nil {
			return nil, errors.Wrap(err, "error scan user")
		}

		uu = append(uu, u)
	}

	return uu, nil
}

func buildListQueryParams(query string, filters map[string]interface{}) (string, []interface{}) {
	params := make([]interface{}, 0, 2)
	where := make([]string, 0)

	for key, val := range filters {
		switch key {
		case "status":
			if val.(user.Status) == user.Status(-1) {
				continue
			}

			where = append(where, fmt.Sprintf("status=$%d", len(where)+1))
			params = append(params, val)
		case "hotel":
			if val.(int64) == int64(-1) {
				continue
			}

			where = append(where, fmt.Sprintf("hotel_id=$%d", len(where)+1))
			params = append(params, val)
		}
	}

	r := ""
	if len(params) > 0 {
		r = "WHERE " + strings.Join(where, " AND ")
	}

	query = strings.Replace(query, field, r, 1)

	return query, params
}

func (p *sqlite) Password(ID int64, passwordHash string) error {
	_, err := p.db.Exec(updatePasswordUser, passwordHash, ID)
	if err != nil {
		return errors.Wrap(err, "error update user")
	}

	return nil
}

func (p *sqlite) Save(u *user.User) error {
	if u.ID == 0 {
		return p.insert(u)
	}

	return p.update(u)
}

func (p *sqlite) Delete(ID int64) error {
	_, err := p.db.Exec(deleteUser, ID)
	if err != nil {
		return errors.Wrap(err, "error delete user")
	}

	return nil
}

func (p *sqlite) insert(u *user.User) error {
	tx, txErr := p.db.Begin()
	if txErr != nil {
		return errors.Wrap(txErr, "error create tx insert user")
	}
	defer func() { _ = tx.Rollback() }()

	result, err := tx.Exec(insertUser, u.FirstName, u.SecondName, u.LastName, insertHotelName(u.Hotel), u.Hotel.ID,
		u.Note, u.Email, u.Phone, u.PasswordHash, u.Status, u.Permissions)
	if err != nil {
		sqliteErr := err.(sqlite3.Error)
		if sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {
			return errors.Wrap(user.PhoneEmailUniqueError, "error insert user")
		}

		return errors.Wrap(err, "error insert user")
	}

	var IDErr error
	u.ID, IDErr = result.LastInsertId()
	if IDErr != nil {
		return errors.Wrap(IDErr, "error get last insert user id")
	}

	if err := tx.Commit(); err != nil {
		return errors.Wrap(err, "error commit tx insert user")
	}

	return nil
}

func (p *sqlite) update(u *user.User) error {
	q := buildUpdateQuery(u.PasswordHash != "")
	params := buildUpdateParams(u)

	if _, err := p.db.Exec(q, params...); err != nil {
		return errors.Wrap(err, "error update user")
	}

	return nil
}

func buildUpdateQuery(password bool) string {
	params := []string{
		"first_name", "second_name", "last_name", "hotel_name",
		"hotel_id", "note", "email", "phone", "status", "role",
	}
	if password {
		params = append(params, "password_hash")
	}

	query := "UPDATE users SET "
	for i := 1; i <= len(params); i++ {
		query += params[i-1] + "=$" + strconv.Itoa(i) + ", "
	}
	query += "updated_at=CURRENT_TIMESTAMP WHERE id=$" + strconv.Itoa(len(params)+1)

	return query
}

func buildUpdateParams(u *user.User) []interface{} {
	params := []interface{}{
		u.FirstName, u.SecondName, u.LastName, insertHotelName(u.Hotel),
		u.Hotel.ID, u.Note, u.Email, u.Phone, u.Status, u.Permissions,
	}
	if u.PasswordHash != "" {
		params = append(params, u.PasswordHash)
	}
	params = append(params, u.ID)

	return params
}

func insertHotelName(h *hotel.Hotel) string {
	if h.ID != 0 {
		return ""
	}

	return h.Name

}
