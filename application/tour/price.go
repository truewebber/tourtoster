package tour

import (
	"math"

	"github.com/truewebber/tourtoster/currency"
)

type (
	Price int
)

func NewRUB(p int) Price {
	return Price(p)
}

func (p Price) USD() int {
	return p.calc(currency.USD)
}

func (p Price) EUR() int {
	return p.calc(currency.EUR)
}

func (p Price) RUB() int {
	return int(p)
}

func (p Price) calc(currency float64) int {
	return int(math.Round(float64(p) / currency))
}
