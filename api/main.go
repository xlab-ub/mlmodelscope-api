package main

import (
	"api/api_mq"
	"api/endpoints"
	"api/status"
	"log"
)

var trackerDone chan bool

func main() {
	api_mq.ConnectToMq()
	trackerDone = make(chan bool)
	go status.StartTracker(trackerDone)
	go keepMqAlive()

	r := endpoints.SetupRoutes()
	r.Run()

	trackerDone <- true
}

func keepMqAlive() {
	for {
		closed := make(chan error)
		api_mq.GetMessageQueue().NotifyClose(closed)

		for err := range closed {
			log.Printf("[INFO] MessageQueue connection closed unexpectedly: %s", err.Error())
		}

		trackerDone <- true
		api_mq.ConnectToMq()
		go status.StartTracker(trackerDone)
	}
}
