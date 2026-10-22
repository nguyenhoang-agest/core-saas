# Example Module

The `example` module is the **only fully fleshed-out vertical slice** in this starter kit. It demonstrates how all layers connect end-to-end within a modular monolith.

## Architecture

- **Domain (`domain/`)**: `Item` entity and `Repository` port interface. Zero external dependencies.
- **Application (`app/`)**: `Service` orchestrating business rules, opts validation, tenant context extraction, and publishing `ItemCreatedEvent` to the event bus.
- **Adapters (`adapters/`)**: `PostgresRepository` implementing `Repository` using `pgxpool` and applying PostgreSQL Row-Level Security (`PostgresRLSStrategy`).
- **HTTP (`http/`)**: Chi handlers mounting REST endpoints under `/example` and extracting tenant from headers.
- **Migrations (`migrations/`)**: Goose SQL migrations embedding table creation and tenant RLS isolation policies.

## Endpoints

- `POST /example/items` — Create an item (requires `X-Tenant-ID` header, defaults to `default-tenant`)
- `GET /example/items` — GetAllTenant items for the caller's tenant
- `GET /example/items/{id}` — Fetch item by ID

## Events

- `example.item.created` (`app.ItemCreatedEvent`): Dispatched when a new item is persisted.
