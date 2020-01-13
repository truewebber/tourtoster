package user

//go:generate mockgen -source=repository.go -destination=repository/mock.go -package=repository
type (
	Repository interface {
		User(ID int64) (*User, error)
		UserWithEmail(email string) (*User, error)
		Save(u *User) error
		Delete(ID int64) error
	}
)
