.PHONY: help test test-coverage lint build clean fmt vet

help: ## Display this help message
	@echo "Available targets:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  %-20s %s\n", $$1, $$2}'

test: ## Run tests
	go test -v -race ./...

test-coverage: ## Run tests with coverage
	go test -v -race -coverprofile=coverage.out -covermode=atomic ./...
	go tool cover -html=coverage.out -o coverage.html

lint: ## Run linters
	golangci-lint run ./...

fmt: ## Format code
	go fmt ./...

vet: ## Run go vet
	go vet ./...

build: ## Build the project
	go build ./...

clean: ## Clean build artifacts
	rm -f coverage.out coverage.html
	go clean ./...

deps: ## Download dependencies
	go mod download
	go mod tidy

.DEFAULT_GOAL := help
