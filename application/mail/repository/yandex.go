package repository

import (
	"crypto/tls"
	"net/mail"
	"net/smtp"

	"github.com/pkg/errors"
)

type (
	yandex struct {
		tlsConfig *tls.Config
		auth      smtp.Auth
		from      mail.Address
		host      string
	}
)

const (
	yandexHost = "smtp.yandex.com"
	yandexPort = "465"
	YandexName = "Yandex"
)

func NewYandex(user, password string) *yandex {
	auth := smtp.PlainAuth(
		"",
		user,
		password,
		yandexHost,
	)

	tlsConfig := &tls.Config{
		InsecureSkipVerify: false,
		ServerName:         yandexHost,
	}

	return &yandex{
		tlsConfig: tlsConfig,
		auth:      auth,
		from:      mail.Address{Name: "Tourtoster", Address: user},
		host:      yandexHost,
	}
}

func (g *yandex) Name() string {
	return YandexName
}

func (g *yandex) Send(to string, title, body string) error {
	msg := "From: " + g.from.String() + "\n" +
		"To: " + to + "\n" +
		"Subject: " + title + "\n\n" +
		body

	conn, err := tls.Dial("tcp", yandexHost+":"+yandexPort, g.tlsConfig)
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
