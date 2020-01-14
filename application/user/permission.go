package user

type (
	Permission int
)

const (
	CreateNewUserPermission      Permission = 2
	EditToursPermission          Permission = 4
	EditAllNewBookingsPermission Permission = 8
	EditUserBookingsPermission   Permission = 16

	CreateNewUserPermissionName      = "Create Users"
	EditToursPermissionName          = "Edit Tours"
	EditAllNewBookingsPermissionName = "Edit All Bookings"
	EditUserBookingsPermissionName   = "Edit User Bookings"
)

func (p Permission) String() string {
	switch p {
	case CreateNewUserPermission:
		return CreateNewUserPermissionName
	case EditToursPermission:
		return EditToursPermissionName
	case EditAllNewBookingsPermission:
		return EditAllNewBookingsPermissionName
	case EditUserBookingsPermission:
		return EditUserBookingsPermissionName
	default:
		return ""
	}
}
