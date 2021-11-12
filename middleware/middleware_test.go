package middleware

import (
	"github.com/oleksiivelychko/go-jwt-issuer/env"
	"github.com/oleksiivelychko/go-jwt-issuer/issuer"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestAllowToEndpointMiddleware(t *testing.T) {
	_ = os.Setenv("SECRET_KEY", "secretkey")
	token, _ := issuer.IssueJWT(env.GetSecretKey(), "", "")

	closure := func(writer http.ResponseWriter, request *http.Request) {}

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", token)

	res := httptest.NewRecorder()

	closure(res, req)

	handler := AllowToEndpoint(closure)
	handler.ServeHTTP(res, req)

	if string(res.Body.Bytes()) == "environment variable `SECRET_KEY` is not defined" {
		t.Errorf("failed to get access to resource using token")
	}
}

func TestJwtAuthenticationMiddleware(t *testing.T) {
	_ = os.Setenv("SECRET_KEY", "secretkey")
	token, _ := issuer.IssueJWT(env.GetSecretKey(), "", "")

	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", token)

	res := httptest.NewRecorder()

	handler := JwtAuthentication(nextHandler)
	handler.ServeHTTP(res, req)

	if string(res.Body.Bytes()) == "environment variable `SECRET_KEY` is not defined" {
		t.Errorf("failed to get access to resource using token")
	}
}
