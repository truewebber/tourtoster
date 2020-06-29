package tour

import (
	"fmt"
	"strings"
)

type (
	Order struct {
		Field string
		Order string
	}
)

func NewOrder(field, order string) *Order {
	return &Order{
		Field: field,
		Order: filterOrder(order),
	}
}

func (o Order) Build() string {
	return fmt.Sprintf("%s %s", o.Field, o.Order)
}

func filterOrder(o string) string {
	o = strings.ToUpper(o)

	switch o {
	case "ASC", "DESC":
		return o
	default:
		return "ASC"
	}
}
