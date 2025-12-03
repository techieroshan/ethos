# Ethos Platform Makefile

.PHONY: help build test clean docker-build docker-run docker-stop dev setup db-migrate db-rollback lint format

# Default target
help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

# Development
dev: ## Start development environment with docker-compose
	docker-compose up --build

dev-detached: ## Start development environment in background
	docker-compose up --build -d

dev-logs: ## Show development environment logs
	docker-compose logs -f

dev-stop: ## Stop development environment
	docker-compose down

# Building
build: ## Build the Go backend
	go build -o bin/ethos ./cmd/api

build-frontend: ## Build the React frontend
	cd ethos-ui && npm run build

build-all: build build-frontend ## Build both backend and frontend

# Testing
test: ## Run Go backend tests
	go test -v -race -cover ./...

test-frontend: ## Run React frontend tests
	cd ethos-ui && npm test -- --watchAll=false

test-e2e: ## Run E2E tests
	cd ethos-ui && npm run test:e2e

test-all: test test-frontend test-e2e ## Run all tests

# Docker
docker-build: ## Build Docker images
	docker-compose build

docker-run: ## Run Docker containers
	docker-compose up -d

docker-stop: ## Stop Docker containers
	docker-compose down

docker-logs: ## Show Docker container logs
	docker-compose logs -f

docker-clean: ## Remove Docker containers and volumes
	docker-compose down -v --remove-orphans

# Database
db-migrate: ## Run database migrations
	@echo "Running database migrations..."
	# Add your migration command here
	@echo "Migrations completed"

db-rollback: ## Rollback database migrations
	@echo "Rolling back database migrations..."
	# Add your rollback command here
	@echo "Rollback completed"

db-seed: ## Seed database with test data
	@echo "Seeding database..."
	# Add your seeding command here
	@echo "Database seeded"

# Code Quality
lint: ## Run Go linting
	go vet ./...
	golangci-lint run

lint-frontend: ## Run frontend linting
	cd ethos-ui && npm run lint

lint-all: lint lint-frontend ## Run all linting

format: ## Format Go code
	go fmt ./...
	goimports -w .

format-frontend: ## Format frontend code
	cd ethos-ui && npm run format

format-all: format format-frontend ## Format all code

# Setup
setup: ## Set up development environment
	@echo "Setting up development environment..."
	# Install Go dependencies
	go mod download
	# Install frontend dependencies
	cd ethos-ui && npm install
	# Create necessary directories
	mkdir -p logs uploads
	@echo "Setup complete"

setup-hooks: ## Set up Git hooks
	@echo "Setting up Git hooks..."
	cp scripts/pre-commit .git/hooks/pre-commit
	chmod +x .git/hooks/pre-commit
	@echo "Git hooks set up"

# Deployment
deploy-staging: ## Deploy to staging environment
	@echo "Deploying to staging..."
	# Add your staging deployment commands here
	@echo "Staging deployment complete"

deploy-production: ## Deploy to production environment
	@echo "Deploying to production..."
	# Add your production deployment commands here
	@echo "Production deployment complete"

# Monitoring
health-check: ## Check application health
	curl -f http://localhost:8000/api/v1/health || echo "Backend not healthy"
	curl -f http://localhost:3000/health || echo "Frontend not healthy"

logs: ## Show application logs
	tail -f logs/*.log

# Cleanup
clean: ## Clean build artifacts
	rm -rf bin/
	rm -rf ethos-ui/dist/
	rm -rf ethos-ui/node_modules/.cache
	go clean

clean-all: clean ## Clean everything including containers
	docker system prune -f
	docker volume prune -f

# CI/CD
ci: lint-all test-all docker-build ## Run CI pipeline locally

# Utility
version: ## Show version information
	@echo "Ethos Platform v1.0.0"
	@go version
	@cd ethos-ui && npm --version

# Database operations (requires migrate tool)
migrate-up: ## Run database migrations up
	migrate -path internal/database/migrations -database $(DATABASE_URL) up

migrate-down: ## Run database migrations down
	migrate -path internal/database/migrations -database $(DATABASE_URL) down 1

migrate-create: ## Create a new migration file
	@echo "Usage: make migrate-create name=migration_name"
	@if [ -z "$(name)" ]; then echo "Please provide name=migration_name"; exit 1; fi
	migrate create -ext sql -dir internal/database/migrations -seq $(name)