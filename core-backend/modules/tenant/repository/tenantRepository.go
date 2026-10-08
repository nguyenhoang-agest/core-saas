package repository

import (
	"context"

	"github.com/example/core-backend/modules/tenant/domain"
)

type ITenantRepository interface {
	IRepository[domain.Tenant]
	GetTenantByID(ctx context.Context, id string) (*domain.Tenant, error)
	GetTenantByName(ctx context.Context, name string) (*domain.Tenant, error)
}
