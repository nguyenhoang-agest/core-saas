package service

import (
	"context"
	"log/slog"
	"strings"

	"github.com/example/core-backend/modules/user/internal/dto"
	"github.com/example/core-backend/modules/user/internal/errcode"
	"github.com/example/core-backend/modules/user/internal/model"
	"github.com/example/core-backend/modules/user/internal/repository"
	"github.com/example/core-backend/pkg/constants/response_code"
	"github.com/example/core-backend/pkg/eventbus"
	"github.com/example/core-backend/pkg/exceptions"
	"github.com/example/core-backend/pkg/rbac"
	"github.com/example/core-backend/pkg/tenantctx"
	"github.com/google/uuid"
)

const (
	RoleCreatedTopic      = "iam.role.created"
	RoleUpdatedTopic      = "iam.role.updated"
	RoleDeletedTopic      = "iam.role.deleted"
	UserRolesChangedTopic = "iam.user_roles.changed"
)

type IRoleService interface {
	ListRoles(ctx context.Context) ([]dto.RoleResponse, error)
	GetRoleByID(ctx context.Context, id string) (*dto.RoleResponse, error)
	CreateRole(ctx context.Context, actorPermissions []string, isHost bool, req *dto.CreateRoleRequest) (*dto.RoleResponse, error)
	UpdateRole(ctx context.Context, actorPermissions []string, isHost bool, id string, req *dto.UpdateRoleRequest) (*dto.RoleResponse, error)
	DeleteRole(ctx context.Context, id string) (int, error)
	CloneRole(ctx context.Context, actorPermissions []string, isHost bool, id string, req *dto.CloneRoleRequest) (*dto.RoleResponse, error)
	AssignUserRoles(ctx context.Context, actorRole string, actorPermissions []string, isHost bool, targetUserID string, req *dto.AssignUserRolesRequest, assignedBy string) error
	GetPermissionsCatalog(ctx context.Context, isHost bool) []dto.PermissionResponse
}

type RoleService struct {
	repo repository.IRoleRepository
	bus  *eventbus.Bus
}

func NewRoleService(repo repository.IRoleRepository, bus *eventbus.Bus) *RoleService {
	return &RoleService{
		repo: repo,
		bus:  bus,
	}
}

func (s *RoleService) ListRoles(ctx context.Context) ([]dto.RoleResponse, error) {
	tenantID := tenantctx.GetTenantID(ctx)
	roles, err := s.repo.ListRoles(ctx, tenantID)
	if err != nil {
		return nil, exceptions.NewInternalServerError(response_code.DatabaseError)
	}

	result := make([]dto.RoleResponse, 0, len(roles))
	for _, r := range roles {
		result = append(result, toRoleResponse(&r))
	}
	return result, nil
}

func (s *RoleService) GetRoleByID(ctx context.Context, id string) (*dto.RoleResponse, error) {
	tenantID := tenantctx.GetTenantID(ctx)
	role, err := s.repo.GetRoleByID(ctx, tenantID, id)
	if err != nil {
		return nil, err
	}
	resp := toRoleResponse(role)
	return &resp, nil
}

func (s *RoleService) CreateRole(ctx context.Context, actorPermissions []string, isHost bool, req *dto.CreateRoleRequest) (*dto.RoleResponse, error) {
	if req == nil {
		return nil, exceptions.NewBadRequestError(response_code.InvalidParameter, nil)
	}

	name := strings.TrimSpace(req.Name)
	if name == "" {
		return nil, exceptions.NewBadRequestError(errcode.RoleNameRequired, nil)
	}

	slug := strings.ToLower(strings.TrimSpace(req.Slug))
	if slug == "" {
		return nil, exceptions.NewBadRequestError(errcode.RoleSlugRequired, nil)
	}

	// Reject reserved system role slugs
	if rbac.IsSystemRoleSlug(slug) {
		return nil, exceptions.NewBadRequestError(errcode.CannotModifySystemRole, nil)
	}

	// Validate concrete permission keys against in-code registry (no wildcards allowed)
	reg := rbac.GetRegistry()
	if err := reg.ValidateConcreteKeys(req.Permissions); err != nil {
		return nil, exceptions.NewBadRequestError(errcode.UnknownPermissionKey, []exceptions.FieldErrorDetail{
			{Field: "permissions", MessageCode: err.Error()},
		})
	}

	// Enforce Scope Boundaries: non-host tenants cannot include ScopePlatform permissions
	if !isHost {
		for _, k := range req.Permissions {
			if perm, ok := reg.Get(k); ok && perm.Scope == rbac.ScopePlatform {
				return nil, exceptions.NewBadRequestError(errcode.OutOfScopePermission, nil)
			}
		}
	}

	// Enforce Subset Rule: actor must possess all permissions they are attempting to grant
	if !hasAllPermissions(actorPermissions, req.Permissions) {
		return nil, exceptions.NewBadRequestError(errcode.PrivilegeEscalation, nil)
	}

	tenantID := tenantctx.GetTenantID(ctx)
	role := &model.Role{
		ID:          uuid.NewString(),
		TenantID:    tenantID,
		Slug:        slug,
		Name:        name,
		Description: strings.TrimSpace(req.Description),
		IsSystem:    false,
		Version:     1,
	}

	if err := s.repo.CreateRole(ctx, role, req.Permissions); err != nil {
		return nil, err
	}

	// Audit trail event
	if s.bus != nil {
		_ = eventbus.Publish(ctx, s.bus, RoleCreatedTopic, map[string]string{
			"tenant_id": tenantID,
			"role_id":   role.ID,
			"slug":      role.Slug,
			"name":      role.Name,
		})
		_ = eventbus.Publish(ctx, s.bus, AuthzInvalidatedTopic, map[string]string{
			"tenant_id": tenantID,
		})
	}

	created, err := s.repo.GetRoleByID(ctx, tenantID, role.ID)
	if err != nil {
		return nil, err
	}
	resp := toRoleResponse(created)
	return &resp, nil
}

