package testingUtils

import (
	"io"
	"net/http"
	"net/http/httptest"
	"os"

	"github.com/thomaspoignant/user-microservice/config"
)

// PrepareTest set up the environnement before the test
func PrepareTest() {
	os.Setenv("TEST", "true") // allow to load correct config file
	os.Setenv("ENV", "")      // allow to set the env config file

	os.Setenv("APP_PORT", "8080")
	os.Setenv("GIN_MODE", "debug")
	os.Setenv("RUNNING_MODE", "test")
	os.Setenv("AWS_REGION", "eu-west-1")
	os.Setenv("DYNAMODB_ENDPOINT", "http://localhost:9000")
	os.Setenv("DYNAMODB_TABLE_NAME", "user_test")

	config.LoadConfigFile()
}

// PerformHTTPRequest utility func who make the request
func PerformHTTPRequest(r http.Handler, method, path string, body io.Reader) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(method, path, body)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}
