package service

import (
	"context"
	"strings"

	"github.com/example/core-backend/modules/tenant/domain"
	"github.com/example/core-backend/modules/tenant/domain/dto"
	"github.com/example/core-backend/modules/tenant/domain/exception"
	"github.com/example/core-backend/modules/tenant/repository"
	"github.com/example/core-backend/pkg/eventbus"
)

const (
	TenantCreatedTopic = "tenant.created"
	TenantUpdatedTopic = "tenant.updated"
)

type ITenantService interface {
	Create(ctx context.Context, req *dto.CreateTenantRequest) (*domain.Tenant, error)
	GetByID(ctx context.Context, id string) (*domain.Tenant, error)
}

type tenantService struct {
	tenantRepository repository.ITenantRepository
	bus              *eventbus.Bus
}

func NewTenantService(tenantRepository repository.ITenantRepository, bus *eventbus.Bus) ITenantService {
	return &tenantService{
		tenantRepository: tenantRepository,
		bus:              bus,
	}
}

func (t *tenantService) Create(ctx context.Context, req *dto.CreateTenantRequest) (*domain.Tenant, error) {
	if req == nil || strings.TrimSpace(req.Name) == "" {
		return nil, exception.ErrInvalidInput
	}

	name := strings.TrimSpace(req.Name)

	existing, err := t.tenantRepository.GetTenantByName(ctx, name)
	if err == nil && existing != nil {
		return nil, exception.ErrAlreadyExists
	}

	tenant := &domain.Tenant{
		Name:        name,
		Description: strings.TrimSpace(req.Description),
		Status:      domain.StatusActive,
	}

	if err := t.tenantRepository.Create(ctx, tenant); err != nil {
		return nil, err
	}

	if t.bus != nil {
		_ = eventbus.Publish(ctx, t.bus, TenantCreatedTopic, tenant)
	}

	return tenant, nil
}

func (t *tenantService) GetByID(ctx context.Context, id string) (*domain.Tenant, error) {
	if strings.TrimSpace(id) == "" {
		return nil, exception.ErrInvalidInput
	}

	tenant, err := t.tenantRepository.GetTenantByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if tenant == nil {
		return nil, exception.ErrNotFound
	}

	return tenant, nil
}
