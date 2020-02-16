package repository

import (
	"database/sql"
	"testing"
	"tourtoster/hotel"

	"github.com/pkg/errors"

	"tourtoster/user"
)

func Test_postgres_Save(t *testing.T) {
	db, err := newConn("/Users/truewebber/tourtoster/ttdb.sqlite")
	if err != nil {
		t.Error(err)

		return
	}

	type fields struct {
		db *sql.DB
	}
	type args struct {
		u *user.User
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:   "name1",
			fields: fields{db: db},
			args: args{&user.User{
				ID:           2,
				FirstName:    "Alex",
				SecondName:   "",
				LastName:     "Alex",
				Hotel:        &hotel.Hotel{ID: 4, Name: "Hotel"},
				Note:         "",
				Email:        "kish99@mail.ru",
				Phone:        "12345678912",
				Status:       3,
				Permissions:  6,
				PasswordHash: "hash",
				Token:        nil,
			}},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &postgres{
				db: tt.fields.db,
			}
			if err := p.Save(tt.args.u); (err != nil) != tt.wantErr {
				t.Errorf("Save() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func newConn(dbFilePath string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dbFilePath)
	if err != nil {
		return nil, errors.Wrap(err, "error connect to db")
	}
	if err := db.Ping(); err != nil {
		return nil, errors.Wrap(err, "error ping db")
	}

	return db, nil
}
