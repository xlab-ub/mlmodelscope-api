package main

import (
	"api/api_mq"
	"api/endpoints"
)

func main() {
	go api_mq.ConnectToMq()

	r := endpoints.SetupRoutes()
	r.Run()
}
