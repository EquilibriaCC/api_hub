package tasks

import (
	"time"
)

var (
	NetInfo NetworkInfoStruct
	CoingeckoInfo CongeckoMarketResponse
	TradeOgreInfo TradeOgrePriceStruct
	HOTBITBTC HotBitOrderBook
	HOTBITUSDT HotBitOrderBook
)

func RunTasks() {
	go networkInfo()
	go marketInfo()

	geckoTime, networkInfoTime, exchangePricesTime := time.Now(), time.Now(), time.Now()
	for {
		if time.Since(geckoTime) > time.Second*10 {
			go marketInfo()
		}
		if time.Since(networkInfoTime) > time.Second*30 {
			go networkInfo()
		}
		if time.Since(exchangePricesTime) > time.Second * 5 {
			go getPriceFromExchanges()
		}
	}
}