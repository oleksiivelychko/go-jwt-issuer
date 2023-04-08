package handlers

import (
	"encoding/json"
	"github.com/oleksiivelychko/go-jwt-issuer/config"
	"github.com/oleksiivelychko/go-jwt-issuer/issuer"
	"github.com/oleksiivelychko/go-jwt-issuer/middleware"
	"github.com/oleksiivelychko/go-jwt-issuer/service"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandler_ClearToken(t *testing.T) {
	tokenService := service.TokenService{
		Config:      config.NewConfig(),
		RedisClient: config.InitRedis(),
	}

	req, _ := http.NewRequest(http.MethodGet, "/access-token/?userID=1", nil)
	resp := httptest.NewRecorder()

	accessTokenHandler := NewAccessToken(&tokenService)
	accessTokenHandler.ServeHTTP(resp, req)

	if resp.Code != 201 {
		t.Fatalf("non-expected status code: %d\nbody: %v", resp.Code, resp.Body)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("unable to read response body: %s", err.Error())
	}

	responseJWT := &issuer.ResponseJWT{}
	err = json.Unmarshal(body, &responseJWT)
	if err != nil {
		t.Fatalf("unable to unmarshal response body: %s", err.Error())
	}

	req, _ = http.NewRequest(http.MethodPost, "/clear-token", nil)
	req.Header.Set("Authorization", responseJWT.AccessToken)
	req.Header.Set("Expires", responseJWT.ExpirationTime)
	resp = httptest.NewRecorder()

	clearTokenHandler := NewClearToken(&tokenService)
	middlewareHandler := middleware.JWT(clearTokenHandler)
	middlewareHandler.ServeHTTP(resp, req)

	if resp.Code != 200 {
		t.Fatalf("non-expected status code: %d\nbody: %v", resp.Code, resp.Body)
	}
}

func TestHandler_AuthorizeByRemovedToken(t *testing.T) {
	tokenService := service.TokenService{
		Config:      config.NewConfig(),
		RedisClient: config.InitRedis(),
	}

	req, _ := http.NewRequest(http.MethodGet, "/access-token/?userID=1", nil)
	resp := httptest.NewRecorder()

	accessTokenHandler := NewAccessToken(&tokenService)
	accessTokenHandler.ServeHTTP(resp, req)

	if resp.Code != 201 {
		t.Fatalf("non-expected status code: %d\nbody: %v", resp.Code, resp.Body)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("unable to read response body: %s", err.Error())
	}

	responseJWT := &issuer.ResponseJWT{}
	err = json.Unmarshal(body, &responseJWT)
	if err != nil {
		t.Fatalf("unable to unmarshal response body: %s", err.Error())
	}

	req, _ = http.NewRequest(http.MethodPost, "/clear-token", nil)
	req.Header.Set("Authorization", responseJWT.AccessToken)
	req.Header.Set("Expires", responseJWT.ExpirationTime)
	resp = httptest.NewRecorder()

	clearTokenHandler := NewClearToken(&tokenService)
	middlewareHandler := middleware.JWT(clearTokenHandler)
	middlewareHandler.ServeHTTP(resp, req)

	if resp.Code != 200 {
		t.Fatalf("non-expected status code: %d\nbody: %v", resp.Code, resp.Body)
	}

	req, _ = http.NewRequest(http.MethodPost, "/authorize-token", nil)
	req.Header.Set("Authorization", responseJWT.AccessToken)
	req.Header.Set("Expires", responseJWT.ExpirationTime)

	resp = httptest.NewRecorder()

	authorizeTokenHandler := NewAuthorizeToken(&tokenService)

	middlewareHandler = middleware.JWT(authorizeTokenHandler)
	middlewareHandler.ServeHTTP(resp, req)

	if resp.Code != 400 {
		t.Fatalf("non-expected status code: %d\nbody: %v", resp.Code, resp.Body)
	}

	body, err = io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("unable to read response body: %s", err.Error())
	}

	err = json.Unmarshal(body, &responseJWT)
	if err != nil {
		t.Fatalf("unable to unmarshal response body: %s", err.Error())
	}

	if responseJWT.ErrorMessage != "unable to validate cached token" {
		t.Fatalf("got error message: %s", responseJWT.ErrorMessage)
	}
}
