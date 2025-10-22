package req

import (
	"mad-scanner.com/scriner/internal/proxy"
)

var clients = make(map[string]*ExchangeClient)

func InitExchangeClients(proxyPool *proxy.ProxyPool) {
	// Конфигурации rate limit для разных бирж
	exchangeConfigs := map[string]struct {
		rps   int
		burst int
	}{
		"kucoin":   {rps: 20, burst: 40},
		"mexc":     {rps: 20, burst: 40},
		"gate":     {rps: 20, burst: 40},
		"bybit":    {rps: 20, burst: 40},
		"bingx":    {rps: 20, burst: 40},
		"bitget":   {rps: 20, burst: 40},
		"htx":      {rps: 20, burst: 40},
		"okx":      {rps: 20, burst: 40},
		"poloniex": {rps: 20, burst: 40},
	}

	for name, config := range exchangeConfigs {
		clients[name] = NewExchangeClient(ClientConfig{
			Name:      name,
			RPS:       config.rps,
			Burst:     config.burst,
			ProxyPool: proxyPool,
		})
	}
}

func GetExchangeClient(name string) *ExchangeClient {
	return clients[name]
}
