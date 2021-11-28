package env

import "os"

type JsonJwt struct {
	AccessToken    string `json:"access-token"`
	RefreshToken   string `json:"refresh-token"`
	ExpirationTime string `json:"expiration-time"`
}

func InitEnv() {
	_ = os.Setenv("REDIS_HOST", "localhost")
	_ = os.Setenv("REDIS_PORT", "6379")
	_ = os.Setenv("REDIS_PASSWORD", "secret")
	_ = os.Setenv("REDIS_DB", "0")
	_ = os.Setenv("SECRET_KEY", "secretkey")
	_ = os.Setenv("EXPIRES_MINUTES", "5")
	_ = os.Setenv("AUDIENCE_AUD", "account.jwt.local")
	_ = os.Setenv("ISSUER_ISS", "jwt.local")
}
