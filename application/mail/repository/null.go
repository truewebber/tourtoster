package repository

import (
	"tourtoster/log"
)

type (
	null struct {
		logger log.Logger
	}
)

const (
	NullName = "Null"
)

func NewNull(logger log.Logger) *null {
	return &null{
		logger: logger,
	}
}

func (n *null) Name() string {
	return NullName
}

func (n *null) Send(to string, title, body string) error {
	n.logger.Info("Null mailer send email", "to", to, title, "title", "body", body)
	return nil
}
