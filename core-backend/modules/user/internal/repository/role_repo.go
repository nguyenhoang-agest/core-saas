package repository

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/example/core-backend/modules/user/internal/errcode"
	"github.com/example/core-backend/modules/user/internal/model"
	"github.com/example/core-backend/pkg/exceptions"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type IRoleRepository interface {
	ListRoles(ctx context.Context, tenantID string) ([]model.Role, error)
	GetRoleByID(ctx context.Context, tenantID, roleID string) (*model.Role, error)
	GetRoleBySlug(ctx context.Context, tenantID, slug string) (*model.Role, error)
	CreateRole(ctx context.Context, role *model.Role, permissionKeys []string) error
	UpdateRole(ctx context.Context, tenantID, roleID string, name, description string, permissionKeys []string, expectedVersion int) (*model.Role, error)
	DeleteRole(ctx context.Context, tenantID, roleID string) (int, error)
	CloneRole(ctx context.Context, tenantID, sourceRoleID, newName, newSlug, newDescription string) (*model.Role, error)
	GetUserRoles(ctx context.Context, tenantID, userID string) ([]model.Role, error)
	AssignUserRoles(ctx context.Context, tenantID, userID string, roleIDs []string, assignedBy string, isActorOwner bool) error
	BumpAuthzVersion(ctx context.Context, tenantID string) (int64, error)
	CountActiveOwners(ctx context.Context, tenantID string) (int64, error)
}

type PostgresRoleRepository struct {
	db *gorm.DB
}

func NewPostgresRoleRepository(db *gorm.DB) *PostgresRoleRepository {
	return &PostgresRoleRepository{db: db}
}

func (r *PostgresRoleRepository) ListRoles(ctx context.Context, tenantID string) ([]model.Role, error) {
	var roles []model.Role
	err := r.db.WithContext(ctx).
		Preload("Permissions").
		Where("tenant_id = ? AND delete_at IS NULL", tenantID).
		Order("is_system DESC, name ASC").
		Find(&roles).Error
	if err != nil {
		return nil, err
	}

	// Aggregate active user counts per role
	type roleCount struct {
		RoleID string
		Count  int
	}
	var counts []roleCount
	_ = r.db.WithContext(ctx).
		Table("user_roles").
		Select("role_id, count(*) as count").
		Where("tenant_id = ?", tenantID).
		Group("role_id").
		Scan(&counts).Error

	countMap := make(map[string]int)
	for _, c := range counts {
		countMap[c.RoleID] = c.Count
	}

	for i := range roles {
		roles[i].UserCount = countMap[roles[i].ID]
	}

	return roles, nil
}

func (r *PostgresRoleRepository) GetRoleByID(ctx context.Context, tenantID, roleID string) (*model.Role, error) {
	var role model.Role
	err := r.db.WithContext(ctx).
		Preload("Permissions").
		Where("tenant_id = ? AND id = ? AND delete_at IS NULL", tenantID, roleID).
		First(&role).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, exceptions.NewNotFoundError(errcode.NotFoundRole)
		}
		return nil, err
	}

	var count int64
	_ = r.db.WithContext(ctx).
		Table("user_roles").
		Where("tenant_id = ? AND role_id = ?", tenantID, roleID).
		Count(&count).Error
	role.UserCount = int(count)

	return &role, nil
}

func (r *PostgresRoleRepository) GetRoleBySlug(ctx context.Context, tenantID, slug string) (*model.Role, error) {
	var role model.Role
	err := r.db.WithContext(ctx).
		Preload("Permissions").
		Where("tenant_id = ? AND slug = ? AND delete_at IS NULL", tenantID, slug).
		First(&role).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, exceptions.NewNotFoundError(errcode.NotFoundRole)
		}
		return nil, err
	}
	return &role, nil
}

func (r *PostgresRoleRepository) CreateRole(ctx context.Context, role *model.Role, permissionKeys []string) error {
	if role.ID == "" {
		role.ID = uuid.NewString()
	}
	role.Version = 1

	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Check slug uniqueness
		var existing model.Role
		err := tx.Where("tenant_id = ? AND slug = ? AND delete_at IS NULL", role.TenantID, role.Slug).First(&existing).Error
		if err == nil {
			return exceptions.NewBadRequestError(errcode.RoleSlugAlreadyExists, nil)
		} else if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}

		if err := tx.Create(role).Error; err != nil {
			return err
		}

		if len(permissionKeys) > 0 {
			perms := make([]model.RolePermission, 0, len(permissionKeys))
			for _, k := range permissionKeys {
				perms = append(perms, model.RolePermission{
					TenantID:      role.TenantID,
					RoleID:        role.ID,
					PermissionKey: k,
				})
			}
			if err := tx.Create(&perms).Error; err != nil {
				return err
			}
		}

		return r.bumpAuthzVersionTx(tx, role.TenantID)
	})
}

