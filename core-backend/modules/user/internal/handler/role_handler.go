package handler

import (
	"log/slog"

	"github.com/example/core-backend/modules/user/internal/dto"
	"github.com/example/core-backend/modules/user/internal/service"
	"github.com/example/core-backend/pkg/helpers"
	"github.com/example/core-backend/pkg/httpx"
	"github.com/example/core-backend/pkg/rbac"
	"github.com/gofiber/fiber/v3"
)

type RoleHandler struct {
	roleService service.IRoleService
}

func NewRoleHandler(roleService service.IRoleService) *RoleHandler {
	return &RoleHandler{roleService: roleService}
}

func (h *RoleHandler) RegisterRoutes(r fiber.Router) {
	// Catalog metadata
	r.Get("/permissions", rbac.RequirePermission("iam:roles:read"), h.GetPermissionsCatalog)

	// Role management
	r.Get("/roles", rbac.RequirePermission("iam:roles:read"), h.ListRoles)
	r.Get("/roles/:id", rbac.RequirePermission("iam:roles:read"), h.GetRoleByID)
	r.Post("/roles", rbac.RequirePermission("iam:roles:manage"), h.CreateRole)
	r.Put("/roles/:id", rbac.RequirePermission("iam:roles:manage"), h.UpdateRole)
	r.Delete("/roles/:id", rbac.RequirePermission("iam:roles:manage"), h.DeleteRole)
	r.Post("/roles/:id/clone", rbac.RequirePermission("iam:roles:manage"), h.CloneRole)

	// User-role assignment
	r.Put("/:id/roles", rbac.RequirePermission("iam:users:manage"), h.AssignUserRoles)
}

func (h *RoleHandler) GetPermissionsCatalog(c fiber.Ctx) error {
	isHost, _ := c.Locals("is_host").(bool)
	catalog := h.roleService.GetPermissionsCatalog(c.Context(), isHost)
	return c.Status(fiber.StatusOK).JSON(httpx.NewResponse(catalog))
}

func (h *RoleHandler) ListRoles(c fiber.Ctx) error {
	roles, err := h.roleService.ListRoles(c.Context())
	if err != nil {
		return err
	}
	return c.Status(fiber.StatusOK).JSON(httpx.NewResponse(roles))
}

func (h *RoleHandler) GetRoleByID(c fiber.Ctx) error {
	id := c.Params("id")
	role, err := h.roleService.GetRoleByID(c.Context(), id)
	if err != nil {
		return err
	}
	return c.Status(fiber.StatusOK).JSON(httpx.NewResponse(role))
}

func (h *RoleHandler) CreateRole(c fiber.Ctx) error {
	var input dto.CreateRoleRequest
	if err := c.Bind().Body(&input); err != nil {
		slog.ErrorContext(c.Context(), "body validation error in CreateRole", "error", err)
		return helpers.ValidationErrorFromBind(err, input)
	}

	role, perms := rbac.ExtractCallerContext(c)
	isHost, _ := c.Locals("is_host").(bool)
	expandedActorPerms := rbac.GetEffectivePermissionsWithScope(role, perms, isHost)

	res, err := h.roleService.CreateRole(c.Context(), expandedActorPerms, isHost, &input)
	if err != nil {
		return err
	}
	return c.Status(fiber.StatusCreated).JSON(httpx.NewResponse(res))
}

func (h *RoleHandler) UpdateRole(c fiber.Ctx) error {
	id := c.Params("id")
	var input dto.UpdateRoleRequest
	if err := c.Bind().Body(&input); err != nil {
		slog.ErrorContext(c.Context(), "body validation error in UpdateRole", "error", err)
		return helpers.ValidationErrorFromBind(err, input)
	}

	role, perms := rbac.ExtractCallerContext(c)
	isHost, _ := c.Locals("is_host").(bool)
	expandedActorPerms := rbac.GetEffectivePermissionsWithScope(role, perms, isHost)

	res, err := h.roleService.UpdateRole(c.Context(), expandedActorPerms, isHost, id, &input)
	if err != nil {
		return err
	}
	return c.Status(fiber.StatusOK).JSON(httpx.NewResponse(res))
}

func (h *RoleHandler) DeleteRole(c fiber.Ctx) error {
	id := c.Params("id")
	affectedUsers, err := h.roleService.DeleteRole(c.Context(), id)
	if err != nil {
		if affectedUsers > 0 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"status":         "error",
				"message_code":   "E032006",
				"message":        "Cannot delete role because it is currently assigned to active member(s).",
				"affected_users": affectedUsers,
			})
		}
		return err
	}
	return c.Status(fiber.StatusOK).JSON(httpx.NewResponse(map[string]string{"deleted_id": id}))
}

func (h *RoleHandler) CloneRole(c fiber.Ctx) error {
	id := c.Params("id")
	var input dto.CloneRoleRequest
	if err := c.Bind().Body(&input); err != nil {
		slog.ErrorContext(c.Context(), "body validation error in CloneRole", "error", err)
		return helpers.ValidationErrorFromBind(err, input)
	}

	role, perms := rbac.ExtractCallerContext(c)
	isHost, _ := c.Locals("is_host").(bool)
	expandedActorPerms := rbac.GetEffectivePermissionsWithScope(role, perms, isHost)

	res, err := h.roleService.CloneRole(c.Context(), expandedActorPerms, isHost, id, &input)
	if err != nil {
		return err
	}
	return c.Status(fiber.StatusCreated).JSON(httpx.NewResponse(res))
}

func (h *RoleHandler) AssignUserRoles(c fiber.Ctx) error {
	targetUserID := c.Params("id")
	var input dto.AssignUserRolesRequest
	if err := c.Bind().Body(&input); err != nil {
		slog.ErrorContext(c.Context(), "body validation error in AssignUserRoles", "error", err)
		return helpers.ValidationErrorFromBind(err, input)
	}

	role, perms := rbac.ExtractCallerContext(c)
	isHost, _ := c.Locals("is_host").(bool)
	expandedActorPerms := rbac.GetEffectivePermissionsWithScope(role, perms, isHost)

	assignedBy, _ := c.Locals("user_id").(string)
	if assignedBy == "" {
		assignedBy = "admin"
	}

	if err := h.roleService.AssignUserRoles(c.Context(), role, expandedActorPerms, isHost, targetUserID, &input, assignedBy); err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(httpx.NewResponse(map[string]any{
		"user_id":  targetUserID,
		"role_ids": input.RoleIDs,
	}))
}
