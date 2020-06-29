package tour

import (
	"reflect"
	"testing"
)

func TestFilter_Build(t *testing.T) {
	type fields struct {
		Field string
		Value interface{}
	}
	type args struct {
		index int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
		want1  []interface{}
	}{
		{
			name:   "test1",
			fields: fields{Field: "field", Value: []Status{1, 2}},
			args:   args{index: 1},
			want:   "(field=$1 OR field=$2)",
			want1:  []interface{}{Status(1), Status(2)},
		},
		{
			name:   "test2",
			fields: fields{Field: "field", Value: int64(1010)},
			args:   args{index: 1},
			want:   "field=$1",
			want1:  []interface{}{int64(1010)},
		},
		{
			name:   "test3",
			fields: fields{Field: "field", Value: "tosibosi"},
			args:   args{index: 6},
			want:   "field=$6",
			want1:  []interface{}{"tosibosi"},
		},
		{
			name:   "test4",
			fields: fields{Field: "field", Value: []string{}},
			args:   args{index: 6},
			want:   "",
			want1:  []interface{}{},
		},
		{
			name:   "test5",
			fields: fields{Field: "field", Value: []string{"tosi"}},
			args:   args{index: 6},
			want:   "field=$6",
			want1:  []interface{}{"tosi"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := Filter{
				Field: tt.fields.Field,
				Value: tt.fields.Value,
			}
			got, got1 := f.Build(tt.args.index)
			if got != tt.want {
				t.Errorf("Build() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("Build() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
