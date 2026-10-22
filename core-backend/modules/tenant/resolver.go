package tenant

import (
	"context"

	"github.com/example/core-backend/modules/tenant/internal/service"
	"github.com/example/core-backend/pkg/db/router"
)

type tenantResolver struct {
	service service.ITenantService
}

// NewTenantResolver constructs a router.TenantResolver backed by tenant ITenantService.
func NewTenantResolver(svc service.ITenantService) router.TenantResolver {
	return &tenantResolver{service: svc}
}

// ResolveTenant maps tenant entity to router.TenantRouteInfo.
func (r *tenantResolver) ResolveTenant(ctx context.Context, tenantID string) (*router.TenantRouteInfo, error) {
	tenant, err := r.service.GetByID(ctx, tenantID)
	if err != nil {
		return nil, err
	}
	if tenant == nil {
		return nil, nil
	}

	info := &router.TenantRouteInfo{
		TenantID:    tenant.Id,
		IsolateMode: string(tenant.IsolateMode),
	}

	if tenant.DBConfig != nil {
		info.DSN = tenant.DBConfig.DSN
		info.MaxConns = tenant.DBConfig.MaxConns
	}

	return info, nil
}
