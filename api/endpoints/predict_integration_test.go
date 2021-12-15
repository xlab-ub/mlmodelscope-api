// +build integration

package endpoints

import (
	"api/api_mq"
	"encoding/json"
	"github.com/c3sr/mq"
	"github.com/c3sr/mq/interfaces"
	"github.com/c3sr/mq/rabbit"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http/httptest"
	"os"
	"testing"
	"time"
)

var messageQueue interfaces.MessageQueue

func setupForIntegrationTest() {
	dialer, _ := rabbit.NewRabbitDialer()
	mq.SetDialer(dialer)
	go api_mq.ConnectToMq()

	time.Sleep(time.Millisecond * 100)
	messageQueue = api_mq.GetMessageQueue()
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
	requestBody := validPredictRequestBody("pytorch")
	w := httptest.NewRecorder()
	req := NewJsonRequest("POST", "/predict", requestBody)
	router.ServeHTTP(w, req)

	message := <-channel

	assert.Equal(t, "do some work", string(message.Body))
	messageQueue.Acknowledge(message)
}

func TestPredictEndpointAgentRoundTrip(t *testing.T) {
	setupForIntegrationTest()
	channel, _ := messageQueue.SubscribeToChannel("API")

	router := SetupRoutes()
	requestBody := validPredictRequestBody("mock")
	w := httptest.NewRecorder()
	req := NewJsonRequest("POST", "/predict", requestBody)
	router.ServeHTTP(w, req)

	message := <-channel
	prediction := &mockPredictionResponse{}
	json.Unmarshal(message.Body, prediction)

	assert.Equal(t, "ec1578ee-4ad8-46af-b7e7-10d6d1570abc", prediction.Id)
	messageQueue.Acknowledge(message)
}

type mockPredictionResponse struct {
	Id string `json:"id"`
}
