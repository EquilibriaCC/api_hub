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
	HardforkInfo HardForkInfoStruct
	StakingRequirement StakingRequirementStruct
	BlockHeader BlockHeaderStruct
	ServiceNodes ServiceNodeStruct
	QuorumServiceNodes QuorumServiceNodesStruct
	Height int
	Supply int
)

func RunTasks() {
	networkInfo()
	Emissions()
	//marketInfo()
	//hardforkInfo()

	geckoTime, networkInfoTime, exchangePricesTime, hfInfoTimer := time.Now(), time.Now(), time.Now(), time.Now()
	stakingReqTimer, blockHeadTimer, snTimer, emissionTimer := time.Now(), time.Now(), time.Now(), time.Now()
	for {
		if time.Since(geckoTime) > time.Second*10 {
			go marketInfo()
			geckoTime = time.Now()
		}
		if time.Since(networkInfoTime) > time.Second*30 {
			go networkInfo()
			networkInfoTime = time.Now()
		}
		if time.Since(exchangePricesTime) > time.Second * 5 {
			go getPriceFromExchanges()
			exchangePricesTime = time.Now()
		}
		if time.Since(hfInfoTimer) > time.Hour*24 {
			go getPriceFromExchanges()
			hfInfoTimer = time.Now()
		}
		if time.Since(blockHeadTimer) > time.Second*5 {
			go getPriceFromExchanges()
			hfInfoTimer = time.Now()
		}
		if time.Since(stakingReqTimer) > time.Second*15 {
			go getPriceFromExchanges()
			hfInfoTimer = time.Now()
		}
		if time.Since(snTimer) > time.Minute {
			go serviceNodes()
			snTimer = time.Now()
		}
		if time.Since(emissionTimer) > time.Hour*12 {
			go Emissions()
			emissionTimer = time.Now()
		}
	}
}