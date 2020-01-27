package mail

type (
	Mailer interface {
		Send(to string, body []byte) error
	}
)
