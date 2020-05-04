package tour

//go:generate mockgen -source=repository.go -destination=repository/mock.go -package=repository
type (
	Repository interface {
		List(currencies ...Currency) (map[Currency]float64, error)
	}
)
