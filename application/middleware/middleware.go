package middleware

import (
	"tourtoster/token"
	"tourtoster/user"
)

type (
	Middleware struct {
		token token.Repository
		user  user.Repository
	}
)

func New(token token.Repository, user user.Repository) *Middleware {
	return &Middleware{
		token: token,
		user:  user,
	}
}
