package user

type (
	Role int
)

const (
	RoleSuperAdmin = 1
	RoleAdmin      = 2
	RoleEditor     = 3

	RoleSuperAdminName = "SuperAdmin"
	RoleAdminName      = "Admin"
	RoleEditorName     = "Editor"
)

func (r Role) String() string {
	switch r {
	case RoleSuperAdmin:
		return RoleSuperAdminName
	case RoleAdmin:
		return RoleAdminName
	case RoleEditor:
		return RoleEditorName
	default:
		return ""
	}
}
