package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/example/core-backend/modules/tenant/internal/model"
	"github.com/example/core-backend/modules/tenant/internal/service"
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
