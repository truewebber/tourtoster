package tour

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/truewebber/tourtoster/user"
)

type (
	Filter struct {
		Field string
		Value interface{}
	}
)

const (
	tourTypeIDField = "tour_type_id"
	tourIDField     = "id"
	creatorIDField  = "creator_id"
	tourStatusField = "status"
)

func (f Filter) Build(index int) (string, []interface{}) {
	out := ""
	outArgs := make([]interface{}, 0, 1)

	s := reflect.ValueOf(f.Value)

	switch reflect.TypeOf(f.Value).Kind() {
	case reflect.Slice:
		if s.Len() == 0 {
			break
		}

		ss := make([]string, 0, 2)
		for i := 0; i < s.Len(); i++ {
			ss = append(ss, fmt.Sprintf("%s=$%d", f.Field, index+i))
			outArgs = append(outArgs, s.Index(i).Interface())
		}
		out = strings.Join(ss, " OR ")

		if s.Len() > 1 {
			out = "(" + out + ")"
		}
	default:
		out = fmt.Sprintf("%s=$%d", f.Field, index)
		outArgs = append(outArgs, s.Interface())
	}

	return out, outArgs
}

func FilterTourType(t Type) Filter {
	return Filter{
		Field: tourTypeIDField,
		Value: t,
	}
}

func FilterTourID(ID int64) Filter {
	return Filter{
		Field: tourIDField,
		Value: ID,
	}
}

func FilterTourCreator(u *user.User) Filter {
	return Filter{
		Field: creatorIDField,
		Value: u.ID,
	}
}

func FilterTourStatus(s ...Status) Filter {
	return Filter{
		Field: tourStatusField,
		Value: s,
	}
}
