docker-build:
	[[ -z "$(docker images -q local/dlv-gojwtissuer)" ]] || docker image rm local/dlv-gojwtissuer
	docker build --tag local/dlv-gojwtissuer . -f dlv.dockerfile

docker-run-redis:
	docker run --rm --name redis-server -p 6379:6379 redis --requirepass "secret"

go-update:
	go mod edit -go=1.18
	go get -u ./... && go mod tidy

go-test:
	$(warning main:18 env.SetDefaults - uncomment for local testing)
	go clean -testcache && go test ./*/

kube-apply-all:
	kubectl apply -f .kubernetes

kube-delete-all:
	kubectl delete -f .kubernetes

kube-redis-cli:
	kubectl -n gojwtissuer exec -it redis -- redis-cli

secret-create:
	kubectl delete secret redis -n gojwtissuer
	@/bin/echo -n 'secret' > .kubernetes/secrets/password.txt
	kubectl create secret generic redis --from-file=password=.kubernetes/secrets/password.txt -n gojwtissuer

secret-verify:
	kubectl get secrets/redis -o yaml -n gojwtissuer
	kubectl get secret redis -o jsonpath='{.data.password}' -n gojwtissuer | base64 --decode
