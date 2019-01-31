package api

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/thomaspoignant/user-microservice/testingUtils"
)

// TestHealthCheck check if the health check is answering correctly
func TestHealthCheck(t *testing.T) {
	testingUtils.PrepareTest()
	router := SetupRouter()

	w := testingUtils.PerformHTTPRequest(router, "GET", "/health", nil)
	assert.Equal(t, http.StatusOK, w.Code)

	expected := healthCheck{
		Health: "RUNNING",
		Code:   Success,
	}

	var response healthCheck
	err := json.Unmarshal([]byte(w.Body.String()), &response)
	assert.Nil(t, err)
	assert.Equal(t, expected, response)
}
