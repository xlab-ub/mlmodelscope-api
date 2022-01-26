package main

import (
	"api/api_db"
	"api/api_mq"
	"api/endpoints"
	"api/status"
	"log"
)

var trackerDone chan bool

func main() {
	migrateDatabase()
	api_mq.ConnectToMq()
	trackerDone = make(chan bool)
	go status.StartTracker(trackerDone)
	go keepMqAlive()

	r := endpoints.SetupRoutes()
	r.Run()

	trackerDone <- true
}

func migrateDatabase() {
	db, err := api_db.GetDatabase()
	if err != nil {
		log.Fatalf("[FATAL] failed to get database instance for migration: %s", err.Error())
	}

	err = db.Migrate()
	if err != nil {
		log.Printf("[WARN] failed to migrate database: %s", err.Error())
	}
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
