.PHONY: build run test clean install

BINARY_SERVER=bin/h-ai-server
BINARY_MCP=bin/h-ai-mcp

build:
	@echo "Building H-AI..."
	@mkdir -p bin
	@go build -o $(BINARY_SERVER) ./main.go
	@go build -o $(BINARY_MCP) ./cmd/mcp
	@echo "Build complete!"

run: build
	@echo "Starting H-AI server..."
	@./$(BINARY_SERVER)

test:
	@echo "Running tests..."
	@go test -v ./...

clean:
	@echo "Cleaning..."
	@rm -rf bin/
	@go clean

install:
	@echo "Installing dependencies..."
	@go mod download
	@go mod tidy

fmt:
	@echo "Formatting code..."
	@go fmt ./...

lint:
	@echo "Linting code..."
	@golangci-lint run

.DEFAULT_GOAL := build
