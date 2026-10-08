package example

import (
	"context"
	"embed"

	"github.com/example/core-backend/modules/example/adapters"
	"github.com/example/core-backend/modules/example/app"
	exithttp "github.com/example/core-backend/modules/example/http"
	"github.com/example/core-backend/pkg/eventbus"
	"github.com/example/core-backend/pkg/kernel"
	"github.com/example/core-backend/pkg/migrate"
	"github.com/example/core-backend/pkg/tenantctx"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

//go:embed migrations/*.sql
var migrationsFS embed.FS

var _ kernel.Module = (*Module)(nil)

// Module encapsulates the example domain slice.
type Module struct {
	pool    *pgxpool.Pool
	bus     *eventbus.Bus
	service *app.Service
	handler *exithttp.Handler
}

// New constructs and wires the example module dependencies explicitly.
func New(pool *pgxpool.Pool, bus *eventbus.Bus) *Module {
	strategy := tenantctx.NewPostgresRLSStrategy()
	repo := adapters.NewPostgresRepository(pool, strategy)
	service := app.NewService(repo, bus)
	handler := exithttp.NewHandler(service)

	return &Module{
		pool:    pool,
		bus:     bus,
		service: service,
		handler: handler,
	}
}

// Name returns the module identifier.
func (m *Module) Name() string {
	return "example"
}

// Register mounts the module HTTP endpoints under the given chi.Router.
func (m *Module) Register(r chi.Router) error {
	r.Route("/example", func(sub chi.Router) {
		m.handler.RegisterRoutes(sub)
	})
	return nil
}

// Migrate executes embedded database migrations for the example module.
func (m *Module) Migrate(ctx context.Context) error {
	return migrate.Run(ctx, m.pool, m.Name(), migrationsFS, "migrations")
}
