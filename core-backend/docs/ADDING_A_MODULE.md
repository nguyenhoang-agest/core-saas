# Adding a New Module

This guide details how to add a new business domain module into the starter kit.

---

## 1. Run the Generator

Run `make new-module` with your module's name:

```bash
make new-module NAME=billing
```

This creates the directory `modules/billing/` and registers the module in `cmd/monolith/modules.go`.

---

## 2. Implement the Layers

### Step 2.1: Domain Layer (`domain/`)
Define your domain entities, values, and outbound ports in `modules/billing/domain/`:
- Entities: plain Go structs without JSON/database tags or framework imports.
- Ports: repository and external service interfaces.

### Step 2.2: Application Layer (`app/`)
Implement application use cases and business workflows in `modules/billing/app/`:
- Coordinates data access via domain interfaces.
- Validates inputs and returns `pkg/errors` types.
- Publishes domain events to `pkg/eventbus`.

### Step 2.3: Adapters Layer (`adapters/`)
Implement outbound ports defined in domain:
- Database repositories using `router.DBProvider` (or `pgxpool.Pool`) and `tenantctx.Strategy`. Using `router.DBProvider` enables automatic dynamic routing to dedicated databases when tenants use `isolate_mode: "dedicated"`.
- Third-party client integrations (e.g. Stripe, SendGrid).

### Step 2.4: HTTP Layer (`http/`)
Expose inbound HTTP REST endpoints using Fiber:
- Bind requests with `c.Bind().Body(&opts)`.
- Extract tenant context with `tenantctx.FromContext(c.Context())`.
- Return responses using `c.JSON(httpx.NewResponse(data))` or `httpx.NewResponseError(code, err)`.

### Step 2.5: Database Migrations (`migrations/`)
Add Goose SQL migration files in `modules/billing/migrations/`:
- Create tables with a `tenant_id` column.
- Enable PostgreSQL Row Level Security (RLS) policies.

---

## 3. Verify Compilation

Verify your new module builds cleanly:

```bash
make build
make test
```
