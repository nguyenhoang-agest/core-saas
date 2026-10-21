# Core SaaS — Modular Multi-Tenant Enterprise Platform

An enterprise-grade, production-ready **Modular Multi-Tenant SaaS Platform**. Engineered with a hardened **Go Fiber v3** backend with dynamic database connection pool routing, alongside a **Next.js 16** admin console designed with an **Outlined Minimalist / Flat 2.0 Monochrome UI**.

---

## Monorepo Architecture

```
Core-saas/
├── core-backend/                   # Go 1.27 Modular Monolith
│   ├── cmd/
│   │   ├── monolith/               # Full application server (loads all modules)
│   │   └── example-service/        # Standalone microservice proof-of-concept
│   ├── pkg/                        # Domain-agnostic shared infrastructure
│   │   ├── config/                 # .env & YAML loader with fallback resolution
│   │   ├── db/                     # pgxpool connection, GORM wrapper & security policies
│   │   │   ├── security.go         # DSN validation, SSRF blocklist, pool & TLS policies
│   │   │   └── router/             # Dynamic multi-tenant DB pool router & tenant resolver
│   │   ├── exceptions/             # 3-layer typed domain errors (AppError, ErrorCode)
│   │   ├── httpx/                  # Fiber server, health checks, centralized error handler
│   │   ├── i18n/                   # Decoupled i18n engine with go-i18n & Fiber middleware
│   │   └── eventbus/               # Type-safe generic in-process event bus
│   ├── modules/
│   │   ├── tenant/                 # Tenant workspace catalog, isolation mode & DB configs
│   │   ├── user/                   # Multi-tenant user directory, RBAC roles, ACL & Postgres NOTIFY
│   │   ├── example/                # Reference domain vertical slice
│   │   ├── iam/                    # Auth and session management
│   │   ├── notification/           # Notification dispatcher & log adapter
│   │   └── audit/                  # Append-only audit logging
│   └── deploy/
│       └── docker-compose.yml      # Local PostgreSQL 17 container
│
├── core-frontend/                  # Next.js 16 Admin Console (AdminCN Style)
│   ├── src/
│   │   ├── app/                    # Next.js App Router (/, /tenants, /users, /billing, /settings)
│   │   ├── components/             # Minimalist UI primitives, TenantSwitcher, LocaleSwitcher
│   │   ├── core/
│   │   │   ├── i18n/               # React i18n context, LocaleSwitcher, parameter interpolation
│   │   │   └── tenant/             # Active tenant context provider & cookie sync
│   │   ├── locales/                # Dictionaries: en/ & vi/ (errors.json, common.json)
│   │   └── modules/
│   │       ├── tenant/             # Tenant Inspector, KPI Cards, Create/Edit Dialogs
│   │       └── user/               # User directory, Role management, Permission matrix
│   └── config/app.config.ts        # Modular UI registry configuration
│
├── scripts/
│   └── setup.sh                    # One-click environment bootstrap script
├── Makefile                        # Root unified command orchestrator
└── README.md                       # Workspace documentation (this file)
```

---

## Quick Start in 3 Minutes

### 1. Prerequisites
- **Go** >= 1.24 (Recommended: Go 1.27)
- **Node.js** >= 18.x & **npm**
- **Docker** & **Docker Compose**

### 2. One-Click Bootstrap
From the repository root, run:

```bash
make setup
```

This single command will:
1. Generate `core-backend/.env` from `.env.example`.
2. Generate `core-frontend/.env.local` from `.env.example`.
3. Start the local **PostgreSQL 17** container via Docker Compose and wait until it is healthy.

### 3. Launch Development Servers

In **Terminal 1** (Start Backend Monolith on `:8080`):
```bash
make dev-backend
```
*(Backend automatically connects to PostgreSQL and applies all pending module Goose migrations!)*

In **Terminal 2** (Start Frontend Admin Console on `:3000`):
```bash
make dev-frontend
```

