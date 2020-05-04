package tour

import (
	"math"
)

type (
	Price int64
)

func (p Price) ToUSD() int64 {
	return p.calc(USD)
}

func (p Price) ToEUR() int64 {
	return p.calc(EUR)
}

func (p Price) calc(currency float64) int64 {
	return int64(math.Round(float64(p) / currency))
}
