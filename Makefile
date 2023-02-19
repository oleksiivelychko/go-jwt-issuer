docker-build:
	[[ -z "$(docker images -q local/gojwtissuerdlv)" ]] || docker image rm local/gojwtissuerdlv
	docker build --no-cache --tag local/gojwtissuerdlv . -f dlv.dockerfile

docker-run-redis:
	docker run --rm --name redis-server -p 6379:6379 redis --requirepass "secret"

go-edit:
	# go clean --modcache
	go mod edit -go=1.20
	go mod edit -module github.com/oleksiivelychko/go-jwt-issuer
	go get -u ./... && go mod tidy

go-test:
	go clean -testcache && go test ./*/

kube-apply-all:
	kubectl apply -f .kubernetes

kube-delete-all:
	kubectl delete -f .kubernetes

kube-redis-cli-all:
	kubectl -n gons exec -it redis -- redis-cli --pass secret --no-auth-warning keys token-*

# make kube-redis-cli-get [userId]
kube-redis-cli-get:
	kubectl -n gons exec -it redis -- redis-cli --pass secret --no-auth-warning get token-$(id)

# make kube-redis-cli-del [userId]
kube-redis-cli-del:
	kubectl -n gons exec -it redis -- redis-cli --pass secret --no-auth-warning del token-$(id)

# make secret-redis [password]
secret-redis:
	@/bin/echo -n '$(password)' > .kubernetes/secrets/password.txt
	kubectl create secret generic redis --from-file=password=.kubernetes/secrets/password.txt -n gons

secret-redis-del:
	kubectl delete secret redis -n gons

secret-redis-get:
	kubectl get secrets/redis -o yaml -n gons
	kubectl get secret redis -o jsonpath='{.data.password}' -n gons | base64 --decode
