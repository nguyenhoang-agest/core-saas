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
	Update(ctx context.Context, id string, req *dto.UpdateTenantRequest) (*model.Tenant, error)
	UpdateStatus(ctx context.Context, id string, status model.Status) (*model.Tenant, error)
	Delete(ctx context.Context, id string) error
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
		evt := api.TenantCreatedEvent{
			ID:          tenant.Id,
			Name:        tenant.Name,
			IsolateMode: string(tenant.IsolateMode),
			IsHost:      tenant.IsHost,
		}
		if req.Owner != nil {
			evt.OwnerEmail = strings.TrimSpace(req.Owner.Email)
			evt.OwnerName = strings.TrimSpace(req.Owner.Name)
			evt.OwnerPassword = req.Owner.Password
		}
		_ = eventbus.Publish(ctx, t.bus, api.TenantCreatedTopic, evt)
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

func (t *tenantService) Update(ctx context.Context, id string, req *dto.UpdateTenantRequest) (*model.Tenant, error) {
	if strings.TrimSpace(id) == "" || req == nil || strings.TrimSpace(req.Name) == "" {
		return nil, exceptions.NewBadRequestError(response_code.InvalidParameter, nil)
	}

	tenant, err := t.tenantRepository.GetTenantByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, exceptions.NewNotFoundError(errcode.NotFoundTenant)
		}
		return nil, exceptions.NewInternalServerError(errcode.FailToGetTenant)
	}

	newName := strings.TrimSpace(req.Name)
	if newName != tenant.Name {
		existing, err := t.tenantRepository.GetTenantByName(ctx, newName)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, exceptions.NewInternalServerError(errcode.FailToGetTenant)
		}
		if existing != nil && existing.Id != id {
			return nil, exceptions.NewBadRequestError(errcode.TenantAlreadyExists, nil)
		}
	}

	tenant.Name = newName
	tenant.Description = strings.TrimSpace(req.Description)

	if err := t.tenantRepository.Update(ctx, tenant); err != nil {
		slog.ErrorContext(ctx, "failed to update tenant", "id", id, "error", err)
		return nil, exceptions.NewInternalServerError(errcode.FailToCreateTenant)
	}

	if t.bus != nil {
		_ = eventbus.Publish(ctx, t.bus, api.TenantUpdatedTopic, api.TenantUpdatedEvent{
			ID:          tenant.Id,
			Name:        tenant.Name,
			Status:      string(tenant.Status),
			IsolateMode: string(tenant.IsolateMode),
		})
	}

	return tenant, nil
}

func (t *tenantService) UpdateStatus(ctx context.Context, id string, status model.Status) (*model.Tenant, error) {
	if strings.TrimSpace(id) == "" || (status != model.StatusActive && status != model.StatusSuspended) {
		return nil, exceptions.NewBadRequestError(response_code.InvalidParameter, nil)
	}

	tenant, err := t.tenantRepository.GetTenantByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, exceptions.NewNotFoundError(errcode.NotFoundTenant)
		}
		return nil, exceptions.NewInternalServerError(errcode.FailToGetTenant)
	}

	tenant.Status = status
	if err := t.tenantRepository.Update(ctx, tenant); err != nil {
		slog.ErrorContext(ctx, "failed to update tenant status", "id", id, "status", status, "error", err)
		return nil, exceptions.NewInternalServerError(errcode.FailToGetTenant)
	}

	if t.bus != nil {
		_ = eventbus.Publish(ctx, t.bus, api.TenantUpdatedTopic, api.TenantUpdatedEvent{
			ID:          tenant.Id,
			Name:        tenant.Name,
			Status:      string(tenant.Status),
			IsolateMode: string(tenant.IsolateMode),
		})
		if status == model.StatusSuspended {
			_ = eventbus.Publish(ctx, t.bus, api.TenantSuspendedTopic, api.TenantStatusEvent{
				ID:     tenant.Id,
				Status: string(status),
			})
		}
	}

	return tenant, nil
}

func (t *tenantService) Delete(ctx context.Context, id string) error {
	if strings.TrimSpace(id) == "" {
		return exceptions.NewBadRequestError(response_code.InvalidParameter, nil)
	}

	tenant, err := t.tenantRepository.GetTenantByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return exceptions.NewNotFoundError(errcode.NotFoundTenant)
		}
		return exceptions.NewInternalServerError(errcode.FailToGetTenant)
	}

	if err := t.tenantRepository.Delete(ctx, id); err != nil {
		slog.ErrorContext(ctx, "failed to delete tenant", "id", id, "error", err)
		return exceptions.NewInternalServerError(errcode.FailToGetTenant)
	}

	if t.bus != nil {
		_ = eventbus.Publish(ctx, t.bus, api.TenantDeletedTopic, api.TenantStatusEvent{
			ID:     tenant.Id,
			Status: "deleted",
		})
	}

	return nil
}
