package req

import (
	"testing"
)

func TestGetOrderbook(t *testing.T) {
	tests := []struct {
		name   string
		symbol string
	}{
		{
			name:   "test",
			symbol: "BTC/USDT",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := Bitrue_get_orderbook(tt.symbol)
			t.Logf("res - %+v", res)
			if len(res.Asks) == 0 || len(res.Bids) == 0 {
				t.Errorf("Error: %+v", res)
			}
		})
	}
}
