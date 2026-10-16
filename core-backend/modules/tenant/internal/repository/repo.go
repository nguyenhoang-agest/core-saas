package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/example/core-backend/modules/tenant/internal/dto"
	"github.com/example/core-backend/modules/tenant/internal/model"
	"github.com/example/core-backend/modules/tenant/internal/service"
	"github.com/example/core-backend/pkg/db"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

var _ service.Repository = (*PostgresTenantRepository)(nil)

type PostgresTenantRepository struct {
	db *gorm.DB
}

func NewPostgresTenantRepository(db *gorm.DB) *PostgresTenantRepository {
	return &PostgresTenantRepository{db: db}
}

func (r *PostgresTenantRepository) Create(ctx context.Context, tenant *model.Tenant) error {
	if r.db == nil {
		return fmt.Errorf("database connection is not initialized")
	}

	if tenant.Id == "" {
		tenant.Id = uuid.NewString()
	}
	if tenant.IsolateMode == "" {
		tenant.IsolateMode = model.IsolateModeShared
	}

	return r.db.WithContext(ctx).Create(tenant).Error
}

func (r *PostgresTenantRepository) GetTenantByID(ctx context.Context, id string) (*model.Tenant, error) {
	if r.db == nil {
		return nil, fmt.Errorf("database connection is not initialized")
	}

	var t model.Tenant
	if err := r.db.WithContext(ctx).
		Where("id = ? AND delete_at IS NULL", id).
		First(&t).Error; err != nil {
		return nil, err
	}

	if t.IsolateMode == "" {
		t.IsolateMode = model.IsolateModeShared
	}

	return &t, nil
}

func (r *PostgresTenantRepository) GetTenantByName(ctx context.Context, name string) (*model.Tenant, error) {
	if r.db == nil {
		return nil, fmt.Errorf("database connection is not initialized")
	}

	var t model.Tenant
	if err := r.db.WithContext(ctx).
		Where("name = ? AND delete_at IS NULL", name).
		First(&t).Error; err != nil {
		return nil, err
	}

	if t.IsolateMode == "" {
		t.IsolateMode = model.IsolateModeShared
	}

	return &t, nil
}

func (r *PostgresTenantRepository) Update(ctx context.Context, tenant *model.Tenant) error {
	if r.db == nil {
		return fmt.Errorf("database connection is not initialized")
	}

	result := r.db.WithContext(ctx).
		Model(&model.Tenant{}).
		Where("id = ? AND delete_at IS NULL", tenant.Id).
		Select("name", "description", "status", "isolate_mode", "db_config", "update_at").
		Updates(tenant)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

func (r *PostgresTenantRepository) Delete(ctx context.Context, id string) error {
	if r.db == nil {
		return fmt.Errorf("database connection is not initialized")
	}

	now := time.Now().UTC()
	result := r.db.WithContext(ctx).
		Model(&model.Tenant{}).
		Where("id = ? AND delete_at IS NULL", id).
		Updates(map[string]any{
			"delete_at": now,
			"update_at": now,
		})
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

func (r *PostgresTenantRepository) GetAllTenant(ctx context.Context, req *dto.PaginationOptions) ([]model.Tenant, int64, error) {
	if r.db == nil {
		return nil, 0, fmt.Errorf("database connection is not initialized")
	}

	tx := r.db.WithContext(ctx).Model(&model.Tenant{}).Where("delete_at IS NULL")

	if req != nil {
		if req.Keyword != "" {
			tx = tx.Where("name ILIKE ?", "%"+req.Keyword+"%")
		}
		if req.Status != "" {
			tx = tx.Where("status = ?", req.Status)
		}
		if req.IsolateMode != "" {
			tx = tx.Where("isolate_mode = ?", req.IsolateMode)
		}
	}

	var total int64
	if err := tx.Session(&gorm.Session{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	sortBy := "create_at"
	order := "desc"
	offset := 0
	limit := 20

	if req != nil {
		if req.SortBy == "name" {
			sortBy = "name"
		}
		if req.Order == "asc" || req.Order == "desc" {
			order = req.Order
		}
		offset = req.Offset()
		limit = req.Limit()
	}

	tx = tx.Order(fmt.Sprintf("%s %s", sortBy, order))

	var tenants []model.Tenant
	if err := tx.Scopes(db.PaginateOffset(offset, limit)).Find(&tenants).Error; err != nil {
		return nil, 0, err
	}

	for i := range tenants {
		if tenants[i].IsolateMode == "" {
			tenants[i].IsolateMode = model.IsolateModeShared
		}
	}

	return tenants, total, nil
}
