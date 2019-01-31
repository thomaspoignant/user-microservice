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
	config.LoadConfigFile()
}

// utility func who make the request
func PerformHTTPRequest(r http.Handler, method, path string, body io.Reader) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(method, path, body)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}
