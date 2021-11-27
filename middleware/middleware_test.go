package middleware

import (
	"github.com/oleksiivelychko/go-jwt-issuer/env"
	"github.com/oleksiivelychko/go-jwt-issuer/issuer"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"
)

func initEnv() {
	_ = os.Setenv("SECRET_KEY", "secretkey")
	_ = os.Setenv("EXPIRES_MINUTES", "10")
	_ = os.Setenv("AUDIENCE_AUD", "account.jwt.local")
	_ = os.Setenv("ISSUER_ISS", "jwt.local")
}

func TestAllowToEndpointMiddleware(t *testing.T) {
	initEnv()

	var secretKey = env.GetSecretKey()
	var aud = env.GetAUD()
	var iss = env.GetISS()
	var expiresMinutes = env.GetExpiresMinutes()

	token, _, exp, _ := issuer.IssueUserJWT(secretKey, aud, iss, expiresMinutes, 1)

	// to validate expiration time
	time.Sleep(1 * time.Second)

	closure := func(writer http.ResponseWriter, request *http.Request) {}

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", token)
	req.Header.Set("ExpirationTime", strconv.FormatInt(exp, 10))

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
	initEnv()

	var secretKey = env.GetSecretKey()
	var aud = env.GetAUD()
	var iss = env.GetISS()
	var expiresMinutes = env.GetExpiresMinutes()

	token, _, exp, _ := issuer.IssueUserJWT(secretKey, aud, iss, expiresMinutes, 1)

	// to validate expiration time
	time.Sleep(1 * time.Second)

	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", token)
	req.Header.Set("ExpirationTime", strconv.FormatInt(exp, 10))

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
