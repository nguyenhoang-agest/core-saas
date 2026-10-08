package template

import (
	"context"
	"embed"

	"github.com/example/core-backend/modules/_template/adapters"
	"github.com/example/core-backend/modules/_template/app"
	templatehttp "github.com/example/core-backend/modules/_template/http"
	"github.com/example/core-backend/pkg/eventbus"
	"github.com/example/core-backend/pkg/kernel"
	"github.com/example/core-backend/pkg/migrate"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

//go:embed migrations/*.sql
var migrationsFS embed.FS

var _ kernel.Module = (*Module)(nil)

// Module encapsulates the template module.
type Module struct {
	pool    *pgxpool.Pool
	bus     *eventbus.Bus
	service *app.Service
	handler *templatehttp.Handler
}

// New constructs the template module dependencies.
func New(pool *pgxpool.Pool, bus *eventbus.Bus) *Module {
	repo := adapters.NewRepository(pool)
	service := app.NewService(repo, bus)
	handler := templatehttp.NewHandler(service)

	return &Module{
		pool:    pool,
		bus:     bus,
		service: service,
		handler: handler,
	}
}

// Name returns the module identifier.
func (m *Module) Name() string {
	return "template"
}

// Register mounts HTTP endpoints.
func (m *Module) Register(r chi.Router) error {
	r.Route("/template", func(sub chi.Router) {
		m.handler.RegisterRoutes(sub)
	})
	return nil
}

// Migrate executes migrations.
func (m *Module) Migrate(ctx context.Context) error {
	return migrate.Run(ctx, m.pool, m.Name(), migrationsFS, "migrations")
}
