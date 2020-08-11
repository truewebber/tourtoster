package middleware

import (
	"github.com/truewebber/tourtoster/hotel"
	"github.com/truewebber/tourtoster/log"
	"github.com/truewebber/tourtoster/token"
	"github.com/truewebber/tourtoster/user"
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
