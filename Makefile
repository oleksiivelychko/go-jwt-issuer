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

kube-apply:
	kubectl apply -f .kubernetes/deployment.yml \
	&& kubectl apply -f .kubernetes/service.yml

kube-delete:
	kubectl delete deployment gojwtissuer-dpl --namespace=gojwtissuer \
	&& kubectl delete service gojwtissuer-srv --namespace=gojwtissuer

build:
	[[ -z "$(docker images -q local/gojwtissuer)" ]] || docker image rm local/gojwtissuer && \
	docker buildx build --platform linux/amd64 --tag local/gojwtissuer .

push-to-dockerhub: build
	$(warning instead of `local` prefix use dockerhub account name and change `imagePullPolicy`.)
	docker push local/gojwtissuer
