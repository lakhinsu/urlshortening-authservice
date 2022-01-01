package tests

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/lakhinsu/urlshortening-authservice/app"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
)

func init() {
	gin.SetMode(gin.TestMode)
	zerolog.SetGlobalLevel(zerolog.Disabled)
}

func TestPingRoute(t *testing.T) {
	router := app.SetupApp()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/v1/", nil)
	router.ServeHTTP(w, req)

	body, _ := io.ReadAll(w.Result().Body)
	var f interface{}
	json.Unmarshal(body, &f)

	myMap := f.(map[string]interface{})
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, myMap["message"], "Hello world")
}
