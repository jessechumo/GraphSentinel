.PHONY: help build run test tidy clean fmt vet check

help:
	@echo "GraphSentinel — targets: build, run, test, check, fmt, vet, tidy, clean"

build:
	go build -o bin/graphsentinel ./cmd/server

run:
	go run ./cmd/server

test:
	go test ./...

fmt:
	go fmt ./...

vet:
	go vet ./...

check: fmt vet test

tidy:
	go mod tidy

clean:
	rm -rf bin/
