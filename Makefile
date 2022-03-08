go-build:
	go build -o bin/app -v .

go-update:
	go get -u && go mod tidy

go-test:
	go clean -testcache && go test ./*/

kube-apply:
	kubectl apply -f .kubernetes/deployment.yml && kubectl apply -f .kubernetes/service.yml

kube-stop:
	kubectl delete deployment gojwtissuer-instance && kubectl delete service gojwtissuer-service

build-and-push-docker:
	[[ -z "$(docker images -q alexvelychko/gojwtissuer)" ]] || docker image rm alexvelychko/gojwtissuer && docker buildx build --platform linux/amd64 --tag alexvelychko/gojwtissuer . && docker push alexvelychko/gojwtissuer
