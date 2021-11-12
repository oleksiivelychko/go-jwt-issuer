package env

import (
	"os"
)

const defaultBindPort = ":8080"

func GetSecretKey() []byte {
	var secretKey = []byte(os.Getenv("SECRET_KEY"))
	if len(secretKey) == 0 {
		return []byte("")
	}

	return secretKey
}

func GetAUD() string {
	return os.Getenv("AUDIENCE_AUD")
}

func GetISS() string {
	return os.Getenv("ISSUER_ISS")
}

func GetPort() string {
	var port = []byte(os.Getenv("PORT"))
	if len(port) == 0 {
		return defaultBindPort
	}

	return ":" + string(port)
}
