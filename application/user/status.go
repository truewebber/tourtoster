package user

import (
	"github.com/pkg/errors"
)

type (
	Status int
)

const (
	StatusNew      Status = 1
	StatusEnabled  Status = 2
	StatusDisabled Status = 3

	StatusNewName      = "new"
	StatusEnabledName  = "enabled"
	StatusDisabledName = "disabled"
)

var (
	allowedStatuses = map[Status]struct{}{
		StatusNew:      {},
		StatusEnabled:  {},
		StatusDisabled: {},
	}
)

func (s Status) String() string {
	switch s {
	case StatusNew:
		return StatusNewName
	case StatusEnabled:
		return StatusEnabledName
	case StatusDisabled:
		return StatusDisabledName
	}

	return ""
}

func ValidationStatus(s Status) error {
	if _, ok := allowedStatuses[s]; !ok {
		return errors.New("status invalid")
	}

	return nil
}

func ValidationFilterStatus(s Status) error {
	if err := ValidationStatus(s); err != nil {
		if s == Status(-1) {
			return nil
		}

		return err
	}

	return nil
}
