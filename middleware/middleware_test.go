package middleware

import (
	"github.com/oleksiivelychko/go-jwt-issuer/config"
	"github.com/oleksiivelychko/go-jwt-issuer/issuer"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
)

func TestMiddleware_ValidateJWT(t *testing.T) {
	var secretKey = config.GetSecretKey()
	var aud = config.GetAudience()
	var iss = config.GetIssuer()
	var expiresMinutes = config.GetExpirationTimeMinutes()

	token, _, exp, _ := issuer.IssueJWT(secretKey, aud, iss, expiresMinutes, 1)

	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", token)
	req.Header.Set("Expires", strconv.FormatInt(exp, 10))

	resp := httptest.NewRecorder()

	middlewareHandler := JWT(nextHandler)
	middlewareHandler.ServeHTTP(resp, req)

	respMessage := string(resp.Body.Bytes())

	if respMessage == "environment variable 'SECRET_KEY' is not defined" {
		t.Errorf("response code: %d, message: %s", resp.Code, respMessage)
	}
	if respMessage == "unable to get token from 'Authorization' header" {
		t.Errorf("response code: %d, message: %s", resp.Code, respMessage)
	}
	if strings.HasPrefix(respMessage, "unexpected signing method") {
		t.Errorf("response code: %d, message: %s", resp.Code, respMessage)
	}
	if respMessage == "unable to verify 'aud' claim" {
		t.Errorf("response code: %d, message: %s", resp.Code, respMessage)
	}
	if respMessage == "unable to verify 'iss' claim" {
		t.Errorf("response code: %d, message: %s", resp.Code, respMessage)
	}
	if respMessage == "unable to verify 'exp' claim" {
		t.Errorf("response code: %d, message: %s", resp.Code, respMessage)
	}
}
