.PHONY: help dev build run test clean docker-up docker-down migrate-up migrate-down migrate-create migrate-install

help: ## Display this help screen
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

dev: ## Run the application with hot reload using Air
	@if command -v air > /dev/null; then \
		air; \
	else \
		echo "Installing air..."; \
		go install github.com/cosmtrek/air@latest; \
		air; \
	fi

build: ## Build the application
	@echo "Building..."
	@go build -o bin/main cmd/api/main.go

run: ## Run the application
	@echo "Running..."
	@go run cmd/api/main.go

test: ## Run tests
	@echo "Running tests..."
	@go test -v ./...

clean: ## Clean build files
	@echo "Cleaning..."
	@rm -rf bin tmp

docker-up: ## Start Docker containers
	@echo "Starting Docker containers..."
	@docker compose up -d

docker-down: ## Stop Docker containers
	@echo "Stopping Docker containers..."
	@docker compose down

docker-logs: ## View Docker logs
	@docker compose logs -f

migrate-up: ## Run database migrations up
	@echo "Running migrations..."
	@if command -v migrate > /dev/null; then \
		migrate -path migrations -database "postgresql://postgres:postgres@localhost:5432/go_clean_boiler?sslmode=disable" up; \
	else \
		echo "migrate CLI not found. Using Docker..."; \
		docker run --rm -v $(PWD)/migrations:/migrations --network host migrate/migrate \
			-path=/migrations -database "postgresql://postgres:postgres@localhost:5432/go_clean_boiler?sslmode=disable" up; \
	fi

migrate-down: ## Rollback database migrations
	@echo "Rolling back migrations..."
	@if command -v migrate > /dev/null; then \
		migrate -path migrations -database "postgresql://postgres:postgres@localhost:5432/go_clean_boiler?sslmode=disable" down; \
	else \
		echo "migrate CLI not found. Using Docker..."; \
		docker run --rm -v $(PWD)/migrations:/migrations --network host migrate/migrate \
			-path=/migrations -database "postgresql://postgres:postgres@localhost:5432/go_clean_boiler?sslmode=disable" down; \
	fi

migrate-create: ## Create a new migration file (usage: make migrate-create name=create_users_table)
	@if [ -z "$(name)" ]; then \
		echo "Error: name is required. Usage: make migrate-create name=your_migration_name"; \
		exit 1; \
	fi
	@if command -v migrate > /dev/null; then \
		migrate create -ext sql -dir migrations -seq $(name); \
	else \
		echo "migrate CLI not found. Using Docker..."; \
		docker run --rm -v $(PWD)/migrations:/migrations migrate/migrate \
			create -ext sql -dir /migrations -seq $(name); \
	fi

migrate-install: ## Install golang-migrate CLI
	@echo "Installing golang-migrate..."
	@go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
	@echo "migrate installed successfully!"


tidy: ## Tidy go modules
	@echo "Tidying go modules..."
	@go mod tidy

deps: ## Download dependencies
	@echo "Downloading dependencies..."
	@go mod download
