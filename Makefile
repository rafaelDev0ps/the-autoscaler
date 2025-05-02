build-server:
	docker buildx build -f Dockerfile.server . -t server

run-server:
	docker run -d -v /var/run/docker.sock:/var/run/docker.sock server

build-client:
	docker buildx build -f Dockerfile.client . -t client --build-arg ORCHESTRATOR_URL=$(ifconfig | grep -w inet | awk '{print $2} ' | tail -n 1)

go-run-server:
	go run server/main.go

go-run-client:
	go run client/main.go

.PHONY: \
	build-server \
	run-server  \
	build-client \
	go-run-client \ 
	go-run-server