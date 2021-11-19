# go-jwt-issuer

### Microservice generates JSON web tokens signed by user identifier.

💡 Deployed on <a href="https://oleksiivelychkogojwtissuer.herokuapp.com">Heroku</a>

Run tests:
```
export SECRET_KEY=secretkey && go test ./*/
```

To generate a new token:
```
curl http://0.0.0.0:8080/issue
curl http://127.0.0.1:30000/issue
```
