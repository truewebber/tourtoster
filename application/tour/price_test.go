package tour

import (
	"testing"

	"github.com/truewebber/tourtoster/currency"
)

func TestPrice_USD(t *testing.T) {
	currency.USD = 61.5

	tests := []struct {
		name string
		p    Price
		want int
	}{
		{
			name: "test1",
			p:    Price(1000),
			want: 16,
		},
		{
			name: "test2",
			p:    Price(2323),
			want: 38,
		},
		{
			name: "test3",
			p:    Price(60),
			want: 1,
		},
		{
			name: "test4",
			p:    Price(30),
			want: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.p.USD(); got != tt.want {
				t.Errorf("USD() = %v, want %v", got, tt.want)
			}
		})
	}
}
