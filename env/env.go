package env

import (
	"os"
	"strconv"
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

func GetEXP() int {
	expirationTimeExp := os.Getenv("EXPIRATION_TIME_EXP")
	if expirationTimeExp != "" {
		exp, err := strconv.ParseInt(expirationTimeExp, 10, 32)
		if err == nil {
			return int(exp)
		}
	}
	return 60
}

func GetPort() string {
	var port = []byte(os.Getenv("PORT"))
	if len(port) == 0 {
		return defaultBindPort
	}

	return ":" + string(port)
}
