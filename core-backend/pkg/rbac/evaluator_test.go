package rbac_test

import (
	"testing"

	"github.com/example/core-backend/pkg/rbac"
)

func TestMatchPermission(t *testing.T) {
	tests := []struct {
		pattern    string
		permission string
		expected   bool
	}{
		{"*", "iam:view", true},
		{"*", "billing:plans:update", true},
		{"iam:*", "iam:view", true},
		{"iam:*", "iam:users:manage", true},
		{"iam:*", "tenant:view", false},
		{"iam:view", "iam:view", true},
		{"iam:view", "iam:users:manage", false},
	}

	for _, tc := range tests {
		got := rbac.MatchPermission(tc.pattern, tc.permission)
		if got != tc.expected {
			t.Errorf("MatchPermission(%q, %q) = %v; expected %v", tc.pattern, tc.permission, got, tc.expected)
		}
	}
}

func TestHasPermission_Owner(t *testing.T) {
	if !rbac.HasPermission(string(rbac.RoleOwner), nil, "iam:users:manage") {
		t.Error("expected owner to have any permission")
	}
	if !rbac.HasPermission(string(rbac.RoleOwner), nil, "unknown:custom:perm") {
		t.Error("expected owner to have unknown permissions via wildcard")
	}
}

func TestHasPermission_Admin(t *testing.T) {
	if !rbac.HasPermission(string(rbac.RoleAdmin), nil, "iam:users:manage") {
		t.Error("expected admin to have iam:users:manage via iam:*")
	}
	if !rbac.HasPermission(string(rbac.RoleAdmin), nil, "tenant:delete") {
		t.Error("expected admin to have tenant:delete via tenant:*")
	}
}

func TestHasPermission_Member_And_CustomACL(t *testing.T) {
	// Member by default has iam:view, but NOT iam:users:manage
	if !rbac.HasPermission(string(rbac.RoleMember), nil, "iam:view") {
		t.Error("expected member to have iam:view")
	}
	if rbac.HasPermission(string(rbac.RoleMember), nil, "iam:users:manage") {
		t.Error("expected member not to have iam:users:manage by default")
	}

	// With custom ACL grant, member gets specific permission
	userACL := []string{"iam:users:manage"}
	if !rbac.HasPermission(string(rbac.RoleMember), userACL, "iam:users:manage") {
		t.Error("expected member with ACL grant to have iam:users:manage")
	}
}

func TestGetEffectivePermissions(t *testing.T) {
	ownerPerms := rbac.GetEffectivePermissions(string(rbac.RoleOwner), nil)
	if len(ownerPerms) != 1 || ownerPerms[0] != "*" {
		t.Errorf("expected owner to have ['*'], got %v", ownerPerms)
	}

	viewerPerms := rbac.GetEffectivePermissions(string(rbac.RoleViewer), nil)
	if len(viewerPerms) == 0 {
		t.Errorf("expected viewer to have effective permissions")
	}

	// Viewer with custom ACL
	viewerWithACL := rbac.GetEffectivePermissions(string(rbac.RoleViewer), []string{"billing:plans:update"})
	found := false
	for _, p := range viewerWithACL {
		if p == "billing:plans:update" {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("expected viewerWithACL to contain custom granted permission")
	}
}
