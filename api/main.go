package main

import (
	"api/endpoints"
	"github.com/c3sr/mq"
	"github.com/c3sr/mq/interfaces"
	"github.com/c3sr/mq/rabbit"
	"log"
	"time"
)

var messageQueue interfaces.MessageQueue

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("[FATAL] %s: %s", msg, err)
	}
}

func main() {
	go connectToMq()

	r := endpoints.SetupRoutes()
	r.Run()
}

func connectToMq() {
	ready := make(chan bool)
	defer close(ready)
	dialer, err := rabbit.NewRabbitDialer()
	failOnError(err, "Failed to initialize RabbitMQ dialer")
	mq.SetDialer(dialer)

	go func() {
		for i := 0; i < 5; i++ {
			mq, err := mq.NewMessageQueue()
			if err != nil {
				log.Printf("[INFO] Waiting for message queue")
				time.Sleep(time.Second * 5)
			} else {
				messageQueue = mq
				ready <- true
				return
			}
		}
	}()

	select {
	case isReady := <-ready:
		if isReady {
			defer messageQueue.Shutdown()
			log.Printf("[INFO] Connected to message queue")
		} else {
			log.Fatalf("[FATAL] Could not connect to message queue server")
		}
	}
}
