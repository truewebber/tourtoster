package repository

import (
	"net/smtp"

	"github.com/pkg/errors"
)

type (
	gMail struct {
		auth smtp.Auth
		addr string
		from string
	}
)

func NewGMail(user, password, host, port string) *gMail {
	auth := smtp.PlainAuth(
		"",
		user,
		password,
		host,
	)

	return &gMail{
		auth: auth,
		addr: host + ":" + port,
		from: user,
	}
}

func (g *gMail) Name() string {
	return "GMail"
}

func (g *gMail) Send(to string, body []byte) error {
	if err := smtp.SendMail(g.addr, g.auth, g.from, []string{to}, body); err != nil {
		return errors.Wrap(err, "error send mail")
	}

	return nil
}
