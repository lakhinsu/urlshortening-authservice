package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/lakhinsu/urlshortening-authservice/app"
	"github.com/lakhinsu/urlshortening-authservice/utils"
	"github.com/stretchr/testify/assert"
)

func init() {
	gin.SetMode(gin.TestMode)
	//zerolog.SetGlobalLevel(zerolog.Disabled)
}

func TestLoginSuccess(t *testing.T) {
	router := app.SetupApp()
	w := httptest.NewRecorder()

	testUser := utils.GetEnvVar("TEST_USER")
	testUserPassword := utils.GetEnvVar("TEST_USER_PASSWORD")

	b := []byte(fmt.Sprintf(`{"username": "%s", "password":"%s"}`, testUser, testUserPassword))

	reader := bytes.NewReader(b)

	req, _ := http.NewRequest("POST", "/v1/login", reader)
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	router.ServeHTTP(w, req)

	body, _ := io.ReadAll(w.Result().Body)
	var f interface{}
	json.Unmarshal(body, &f)

	myMap := f.(map[string]interface{})
	// validate response code
	assert.Equal(t, 200, w.Code)
	// validate token
	tokenResponse := fmt.Sprintf("%s", myMap["token"])
	jwtSecret := []byte(utils.GetEnvVar("JWT_SECRET"))
	token, err := jwt.Parse(tokenResponse, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return jwtSecret, nil
	})
	if err != nil || !token.Valid {
		fmt.Println("Could not verify the token")
		t.Fail()
	}
}

func TestInvalidCredentials(t *testing.T) {
	router := app.SetupApp()
	w := httptest.NewRecorder()

	testUser := utils.GetEnvVar("TEST_USER")
	testUserPassword := "invalid"

	b := []byte(fmt.Sprintf(`{"username": "%s", "password":"%s"}`, testUser, testUserPassword))

	reader := bytes.NewReader(b)

	req, _ := http.NewRequest("POST", "/v1/login", reader)
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	router.ServeHTTP(w, req)

	// validate response code
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestUnproccessableEntity(t *testing.T) {
	router := app.SetupApp()
	w := httptest.NewRecorder()

	testUser := utils.GetEnvVar("TEST_USER")
	testUserPassword := utils.GetEnvVar("TEST_USER_PASSWORD")

	b := []byte(fmt.Sprintf(`{"user": "%s", "password":"%s"}`, testUser, testUserPassword))

	reader := bytes.NewReader(b)

	req, _ := http.NewRequest("POST", "/v1/login", reader)
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	router.ServeHTTP(w, req)

	// validate response code
	assert.Equal(t, http.StatusUnprocessableEntity, w.Code)
}
