package finder

import (
	"math"
	"testing"

	"mad-scanner.com/scriner/common"
)

func TestCalcMaxVolume(t *testing.T) {
	tests := []struct {
		name         string
		asks         []common.Order
		bids         []common.Order
		thenMakerFee float64
		withdrawFee  float64
		expected     common.VolumeData
	}{
		// BID < ASK
		// 5 5.1
		// 4 6
		// 3 8
		// BID всегда падает
		// ASK всегда растёт
		// Арбитраж когда наоборот BID > ASK

		// ответы на тесты были посчитаны не вручную
		// взяты результаты алгоритма, который был наиболее объективен

		// Positive spread with enough depth
		{
			name: "Positive spread with enough depth",
			asks: []common.Order{
				{Price: 100, Amount: 1},
				{Price: 101, Amount: 1},
			},
			bids: []common.Order{
				{Price: 105, Amount: 1},
				{Price: 104, Amount: 1},
			},
			thenMakerFee: 0,
			withdrawFee:  0,
			expected: common.VolumeData{
				Up: common.UpperLimit{
					AskVolUsdt: 201,
					BidVolUsdt: 209,
					BuyPrice:   100.5,
					SellPrice:  104.5,
					Profit:     8,
					AskDepth:   2,
					BidDepth:   2,
					Spread:     3.9800995024875623,
				},
				Low: common.LowerLimit{
					AskVolUsdt: 100,
					BidVolUsdt: 105,
					BuyPrice:   100,
					SellPrice:  105,
					Profit:     5,
					AskDepth:   1,
					BidDepth:   1,
					Spread:     5,
				},
			},
		},
		// No profit
		{
			name: "No profit",
			asks: []common.Order{
				{Price: 100, Amount: 50},
				{Price: 100, Amount: 50},
			},
			bids: []common.Order{
				{Price: 100, Amount: 50},
				{Price: 100, Amount: 50},
			},
			thenMakerFee: 0.1,
			withdrawFee:  0.1,
			expected: common.VolumeData{
				Up: common.UpperLimit{
					// all 0
				},
				Low: common.LowerLimit{
					// all 0
				},
			},
		},
		// High withdraw-fee
		{
			name: "High withdraw-fee",
			asks: []common.Order{
				{Price: 100, Amount: 10},
			},
			bids: []common.Order{
				{Price: 105, Amount: 10},
			},
			thenMakerFee: 0.1,
			withdrawFee:  5,
			expected: common.VolumeData{
				Up: common.UpperLimit{
					// all 0
				},
				Low: common.LowerLimit{
					// all 0
				},
			},
		},
		// Deep orderbook
		{
			name: "Deep orderbook",
			asks: []common.Order{
				{Price: 100, Amount: 2},
				{Price: 101, Amount: 3},
				{Price: 102, Amount: 5},
				{Price: 103, Amount: 10},
			},
			bids: []common.Order{
				{Price: 105, Amount: 1},
				{Price: 104, Amount: 4},
				{Price: 103, Amount: 5},
			},
			thenMakerFee: 0,
			withdrawFee:  0,
			// asks (100 * 2 + 101 * 3 + 102 * 5)/10 = 101.3
			// bids (105 * 1 + 104 * 4 + 103 * 3)/8 = 103.75
			expected: common.VolumeData{
				Up: common.UpperLimit{
					AskVolUsdt: 1013,
					BidVolUsdt: 1036,
					BuyPrice:   101.3,
					SellPrice:  103.6,
					Profit:     23,
					AskDepth:   3,
					BidDepth:   3,
					Spread:     2.2704837117472825,
				},
				Low: common.LowerLimit{
					AskVolUsdt: 100,
					BidVolUsdt: 105,
					BuyPrice:   100,
					SellPrice:  105,
					Profit:     5,
					AskDepth:   1,
					BidDepth:   1,
					Spread:     5,
				},
			},
		},
		// Hard orderbook
		{
			name: "Hard orderbook",
			asks: []common.Order{
				{Price: 110, Amount: 10},
				{Price: 115, Amount: 5},
				{Price: 116, Amount: 8},
				{Price: 117, Amount: 11},
				{Price: 120, Amount: 20},
			},
			bids: []common.Order{
				{Price: 120, Amount: 13},
				{Price: 115, Amount: 8},
				{Price: 113, Amount: 2},
				{Price: 111, Amount: 4},
				{Price: 108, Amount: 16},
			},
			thenMakerFee: 0,
			withdrawFee:  0.1,
			expected: common.VolumeData{
				Up: common.UpperLimit{
					AskVolUsdt: 3890,
					BidVolUsdt: 3895.2000000000003,
					BuyPrice:   114.41176470588235,
					SellPrice:  114.90265486725664,
					Profit:     5.200000000000273,
					AskDepth:   4,
					BidDepth:   5,
					Spread:     0.42905566803923006,
				},
				Low: common.LowerLimit{
					AskVolUsdt: 220,
					BidVolUsdt: 228,
					BuyPrice:   110,
					SellPrice:  120,
					Profit:     8,
					AskDepth:   1,
					BidDepth:   1,
					Spread:     9.090909090909092,
				},
				Best: common.BestVolume{
					AskVolUsdt: 1560,
					BidVolUsdt: 1663.5,
					BuyPrice:   111.42857142857143,
					SellPrice:  119.67625899280576,
					Profit:     103.5,
					AskDepth:   2,
					BidDepth:   2,
					Spread:     7.401770890979523,
				},
			},
		},
		// TO-DO
		// Нет профита из-за fee вначале стакана, но потом он появляется
		// !! не знаю возможно ли это
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Logf("Asks: %+v, Bids: %+v, thenMakerFee: %f, withdrawFee: %f",
				tt.asks, tt.bids, tt.thenMakerFee, tt.withdrawFee)
			volData := CalcVolume(&tt.asks, &tt.bids, tt.thenMakerFee, tt.withdrawFee)
			t.Logf("Got max volume: %f < %f, Expected: %v", volData.Up.AskVolUsdt, volData.Up.BidVolUsdt, tt.expected)
			t.Logf("res - %+v", volData)
			// проверка  верхней точки
			if volData.Up != tt.expected.Up {
				t.Errorf("CalcVolume().Up = %+v\nwant %+v", volData.Up, tt.expected.Up)
			}
			// проверка  нижней точки
			if volData.Low != tt.expected.Low {
				t.Errorf("CalcVolume().Low = %+v\nwant %+v", volData.Low, tt.expected.Low)
			}
		})
	}
}

