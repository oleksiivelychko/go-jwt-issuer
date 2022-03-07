go-build:
	go build -o bin/go-jwt-issuer -v .

heroku-bash:
	heroku run bash -a oleksiivelychkogojwtissuer

heroku-logs:
	heroku logs -n 200 -a oleksiivelychkogojwtissuer --tail

kube-apply:
	kubectl apply -f .kubernetes/minikube/deployment.yml && kubectl apply -f .kubernetes/minikube/service.yml

kube-apply-gke:
	kubectl apply -f .kubernetes/gke/deployment.yml && kubectl apply -f .kubernetes/gke/service.yml

kube-stop:
	kubectl delete deployment gojwtissuer-instance && kubectl delete service gojwtissuer-service

build-and-push-docker-hub:
	[[ -z "$(docker images -q alexvelychko/gojwtissuer)" ]] || docker image rm alexvelychko/gojwtissuer && docker buildx build --platform linux/amd64 --tag alexvelychko/gojwtissuer . && docker push alexvelychko/gojwtissuer

build-and-push-gcp-artifact-registry:
	docker buildx build --platform linux/amd64 --pull --no-cache --tag europe-central2-docker.pkg.dev/PROJECT-NAME/repository/gojwtissuer . && docker push europe-central2-docker.pkg.dev/PROJECT-NAME/repository/gojwtissuer
