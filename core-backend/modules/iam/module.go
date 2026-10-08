package iam

import (
	"context"
	"embed"

	"github.com/example/core-backend/modules/iam/adapters"
	"github.com/example/core-backend/modules/iam/app"
	iamhttp "github.com/example/core-backend/modules/iam/http"
	"github.com/example/core-backend/pkg/eventbus"
	"github.com/example/core-backend/pkg/kernel"
	"github.com/example/core-backend/pkg/migrate"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

//go:embed migrations/*.sql
var migrationsFS embed.FS

var _ kernel.Module = (*Module)(nil)

// Module encapsulates the IAM module stub.
type Module struct {
	pool    *pgxpool.Pool
	bus     *eventbus.Bus
	service *app.Service
	handler *iamhttp.Handler
}

// New constructs the IAM module stub.
func New(pool *pgxpool.Pool, bus *eventbus.Bus) *Module {
	repo := adapters.NewPostgresUserRepository(pool)
	hasher := adapters.NewBcryptHasher()
	tokens := adapters.NewTokenStub()
	service := app.NewService(repo, hasher, tokens, bus)
	handler := iamhttp.NewHandler(service)

	return &Module{
		pool:    pool,
		bus:     bus,
		service: service,
		handler: handler,
	}
}

func (m *Module) Name() string {
	return "iam"
}

func (m *Module) Register(r chi.Router) error {
	r.Route("/iam", func(sub chi.Router) {
		m.handler.RegisterRoutes(sub)
	})
	return nil
}

func (m *Module) Migrate(ctx context.Context) error {
	return migrate.Run(ctx, m.pool, m.Name(), migrationsFS, "migrations")
}
