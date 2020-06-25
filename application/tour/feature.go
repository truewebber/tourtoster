package tour

var (
	FeaturesByType = make(map[Type][]Feature)
	FeaturesByID   = make(map[int]Feature)
)

type (
	Feature struct {
		ID       int    `json:"-"`
		TourType Type   `json:"-"`
		Icon     string `json:"icon"`
		Title    string `json:"title"`
	}
)
