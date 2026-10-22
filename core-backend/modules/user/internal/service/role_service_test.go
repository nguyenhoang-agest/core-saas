package service_test

import (
	"context"
	"testing"

	"github.com/example/core-backend/modules/user/internal/dto"
	"github.com/example/core-backend/modules/user/internal/errcode"
	"github.com/example/core-backend/modules/user/internal/model"
	"github.com/example/core-backend/modules/user/internal/service"
	"github.com/example/core-backend/pkg/exceptions"
	"github.com/example/core-backend/pkg/tenantctx"
)

type mockRoleRepo struct {
	roles         map[string]*model.Role
	rolePerms     map[string][]string
	userRoles     map[string][]string // userID -> roleIDs
	affectedUsers int
	deleteErr     error
}

func newMockRoleRepo() *mockRoleRepo {
	return &mockRoleRepo{
		roles:     make(map[string]*model.Role),
		rolePerms: make(map[string][]string),
		userRoles: make(map[string][]string),
	}
}

func (m *mockRoleRepo) ListRoles(ctx context.Context, tenantID string) ([]model.Role, error) {
	var list []model.Role
	for _, r := range m.roles {
		if r.TenantID == tenantID {
			list = append(list, *r)
		}
	}
	return list, nil
}

func (m *mockRoleRepo) GetRoleByID(ctx context.Context, tenantID, roleID string) (*model.Role, error) {
	r, ok := m.roles[roleID]
	if !ok || r.TenantID != tenantID {
		return nil, exceptions.NewNotFoundError(errcode.NotFoundRole)
	}
	return r, nil
}

func (m *mockRoleRepo) GetRoleBySlug(ctx context.Context, tenantID, slug string) (*model.Role, error) {
	for _, r := range m.roles {
		if r.TenantID == tenantID && r.Slug == slug {
			return r, nil
		}
	}
	return nil, exceptions.NewNotFoundError(errcode.NotFoundRole)
}

func (m *mockRoleRepo) CreateRole(ctx context.Context, role *model.Role, permissionKeys []string) error {
	m.roles[role.ID] = role
	m.rolePerms[role.ID] = permissionKeys
	return nil
}

func (m *mockRoleRepo) UpdateRole(ctx context.Context, tenantID, roleID string, name, description string, permissionKeys []string, expectedVersion int) (*model.Role, error) {
	r, ok := m.roles[roleID]
	if !ok || r.TenantID != tenantID {
		return nil, exceptions.NewNotFoundError(errcode.NotFoundRole)
	}
	if r.IsSystem {
		return nil, exceptions.NewBadRequestError(errcode.CannotModifySystemRole, nil)
	}
	if r.Version != expectedVersion {
		return nil, exceptions.NewBadRequestError(errcode.ConcurrencyConflict, nil)
	}
	r.Name = name
	r.Description = description
	r.Version++
	m.rolePerms[roleID] = permissionKeys
	return r, nil
}

func (m *mockRoleRepo) DeleteRole(ctx context.Context, tenantID, roleID string) (int, error) {
	if m.affectedUsers > 0 {
		return m.affectedUsers, exceptions.NewBadRequestError(errcode.RoleInUse, nil)
	}
	delete(m.roles, roleID)
	delete(m.rolePerms, roleID)
	return 0, nil
}

func (m *mockRoleRepo) CloneRole(ctx context.Context, tenantID, sourceRoleID, newName, newSlug, newDescription string) (*model.Role, error) {
	cloned := &model.Role{
		ID:          "cloned-id",
		TenantID:    tenantID,
		Slug:        newSlug,
		Name:        newName,
		Description: newDescription,
		IsSystem:    false,
		Version:     1,
	}
	m.roles[cloned.ID] = cloned
	m.rolePerms[cloned.ID] = m.rolePerms[sourceRoleID]
	return cloned, nil
}

func (m *mockRoleRepo) GetUserRoles(ctx context.Context, tenantID, userID string) ([]model.Role, error) {
	var res []model.Role
	for _, rid := range m.userRoles[userID] {
		if r, ok := m.roles[rid]; ok {
			res = append(res, *r)
		}
	}
	return res, nil
}

func (m *mockRoleRepo) AssignUserRoles(ctx context.Context, tenantID, userID string, roleIDs []string, assignedBy string, isActorOwner bool) error {
	m.userRoles[userID] = roleIDs
	return nil
}

func (m *mockRoleRepo) BumpAuthzVersion(ctx context.Context, tenantID string) (int64, error) {
	return 2, nil
}

func (m *mockRoleRepo) CountActiveOwners(ctx context.Context, tenantID string) (int64, error) {
	return 1, nil
}

func TestRoleService_SubsetRule(t *testing.T) {
	repo := newMockRoleRepo()
	svc := service.NewRoleService(repo, nil)
	ctx := tenantctx.WithContext(context.Background(), "test-tenant")

	// Actor has only iam:view
	actorPerms := []string{"iam:view"}

	// Attempting to grant billing:view (not held by actor) -> must fail with PrivilegeEscalation
	req := &dto.CreateRoleRequest{
		Name:        "Billing Viewer",
		Slug:        "billing-viewer",
		Permissions: []string{"billing:view"},
	}

	_, err := svc.CreateRole(ctx, actorPerms, false, req)
	if err == nil {
		t.Fatal("expected subset rule violation error, got nil")
	}

	customErr, ok := err.(exceptions.BadRequestError)
	if !ok || customErr.MessageCode != errcode.PrivilegeEscalation.Code() {
		t.Fatalf("expected error code %s (PrivilegeEscalation), got %v", errcode.PrivilegeEscalation.Code(), err)
	}

	// Actor with universal wildcard '*' or matching permission can create the role
	actorAdminPerms := []string{"*"}
	created, errSuccess := svc.CreateRole(ctx, actorAdminPerms, false, req)
	if errSuccess != nil {
		t.Fatalf("expected success for authorized actor, got %v", errSuccess)
	}
	if created.Slug != "billing-viewer" {
		t.Fatalf("expected slug billing-viewer, got %s", created.Slug)
	}
}

