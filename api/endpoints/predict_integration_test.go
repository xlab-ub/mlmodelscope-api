// +build integration

package endpoints

import (
	"api/api_mq"
	"github.com/c3sr/mq"
	"github.com/c3sr/mq/interfaces"
	"github.com/c3sr/mq/rabbit"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http/httptest"
	"os"
	"testing"
)

var messageQueue interfaces.MessageQueue

func setupForIntegrationTest() {
	dialer, _ := rabbit.NewRabbitDialer()
	mq.SetDialer(dialer)
	messageQueue, _ = mq.NewMessageQueue()
	api_mq.SetMessageQueue(messageQueue)
}

func TestMain(m *testing.M) {
	gin.DefaultWriter = &nullWriter{}
	code := m.Run()
	os.Exit(code)
}

func TestPredictEndpointQueuesMessage(t *testing.T) {
	setupForIntegrationTest()
	channel, err := messageQueue.SubscribeToChannel("agent-pytorch-amd64")
	assert.Nil(t, err, "SubscribeToChannel should succeed")

	router := SetupRoutes()
	requestBody := validPredictRequestBody()
	w := httptest.NewRecorder()
	req := NewJsonRequest("POST", "/predict", requestBody)
	router.ServeHTTP(w, req)

	message := <-channel

	assert.Equal(t, "do some work", string(message.Body))
}
