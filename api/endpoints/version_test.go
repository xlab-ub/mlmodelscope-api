// +build !integration

package endpoints

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVersionRoute(t *testing.T) {
	router := SetupRoutes()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/", nil)
	router.ServeHTTP(w, req)
	res := versionResponse{}
	json.Unmarshal(w.Body.Bytes(), &res)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, version, res.Version)
}