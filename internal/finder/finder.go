package finder

import (
	"fmt"
	"log"
	"sync"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"mad-scanner.com/scriner/common"
	"mad-scanner.com/scriner/internal/req"
)

func Start() {
	// Отладочная информация
	// go func() {
	// 	ticker := time.NewTicker(100 * time.Millisecond)
	// 	defer ticker.Stop()

	// 	for range ticker.C {
	// 		var m runtime.MemStats
	// 		runtime.ReadMemStats(&m)
	// 		fmt.Printf("\nnum goroutines: %d\n", runtime.NumGoroutine())
	// 		fmt.Printf("\tmemory: %d MB\n", m.Alloc/1024/1024)
	// 		fmt.Printf("\tTotalAlloc: %d MB\n", m.TotalAlloc/1024/1024)
	// 		fmt.Printf("\tHeapAlloc: %d MB\n", m.HeapAlloc/1024/1024)
	// 	}
	// }()

	type SNKey struct {
		Symbol  string
		Network string
	}
	var wg sync.WaitGroup

	for {
		var exchangeRates []common.ExchangeRate
		sourcesMap := make(map[SNKey][]common.ExchangeRate)
		destinationsMap := make(map[SNKey][]common.ExchangeRate)

		dsn := "host=postgres user=admin password=postgres dbname=scriner port=5432 sslmode=disable TimeZone=Europe/Moscow"
		db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err != nil {
			log.Fatalf("Error: Unable to connect to database: %v", err)
		}
		db.DB()

		// ORM-запрос
		err = db.Table("spot_exchange_rates").Where("timestamp > ?", time.Now().Add(-time.Minute)).Find(&exchangeRates).Error
		if err != nil {
			log.Fatalf("Error: Unable to fetch exchange rates: %v", err)
		}

		fmt.Println("Найдено записей:", len(exchangeRates))

		processed := make(map[string]bool)

		for _, er := range exchangeRates {
			if er.Ask == 0 || er.Bid == 0 {
				continue
			}
			key := SNKey{er.Symbol, er.Network[0]}
			if er.WithdrawOpen {
				sourcesMap[key] = append(sourcesMap[key], er)
			}
			if er.DepositOpen {
				destinationsMap[key] = append(destinationsMap[key], er)
			}
		}

		for _, now := range exchangeRates {
			key := SNKey{now.Symbol, now.Network[0]}
			if processed[fmt.Sprintf("%s:%s", now.Symbol, now.Network)] {
				continue
			}
			processed[fmt.Sprintf("%s:%s", now.Symbol, now.Network)] = true

			sources := sourcesMap[key]
			destinations := destinationsMap[key]

			for _, src := range sources {
				for _, dst := range destinations {
					// src - биржа с которой выводим
					// dst - биржа на которой получаем
					// src.Ask - buy
					// dst.Bid - sell
					if src.Exchange == dst.Exchange {
						continue
					}
					if !src.WithdrawOpen || !dst.DepositOpen {
						continue
					}

					testUsdtAmount := 300.0
					coins := testUsdtAmount / (src.Bid * (1 + src.TakerFee))
					profit := dst.Bid*(coins-src.Fee) - src.Ask*coins

					spread := ((dst.Bid - src.Ask) / src.Ask) * 100
					if profit <= 0 || spread <= 0.5 {
						continue
					}

					SendData(src, dst)
				}
			}
		}

		fmt.Println("Обход завершён")
		wg.Wait()
		fmt.Println("Все горутины завершены.")
		select {}
	}
}

func SendData(now, then common.ExchangeRate) {
	// now - биржа с которой выводим
	// then - биржа на которой получаем
	// now.Ask - buy
	// then.Bid - sell
	resultCh_now := make(chan common.OrderBook)
	resultCh_then := make(chan common.OrderBook)

	go func() {
		resultCh_now <- GetOrderbook(now.Exchange, now.Symbol)
	}()
	go func() {
		resultCh_then <- GetOrderbook(then.Exchange, then.Symbol)
	}()

	var nowOrderBook, thenOrderBook common.OrderBook
	receivedNow, receivedThen := false, false

	for !receivedNow || !receivedThen {
		select {
		case ob := <-resultCh_now:
			nowOrderBook = ob
			receivedNow = true
		case ob := <-resultCh_then:
			thenOrderBook = ob
			receivedThen = true
		}
	}

	close(resultCh_now)
	close(resultCh_then)

	if len(nowOrderBook.Asks) == 0 || len(thenOrderBook.Bids) == 0 {
		fmt.Println("ERROR:finder::Order book is null")
		return
	}

	volData := CalcVolume(&nowOrderBook.Asks, &thenOrderBook.Bids, then.MakerFee, now.Fee)

	if volData.Up.BidVolUsdt <= 0 || volData.Up.AskVolUsdt <= 0 || volData.Low.Profit <= 0 {
		fmt.Println("WARNING:finder:Арбитража нет")
		return
	}

	serverData := common.ServerData{
		Symbol:             now.Symbol,
		DepositExchange:    then.Exchange,
		WithdrawalExchange: now.Exchange,
		WithdrawalNetwork:  now.Network[0],
		WithdrawalFee:      now.Fee,
		// WithdrawalTime string
		MakerFee: now.MakerFee,
		// TimeLife       string
		Data: volData,
		OrderBook: common.OrderBook{
			Asks: nowOrderBook.Asks,
			Bids: thenOrderBook.Bids,
		},
	}

	common.ServerChannel <- serverData
}

func GetOrderbook(exchange, symbol string) common.OrderBook {
	var result common.OrderBook

	switch exchange {
	case "MEXC":
		result = req.Mexc_get_orderbook(symbol)
	case "GATE":
		result = req.Gate_get_orderbook(symbol)
	case "BYBIT":
		result = req.Bybyit_get_orderbook(symbol)
	case "BINGX":
		result = req.Bingx_get_orderbook(symbol)
	case "BITGET":
		result = req.Bitget_get_orderbook(symbol)
	case "KUCOIN":
		result = req.Kucoin_get_orderbook(symbol)
	case "HTX":
		result = req.Htx_get_orderbook(symbol)
	case "OKX":
		result = req.Okx_get_orderbook(symbol)
	case "POLONIEX":
		result = req.Poloniex_get_orderbook(symbol)
	case "BITMART":
		result = req.Bitmart_get_orderbook(symbol)
	case "BINANCE":
		result = req.Binance_get_orderbook(symbol)
	case "COINEX":
		result = req.Coinex_get_orderbook(symbol)
	case "COINW":
		result = req.Coinw_get_orderbook(symbol)
	case "XT":
		result = req.Xt_get_orderbook(symbol)
	case "ASCENDEX":
		result = req.Ascendex_get_orderbook(symbol)
	case "DIGIFINEX":
		result = req.Digifinex_get_orderbook(symbol)
	case "BITRUE":
		result = req.Bitrue_get_orderbook(symbol)
	}

	return result
}
