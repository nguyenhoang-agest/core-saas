package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/example/core-backend/modules/user/internal/dto"
	"github.com/example/core-backend/modules/user/internal/model"
	"github.com/example/core-backend/modules/user/internal/service"
	"github.com/example/core-backend/pkg/db"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

var _ service.Repository = (*PostgresUserRepository)(nil)

type PostgresUserRepository struct {
	db *gorm.DB
}

func NewPostgresUserRepository(db *gorm.DB) *PostgresUserRepository {
	return &PostgresUserRepository{db: db}
}

func (r *PostgresUserRepository) Create(ctx context.Context, user *model.User) error {
	if r.db == nil {
		return fmt.Errorf("database connection is not initialized")
	}

	if user.Id == "" {
		user.Id = uuid.NewString()
	}
	if user.Status == "" {
		user.Status = model.StatusActive
	}

	return r.db.WithContext(ctx).Create(user).Error
}

func (r *PostgresUserRepository) GetByID(ctx context.Context, id string) (*model.User, error) {
	if r.db == nil {
		return nil, fmt.Errorf("database connection is not initialized")
	}

	var t model.User
	if err := r.db.WithContext(ctx).
		Where("id = ? AND delete_at IS NULL", id).
		First(&t).Error; err != nil {
		return nil, err
	}

	return &t, nil
}

func (r *PostgresUserRepository) GetByName(ctx context.Context, name string) (*model.User, error) {
	if r.db == nil {
		return nil, fmt.Errorf("database connection is not initialized")
	}

	var t model.User
	if err := r.db.WithContext(ctx).
		Where("name = ? AND delete_at IS NULL", name).
		First(&t).Error; err != nil {
		return nil, err
	}

	return &t, nil
}

func (r *PostgresUserRepository) Update(ctx context.Context, user *model.User) error {
	if r.db == nil {
		return fmt.Errorf("database connection is not initialized")
	}

	result := r.db.WithContext(ctx).
		Model(&model.User{}).
		Where("id = ? AND delete_at IS NULL", user.Id).
		Select("name", "description", "status", "update_at").
		Updates(user)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

func (r *PostgresUserRepository) Delete(ctx context.Context, id string) error {
	if r.db == nil {
		return fmt.Errorf("database connection is not initialized")
	}

	now := time.Now().UTC()
	result := r.db.WithContext(ctx).
		Model(&model.User{}).
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

func (r *PostgresUserRepository) GetAll(ctx context.Context, req *dto.PaginationOptions) ([]model.User, int64, error) {
	if r.db == nil {
		return nil, 0, fmt.Errorf("database connection is not initialized")
	}

	tx := r.db.WithContext(ctx).Model(&model.User{}).Where("delete_at IS NULL")

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

	var items []model.User
	if err := tx.Scopes(db.PaginateOffset(offset, limit)).Find(&items).Error; err != nil {
		return nil, 0, err
	}

	return items, total, nil
}
