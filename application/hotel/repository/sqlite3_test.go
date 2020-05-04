package repository

import (
	"database/sql"
	"testing"

	"github.com/pkg/errors"

	"tourtoster/hotel"
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
		h *hotel.Hotel
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:   "test1",
			fields: fields{db: db},
			args: args{&hotel.Hotel{
				ID:   4,
				Name: "Update1 Blahotel",
			}},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &sqlite{
				db: tt.fields.db,
			}
			if err := p.Save(tt.args.h); (err != nil) != tt.wantErr {
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
