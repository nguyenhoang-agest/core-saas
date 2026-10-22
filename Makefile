# ==============================================================================
# Core SaaS — Unified Workspace Makefile
# ==============================================================================

SHELL := /usr/bin/env bash
.DEFAULT_GOAL := help

# Colors for terminal output
CYAN  := \033[36m
GREEN := \033[32m
YELLOW:= \033[33m
RESET := \033[0m

.PHONY: help env setup db-up db-down db-status db-logs db-psql dev-backend dev-frontend test test-backend test-frontend lint build clean new-module

help: ## Show this help message with available commands
	@echo -e ""
	@echo -e "$(CYAN)Core SaaS Monorepo — Available Make Commands:$(RESET)"
	@echo -e "$(CYAN)------------------------------------------------------------$(RESET)"
	@echo -e ""
	@echo -e "$(YELLOW) Quickstart & Setup:$(RESET)"
	@echo -e "  $(GREEN)make setup$(RESET)          Run one-click setup (env + docker db + deps check)"
	@echo -e "  $(GREEN)make env$(RESET)            Generate .env and .env.local from .env.example files"
	@echo -e ""
	@echo -e "$(YELLOW) Database (PostgreSQL):$(RESET)"
	@echo -e "  $(GREEN)make db-up$(RESET)              Start local PostgreSQL 17 container in background"
	@echo -e "  $(GREEN)make db-down$(RESET)            Stop and remove PostgreSQL container"
	@echo -e "  $(GREEN)make db-status$(RESET)          Check status and health of PostgreSQL container"
	@echo -e "  $(GREEN)make db-logs$(RESET)            View and follow PostgreSQL container logs"
	@echo -e "  $(GREEN)make db-psql$(RESET)            Open interactive psql CLI inside PostgreSQL container"
	@echo -e "  $(GREEN)make db-migrate-status$(RESET)  Display applied migration versions for all modules"
	@echo -e ""
	@echo -e "$(YELLOW) Development Servers:$(RESET)"
	@echo -e "  $(GREEN)make dev-backend$(RESET)        Start Go Fiber backend monolith (http://localhost:8080)"
	@echo -e "  $(GREEN)make dev-frontend$(RESET)       Start Next.js 16 frontend admin (http://localhost:3000)"
	@echo -e ""
	@echo -e "$(YELLOW) Testing & Quality:$(RESET)"
	@echo -e "  $(GREEN)make test$(RESET)               Run all tests across backend (Go) and frontend (Vitest)"
	@echo -e "  $(GREEN)make test-backend$(RESET)       Run backend unit test suite (with race detector)"
	@echo -e "  $(GREEN)make test-frontend$(RESET)      Run frontend unit and integration test suite"
	@echo -e "  $(GREEN)make lint$(RESET)               Run linter on both backend and frontend"
	@echo -e "  $(GREEN)make build$(RESET)              Build production binaries and Next.js Turbopack bundle"
	@echo -e ""
	@echo -e "$(YELLOW) Module & Migration Generators:$(RESET)"
	@echo -e "  $(GREEN)make new-module NAME=<name>$(RESET)               Scaffold new backend domain module"
	@echo -e "  $(GREEN)make new-migration MODULE=<mod> NAME=<name>$(RESET) Create sequential Goose migration"
	@echo -e ""

# ------------------------------------------------------------------------------
# Environment Setup
# ------------------------------------------------------------------------------

env: ## Create .env and .env.local from templates if they do not exist
	@echo -e "$(CYAN)Setting up environment files...$(RESET)"
	@if [ ! -f core-backend/.env ]; then \
		cp core-backend/.env.example core-backend/.env; \
		echo -e "$(GREEN)✓ Created core-backend/.env$(RESET)"; \
	else \
		echo -e "$(YELLOW)• core-backend/.env already exists (kept intact)$(RESET)"; \
	fi
	@if [ ! -f core-frontend/.env.local ]; then \
		cp core-frontend/.env.example core-frontend/.env.local; \
		echo -e "$(GREEN)✓ Created core-frontend/.env.local$(RESET)"; \
	else \
		echo -e "$(YELLOW)• core-frontend/.env.local already exists (kept intact)$(RESET)"; \
	fi

setup: ## One-click environment bootstrap
	@./scripts/setup.sh

# ------------------------------------------------------------------------------
# Database Management
# ------------------------------------------------------------------------------

db-up: ## Start PostgreSQL container
	@echo -e "$(CYAN)Starting PostgreSQL 17 container...$(RESET)"
	@docker compose -f core-backend/deploy/docker-compose.yml up -d
	@echo -e "$(GREEN)Waiting for database readiness...$(RESET)"
	@until docker exec saas-starter-postgres pg_isready -U postgres -d saas_dev >/dev/null 2>&1; do sleep 1; done
	@echo -e "$(GREEN)✓ PostgreSQL is ready to accept connections on port 5432.$(RESET)"

db-down: ## Stop PostgreSQL container
	@echo -e "$(CYAN)Stopping PostgreSQL container...$(RESET)"
	@docker compose -f core-backend/deploy/docker-compose.yml down
	@echo -e "$(GREEN)✓ PostgreSQL container stopped.$(RESET)"

db-status: ## Check PostgreSQL container status
	@docker ps --filter "name=saas-starter-postgres"

db-logs: ## View PostgreSQL logs
	@docker compose -f core-backend/deploy/docker-compose.yml logs -f postgres

db-psql: ## Open interactive psql session
	@docker exec -it saas-starter-postgres psql -U postgres -d saas_dev

db-migrate-status: ## Check applied migration status across all modules
	@cd core-backend && $(MAKE) db-migrate-status

# ------------------------------------------------------------------------------
# Local Development
# ------------------------------------------------------------------------------

dev-backend: ## Run Go Fiber backend monolith locally
	@echo -e "$(CYAN)Starting Core SaaS backend monolith on :8080...$(RESET)"
	@cd core-backend && go run ./cmd/monolith

dev-frontend: ## Run Next.js 16 frontend admin dev server
	@echo -e "$(CYAN)Starting Core SaaS frontend admin on :3000...$(RESET)"
	@cd core-frontend && npm run dev

# ------------------------------------------------------------------------------
# Testing & Verification
# ------------------------------------------------------------------------------

test-backend: ## Run backend tests with race detection
	@echo -e "$(CYAN)Running Go backend tests...$(RESET)"
	@cd core-backend && go test -race -v ./...

test-frontend: ## Run frontend Vitest test suite
	@echo -e "$(CYAN)Running Next.js frontend tests...$(RESET)"
	@cd core-frontend && npm test

test: test-backend test-frontend ## Run all test suites across the repository
	@echo -e "$(GREEN)✓ All backend and frontend test suites passed!$(RESET)"

lint: ## Run linters on both backend and frontend
	@echo -e "$(CYAN)Linting Go backend...$(RESET)"
	@cd core-backend && (which golangci-lint >/dev/null 2>&1 && golangci-lint run || echo "golangci-lint not installed; running go vet...") && go vet ./...
	@echo -e "$(CYAN)Linting Next.js frontend...$(RESET)"
	@cd core-frontend && npm run lint

build: ## Build both backend binary and frontend production bundle
	@echo -e "$(CYAN)Building Go backend binary...$(RESET)"
	@cd core-backend && go build -o bin/monolith ./cmd/monolith
	@echo -e "$(CYAN)Building Next.js frontend production bundle...$(RESET)"
	@cd core-frontend && npm run build
	@echo -e "$(GREEN)✓ Full workspace build succeeded!$(RESET)"

clean: ## Clean build binaries and temporary files
	@rm -rf core-backend/bin core-frontend/.next
	@echo -e "$(GREEN)✓ Cleaned build artifacts.$(RESET)"

# ------------------------------------------------------------------------------
# Generators
# ------------------------------------------------------------------------------

new-module: ## Scaffold a new backend module (usage: make new-module NAME=billing [DIRS_ONLY=true])
	@if [ -z "$(NAME)" ]; then echo "Error: NAME is required. Example: make new-module NAME=billing [DIRS_ONLY=true]"; exit 1; fi
	@cd core-backend && ./scripts/new-module.sh $(NAME) $(if $(filter true 1 yes,$(DIRS_ONLY)$(BARE)$(EMPTY)),--dirs-only,)

new-migration: ## Scaffold a new Goose migration (usage: make new-migration MODULE=tenant NAME=add_domain)
	@if [ -z "$(MODULE)" ] || [ -z "$(NAME)" ]; then echo "Error: MODULE and NAME are required. Example: make new-migration MODULE=tenant NAME=add_domain"; exit 1; fi
	@cd core-backend && ./scripts/new-migration.sh $(MODULE) $(NAME)
