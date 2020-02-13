package repository

import (
	"net/smtp"

	"github.com/mgutz/logxi/v1"
)

type (
	null struct {
		auth smtp.Auth
		addr string
		from string
	}
)

func NewNull() *null {
	return &null{}
}

func (n *null) Name() string {
	return "Null"
}

func (n *null) Send(to string, body []byte) error {
	log.Debug("Null mailer send email", "to", to, "body", string(body))
	return nil
}