func (r *PostgresRoleRepository) UpdateRole(ctx context.Context, tenantID, roleID string, name, description string, permissionKeys []string, expectedVersion int) (*model.Role, error) {
	var updatedRole model.Role

	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Pessimistic lock on role row
		var existing model.Role
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("tenant_id = ? AND id = ? AND delete_at IS NULL", tenantID, roleID).
			First(&existing).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return exceptions.NewNotFoundError(errcode.NotFoundRole)
			}
			return err
		}

		// Protect immutable system roles
		if existing.IsSystem {
			return exceptions.NewBadRequestError(errcode.CannotModifySystemRole, nil)
		}

		// Optimistic concurrency verification
		if existing.Version != expectedVersion {
			return exceptions.NewBadRequestError(errcode.ConcurrencyConflict, nil)
		}

		existing.Name = strings.TrimSpace(name)
		existing.Description = strings.TrimSpace(description)
		existing.Version++
		existing.UpdateAt = time.Now()

		if err := tx.Save(&existing).Error; err != nil {
			return err
		}

		// Synchronize permissions: clear old and insert updated concrete keys
		if err := tx.Where("role_id = ?", roleID).Delete(&model.RolePermission{}).Error; err != nil {
			return err
		}

		if len(permissionKeys) > 0 {
			perms := make([]model.RolePermission, 0, len(permissionKeys))
			for _, k := range permissionKeys {
				perms = append(perms, model.RolePermission{
					TenantID:      tenantID,
					RoleID:        roleID,
					PermissionKey: k,
				})
			}
			if err := tx.Create(&perms).Error; err != nil {
				return err
			}
		}

		if err := r.bumpAuthzVersionTx(tx, tenantID); err != nil {
			return err
		}

		updatedRole = existing
		return nil
	})

	if err != nil {
		return nil, err
	}
	return &updatedRole, nil
}

func (r *PostgresRoleRepository) DeleteRole(ctx context.Context, tenantID, roleID string) (int, error) {
	var affectedUsers int

	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Lock role row for update
		var existing model.Role
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("tenant_id = ? AND id = ? AND delete_at IS NULL", tenantID, roleID).
			First(&existing).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return exceptions.NewNotFoundError(errcode.NotFoundRole)
			}
			return err
		}

		// Protect system roles
		if existing.IsSystem {
			return exceptions.NewBadRequestError(errcode.CannotModifySystemRole, nil)
		}

		// Count active members assigned to this role
		var userCount int64
		if err := tx.Table("user_roles").
			Where("tenant_id = ? AND role_id = ?", tenantID, roleID).
			Count(&userCount).Error; err != nil {
			return err
		}

		if userCount > 0 {
			affectedUsers = int(userCount)
			return exceptions.NewBadRequestError(errcode.RoleInUse, nil)
		}

		// Hard delete custom role to trigger native DB engine ON DELETE RESTRICT
		if err := tx.Where("id = ? AND tenant_id = ?", roleID, tenantID).Delete(&model.Role{}).Error; err != nil {
			return err
		}

		return r.bumpAuthzVersionTx(tx, tenantID)
	})

	return affectedUsers, err
}

func (r *PostgresRoleRepository) CloneRole(ctx context.Context, tenantID, sourceRoleID, newName, newSlug, newDescription string) (*model.Role, error) {
	sourceRole, err := r.GetRoleByID(ctx, tenantID, sourceRoleID)
	if err != nil {
		return nil, err
	}

	var permKeys []string
	for _, p := range sourceRole.Permissions {
		permKeys = append(permKeys, p.PermissionKey)
	}

	newRole := &model.Role{
		ID:          uuid.NewString(),
		TenantID:    tenantID,
		Slug:        strings.ToLower(strings.TrimSpace(newSlug)),
		Name:        strings.TrimSpace(newName),
		Description: strings.TrimSpace(newDescription),
		IsSystem:    false,
		Version:     1,
	}

	if err := r.CreateRole(ctx, newRole, permKeys); err != nil {
		return nil, err
	}

	return r.GetRoleByID(ctx, tenantID, newRole.ID)
}

func (r *PostgresRoleRepository) GetUserRoles(ctx context.Context, tenantID, userID string) ([]model.Role, error) {
	var userRoles []model.UserRole
	err := r.db.WithContext(ctx).
		Preload("Role.Permissions").
		Where("tenant_id = ? AND user_id = ?", tenantID, userID).
		Find(&userRoles).Error
	if err != nil {
		return nil, err
	}

	var roles []model.Role
	for _, ur := range userRoles {
		if ur.Role != nil {
			roles = append(roles, *ur.Role)
		}
	}
	return roles, nil
}

