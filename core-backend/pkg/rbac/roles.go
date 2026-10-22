package rbac

type Role string

const (
	RoleOwner  Role = "owner"
	RoleAdmin  Role = "admin"
	RoleMember Role = "member"
	RoleViewer Role = "viewer"
)

func (r Role) String() string {
	return string(r)
}

func (r Role) IsValid() bool {
	switch r {
	case RoleOwner, RoleAdmin, RoleMember, RoleViewer:
		return true
	default:
		return false
	}
}

// IsSystemRoleSlug returns true if the slug matches a reserved built-in system role.
func IsSystemRoleSlug(slug string) bool {
	switch Role(slug) {
	case RoleOwner, RoleAdmin, RoleMember, RoleViewer:
		return true
	default:
		return false
	}
}

// GetSystemRoleSlugs returns the list of all reserved system role slugs.
func GetSystemRoleSlugs() []string {
	return []string{string(RoleOwner), string(RoleAdmin), string(RoleMember), string(RoleViewer)}
}

// RoleDefinition encapsulates metadata and baseline permissions for a role.
type RoleDefinition struct {
	Role               Role     `json:"role"`
	Name               string   `json:"name"`
	Description        string   `json:"description"`
	DefaultPermissions []string `json:"default_permissions"`
}

// PermissionDefinition represents a cataloged system permission (maintained for backward-compatibility).
type PermissionDefinition struct {
	ID          string `json:"id"`
	Module      string `json:"module"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Scope       Scope  `json:"scope,omitempty"`
}

// Default role permission template mappings
var defaultRolePermissions = map[Role][]string{
	RoleOwner: {"*"},
	RoleAdmin: {
		"iam:*",
		"tenant:*",
		"billing:*",
		"settings:*",
		"audit:*",
		"example:*",
	},
	RoleMember: {
		"iam:view",
		"tenant:view",
		"billing:view",
		"billing:invoices:view",
		"settings:view",
		"example:*",
	},
	RoleViewer: {
		"iam:view",
		"tenant:view",
		"billing:view",
		"settings:view",
		"example:view",
	},
}

// SystemPermissionsCatalog provides the legacy list of permissions mapped from the registry.
var SystemPermissionsCatalog = []PermissionDefinition{
	{ID: "iam:view", Module: "iam", Name: "View IAM & Users", Description: "Browse organization directory and members", Scope: ScopeTenant},
	{ID: "iam:users:manage", Module: "iam", Name: "Manage Users", Description: "Invite, create, update, and remove workspace members", Scope: ScopeTenant},
	{ID: "iam:roles:read", Module: "iam", Name: "Read Roles", Description: "View role details and permission mappings", Scope: ScopeTenant},
	{ID: "iam:roles:manage", Module: "iam", Name: "Manage Roles & ACL", Description: "Modify role assignments and custom ACL permissions", Scope: ScopeTenant},

	{ID: "tenant:view", Module: "tenant", Name: "View Tenants", Description: "Inspect tenant workspaces and isolation modes", Scope: ScopeTenant},
	{ID: "tenant:list", Module: "tenant", Name: "List All Tenants", Description: "Browse all tenant workspaces across the platform", Scope: ScopePlatform},
	{ID: "tenant:update", Module: "tenant", Name: "Update Tenants", Description: "Modify tenant settings and connection configurations", Scope: ScopeTenant},
	{ID: "tenant:create", Module: "tenant", Name: "Create Tenants", Description: "Provision new tenant workspaces", Scope: ScopePlatform},
	{ID: "tenant:delete", Module: "tenant", Name: "Delete Tenants", Description: "Suspend or de-provision tenant workspaces", Scope: ScopePlatform},

	{ID: "billing:view", Module: "billing", Name: "View Billing", Description: "Access billing dashboard and subscription state", Scope: ScopeTenant},
	{ID: "billing:plans:update", Module: "billing", Name: "Update Plans", Description: "Change subscription tiers and add-ons", Scope: ScopeTenant},
	{ID: "billing:invoices:view", Module: "billing", Name: "View Invoices", Description: "Inspect and download invoices", Scope: ScopeTenant},

	{ID: "settings:view", Module: "settings", Name: "View Settings", Description: "View organization configuration", Scope: ScopeTenant},
	{ID: "settings:branding:update", Module: "settings", Name: "Update Branding", Description: "Update custom logos, titles, and themes", Scope: ScopeTenant},

	{ID: "audit:view", Module: "audit", Name: "View Audit Logs", Description: "Inspect historical audit activity logs", Scope: ScopeTenant},

	{ID: "example:view", Module: "example", Name: "View Example Items", Description: "Read sample module domain items", Scope: ScopeTenant},
	{ID: "example:write", Module: "example", Name: "Write Example Items", Description: "Create, update, or remove sample domain items", Scope: ScopeTenant},
}

// GetDefaultRolePermissions returns the baseline permissions assigned to a given role.
func GetDefaultRolePermissions(role Role) []string {
	if perms, ok := defaultRolePermissions[role]; ok {
		res := make([]string, len(perms))
		copy(res, perms)
		return res
	}
	return []string{}
}

// GetSystemRoles returns metadata for all standard roles.
func GetSystemRoles() []RoleDefinition {
	return []RoleDefinition{
		{
			Role:               RoleOwner,
			Name:               "Owner",
			Description:        "Full unrestricted administrative access to the entire workspace and all modules.",
			DefaultPermissions: GetDefaultRolePermissions(RoleOwner),
		},
		{
			Role:               RoleAdmin,
			Name:               "Administrator",
			Description:        "Manage organization members, workspace configurations, and billing settings.",
			DefaultPermissions: GetDefaultRolePermissions(RoleAdmin),
		},
		{
			Role:               RoleMember,
			Name:               "Member",
			Description:        "Standard team collaborator with operational access and view rights.",
			DefaultPermissions: GetDefaultRolePermissions(RoleMember),
		},
		{
			Role:               RoleViewer,
			Name:               "Viewer",
			Description:        "Read-only visibility across organization workspaces and domain resources.",
			DefaultPermissions: GetDefaultRolePermissions(RoleViewer),
		},
	}
}
