package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/example/core-backend/modules/_template/internal/dto"
	"github.com/example/core-backend/modules/_template/internal/model"
	"github.com/example/core-backend/modules/_template/internal/service"
	"github.com/example/core-backend/pkg/db"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

var _ service.Repository = (*PostgresTemplateRepository)(nil)

type PostgresTemplateRepository struct {
	db *gorm.DB
}

func NewPostgresTemplateRepository(db *gorm.DB) *PostgresTemplateRepository {
	return &PostgresTemplateRepository{db: db}
}

func (r *PostgresTemplateRepository) Create(ctx context.Context, template *model.Template) error {
	if r.db == nil {
		return fmt.Errorf("database connection is not initialized")
	}

	if template.Id == "" {
		template.Id = uuid.NewString()
	}
	if template.Status == "" {
		template.Status = model.StatusActive
	}

	return r.db.WithContext(ctx).Create(template).Error
}

func (r *PostgresTemplateRepository) GetByID(ctx context.Context, id string) (*model.Template, error) {
	if r.db == nil {
		return nil, fmt.Errorf("database connection is not initialized")
	}

	var t model.Template
	if err := r.db.WithContext(ctx).
		Where("id = ? AND delete_at IS NULL", id).
		First(&t).Error; err != nil {
		return nil, err
	}

	return &t, nil
}

func (r *PostgresTemplateRepository) GetByName(ctx context.Context, name string) (*model.Template, error) {
	if r.db == nil {
		return nil, fmt.Errorf("database connection is not initialized")
	}

	var t model.Template
	if err := r.db.WithContext(ctx).
		Where("name = ? AND delete_at IS NULL", name).
		First(&t).Error; err != nil {
		return nil, err
	}

	return &t, nil
}

func (r *PostgresTemplateRepository) Update(ctx context.Context, template *model.Template) error {
	if r.db == nil {
		return fmt.Errorf("database connection is not initialized")
	}

	result := r.db.WithContext(ctx).
		Model(&model.Template{}).
		Where("id = ? AND delete_at IS NULL", template.Id).
		Select("name", "description", "status", "update_at").
		Updates(template)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

func (r *PostgresTemplateRepository) Delete(ctx context.Context, id string) error {
	if r.db == nil {
		return fmt.Errorf("database connection is not initialized")
	}

	now := time.Now().UTC()
	result := r.db.WithContext(ctx).
		Model(&model.Template{}).
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

func (r *PostgresTemplateRepository) GetAll(ctx context.Context, req *dto.PaginationOptions) ([]model.Template, int64, error) {
	if r.db == nil {
		return nil, 0, fmt.Errorf("database connection is not initialized")
	}

	tx := r.db.WithContext(ctx).Model(&model.Template{}).Where("delete_at IS NULL")

	if req != nil {
		if req.Keyword != "" {
			tx = tx.Where("name ILIKE ?", "%"+req.Keyword+"%")
		}
		if req.Status != "" {
			tx = tx.Where("status = ?", req.Status)
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

	var items []model.Template
	if err := tx.Scopes(db.PaginateOffset(offset, limit)).Find(&items).Error; err != nil {
		return nil, 0, err
	}

	return items, total, nil
}
