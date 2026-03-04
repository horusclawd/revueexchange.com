# RevUExchange Makefile

.PHONY: help local-start local-stop local-api local-ui test test-backend test-frontend build build-docker lint

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
	@echo "Testing & Linting:"
	@echo "  make test              Run all tests"
	@echo "  make test-backend      Run backend tests"
	@echo "  make test-frontend     Run frontend tests"
	@echo "  make lint              Run linters"
	@echo ""
	@echo "Building:"
	@echo "  make build             Build for production"
	@echo "  make build-docker      Build Docker images"
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
	cd api && go run ./cmd/api/main.go

local-ui:
	cd frontend && npm run dev

# Testing & Linting
test: test-backend test-frontend

test-backend:
	cd api && go test -v -race ./...

test-frontend:
	cd frontend && npm test -- --run

lint:
	cd api && go install golang.org/x/lint/golint@latest && golint -set_exit_status ./...
	cd frontend && npm run lint

# Building
build:
	cd api && go build -o bin/server ./cmd/api
	cd frontend && npm run build

build-docker:
	docker build -t revueexchange-backend:latest ./api
	docker build -t revueexchange-frontend:latest ./frontend
