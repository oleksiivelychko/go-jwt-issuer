package handlers

import (
	"encoding/json"
	"github.com/oleksiivelychko/go-jwt-issuer/env"
	"github.com/oleksiivelychko/go-jwt-issuer/middleware"
	"github.com/oleksiivelychko/go-jwt-issuer/service"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestClearTokenHandler(t *testing.T) {
	env.SetDefaults()

	cfg := env.InitConfig()
	tokenService := service.Service{
		Env:   cfg,
		Redis: cfg.InitRedis(),
	}

	request, _ := http.NewRequest(http.MethodGet, "/access-token/?userId=1", nil)
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

	request, _ = http.NewRequest(http.MethodPost, "/clear-token", nil)
	request.Header.Set("Authorization", jsonJwt.AccessToken)
	request.Header.Set("Expires", jsonJwt.ExpirationTime)
	response = httptest.NewRecorder()

	clearTokenHandler := NewClearTokenHandler(&tokenService)
	handler := middleware.JWT(clearTokenHandler)
	handler.ServeHTTP(response, request)

	if response.Code != 200 {
		t.Fatalf("non-expected status code: %d\nbody: %v", response.Code, response.Body)
	}
}

func TestAuthorizeByRemovedTokenHandler(t *testing.T) {
	env.SetDefaults()

	cfg := env.InitConfig()
	tokenService := service.Service{
		Env:   cfg,
		Redis: cfg.InitRedis(),
	}

	request, _ := http.NewRequest(http.MethodGet, "/access-token/?userId=1", nil)
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

	request, _ = http.NewRequest(http.MethodPost, "/clear-token", nil)
	request.Header.Set("Authorization", jsonJwt.AccessToken)
	request.Header.Set("Expires", jsonJwt.ExpirationTime)
	response = httptest.NewRecorder()

	clearTokenHandler := NewClearTokenHandler(&tokenService)
	handler := middleware.JWT(clearTokenHandler)
	handler.ServeHTTP(response, request)

	if response.Code != 200 {
		t.Fatalf("non-expected status code: %d\nbody: %v", response.Code, response.Body)
	}

	request, _ = http.NewRequest(http.MethodPost, "/authorize-token", nil)
	request.Header.Set("Authorization", jsonJwt.AccessToken)
	request.Header.Set("Expires", jsonJwt.ExpirationTime)

	response = httptest.NewRecorder()

	authorizeTokenHandler := NewAuthorizeTokenHandler(&tokenService)

	handler = middleware.JWT(authorizeTokenHandler)
	handler.ServeHTTP(response, request)

	if response.Code != 400 {
		t.Fatalf("non-expected status code: %d\nbody: %v", response.Code, response.Body)
	}

	body, err = ioutil.ReadAll(response.Body)
	if err != nil {
		t.Fatalf("unable to read response body: %s", err.Error())
	}

	err = json.Unmarshal(body, &jsonJwt)
	if err != nil {
		t.Fatalf("unable to unmarshal response body: %s", err.Error())
	}

	if jsonJwt.ErrorMessage != "unable to validate cached token" {
		t.Fatalf("got valid cached token")
	}
}
