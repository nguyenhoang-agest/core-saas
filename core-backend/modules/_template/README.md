# Module Template

This directory serves as the template scaffold for generating new modules with `scripts/new-module.sh <name>`.

## Structure

- `module.go`: Wires dependencies and satisfies `kernel.Module`.
- `domain/`: Pure domain types and port interfaces.
- `app/`: Application use-case services.
- `adapters/`: Outbound adapter implementations (database, external clients).
- `http/`: Inbound HTTP handlers and route registration.
- `migrations/`: Embedded goose SQL migrations.

## Planned Ports & Events

- Inbound: HTTP REST endpoints.
- Outbound: Database repository.
- Events: Domain events published to `pkg/eventbus`.

## TODO Checklist

- [ ] Define domain entities and repository interfaces in `domain/`
- [ ] Implement use cases in `app/`
- [ ] Implement data access in `adapters/`
- [ ] Define HTTP request/response DTOs and routes in `http/`
- [ ] Add SQL schema migrations in `migrations/`
- [ ] Wire dependencies in `module.go`
