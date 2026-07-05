.PHONY: build test lint install clean

build:
	go build -o bin/leads-brasil-pp-cli ./cmd/leads-brasil-pp-cli

test:
	go test ./...

lint:
	golangci-lint run

install:
	go install ./cmd/leads-brasil-pp-cli

clean:
	rm -rf bin/

build-mcp:
	go build -o bin/leads-brasil-pp-mcp ./cmd/leads-brasil-pp-mcp

install-mcp:
	go install ./cmd/leads-brasil-pp-mcp

build-all: build build-mcp
