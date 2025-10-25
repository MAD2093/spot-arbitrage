package main

import (
	"math/rand"
	"strconv"
	"time"

	"mad-scanner.com/scriner/common"
	"mad-scanner.com/scriner/pkg/clients/redis"
)

func main() {
	test_webserver()
	select {}
}

// TESTS
func test_webserver() {
	go redis.RedisChannelListener(common.ServerChannel)
	go generate_websocket_data(100, 3*time.Second)
}

// count - кол-во сгенерированных данных;
// cooldown - задержка между пушами в канал;
func generate_websocket_data(count int, cooldown time.Duration) {
	for i := 0; i < count; i++ {
		common.ServerChannel <- MockData()
		time.Sleep(cooldown)
	}
}

func MockData() common.ServerData {
	return common.ServerData{
		Symbol:             "BTC",
		DepositExchange:    "MEXC",
		WithdrawalExchange: "OKX",
		WithdrawalNetwork:  "ERC20",
		WithdrawalFee:      0.1,
		WithdrawalTime:     strconv.Itoa(rand.Intn(50)+10) + "s",
		MakerFee:           0,
		TimeLife:           strconv.Itoa(rand.Intn(50)+10) + "m",
		Volume24h:          142986,
		Data: common.VolumeData{
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
		OrderBook: common.OrderBook{
			Asks: []common.Order{
				{Price: 110, Amount: 10},
				{Price: 115, Amount: 5},
				{Price: 116, Amount: 8},
				{Price: 117, Amount: 11},
				{Price: 120, Amount: 20},
			},
			Bids: []common.Order{
				{Price: 120, Amount: 13},
				{Price: 115, Amount: 8},
				{Price: 113, Amount: 2},
				{Price: 111, Amount: 4},
				{Price: 108, Amount: 16},
			},
		},
	}
}
