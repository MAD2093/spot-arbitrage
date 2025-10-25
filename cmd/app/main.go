package main

import (
	"flag"
	"net/url"

	"github.com/rs/zerolog"
	"mad-scanner.com/scriner/common"
	"mad-scanner.com/scriner/internal/finder"
	"mad-scanner.com/scriner/internal/req"
	"mad-scanner.com/scriner/pkg/clients/http/proxy"
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

	log := log.Init(log.Config{
		Level:  logLevel,
		Pretty: *pretty,
	})

	// Инициализация пула прокси
	proxies := []string{
		"https://PrsUSPMFXVEO:kJmmcjIF@net-193-233-223-200.mcccx.com:8443",
	}
	var pp []*url.URL
	for k, v := range proxies {
		url, err := url.Parse(v)
		if err != nil {
			log.Panic().Msgf("не получилось спарсить прокси под индексом %d", k)
		}
		pp = append(pp, url)
	}
	proxyPool := proxy.NewProxyPool(pp)

	req.InitExchangeClients(proxyPool)

	finder.Start()
}
