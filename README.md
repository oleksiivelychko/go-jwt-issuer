# go-jwt-issuer

### Microservice generates pair JSON web tokens - access-token and refresh-token signed by user identifier.

💡 Deployed on <a href="https://oleksiivelychkogojwtissuer.herokuapp.com/access-token/?userId=1">Heroku</a>

Run tests:
```
go clean -testcache && go test ./*/
```

To generate a new tokens pair for user with identifier 1:
```
GET http://0.0.0.0:8080/access-token/?userId=1
GET http://127.0.0.1:30000/access-token/?userId=1
GET https://oleksiivelychkogojwtissuer.herokuapp.com/access-token/?userId=1
```

To re-generate the tokens pair for user with identifier 1:
```
POST http://0.0.0.0:8080/refresh-token
POST http://127.0.0.1:30000/refresh-token
POST https://oleksiivelychkogojwtissuer.herokuapp.com/refresh-token

Accept: application/json
Authorization: <refresh-token>
Expires: <expiration-time>
```

To remove the tokens pair for user with identifier 1:
```
POST http://0.0.0.0:8080/clear-token
POST http://127.0.0.1:30000/clear-token
POST https://oleksiivelychkogojwtissuer.herokuapp.com/clear-token

Accept: application/json
Authorization: <access-token>
Expires: <expiration-time>
```

Available .env variables with default values:
```
SECRET_KEY=secretkey
AUDIENCE_AUD=oleksiivelychkogoaccount.herokuapp.com
ISSUER_ISS=oleksiivelychkogojwtissuer.herokuapp.com
EXPIRES_MINUTES=5
PORT=8080
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=secret
REDIS_DB=0
```

To get/check Redis data:
```
redis-cli --pass secret --no-auth-warning keys token-*
redis-cli --pass secret --no-auth-warning get token-1
redis-cli --pass secret --no-auth-warning del token-1
```
