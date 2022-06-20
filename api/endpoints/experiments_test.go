// +build !integration

package endpoints

import (
	"api/db/models"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestExperimentRoute(t *testing.T) {
	openDatabase()
	createTestModelAndFramework()
	createTestExperiment()
	defer cleanupTestDatabase()
	router := SetupRoutes()

	t.Run("GetMissingExperiment", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/experiments/does_not_exist", nil)

		router.ServeHTTP(w, req)

		assert.Equal(t, 404, w.Code)
	})

	t.Run("GetEmptyExperiment", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/experiments/test", nil)

		router.ServeHTTP(w, req)
		var experiment models.Experiment
		err := json.Unmarshal(w.Body.Bytes(), &experiment)

		assert.Nil(t, err)
		assert.Equal(t, 0, len(experiment.Trials))
	})

	t.Run("GetExperimentWithTrials", func(t *testing.T) {
		testDb.CreateTrial(&models.Trial{ID: "trial1", ModelID: 1, ExperimentID: "test"})
		testDb.CreateTrial(&models.Trial{ID: "trial2", ModelID: 1, ExperimentID: "test"})
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/experiments/test", nil)

		router.ServeHTTP(w, req)
		var experiment models.Experiment
		err := json.Unmarshal(w.Body.Bytes(), &experiment)

		assert.Nil(t, err)
		assert.Equal(t, 2, len(experiment.Trials))
	})
}
