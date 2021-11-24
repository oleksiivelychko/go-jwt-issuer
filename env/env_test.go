package env

import (
	"github.com/oleksiivelychko/go-helper/in"
	"os"
	"testing"
)

func TestEnvConfig(t *testing.T) {
	_ = os.Setenv("SECRET_KEY", "secretkey")

	cfg := Config{
		Port:           GetPort(),
		SecretKey:      GetSecretKey(),
		ISS:            GetISS(),
		AUD:            GetAUD(),
		ExpiresMinutes: GetExpiresMinutes(),
	}

	var allowedRangePorts = []string{":80", ":8080", ":443"}
	_, result := in.StringIn(allowedRangePorts, cfg.Port)
	if !result {
		t.Errorf("given $PORT %s is not acceptable", cfg.Port)
	}

	if cfg.SecretKey == "" {
		t.Errorf("environment variable `SECRET_KEY` is not defined")
	}
}
