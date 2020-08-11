package handler

import (
	"github.com/truewebber/tourtoster/user"
)

type (
	me struct {
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

func templateMe(u *user.User) *me {
	return &me{
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
