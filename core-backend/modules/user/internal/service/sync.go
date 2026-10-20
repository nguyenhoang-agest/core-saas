package service

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/example/core-backend/modules/user/internal/model"
	"github.com/example/core-backend/pkg/rbac"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// SyncSystemRoles ensures that the 4 baseline system roles (owner, admin, member, viewer)
// exist in the database for the given tenant and that their role_permissions match the in-code templates.
func SyncSystemRoles(ctx context.Context, db *gorm.DB, tenantID string, isHost bool) error {
	if strings.TrimSpace(tenantID) == "" {
		return fmt.Errorf("tenantID cannot be empty")
	}

	reg := rbac.GetRegistry()
	systemRoles := rbac.GetSystemRoles()

	return db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for _, def := range systemRoles {
			slug := string(def.Role)
			var existingRole model.Role
			err := tx.Where("tenant_id = ? AND slug = ?", tenantID, slug).First(&existingRole).Error
			if err != nil && err != gorm.ErrRecordNotFound {
				return fmt.Errorf("checking system role %s for tenant %s: %w", slug, tenantID, err)
			}

			roleID := existingRole.ID
			if err == gorm.ErrRecordNotFound {
				roleID = uuid.NewString()
				newRole := model.Role{
					ID:          roleID,
					TenantID:    tenantID,
					Slug:        slug,
					Name:        def.Name,
					Description: def.Description,
					IsSystem:    true,
					Version:     1,
				}
				if createErr := tx.Create(&newRole).Error; createErr != nil {
					return fmt.Errorf("creating system role %s for tenant %s: %w", slug, tenantID, createErr)
				}
			} else {
				// Ensure metadata matches system definition
				if !existingRole.IsSystem || existingRole.Name != def.Name {
					existingRole.IsSystem = true
					existingRole.Name = def.Name
					existingRole.Description = def.Description
					if saveErr := tx.Save(&existingRole).Error; saveErr != nil {
						return fmt.Errorf("updating system role metadata %s: %w", slug, saveErr)
					}
				}
			}

			// Expand template permissions filtered by host/non-host scope
			concreteKeys := reg.ExpandPermissions(def.DefaultPermissions, isHost)

			// Remove existing permissions for this system role
			if delErr := tx.Where("role_id = ?", roleID).Delete(&model.RolePermission{}).Error; delErr != nil {
				return fmt.Errorf("clearing permissions for system role %s: %w", slug, delErr)
			}

			// Insert updated concrete permission keys
			if len(concreteKeys) > 0 {
				permEntities := make([]model.RolePermission, 0, len(concreteKeys))
				for _, k := range concreteKeys {
					permEntities = append(permEntities, model.RolePermission{
						TenantID:      tenantID,
						RoleID:        roleID,
						PermissionKey: k,
					})
				}
				if batchErr := tx.Create(&permEntities).Error; batchErr != nil {
					return fmt.Errorf("inserting permissions for system role %s: %w", slug, batchErr)
				}
			}
		}

		// Increment or initialize authz_version for tenant
		var authzVer model.AuthzVersion
		if findErr := tx.Where("tenant_id = ?", tenantID).First(&authzVer).Error; findErr != nil {
			if findErr == gorm.ErrRecordNotFound {
				authzVer = model.AuthzVersion{
					TenantID:  tenantID,
					Version:   1,
					UpdatedAt: time.Now(),
				}
				if createErr := tx.Create(&authzVer).Error; createErr != nil {
					return fmt.Errorf("initializing authz_version for tenant %s: %w", tenantID, createErr)
				}
			} else {
				return fmt.Errorf("querying authz_version for tenant %s: %w", tenantID, findErr)
			}
		} else {
			authzVer.Version++
			authzVer.UpdatedAt = time.Now()
			if saveErr := tx.Save(&authzVer).Error; saveErr != nil {
				return fmt.Errorf("updating authz_version for tenant %s: %w", tenantID, saveErr)
			}
		}

		return nil
	})
}

// BackfillLegacyUserRoles inspects users for the specified tenant, identifies any user without
// an entry in user_roles, and creates a user_roles record mapping their legacy users.role string.
func BackfillLegacyUserRoles(ctx context.Context, db *gorm.DB, tenantID string) error {
	var users []model.User
	if err := db.WithContext(ctx).Where("tenant_id = ?", tenantID).Find(&users).Error; err != nil {
		return fmt.Errorf("fetching users for backfill in tenant %s: %w", tenantID, err)
	}

	for _, u := range users {
		var count int64
		if err := db.WithContext(ctx).Model(&model.UserRole{}).
			Where("tenant_id = ? AND user_id = ?", tenantID, u.Id).
			Count(&count).Error; err != nil {
			return fmt.Errorf("checking user_roles for user %s: %w", u.Id, err)
		}

		if count > 0 {
			continue // Already has user_roles
		}

		// Look up role by slug in this tenant
		roleSlug := strings.ToLower(strings.TrimSpace(u.Role))
		if roleSlug == "" {
			roleSlug = string(rbac.RoleMember)
		}

		var targetRole model.Role
		err := db.WithContext(ctx).
			Where("tenant_id = ? AND slug = ? AND delete_at IS NULL", tenantID, roleSlug).
			First(&targetRole).Error
		if err != nil {
			slog.WarnContext(ctx, "backfill: orphaned user with unrecognized role slug",
				"user_id", u.Id, "email", u.Email, "role", u.Role, "tenant_id", tenantID)
			continue
		}

		userRole := model.UserRole{
			TenantID:   tenantID,
			UserID:     u.Id,
			RoleID:     targetRole.ID,
			AssignedBy: "system_migration",
			AssignedAt: time.Now(),
		}

		if insErr := db.WithContext(ctx).Clauses(clause.OnConflict{DoNothing: true}).Create(&userRole).Error; insErr != nil {
			slog.ErrorContext(ctx, "failed to insert user_role backfill", "user_id", u.Id, "role_id", targetRole.ID, "error", insErr)
		}
	}

	return nil
}
