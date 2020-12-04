package tasks

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
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
	OracleNodeList [][]string
	OracleNodeHistory []int
	EmissionHistory []int
	Height int
	Supply int
)

func RunScheduledTasks() {
	go writeToStorageTasks()
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
			blockHeadTimer = time.Now()
		}
		if time.Since(stakingReqTimer) > time.Second*15 {
			go getPriceFromExchanges()
			stakingReqTimer = time.Now()
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

func writeToStorageTasks() {
	emissionTimer, numNodesTimer := time.Now(), time.Now()
	for {
		if time.Since(emissionTimer) > time.Hour * 12 {
			emissionTimer = time.Now()
		}
		if time.Since(numNodesTimer) > time.Hour * 12 {
			numNodesTimer = time.Now()
		}
	}
}

func init() {
	nodes := checkFiles("oraclenodehistory.json")
	var data []int
	err := json.Unmarshal(nodes, &data)
	if err != nil {
		log.Fatal("Could not initialise files")
	}
	OracleNodeHistory = data

	emissions := checkFiles("emissionhistory.json")
	err = json.Unmarshal(emissions, &data)
	if err != nil {
		log.Fatal("Could not initialise files")
	}
	EmissionHistory = data

	//networkInfo()
	//Emissions()
	//marketInfo()
	//hardforkInfo()
}

func checkFiles(fileName string) []byte {
	storageFile, err := os.Open("tempstorage/"+fileName)
	if err != nil {
		file, _ := json.MarshalIndent([]int{}, "", " ")
		err = ioutil.WriteFile("tempstorage/"+fileName, file, 0644)
		if err != nil {
			log.Fatal("Couldnt write " + fileName)
		}
		storageFile, err = os.Open("tempstorage/"+fileName)
		if err != nil {
			log.Fatal("Couldnt initialize file " + fileName)

		}
	}
	fileBytes, err := ioutil.ReadAll(storageFile)
	if err != nil {
		log.Fatal("Couldnt initialize file " + fileName)

	}
	return fileBytes
}