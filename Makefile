.PHONY: bin/server
bin/server:
	go build -o bin/server ./cmd/server/

.PHONY: bin/server.linux
bin/server.linux:
	GOOS=linux go build -o bin/server.linux ./cmd/server/

.PHONY: docker
docker: bin/server.linux
	docker build -t jhedev/prom-timestream -f Dockerfile .
