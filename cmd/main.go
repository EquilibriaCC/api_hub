package main

import (
	"teamAPI/config"
	"teamAPI/tasks"
)

func main() {
	config.GenerateConfig()
	go tasks.RunTasks()
	//api.API()
	select {

	}
}
