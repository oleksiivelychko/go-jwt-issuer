package handlers

import (
	"encoding/json"
	"github.com/oleksiivelychko/go-jwt-issuer/env"
	"github.com/oleksiivelychko/go-jwt-issuer/middleware"
	"github.com/oleksiivelychko/go-jwt-issuer/service"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"
)

func TestRefreshTokenHandler(t *testing.T) {
	env.SetDefaults()

	cfg := env.InitConfig()
	tokenService := service.Service{
		Env:   cfg,
		Redis: cfg.InitRedis(),
	}

	request, _ := http.NewRequest("GET", "/access-token?userId=1", nil)
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

	request, _ = http.NewRequest("POST", "/refresh-token", nil)
	request.Header.Set("Authorization", jsonJwt.RefreshToken)
	request.Header.Set("Expires", jsonJwt.ExpirationTime)

	response = httptest.NewRecorder()

	refreshTokenHandler := NewRefreshTokenHandler(&tokenService)
	handler := middleware.JWT(refreshTokenHandler)
	handler.ServeHTTP(response, request)

	if response.Code != 200 {
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

	if jsonJwt.AccessToken == "" {
		t.Fatalf("got empty `access_token`")
	}

	if jsonJwt.RefreshToken == "" {
		t.Fatalf("got empty `refresh_token`")
	}

	if jsonJwt.ExpirationTime == "" {
		t.Fatalf("got empty `expiration_time`")
	}
}

func TestAuthorizeByExpiredTokenHandler(t *testing.T) {
	env.SetDefaults()

	cfg := env.InitConfig()
	tokenService := service.Service{
		Env:   cfg,
		Redis: cfg.InitRedis(),
	}

	request, _ := http.NewRequest(http.MethodGet, "/access-token?userId=1", nil)
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

	request, _ = http.NewRequest(http.MethodPost, "/authorize-token", nil)
	request.Header.Set("Authorization", jsonJwt.AccessToken)
	request.Header.Set("Expires", jsonJwt.ExpirationTime)

	response = httptest.NewRecorder()

	authorizeTokenHandler := NewAuthorizeTokenHandler(&tokenService)

	handler := middleware.JWT(authorizeTokenHandler)
	handler.ServeHTTP(response, request)

	if response.Code != 200 {
		t.Fatalf("non-expected status code: %d\nbody: %v", response.Code, response.Body)
	}

	exp, err := strconv.ParseInt(jsonJwt.ExpirationTime, 10, 64)
	ts := time.Unix(exp, 0).Add(time.Second * time.Duration(1)).Unix()
	expiredTs := strconv.FormatInt(ts, 10)

	request, _ = http.NewRequest(http.MethodPost, "/refresh-token", nil)
	request.Header.Set("Authorization", jsonJwt.RefreshToken)
	request.Header.Set("Expires", expiredTs)

	response = httptest.NewRecorder()

	refreshTokenHandler := NewRefreshTokenHandler(&tokenService)
	handler = middleware.JWT(refreshTokenHandler)
	handler.ServeHTTP(response, request)

	if response.Code != 401 {
		t.Fatalf("non-expected status code: %d\nbody: %v", response.Code, response.Body)
	}

	body, err = ioutil.ReadAll(response.Body)
	if err != nil {
		t.Fatalf("unable to read response body: %s", err.Error())
	}

	if string(body) != "failed to verify `exp` claim\n" {
		t.Fatalf("got valid `expiration_time`")
	}
}
