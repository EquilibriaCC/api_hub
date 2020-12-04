package tasks

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"time"
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

var (
	writeNumNodesTimer = time.Now()
	emissionCountTimer = time.Now()
)

func resetTimer(t *time.Time) {
	*t = time.Now()
}

func updateHistoryFiles(fileName string, t *time.Time) {
	defer resetTimer(t)
	storageFile, err := os.Open("tempstorage/" + fileName)
	if err != nil {
		file, _ := json.MarshalIndent([]int{}, "", " ")
		err = ioutil.WriteFile("tempstorage/"+fileName, file, 0644)
		if err != nil {
			log.Println("Couldnt write " + fileName)
			return
		}
		storageFile, err = os.Open("tempstorage/" + fileName)
		if err != nil {
			log.Println("Couldnt initialize file " + fileName)
			return
		}
	}
	fileBytes, err := ioutil.ReadAll(storageFile)
	if err != nil {
		log.Println("Couldnt initialize file " + fileName)

	}

	var data []int
	err = json.Unmarshal(fileBytes, &data)
	if len(data) > 720 {
		data = data[len(data)-720:]
	}
	data = append(data, NumberOfOracleNodes())
	newData, err := json.Marshal(data)
	file, _ := json.MarshalIndent(newData, "", " ")
	err = ioutil.WriteFile("tempstorage/"+fileName, file, 0644)
	if err != nil {
		log.Println("Couldnt write " + fileName)
	}

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

	var sns []string
	for _, v := range ServiceNodes.Result.ServiceNodeStates {
		sns = append(sns, v.OperatorAddress)
	}
	QuorumServiceNodes = resQ
	OracleNodeList = [][]string{sns, QuorumServiceNodes.Result.NodesToTest, QuorumServiceNodes.Result.QuorumNodes}
	if time.Since(writeNumNodesTimer) > time.Hour*12 {
		updateHistoryFiles(config.OracleNodeHistoryFileName, &writeNumNodesTimer)
	}
}

func Emissions() {
	supply := 0
	for x, y := 0, 5000; x < Height+5000; x, y = x+5000, y+5000 {
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
	if time.Since(emissionCountTimer) > time.Hour*12 {
		updateHistoryFiles(config.EmissionHistoryFileName, &emissionCountTimer)
	}
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
