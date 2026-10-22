package template

import (
	"context"
	"embed"

	"github.com/example/core-backend/modules/_template/internal/handler"
	"github.com/example/core-backend/modules/_template/internal/repository"
	"github.com/example/core-backend/modules/_template/internal/service"
	"github.com/example/core-backend/pkg/db"
	"github.com/example/core-backend/pkg/eventbus"
	"github.com/example/core-backend/pkg/kernel"
	"github.com/example/core-backend/pkg/migrate"
	"github.com/gofiber/fiber/v3"
	"github.com/jackc/pgx/v5/pgxpool"
)

//go:embed migrations/*.sql
var migrationsFS embed.FS

var _ kernel.Module = (*Module)(nil)

type Module struct {
	pool    *pgxpool.Pool
	bus     *eventbus.Bus
	service service.ITemplateService
	handler *handler.Handler
}

func New(pool *pgxpool.Pool, bus *eventbus.Bus) *Module {
	gormDB, _ := db.NewGorm(pool)
	repo := repository.NewPostgresTemplateRepository(gormDB)
	svc := service.NewTemplateService(repo, bus)
	h := handler.NewHandler(svc)

	return &Module{
		pool:    pool,
		bus:     bus,
		service: svc,
		handler: h,
	}
}

func (m *Module) Name() string {
	return "template"
}

func (m *Module) Register(r fiber.Router) error {
	group := r.Group("/templates")
	m.handler.RegisterRoutes(group)
	return nil
}

func (m *Module) Migrate(ctx context.Context) error {
	return migrate.Run(ctx, m.pool, m.Name(), migrationsFS, "migrations")
}
