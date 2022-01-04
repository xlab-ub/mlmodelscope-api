// +build integration

package endpoints

import (
	"api/api_db"
	"api/api_mq"
	"api/db/models"
	"encoding/json"
	"github.com/c3sr/mq/interfaces"
	"github.com/c3sr/mq/messages"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http/httptest"
	"os"
	"testing"
	"time"
)

var messageQueue interfaces.MessageQueue

func setupForIntegrationTest() {
	db, _ := api_db.GetDatabase()
	db.Migrate()
	db.CreateModel(&models.Model{
		Attributes:  models.ModelAttributes{},
		Description: "for integration test",
		Details:     models.ModelDetails{},
		Framework:   &models.Framework{
			Name:      "PyTorch",
			Version:   "1.0",
		},
		Input:       models.ModelOutput{},
		License:     "",
		Name:        "integrate",
		Output:      models.ModelOutput{},
		Version:     "1.0",
	})

	db.CreateModel(&models.Model{
		Attributes:  models.ModelAttributes{},
		Description: "for integration test",
		Details:     models.ModelDetails{},
		Framework:   &models.Framework{
			Name:      "Mock",
			Version:   "1.0",
		},
		Input:       models.ModelOutput{},
		License:     "",
		Name:        "Mock",
		Output:      models.ModelOutput{},
		Version:     "1.0",
	})

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
	requestBody := validPredictRequestBody()
	w := httptest.NewRecorder()
	req := NewJsonRequest("POST", "/predict", requestBody)
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code, "predict endpoint should return success")

	message := <-channel

	var response messages.PredictByModelName
	json.Unmarshal(message.Body, &response)
	assert.Equal(t, "integrate_1.0", response.ModelName)
	messageQueue.Acknowledge(message)
}

func TestPredictEndpointAgentRoundTrip(t *testing.T) {
	setupForIntegrationTest()
	channel, _ := messageQueue.SubscribeToChannel("API")

	router := SetupRoutes()
	requestBody := validPredictRequestBody()
	// Predictions for model ID 2 should be processed by the mock agent
	requestBody.Model = 2
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
