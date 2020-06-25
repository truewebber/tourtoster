package tour

import (
	"testing"

	"tourtoster/currency"
)

func TestPrice_ToUSD(t *testing.T) {
	currency.USD = 61.5

	tests := []struct {
		name string
		p    Price
		want int64
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
			if got := tt.p.ToUSD(); got != tt.want {
				t.Errorf("ToUSD() = %v, want %v", got, tt.want)
			}
		})
	}
}
