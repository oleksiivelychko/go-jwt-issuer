package handlers

import (
	"encoding/json"
	"github.com/oleksiivelychko/go-jwt-issuer/config"
	"github.com/oleksiivelychko/go-jwt-issuer/issuer"
	"github.com/oleksiivelychko/go-jwt-issuer/service"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandler_AccessToken(t *testing.T) {
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
