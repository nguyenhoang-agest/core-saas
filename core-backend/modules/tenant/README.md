# Tenant Module (STUB)

## Purpose
Manages SaaS organization workspaces (tenants), provisioning, lifecycle states (active, suspended), and tier/subscription mapping.

## Layout
- `module.go`: Module entry point, wires `internal/repository` -> `internal/service` -> `internal/handler`.
- `api/`: Public contract (`TenantCreatedTopic`, `TenantUpdatedTopic`).
- `internal/model/`: Tenant entity, BaseModel, and domain errors.
- `internal/service/`: Business use cases and `Repository` interface definition.
- `internal/repository/`: `PostgresTenantRepository` implementing `service.Repository`.
- `internal/handler/`: HTTP handlers and routing.
- `migrations/`: Embedded SQL migrations.

## Multi-Tenancy Strategy Notes
This starter kit defines the `tenantctx.Strategy` interface and implements the **Discriminator + Postgres RLS** pattern.
Future alternative strategies:
- Schema per tenant (`SET search_path = <tenant_schema>`)
- Database per tenant (dynamic connection pool router)

## TODO Checklist
- [ ] Implement `PostgresTenantRepository` with full CRUD
- [ ] Connect tenant creation to workspace slug validation (lowercase, alphanumeric, dashes)
- [ ] Implement tenant member association
- [ ] Integrate billing/stripe subscription webhook handling
