package user

import (
	"tourtoster/hotel"
	"tourtoster/token"
)

type (
	User struct {
		ID           int64        `json:"id"`
		FirstName    string       `json:"first_name"`
		SecondName   string       `json:"second_name"`
		LastName     string       `json:"last_name"`
		Hotel        *hotel.Hotel `json:"hotel"`
		Note         string       `json:"note"`
		Email        string       `json:"email"`
		Phone        string       `json:"phone"`
		Status       Status       `json:"status"`
		Role         Permission   `json:"role"`
		PasswordHash string       `json:"-"`
		Token        *token.Token `json:"-"`
	}
)
