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

func TestAuthorizeRefreshTokenHandler(t *testing.T) {
	env.InitEnv()

	cfg := env.InitConfig()
	tokenService := service.Service{
		Env:   cfg,
		Redis: cfg.InitRedis(),
	}

	request, _ := http.NewRequest("GET", "/access-token?userId=1", nil)
	response := httptest.NewRecorder()

	accessTokenHandler := NewAccessTokenHandler(&tokenService)
	accessTokenHandler.Generate(response, request)

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

	request, _ = http.NewRequest("POST", "/authorize-token", nil)
	request.Header.Set("Authorization", jsonJwt.AccessToken)
	request.Header.Set("Expires", jsonJwt.ExpirationTime)

	response = httptest.NewRecorder()

	authorizeTokenHandler := NewAuthorizeTokenHandler(&tokenService)
	authorizeTokenHandler.Auth(response, request)

	if response.Code != 200 {
		t.Fatalf("non-expected status code %v:\n\tbody: %v", "201", response.Code)
	}
}
