# Tenant Module

## Purpose
Manages SaaS organization workspaces (tenants), provisioning lifecycle, workspace states (`active`, `suspended`), isolation modes (`shared`, `schema`, `dedicated`), and dynamic database routing configuration.

---

## Layout
- `module.go`: Module entry point, wires `internal/repository` -> `internal/service` -> `internal/handler` and exports `Resolver()`.
- `resolver.go`: Adapter implementing `router.TenantResolver` to provide routing information to `pkg/db/router`.
- `api/events.go`: Public contracts and topics (`TenantCreatedTopic`, `TenantUpdatedTopic`, `TenantSuspendedTopic`, `TenantDeletedTopic`).
- `internal/model/`: Tenant entity, `IsolateMode`, `TenantDBConfig`, BaseModel, and domain errors.
- `internal/dto/`: Request/Response data transfer objects (`CreateTenantRequest`, `UpdateTenantRequest`, `UpdateTenantStatusRequest`).
- `internal/service/`: Business use cases and `Repository` interface definition.
- `internal/repository/`: `PostgresTenantRepository` implementing `service.Repository`.
- `internal/handler/`: HTTP handlers and Fiber routing guarded by RBAC permissions.
- `migrations/`: Embedded SQL migrations (`00001_create_tenants_table.sql`, `00002_add_isolate_mode_to_tenants.sql`).

---

## Multi-Tenancy Isolation Modes (`isolate_mode`)

When registering a tenant, the client can specify an `isolate_mode`:

| Mode | Value | Description |
|---|---|---|
| **Shared** (Default) | `"shared"` | Tenant shares the primary database and schema. Tenant rows are filtered using PostgreSQL Row-Level Security (RLS) via `tenantctx.PostgresRLSStrategy`. |
| **Schema-per-tenant** | `"schema"` | Tenant data is placed in a dedicated schema within the shared database instance (`search_path`). |
| **Database-per-tenant** | `"dedicated"` | Tenant data resides in an isolated PostgreSQL database. Queries are dynamically routed via `pkg/db/router.DynamicDBRouter`. |

---

## API Endpoints & RBAC Permissions

### 1. Create Tenant (`POST /tenants`)
- **Permission**: `tenant:create` (Platform-scoped)
- Automatically triggers `TenantCreatedTopic` event containing `owner` metadata.
- Inter-module subscriber in `modules/user` automatically executes `SyncSystemRoles` and provisions the initial `owner` user account.

#### Option A: Shared Multi-Tenant (Standard)
```json
{
  "name": "Acme Corp",
  "description": "Standard organization using shared database",
  "isolate_mode": "shared",
  "owner": {
    "name": "Acme Admin",
    "email": "admin@acme.com"
  }
}
```

#### Option B: Dedicated Database (Enterprise / Isolated)
```json
{
  "name": "Enterprise Global",
  "description": "Enterprise customer with dedicated database",
  "isolate_mode": "dedicated",
  "db_config": {
    "dsn": "postgres://user:password@dedicated-db-host:5432/enterprise_db?sslmode=require"
  },
  "owner": {
    "name": "Enterprise Owner",
    "email": "owner@enterprise.com"
  }
}
```

**Security & Governance Policies for Dedicated Mode:**
- **DSN Validation & SSRF Prevention (`pkg/db.ValidateDSN`)**: 
  - `db_config.dsn` is required, must be a valid PostgreSQL connection string, and cannot be empty.
  - Prohibits loopback addresses (`localhost`, `127.0.0.0/8`, `::1`, `0.0.0.0`) and cloud instance metadata (`169.254.169.254`) to prevent Server-Side Request Forgery (SSRF).
- **Connection Pool Policy**: Strictly platform-managed. Any user-provided connection pool parameter is denied/overridden and fixed to `DefaultTenantMaxConns = 10` to prevent database connection exhaustion.
- **SSL Security Mode**: Strictly platform-managed. Changing SSL mode (e.g. to `disable`) is prohibited. All dedicated databases are enforced to `DefaultTenantSSLMode = "require"` (TLS encryption required).

---

### 2. Update Tenant (`PUT /tenants/:id`)
- **Permission**: `tenant:update`
- Modifies workspace name and description.
- Emits `TenantUpdatedTopic` over `eventbus`.

```json
{
  "name": "Acme Corporation Ltd",
  "description": "Updated global headquarters workspace"
}
```

---

### 3. Update Tenant Status (`PATCH /tenants/:id/status`)
- **Permission**: `tenant:update`
- Sets workspace status to `"active"` or `"suspended"`.
- Emits `TenantUpdatedTopic` or `TenantSuspendedTopic`.
- When suspended, active user tokens and pooled connections are invalidated.

```json
{
  "status": "suspended"
}
```

---

### 4. Delete Tenant (`DELETE /tenants/:id`)
- **Permission**: `tenant:delete` (Platform-scoped)
- Soft-deletes tenant record and emits `TenantDeletedTopic`.
- Clears routing cache and closes any open dedicated connection pools.

---

### 5. Get Tenant by ID (`GET /tenants/:id`)
- **Permission**: `tenant:view`
- Retrieves tenant workspace details including current `status` and `isolate_mode`.

---

### 6. List Tenants (`GET /tenants`)
- **Permission**: `tenant:view`
- Retrieves paginated tenant workspaces with optional status and name filtering.

---

## Dynamic Router Integration

The tenant module exposes:
```go
func (m *Module) Resolver() router.TenantResolver
```
This is connected to `router.NewCachedTenantResolver` in `cmd/monolith/modules.go`. Whenever any business module executes a query scoped to an active `X-Tenant-ID`:
1. The router resolves the tenant's `isolate_mode`.
2. For `dedicated` tenants, it lazily creates and caches a dedicated `*pgxpool.Pool` directed at the tenant's `db_config.dsn`.
3. Background idle pool reaper closes connection pools that have been idle past 30 minutes.
