package user

import (
	"context"
	"embed"

	"github.com/example/core-backend/modules/user/internal/handler"
	"github.com/example/core-backend/modules/user/internal/repository"
	"github.com/example/core-backend/modules/user/internal/service"
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
	service service.IUserService
	handler *handler.Handler
}

func New(pool *pgxpool.Pool, bus *eventbus.Bus) *Module {
	gormDB, _ := db.NewGorm(pool)
	repo := repository.NewPostgresUserRepository(gormDB)
	svc := service.NewUserService(repo, bus)
	h := handler.NewHandler(svc)

	return &Module{
		pool:    pool,
		bus:     bus,
		service: svc,
		handler: h,
	}
}

func (m *Module) Name() string {
	return "user"
}

func (m *Module) Register(r fiber.Router) error {
	group := r.Group("/users")
	m.handler.RegisterRoutes(group)
	return nil
}

func (m *Module) Migrate(ctx context.Context) error {
	return migrate.Run(ctx, m.pool, m.Name(), migrationsFS, "migrations")
}
