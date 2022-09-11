package env

import "os"

type JsonJwt struct {
	AccessToken    string `json:"access-token,omitempty"`
	RefreshToken   string `json:"refresh-token,omitempty"`
	ExpirationTime string `json:"expiration-time,omitempty"`
	ErrorMessage   string `json:"error-message,omitempty"`
	ErrorCode      uint8  `json:"error-code,omitempty"`
}

func SetDefaults() {
	if os.Getenv("HOST") == "" {
		_ = os.Setenv("HOST", "localhost")
	}
	if os.Getenv("PORT") == "" {
		_ = os.Setenv("PORT", "8080")
	}
	if os.Getenv("REDIS_HOST") == "" {
		_ = os.Setenv("REDIS_HOST", "localhost")
	}
	if os.Getenv("REDIS_PORT") == "" {
		_ = os.Setenv("REDIS_PORT", "6379")
	}
	if os.Getenv("REDIS_PASSWORD") == "" {
		_ = os.Setenv("REDIS_PASSWORD", "secret")
	}
	if os.Getenv("SECRET_KEY") == "" {
		_ = os.Setenv("SECRET_KEY", "secretkey")
	}
	if os.Getenv("EXPIRES_MINUTES") == "" {
		_ = os.Setenv("EXPIRES_MINUTES", "1")
	}
	if os.Getenv("AUDIENCE_AUD") == "" {
		_ = os.Setenv("AUDIENCE_AUD", "account.jwt.local")
	}
	if os.Getenv("ISSUER_ISS") == "" {
		_ = os.Setenv("ISSUER_ISS", "jwt.local")
	}
}
