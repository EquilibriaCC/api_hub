package tasks

import (
	"encoding/json"
	"log"
	"teamAPI/config"
)

func getHeight() {

}

func getPriceFromExchanges() {
	go tradeogreInfo()
	go hotbitInfo()
}

func marketInfo() CongeckoMarketResponse {
	data := apiRequest("GET", "https://api.coingecko.com/api/v3/coins/triton?tickers=true&market_data=true&community_data=false&developer_data=false&sparkline=false", nil)
	var res CongeckoMarketResponse
	err := json.Unmarshal(data, &res)
	if err != nil {
		log.Println(err.Error())
		return CongeckoMarketResponse{}
	}
	CoingeckoInfo = res
	return res
}

func networkInfo() NetworkInfoStruct {
	params := map[string]string{"jsonrpc":"2.0","id":"0","method":"get_info"}
	data := apiRequest("POST", config.DaemonURL+"/json_rpc", params)

	var res NetworkInfoStruct
	err := json.Unmarshal(data, &res)
	if err != nil {
		log.Println(err.Error())
		return NetworkInfoStruct{}
	}
	NetInfo.Result.HashRate = res.Result.Difficulty / 120
	NetInfo = res
	return res
}

func tradeogreInfo() TradeOgrePriceStruct {
	data := apiRequest("GET", "https://tradeogre.com/api/v1/ticker/BTC-XEQ", nil)
	var res TradeOgrePriceStruct
	err := json.Unmarshal(data, &res)
	if err != nil {
		log.Println(err.Error())
		return TradeOgrePriceStruct{}
	}
	TradeOgreInfo = res
	return res
}

func hotbitInfo() {
	data := apiRequest("GET", "https://api.hotbit.io/api/v1/order.depth?market=XEQ/USDT&limit=100&interval=1e-8", nil)
	var res HotBitOrderBook
	err := json.Unmarshal(data, &res)
	if err != nil {
		log.Println(err.Error())
		return
	}
	HOTBITBTC = res
	data = apiRequest("GET", "https://api.hotbit.io/api/v1/order.depth?market=XEQ/USDT&limit=100&interval=1e-8", nil)
	err = json.Unmarshal(data, &res)
	if err != nil {
		log.Println(err.Error())
		return
	}
	HOTBITUSDT = res
}


