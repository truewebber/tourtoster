package hotel

import (
	"github.com/pkg/errors"
)

type (
	Hotel struct {
		ID   int64  `json:"id"`
		Name string `json:"name"`
	}
)

func ValidationFilterHotelID(ID int64) error {
	// 0 as all hotel_id is not set
	if ID >= -1 {
		return nil
	}

	return errors.New("invalid hotel id")
}
