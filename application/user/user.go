package user

import (
	"tourtoster/token"
)

type (
	User struct {
		ID           int64        `json:"id"`
		Email        string       `json:"email"`
		Role         Role         `json:"role"`
		Status       Status       `json:"status"`
		PasswordHash string       `json:"-"`
		Token        *token.Token `json:"-"`
	}
)
