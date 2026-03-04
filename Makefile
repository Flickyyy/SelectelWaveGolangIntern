.DEFAULT_GOAL := build

build:
	go build -o bin/loglint ./cmd/loglint

test:
	go test -v -race ./...

lint:
	golangci-lint run ./...

vet:
	go vet ./...

clean:
	rm -rf bin/

.PHONY: build test lint vet clean
