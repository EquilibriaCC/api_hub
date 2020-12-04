package main

import (
	"xeq_hub/api"
	"xeq_hub/config"
)

func main() {
	//go tasks.RunScheduledTasks()
	api.API()
	select {

	}
}

func init() {
	config.GenerateConfig()
}
