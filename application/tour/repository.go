package tour

//go:generate mockgen -source=repository.go -destination=repository/mock.go -package=repository
type (
	Repository interface {
		Features() ([]Feature, error)

		List(*Order, ...Filter) ([]Tour, error)
		Tour(ID int64) (*Tour, error)
		Save(t *Tour) error
		Delete(ID int64) error
	}
)
