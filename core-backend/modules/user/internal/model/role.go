package model

import (
	"time"

	"gorm.io/gorm"
)

// Role represents a system-defined or tenant-custom role.
type Role struct {
	ID          string         `gorm:"column:id;type:varchar(64);primaryKey" json:"id"`
	TenantID    string         `gorm:"column:tenant_id;type:varchar(64);not null;index:idx_roles_tenant" json:"tenant_id"`
	Slug        string         `gorm:"column:slug;type:varchar(50);not null" json:"slug"`
	Name        string         `gorm:"column:name;type:varchar(255);not null" json:"name"`
	Description string         `gorm:"column:description;type:varchar(255)" json:"description"`
	IsSystem    bool           `gorm:"column:is_system;type:boolean;not null;default:false" json:"is_system"`
	Version     int            `gorm:"column:version;type:int;not null;default:1" json:"version"`
	CreateAt    time.Time      `gorm:"column:create_at;not null;autoCreateTime" json:"create_at"`
	UpdateAt    time.Time      `gorm:"column:update_at;not null;autoUpdateTime" json:"update_at"`
	DeleteAt    gorm.DeletedAt `gorm:"column:delete_at;index" json:"delete_at,omitempty"`

	// Associations
	Permissions []RolePermission `gorm:"foreignKey:RoleID;references:ID" json:"permissions,omitempty"`
	UserCount   int              `gorm:"-" json:"user_count"`
}

func (Role) TableName() string {
	return "roles"
}

// RolePermission associates a role with a concrete permission key.
type RolePermission struct {
	TenantID      string `gorm:"column:tenant_id;type:varchar(64);not null;index:idx_role_permissions_tenant" json:"tenant_id"`
	RoleID        string `gorm:"column:role_id;type:varchar(64);primaryKey" json:"role_id"`
	PermissionKey string `gorm:"column:permission_key;type:varchar(64);primaryKey" json:"permission_key"`
}

func (RolePermission) TableName() string {
	return "role_permissions"
}

// UserRole represents a many-to-many relationship assigning a role to a user.
type UserRole struct {
	TenantID   string    `gorm:"column:tenant_id;type:varchar(64);not null;index:idx_user_roles_tenant" json:"tenant_id"`
	UserID     string    `gorm:"column:user_id;type:varchar(64);primaryKey" json:"user_id"`
	RoleID     string    `gorm:"column:role_id;type:varchar(64);primaryKey;index:idx_user_roles_role" json:"role_id"`
	AssignedBy string    `gorm:"column:assigned_by;type:varchar(64)" json:"assigned_by,omitempty"`
	AssignedAt time.Time `gorm:"column:assigned_at;not null;autoCreateTime" json:"assigned_at"`

	Role *Role `gorm:"foreignKey:RoleID;references:ID" json:"role,omitempty"`
}

func (UserRole) TableName() string {
	return "user_roles"
}

// AuthzVersion tracks the authorization epoch for a tenant to invalidate caches.
type AuthzVersion struct {
	TenantID  string    `gorm:"column:tenant_id;type:varchar(64);primaryKey" json:"tenant_id"`
	Version   int64     `gorm:"column:version;type:bigint;not null;default:1" json:"version"`
	UpdatedAt time.Time `gorm:"column:updated_at;not null;autoUpdateTime" json:"updated_at"`
}

func (AuthzVersion) TableName() string {
	return "authz_version"
}
