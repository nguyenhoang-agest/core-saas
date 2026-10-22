package dto

import "time"

// CreateRoleRequest defines payload for creating a custom role.
type CreateRoleRequest struct {
	Name        string   `json:"name" validate:"required,min=2,max=100"`
	Slug        string   `json:"slug" validate:"required,min=2,max=50"`
	Description string   `json:"description,omitempty" validate:"max=255"`
	Permissions []string `json:"permissions" validate:"required,min=1"`
}

// UpdateRoleRequest defines payload for updating a custom role.
type UpdateRoleRequest struct {
	Name        string   `json:"name" validate:"required,min=2,max=100"`
	Description string   `json:"description,omitempty" validate:"max=255"`
	Permissions []string `json:"permissions" validate:"required,min=1"`
	Version     int      `json:"version" validate:"required,min=1"`
}

// CloneRoleRequest defines payload for cloning an existing role into a new custom role.
type CloneRoleRequest struct {
	Name        string `json:"name" validate:"required,min=2,max=100"`
	Slug        string `json:"slug" validate:"required,min=2,max=50"`
	Description string `json:"description,omitempty" validate:"max=255"`
}

// AssignUserRolesRequest defines payload for assigning multiple roles to a user.
type AssignUserRolesRequest struct {
	RoleIDs []string `json:"role_ids" validate:"required,min=1"`
}

// RoleResponse formats role details for API clients.
type RoleResponse struct {
	ID          string    `json:"id"`
	TenantID    string    `json:"tenant_id"`
	Slug        string    `json:"slug"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	IsSystem    bool      `json:"is_system"`
	Version     int       `json:"version"`
	Permissions []string  `json:"permissions"`
	UserCount   int       `json:"user_count"`
	CreateAt    time.Time `json:"create_at"`
	UpdateAt    time.Time `json:"update_at"`
}

// PermissionResponse formats in-code registry permission metadata for API clients.
type PermissionResponse struct {
	Key         string `json:"key"`
	Module      string `json:"module"`
	Scope       string `json:"scope"`
	I18nNameKey string `json:"i18n_name_key"`
	I18nDescKey string `json:"i18n_desc_key"`
}
