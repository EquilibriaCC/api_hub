package main

import (
	"xeq_hub/config"
	"xeq_hub/tasks"
)

func main() {
	config.GenerateConfig()
	go tasks.RunTasks()
	//api.API()
	select {

	}
}
