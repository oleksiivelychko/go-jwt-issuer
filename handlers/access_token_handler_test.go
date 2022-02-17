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

func TestAccessTokenHandler(t *testing.T) {
	env.InitEnv()

	cfg := env.InitConfig()
	tokenService := service.Service{
		Env:   cfg,
		Redis: cfg.InitRedis(),
	}

	request, _ := http.NewRequest("GET", "/access-token/?userId=1", nil)
	response := httptest.NewRecorder()

	accessTokenHandler := NewAccessTokenHandler(&tokenService)
	accessTokenHandler.ServeHTTP(response, request)

	if response.Code != 201 {
		t.Fatalf("non-expected status code: %d\nbody: %v", response.Code, response.Body)
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
