package rbac

import (
	"fmt"
	"sort"
	"strings"
	"sync"
)

// Scope designates the boundary where a permission applies.
type Scope string

const (
	ScopeTenant   Scope = "tenant"
	ScopePlatform Scope = "platform"
)

// Permission defines a strongly-typed entry in the permission catalog.
type Permission struct {
	Key         string `json:"key"`
	Module      string `json:"module"`
	Scope       Scope  `json:"scope"`
	I18nNameKey string `json:"i18n_name_key"`
	I18nDescKey string `json:"i18n_desc_key"`
}

// Registry maintains the in-code catalog of system capabilities.
type Registry struct {
	mu          sync.RWMutex
	permissions map[string]Permission
}

var globalRegistry = NewRegistry()

// GetRegistry returns the singleton in-code permissions registry.
func GetRegistry() *Registry {
	return globalRegistry
}

// NewRegistry constructs a new isolated Registry instance and initializes standard permissions.
func NewRegistry() *Registry {
	r := &Registry{
		permissions: make(map[string]Permission),
	}
	r.initDefaultCatalog()
	return r
}

func (r *Registry) initDefaultCatalog() {
	standardCatalog := []Permission{
		// IAM / User Management (ScopeTenant)
		{Key: "iam:view", Module: "iam", Scope: ScopeTenant, I18nNameKey: "iam.permissions.view.name", I18nDescKey: "iam.permissions.view.desc"},
		{Key: "iam:users:manage", Module: "iam", Scope: ScopeTenant, I18nNameKey: "iam.permissions.users_manage.name", I18nDescKey: "iam.permissions.users_manage.desc"},
		{Key: "iam:roles:read", Module: "iam", Scope: ScopeTenant, I18nNameKey: "iam.permissions.roles_read.name", I18nDescKey: "iam.permissions.roles_read.desc"},
		{Key: "iam:roles:manage", Module: "iam", Scope: ScopeTenant, I18nNameKey: "iam.permissions.roles_manage.name", I18nDescKey: "iam.permissions.roles_manage.desc"},

		// Tenant Management
		{Key: "tenant:view", Module: "tenant", Scope: ScopeTenant, I18nNameKey: "tenant.permissions.view.name", I18nDescKey: "tenant.permissions.view.desc"},
		{Key: "tenant:list", Module: "tenant", Scope: ScopePlatform, I18nNameKey: "tenant.permissions.list.name", I18nDescKey: "tenant.permissions.list.desc"},
		{Key: "tenant:update", Module: "tenant", Scope: ScopeTenant, I18nNameKey: "tenant.permissions.update.name", I18nDescKey: "tenant.permissions.update.desc"},
		{Key: "tenant:create", Module: "tenant", Scope: ScopePlatform, I18nNameKey: "tenant.permissions.create.name", I18nDescKey: "tenant.permissions.create.desc"},
		{Key: "tenant:delete", Module: "tenant", Scope: ScopePlatform, I18nNameKey: "tenant.permissions.delete.name", I18nDescKey: "tenant.permissions.delete.desc"},

		// Billing (ScopeTenant)
		{Key: "billing:view", Module: "billing", Scope: ScopeTenant, I18nNameKey: "billing.permissions.view.name", I18nDescKey: "billing.permissions.view.desc"},
		{Key: "billing:plans:update", Module: "billing", Scope: ScopeTenant, I18nNameKey: "billing.permissions.plans_update.name", I18nDescKey: "billing.permissions.plans_update.desc"},
		{Key: "billing:invoices:view", Module: "billing", Scope: ScopeTenant, I18nNameKey: "billing.permissions.invoices_view.name", I18nDescKey: "billing.permissions.invoices_view.desc"},

		// Settings (ScopeTenant)
		{Key: "settings:view", Module: "settings", Scope: ScopeTenant, I18nNameKey: "settings.permissions.view.name", I18nDescKey: "settings.permissions.view.desc"},
		{Key: "settings:branding:update", Module: "settings", Scope: ScopeTenant, I18nNameKey: "settings.permissions.branding_update.name", I18nDescKey: "settings.permissions.branding_update.desc"},

		// Audit (ScopeTenant)
		{Key: "audit:view", Module: "audit", Scope: ScopeTenant, I18nNameKey: "audit.permissions.view.name", I18nDescKey: "audit.permissions.view.desc"},

		// Example / Sample Operations (ScopeTenant)
		{Key: "example:view", Module: "example", Scope: ScopeTenant, I18nNameKey: "example.permissions.view.name", I18nDescKey: "example.permissions.view.desc"},
		{Key: "example:write", Module: "example", Scope: ScopeTenant, I18nNameKey: "example.permissions.write.name", I18nDescKey: "example.permissions.write.desc"},
	}

	for _, p := range standardCatalog {
		r.permissions[p.Key] = p
	}
}

// Register adds or updates a permission definition in the registry.
func (r *Registry) Register(p Permission) error {
	if strings.TrimSpace(p.Key) == "" {
		return fmt.Errorf("permission key cannot be empty")
	}
	if strings.Contains(p.Key, "*") {
		return fmt.Errorf("permission key cannot contain wildcards: %s", p.Key)
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	r.permissions[p.Key] = p
	return nil
}

// Get looks up a permission by key.
func (r *Registry) Get(key string) (Permission, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	p, ok := r.permissions[key]
	return p, ok
}

// List returns all permissions registered in the catalog sorted by key.
func (r *Registry) List() []Permission {
	r.mu.RLock()
	defer r.mu.RUnlock()
	list := make([]Permission, 0, len(r.permissions))
	for _, p := range r.permissions {
		list = append(list, p)
	}
	sort.Slice(list, func(i, j int) bool {
		return list[i].Key < list[j].Key
	})
	return list
}

// ValidateConcreteKeys verifies that every key in keys is a recognized permission in the registry
// and does not contain wildcard patterns. Returns an error on the first invalid key found.
func (r *Registry) ValidateConcreteKeys(keys []string) error {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, k := range keys {
		trimmed := strings.TrimSpace(k)
		if trimmed == "" {
			return fmt.Errorf("empty permission key not allowed")
		}
		if strings.Contains(trimmed, "*") {
			return fmt.Errorf("wildcards not permitted in custom roles: %s", trimmed)
		}
		if _, ok := r.permissions[trimmed]; !ok {
			return fmt.Errorf("unknown permission key: %s", trimmed)
		}
	}
	return nil
}

// ExpandPermissions expands wildcard patterns ("*" and "module:*") into concrete permission keys,
// strictly filtering by scope. Non-host tenants (isHost == false) never receive ScopePlatform permissions.
func (r *Registry) ExpandPermissions(patterns []string, isHost bool) []string {
	r.mu.RLock()
	defer r.mu.RUnlock()

	set := make(map[string]bool)

	for _, p := range patterns {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}

		if p == "*" {
			for _, perm := range r.permissions {
				if perm.Scope == ScopeTenant || isHost {
					set[perm.Key] = true
				}
			}
		} else if strings.HasSuffix(p, ":*") {
			prefix := strings.TrimSuffix(p, "*")
			for _, perm := range r.permissions {
				if strings.HasPrefix(perm.Key, prefix) {
					if perm.Scope == ScopeTenant || isHost {
						set[perm.Key] = true
					}
				}
			}
		} else {
			if perm, ok := r.permissions[p]; ok {
				if perm.Scope == ScopeTenant || isHost {
					set[perm.Key] = true
				}
			}
		}
	}

	result := make([]string, 0, len(set))
	for k := range set {
		result = append(result, k)
	}
	sort.Strings(result)
	return result
}
