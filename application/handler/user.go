package handler

import (
	"tourtoster/user"
)

type (
	htmlUser struct {
		FirstName  string
		SecondName string
		LastName   string
		Allow      allow
	}

	allow struct {
		CreateNewUser   bool
		EditTours       bool
		EditAllBookings bool
	}
)

func templateUser(u *user.User) *htmlUser {
	return &htmlUser{
		FirstName:  u.FirstName,
		SecondName: u.SecondName,
		LastName:   u.LastName,
		Allow: allow{
			CreateNewUser:   u.HasPermission(user.CreateNewUserPermission),
			EditTours:       u.HasPermission(user.EditToursPermission),
			EditAllBookings: u.HasPermission(user.EditAllBookingsPermission),
		},
	}
}
