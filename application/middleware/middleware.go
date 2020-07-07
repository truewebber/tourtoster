package middleware

import (
	"tourtoster/hotel"
	"tourtoster/log"
	"tourtoster/token"
	"tourtoster/user"
)

type (
	Middleware struct {
		token  token.Repository
		user   user.Repository
		hotel  hotel.Repository
		logger log.Logger
	}
)

func New(token token.Repository, user user.Repository, hotel hotel.Repository, logger log.Logger) *Middleware {
	return &Middleware{
		token:  token,
		user:   user,
		hotel:  hotel,
		logger: logger,
	}
}
