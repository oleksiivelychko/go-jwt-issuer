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
	kubectl delete deployment gojwtissuer-instance && kubectl delete service gojwtissuer-service

build-and-push-docker:
	[[ -z "$(docker images -q mydockerhubaccount/gojwtissuer)" ]] || docker image rm mydockerhubaccount/gojwtissuer && docker buildx build --platform linux/amd64 --tag mydockerhubaccount/gojwtissuer . && docker push mydockerhubaccount/gojwtissuer

