.PHONY: help build run test tidy clean

help:
	@echo "GraphSentinel — targets: build, run, test, tidy, clean"

build:
	go build -o bin/graphsentinel ./cmd/server

run:
	go run ./cmd/server

test:
	go test ./...

tidy:
	go mod tidy

clean:
	rm -rf bin/