func (r *PostgresRoleRepository) CountActiveOwners(ctx context.Context, tenantID string) (int64, error) {
	var count int64
	query := `
		SELECT count(DISTINCT ur.user_id)
		FROM user_roles ur
		JOIN roles r ON ur.role_id = r.id
		WHERE ur.tenant_id = ? AND r.slug = 'owner' AND r.delete_at IS NULL
	`
	err := r.db.WithContext(ctx).Raw(query, tenantID).Scan(&count).Error
	return count, err
}

func (r *PostgresRoleRepository) AssignUserRoles(ctx context.Context, tenantID, userID string, roleIDs []string, assignedBy string, isActorOwner bool) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Acquire tenant advisory lock to serialize concurrent owner changes
		var lockAcquired bool
		if err := tx.Raw("SELECT pg_advisory_xact_lock(hashtext(?))", "tenant_owner_lock:"+tenantID).Scan(&lockAcquired).Error; err != nil {
			return fmt.Errorf("acquiring tenant advisory lock: %w", err)
		}

		// Pessimistically lock target roles with FOR SHARE
		var targetRoles []model.Role
		if err := tx.Clauses(clause.Locking{Strength: "SHARE"}).
			Where("tenant_id = ? AND id IN (?) AND delete_at IS NULL", tenantID, roleIDs).
			Find(&targetRoles).Error; err != nil {
			return err
		}

		if len(targetRoles) != len(roleIDs) {
			return exceptions.NewNotFoundError(errcode.NotFoundRole)
		}

		// Owner assignment verification
		assigningOwner := false
		for _, tr := range targetRoles {
			if tr.Slug == "owner" {
				assigningOwner = true
				break
			}
		}

		if assigningOwner && !isActorOwner {
			return exceptions.NewBadRequestError(errcode.PrivilegeEscalation, nil)
		}

		// Check if target user currently has owner role
		var currentOwnerCount int64
		_ = tx.Table("user_roles").
			Joins("JOIN roles ON user_roles.role_id = roles.id").
			Where("user_roles.tenant_id = ? AND user_roles.user_id = ? AND roles.slug = 'owner'", tenantID, userID).
			Count(&currentOwnerCount).Error

		if currentOwnerCount > 0 && !assigningOwner {
			// Revoking owner role from user: verify tenant still has at least 1 other active owner
			var totalOwners int64
			query := `
				SELECT count(DISTINCT ur.user_id)
				FROM user_roles ur
				JOIN roles r ON ur.role_id = r.id
				WHERE ur.tenant_id = ? AND r.slug = 'owner' AND r.delete_at IS NULL
			`
			if err := tx.Raw(query, tenantID).Scan(&totalOwners).Error; err != nil {
				return err
			}
			if totalOwners <= 1 {
				return exceptions.NewBadRequestError(errcode.LastOwnerViolation, nil)
			}
		}

		// Clear existing user_roles for user
		if err := tx.Where("tenant_id = ? AND user_id = ?", tenantID, userID).Delete(&model.UserRole{}).Error; err != nil {
			return err
		}

		// Insert new user_roles
		userRoles := make([]model.UserRole, 0, len(roleIDs))
		for _, rid := range roleIDs {
			userRoles = append(userRoles, model.UserRole{
				TenantID:   tenantID,
				UserID:     userID,
				RoleID:     rid,
				AssignedBy: assignedBy,
				AssignedAt: time.Now(),
			})
		}
		if err := tx.Create(&userRoles).Error; err != nil {
			return err
		}

		// Sync legacy users.role column with primary role for backward compatibility
		if len(targetRoles) > 0 {
			primarySlug := targetRoles[0].Slug
			_ = tx.Table("users").Where("tenant_id = ? AND id = ?", tenantID, userID).Update("role", primarySlug).Error
		}

		return r.bumpAuthzVersionTx(tx, tenantID)
	})
}

func (r *PostgresRoleRepository) BumpAuthzVersion(ctx context.Context, tenantID string) (int64, error) {
	var ver int64
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return r.bumpAuthzVersionTx(tx, tenantID)
	})
	if err == nil {
		var av model.AuthzVersion
		_ = r.db.WithContext(ctx).Where("tenant_id = ?", tenantID).First(&av).Error
		ver = av.Version
	}
	return ver, err
}

func (r *PostgresRoleRepository) bumpAuthzVersionTx(tx *gorm.DB, tenantID string) error {
	var av model.AuthzVersion
	err := tx.Where("tenant_id = ?", tenantID).First(&av).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			av = model.AuthzVersion{
				TenantID:  tenantID,
				Version:   1,
				UpdatedAt: time.Now(),
			}
			return tx.Create(&av).Error
		}
		return err
	}
	av.Version++
	av.UpdatedAt = time.Now()
	if err := tx.Save(&av).Error; err != nil {
		return err
	}

	// Trigger PostgreSQL NOTIFY on channel authz_invalidation
	_ = tx.Exec("NOTIFY authz_invalidation, ?", tenantID).Error
	return nil
}
