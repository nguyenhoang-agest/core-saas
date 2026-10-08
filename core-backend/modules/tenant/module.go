package tenant

import (
	"context"
	"embed"

	"github.com/example/core-backend/modules/tenant/adapters"
	tenanthttp "github.com/example/core-backend/modules/tenant/http"
	"github.com/example/core-backend/modules/tenant/service"
	"github.com/example/core-backend/pkg/eventbus"
	"github.com/example/core-backend/pkg/kernel"
	"github.com/example/core-backend/pkg/migrate"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

//go:embed migrations/*.sql
var migrationsFS embed.FS

var _ kernel.Module = (*Module)(nil)

// Module encapsulates the tenant management module.
type Module struct {
	pool    *pgxpool.Pool
	bus     *eventbus.Bus
	service service.ITenantService
	handler *tenanthttp.Handler
}

func New(pool *pgxpool.Pool, bus *eventbus.Bus) *Module {
	repo := adapters.NewPostgresTenantRepository(pool)
	svc := service.NewTenantService(repo, bus)
	handler := tenanthttp.NewHandler(svc)

	return &Module{
		pool:    pool,
		bus:     bus,
		service: svc,
		handler: handler,
	}
}

func (m *Module) Name() string {
	return "tenant"
}

func (m *Module) Register(r chi.Router) error {
	r.Route("/tenants", func(sub chi.Router) {
		m.handler.RegisterRoutes(sub)
	})
	return nil
}

func (m *Module) Migrate(ctx context.Context) error {
	return migrate.Run(ctx, m.pool, m.Name(), migrationsFS, "migrations")
}
