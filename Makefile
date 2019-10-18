
GOPATH:=$(shell go env GOPATH)


.PHONY: build
build:
	swag init
	go build -o basic-web main.go plugin.go
.PHONY: test
test:
	go test -v ./... -cover

.PHONY: docker
docker:
	docker build . -t basic-web:latest

.PHONY:
run:
	swag init
	go run main.go
