package user

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
