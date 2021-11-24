# go-jwt-issuer

### Microservice generates pair access and refresh JSON web tokens signed by user identifier.

💡 Deployed on <a href="https://oleksiivelychkogojwtissuer.herokuapp.com">Heroku</a>

Run tests:
```
export SECRET_KEY=secretkey && go test ./*/
```

To generate a new token:
```
curl http://0.0.0.0:8080/access-token
curl http://127.0.0.1:30000/access-token
```

Available .env variables:
```
SECRET_KEY=secretkey
CLAIMS_AUD=
CLAIMS_ISS=
CLAIMS_EXP=60
PORT=8080
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=secret
REDIS_DB=0
```
