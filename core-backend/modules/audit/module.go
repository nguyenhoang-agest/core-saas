package audit

import (
	"context"
	"embed"

	"github.com/example/core-backend/modules/audit/adapters"
	"github.com/example/core-backend/modules/audit/app"
	audithttp "github.com/example/core-backend/modules/audit/http"
	"github.com/example/core-backend/pkg/eventbus"
	"github.com/example/core-backend/pkg/kernel"
	"github.com/example/core-backend/pkg/migrate"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

//go:embed migrations/*.sql
var migrationsFS embed.FS

var _ kernel.Module = (*Module)(nil)

// Module encapsulates the audit logging module stub.
type Module struct {
	pool    *pgxpool.Pool
	bus     *eventbus.Bus
	service *app.Service
	handler *audithttp.Handler
}

func New(pool *pgxpool.Pool, bus *eventbus.Bus) *Module {
	repo := adapters.NewPostgresAuditRepository(pool)
	service := app.NewService(repo, bus)
	handler := audithttp.NewHandler(service)

	return &Module{
		pool:    pool,
		bus:     bus,
		service: service,
		handler: handler,
	}
}

func (m *Module) Name() string {
	return "audit"
}

func (m *Module) Register(r chi.Router) error {
	r.Route("/audit", func(sub chi.Router) {
		m.handler.RegisterRoutes(sub)
	})
	return nil
}

func (m *Module) Migrate(ctx context.Context) error {
	return migrate.Run(ctx, m.pool, m.Name(), migrationsFS, "migrations")
}
