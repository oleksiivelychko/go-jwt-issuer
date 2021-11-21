package middleware

import (
	"github.com/oleksiivelychko/go-jwt-issuer/env"
	"github.com/oleksiivelychko/go-jwt-issuer/issuer"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

func TestAllowToEndpointMiddleware(t *testing.T) {
	_ = os.Setenv("SECRET_KEY", "secretkey")
	_ = os.Setenv("EXPIRES_MINUTES", "10")

	var secretKey = env.GetSecretKey()
	var aud = env.GetAUD()
	var iss = env.GetISS()
	var expiresMinutes = env.GetExpiresMinutes()

	token, _, _, _ := issuer.IssueUserJWT(secretKey, aud, iss, expiresMinutes, 1)

	closure := func(writer http.ResponseWriter, request *http.Request) {}

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", token)

	res := httptest.NewRecorder()

	closure(res, req)

	handler := AllowToEndpoint(closure)
	handler.ServeHTTP(res, req)

	if string(res.Body.Bytes()) == "environment variable `SECRET_KEY` is not defined" {
		t.Errorf(string(res.Body.Bytes()))
	}

	if string(res.Body.Bytes()) == "failed to get token from header request" {
		t.Errorf(string(res.Body.Bytes()))
	}

	if strings.HasPrefix(string(res.Body.Bytes()), "unexpected signing method") {
		t.Errorf(string(res.Body.Bytes()))
	}

	if string(res.Body.Bytes()) == "failed to verify `aud` claim" {
		t.Errorf(string(res.Body.Bytes()))
	}

	if string(res.Body.Bytes()) == "failed to verify `iss` claim" {
		t.Errorf(string(res.Body.Bytes()))
	}

	if string(res.Body.Bytes()) == "failed to verify `exp` claim" {
		t.Errorf(string(res.Body.Bytes()))
	}
}

func TestJwtAuthenticationMiddleware(t *testing.T) {
	_ = os.Setenv("SECRET_KEY", "secretkey")
	_ = os.Setenv("EXPIRES_MINUTES", "10")

	var secretKey = env.GetSecretKey()
	var aud = env.GetAUD()
	var iss = env.GetISS()
	var expiresMinutes = env.GetExpiresMinutes()

	token, _, _, _ := issuer.IssueUserJWT(secretKey, aud, iss, expiresMinutes, 1)

	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", token)

	res := httptest.NewRecorder()

	handler := JwtAuthentication(nextHandler)
	handler.ServeHTTP(res, req)

	if string(res.Body.Bytes()) == "environment variable `SECRET_KEY` is not defined" {
		t.Errorf(string(res.Body.Bytes()))
	}

	if string(res.Body.Bytes()) == "failed to get token from header request" {
		t.Errorf(string(res.Body.Bytes()))
	}

	if strings.HasPrefix(string(res.Body.Bytes()), "unexpected signing method") {
		t.Errorf(string(res.Body.Bytes()))
	}

	if string(res.Body.Bytes()) == "failed to verify `aud` claim" {
		t.Errorf(string(res.Body.Bytes()))
	}

	if string(res.Body.Bytes()) == "failed to verify `iss` claim" {
		t.Errorf(string(res.Body.Bytes()))
	}

	if string(res.Body.Bytes()) == "failed to verify `exp` claim" {
		t.Errorf(string(res.Body.Bytes()))
	}
}
