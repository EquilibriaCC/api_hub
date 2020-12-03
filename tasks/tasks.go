package tasks

import (
	"encoding/json"
	"log"
	"xeq_hub/config"
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

func stakingRequirement() StakingRequirementStruct {
	params := map[string]interface{}{"jsonrpc": "2.0", "id": "0", "method": "get_staking_requirement", "params": map[string]int{"height": NetInfo.Result.Height}}
	data := apiRequest("POST", config.DaemonURL+"/json_rpc", params)
	var res StakingRequirementStruct
	err := json.Unmarshal(data, &res)
	if err != nil {
		log.Println(err.Error())
		return StakingRequirementStruct{}
	}
	StakingRequirement = res
	return res
}

func blockHeader() BlockHeaderStruct {
	params := map[string]interface{}{"jsonrpc": "2.0", "id": "0", "method": "get_last_block_header", "params": map[string]bool{"fill_pow_hash": true}}
	data := apiRequest("POST", config.DaemonURL+"/json_rpc", params)
	var res BlockHeaderStruct
	err := json.Unmarshal(data, &res)
	if err != nil {
		log.Println(err.Error())
		return BlockHeaderStruct{}
	}
	BlockHeader = res
	return res
}

func hardforkInfo() HardForkInfoStruct {
	params := map[string]string{"jsonrpc": "2.0", "id": "0", "method": "hard_fork_info"}
	data := apiRequest("POST", config.DaemonURL+"/json_rpc", params)
	var res HardForkInfoStruct
	err := json.Unmarshal(data, &res)
	if err != nil {
		log.Println(err.Error())
		return HardForkInfoStruct{}
	}
	HardforkInfo = res
	return res
}

func networkInfo() NetworkInfoStruct {
	params := map[string]string{"jsonrpc": "2.0", "id": "0", "method": "get_info"}
	data := apiRequest("POST", config.DaemonURL+"/json_rpc", params)

	var res NetworkInfoStruct
	err := json.Unmarshal(data, &res)
	if err != nil {
		log.Println(err.Error())
		return NetworkInfoStruct{}
	}
	NetInfo.Result.HashRate = res.Result.Difficulty / 120
	Height = res.Result.Height
	NetInfo = res
	return res
}

func serviceNodes() {
	params := map[string]string{"jsonrpc": "2.0", "id": "0", "method": "get_service_nodes"}
	data := apiRequest("POST", config.DaemonURL+"/json_rpc", params)
	var res ServiceNodeStruct
	err := json.Unmarshal(data, &res)
	if err != nil {
		log.Println(err.Error())
		return
	}
	ServiceNodes = res

	paramsQuorum := map[string]interface{}{"jsonrpc": "2.0", "id": "0", "method": "get_quorum_state", "params": map[string]int{"height": Height - 1}}
	data = apiRequest("POST", config.DaemonURL+"/json_rpc", paramsQuorum)
	var resQ QuorumServiceNodesStruct
	err = json.Unmarshal(data, &resQ)
	if err != nil {
		log.Println(err.Error())
		return
	}
	QuorumServiceNodes = resQ
}

func Emissions() {

	supply := 0
	for x, y := 0, 5000; x < Height + 5000; x, y = x + 5000, y + 5000 {
		if y > Height {
			y = Height - 1
		}
		params := map[string]interface{}{"jsonrpc": "2.0", "id": "0", "method": "get_coinbase_tx_sum", "params": map[string]int{"height": x, "count": y}}
		data := apiRequest("POST", config.DaemonURL+"/json_rpc", params)
		var res EmissionsStruct
		err := json.Unmarshal(data, &res)
		if err != nil {
			log.Println(err.Error())
			return
		}
		supply += res.Result.EmissionAmount
	}
	Supply = supply
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
