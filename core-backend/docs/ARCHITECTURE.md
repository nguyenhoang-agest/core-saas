# Architecture

## 1. Overview & Diagram

The SaaS Starter Kit follows a **Modular Monolith** architecture designed for velocity and operational simplicity today, with zero-friction microservice extractability tomorrow and flexible **Multi-Tenant Data Isolation** (Shared RLS, Schema-per-tenant, and Dedicated Database-per-tenant).

```
                    +-----------------------------+
                    |        HTTP Clients         |
                    |    (Header: X-Tenant-ID)    |
                    +--------------+--------------+
                                   |
                                   v
                    +-----------------------------+
                    |    Fiber / httpx Router     |
                    |    Middlewares & /healthz   |
                    +--------------+--------------+
                                   |
         +-----------------+-------+-------+-----------------+
         |                 |               |                 |
         v                 v               v                 v
  +--------------+  +--------------+  +--------------+  +--------------+
  |   example    |  |     iam      |  |    tenant    |  | notification |
  |    Module    |  |  (Stub Mod)  |  |  (Catalog)   |  |  (Stub Mod)  |
  +-------+------+  +-------+------+  +-------+------+  +-------+------+
          |                 |                 |                 |
          +--------+--------+                 |                 |
                   |                          |                 |
                   v                          v                 v
     +---------------------------+      +-------------------------------+
     |  Dynamic DB Router        |      |     In-Process Event Bus      |
     |  (pkg/db/router)          |      |        (pkg/eventbus)         |
     +-------------+-------------+      +-------------------------------+
                   |
     +-------------+-------------+
     |                           |
[ isolate_mode: shared ]    [ isolate_mode: dedicated ]
     |                           |
     v                           v
+---------------------+     +---------------------+
| Primary / Default   |     | Dedicated DB Pool   |
| Database Pool       |     | (Tenant PostgreSQL) |
+---------------------+     +---------------------+
| - tenants catalog   |     | - isolated items    |
| - shared data (RLS) |     | - isolated users    |
+---------------------+     +---------------------+
```

---

## 2. Module Boundaries & Rules

1. **Self-Contained Structure**: Every module contains its own:
   - `domain/`: Business entities and ports (repository and service interfaces). Zero dependencies on external frameworks.
   - `app/`: Use-case coordinators and domain event publication.
   - `adapters/`: Outbound implementations (SQL storage, external APIs).
   - `http/`: Inbound HTTP handlers and route registration.
   - `migrations/`: Embedded SQL migrations isolated in a per-module Goose version table.
2. **Explicit Dependency Injection**: Modules expose a `New(pool, bus)` constructor and implement `kernel.Module`. Reflection and dynamic scanning are forbidden.
3. **Internal Package Isolation**: Modules must **NEVER** import another module's internal packages (`domain/`, `app/`, `adapters/`, `http/`). Cross-module interactions occur strictly via:
   - Exported root package contracts or DTOs.
   - Decoupled asynchronous events via `pkg/eventbus`.
4. **Infra Isolation**: Packages in `pkg/` must never import anything from `modules/*`. This boundary is enforced by `golangci-lint` with `depguard`.

---

## 3. Multi-Tenant Architecture & Dynamic DB Routing

The platform supports a hybrid, multi-tier isolation model where tenants can select different isolation strategies based on compliance, performance, or subscription tier.

### 3.1. Tenant Isolation Modes (`isolate_mode`)

Each tenant record in the catalog specifies an `isolate_mode`:

| Mode | Database | Schema | Isolation Mechanism | Best Suited For |
|---|---|---|---|---|
| `shared` | Shared Primary DB | Shared Schema | PostgreSQL Row-Level Security (`app.current_tenant_id`) | Free / Starter / Standard tiers |
| `schema` | Shared Primary DB | Separate Schema | Search path (`SET LOCAL search_path = <tenant>`) | Mid-tier / Team workspaces |
| `dedicated` | Separate DB Instance/Cluster | Isolated DB | Dynamic Connection Pool Router (`pkg/db/router`) | Enterprise / High-compliance / Data sovereignty |

### 3.2. Control Plane vs. Data Plane

1. **Control Plane (Catalog Database)**:
   - Lives in the Primary PostgreSQL database (`DefaultPool`).
   - Hosts system catalog tables: `tenants`, subscriptions, global billing, and platform audit logs.
   - The `tenant` module always connects to the Control Plane.

2. **Data Plane (Business Data)**:
   - Hosts domain tables: `example_items`, `iam_users`, `notifications`, etc.
   - For `shared` tenants, resides in the Primary DB alongside RLS policies.
   - For `dedicated` tenants, resides in the tenant's dedicated PostgreSQL database.

### 3.3. Dynamic DB Connection Pool Router (`pkg/db/router`)

Modules interact with the database via the **`router.DBProvider`** interface:
```go
type DBProvider interface {
    GetPool(ctx context.Context) (*pgxpool.Pool, error)
    GetDefaultPool() *pgxpool.Pool
    Close()
}
```

Key runtime features of `DynamicDBRouter`:
- **Context-Aware Routing**: Inspects `tenantctx.FromContext(ctx)` to determine tenant ID.
- **Tenant Resolver & In-Memory TTL Cache**: Resolves `TenantRouteInfo` (`IsolateMode`, `DSN`, `MaxConns`) via `CachedTenantResolver` without repeating catalog lookups on every request.
- **Lazy Connection**: Pools for dedicated tenants are opened on-demand upon first incoming request.
- **Thread-Safe Pooling**: Managed via `sync.RWMutex` protecting `map[string]*pgxpool.Pool`.
- **Idle Pool Reaper**: A background worker periodically reaps connection pools that have been idle past `idleTimeout` (default: 30 minutes), releasing database file descriptors and memory.
- **Graceful Shutdown**: Automatically drains and closes all opened tenant pools upon application termination.

### 3.4. Tenant Promotion & Migration Workflow (Shared to Dedicated)

When an existing `shared` tenant upgrades to an enterprise `dedicated` database:

1. **Set Tenant State**: Update tenant status to `migrating` (blocks writes, read-only mode).
2. **Extract Tenant Data**: Dump rows matching `tenant_id = X` across all domain tables.
3. **Run Schema Migrations**: Execute Goose migrations for all business modules onto the new dedicated database.
4. **Load Data**: Insert the exported rows into the dedicated database.
5. **Update Routing Catalog**: Update `tenants` with `isolate_mode = 'dedicated'` and `db_config.dsn`.
6. **Invalidate Cache**: Broadcast `tenant.updated` event over `eventbus` to flush `CachedTenantResolver`.
7. **Purge Old Data**: Soft-delete or archive tenant records in the shared database.
8. **Re-activate**: Set tenant status back to `active`.

---

## 4. How to Split a Module into its Own Microservice

The starter kit already demonstrates this pattern with `cmd/example-service/main.go`.

When a module experiences high traffic or requires independent deployment:

1. **Create an Entrypoint**:
   Add `cmd/<module-service>/main.go` which imports only your module and the required `pkg/` shared libraries.
2. **Swap the Event Bus**:
   Replace the in-process `pkg/eventbus` in that entrypoint with an external message broker client (e.g. NATS, RabbitMQ, or Apache Kafka).
3. **Database Extraction**:
   Because migrations already use per-module tables (`goose_version_<module>`), moving the module's schema to a separate PostgreSQL database or instance requires only pointing `DATABASE_URL` to the new database.
4. **Network Routing**:
   Update your reverse proxy or API Gateway (e.g. Kong, Envoy, Traefik) to route traffic for `/<module>` directly to the new service container instead of the monolith.
