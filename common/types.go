package common

import (
	sharedDataBase "github.com/MAD2093/shared-go/pkg/models/database"
	sharedSpot "github.com/MAD2093/shared-go/pkg/models/spotarbitrage"
)

// алиасы
type ServerData = sharedSpot.ServerData
type VolumeData = sharedSpot.VolumeData
type UpperLimit = sharedSpot.UpperLimit
type LowerLimit = sharedSpot.LowerLimit
type BestVolume = sharedSpot.BestVolume
type RedisMessage = sharedSpot.RedisMessage

type OrderBook = sharedSpot.OrderBook
type Order = sharedSpot.Order

type ExchangeRate = sharedDataBase.SpotExchangeRate

var ServerChannel chan ServerData = make(chan ServerData)
