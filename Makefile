.PHONY: help build up down logs clean test

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

build: ## Build all Docker images
	docker-compose build

up: ## Start all services
	docker-compose up -d

down: ## Stop all services
	docker-compose down

logs: ## View logs from all services
	docker-compose logs -f

clean: ## Remove all containers and volumes
	docker-compose down -v

test: ## Run backend tests
	cd backend && go test ./services/... -v

setup-db: ## Setup database locally (requires PostgreSQL)
	@echo "Setting up database..."
	psql -U postgres -c "CREATE DATABASE brute_force_login;" || true
	psql -U postgres -d brute_force_login -f backend/database/schema.sql
	psql -U postgres -d brute_force_login -f backend/database/init.sql
	@echo "Database setup complete!"

run-backend: ## Run backend locally
	cd backend && go mod tidy && go run main.go

run-frontend: ## Run frontend locally
	cd frontend && npm install && npm run dev

