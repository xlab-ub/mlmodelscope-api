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

func TestModelRoutes(t *testing.T) {
	openDatabase()
	defer cleanupTestDatabase()
	router := SetupRoutes()
	req, _ := http.NewRequest("GET", "/models", nil)

	t.Run("ListEmpty", func(t *testing.T) {
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, 200, w.Code)
		assert.Equal(t, "{}", w.Body.String())
	})

	t.Run("ListNotEmpty", func(t *testing.T) {
		testDb.CreateModel(&models.Model{Name: "model1", Framework: &models.Framework{Name: "fw1"}, Output: models.ModelOutput{Type: "test"}})
		testDb.CreateModel(&models.Model{Name: "model2", Framework: &models.Framework{Name: "fw2"}})

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		var result ModelListResponse
		_ = json.Unmarshal(w.Body.Bytes(), &result)

		assert.Equal(t, "model1", result.Models[0].Name)
		assert.Equal(t, uint(1), result.Models[0].ID)
		assert.Equal(t, "fw1", result.Models[0].Framework.Name)
		assert.Equal(t, "model2", result.Models[1].Name)
		assert.Equal(t, uint(2), result.Models[1].ID)
		assert.Equal(t, "fw2", result.Models[1].Framework.Name)
	})

	t.Run("ListByFrameworkId", func(t *testing.T) {
		req, _ = http.NewRequest("GET", "/models/framework/2", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		var result ModelListResponse
		_ = json.Unmarshal(w.Body.Bytes(), &result)

		assert.Equal(t, uint(2), result.Models[0].ID)
		assert.Equal(t, "model2", result.Models[0].Name)
		assert.Equal(t, "fw2", result.Models[0].Framework.Name)
	})

	t.Run("ListByFrameworkId_BadRequest", func(t *testing.T) {
		req, _ = http.NewRequest("GET", "/models/framework/x", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, 400, w.Code)
		assert.Equal(t, "{\"error\":\"invalid Framework ID: x\"}", w.Body.String())
	})

	t.Run("ListByTask", func(t *testing.T) {
		req, _ = http.NewRequest("GET", "/models/task/test", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		var result ModelListResponse
		_ = json.Unmarshal(w.Body.Bytes(), &result)

		assert.Equal(t, uint(1), result.Models[0].ID)
		assert.Equal(t, "model1", result.Models[0].Name)
		assert.Equal(t, "fw1", result.Models[0].Framework.Name)
		assert.Equal(t, "test", result.Models[0].Output.Type)
	})
}
