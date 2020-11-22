GOBIN=$(shell pwd)/bin
GONAME=$(shell basename "$(PWD)")
GOPATH=$(shell go env GOPATH)/bin

.PHONY: format build build_static test all
.DEFAULT_GOAL := all

format:
	go fmt

grpc:
	@echo "Generate go-grpc files"
	GO111MODULE=on go get google.golang.org/protobuf/cmd/protoc-gen-go
	GO111MODULE=on go get google.golang.org/grpc/cmd/protoc-gen-go-grpc
	PATH=${PATH}:$(GOPATH) protoc --go_out=./generated --go_opt=paths=source_relative --go-grpc_out=./generated --go-grpc_opt=paths=source_relative  interface/service.proto

build: grpc format
	@echo "Building Linux ${GOFILES} to ./bin"
	GOBIN=$(GOBIN) CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/$(GONAME) -a -tags netgo -ldflags '-w -extldflags "-static"' -o bin/$(GONAME)

test: build
	go test ./service/ ./service/handler/

coverage:
	go test -cover ./service/

all: test

container:
	docker build -t sgomenyuk/atlant:latest .
