# go-jwt-issuer

### Microservice generates the pair of JSON web tokens - access-token and refresh-token are signed by user identifier.

📌 Create namespace before deploy on Kubernetes cluster:
```
kubectl create ns gons
```

📌 There are available environment variables with default values:
```
SECRET_KEY=secretkey
AUDIENCE_AUD=account.jwt.local
ISSUER_ISS=jwt.local
EXPIRES_MINUTES=60
HOST=127.0.0.1
PORT=8080
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=secret
REDIS_URL=redis://:secret@localhost:6379
```

![Debugging an application](social_preview.png)
