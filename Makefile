.PHONY: build run test clean docker-build docker-run help

# Variables
BINARY_NAME=frostbyte-api
DOCKER_IMAGE=frostbyte-api
PORT?=8080

help: ## Display this help message
	@echo "FrostByte-API - Makefile commands:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}'

build: ## Build the binary
	@echo "Building $(BINARY_NAME)..."
	@go build -o $(BINARY_NAME) .
	@echo "Build complete!"

run: ## Run the application
	@echo "Running $(BINARY_NAME)..."
	@go run .

test: ## Run tests
	@echo "Running tests..."
	@go test -v -race -coverprofile=coverage.out ./...
	@echo "Tests complete!"

coverage: test ## Run tests with coverage report
	@echo "Generating coverage report..."
	@go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

clean: ## Clean build artifacts
	@echo "Cleaning..."
	@rm -f $(BINARY_NAME)
	@rm -f coverage.out coverage.html
	@echo "Clean complete!"

docker-build: ## Build Docker image
	@echo "Building Docker image..."
	@docker build -t $(DOCKER_IMAGE):latest .
	@echo "Docker build complete!"

docker-run: ## Run Docker container
	@echo "Running Docker container..."
	@docker run -p $(PORT):8080 --rm $(DOCKER_IMAGE):latest

fmt: ## Format code
	@echo "Formatting code..."
	@go fmt ./...
	@echo "Format complete!"

vet: ## Run go vet
	@echo "Running go vet..."
	@go vet ./...
	@echo "Vet complete!"

lint: fmt vet ## Run linters

install-deps: ## Install dependencies
	@echo "Installing dependencies..."
	@go mod download
	@go mod verify
	@echo "Dependencies installed!"

all: clean lint test build ## Run all checks and build
