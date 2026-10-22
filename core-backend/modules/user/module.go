package user

import (
	"context"
	"embed"
	"log/slog"
	"os"
	"strings"
	"time"

	tenantapi "github.com/example/core-backend/modules/tenant/api"
	"github.com/example/core-backend/modules/user/internal/dto"
	"github.com/example/core-backend/modules/user/internal/handler"
	"github.com/example/core-backend/modules/user/internal/repository"
	"github.com/example/core-backend/modules/user/internal/service"
	"github.com/example/core-backend/pkg/auth"
	"github.com/example/core-backend/pkg/db"
	"github.com/example/core-backend/pkg/eventbus"
	"github.com/example/core-backend/pkg/kernel"
	"github.com/example/core-backend/pkg/migrate"
	"github.com/example/core-backend/pkg/rbac"
	"github.com/example/core-backend/pkg/tenantctx"
	"github.com/gofiber/fiber/v3"
	"github.com/jackc/pgx/v5/pgxpool"
	"gorm.io/gorm"
)

//go:embed migrations/*.sql
var migrationsFS embed.FS

var _ kernel.Module = (*Module)(nil)

type Module struct {
	pool        *pgxpool.Pool
	bus         *eventbus.Bus
	gormDB      *gorm.DB
	service     service.IUserService
	roleService service.IRoleService
	authService service.IAuthService
	handler     *handler.Handler
	roleHandler *handler.RoleHandler
	authHandler *handler.AuthHandler
	jwtSecret   string
}

func New(pool *pgxpool.Pool, bus *eventbus.Bus) *Module {
	gormDB, _ := db.NewGorm(pool)
	repo := repository.NewPostgresUserRepository(gormDB)
	roleRepo := repository.NewPostgresRoleRepository(gormDB)

	svc := service.NewUserService(repo, bus)
	roleSvc := service.NewRoleService(roleRepo, bus)

	h := handler.NewHandler(svc)
	rh := handler.NewRoleHandler(roleSvc)

	var resolver *service.CachedPermissionResolver
	if gormDB != nil {
		resolver = service.NewCachedPermissionResolver(gormDB, bus)
		rbac.SetPermissionResolver(resolver)
	}

	if bus != nil && gormDB != nil {
		eventbus.Subscribe[tenantapi.TenantCreatedEvent](bus, tenantapi.TenantCreatedTopic, func(ctx context.Context, evt tenantapi.TenantCreatedEvent) error {
			slog.InfoContext(ctx, "seeding system roles for new tenant", "tenant_id", evt.ID, "is_host", evt.IsHost)
			if err := service.SyncSystemRoles(ctx, gormDB, evt.ID, evt.IsHost); err != nil {
				slog.ErrorContext(ctx, "failed to seed system roles for tenant", "tenant_id", evt.ID, "error", err)
				return err
			}

			if evt.OwnerEmail != "" {
				cleanEmail := strings.ToLower(strings.TrimSpace(evt.OwnerEmail))
				existing, _ := repo.GetByEmail(ctx, evt.ID, cleanEmail)
				if existing == nil {
					slog.InfoContext(ctx, "provisioning initial owner user for tenant", "tenant_id", evt.ID, "email", cleanEmail)
					ownerReq := &dto.CreateUserRequest{
						Email:    cleanEmail,
						Name:     strings.TrimSpace(evt.OwnerName),
						Role:     string(rbac.RoleOwner),
						Password: evt.OwnerPassword,
					}
					tenantCtx := tenantctx.WithContext(ctx, evt.ID)
					if _, err := svc.Create(tenantCtx, ownerReq); err != nil {
						slog.ErrorContext(ctx, "failed to provision initial owner for tenant", "tenant_id", evt.ID, "error", err)
					}
				}
				_ = service.BackfillLegacyUserRoles(ctx, gormDB, evt.ID)
			}
			return nil
		})

		eventbus.Subscribe[tenantapi.TenantStatusEvent](bus, tenantapi.TenantSuspendedTopic, func(ctx context.Context, evt tenantapi.TenantStatusEvent) error {
			if resolver != nil {
				resolver.Invalidate(ctx, evt.ID, "")
			}
			return nil
		})
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "dev-insecure-secret-key-change-in-production-12345"
	}

	authSvc := service.NewAuthService(gormDB, repo, bus, resolver, jwtSecret, 15*time.Minute, 7*24*time.Hour)
	authH := handler.NewAuthHandler(authSvc, jwtSecret)

	return &Module{
		pool:        pool,
		bus:         bus,
		gormDB:      gormDB,
		service:     svc,
		roleService: roleSvc,
		authService: authSvc,
		handler:     h,
		roleHandler: rh,
		authHandler: authH,
		jwtSecret:   jwtSecret,
	}
}

func (m *Module) Name() string {
	return "user"
}

func (m *Module) Register(r fiber.Router) error {
	// Top-level auth routes (/auth/register, /auth/login, /auth/refresh, /auth/me, /auth/logout)
	m.authHandler.RegisterRoutes(r)

	// Nested routes under /users with JWT context resolution
	group := r.Group("/users", auth.OptionalJWTAuthMiddleware(m.jwtSecret))
	m.roleHandler.RegisterRoutes(group)
	m.handler.RegisterRoutes(group)

	// Top-level aliases matching root API specification (/permissions and /roles)
	rolesGroup := r.Group("", auth.OptionalJWTAuthMiddleware(m.jwtSecret))
	rolesGroup.Get("/permissions", rbac.RequirePermission("iam:roles:read"), m.roleHandler.GetPermissionsCatalog)
	rolesGroup.Get("/roles", rbac.RequirePermission("iam:roles:read"), m.roleHandler.ListRoles)
	rolesGroup.Get("/roles/:id", rbac.RequirePermission("iam:roles:read"), m.roleHandler.GetRoleByID)
	rolesGroup.Post("/roles", rbac.RequirePermission("iam:roles:manage"), m.roleHandler.CreateRole)
	rolesGroup.Put("/roles/:id", rbac.RequirePermission("iam:roles:manage"), m.roleHandler.UpdateRole)
	rolesGroup.Delete("/roles/:id", rbac.RequirePermission("iam:roles:manage"), m.roleHandler.DeleteRole)
	rolesGroup.Post("/roles/:id/clone", rbac.RequirePermission("iam:roles:manage"), m.roleHandler.CloneRole)

	return nil
}

func (m *Module) Migrate(ctx context.Context) error {
	if err := migrate.Run(ctx, m.pool, m.Name(), migrationsFS, "migrations"); err != nil {
		return err
	}
	if m.gormDB != nil {
		_ = service.SyncSystemRoles(ctx, m.gormDB, "default", true)
		_ = service.BackfillLegacyUserRoles(ctx, m.gormDB, "default")
	}
	return nil
}
