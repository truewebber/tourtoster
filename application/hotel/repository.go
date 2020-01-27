package hotel

//go:generate mockgen -source=repository.go -destination=repository/mock.go -package=repository
type (
	Repository interface {
		Hotel(ID int64) (*Hotel, error)
		HotelByName(name string) (*Hotel, error)
		Save(h *Hotel) error
		Delete(ID int64) error
	}
)
