package main

import (
	"api/api_mq"
	"api/endpoints"
	"api/status"
)

func main() {
	api_mq.ConnectToMq()
	done := make(chan bool)
	go status.StartTracker(done)

	r := endpoints.SetupRoutes()
	r.Run()

	done <- true
}
