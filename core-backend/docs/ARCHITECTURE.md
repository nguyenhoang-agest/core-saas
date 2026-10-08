# Architecture

## 1. Overview & Diagram

The SaaS Starter Kit follows a **Modular Monolith** architecture designed for velocity and operational simplicity today, with zero-friction microservice extractability tomorrow.

```
                    +-----------------------------+
                    |        HTTP Clients         |
                    +--------------+--------------+
                                   |
                                   v
                    +-----------------------------+
                    |   Chi Router (pkg/httpx)    |
                    |   Middlewares & /healthz    |
                    +--------------+--------------+
                                   |
         +-----------------+-------+-------+-----------------+
         |                 |               |                 |
         v                 v               v                 v
  +--------------+  +--------------+  +--------------+  +--------------+
  |   example    |  |     iam      |  |    tenant    |  | notification |
  |    Module    |  |  (Stub Mod)  |  |  (Stub Mod)  |  |  (Stub Mod)  |
  +-------+------+  +-------+------+  +-------+------+  +-------+------+
          |                 |                 |                 |
          +-----------------+--------+--------+-----------------+
                                     |
                                     v
                       +-----------------------------+
                       |    In-Process Event Bus     |
                       |       (pkg/eventbus)        |
                       +-----------------------------+
                                     |
                                     v
                       +-----------------------------+
                       |       PostgreSQL Pool       |
                       |  Multi-Tenant RLS Strategy  |
                       |        (pkg/tenantctx)      |
                       +-----------------------------+
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

## 3. How to Split a Module into its Own Microservice

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
