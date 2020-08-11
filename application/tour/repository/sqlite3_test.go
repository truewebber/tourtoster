package repository

import (
	"database/sql"
	"reflect"
	"strings"
	"testing"

	"github.com/truewebber/tourtoster/tour"
	"github.com/truewebber/tourtoster/user"
)

func Test_sqlite_buildListQuery(t *testing.T) {
	type (
		fields struct {
			db       *sql.DB
			userRepo user.Repository
		}
		args struct {
			o  *tour.Order
			ff []tour.Filter
		}
	)

	defaultFields := fields{db: nil, userRepo: nil}

	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
		want1  []interface{}
	}{
		{
			name:   "test1",
			fields: defaultFields,
			args:   args{o: nil, ff: nil},
			want:   "SELECT id,tour_type_id,creator_id,status,recurrence_rule,title,image,description,map,max_persons,price_per_children_3_6,price_per_children_0_6,price_per_children_7_17,price_per_adults,updated_at,created_at FROM tours;",
			want1:  []interface{}{},
		},
		{
			name:   "test2",
			fields: defaultFields,
			args:   args{o: tour.NewOrder("tosi", "asc"), ff: nil},
			want:   "SELECT id,tour_type_id,creator_id,status,recurrence_rule,title,image,description,map,max_persons,price_per_children_3_6,price_per_children_0_6,price_per_children_7_17,price_per_adults,updated_at,created_at FROM tours ORDER BY tosi ASC;",
			want1:  []interface{}{},
		},
		{
			name:   "test3",
			fields: defaultFields,
			args:   args{o: nil, ff: []tour.Filter{{Field: "tosi", Value: "bosi"}}},
			want:   "SELECT id,tour_type_id,creator_id,status,recurrence_rule,title,image,description,map,max_persons,price_per_children_3_6,price_per_children_0_6,price_per_children_7_17,price_per_adults,updated_at,created_at FROM tours WHERE tosi=$1;",
			want1:  []interface{}{"bosi"},
		},
		{
			name:   "test4",
			fields: defaultFields,
			args:   args{o: nil, ff: []tour.Filter{{Field: "tosi", Value: "bosi"}, {Field: "foo", Value: "bar"}}},
			want:   "SELECT id,tour_type_id,creator_id,status,recurrence_rule,title,image,description,map,max_persons,price_per_children_3_6,price_per_children_0_6,price_per_children_7_17,price_per_adults,updated_at,created_at FROM tours WHERE tosi=$1 AND foo=$2;",
			want1:  []interface{}{"bosi", "bar"},
		},
		{
			name:   "test5",
			fields: defaultFields,
			args:   args{o: nil, ff: []tour.Filter{{Field: "tosi", Value: []string{"bosi", "sisi"}}, {Field: "foo", Value: "bar"}}},
			want:   "SELECT id,tour_type_id,creator_id,status,recurrence_rule,title,image,description,map,max_persons,price_per_children_3_6,price_per_children_0_6,price_per_children_7_17,price_per_adults,updated_at,created_at FROM tours WHERE (tosi=$1 OR tosi=$2) AND foo=$3;",
			want1:  []interface{}{"bosi", "sisi", "bar"},
		},
		{
			name:   "test6",
			fields: defaultFields,
			args:   args{o: tour.NewOrder("tosi", "desc"), ff: []tour.Filter{tour.FilterTourStatus(tour.Enabled, tour.Disabled)}},
			want:   "SELECT id,tour_type_id,creator_id,status,recurrence_rule,title,image,description,map,max_persons,price_per_children_3_6,price_per_children_0_6,price_per_children_7_17,price_per_adults,updated_at,created_at FROM tours WHERE (status=$1 OR status=$2) ORDER BY tosi DESC;",
			want1:  []interface{}{tour.Status(1), tour.Status(2)},
		},
		{
			name:   "test7",
			fields: defaultFields,
			args:   args{o: nil, ff: []tour.Filter{tour.FilterTourStatus(tour.Deleted)}},
			want:   "SELECT id,tour_type_id,creator_id,status,recurrence_rule,title,image,description,map,max_persons,price_per_children_3_6,price_per_children_0_6,price_per_children_7_17,price_per_adults,updated_at,created_at FROM tours WHERE status=$1;",
			want1:  []interface{}{tour.Status(3)},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &sqlite{
				db:       tt.fields.db,
				userRepo: tt.fields.userRepo,
			}
			got, got1 := s.buildListQuery(tt.args.o, tt.args.ff)

			if trimQuery(got) != trimQuery(tt.want) {
				t.Errorf("buildListQuery() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("buildListQuery() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func trimQuery(s string) string {
	s = strings.ReplaceAll(s, "\n", "")
	s = strings.ReplaceAll(s, "\t", " ")
	s = strings.ReplaceAll(s, " ", "")

	return s
}
