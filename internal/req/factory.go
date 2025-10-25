package req

import (
	"mad-scanner.com/scriner/pkg/clients/http"
	"mad-scanner.com/scriner/pkg/clients/http/proxy"
)

var clients = make(map[string]*http.HTTPClient)

func InitExchangeClients(proxyPool *proxy.ProxyPool) {
	exchanges := [17]string{
		"MEXC",
		"GATE",
		"BYBIT",
		"BINGX",
		"BITGET",
		"KUCOIN",
		"HTX",
		"OKX",
		"POLONIEX",
		"BITMART",
		"BINANCE",
		"COINEX",
		"COINW",
		"XT",
		"ASCENDEX",
		"DIGIFINEX",
		"BITRUE",
	}

	for _, name := range exchanges {
		clients[name] = http.NewHTTPClient(proxyPool)
	}
}

func GetExchangeClient(name string) *http.HTTPClient {
	return clients[name]
}
