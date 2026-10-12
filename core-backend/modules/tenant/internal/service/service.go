package service

import (
	"context"
	"strings"

	"github.com/example/core-backend/modules/tenant/api"
	"github.com/example/core-backend/modules/tenant/internal/dto"
	"github.com/example/core-backend/modules/tenant/internal/model"
	"github.com/example/core-backend/pkg/eventbus"
)

type Repository interface {
	Create(ctx context.Context, tenant *model.Tenant) error
	GetTenantByID(ctx context.Context, id string) (*model.Tenant, error)
	GetTenantByName(ctx context.Context, name string) (*model.Tenant, error)
	Update(ctx context.Context, tenant *model.Tenant) error
	Delete(ctx context.Context, id string) error
}

type ITenantService interface {
	Create(ctx context.Context, req *dto.CreateTenantRequest) (*model.Tenant, error)
	GetByID(ctx context.Context, id string) (*model.Tenant, error)
}

type tenantService struct {
	tenantRepository Repository
	bus              *eventbus.Bus
}

func NewTenantService(tenantRepository Repository, bus *eventbus.Bus) ITenantService {
	return &tenantService{
		tenantRepository: tenantRepository,
		bus:              bus,
	}
}

func (t *tenantService) Create(ctx context.Context, req *dto.CreateTenantRequest) (*model.Tenant, error) {
	if req == nil || strings.TrimSpace(req.Name) == "" {
		return nil, model.ErrInvalidInput
	}

	name := strings.TrimSpace(req.Name)

	existing, err := t.tenantRepository.GetTenantByName(ctx, name)
	if err == nil && existing != nil {
		return nil, model.ErrAlreadyExists
	}

	isolateMode := req.IsolateMode
	if isolateMode == "" {
		isolateMode = model.IsolateModeShared
	} else if !isolateMode.IsValid() {
		return nil, model.ErrInvalidInput
	}

	var dbConfig *model.TenantDBConfig
	if isolateMode == model.IsolateModeDedicated {
		if req.DBConfig == nil || strings.TrimSpace(req.DBConfig.DSN) == "" {
			return nil, model.ErrInvalidInput
		}
		dbConfig = &model.TenantDBConfig{
			DSN:      strings.TrimSpace(req.DBConfig.DSN),
			MaxConns: req.DBConfig.MaxConns,
			SSLMode:  strings.TrimSpace(req.DBConfig.SSLMode),
		}
		if dbConfig.MaxConns <= 0 {
			dbConfig.MaxConns = 10
		}
	} else if req.DBConfig != nil {
		dbConfig = req.DBConfig
	}

	tenant := &model.Tenant{
		Name:        name,
		Description: strings.TrimSpace(req.Description),
		Status:      model.StatusActive,
		IsolateMode: isolateMode,
		DBConfig:    dbConfig,
	}

	if err := t.tenantRepository.Create(ctx, tenant); err != nil {
		return nil, err
	}

	if t.bus != nil {
		_ = eventbus.Publish(ctx, t.bus, api.TenantCreatedTopic, tenant)
	}

	return tenant, nil
}

func (t *tenantService) GetByID(ctx context.Context, id string) (*model.Tenant, error) {
	if strings.TrimSpace(id) == "" {
		return nil, model.ErrInvalidInput
	}

	tenant, err := t.tenantRepository.GetTenantByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if tenant == nil {
		return nil, model.ErrNotFound
	}

	return tenant, nil
}
