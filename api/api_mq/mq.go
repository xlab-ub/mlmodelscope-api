package api_mq

import (
	"github.com/c3sr/mq"
	"github.com/c3sr/mq/interfaces"
	"github.com/c3sr/mq/rabbit"
	"log"
	"sync"
	"time"
)

var messageQueue interfaces.MessageQueue
var mutex sync.Mutex

func ConnectToMq() {
	ready := make(chan bool)
	defer close(ready)
	dialer, err := rabbit.NewRabbitDialer()
	if err != nil {
		log.Fatalf("[FATAL] Failed to initialize RabbitMQ Dialer: %s", err.Error())
	}
	mq.SetDialer(dialer)

	go func() {
		for i := 0; i < 5; i++ {
			mq, err := mq.NewMessageQueue()
			if err != nil {
				log.Printf("[INFO] Waiting for message queue")
				time.Sleep(time.Second * 5)
			} else {
				SetMessageQueue(mq)
				ready <- true
				return
			}
		}
		ready <- false
	}()

	select {
	case isReady := <-ready:
		if isReady {
			log.Printf("[INFO] Connected to message queue")
		} else {
			log.Fatalf("[FATAL] Could not connect to message queue server")
		}
	}
}

func SetMessageQueue(queue interfaces.MessageQueue) {
	mutex.Lock()
	defer mutex.Unlock()

	messageQueue = queue
}

func GetMessageQueue() interfaces.MessageQueue {
	mutex.Lock()
	defer mutex.Unlock()

	return messageQueue
}
