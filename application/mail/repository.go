package mail

//go:generate mockgen -source=repository.go -destination=repository/mock.go -package=repository
type (
	Mailer interface {
		Send(to string, title, body string) error
		Name() string
	}
)
