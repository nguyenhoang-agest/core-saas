# Registering a Module

This guide details how to register a business domain module into the monolith.

---

## Option 1: Automated Scaffolding & Registration

You can use `make new-module` with two modes:

### Mode A: Full Scaffolding (Code + Auto-Registration)
Scaffold the recommended Tenant-style architecture, unit tests, and automatically wire the module into the monolith:

```bash
make new-module NAME=billing
```

### Mode B: Directory Structure Only (No Code Generation)
If you want maximum flexibility to develop your own custom architecture and style without generating starter boilerplate code:

```bash
make new-module NAME=billing DIRS_ONLY=true
```

This creates the clean directory skeleton (`api/`, `internal/{dto,errcode,handler,model,repository,service}`, `migrations/`) with `.gitkeep` files and a `README.md`, leaving the implementation completely up to you. (Follow Option 2 below when you are ready to register it in `cmd/monolith/modules.go`).

---

## Option 2: Manual Registration (Step-by-Step)

If you are developing your own module with your own custom architecture, patterns, or structure, you only need to follow these 3 steps to register it with the platform:

### Step 1: Implement the `kernel.Module` Interface

In your module's root (e.g. `modules/<name>/module.go`), implement the `kernel.Module` contract:

```go
package <name>

import (
	"context"

	"github.com/example/core-backend/pkg/kernel"
	"github.com/gofiber/fiber/v3"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/example/core-backend/pkg/eventbus"
)

var _ kernel.Module = (*Module)(nil)

type Module struct {
	pool *pgxpool.Pool
	bus  *eventbus.Bus
}

func New(pool *pgxpool.Pool, bus *eventbus.Bus) *Module {
	return &Module{
		pool: pool,
		bus:  bus,
	}
}

// Name returns the unique module identifier.
func (m *Module) Name() string {
	return "<name>"
}

// Register mounts module HTTP endpoints onto Fiber router.
func (m *Module) Register(r fiber.Router) error {
	group := r.Group("/<name>")
	// Mount your HTTP handlers here:
	// group.Get("/", h.GetAll)
	// group.Post("/", h.Create)
	return nil
}

// Migrate executes database schema migrations (return nil if no migrations).
func (m *Module) Migrate(ctx context.Context) error {
	return nil
}
```

---

### Step 2: Register in `cmd/monolith/modules.go`

Open [`core-backend/cmd/monolith/modules.go`](../cmd/monolith/modules.go) and wire your module into `enabledModules`:

1. **Import your module**:
   ```go
   import (
       "github.com/example/core-backend/modules/<name>"
       // other module imports...
   )
   ```

2. **Add constructor call to `enabledModules` slice**:
   ```go
   func enabledModules(pool *pgxpool.Pool, bus *eventbus.Bus) []kernel.Module {
       tenantMod := tenant.New(pool, bus)
       cachedResolver := router.NewCachedTenantResolver(tenantMod.Resolver(), 5*time.Minute)
       dbRouter := router.NewDynamicDBRouter(pool, cachedResolver)

       return []kernel.Module{
           <name>.New(pool, bus),
           example.NewWithDBProvider(dbRouter, bus),
           iam.New(pool, bus),
           tenantMod,
           notification.New(pool, bus),
           audit.New(pool, bus),
       }
   }
   ```

---

### Step 3: Verify Registration

Verify that the monolith compiles cleanly with your newly registered module:

```bash
make test-backend
make build
```
