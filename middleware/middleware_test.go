package middleware

import (
	"github.com/oleksiivelychko/go-jwt-issuer/env"
	"github.com/oleksiivelychko/go-jwt-issuer/issuer"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
)

func TestJWTMiddleware(t *testing.T) {
	env.InitEnv()

	var secretKey = env.GetSecretKey()
	var aud = env.GetAUD()
	var iss = env.GetISS()
	var expiresMinutes = env.GetExpiresMinutes()

	token, _, exp, _ := issuer.IssueUserJWT(secretKey, aud, iss, expiresMinutes, 1)

	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", token)
	req.Header.Set("Expires", strconv.FormatInt(exp, 10))

	res := httptest.NewRecorder()

	handler := JWT(nextHandler)
	handler.ServeHTTP(res, req)

	responseMessage := string(res.Body.Bytes())
	if responseMessage == "environment variable `SECRET_KEY` is not defined" {
		t.Errorf("status: %d, message: %s", res.Code, responseMessage)
	}
	if responseMessage == "failed to get token from header request" {
		t.Errorf("status: %d, message: %s", res.Code, responseMessage)
	}
	if strings.HasPrefix(responseMessage, "unexpected signing method") {
		t.Errorf("status: %d, message: %s", res.Code, responseMessage)
	}
	if responseMessage == "failed to verify `aud` claim" {
		t.Errorf("status: %d, message: %s", res.Code, responseMessage)
	}
	if responseMessage == "failed to verify `iss` claim" {
		t.Errorf("status: %d, message: %s", res.Code, responseMessage)
	}
	if responseMessage == "failed to verify `exp` claim" {
		t.Errorf("status: %d, message: %s", res.Code, responseMessage)
	}
}
