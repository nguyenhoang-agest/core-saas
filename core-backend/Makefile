.PHONY: help init tidy build test lint run run-example up down new-module

PROJECT ?=
MODULE ?=
KEEP ?=
NAME ?=

help:
	@echo "SaaS Starter Kit — Available Commands:"
	@echo "  make init PROJECT=<name> MODULE=<path> [KEEP=iam,tenant]  Initialize & rename kit"
	@echo "  make new-module NAME=<name>                               Scaffold new module from template"
	@echo "  make build                                                Compile all packages & binaries"
	@echo "  make test                                                 Run test suite"
	@echo "  make lint                                                 Run golangci-lint"
	@echo "  make tidy                                                 Tidy go.mod dependencies"
	@echo "  make run                                                  Run monolith service locally"
	@echo "  make run-example                                          Run standalone example service"
	@echo "  make up                                                   Start local PostgreSQL"
	@echo "  make down                                                 Stop local PostgreSQL"

init:
	@./scripts/init.sh PROJECT=$(PROJECT) MODULE=$(MODULE) KEEP=$(KEEP)

new-module:
	@if [ -z "$(NAME)" ]; then echo "Error: NAME is required. Example: make new-module NAME=billing"; exit 1; fi
	@./scripts/new-module.sh $(NAME)

tidy:
	go mod tidy

build:
	go build ./...

test:
	go test -race ./...

lint:
	@which golangci-lint >/dev/null 2>&1 && golangci-lint run || echo "golangci-lint not installed in PATH; skipping"

run:
	go run ./cmd/monolith/main.go

run-example:
	go run ./cmd/example-service/main.go

up:
	docker compose -f deploy/docker-compose.yml up -d

down:
	docker compose -f deploy/docker-compose.yml down
