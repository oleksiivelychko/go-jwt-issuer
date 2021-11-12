package env

import (
	"github.com/oleksiivelychko/go-helper/in"
	"os"
	"testing"
)

func TestGetSecretKey(t *testing.T) {
	_ = os.Setenv("SECRET_KEY", "secretkey")
	var secretKey = GetSecretKey()
	if string(secretKey) == "" {
		t.Errorf("environment variable `SECRET_KEY` is not defined")
	}
}

func TestGetPort(t *testing.T) {
	var port = GetPort()
	var allowedRangePorts = []string{":80", ":8080", ":443"}
	_, result := in.StringIn(allowedRangePorts, port)
	if !result {
		t.Errorf("given $PORT as %s is not included in allowed range", port)
	}
}
