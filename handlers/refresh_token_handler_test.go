package handlers

import (
	"encoding/json"
	"github.com/oleksiivelychko/go-jwt-issuer/env"
	"github.com/oleksiivelychko/go-jwt-issuer/service"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRefreshTokenHandler(t *testing.T) {
	env.InitEnv()

	cfg := env.InitConfig()
	tokenService := service.Service{
		Env:   cfg,
		Redis: cfg.InitRedis(),
	}

	request, _ := http.NewRequest("GET", "/access-token?userId=1", nil)
	response := httptest.NewRecorder()

	AccessTokenHandler(&tokenService)(response, request)

	if response.Code != 201 {
		t.Fatalf("non-expected status code %v:\n\tbody: %v", "201", response.Code)
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		t.Fatalf("unable to read response body: %s", err.Error())
	}

	jsonJwt := &env.JsonJwt{}
	err = json.Unmarshal(body, &jsonJwt)
	if err != nil {
		t.Fatalf("unable to unmarshal response body: %s", err.Error())
	}

	request, _ = http.NewRequest("GET", "/refresh-token?userId=1", nil)
	request.Header.Set("Authorization", jsonJwt.RefreshToken)
	request.Header.Set("ExpirationTime", jsonJwt.ExpirationTime)

	response = httptest.NewRecorder()

	RefreshTokenHandler(&tokenService)(response, request)

	if response.Code != 201 {
		t.Fatalf("non-expected status code %v:\n\tbody: %v", "201", response.Code)
	}

	body, err = ioutil.ReadAll(response.Body)
	if err != nil {
		t.Fatalf("unable to read response body: %s", err.Error())
	}

	err = json.Unmarshal(body, &jsonJwt)
	if err != nil {
		t.Fatalf("unable to unmarshal response body: %s", err.Error())
	}

	if jsonJwt.AccessToken == "" {
		t.Fatalf("got empty `access-token`")
	}

	if jsonJwt.RefreshToken == "" {
		t.Fatalf("got empty `refresh-token`")
	}

	if jsonJwt.ExpirationTime == "" {
		t.Fatalf("got empty `expiration-time`")
	}
}
