package token

//go:generate mockgen -source=repository.go -destination=repository/mock.go -package=repository
type (
	Repository interface {
		Token(tkn string) (*Token, error)
		Save(tkn *Token) error
		Delete(tkn string) error
	}
)
