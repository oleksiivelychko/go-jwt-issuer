go-build:
	go build -o bin/app -v .

go-update: go-upgrade
	go get -u ./... && go mod tidy

go-upgrade:
	go mod edit -go=1.18

go-test:
	$(warning main:18 env.InitEnv - uncomment for local testing)
	go clean -testcache && go test ./*/

kube-apply:
	kubectl apply -f .kubernetes/deployment.yml && kubectl apply -f .kubernetes/service.yml

kube-stop:
	kubectl delete deployment gojwtissuer-deployment && kubectl delete service gojwtissuer-service

build-and-push-docker:
	$(warning `instead of <local> prefix use dockerhub account name.`)
	[[ -z "$(docker images -q local/gojwtissuer)" ]] || docker image rm local/gojwtissuer && docker buildx build --platform linux/amd64 --tag local/gojwtissuer . && docker push local/gojwtissuer

