package endpoints

import (
	"api/api_mq"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPredictRoute(t *testing.T) {
	openDatabase()
	defer cleanupTestDatabase()
	router := SetupRoutes()

	t.Run("RequiresContentType", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/predict", strings.NewReader("{architecture: \"x\", framework: \"x\", model: \"x\", inputs: []}"))
		req.Header.Set("content-type", "application/json")
		router.ServeHTTP(w, req)

		assert.Equal(t, 400, w.Code)
	})

	t.Run("RequiresArchitecture", func(t *testing.T) {
		requestBody := &predictRequestBody{
			Architecture:          "",
			Framework:             "x",
			Inputs:                []string{"x"},
			Model:                 "x",
			DesiredResultModality: "x",
		}

		w := httptest.NewRecorder()
		req := NewJsonRequest("POST", "/predict", requestBody)
		router.ServeHTTP(w, req)

		assert.Equal(t, 400, w.Code)
	})

	t.Run("RequiresFramework", func(t *testing.T) {
		requestBody := &predictRequestBody{
			Architecture:          "x",
			Framework:             "",
			Inputs:                []string{"x"},
			Model:                 "x",
			DesiredResultModality: "x",
		}
		jsonBody, _ := json.Marshal(requestBody)

		w := httptest.NewRecorder()
		req := NewJsonRequest("POST", "/predict", jsonBody)
		router.ServeHTTP(w, req)

		assert.Equal(t, 400, w.Code)
	})

	t.Run("RequiresInputs", func(t *testing.T) {
		requestBody := &predictRequestBody{
			Architecture:          "x",
			Framework:             "x",
			Inputs:                []string{},
			Model:                 "x",
			DesiredResultModality: "x",
		}

		w := httptest.NewRecorder()
		req := NewJsonRequest("POST", "/predict", requestBody)
		router.ServeHTTP(w, req)

		assert.Equal(t, 400, w.Code)
	})

	t.Run("RequiresModel", func(t *testing.T) {
		requestBody := &predictRequestBody{
			Architecture:          "x",
			Framework:             "x",
			Inputs:                []string{"x"},
			Model:                 "",
			DesiredResultModality: "x",
		}

		w := httptest.NewRecorder()
		req := NewJsonRequest("POST", "/predict", requestBody)
		router.ServeHTTP(w, req)

		assert.Equal(t, 400, w.Code)
	})

	t.Run("RequiresDesiredResultModality", func(t *testing.T) {
		requestBody := &predictRequestBody{
			Architecture:          "x",
			Framework:             "x",
			Inputs:                []string{"x"},
			Model:                 "x",
			DesiredResultModality: "",
		}
		jsonBody, _ := json.Marshal(requestBody)

		w := httptest.NewRecorder()
		req := NewJsonRequest("POST", "/predict", jsonBody)
		router.ServeHTTP(w, req)

		assert.Equal(t, 400, w.Code)
	})

	t.Run("AcceptsValidRequest", func(t *testing.T) {
		api_mq.SetMessageQueue(api_mq.NullMessageQueue())
		router := SetupRoutes()
		requestBody := validPredictRequestBody("pytorch")
		w := httptest.NewRecorder()
		req := NewJsonRequest("POST", "/predict", requestBody)
		router.ServeHTTP(w, req)

		res := predictResponseBody{}
		json.Unmarshal(w.Body.Bytes(), &res)

		assert.Equal(t, 200, w.Code)
		assert.Equal(t, "1", res.ExperimentId)
	})

	t.Run("SendsMessageToAgentQueue", func(t *testing.T) {
		spy := &messageQueueSpy{}
		api_mq.SetMessageQueue(spy)
		router := SetupRoutes()
		requestBody := validPredictRequestBody("pytorch")
		w := httptest.NewRecorder()
		req := NewJsonRequest("POST", "/predict", requestBody)
		router.ServeHTTP(w, req)

		assert.Equal(t, "agent-pytorch-amd64", spy.publishChannel)
		assert.Equal(t, "do some work", spy.channel.message)
	})
}
