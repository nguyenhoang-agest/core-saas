package tenant

import (
	"context"
	"embed"
	"os"

	"github.com/example/core-backend/modules/tenant/internal/handler"
	"github.com/example/core-backend/modules/tenant/internal/repository"
	"github.com/example/core-backend/modules/tenant/internal/service"
	"github.com/example/core-backend/pkg/auth"
	"github.com/example/core-backend/pkg/db"
	"github.com/example/core-backend/pkg/db/router"
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
	pool      *pgxpool.Pool
	bus       *eventbus.Bus
	service   service.ITenantService
	handler   *handler.Handler
	jwtSecret string
}

func New(pool *pgxpool.Pool, bus *eventbus.Bus) *Module {
	gormDB, _ := db.NewGorm(pool)
	repo := repository.NewPostgresTenantRepository(gormDB)
	svc := service.NewTenantService(repo, bus)
	h := handler.NewHandler(svc)

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "dev-insecure-secret-key-change-in-production-12345"
	}

	return &Module{
		pool:      pool,
		bus:       bus,
		service:   svc,
		handler:   h,
		jwtSecret: jwtSecret,
	}
}

func (m *Module) Name() string {
	return "tenant"
}

func (m *Module) Register(r fiber.Router) error {
	group := r.Group("/tenants", auth.OptionalJWTAuthMiddleware(m.jwtSecret))
	m.handler.RegisterRoutes(group)
	return nil
}

func (m *Module) Migrate(ctx context.Context) error {
	return migrate.Run(ctx, m.pool, m.Name(), migrationsFS, "migrations")
}

// Resolver returns a router.TenantResolver backed by this module's tenant service.
func (m *Module) Resolver() router.TenantResolver {
	return NewTenantResolver(m.service)
}