func (s *RoleService) UpdateRole(ctx context.Context, actorPermissions []string, isHost bool, id string, req *dto.UpdateRoleRequest) (*dto.RoleResponse, error) {
	if req == nil {
		return nil, exceptions.NewBadRequestError(response_code.InvalidParameter, nil)
	}

	name := strings.TrimSpace(req.Name)
	if name == "" {
		return nil, exceptions.NewBadRequestError(errcode.RoleNameRequired, nil)
	}

	reg := rbac.GetRegistry()
	if err := reg.ValidateConcreteKeys(req.Permissions); err != nil {
		return nil, exceptions.NewBadRequestError(errcode.UnknownPermissionKey, []exceptions.FieldErrorDetail{
			{Field: "permissions", MessageCode: err.Error()},
		})
	}

	if !isHost {
		for _, k := range req.Permissions {
			if perm, ok := reg.Get(k); ok && perm.Scope == rbac.ScopePlatform {
				return nil, exceptions.NewBadRequestError(errcode.OutOfScopePermission, nil)
			}
		}
	}

	if !hasAllPermissions(actorPermissions, req.Permissions) {
		return nil, exceptions.NewBadRequestError(errcode.PrivilegeEscalation, nil)
	}

	tenantID := tenantctx.GetTenantID(ctx)
	updated, err := s.repo.UpdateRole(ctx, tenantID, id, name, req.Description, req.Permissions, req.Version)
	if err != nil {
		return nil, err
	}

	if s.bus != nil {
		_ = eventbus.Publish(ctx, s.bus, RoleUpdatedTopic, map[string]string{
			"tenant_id": tenantID,
			"role_id":   id,
			"name":      name,
		})
		_ = eventbus.Publish(ctx, s.bus, AuthzInvalidatedTopic, map[string]string{
			"tenant_id": tenantID,
		})
	}

	resp := toRoleResponse(updated)
	return &resp, nil
}

func (s *RoleService) DeleteRole(ctx context.Context, id string) (int, error) {
	tenantID := tenantctx.GetTenantID(ctx)
	affectedUsers, err := s.repo.DeleteRole(ctx, tenantID, id)
	if err != nil {
		return affectedUsers, err
	}

	if s.bus != nil {
		_ = eventbus.Publish(ctx, s.bus, RoleDeletedTopic, map[string]string{
			"tenant_id": tenantID,
			"role_id":   id,
		})
		_ = eventbus.Publish(ctx, s.bus, AuthzInvalidatedTopic, map[string]string{
			"tenant_id": tenantID,
		})
	}

	return 0, nil
}

