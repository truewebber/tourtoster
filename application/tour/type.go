package tour

import (
	"errors"
)

type (
	Type   int
	Status int
)

const (
	GroupType Type = iota + 1
	PrivateType
)

const (
	Enabled Status = iota + 1
	Disabled
	Deleted
)

func ValidateType(t Type) error {
	switch t {
	case GroupType, PrivateType:
	default:
		return errors.New("invalid tour type")
	}

	return nil
}

func ValidateStatus(s Status) error {
	switch s {
	case Enabled, Disabled, Deleted:
	default:
		return errors.New("invalid tour status")
	}

	return nil
}
