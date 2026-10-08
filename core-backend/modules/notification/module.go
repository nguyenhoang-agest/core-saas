package notification

import (
	"context"
	"embed"

	"github.com/example/core-backend/modules/notification/adapters"
	"github.com/example/core-backend/modules/notification/app"
	notificationhttp "github.com/example/core-backend/modules/notification/http"
	"github.com/example/core-backend/pkg/eventbus"
	"github.com/example/core-backend/pkg/kernel"
	"github.com/example/core-backend/pkg/migrate"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

//go:embed migrations/*.sql
var migrationsFS embed.FS

var _ kernel.Module = (*Module)(nil)

// Module encapsulates the notification module stub.
type Module struct {
	pool    *pgxpool.Pool
	bus     *eventbus.Bus
	service *app.Service
	handler *notificationhttp.Handler
}

func New(pool *pgxpool.Pool, bus *eventbus.Bus) *Module {
	sender := adapters.NewLogSender()
	service := app.NewService(sender, bus)
	handler := notificationhttp.NewHandler(service)

	return &Module{
		pool:    pool,
		bus:     bus,
		service: service,
		handler: handler,
	}
}

func (m *Module) Name() string {
	return "notification"
}

func (m *Module) Register(r chi.Router) error {
	r.Route("/notifications", func(sub chi.Router) {
		m.handler.RegisterRoutes(sub)
	})
	return nil
}

func (m *Module) Migrate(ctx context.Context) error {
	return migrate.Run(ctx, m.pool, m.Name(), migrationsFS, "migrations")
}
