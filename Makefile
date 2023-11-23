export IMAGE ?= data_microservices:latest


.PHONY: build
build:
	docker build -t $(IMAGE) -f docker/Dockerfile .

.PHONY: run-app
run-app:
	go build -o build/main
	./build/main

.PHONY: run
run:
	make build
	docker-compose -f docker/docker-compose.yaml up