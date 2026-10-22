package user

import (
	"context"
	"embed"
	"log/slog"

	tenantapi "github.com/example/core-backend/modules/tenant/api"
	"github.com/example/core-backend/modules/user/internal/dto"
	"github.com/example/core-backend/modules/user/internal/handler"
	"github.com/example/core-backend/modules/user/internal/repository"
	"github.com/example/core-backend/modules/user/internal/service"
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
	handler     *handler.Handler
	roleHandler *handler.RoleHandler
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
				slog.InfoContext(ctx, "provisioning initial owner user for tenant", "tenant_id", evt.ID, "email", evt.OwnerEmail)
				ownerReq := &dto.CreateUserRequest{
					Email: evt.OwnerEmail,
					Name:  evt.OwnerName,
					Role:  string(rbac.RoleOwner),
				}
				tenantCtx := tenantctx.WithContext(ctx, evt.ID)
				if _, err := svc.Create(tenantCtx, ownerReq); err != nil {
					slog.ErrorContext(ctx, "failed to provision initial owner for tenant", "tenant_id", evt.ID, "error", err)
				} else {
					_ = service.BackfillLegacyUserRoles(ctx, gormDB, evt.ID)
				}
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

	return &Module{
		pool:        pool,
		bus:         bus,
		gormDB:      gormDB,
		service:     svc,
		roleService: roleSvc,
		handler:     h,
		roleHandler: rh,
	}
}

func (m *Module) Name() string {
	return "user"
}

func (m *Module) Register(r fiber.Router) error {
	// Nested routes under /users (register static role routes before user :id)
	group := r.Group("/users")
	m.roleHandler.RegisterRoutes(group)
	m.handler.RegisterRoutes(group)

	// Top-level aliases matching root API specification (/permissions and /roles)
	r.Get("/permissions", rbac.RequirePermission("iam:roles:read"), m.roleHandler.GetPermissionsCatalog)
	r.Get("/roles", rbac.RequirePermission("iam:roles:read"), m.roleHandler.ListRoles)
	r.Get("/roles/:id", rbac.RequirePermission("iam:roles:read"), m.roleHandler.GetRoleByID)
	r.Post("/roles", rbac.RequirePermission("iam:roles:manage"), m.roleHandler.CreateRole)
	r.Put("/roles/:id", rbac.RequirePermission("iam:roles:manage"), m.roleHandler.UpdateRole)
	r.Delete("/roles/:id", rbac.RequirePermission("iam:roles:manage"), m.roleHandler.DeleteRole)
	r.Post("/roles/:id/clone", rbac.RequirePermission("iam:roles:manage"), m.roleHandler.CloneRole)

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
