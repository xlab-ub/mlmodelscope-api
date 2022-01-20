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

func TestFrameworkRoutes(t *testing.T) {
	openDatabase()
	defer cleanupTestDatabase()
	router := SetupRoutes()
	req, _ := http.NewRequest("GET", "/frameworks", nil)

	t.Run("ListEmpty", func(t *testing.T) {
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, 200, w.Code)
		assert.Equal(t, "{}", w.Body.String())
	})

	t.Run("ListNotEmpty", func(t *testing.T) {
		testDb.CreateFramework(&models.Framework{Name: "fw1", Architectures: []models.Architecture{{Name: "amd64"}}})
		testDb.CreateFramework(&models.Framework{Name: "fw2", Architectures: []models.Architecture{{Name: "arm"}}})

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		var result ListFrameworksResponse
		_ = json.Unmarshal(w.Body.Bytes(), &result)

		assert.Equal(t, uint(1), result.Frameworks[0].ID)
		assert.Equal(t, "fw1", result.Frameworks[0].Name)
		assert.Equal(t, "amd64", result.Frameworks[0].Architectures[0].Name)
		assert.Equal(t, uint(2), result.Frameworks[1].ID)
		assert.Equal(t, "fw2", result.Frameworks[1].Name)
		assert.Equal(t, "arm", result.Frameworks[1].Architectures[0].Name)
	})
}
