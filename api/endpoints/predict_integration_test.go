// +build integration

package endpoints

import (
	"api/api_db"
	"api/api_mq"
	"api/db/models"
	"api/status"
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
var router *gin.Engine
var trackerDone chan bool

func setupForIntegrationTest() {
	trackerDone = make(chan bool)
	db, _ := api_db.GetDatabase()
	db.Migrate()
	db.CreateModel(&models.Model{
		Attributes:  models.ModelAttributes{},
		Description: "for integration test",
		Details:     models.ModelDetails{},
		Framework: &models.Framework{
			Name:    "PyTorch",
			Version: "1.0",
		},
		Input:   models.ModelOutput{},
		License: "",
		Name:    "integrate",
		Output:  models.ModelOutput{},
		Version: "1.0",
	})

	db.CreateModel(&models.Model{
		Attributes:  models.ModelAttributes{},
		Description: "for integration test",
		Details:     models.ModelDetails{},
		Framework: &models.Framework{
			Name:    "Mock",
			Version: "1.0",
		},
		Input:   models.ModelOutput{},
		License: "",
		Name:    "Mock",
		Output:  models.ModelOutput{},
		Version: "1.0",
	})

	reconnectToMq()

	router = SetupRoutes()
}

func reconnectToMq() {
	if messageQueue != nil {
		messageQueue.Shutdown()
		messageQueue = nil
	}

	go api_mq.ConnectToMq()
	time.Sleep(time.Millisecond * 100)
	messageQueue = api_mq.GetMessageQueue()
}

func TestMain(m *testing.M) {
	gin.DefaultWriter = &nullWriter{}
	code := m.Run()
	os.Exit(code)
}

func TestPredictEndpoint(t *testing.T) {
	setupForIntegrationTest()

	t.Run("QueuesMessage", func(t *testing.T) {
		channel, err := messageQueue.SubscribeToChannel("agent-pytorch-amd64")
		assert.Nil(t, err, "SubscribeToChannel should succeed")

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
	})

	t.Run("AgentRoundTrip", func(t *testing.T) {
		db, _ := api_db.GetDatabase()
		channel, _ := messageQueue.SubscribeToChannel("API")

		requestBody := validPredictRequestBody()
		// Predictions for model ID 2 should be processed by the mock agent
		requestBody.Model = 2
		w := httptest.NewRecorder()
		req := NewJsonRequest("POST", "/predict", requestBody)
		router.ServeHTTP(w, req)

		response := &predictResponseBody{}
		json.Unmarshal(w.Body.Bytes(), response)
		message := <-channel
		prediction := &mockPredictionResponse{}
		json.Unmarshal(message.Body, prediction)

		assert.Equal(t, "ec1578ee-4ad8-46af-b7e7-10d6d1570abc", prediction.Id)
		messageQueue.Acknowledge(message)

		trial, _ := db.GetTrialById(message.CorrelationId)
		assert.NotNil(t, trial)
		assert.Equal(t, response.TrialId, trial.ID)
	})

	t.Run("StatusTrackerCompletesTrial", func(t *testing.T) {
		reconnectToMq()
		go status.StartTracker(trackerDone)
		db, _ := api_db.GetDatabase()
		requestBody := validPredictRequestBody()
		// Predictions for model ID 2 should be processed by the mock agent
		requestBody.Model = 2
		w := httptest.NewRecorder()
		req := NewJsonRequest("POST", "/predict", requestBody)
		router.ServeHTTP(w, req)

		response := &predictResponseBody{}
		json.Unmarshal(w.Body.Bytes(), response)

		time.Sleep(time.Millisecond * 100)

		trial, err := db.GetTrialById(response.TrialId)
		assert.Nil(t, err)
		assert.NotNil(t, trial)
		assert.NotEqual(t, "", trial.Result)

		trackerDone <- true
	})
}

type mockPredictionResponse struct {
	Id string `json:"id"`
}
