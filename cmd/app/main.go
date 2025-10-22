package main

import (
	"flag"

	"github.com/rs/zerolog"
	"mad-scanner.com/scriner/common"
	"mad-scanner.com/scriner/internal/finder"
	"mad-scanner.com/scriner/internal/proxy"
	"mad-scanner.com/scriner/internal/req"
	"mad-scanner.com/scriner/pkg/clients/redis"
	"mad-scanner.com/scriner/pkg/log"
)

func main() {
	go redis.RedisChannelListener(common.ServerChannel)

	// получаем флаги запуска
	debug := flag.Bool("debug", false, "sets log level to debug")
	pretty := flag.Bool("pretty", false, "sets log view as pretty")

	flag.Parse()

	var logLevel zerolog.Level

	if *debug {
		logLevel = zerolog.DebugLevel
	} else {
		logLevel = zerolog.InfoLevel
	}

	log.Init(log.Config{
		Level:  logLevel,
		Pretty: *pretty,
	})

	// Инициализация пула прокси
	proxies := []string{
		"https://PrsUSPMFXVEO:kJmmcjIF@net-193-233-223-200.mcccx.com:8443",
	}
	proxyPool := proxy.NewProxyPool(proxies)

	req.InitExchangeClients(proxyPool)

	finder.Start()
}
