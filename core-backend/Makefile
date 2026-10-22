# ==============================================================================
# Core SaaS Backend — Makefile
# ==============================================================================

SHELL := /usr/bin/env bash
.DEFAULT_GOAL := help

# Colors
CYAN  := \033[36m
GREEN := \033[32m
YELLOW:= \033[33m
RESET := \033[0m

PROJECT ?=
MODULE ?=
KEEP ?=
NAME ?=

.PHONY: help env db-up db-down db-status db-logs db-psql up down run run-example build test test-v test-coverage lint tidy clean init new-module

help: ## Show this help message
	@echo -e ""
	@echo -e "$(CYAN)Core SaaS Backend — Available Commands:$(RESET)"
	@echo -e "$(CYAN)------------------------------------------------------------$(RESET)"
	@echo -e ""
	@echo -e "$(YELLOW) Environment & Setup:$(RESET)"
	@echo -e "  $(GREEN)make env$(RESET)                                  Create .env from .env.example if missing"
	@echo -e "  $(GREEN)make init PROJECT=<name> MODULE=<path>$(RESET)     Initialize and rename Go module"
	@echo -e "  $(GREEN)make new-module NAME=<name>$(RESET)                Scaffold new domain module from template"
	@echo -e ""
	@echo -e "$(YELLOW) Database (PostgreSQL):$(RESET)"
	@echo -e "  $(GREEN)make db-up$(RESET) (or $(GREEN)make up$(RESET))           Start PostgreSQL 17 container in background"
	@echo -e "  $(GREEN)make db-down$(RESET) (or $(GREEN)make down$(RESET))       Stop PostgreSQL container"
	@echo -e "  $(GREEN)make db-status$(RESET)                            Check container health and port 5432"
	@echo -e "  $(GREEN)make db-logs$(RESET)                              View PostgreSQL container logs"
	@echo -e "  $(GREEN)make db-psql$(RESET)                              Open interactive psql session in container"
	@echo -e "  $(GREEN)make db-migrate-status$(RESET)                     Show applied migration versions for all modules"
	@echo -e ""
	@echo -e "$(YELLOW) Run Services:$(RESET)"
	@echo -e "  $(GREEN)make run$(RESET)                                  Run monolith service locally (:8080)"
	@echo -e "  $(GREEN)make run-example$(RESET)                          Run standalone example microservice (:8081)"
	@echo -e ""
	@echo -e "$(YELLOW) Testing & Code Quality:$(RESET)"
	@echo -e "  $(GREEN)make test$(RESET)                                 Run test suite with race detector"
	@echo -e "  $(GREEN)make test-v$(RESET)                               Run test suite with verbose output"
	@echo -e "  $(GREEN)make test-coverage$(RESET)                        Run tests and generate coverage.out"
	@echo -e "  $(GREEN)make lint$(RESET)                                 Run golangci-lint / go vet"
	@echo -e "  $(GREEN)make tidy$(RESET)                                 Tidy and verify go.mod dependencies"
	@echo -e "  $(GREEN)make build$(RESET)                                Compile binaries into bin/"
	@echo -e "  $(GREEN)make clean$(RESET)                                Remove compiled binaries and coverage files"
	@echo -e ""

# ------------------------------------------------------------------------------
# Environment
# ------------------------------------------------------------------------------

env: ## Create .env from template
	@if [ ! -f .env ]; then \
		cp .env.example .env; \
		echo -e "$(GREEN)✓ Created .env from .env.example$(RESET)"; \
	else \
		echo -e "$(YELLOW)• .env already exists$(RESET)"; \
	fi

# ------------------------------------------------------------------------------
# Database
# ------------------------------------------------------------------------------

db-up up: ## Start local PostgreSQL container
	@echo -e "$(CYAN)Starting local PostgreSQL container...$(RESET)"
	@docker compose -f deploy/docker-compose.yml up -d
	@echo -e "$(GREEN)Waiting for PostgreSQL to be ready...$(RESET)"
	@until docker exec saas-starter-postgres pg_isready -U postgres -d saas_dev >/dev/null 2>&1; do sleep 1; done
	@echo -e "$(GREEN)✓ PostgreSQL is ready on localhost:5432$(RESET)"

db-down down: ## Stop local PostgreSQL container
	@echo -e "$(CYAN)Stopping PostgreSQL container...$(RESET)"
	@docker compose -f deploy/docker-compose.yml down
	@echo -e "$(GREEN)✓ PostgreSQL stopped$(RESET)"

db-status: ## Check PostgreSQL health
	@docker ps --filter "name=saas-starter-postgres"

db-logs: ## View PostgreSQL logs
	@docker compose -f deploy/docker-compose.yml logs -f postgres

db-psql: ## Open psql CLI inside container
	@docker exec -it saas-starter-postgres psql -U postgres -d saas_dev

db-migrate-status: ## Show applied migration versions per module
	@docker exec saas-starter-postgres psql -U postgres -d saas_dev -c "\
	SELECT \
	    replace(t.table_name, 'goose_version_', '') AS module, \
	    v.version_id, \
	    v.is_applied, \
	    v.tstamp \
	FROM information_schema.tables t \
	CROSS JOIN LATERAL ( \
	    SELECT version_id, is_applied, tstamp \
	    FROM ( \
	        SELECT version_id, is_applied, tstamp FROM goose_version_audit WHERE t.table_name = 'goose_version_audit' \
	        UNION ALL SELECT version_id, is_applied, tstamp FROM goose_version_example WHERE t.table_name = 'goose_version_example' \
	        UNION ALL SELECT version_id, is_applied, tstamp FROM goose_version_iam WHERE t.table_name = 'goose_version_iam' \
	        UNION ALL SELECT version_id, is_applied, tstamp FROM goose_version_notification WHERE t.table_name = 'goose_version_notification' \
	        UNION ALL SELECT version_id, is_applied, tstamp FROM goose_version_tenant WHERE t.table_name = 'goose_version_tenant' \
	    ) sub ORDER BY version_id DESC LIMIT 1 \
	) v \
	WHERE t.table_schema = 'public' AND t.table_name LIKE 'goose_version_%' \
	ORDER BY module;"

# ------------------------------------------------------------------------------
# Run
# ------------------------------------------------------------------------------

run: ## Run monolith service locally
	go run ./cmd/monolith

run-example: ## Run standalone example microservice
	go run ./cmd/example-service

# ------------------------------------------------------------------------------
# Testing & Quality
# ------------------------------------------------------------------------------

test: ## Run test suite with race detection
	go test -race ./...

test-v: ## Run test suite with verbose output
	go test -v ./...

test-coverage: ## Run tests with code coverage profile
	go test -coverprofile=coverage.out ./...
	go tool cover -func=coverage.out

lint: ## Run linter
	@which golangci-lint >/dev/null 2>&1 && golangci-lint run || (echo "golangci-lint not in PATH, running go vet..."; go vet ./...)

tidy: ## Tidy go.mod dependencies
	go mod tidy

build: ## Compile all packages and binaries
	@mkdir -p bin
	go build -o bin/monolith ./cmd/monolith
	go build -o bin/example-service ./cmd/example-service
	@echo -e "$(GREEN)✓ Binaries compiled into bin/$(RESET)"

clean: ## Clean build binaries and coverage
	rm -rf bin coverage.out
	@echo -e "$(GREEN)✓ Cleaned build artifacts$(RESET)"

# ------------------------------------------------------------------------------
# Project & Module Generators
# ------------------------------------------------------------------------------

init: ## Initialize and rename project (usage: make init PROJECT=my-saas MODULE=github.com/org/my-saas)
	@./scripts/init.sh PROJECT=$(PROJECT) MODULE=$(MODULE) KEEP=$(KEEP)

new-module: ## Scaffold new module from template (usage: make new-module NAME=billing [DIRS_ONLY=true])
	@if [ -z "$(NAME)" ]; then echo "Error: NAME is required. Example: make new-module NAME=billing [DIRS_ONLY=true]"; exit 1; fi
	@./scripts/new-module.sh $(NAME) $(if $(filter true 1 yes,$(DIRS_ONLY)$(BARE)$(EMPTY)),--dirs-only,)

new-migration: ## Scaffold new migration file (usage: make new-migration MODULE=tenant NAME=add_domain)
	@if [ -z "$(MODULE)" ] || [ -z "$(NAME)" ]; then echo "Error: MODULE and NAME are required. Example: make new-migration MODULE=tenant NAME=add_domain"; exit 1; fi
	@./scripts/new-migration.sh $(MODULE) $(NAME)