func TestRoleService_ScopeBoundary(t *testing.T) {
	repo := newMockRoleRepo()
	svc := service.NewRoleService(repo, nil)
	ctx := tenantctx.WithContext(context.Background(), "test-tenant")

	actorPerms := []string{"*"}

	// In non-host tenant, tenant:delete is ScopePlatform and must be rejected
	req := &dto.CreateRoleRequest{
		Name:        "Platform Admin",
		Slug:        "platform-admin",
		Permissions: []string{"tenant:delete"},
	}

	_, err := svc.CreateRole(ctx, actorPerms, false, req)
	if err == nil {
		t.Fatal("expected out of scope permission error for non-host tenant, got nil")
	}

	badReq, ok := err.(exceptions.BadRequestError)
	if !ok || badReq.MessageCode != errcode.OutOfScopePermission.Code() {
		t.Fatalf("expected OutOfScopePermission code %s, got %v", errcode.OutOfScopePermission.Code(), err)
	}

	// In host tenant (isHost = true), platform permission is permitted
	created, errHost := svc.CreateRole(ctx, actorPerms, true, req)
	if errHost != nil {
		t.Fatalf("expected success in host tenant, got: %v", errHost)
	}
	if created == nil {
		t.Fatal("expected created role response")
	}
}

func TestRoleService_UnknownOrWildcardKey(t *testing.T) {
	repo := newMockRoleRepo()
	svc := service.NewRoleService(repo, nil)
	ctx := tenantctx.WithContext(context.Background(), "test-tenant")

	// Custom roles MUST NOT store wildcards (e.g. billing:*)
	reqWildcard := &dto.CreateRoleRequest{
		Name:        "Bad Role",
		Slug:        "bad-role",
		Permissions: []string{"billing:*"},
	}

	_, err := svc.CreateRole(ctx, []string{"*"}, false, reqWildcard)
	if err == nil {
		t.Fatal("expected wildcard key to be rejected for custom roles")
	}

	// Unknown permission keys must be rejected
	reqUnknown := &dto.CreateRoleRequest{
		Name:        "Unknown Perm Role",
		Slug:        "unknown-role",
		Permissions: []string{"nonexistent:custom:action"},
	}
	_, errUnknown := svc.CreateRole(ctx, []string{"*"}, false, reqUnknown)
	if errUnknown == nil {
		t.Fatal("expected unknown permission key to be rejected")
	}
}

func TestRoleService_SystemRoleImmutability(t *testing.T) {
	repo := newMockRoleRepo()
	svc := service.NewRoleService(repo, nil)
	ctx := tenantctx.WithContext(context.Background(), "test-tenant")

	// Attempting to create a custom role with reserved slug "admin"
	req := &dto.CreateRoleRequest{
		Name:        "My Admin",
		Slug:        "admin",
		Permissions: []string{"iam:view"},
	}
	_, err := svc.CreateRole(ctx, []string{"*"}, false, req)
	if err == nil {
		t.Fatal("expected creating role with system slug 'admin' to be rejected")
	}

	// Seed system role in repo
	repo.roles["system-admin-id"] = &model.Role{
		ID:       "system-admin-id",
		TenantID: "test-tenant",
		Slug:     "admin",
		Name:     "Administrator",
		IsSystem: true,
		Version:  1,
	}

	// Attempting to update system role
	updateReq := &dto.UpdateRoleRequest{
		Name:        "Hacked Admin",
		Permissions: []string{"iam:view"},
		Version:     1,
	}
	_, errUpdate := svc.UpdateRole(ctx, []string{"*"}, false, "system-admin-id", updateReq)
	if errUpdate == nil {
		t.Fatal("expected modifying system role to be rejected")
	}
}

func TestRoleService_RoleDeletionSafety(t *testing.T) {
	repo := newMockRoleRepo()
	svc := service.NewRoleService(repo, nil)
	ctx := tenantctx.WithContext(context.Background(), "test-tenant")

	repo.roles["custom-role-1"] = &model.Role{
		ID:       "custom-role-1",
		TenantID: "test-tenant",
		Slug:     "compliance-auditor",
		Name:     "Compliance Auditor",
		IsSystem: false,
		Version:  1,
	}

	// When role is assigned to 3 active users
	repo.affectedUsers = 3
	affected, err := svc.DeleteRole(ctx, "custom-role-1")
	if err == nil {
		t.Fatal("expected deleting role in use to fail")
	}
	if affected != 3 {
		t.Fatalf("expected 3 affected users returned, got %d", affected)
	}

	// When role has 0 active users
	repo.affectedUsers = 0
	_, errSuccess := svc.DeleteRole(ctx, "custom-role-1")
	if errSuccess != nil {
		t.Fatalf("expected successful deletion of unused role, got: %v", errSuccess)
	}
}
