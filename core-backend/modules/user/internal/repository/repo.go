package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/example/core-backend/modules/user/internal/dto"
	"github.com/example/core-backend/modules/user/internal/model"
	"github.com/example/core-backend/pkg/db"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type IUserRepository interface {
	Create(ctx context.Context, user *model.User) error
	GetByID(ctx context.Context, tenantID, id string) (*model.User, error)
	GetByEmail(ctx context.Context, tenantID, email string) (*model.User, error)
	GetByEmailAnyTenant(ctx context.Context, email string) ([]model.User, error)
	Update(ctx context.Context, user *model.User) error
	UpdateLastLogin(ctx context.Context, tenantID, userID string, loginAt time.Time) error
	Delete(ctx context.Context, tenantID, id string) error
	GetAll(ctx context.Context, tenantID string, req *dto.PaginationOptions) ([]model.User, int64, error)
}

var _ IUserRepository = (*PostgresUserRepository)(nil)

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
	if user.Role == "" {
		user.Role = "member"
	}
	if user.Permissions == nil {
		user.Permissions = model.StringArray{}
	}

	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(user).Error; err != nil {
			return err
		}

		var role model.Role
		if err := tx.Where("tenant_id = ? AND slug = ? AND delete_at IS NULL", user.TenantID, user.Role).First(&role).Error; err == nil {
			_ = tx.Create(&model.UserRole{
				TenantID:   user.TenantID,
				UserID:     user.Id,
				RoleID:     role.ID,
				AssignedBy: "system",
				AssignedAt: time.Now(),
			}).Error
		}

		return nil
	})
}

func (r *PostgresUserRepository) GetByID(ctx context.Context, tenantID, id string) (*model.User, error) {
	if r.db == nil {
		return nil, fmt.Errorf("database connection is not initialized")
	}

	var u model.User
	if err := r.db.WithContext(ctx).
		Where("tenant_id = ? AND id = ? AND delete_at IS NULL", tenantID, id).
		First(&u).Error; err != nil {
		return nil, err
	}

	return &u, nil
}

func (r *PostgresUserRepository) GetByEmail(ctx context.Context, tenantID, email string) (*model.User, error) {
	if r.db == nil {
		return nil, fmt.Errorf("database connection is not initialized")
	}

	var u model.User
	if err := r.db.WithContext(ctx).
		Where("tenant_id = ? AND email = ? AND delete_at IS NULL", tenantID, email).
		First(&u).Error; err != nil {
		return nil, err
	}

	return &u, nil
}

func (r *PostgresUserRepository) GetByEmailAnyTenant(ctx context.Context, email string) ([]model.User, error) {
	if r.db == nil {
		return nil, fmt.Errorf("database connection is not initialized")
	}

	var users []model.User
	if err := r.db.WithContext(ctx).
		Where("email = ? AND delete_at IS NULL", email).
		Find(&users).Error; err != nil {
		return nil, err
	}

	return users, nil
}

func (r *PostgresUserRepository) UpdateLastLogin(ctx context.Context, tenantID, userID string, loginAt time.Time) error {
	if r.db == nil {
		return fmt.Errorf("database connection is not initialized")
	}

	return r.db.WithContext(ctx).
		Model(&model.User{}).
		Where("tenant_id = ? AND id = ? AND delete_at IS NULL", tenantID, userID).
		Update("last_login_at", loginAt).Error
}

func (r *PostgresUserRepository) Update(ctx context.Context, user *model.User) error {
	if r.db == nil {
		return fmt.Errorf("database connection is not initialized")
	}

	now := time.Now().UTC()
	user.UpdateAt = now

	result := r.db.WithContext(ctx).
		Model(&model.User{}).
		Where("tenant_id = ? AND id = ? AND delete_at IS NULL", user.TenantID, user.Id).
		Select("name", "role", "permissions", "status", "phone", "avatar_url", "update_at").
		Updates(user)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

func (r *PostgresUserRepository) Delete(ctx context.Context, tenantID, id string) error {
	if r.db == nil {
		return fmt.Errorf("database connection is not initialized")
	}

	now := time.Now().UTC()
	result := r.db.WithContext(ctx).
		Model(&model.User{}).
		Where("tenant_id = ? AND id = ? AND delete_at IS NULL", tenantID, id).
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

func (r *PostgresUserRepository) GetAll(ctx context.Context, tenantID string, req *dto.PaginationOptions) ([]model.User, int64, error) {
	if r.db == nil {
		return nil, 0, fmt.Errorf("database connection is not initialized")
	}

	tx := r.db.WithContext(ctx).
		Model(&model.User{}).
		Where("tenant_id = ? AND delete_at IS NULL", tenantID)

	if req != nil {
		if req.Keyword != "" {
			pattern := "%" + req.Keyword + "%"
			tx = tx.Where("(name ILIKE ? OR email ILIKE ?)", pattern, pattern)
		}
		if req.Role != "" {
			tx = tx.Where("role = ?", req.Role)
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
		switch req.SortBy {
		case "name", "email", "role", "status", "create_at":
			sortBy = req.SortBy
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
