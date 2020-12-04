package config

import (
	"gopkg.in/yaml.v2"
	"log"
	"os"
)

type Config struct {
	DaemonURL         string  `yaml:"daemonURL"`
	APIPort           string  `yaml:"APIPort"`
	RateLimitTime     float64 `yaml:"RateLimitTime`
	RateLimitRequests int     `yaml:"RateLimitRequests"`
	DatabaseName      string  `yaml:"databaseName"`
	DatabaseHost      string  `yaml:"databaseHost"`
	DatabaseUsername  string  `yaml:"databaseUsername"`
	DatabasePassword  string  `yaml:"databasePassword"`
}

func GenerateConfig() {
	var conf Config
	file, err := os.Open("config.yml")
	if err != nil {
		file, err = os.Open("build/config.yml")
		if err != nil {
			log.Println(err.Error())
			log.Fatal("Config is invalid")
		}
	}
	defer file.Close()

	dec := yaml.NewDecoder(file)
	if err := dec.Decode(&conf); err != nil {
		log.Println(err.Error())
	}

	DaemonURL = conf.DaemonURL
	APIPort = conf.APIPort
	RateLimitTime = conf.RateLimitTime
	RateLimitRequests = int64(conf.RateLimitRequests)
	DatabaseName = conf.DatabaseName
	DatabaseHost = conf.DatabaseHost
	DatabaseUsername = conf.DatabaseUsername
	DatabasePassword = conf.DatabasePassword
	OracleNodeHistoryFileName = "emissionhistory.json"
	EmissionHistoryFileName = "oraclenodehistory.json"
}
