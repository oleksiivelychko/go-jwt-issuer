package env

import (
	"os"
	"strconv"
)

type Config struct {
	Port           string
	SecretKey      string
	AUD            string
	ISS            string
	ExpiresMinutes uint
}

func GetSecretKey() string {
	var secretKey = []byte(os.Getenv("SECRET_KEY"))
	if len(secretKey) == 0 {
		return ""
	}

	return string(secretKey)
}

func GetAUD() string {
	return os.Getenv("AUDIENCE_AUD")
}

func GetISS() string {
	return os.Getenv("ISSUER_ISS")
}

func GetExpiresMinutes() uint {
	expirationTimeExp := os.Getenv("EXPIRES_MINUTES")
	if expirationTimeExp != "" {
		minutes, err := strconv.ParseUint(expirationTimeExp, 10, 32)
		if err == nil {
			return uint(minutes)
		}
	}
	return 60
}

func GetPort() string {
	var port = []byte(os.Getenv("PORT"))
	if len(port) == 0 {
		return ":8080"
	}

	return ":" + string(port)
}
