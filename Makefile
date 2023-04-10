docker-build:
	[[ -z "$(docker images -q local/gojwtissuerdlv)" ]] || docker image rm local/gojwtissuerdlv
	docker build --no-cache --tag local/gojwtissuerdlv . -f dlv.dockerfile

docker-redis:
	docker run --rm --name redis-server -p 6379:6379 redis --requirepass "secret"

kube-ns:
	kubectl create ns gons

kube-apply:
	kubectl apply -f .kubernetes

kube-delete:
	kubectl delete -f .kubernetes

kube-redis-all:
	kubectl -n gons exec -it redis -- redis-cli --pass secret --no-auth-warning keys token-*

# make kube-redis-get id=userID
kube-redis-get:
	kubectl -n gons exec -it redis -- redis-cli --pass secret --no-auth-warning get token-$(id)

# make kube-redis-del id=userID
kube-redis-del:
	kubectl -n gons exec -it redis -- redis-cli --pass secret --no-auth-warning del token-$(id)

# make redis-secret password=PASSWORD
redis-secret: redis-secret-del
	@/bin/echo -n '$(password)' > .kubernetes/secrets/password.txt
	kubectl create secret generic redis --from-file=password=.kubernetes/secrets/password.txt -n gons

redis-secret-del:
	kubectl delete secret redis -n gons

redis-secret-get:
	kubectl get secrets/redis -o yaml -n gons
	kubectl get secret redis -o jsonpath='{.data.password}' -n gons | base64 --decode

run:
	HOST=localhost \
	PORT=8080 \
	REDIS_HOST=localhost \
	REDIS_PORT=6379 \
	REDIS_PASSWORD=secret \
	SECRET_KEY=secretkey \
	ISSUER_ISS=jwt.local \
	AUDIENCE_AUD=jwt.account.local \
	EXPIRATION_TIME_EXP=1 \
	go run main.go

test:
	go clean -testcache && go test ./*/
