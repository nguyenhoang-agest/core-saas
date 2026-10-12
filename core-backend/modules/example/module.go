package example

import (
	"context"
	"embed"

	"github.com/example/core-backend/modules/example/adapters"
	"github.com/example/core-backend/modules/example/app"
	exithttp "github.com/example/core-backend/modules/example/http"
	"github.com/example/core-backend/pkg/db/router"
	"github.com/example/core-backend/pkg/eventbus"
	"github.com/example/core-backend/pkg/kernel"
	"github.com/example/core-backend/pkg/migrate"
	"github.com/example/core-backend/pkg/tenantctx"
	"github.com/gofiber/fiber/v3"
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
	return NewWithDBProvider(router.NewStaticDBProvider(pool), bus)
}

// NewWithDBProvider constructs the example module with dynamic DBProvider routing.
func NewWithDBProvider(dbProvider router.DBProvider, bus *eventbus.Bus) *Module {
	strategy := tenantctx.NewPostgresRLSStrategy()
	repo := adapters.NewPostgresRepositoryWithDBProvider(dbProvider, strategy)
	service := app.NewService(repo, bus)
	handler := exithttp.NewHandler(service)

	var pool *pgxpool.Pool
	if dbProvider != nil {
		pool = dbProvider.GetDefaultPool()
	}

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

// Register mounts the module HTTP endpoints under the given fiber.Router.
func (m *Module) Register(r fiber.Router) error {
	return nil
}

// Migrate executes embedded database migrations for the example module.
func (m *Module) Migrate(ctx context.Context) error {
	return migrate.Run(ctx, m.pool, m.Name(), migrationsFS, "migrations")
}
