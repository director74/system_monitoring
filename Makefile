BIN := "./bin/agent"

GIT_HASH := $(shell git log --format="%h" -n 1)
LDFLAGS := -X main.release="develop" -X main.buildDate=$(shell date -u +%Y-%m-%dT%H:%M:%S) -X main.gitHash=$(GIT_HASH)

build:
	go build -v -o $(BIN) -ldflags "$(LDFLAGS)" ./cmd/agent

run: build
	$(BIN) -config ./configs/measure.yml

generate:
	rm -rf pkg/grpc/ps
	mkdir -p pkg/grpc/ps

	protoc \
		--proto_path=api/ \
		--go_out=pkg/grpc/ps \
		--go-grpc_out=pkg/grpc/ps \
		api/*.proto