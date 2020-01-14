package middleware

import (
	"tourtoster/hotel"
	"tourtoster/token"
	"tourtoster/user"
)

type (
	Middleware struct {
		token token.Repository
		user  user.Repository
		hotel hotel.Repository
	}
)

func New(token token.Repository, user user.Repository, hotel hotel.Repository) *Middleware {
	return &Middleware{
		token: token,
		user:  user,
		hotel: hotel,
	}
}