func TestBookSum(t *testing.T) {
	tests := []struct {
		name     string
		orders   []common.Order
		expected float64
	}{
		{
			name:     "Empty orderbook",
			orders:   []common.Order{},
			expected: 0,
		},
		{
			name: "One order",
			orders: []common.Order{
				{Price: 100, Amount: 1},
			},
			expected: 1,
		},
		{
			name: "Multiple orders",
			orders: []common.Order{
				{Price: 100, Amount: 1},
				{Price: 101, Amount: 2},
				{Price: 102, Amount: 3},
			},
			expected: 6,
		},
		{
			name: "неровный стакан",
			orders: []common.Order{
				{Price: 100, Amount: 1.4},
				{Price: 101, Amount: 2.9},
				{Price: 102, Amount: 11.6},
			},
			expected: 1.4 + 2.9 + 11.6,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Logf("Orders: %+v", tt.orders)
			got := BookSum(&tt.orders)
			t.Logf("Got: %f, Expected: %f", got, tt.expected)
			if math.Abs(got-tt.expected) > 1e-9 {
				t.Errorf("bookSum() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestCalculatePrice(t *testing.T) {
	tests := []struct {
		name     string
		orders   []common.Order
		vol      float64
		expected float64
	}{
		{
			name: "Exact volume",
			orders: []common.Order{
				{Price: 100, Amount: 1},
				{Price: 101, Amount: 1},
			},
			vol:      2,
			expected: (100.0 + 101.0) / 2.0,
		},
		{
			name: "неровный стакан",
			orders: []common.Order{
				{Price: 100, Amount: 3.3},
				{Price: 101, Amount: 5.8},
			},
			vol:      4.2,
			expected: (100.0*3.3 + 101.0*0.9) / 4.2,
		},
		{
			name: "Partial last order",
			orders: []common.Order{
				{Price: 100, Amount: 1},
				{Price: 101, Amount: 2},
			},
			vol:      2,
			expected: (100.0*1 + 101.0*1) / 2.0,
		},
		{
			name: "Not enough orders",
			orders: []common.Order{
				{Price: 100, Amount: 1},
			},
			vol:      2,
			expected: -1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Logf("Orders: %+v, Volume: %f", tt.orders, tt.vol)
			got := CalculatePrice(&tt.orders, tt.vol)
			t.Logf("Got: %f, Expected: %f", got, tt.expected)
			if math.Abs(got-tt.expected) > 1e-9 {
				t.Errorf("calculatePrice() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestCalculateDepth(t *testing.T) {
	tests := []struct {
		name     string
		coins    float64
		orders   []common.Order
		expected int
	}{
		{
			name:  "Simple test",
			coins: 16,
			orders: []common.Order{
				{Price: 110, Amount: 10}, // 10
				{Price: 115, Amount: 5},  // 15
				{Price: 116, Amount: 8},  // 16
				{Price: 117, Amount: 11},
				{Price: 120, Amount: 20},
			},
			expected: 3,
		},
		{
			name:  "неровный стакан",
			coins: 24.6,
			orders: []common.Order{
				{Price: 110, Amount: 10.5}, // 10.5
				{Price: 115, Amount: 5.2},  // 15.7
				{Price: 116, Amount: 8.9},  // 24.6
				{Price: 117, Amount: 11},
				{Price: 120, Amount: 20},
			},
			expected: 3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Logf("orders: %+v, expected: %d",
				tt.orders, tt.expected)
			depth := CalculateDepth(&tt.orders, tt.coins)
			t.Logf("res - %v", depth)
			if tt.expected != depth {
				t.Errorf("Got depth: %d, Expected: %d", depth, tt.expected)
			}
		})
	}
}
