package service

import (
	"context"
	"errors"
	"log/slog"
	"strings"

	"github.com/example/core-backend/modules/tenant/api"
	"github.com/example/core-backend/modules/tenant/internal/dto"
	"github.com/example/core-backend/modules/tenant/internal/errcode"
	"github.com/example/core-backend/modules/tenant/internal/model"
	"github.com/example/core-backend/pkg/constants/response_code"
	"github.com/example/core-backend/pkg/db"
	"github.com/example/core-backend/pkg/eventbus"
	"github.com/example/core-backend/pkg/exceptions"
	"github.com/example/core-backend/pkg/httpx"
	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
)

type Repository interface {
	Create(ctx context.Context, tenant *model.Tenant) error
	GetTenantByID(ctx context.Context, id string) (*model.Tenant, error)
	GetTenantByName(ctx context.Context, name string) (*model.Tenant, error)
	Update(ctx context.Context, tenant *model.Tenant) error
	Delete(ctx context.Context, id string) error
	GetAllTenant(ctx context.Context, req *dto.PaginationOptions) ([]model.Tenant, int64, error)
}

type ITenantService interface {
	Create(ctx context.Context, req *dto.CreateTenantRequest) (*model.Tenant, error)
	GetByID(ctx context.Context, id string) (*model.Tenant, error)
	GetAllTenant(ctx context.Context, req *dto.PaginationOptions) (*httpx.PaginatedData[model.Tenant], error)
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
		return nil, exceptions.NewBadRequestError(errcode.TenantNameRequired, nil)
	}

	name := strings.TrimSpace(req.Name)

	existing, err := t.tenantRepository.GetTenantByName(ctx, name)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		slog.ErrorContext(ctx, "failed to get tenant by name", "name", name, "error", err)
		return nil, exceptions.NewInternalServerError(errcode.FailToCreateTenant)
	}
	if existing != nil {
		return nil, exceptions.NewBadRequestError(errcode.TenantAlreadyExists, nil)
	}

	isolateMode := req.IsolateMode
	if isolateMode == "" {
		isolateMode = model.IsolateModeShared
	} else if !isolateMode.IsValid() {
		return nil, exceptions.NewBadRequestError(errcode.InvalidIsolateMode, nil)
	}

	var dbConfig *model.TenantDBConfig
	if isolateMode == model.IsolateModeDedicated {
		if req.DBConfig == nil || strings.TrimSpace(req.DBConfig.DSN) == "" {
			return nil, exceptions.NewBadRequestError(errcode.MissingDedicatedDBConfig, nil)
		}

		dsn := strings.TrimSpace(req.DBConfig.DSN)
		if err := db.ValidateDSN(dsn); err != nil {
			return nil, exceptions.NewBadRequestError(errcode.InvalidDedicatedDBConfig, nil)
		}

		dbConfig = &model.TenantDBConfig{
			DSN:      dsn,
			MaxConns: db.DefaultTenantMaxConns,
			SSLMode:  db.DefaultTenantSSLMode,
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
		slog.ErrorContext(ctx, "failed to create tenant", "name", name, "error", err)
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return nil, exceptions.NewBadRequestError(errcode.TenantAlreadyExists, nil)
		}
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return nil, exceptions.NewBadRequestError(errcode.TenantAlreadyExists, nil)
		}
		return nil, exceptions.NewInternalServerError(errcode.FailToCreateTenant)
	}

	if t.bus != nil {
		_ = eventbus.Publish(ctx, t.bus, api.TenantCreatedTopic, tenant)
	}

	return tenant, nil
}

func (t *tenantService) GetByID(ctx context.Context, id string) (*model.Tenant, error) {
	if strings.TrimSpace(id) == "" {
		return nil, exceptions.NewBadRequestError(response_code.InvalidParameter, nil)
	}

	tenant, err := t.tenantRepository.GetTenantByID(ctx, id)
	if err != nil {
		slog.ErrorContext(ctx, "failed to get tenant by id", "id", id, "error", err)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, exceptions.NewNotFoundError(errcode.NotFoundTenant)
		}
		return nil, exceptions.NewInternalServerError(errcode.FailToGetTenant)
	}
	if tenant == nil {
		return nil, exceptions.NewNotFoundError(errcode.NotFoundTenant)
	}

	return tenant, nil
}

func (t *tenantService) GetAllTenant(ctx context.Context, req *dto.PaginationOptions) (*httpx.PaginatedData[model.Tenant], error) {
	if req == nil {
		req = &dto.PaginationOptions{}
	}
	req.Normalize()

	tenants, total, err := t.tenantRepository.GetAllTenant(ctx, req)
	if err != nil {
		slog.ErrorContext(ctx, "failed to list tenants", "error", err)
		return nil, exceptions.NewInternalServerError(errcode.FailToGetTenant)
	}

	result := httpx.NewPaginatedData(tenants, total, req.Page, req.PageSize)
	return &result, nil
}
