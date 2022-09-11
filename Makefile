go-build:
	go build -o bin/app -v .

go-update:
	go mod edit -go=1.18
	go get -u ./... && go mod tidy

go-test:
	$(warning main:19 env.SetDefaults - uncomment for local testing)
	go clean -testcache && go test ./*/

kube-apply-all:
	kubectl apply -f .kubernetes
kube-delete-all:
	kubectl delete -f .kubernetes

kubectl-redis-cli:
	kubectl -n gojwtissuer exec -it redis -- redis-cli

secret-create:
	@/bin/echo -n 'secret' > .kubernetes/secrets/password.txt
	kubectl create secret generic redis --from-file=password=.kubernetes/secrets/password.txt -n gojwtissuer
secret-get:
	kubectl get secrets/redis -o yaml -n gojwtissuer
	kubectl get secret redis -o jsonpath='{.data.password}' -n gojwtissuer | base64 --decode
secret-delete:
	kubectl delete secret redis -n gojwtissuer

docker-build:
	[[ -z "$(docker images -q local/gojwtissuer)" ]] || docker image rm local/gojwtissuer
	docker build --tag local/gojwtissuer .

run-redis:
	docker run --rm --name redis-server -p 6379:6379 redis --requirepass "secret"

push-to-dockerhub: docker-build
	$(warning instead of `local` prefix use dockerhub account name and change `imagePullPolicy`.)
	docker buildx build --platform linux/amd64 --tag local/gojwtissuer .
	docker push local/gojwtissuer
