package rbac

import (
	"sort"
	"strings"
)

// MatchPermission checks if a given permission pattern matches a target permission.
// Supports:
// - Universal wildcard: "*" matches any permission.
// - Module wildcard: "iam:*" matches "iam:view", "iam:users:manage", etc.
// - Exact match: "iam:view" matches "iam:view".
func MatchPermission(pattern, permission string) bool {
	if pattern == "*" {
		return true
	}
	if pattern == permission {
		return true
	}
	if strings.HasSuffix(pattern, ":*") {
		prefix := strings.TrimSuffix(pattern, "*")
		return strings.HasPrefix(permission, prefix)
	}
	return false
}

// HasPermission determines whether a user with a given role and direct ACL grants
// possesses the requested permission (defaults isHost to true for backward compatibility).
func HasPermission(userRole string, userACL []string, requiredPermission string) bool {
	return HasPermissionWithScope(userRole, userACL, true, requiredPermission)
}

// HasPermissionWithScope determines whether a user possesses the requested permission,
// enforcing scope boundaries: non-host tenants (isHost == false) can never access
// ScopePlatform permissions (such as tenant:delete or tenant:create) even if holding "*".
func HasPermissionWithScope(userRole string, userACL []string, isHost bool, requiredPermission string) bool {
	// Look up target permission in registry to enforce scope boundary
	if perm, ok := GetRegistry().Get(requiredPermission); ok {
		if perm.Scope == ScopePlatform && !isHost {
			return false
		}
	}

	// Owner evaluation: Host owners are fully unrestricted; Non-host owners are restricted to ScopeTenant
	if userRole == string(RoleOwner) {
		if isHost {
			return true
		}
		if perm, ok := GetRegistry().Get(requiredPermission); ok {
			return perm.Scope == ScopeTenant
		}
		return false
	}

	// Gather role baseline permissions and custom user ACL grants
	var patterns []string
	if userRole != "" {
		patterns = append(patterns, GetDefaultRolePermissions(Role(userRole))...)
	}
	patterns = append(patterns, userACL...)

	expanded := GetRegistry().ExpandPermissions(patterns, isHost)
	for _, p := range expanded {
		if p == requiredPermission {
			return true
		}
	}

	return false
}

// GetEffectivePermissions computes the complete resolved set of concrete permissions
// granted to the user, returning ["*"] for RoleOwner in host tenant mode for backward compatibility.
func GetEffectivePermissions(userRole string, userACL []string) []string {
	if userRole == string(RoleOwner) {
		return []string{"*"}
	}
	return GetEffectivePermissionsWithScope(userRole, userACL, true)
}

// GetEffectivePermissionsWithScope computes the complete resolved set of concrete permissions
// granted to the user, strictly filtering out platform-scoped permissions for non-host tenants.
func GetEffectivePermissionsWithScope(userRole string, userACL []string, isHost bool) []string {
	var patterns []string
	if userRole != "" {
		patterns = append(patterns, GetDefaultRolePermissions(Role(userRole))...)
	}
	patterns = append(patterns, userACL...)

	expanded := GetRegistry().ExpandPermissions(patterns, isHost)
	sort.Strings(expanded)
	return expanded
}
