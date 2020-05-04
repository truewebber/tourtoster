package user

import (
	"github.com/pkg/errors"
)

type (
	Status int
)

const (
	StatusNew Status = iota + 1
	StatusEnabled
	StatusDisabled

	StatusNewName      = "new"
	StatusEnabledName  = "enabled"
	StatusDisabledName = "disabled"
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
	switch s {
	case StatusNew, StatusEnabled, StatusDisabled:
		return nil
	}

	return errors.New("status invalid")
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
