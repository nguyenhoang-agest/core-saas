# Tenant Module (STUB)

## Purpose
Manages SaaS organization workspaces (tenants), provisioning, lifecycle states (active, suspended), and tier/subscription mapping.

## Planned Ports & Interfaces
- **Ports (Domain)**:
  - `TenantRepository`: Querying and mutating tenant records.
  - `BillingGateway`: Provisioning customer billing profiles.
- **Adapters**:
  - `PostgresTenantRepository`: Persists tenant rows in PostgreSQL.
- **Events Emitted**:
  - `tenant.created`: When a new organization registers.
  - `tenant.suspended`: When an organization is suspended.

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
