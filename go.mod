module github.com/oleksiivelychko/go-jwt-issuer

// +heroku goVersion go1.16
go 1.17

require (
	github.com/go-redis/redis/v8 v8.11.4
	github.com/golang-jwt/jwt v3.2.2+incompatible
	github.com/google/uuid v1.3.0
	github.com/gorilla/mux v1.8.0
	github.com/oleksiivelychko/go-helper v0.0.0-20220403045951-f3fd40772938
)

require (
	github.com/cespare/xxhash/v2 v2.1.2 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
)
