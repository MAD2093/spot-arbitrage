package benchs

import (
	"sync"
	"testing"

	"mad-scanner.com/scriner/common"
	"mad-scanner.com/scriner/internal/finder"
)

// Генератор тестовых ExchangeRate
func mock_ExchangeRate() common.ExchangeRate {
	return common.ExchangeRate{
		Exchange: "MEXC",
		Symbol:   "BTC/USDT",
		TakerFee: 0.001,
		MakerFee: 0.001,
		Fee:      0.0005,
		Bid:      99000,
		Ask:      99100,
		Network:  []string{"ERC20"},
	}
}

func Benchmark_SendData(b *testing.B) {
	now := mock_ExchangeRate()
	then := mock_ExchangeRate()
	then.Exchange = "BYBIT"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		finder.SendData(now, then)
	}
}

func Benchmark_SendDataWithGoroutine(b *testing.B) {
	now := mock_ExchangeRate()
	then := mock_ExchangeRate()
	then.Exchange = "KUCOIN"

	var wg sync.WaitGroup

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			finder.SendData(now, then)
		}()
	}
	wg.Wait()
}

func Benchmark_GetOrderbook(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		finder.GetOrderbook("OKX", "BTC/USDT")
	}
}