func (s *RoleService) CloneRole(ctx context.Context, actorPermissions []string, isHost bool, id string, req *dto.CloneRoleRequest) (*dto.RoleResponse, error) {
	if req == nil {
		return nil, exceptions.NewBadRequestError(response_code.InvalidParameter, nil)
	}

	name := strings.TrimSpace(req.Name)
	if name == "" {
		return nil, exceptions.NewBadRequestError(errcode.RoleNameRequired, nil)
	}

	slug := strings.ToLower(strings.TrimSpace(req.Slug))
	if slug == "" {
		return nil, exceptions.NewBadRequestError(errcode.RoleSlugRequired, nil)
	}

	if rbac.IsSystemRoleSlug(slug) {
		return nil, exceptions.NewBadRequestError(errcode.CannotModifySystemRole, nil)
	}

	tenantID := tenantctx.GetTenantID(ctx)
	source, err := s.repo.GetRoleByID(ctx, tenantID, id)
	if err != nil {
		return nil, err
	}

	var sourcePerms []string
	for _, p := range source.Permissions {
		sourcePerms = append(sourcePerms, p.PermissionKey)
	}

	if !hasAllPermissions(actorPermissions, sourcePerms) {
		return nil, exceptions.NewBadRequestError(errcode.PrivilegeEscalation, nil)
	}

	cloned, err := s.repo.CloneRole(ctx, tenantID, id, name, slug, req.Description)
	if err != nil {
		return nil, err
	}

	if s.bus != nil {
		_ = eventbus.Publish(ctx, s.bus, RoleCreatedTopic, map[string]string{
			"tenant_id": tenantID,
			"role_id":   cloned.ID,
			"slug":      cloned.Slug,
			"name":      cloned.Name,
		})
		_ = eventbus.Publish(ctx, s.bus, AuthzInvalidatedTopic, map[string]string{
			"tenant_id": tenantID,
		})
	}

	resp := toRoleResponse(cloned)
	return &resp, nil
}

func (s *RoleService) AssignUserRoles(ctx context.Context, actorRole string, actorPermissions []string, isHost bool, targetUserID string, req *dto.AssignUserRolesRequest, assignedBy string) error {
	if req == nil || len(req.RoleIDs) == 0 {
		return exceptions.NewBadRequestError(response_code.InvalidParameter, nil)
	}

	tenantID := tenantctx.GetTenantID(ctx)

	// Fetch target roles to verify subset rule against assigned roles' permissions
	for _, rid := range req.RoleIDs {
		targetRole, err := s.repo.GetRoleByID(ctx, tenantID, rid)
		if err != nil {
			return err
		}

		var perms []string
		for _, p := range targetRole.Permissions {
			perms = append(perms, p.PermissionKey)
		}

		// Subset Rule: actor must possess all permissions of the roles they assign
		if !hasAllPermissions(actorPermissions, perms) {
			slog.WarnContext(ctx, "subset rule violation on role assignment",
				"actor", assignedBy, "target_user", targetUserID, "role", targetRole.Slug)
			return exceptions.NewBadRequestError(errcode.PrivilegeEscalation, nil)
		}
	}

	isActorOwner := actorRole == string(rbac.RoleOwner)
	if err := s.repo.AssignUserRoles(ctx, tenantID, targetUserID, req.RoleIDs, assignedBy, isActorOwner); err != nil {
		return err
	}

	if s.bus != nil {
		_ = eventbus.Publish(ctx, s.bus, UserRolesChangedTopic, map[string]string{
			"tenant_id":   tenantID,
			"user_id":     targetUserID,
			"assigned_by": assignedBy,
		})
		_ = eventbus.Publish(ctx, s.bus, AuthzInvalidatedTopic, map[string]string{
			"tenant_id": tenantID,
			"user_id":   targetUserID,
		})
	}

	return nil
}

func (s *RoleService) GetPermissionsCatalog(ctx context.Context, isHost bool) []dto.PermissionResponse {
	all := rbac.GetRegistry().List()
	res := make([]dto.PermissionResponse, 0, len(all))

	for _, p := range all {
		if p.Scope == rbac.ScopeTenant || isHost {
			res = append(res, dto.PermissionResponse{
				Key:         p.Key,
				Module:      p.Module,
				Scope:       string(p.Scope),
				I18nNameKey: p.I18nNameKey,
				I18nDescKey: p.I18nDescKey,
			})
		}
	}
	return res
}

func hasAllPermissions(actorPerms, requestedPerms []string) bool {
	permSet := make(map[string]bool)
	for _, p := range actorPerms {
		if p == "*" {
			return true // Universal wildcard covers everything
		}
		permSet[p] = true
	}

	for _, req := range requestedPerms {
		if !permSet[req] {
			return false
		}
	}
	return true
}

func toRoleResponse(r *model.Role) dto.RoleResponse {
	perms := make([]string, 0, len(r.Permissions))
	for _, p := range r.Permissions {
		perms = append(perms, p.PermissionKey)
	}

	return dto.RoleResponse{
		ID:          r.ID,
		TenantID:    r.TenantID,
		Slug:        r.Slug,
		Name:        r.Name,
		Description: r.Description,
		IsSystem:    r.IsSystem,
		Version:     r.Version,
		Permissions: perms,
		UserCount:   r.UserCount,
		CreateAt:    r.CreateAt,
		UpdateAt:    r.UpdateAt,
	}
}
