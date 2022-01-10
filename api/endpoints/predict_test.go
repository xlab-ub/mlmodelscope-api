// +build !integration

package endpoints

import (
	"api/api_mq"
	"encoding/json"
	"github.com/c3sr/mq/messages"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPredictRoute(t *testing.T) {
	openDatabase()
	createTestModelAndFramework()
	defer cleanupTestDatabase()
	router := SetupRoutes()

	t.Run("RequiresContentType", func(t *testing.T) {
		w := httptest.NewRecorder()
		requestBody := validPredictRequestBody()
		jsonBody, _ := json.Marshal(requestBody)
		req, _ := http.NewRequest("POST", "/predict", strings.NewReader(string(jsonBody)))
		router.ServeHTTP(w, req)

		assert.Equal(t, 400, w.Code)
	})

	t.Run("RequiresArchitecture", func(t *testing.T) {
		requestBody := &predictRequestBody{
			Architecture:          "",
			BatchSize:             1,
			DesiredResultModality: "x",
			Inputs:                []string{"x"},
			Model:                 1,
			TraceLevel:            "NO_TRACE",
		}

		w := httptest.NewRecorder()
		req := NewJsonRequest("POST", "/predict", requestBody)
		router.ServeHTTP(w, req)

		assert.Equal(t, 400, w.Code)
	})

	t.Run("RequiresInputs", func(t *testing.T) {
		requestBody := &predictRequestBody{
			Architecture:          "x",
			BatchSize:             1,
			DesiredResultModality: "x",
			Inputs:                []string{},
			Model:                 1,
			TraceLevel:            "NO_TRACE",
		}

		w := httptest.NewRecorder()
		req := NewJsonRequest("POST", "/predict", requestBody)
		router.ServeHTTP(w, req)

		assert.Equal(t, 400, w.Code)
	})

	t.Run("RequiresValidModelId", func(t *testing.T) {
		requestBody := &predictRequestBody{
			Architecture:          "x",
			BatchSize:             1,
			DesiredResultModality: "x",
			Inputs:                []string{"x"},
			Model:                 2,
			TraceLevel:            "NO_TRACE",
		}

		w := httptest.NewRecorder()
		req := NewJsonRequest("POST", "/predict", requestBody)
		router.ServeHTTP(w, req)

		assert.Equal(t, 400, w.Code)
	})

	t.Run("RequiresDesiredResultModality", func(t *testing.T) {
		requestBody := &predictRequestBody{
			Architecture:          "x",
			BatchSize:             1,
			DesiredResultModality: "",
			Inputs:                []string{"x"},
			Model:                 1,
			TraceLevel:            "NO_TRACE",
		}

		w := httptest.NewRecorder()
		req := NewJsonRequest("POST", "/predict", requestBody)
		router.ServeHTTP(w, req)

		assert.Equal(t, 400, w.Code)
	})

	t.Run("RequiresBatchSize", func(t *testing.T) {
		requestBody := &predictRequestBody{
			Architecture:          "x",
			DesiredResultModality: "",
			Inputs:                []string{"x"},
			Model:                 1,
			TraceLevel:            "NO_TRACE",
		}

		w := httptest.NewRecorder()
		req := NewJsonRequest("POST", "/predict", requestBody)
		router.ServeHTTP(w, req)

		assert.Equal(t, 400, w.Code)
	})

	t.Run("RequiresTraceLevel", func(t *testing.T) {
		requestBody := &predictRequestBody{
			Architecture:          "x",
			BatchSize:             1,
			DesiredResultModality: "",
			Inputs:                []string{"x"},
			Model:                 1,
		}

		w := httptest.NewRecorder()
		req := NewJsonRequest("POST", "/predict", requestBody)
		router.ServeHTTP(w, req)

		assert.Equal(t, 400, w.Code)
	})

	t.Run("AcceptsValidRequest", func(t *testing.T) {
		api_mq.SetMessageQueue(api_mq.NullMessageQueue())
		requestBody := validPredictRequestBody()
		w := httptest.NewRecorder()
		req := NewJsonRequest("POST", "/predict", requestBody)
		router.ServeHTTP(w, req)

		res := predictResponseBody{}
		json.Unmarshal(w.Body.Bytes(), &res)

		assert.Equal(t, 200, w.Code)
		assert.Equal(t, "", res.ExperimentId)
	})

	t.Run("SendsMessageToAgentQueue", func(t *testing.T) {
		spy := &messageQueueSpy{}
		api_mq.SetMessageQueue(spy)
		requestBody := validPredictRequestBody()
		w := httptest.NewRecorder()
		req := NewJsonRequest("POST", "/predict", requestBody)
		router.ServeHTTP(w, req)
		message := messages.PredictByModelName{}
		err := json.Unmarshal([]byte(spy.channel.message), &message)

		assert.Nil(t, err)
		assert.Equal(t, "agent-pytorch-amd64", spy.publishChannel)
		assert.Equal(t, "test_model_1.0", message.ModelName)
	})

	t.Run("CreatesTrial", func(t *testing.T) {
		spy := &messageQueueSpy{
			correlationId: "trial1",
		}
		api_mq.SetMessageQueue(spy)
		requestBody := validPredictRequestBody()
		w := httptest.NewRecorder()
		req := NewJsonRequest("POST", "/predict", requestBody)
		router.ServeHTTP(w, req)

		trial, _ := testDb.GetTrialById("trial1")

		assert.NotNil(t, trial)
		assert.Equal(t, uint(1), trial.ModelID)
		assert.Equal(t, 1, len(trial.Inputs))
		assert.Equal(t, "input_url", trial.Inputs[0].URL)
	})
}
