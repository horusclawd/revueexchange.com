# RevUExchange Makefile

.PHONY: help local-start local-stop local-api local-ui test build

# Help
help:
	@echo "RevUExchange - Available Commands"
	@echo ""
	@echo "Local Development:"
	@echo "  make local-start         Start LocalStack and local services"
	@echo "  make local-stop         Stop local services"
	@echo ""
	@echo "Development Servers:"
	@echo "  make local-api          Start Go API locally"
	@echo "  make local-ui           Start React UI locally"
	@echo ""
	@echo "Testing & Build:"
	@echo "  make test              Run tests"
	@echo "  make build             Build for production"
	@echo ""

# Local Development
local-start:
	@chmod +x scripts/start-localstack.sh
	./scripts/start-localstack.sh

local-stop:
	@chmod +x scripts/stop-localstack.sh
	./scripts/stop-localstack.sh

# Development Servers
local-api:
	cd api && go run cmd/api/main.go

local-ui:
	cd frontend && npm run dev

# Testing & Build
test:
	cd api && go test ./...
	cd frontend && npm run test

build:
	cd api && go build -o bin/api ./cmd/api
	cd frontend && npm run build
