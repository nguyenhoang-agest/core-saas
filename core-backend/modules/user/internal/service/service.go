package service

import (
	"context"
	"errors"
	"log/slog"
	"strings"

	"github.com/example/core-backend/modules/user/api"
	"github.com/example/core-backend/modules/user/internal/dto"
	"github.com/example/core-backend/modules/user/internal/errcode"
	"github.com/example/core-backend/modules/user/internal/model"
	"github.com/example/core-backend/pkg/auth"
	"github.com/example/core-backend/pkg/constants/response_code"
	"github.com/example/core-backend/pkg/eventbus"
	"github.com/example/core-backend/pkg/exceptions"
	"github.com/example/core-backend/pkg/httpx"
	"github.com/example/core-backend/pkg/rbac"
	"github.com/example/core-backend/pkg/tenantctx"
	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
)

type Repository interface {
	Create(ctx context.Context, user *model.User) error
	GetByID(ctx context.Context, tenantID, id string) (*model.User, error)
	GetByEmail(ctx context.Context, tenantID, email string) (*model.User, error)
	Update(ctx context.Context, user *model.User) error
	Delete(ctx context.Context, tenantID, id string) error
	GetAll(ctx context.Context, tenantID string, req *dto.PaginationOptions) ([]model.User, int64, error)
}

type IUserService interface {
	Create(ctx context.Context, req *dto.CreateUserRequest) (*dto.UserResponse, error)
	GetByID(ctx context.Context, id string) (*dto.UserResponse, error)
	GetAll(ctx context.Context, req *dto.PaginationOptions) (*httpx.PaginatedData[dto.UserResponse], error)
	Update(ctx context.Context, id string, req *dto.UpdateUserRequest) (*dto.UserResponse, error)
	UpdateStatus(ctx context.Context, id string, status model.Status) (*dto.UserResponse, error)
	Delete(ctx context.Context, id string) error
	GetRolesMeta() []rbac.RoleDefinition
	GetPermissionsMeta() []rbac.PermissionDefinition
}

type userService struct {
	repository Repository
	bus        *eventbus.Bus
}

func NewUserService(repository Repository, bus *eventbus.Bus) IUserService {
	return &userService{
		repository: repository,
		bus:        bus,
	}
}

func (s *userService) Create(ctx context.Context, req *dto.CreateUserRequest) (*dto.UserResponse, error) {
	if req == nil || strings.TrimSpace(req.Email) == "" {
		return nil, exceptions.NewBadRequestError(errcode.UserEmailRequired, nil)
	}
	if strings.TrimSpace(req.Name) == "" {
		return nil, exceptions.NewBadRequestError(errcode.UserNameRequired, nil)
	}

	role := rbac.Role(req.Role)
	if !role.IsValid() {
		return nil, exceptions.NewBadRequestError(errcode.UserInvalidRole, nil)
	}

	status := req.Status
	if status == "" {
		status = model.StatusActive
	} else if !status.IsValid() {
		return nil, exceptions.NewBadRequestError(errcode.UserInvalidStatus, nil)
	}

	tenantID := tenantctx.GetTenantID(ctx)
	email := strings.ToLower(strings.TrimSpace(req.Email))

	existing, err := s.repository.GetByEmail(ctx, tenantID, email)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		slog.ErrorContext(ctx, "failed to query user by email", "email", email, "tenant_id", tenantID, "error", err)
		return nil, exceptions.NewInternalServerError(errcode.FailToGetUser)
	}
	if existing != nil {
		return nil, exceptions.NewBadRequestError(errcode.UserAlreadyExists, nil)
	}

	var passwordHash string
	if req.Password != "" {
		hash, err := auth.HashPassword(req.Password)
		if err != nil {
			slog.ErrorContext(ctx, "failed to hash user password", "error", err)
			return nil, exceptions.NewInternalServerError(response_code.InternalServerGeneral)
		}
		passwordHash = hash
	}

	user := &model.User{
		TenantID:     tenantID,
		Email:        email,
		Name:         strings.TrimSpace(req.Name),
		Role:         string(role),
		Permissions:  model.StringArray(req.Permissions),
		Status:       status,
		Phone:        strings.TrimSpace(req.Phone),
		PasswordHash: passwordHash,
	}

	if err := s.repository.Create(ctx, user); err != nil {
		slog.ErrorContext(ctx, "failed to create user", "email", email, "error", err)
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return nil, exceptions.NewBadRequestError(errcode.UserAlreadyExists, nil)
		}
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return nil, exceptions.NewBadRequestError(errcode.UserAlreadyExists, nil)
		}
		return nil, exceptions.NewInternalServerError(errcode.FailToCreateUser)
	}

	resp := dto.ToUserResponse(user)

	if s.bus != nil {
		_ = eventbus.Publish(ctx, s.bus, api.UserCreatedTopic, resp)
	}

	return resp, nil
}

