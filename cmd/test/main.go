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
		common.ServerChannel <- GenerateRandomServerData()
		time.Sleep(cooldown) // Опциональная задержка
	}
}

func GenerateRandomServerData() common.ServerData {
	symbols := []string{"BTC", "ETH", "BNB", "XRP", "SOL", "ADA", "DOT", "DOGE", "AVAX", "MATIC"}
	exchanges := []string{"Binance", "OKX", "Bybit", "Kucoin", "Gate", "Mexc", "Bitget", "HTX", "Poloniex", "Bitmart"}
	networks := []string{"ERC20", "TRC20", "BEP20", "SOL", "AVAX", "Polygon", "Arbitrum", "Optimism", "Lightning"}

	symbol := symbols[rand.Intn(len(symbols))]

	// Генерация разных бирж для депозита и вывода
	depositExchange := exchanges[rand.Intn(len(exchanges))]
	withdrawalExchange := exchanges[rand.Intn(len(exchanges))]
	for withdrawalExchange == depositExchange {
		withdrawalExchange = exchanges[rand.Intn(len(exchanges))]
	}

	network := networks[rand.Intn(len(networks))]

	// Генерация правдоподобных значений
	return common.ServerData{
		Symbol:             symbol,
		DepositExchange:    depositExchange,
		WithdrawalExchange: withdrawalExchange,
		WithdrawalNetwork:  network,
		WithdrawalFee:      rand.Float64()*10 + 0.1,
		WithdrawalTime:     strconv.Itoa(rand.Intn(50)+10) + "s",
		MakerFee:           rand.Float64()*0.1 + 0.01,
		TimeLife:           strconv.Itoa(rand.Intn(50)+10) + "m",
		Data: common.VolumeData{
			Up: common.UpperLimit{
				AskVolUsdt: rand.Float64()*1000000 + 10000,
				BidVolUsdt: rand.Float64()*1000000 + 10000,
				BuyPrice:   rand.Float64()*10000 + 100,
				SellPrice:  rand.Float64()*10000 + 100,
				Profit:     rand.Float64()*20 - 10,
				AskDepth:   rand.Intn(1000) + 100,
				BidDepth:   rand.Intn(1000) + 100,
				Spread:     rand.Float64()*0.5 + 0.1,
			},
			Low: common.LowerLimit{
				AskVolUsdt: rand.Float64()*100000 + 1000,
				BidVolUsdt: rand.Float64()*100000 + 1000,
				BuyPrice:   rand.Float64()*10000 + 100,
				SellPrice:  rand.Float64()*10000 + 100,
				Profit:     rand.Float64()*20 - 10,
				AskDepth:   rand.Intn(500) + 50,
				BidDepth:   rand.Intn(500) + 50,
				Spread:     rand.Float64()*0.5 + 0.1,
			},
		},
		OrderBook: common.OrderBook{
			Asks: []common.Order{
				{Price: rand.Float64()*10000 + 100, Amount: rand.Float64()*10 + 1},
				{Price: rand.Float64()*10000 + 100, Amount: rand.Float64()*10 + 1},
				{Price: rand.Float64()*10000 + 100, Amount: rand.Float64()*10 + 1},
			},
			Bids: []common.Order{
				{Price: rand.Float64()*10000 + 100, Amount: rand.Float64()*10 + 1},
				{Price: rand.Float64()*10000 + 100, Amount: rand.Float64()*10 + 1},
				{Price: rand.Float64()*10000 + 100, Amount: rand.Float64()*10 + 1},
			},
		},
	}
}
