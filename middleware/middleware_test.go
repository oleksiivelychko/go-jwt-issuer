package middleware

import (
	"github.com/oleksiivelychko/go-jwt-issuer/issuer"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"strings"
	"testing"
)

func TestMiddleware_ValidateJWT(t *testing.T) {
	_ = os.Setenv("SECRET_KEY", "secretkey")
	token, _, expMinutes, _ := issuer.IssueJWT("secretkey", "", "", 0, 1)

	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", token)
	req.Header.Set("Expires", strconv.FormatInt(expMinutes, 10))

	resp := httptest.NewRecorder()

	middlewareHandler := JWT(nextHandler)
	middlewareHandler.ServeHTTP(resp, req)

	respMessage := string(resp.Body.Bytes())

	if respMessage == "environment variable 'SECRET_KEY' is not defined\n" {
		t.Errorf("response code: %d, message: %s", resp.Code, respMessage)
	}
	if respMessage == "unable to get token from 'Authorization' header\n" {
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

	_ = os.Unsetenv("SECRET_KEY")
}
