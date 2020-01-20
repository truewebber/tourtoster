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
		Status       Status       `json:"-"`
		Permissions  Permission   `json:"-"`
		PasswordHash string       `json:"-"`
		Token        *token.Token `json:"-"`
	}
)

func (u *User) HasPermission(p Permission) bool {
	return u.Permissions&p > 0
}

func (u *User) AdminPage() bool {
	return u.HasPermission(CreateNewUserPermission) || u.HasPermission(EditToursPermission) || u.HasPermission(EditAllBookingsPermission)
}
