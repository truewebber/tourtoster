package repository

import (
	"github.com/mgutz/logxi/v1"
)

type (
	null struct{}
)

const (
	NullName = "Null"
)

func NewNull() *null {
	return &null{}
}

func (n *null) Name() string {
	return NullName
}

func (n *null) Send(to string, title, body string) error {
	log.Debug("Null mailer send email", "to", to, title, "title", "body", body)
	return nil
}
