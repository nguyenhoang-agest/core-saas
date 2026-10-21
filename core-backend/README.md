# SaaS Starter Kit in Go 1.27

A production-grade **SaaS Starter Kit** scaffold written in Go 1.27. Engineered as a modular monolith by default with first-class microservice extraction capability, PostgreSQL Row-Level Security (RLS) multi-tenancy, and generator tooling.

---

## Highlights

- **Go 1.27 & Pure Idiomatic Go**: Zero dynamic reflection or classpath scanning; explicit dependency injection.
- **Modular Monolith**: Modules live in `modules/<name>` with strict hexagonal/clean isolation (`domain`, `app`, `adapters`, `http`, `migrations`).
- **Microservice-Ready**: Demonstrates extracting any module into a standalone binary (`cmd/example-service`).
- **Flexible Multi-Tenancy (Shared RLS & Database-per-tenant)**: Built-in discriminator tenant isolation via PostgreSQL Row-Level Security (`app.current_tenant_id`) alongside dynamic connection pool routing (`pkg/db/router`) for dedicated tenant databases with platform-enforced connection pooling (10 conns), mandatory TLS encryption (`sslmode=require`), and SSRF host protections.
- **Granular Multi-Tenant RBAC & ACL**: Complete role engine supporting immutable system roles (`owner`, `admin`, `member`, `viewer`), custom tenant roles, Anti-Lockout Owner protection, Privilege Escalation prevention (Subset Rule), and instant distributed cache invalidation via PostgreSQL `pg_notify`.
- **Event-Driven Lifecycle Orchestration**: Decoupled domain events (`tenant.created`, `tenant.suspended`) automating cross-module workflows (e.g. auto-syncing system roles and provisioning the initial Owner user).
- **Per-Module Migrations**: Embedded Goose SQL migrations with isolated migration tables (`goose_version_<module>`).
- **Typed In-Process Event Bus**: Pub/Sub leveraging Go generics for decoupled inter-module messaging.
- **Generator Tooling**: 
  - `make init PROJECT=<name> MODULE=<path> KEEP=<stubs>` to rename and prune.
  - `make new-module NAME=<name>` to scaffold new domain modules.
- **Architecture Boundaries**: Enforced via `golangci-lint` depguard rules.

---

## Directory Structure

```
.
├── cmd/
│   ├── monolith/main.go            # Loads all enabled modules
│   ├── monolith/modules.go         # Explicit registration list
│   └── example-service/main.go     # Loads only 'example' (proves split)
├── pkg/                            # Shared infra (zero domain logic)
│   ├── config/                     # YAML & environment loader
│   ├── db/                         # pgxpool connection, GORM, & security policies
│   │   ├── security.go             # DSN validation, SSRF blocklist, pool & SSL policies
│   │   └── router/                 # Dynamic multi-tenant DB pool router
│   ├── errors/                     # Structured errors with HTTP mappings
│   ├── eventbus/                   # Typed generic pub/sub
│   ├── httpx/                      # Fiber router, /healthz, responses
│   ├── kernel/                     # Module interface contract
│   ├── logger/                     # slog JSON logger
│   ├── migrate/                    # Goose per-module runner
│   └── tenantctx/                  # Tenant context & RLS strategy
├── modules/
│   ├── tenant/                     # Workspaces, isolate_mode & DB routing catalog (Full CRUD & Events)
│   ├── user/                       # Multi-tenant user directory, RBAC roles, ACL & Postgres NOTIFY
│   ├── example/                    # Reference domain vertical slice
│   ├── iam/                        # STUB: auth, sessions
│   ├── notification/               # STUB: dispatcher port + log adapter
│   ├── audit/                      # STUB: append-only audit trail
│   └── _template/                  # Canonical generator template
├── scripts/
│   ├── init.sh                     # Renames module/project, prunes stubs
│   └── new-module.sh               # Scaffolds new module from _template
├── deploy/
│   ├── docker-compose.yml          # Local PostgreSQL 17
│   └── Dockerfile                  # Distroless production container
├── docs/
│   ├── ARCHITECTURE.md             # System design & microservice split
│   ├── CONVENTIONS.md              # Code, errors, events & migrations
│   ├── ADDING_A_MODULE.md          # Guide to adding new domain modules
│   └── USING_THIS_KIT.md           # Getting started & upstream sync
├── Makefile
├── .golangci.yml
└── .github/workflows/ci.yml
```

---

## Quickstart

### 1. Initialize
```bash
make env      # Generate .env from .env.example
```

### 2. Start PostgreSQL
```bash
make db-up    # Starts container and waits for healthy status
```

### 3. Run Monolith
```bash
make run      # Automatically runs migrations and starts on :8080
```

### 4. Verify
```bash
curl http://localhost:8080/healthz
# {"checks":{"database":"ok"},"status":"ok"}
```

---

## Commands

| Command | Description |
|---|---|
| `make env` | Create `.env` from template if missing |
| `make db-up` (or `make up`) | Start local PostgreSQL 17 container in background |
| `make db-down` (or `make down`) | Stop local PostgreSQL container |
| `make db-status` | Check PostgreSQL container health |
| `make db-logs` | View PostgreSQL container logs |
| `make db-psql` | Open interactive `psql` session inside PostgreSQL container |
| `make run` | Start the monolith locally on `http://localhost:8080` |
| `make run-example` | Start the standalone example service on `http://localhost:8081` |
| `make test` | Run tests with race detection (`go test -race ./...`) |
| `make test-v` | Run tests with verbose output |
| `make test-coverage` | Run tests and generate coverage report |
| `make lint` | Run `golangci-lint` or fallback to `go vet` |
| `make tidy` | Tidy Go dependencies |
| `make build` | Compile all packages and binaries into `bin/` |
| `make clean` | Clean compiled binaries and test coverage profiles |
| `make new-module NAME=<name>` | Scaffold a new module from `_template` |
| `make init` | Rename module and project path, prune unused stubs |

---

## Documentation

- [Using This Kit](docs/USING_THIS_KIT.md)
- [Architecture](docs/ARCHITECTURE.md)
- [Conventions](docs/CONVENTIONS.md)
- [Adding a Module](docs/ADDING_A_MODULE.md)
