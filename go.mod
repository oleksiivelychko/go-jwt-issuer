module github.com/oleksiivelychko/go-jwt-issuer

// +heroku goVersion go1.16
go 1.17

require (
	github.com/go-redis/redis/v8 v8.11.4
	github.com/golang-jwt/jwt v3.2.2+incompatible
	github.com/joho/godotenv v1.4.0
	github.com/oleksiivelychko/go-helper v0.0.0-20211111054117-9aa512a7810b
)

require (
	github.com/cespare/xxhash/v2 v2.1.2 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
)