func (s *userService) GetByID(ctx context.Context, id string) (*dto.UserResponse, error) {
	if strings.TrimSpace(id) == "" {
		return nil, exceptions.NewBadRequestError(response_code.InvalidParameter, nil)
	}

	tenantID := tenantctx.GetTenantID(ctx)
	user, err := s.repository.GetByID(ctx, tenantID, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, exceptions.NewNotFoundError(errcode.NotFoundUser)
		}
		slog.ErrorContext(ctx, "failed to get user by id", "id", id, "error", err)
		return nil, exceptions.NewInternalServerError(errcode.FailToGetUser)
	}

	return dto.ToUserResponse(user), nil
}

func (s *userService) GetAll(ctx context.Context, req *dto.PaginationOptions) (*httpx.PaginatedData[dto.UserResponse], error) {
	if req == nil {
		req = &dto.PaginationOptions{}
	}
	req.Normalize()

	tenantID := tenantctx.GetTenantID(ctx)
	items, total, err := s.repository.GetAll(ctx, tenantID, req)
	if err != nil {
		slog.ErrorContext(ctx, "failed to list users", "tenant_id", tenantID, "error", err)
		return nil, exceptions.NewInternalServerError(errcode.FailToGetUser)
	}

	responses := make([]dto.UserResponse, len(items))
	for i := range items {
		responses[i] = *dto.ToUserResponse(&items[i])
	}

	res := httpx.NewPaginatedData(responses, total, req.Page, req.PageSize)
	return &res, nil
}

func (s *userService) Update(ctx context.Context, id string, req *dto.UpdateUserRequest) (*dto.UserResponse, error) {
	if strings.TrimSpace(id) == "" {
		return nil, exceptions.NewBadRequestError(response_code.InvalidParameter, nil)
	}

	tenantID := tenantctx.GetTenantID(ctx)
	user, err := s.repository.GetByID(ctx, tenantID, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, exceptions.NewNotFoundError(errcode.NotFoundUser)
		}
		return nil, exceptions.NewInternalServerError(errcode.FailToGetUser)
	}

	if req.Name != "" {
		user.Name = strings.TrimSpace(req.Name)
	}
	if req.Role != "" {
		role := rbac.Role(req.Role)
		if !role.IsValid() {
			return nil, exceptions.NewBadRequestError(errcode.UserInvalidRole, nil)
		}
		user.Role = string(role)
	}
	if req.Permissions != nil {
		user.Permissions = model.StringArray(*req.Permissions)
	}
	if req.Status != "" {
		if !req.Status.IsValid() {
			return nil, exceptions.NewBadRequestError(errcode.UserInvalidStatus, nil)
		}
		user.Status = req.Status
	}
	if req.Phone != "" {
		user.Phone = strings.TrimSpace(req.Phone)
	}

	if err := s.repository.Update(ctx, user); err != nil {
		slog.ErrorContext(ctx, "failed to update user", "id", id, "error", err)
		return nil, exceptions.NewInternalServerError(errcode.FailToUpdateUser)
	}

	resp := dto.ToUserResponse(user)

	if s.bus != nil {
		_ = eventbus.Publish(ctx, s.bus, api.UserUpdatedTopic, resp)
	}

	return resp, nil
}

func (s *userService) UpdateStatus(ctx context.Context, id string, status model.Status) (*dto.UserResponse, error) {
	if !status.IsValid() {
		return nil, exceptions.NewBadRequestError(errcode.UserInvalidStatus, nil)
	}

	return s.Update(ctx, id, &dto.UpdateUserRequest{
		Status: status,
	})
}

func (s *userService) Delete(ctx context.Context, id string) error {
	if strings.TrimSpace(id) == "" {
		return exceptions.NewBadRequestError(response_code.InvalidParameter, nil)
	}

	tenantID := tenantctx.GetTenantID(ctx)
	if err := s.repository.Delete(ctx, tenantID, id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return exceptions.NewNotFoundError(errcode.NotFoundUser)
		}
		slog.ErrorContext(ctx, "failed to delete user", "id", id, "error", err)
		return exceptions.NewInternalServerError(errcode.FailToDeleteUser)
	}

	if s.bus != nil {
		_ = eventbus.Publish(ctx, s.bus, api.UserDeletedTopic, map[string]string{
			"id":        id,
			"tenant_id": tenantID,
		})
	}

	return nil
}

func (s *userService) GetRolesMeta() []rbac.RoleDefinition {
	return rbac.GetSystemRoles()
}

func (s *userService) GetPermissionsMeta() []rbac.PermissionDefinition {
	return rbac.SystemPermissionsCatalog
}
