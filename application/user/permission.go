package user

import (
	"github.com/pkg/errors"
)

type (
	Permission int
)

const (
	CreateNewUserPermission    Permission = 2
	EditToursPermission        Permission = 4
	EditAllBookingsPermission  Permission = 8
	EditUserBookingsPermission Permission = 16

	CreateNewUserPermissionName    = "Create Users"
	EditToursPermissionName        = "Edit Tours"
	EditAllBookingsPermissionName  = "Edit All Bookings"
	EditUserBookingsPermissionName = "Edit User Bookings"
)

var (
	allowedPermissions = map[Permission]struct{}{
		CreateNewUserPermission:    {},
		EditToursPermission:        {},
		EditAllBookingsPermission:  {},
		EditUserBookingsPermission: {},
	}
)

func (p Permission) String() string {
	switch p {
	case CreateNewUserPermission:
		return CreateNewUserPermissionName
	case EditToursPermission:
		return EditToursPermissionName
	case EditAllBookingsPermission:
		return EditAllBookingsPermissionName
	case EditUserBookingsPermission:
		return EditUserBookingsPermissionName
	default:
		return ""
	}
}

func ValidationPermission(p Permission) error {
	if _, ok := allowedPermissions[p]; !ok {
		return errors.New("permission invalid")
	}

	return nil
}