Open your browser at **[http://localhost:3000](http://localhost:3000)** to access the Admin Console.

---

## Complete Makefile Cheat Sheet

### Root Commands (`/Makefile`)

| Command | Description |
| :--- | :--- |
| `make help` | Display colored interactive help banner with all available targets |
| `make setup` | **One-click setup**: generate `.env` files and start PostgreSQL container |
| `make env` | Create `.env` and `.env.local` from templates if not already present |
| `make db-up` | Start local PostgreSQL 17 container in background |
| `make db-down` | Stop and tear down local PostgreSQL container |
| `make db-status` | Inspect PostgreSQL container health and port 5432 bindings |
| `make db-logs` | Stream live PostgreSQL container logs |
| `make db-psql` | Launch interactive `psql` CLI session inside PostgreSQL container |
| `make db-migrate-status` | Display applied migration versions for all modules |
| `make dev-backend` | Start Go Fiber backend monolith on `http://localhost:8080` |
| `make dev-frontend` | Start Next.js 16 frontend admin on `http://localhost:3000` |
| `make test` | **Run all test suites** across both backend (Go) and frontend (Vitest) |
| `make test-backend` | Run Go backend tests with race detection (`go test -race ./...`) |
| `make test-frontend` | Run frontend Vitest test suite (`npm test`) |
| `make lint` | Run code quality linters on both backend and frontend |
| `make build` | Compile backend binaries into `bin/` and build frontend Next.js bundle |
| `make clean` | Remove compiled binaries, coverage profiles, and `.next` cache |
| `make new-module NAME=<name>` | Scaffold a new backend domain module from template |
| `make new-migration MODULE=<mod> NAME=<name>` | Create a sequential Goose migration file |

---

## Environment Configuration

### Backend (`core-backend/.env`)

| Variable | Default Value | Description |
| :--- | :--- | :--- |
| `PORT` | `8080` | HTTP port for Fiber server |
| `ENV` | `development` | Environment mode (`development`, `production`, `test`) |
| `LOG_LEVEL` | `debug` | Slog logging verbosity (`debug`, `info`, `warn`, `error`) |
| `DATABASE_URL` | `postgres://postgres:postgres@localhost:5432/saas_dev?sslmode=disable` | Primary PostgreSQL connection DSN |
| `DB_MAX_CONNS` | `10` | Maximum connection pool size |
| `DB_HOST`, `DB_PORT`, `DB_USER`, `DB_PASSWORD`, `DB_NAME` | *(Optional)* | Granular settings used if `DATABASE_URL` is omitted |

### Frontend (`core-frontend/.env.local`)

| Variable | Default Value | Description |
| :--- | :--- | :--- |
| `PORT` | `3000` | Frontend web server port |
| `BACKEND_INTERNAL_URL` | `http://localhost:8080` | Target URL for Next.js internal `/api/*` rewrite proxy |
| `NEXT_PUBLIC_API_URL` | `/api` | Base API path consumed by client React Query hooks |

---

## Multi-Tenant Data Isolation

Core SaaS provides two levels of tenant isolation out of the box:

```
                  ┌────────────────────────────────────────┐
                  │          HTTP Request Header           │
                  │             X-Tenant-ID: <id>          │
                  └──────────────────┬─────────────────────┘
                                     │
                                     ▼
                  ┌────────────────────────────────────────┐
                  │       pkg/db/router.DynamicDBRouter    │
                  │    Resolves tenant isolation mode      │
                  └─────────┬────────────────────┬─────────┘
                            │                    │
              isolate_mode="shared"     isolate_mode="dedicated"
                            │                    │
                            ▼                    ▼
                 ┌───────────────────┐  ┌───────────────────┐
                 │ Shared Pool       │  │ Dedicated Pool    │
                 │ Multi-tenant DB   │  │ Isolated DB DSN   │
                 │ Row-Level Scoped  │  │ Idle Pool Reaping │
                 └───────────────────┘  └───────────────────┘
```

1. **Shared Multi-Tenant (Default)**:
   - Tenants share a unified PostgreSQL database pool (`saas_dev`).
   - Every operation is automatically scoped with tenant context extracted from the `X-Tenant-ID` header.
2. **Dedicated Database Mode & Platform Security Governance**:
   - The dynamic DB router (`pkg/db/router`) resolves the tenant record and establishes an isolated connection pool using the tenant's dedicated connection DSN.
   - **Platform-Managed Connection Pool**: Users are strictly denied from configuring `max_conns`; all dedicated pools enforce `DefaultTenantMaxConns = 10` to prevent database connection exhaustion.
   - **Mandatory TLS Encryption**: Users are strictly denied from configuring or disabling SSL; all dedicated connections enforce `DefaultTenantSSLMode = "require"`.
   - **SSRF & Network Defenses**: Target hostnames are inspected by `pkg/db.ValidateDSN` to block loopback targets (`localhost`, `127.0.0.0/8`, `::1`, `0.0.0.0`) and cloud instance metadata services (`169.254.169.254`).
   - Built-in idle pool reaper automatically closes connections for idle tenants to conserve database file descriptors and memory.

---

## Multi-Tenant RBAC & Real-Time Authorization

The user & authorization module (`core-backend/modules/user`) provides an enterprise-ready RBAC engine:

1. **System Roles vs. Custom Tenant Roles (`is_system`)**:
   - Each tenant contains 4 built-in system roles: `owner`, `admin`, `member`, and `viewer` (`is_system = true`).
   - System roles are **immutable**: edit and delete operations are rejected at both API and repository levels (`E032007`), preventing workspace lockout.
   - Administrators can create custom roles (`is_system = false`) or clone from system templates.
   - System roles automatically sync newly added module permissions via `SyncSystemRoles` when the platform upgrades.
2. **Security & Safety Invariants**:
   - **Anti-Lockout Owner Protection**: Guaranteed invariant that every tenant maintains at least one active user with the `owner` role. Revocation or deletion is denied if it violates this rule (`E032013`).
   - **Privilege Escalation Guard (Subset Rule)**: An administrator cannot grant roles or permissions they do not possess (`E032009`).
   - **Scope Boundaries**: Distinguishes between `ScopeTenant` permissions and `ScopePlatform` permissions.
3. **Real-Time Distributed Invalidation via PostgreSQL `pg_notify`**:
   - Updates atomically bump `authz_version` and broadcast `pg_notify('authz_invalidation', tenant_id)`.
   - Clustered instances immediately evict local permission caches with zero cache-staleness delay.
4. **Header Propagation & Switcher**:
   - Next.js edge proxy (`src/proxy.ts`) prioritizes incoming `X-Tenant-ID` headers before falling back to subdomains and cookies.
   - Topbar `TenantSwitcher` allows instant workspace switching with TanStack Query cache purging.

---

## i18n & Two-Tier Error Handling

The platform features an end-to-end **internationalization (i18n)** system supporting **English (`en`)** and **Vietnamese (`vi`)**:

- **Backend ([`pkg/exceptions`](core-backend/pkg/exceptions), [`pkg/httpx`](core-backend/pkg/httpx))**:
  - Emits standardized, machine-readable `message_code` (e.g. `E031001` - `E031006`).
  - Zero module locale duplication: domain modules do not maintain localized string dictionaries.
  - Fiber middleware automatically negotiates locale via query `?lang=`, header `X-Locale`, or `Accept-Language`.
- **Frontend ([`src/core/i18n`](core-frontend/src/core/i18n))**:
  - Two-tier error resolution: Client dictionary lookup (`src/locales/{en,vi}/errors.json`) $\rightarrow$ Backend returned message $\rightarrow$ Generic fallback.
  - Function `t(key, params)` supports parameter interpolation (e.g. `t('tenant.page.created_success', { name: 'Acme' })`).
  - Outlined Minimalist `LocaleSwitcher` on the topbar for instant language switching with cookie & `localStorage` persistence.

---

## Verification & Testing

Verify the entire repository with a single command:

```bash
make test
```

Or verify components individually:
```bash
# Backend tests
cd core-backend && go test -v ./...

# Frontend tests & linting
cd core-frontend
npx tsc --noEmit    # TypeScript check (0 errors)
npm run lint        # ESLint check (0 errors)
npm test            # Vitest suite (14 test files, 51 tests pass)
npm run build       # Next.js 16 Turbopack build
```
