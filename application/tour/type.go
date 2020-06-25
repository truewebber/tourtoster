package tour

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
