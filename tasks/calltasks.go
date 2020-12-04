package tasks

import (
	"encoding/json"
	"fmt"
	"log"
	"sort"
	"xeq_hub/config"
)

func search(hash []string) []map[string]interface{} {

	hashList, _ := json.Marshal(TxsHashes{hash, true})
	data := apiRequest("POST", "http://sanfran.equilibria.network:9231/get_transactions", hashList)

	var TxSearch TxResponseStruct

	err := json.Unmarshal(data, &TxSearch)
	if err != nil {
		panic(err)
	}

	var responses []map[string]interface{}
	for i := 0; i < len(TxSearch.Txs); i++ {
		var asJSON ExtraTxJSONStruct
		err = json.Unmarshal([]byte(TxSearch.Txs[i].AsJSON), &asJSON)
		if err != nil {
			log.Println(err)
			return nil
		}

		var totalOutput int
		for i := 0; i < len(asJSON.Vout); i++ {
			totalOutput = asJSON.Vout[i].Amount + totalOutput
		}

		var obj = map[string]interface{}{
			"block":        TxSearch.Txs[i].BlockHeight,
			"timestamp":    TxSearch.Txs[i].BlockTimestamp,
			"tx_hash":      TxSearch.Txs[i].TxHash,
			"unlock_block": asJSON.UnlockTime,
			"version":      asJSON.UnlockTime,
			"outputs":      asJSON.Vout,
			"inputs":       asJSON.Vin,
			"total_output": totalOutput,
			"ringCT_type":  asJSON.RctSignatures.Type,
			"txnFee":       asJSON.RctSignatures.TxnFee,
		}
		responses = append(responses, obj)
	}
	return responses
}

func getHashList(height int) HashesStruct {

	params := map[string]interface{}{"jsonrpc": "2.0", "id": "0", "method": "get_staking_requirement", "params": map[string]int{"height": height}}
	data := apiRequest("POST", config.DaemonURL+"/json_rpc", params)

	type Response struct {
		Result HashesStruct `json:"result"`
	}
	var response Response
	err := json.Unmarshal(data, &response)
	if err != nil {
		log.Println(err)
		return HashesStruct{}
	}

	return response.Result

}

func Transactions(startHeight int) []interface{} {

	height := Height - 1 - startHeight
	var blocks []int
	for i := 0; i < 11; i++ {
		blocks = append(blocks, height-i)
	}

	var hashes []string
	for _, v := range blocks {
		newHash := getHashList(v)
		hashes = append(hashes, newHash.MinerTxHash)
		for i := 0; i < len(newHash.TxHashes); i++ {
			hashes = append(hashes, newHash.TxHashes[i])
		}
		height--
	}

	data := search(hashes)

	var cleanData = make(map[int][]interface{})
	for i := 0; i < len(data); i++ {
		var blockHeight = data[i]["block"].(int)
		obj := map[string]interface{}{
			"confirmation": height - blockHeight,
		}
		for k, v := range data[i] {
			obj[k] = v
		}
		if _, ok := cleanData[blockHeight];ok {
			cleanData[blockHeight] = append(cleanData[blockHeight], obj)
		} else {
			cleanData[blockHeight] = []interface{}{obj}
		}
	}

	keys := make([]int, len(cleanData))
	i := 0
	for k := range cleanData {
		keys[i] = k
		i++
	}
	sort.Ints(keys)
	var sortedData []interface{}
	for _, k := range keys {
		sortedData = append(sortedData, map[int]interface{}{k: []interface{}{cleanData[k]}})
	}
	return sortedData
}

func GetTxPool() []map[string]interface{} {

	data := apiRequest("POST", "http://sanfran.equilibria.network:9231/get_transaction_pool", nil)
	log.Println(string(data))
	var txs TxResponseStruct
	err := json.Unmarshal(data, &txs)
	if err != nil {
		fmt.Println(err)
	}

	var output []map[string]interface{}
	for i := 0; i < len(txs.Transactions); i++ {
		newTxJson := []byte(txs.Transactions[i].TxJSON)
		var txsInJson ExtraTxJSONStruct
		err := json.Unmarshal(newTxJson, &txsInJson)
		if err != nil {
			fmt.Println(err)
		}
		var obj = map[string]interface{}{
			"hash":        txs.Transactions[i].IDHash,
			"timestamp":   txs.Transactions[i].ReceiveTime,
			"version":     txsInJson.Version,
			"outputs":     txsInJson.Vout,
			"inputs":      txsInJson.Vin,
			"ringCT_type": txsInJson.RctSignatures.Type,
			"fee":         txsInJson.RctSignatures.TxnFee,
		}
		output = append(output, obj)
	}
	return output
}

func SearchTx(hash string) map[string]interface{} {
	data := search([]string{hash})

	blockHeight := data[0]["block"].(int)
	obj := map[string]interface{}{"confirmation": Height - blockHeight}
	for k, v := range data[0] {
		obj[k] = v
	}
	return obj
}

func NumberOfOracleNodes() int {
	numNodes := 0
	if len(OracleNodeList) < 1 {
		return 0
	}
	for range OracleNodeList[0] {
		numNodes++
	}
	return numNodes
}
