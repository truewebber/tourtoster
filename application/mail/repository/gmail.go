package repository

import (
	"crypto/tls"
	"net/mail"
	"net/smtp"

	"github.com/pkg/errors"
)

type (
	gMail struct {
		tlsConfig *tls.Config
		auth      smtp.Auth
		from      mail.Address
		host      string
	}
)

const (
	gMailHost = "smtp.gmail.com"
	gMailPort = "465"
	GMailName = "GMail"
)

func NewGMail(user, password string) *gMail {
	auth := smtp.PlainAuth(
		"",
		user,
		password,
		gMailHost,
	)

	tlsConfig := &tls.Config{
		InsecureSkipVerify: false,
		ServerName:         gMailHost,
	}

	return &gMail{
		tlsConfig: tlsConfig,
		auth:      auth,
		from:      mail.Address{Name: "Tourtoster", Address: user},
		host:      gMailHost,
	}
}

func (g *gMail) Name() string {
	return GMailName
}

func (g *gMail) Send(to string, title, body string) error {
	msg := "From: " + g.from.String() + "\n" +
		"To: " + to + "\n" +
		"Subject: " + title + "\n\n" +
		body

	conn, err := tls.Dial("tcp", gMailHost+":"+gMailPort, g.tlsConfig)
	if err != nil {
		return errors.Wrap(err, "error dial")
	}

	client, clientErr := smtp.NewClient(conn, g.host)
	if clientErr != nil {
		return errors.Wrap(clientErr, "error create new client with conn")
	}

	if err := client.Auth(g.auth); err != nil {
		return errors.Wrap(err, "error set client auth")
	}

	if err := client.Mail(g.from.Address); err != nil {
		return errors.Wrap(err, "error set from")
	}

	if err := client.Rcpt(to); err != nil {
		return errors.Wrap(err, "error set recipient")
	}

	w, err := client.Data()
	if err != nil {
		return errors.Wrap(err, "error get client writer")
	}

	if _, err := w.Write([]byte(msg)); err != nil {
		return errors.Wrap(err, "error write message")
	}

	if err := w.Close(); err != nil {
		return errors.Wrap(err, "error close writer")
	}

	if err := client.Quit(); err != nil {
		return errors.Wrap(err, "error quit client")
	}

	return nil
}
