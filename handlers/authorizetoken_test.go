package handlers

import (
	"encoding/json"
	"github.com/oleksiivelychko/go-jwt-issuer/config"
	"github.com/oleksiivelychko/go-jwt-issuer/issuer"
	"github.com/oleksiivelychko/go-jwt-issuer/middleware"
	"github.com/oleksiivelychko/go-jwt-issuer/token"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandler_AuthorizeToken(t *testing.T) {
	tokenService := token.Service{
		Config:      config.NewConfig("secretkey", "jwt.account.local", "jwt.local", "1"),
		RedisClient: config.NewRedisClient("localhost", "6379", "secret"),
	}

	req, _ := http.NewRequest("GET", "/access-token?userID=1", nil)
	resp := httptest.NewRecorder()

	accessTokenHandler := NewAccessToken(&tokenService)
	accessTokenHandler.ServeHTTP(resp, req)

	if resp.Code != 201 {
		t.Fatalf("non-expected status code: %d\nbody: %v", resp.Code, resp.Body)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf(err.Error())
	}

	responseJWT := &issuer.ResponseJWT{}
	err = json.Unmarshal(body, &responseJWT)
	if err != nil {
		t.Fatalf(err.Error())
	}

	req, _ = http.NewRequest("POST", "/authorize-token", nil)
	req.Header.Set("Authorization", responseJWT.AccessToken)
	req.Header.Set("Expires", responseJWT.ExpirationTime)

	resp = httptest.NewRecorder()

	authorizeTokenHandler := NewAuthorizeToken(&tokenService)

	middlewareHandler := middleware.JWT(authorizeTokenHandler)
	middlewareHandler.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Fatalf("non-expected status code: %d\nbody: %v", resp.Code, resp.Body)
	}

	req, _ = http.NewRequest("POST", "/refresh-token", nil)
	req.Header.Set("Authorization", responseJWT.RefreshToken)
	req.Header.Set("Expires", responseJWT.ExpirationTime)

	resp = httptest.NewRecorder()

	refreshTokenHandler := NewRefreshToken(&tokenService)
	middlewareHandler = middleware.JWT(refreshTokenHandler)
	middlewareHandler.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Fatalf("non-expected status code: %d\nbody: %v", resp.Code, resp.Body)
	}

	body, err = io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf(err.Error())
	}

	err = json.Unmarshal(body, &responseJWT)
	if err != nil {
		t.Fatalf(err.Error())
	}

	if responseJWT.AccessToken == "" {
		t.Errorf("got empty accessToken")
	}

	if responseJWT.RefreshToken == "" {
		t.Errorf("got empty refreshToken")
	}

	if responseJWT.ExpirationTime == "" {
		t.Errorf("got empty expirationTime")
	}
}
