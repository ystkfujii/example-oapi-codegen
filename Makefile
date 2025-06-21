.PHONY: help build run test clean generate lint format

# Default target
help: ## Show this help message
	@echo "Available targets:"
	@awk 'BEGIN {FS = ":.*##"; printf "\n"} /^[a-zA-Z_-]+:.*##/ { printf "  %-15s %s\n", $$1, $$2 }' $(MAKEFILE_LIST)

test: ## Run tests
	go test -v ./...

generate: ## Generate OpenAPI code from spec
	go tool oapi-codegen -config oapi-codegen.config.yaml ./openapi/openapi.yaml 

lint: ## Run linter
	golangci-lint run

format: ## Format code
	go fmt ./...

clean: ## Clean build artifacts
	rm -f bin/server coverage.out coverage.html

docker-compose-up:
	docker compose up -d
