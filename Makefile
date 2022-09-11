go-build:
	go build -o bin/app -v .

go-update: go-upgrade
	go get -u ./... && go mod tidy

go-upgrade:
	go mod edit -go=1.18

go-test:
	$(warning main:19 env.SetDefaults - uncomment for local testing)
	go clean -testcache && go test ./*/

kube-ns:
	kubectl create ns gojwtissuer

kube-apply-all:
	kubectl apply -f .kubernetes

kube-delete-all:
	kubectl delete -f .kubernetes

kube-get-all:
	kubectl get pod/redis configmap/redis -n gojwtissuer

kube-describe-redis:
	kubectl describe configmap/redis -n gojwtissuer

kubectl-exec-redis:
	kubectl -n gojwtissuer exec -it redis -- redis-cli

create-secret:
	@/bin/echo -n 'secret' > .kubernetes/secrets/password.txt
	kubectl create secret generic redis --from-file=password=.kubernetes/secrets/password.txt -n gojwtissuer

get-secrets:
	kubectl get secrets/redis -o yaml -n gojwtissuer
	kubectl get secret redis -o jsonpath='{.data.password}' -n gojwtissuer | base64 --decode

delete-secret:
	kubectl delete secret redis -n gojwtissuer

build:
	[[ -z "$(docker images -q local/gojwtissuer)" ]] || docker image rm local/gojwtissuer && \
	docker buildx build --platform linux/amd64 --tag local/gojwtissuer .

run-redis:
	docker run --rm --name redis-server -p 6379:6379 redis --requirepass "secret"

push-to-dockerhub: build
	$(warning instead of `local` prefix use dockerhub account name and change `imagePullPolicy`.)
	docker push local/gojwtissuer
