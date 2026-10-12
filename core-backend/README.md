# SaaS Starter Kit in Go 1.27

A production-grade **SaaS Starter Kit** scaffold written in Go 1.27. Engineered as a modular monolith by default with first-class microservice extraction capability, PostgreSQL Row-Level Security (RLS) multi-tenancy, and generator tooling.

---

## Highlights

- **Go 1.27 & Pure Idiomatic Go**: Zero dynamic reflection or classpath scanning; explicit dependency injection.
- **Modular Monolith**: Modules live in `modules/<name>` with strict hexagonal/clean isolation (`domain`, `app`, `adapters`, `http`, `migrations`).
- **Microservice-Ready**: Demonstrates extracting any module into a standalone binary (`cmd/example-service`).
- **Flexible Multi-Tenancy (Shared RLS & Database-per-tenant)**: Built-in discriminator tenant isolation via PostgreSQL Row-Level Security (`app.current_tenant_id`) alongside dynamic connection pool routing (`pkg/db/router`) for dedicated tenant databases.
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
│   ├── db/                         # pgxpool connection & transactions
│   │   └── router/                 # Dynamic multi-tenant DB pool router
│   ├── errors/                     # Structured errors with HTTP mappings
│   ├── eventbus/                   # Typed generic pub/sub
│   ├── httpx/                      # Fiber router, /healthz, responses
│   ├── kernel/                     # Module interface contract
│   ├── logger/                     # slog JSON logger
│   ├── migrate/                    # Goose per-module runner
│   └── tenantctx/                  # Tenant context & RLS strategy
├── modules/
│   ├── example/                    # The ONLY fleshed-out vertical slice
│   ├── iam/                        # STUB: auth, users, sessions
│   ├── tenant/                     # Workspaces, isolate_mode & DB routing catalog
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
make init PROJECT=my-saas MODULE=github.com/my-org/my-saas KEEP=iam,tenant
```

### 2. Start PostgreSQL
```bash
make up
```

### 3. Run Monolith
```bash
make run
```

### 4. Verify
```bash
curl http://localhost:8080/healthz
```

---

## Commands

| Command | Description |
|---|---|
| `make init` | Rename module and project path, prune unused stubs |
| `make new-module NAME=<name>` | Scaffold a new module from `_template` |
| `make build` | Compile all packages and binaries |
| `make test` | Run tests with race detection |
| `make lint` | Run `golangci-lint` |
| `make tidy` | Tidy Go dependencies |
| `make run` | Start the monolith locally |
| `make run-example` | Start the standalone example service |
| `make up` | Start local PostgreSQL |
| `make down` | Stop local PostgreSQL |

---

## Documentation

- [Using This Kit](docs/USING_THIS_KIT.md)
- [Architecture](docs/ARCHITECTURE.md)
- [Conventions](docs/CONVENTIONS.md)
- [Adding a Module](docs/ADDING_A_MODULE.md)
