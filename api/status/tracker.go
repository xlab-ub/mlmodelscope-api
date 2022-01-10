package status

import (
	"api/api_db"
	"api/api_mq"
	db "api/db"
	"fmt"
	"github.com/c3sr/mq/interfaces"
	"log"
)

var (
	channel  <-chan interfaces.Message
	database *db.Db
	mq       interfaces.MessageQueue
)

func StartTracker(done chan bool) {
	mq = api_mq.GetMessageQueue()
	d, err := api_db.GetDatabase()
	if err != nil {
		log.Fatalf("[FATAL] failed to get database instance in Status Tracker: %s", err.Error())
	}
	database = d

	connectChannel()

	running := true
	for running {
		select {
		case message, ok := <-channel:
			if !ok {
				log.Println("[INFO] channel closed unexpectedly, reconnecting")
				connectChannel()
			} else {
				processMessage(message)
			}
		case <-done:
			running = false
			fmt.Println("[INFO] done")
		}
	}
}

func connectChannel() {
	c, err := mq.SubscribeToChannel("API")
	channel = c
	if err != nil {
		log.Fatalf("[FATAL] failed to subscribe to API channel in Status Tracker: %s", err.Error())
	}
}

func processMessage(message interfaces.Message) {
	trial, err := database.GetTrialById(message.CorrelationId)
	if err != nil {
		log.Printf("[WARN] failed to retrieve Trial ID %s in Status Tracker: %s", message.CorrelationId, err.Error())
		return
	}

	if err = database.CompleteTrial(trial, string(message.Body)); err != nil {
		log.Printf("[WARN] failed to complete Trial ID %s in Status Tracker: %s", message.CorrelationId, err.Error())
	}

	if err = mq.Acknowledge(message); err != nil {
		log.Printf("[WARN] failed to acknowledge message %s in Status Tracker: %s", message.CorrelationId, err.Error())
	}
}
