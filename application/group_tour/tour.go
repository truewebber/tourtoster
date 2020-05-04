package group_tour

import (
	"tourtoster/tour"
	"tourtoster/user"
)

type (
	Tour struct {
		ID            int64       `json:"id"`
		Type          tour.Type   `json:"type"`
		Creator       user.User   `json:"creator"`
		Status        tour.Status `json:"status"`
		Title         string      `json:"title"`
		Description   string      `json:"description"`
		Image         string      `json:"image"`
		Map           string      `json:"map"`
		MaxPersons    int         `json:"max_persons"`
		PricePerAdult tour.Price  `json:"price_per_adult"`
		PricePerChild tour.Price  `json:"price_per_child"`
	}
)
