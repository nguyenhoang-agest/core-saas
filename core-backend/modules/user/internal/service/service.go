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
	"github.com/example/core-backend/pkg/constants/response_code"
	"github.com/example/core-backend/pkg/eventbus"
	"github.com/example/core-backend/pkg/exceptions"
	"github.com/example/core-backend/pkg/httpx"
	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
)

type Repository interface {
	Create(ctx context.Context, user *model.User) error
	GetByID(ctx context.Context, id string) (*model.User, error)
	GetByName(ctx context.Context, name string) (*model.User, error)
	Update(ctx context.Context, user *model.User) error
	Delete(ctx context.Context, id string) error
	GetAll(ctx context.Context, req *dto.PaginationOptions) ([]model.User, int64, error)
}

type IUserService interface {
	Create(ctx context.Context, req *dto.CreateUserRequest) (*model.User, error)
	GetByID(ctx context.Context, id string) (*model.User, error)
	GetAll(ctx context.Context, req *dto.PaginationOptions) (*httpx.PaginatedData[model.User], error)
	Delete(ctx context.Context, id string) error
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

func (s *userService) Create(ctx context.Context, req *dto.CreateUserRequest) (*model.User, error) {
	if req == nil || strings.TrimSpace(req.Name) == "" {
		return nil, exceptions.NewBadRequestError(errcode.UserNameRequired, nil)
	}

	name := strings.TrimSpace(req.Name)

	existing, err := s.repository.GetByName(ctx, name)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		slog.ErrorContext(ctx, "failed to get user by name", "name", name, "error", err)
		return nil, exceptions.NewInternalServerError(errcode.FailToCreateUser)
	}
	if existing != nil {
		return nil, exceptions.NewBadRequestError(errcode.UserAlreadyExists, nil)
	}

	status := req.Status
	if status == "" {
		status = model.StatusActive
	}

	user := &model.User{
		Name:        name,
		Description: strings.TrimSpace(req.Description),
		Status:      status,
	}

	if err := s.repository.Create(ctx, user); err != nil {
		slog.ErrorContext(ctx, "failed to create user", "name", name, "error", err)
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return nil, exceptions.NewBadRequestError(errcode.UserAlreadyExists, nil)
		}
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return nil, exceptions.NewBadRequestError(errcode.UserAlreadyExists, nil)
		}
		return nil, exceptions.NewInternalServerError(errcode.FailToCreateUser)
	}

	if s.bus != nil {
		_ = eventbus.Publish(ctx, s.bus, api.UserCreatedTopic, user)
	}

	return user, nil
}

func (s *userService) GetByID(ctx context.Context, id string) (*model.User, error) {
	if strings.TrimSpace(id) == "" {
		return nil, exceptions.NewBadRequestError(response_code.InvalidParameter, nil)
	}

	user, err := s.repository.GetByID(ctx, id)
	if err != nil {
		slog.ErrorContext(ctx, "failed to get user by id", "id", id, "error", err)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, exceptions.NewNotFoundError(errcode.NotFoundUser)
		}
		return nil, exceptions.NewInternalServerError(errcode.FailToGetUser)
	}

	return user, nil
}

func (s *userService) GetAll(ctx context.Context, req *dto.PaginationOptions) (*httpx.PaginatedData[model.User], error) {
	page := 1
	pageSize := 20
	if req != nil {
		if req.Page > 0 {
			page = req.Page
		}
		if req.PageSize > 0 {
			pageSize = req.PageSize
		}
	}

	items, total, err := s.repository.GetAll(ctx, req)
	if err != nil {
		slog.ErrorContext(ctx, "failed to list users", "error", err)
		return nil, exceptions.NewInternalServerError(errcode.FailToGetUser)
	}

	res := httpx.NewPaginatedData(items, total, page, pageSize)
	return &res, nil
}

func (s *userService) Delete(ctx context.Context, id string) error {
	if strings.TrimSpace(id) == "" {
		return exceptions.NewBadRequestError(response_code.InvalidParameter, nil)
	}

	if err := s.repository.Delete(ctx, id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return exceptions.NewNotFoundError(errcode.NotFoundUser)
		}
		return exceptions.NewInternalServerError(errcode.FailToDeleteUser)
	}

	return nil
}
