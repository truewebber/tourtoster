package user

import (
	"github.com/pkg/errors"
)

var (
	PhoneEmailUniqueError = errors.New("phone or e-mail is not unique")
)

//go:generate mockgen -source=repository.go -destination=repository/mock.go -package=repository
type (
	Repository interface {
		User(ID int64) (*User, error)
		UserWithEmail(email string) (*User, error)
		Save(u *User) error
		Password(ID int64, passwordHash string) error
		Delete(ID int64) error
		List() ([]User, error)
	}
)
