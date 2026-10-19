package service

import (
	"context"
	"errors"
	"log/slog"
	"strings"

	"github.com/example/core-backend/modules/_template/api"
	"github.com/example/core-backend/modules/_template/internal/dto"
	"github.com/example/core-backend/modules/_template/internal/errcode"
	"github.com/example/core-backend/modules/_template/internal/model"
	"github.com/example/core-backend/pkg/constants/response_code"
	"github.com/example/core-backend/pkg/eventbus"
	"github.com/example/core-backend/pkg/exceptions"
	"github.com/example/core-backend/pkg/httpx"
	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
)

type Repository interface {
	Create(ctx context.Context, template *model.Template) error
	GetByID(ctx context.Context, id string) (*model.Template, error)
	GetByName(ctx context.Context, name string) (*model.Template, error)
	Update(ctx context.Context, template *model.Template) error
	Delete(ctx context.Context, id string) error
	GetAll(ctx context.Context, req *dto.PaginationOptions) ([]model.Template, int64, error)
}

type ITemplateService interface {
	Create(ctx context.Context, req *dto.CreateTemplateRequest) (*model.Template, error)
	GetByID(ctx context.Context, id string) (*model.Template, error)
	GetAll(ctx context.Context, req *dto.PaginationOptions) (*httpx.PaginatedData[model.Template], error)
	Delete(ctx context.Context, id string) error
}

type templateService struct {
	repository Repository
	bus        *eventbus.Bus
}

func NewTemplateService(repository Repository, bus *eventbus.Bus) ITemplateService {
	return &templateService{
		repository: repository,
		bus:        bus,
	}
}

func (s *templateService) Create(ctx context.Context, req *dto.CreateTemplateRequest) (*model.Template, error) {
	if req == nil || strings.TrimSpace(req.Name) == "" {
		return nil, exceptions.NewBadRequestError(errcode.TemplateNameRequired, nil)
	}

	name := strings.TrimSpace(req.Name)

	existing, err := s.repository.GetByName(ctx, name)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		slog.ErrorContext(ctx, "failed to get template by name", "name", name, "error", err)
		return nil, exceptions.NewInternalServerError(errcode.FailToCreateTemplate)
	}
	if existing != nil {
		return nil, exceptions.NewBadRequestError(errcode.TemplateAlreadyExists, nil)
	}

	status := req.Status
	if status == "" {
		status = model.StatusActive
	}

	template := &model.Template{
		Name:        name,
		Description: strings.TrimSpace(req.Description),
		Status:      status,
	}

	if err := s.repository.Create(ctx, template); err != nil {
		slog.ErrorContext(ctx, "failed to create template", "name", name, "error", err)
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return nil, exceptions.NewBadRequestError(errcode.TemplateAlreadyExists, nil)
		}
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return nil, exceptions.NewBadRequestError(errcode.TemplateAlreadyExists, nil)
		}
		return nil, exceptions.NewInternalServerError(errcode.FailToCreateTemplate)
	}

	if s.bus != nil {
		_ = eventbus.Publish(ctx, s.bus, api.TemplateCreatedTopic, template)
	}

	return template, nil
}

func (s *templateService) GetByID(ctx context.Context, id string) (*model.Template, error) {
	if strings.TrimSpace(id) == "" {
		return nil, exceptions.NewBadRequestError(response_code.InvalidParameter, nil)
	}

	template, err := s.repository.GetByID(ctx, id)
	if err != nil {
		slog.ErrorContext(ctx, "failed to get template by id", "id", id, "error", err)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, exceptions.NewNotFoundError(errcode.NotFoundTemplate)
		}
		return nil, exceptions.NewInternalServerError(errcode.FailToGetTemplate)
	}

	return template, nil
}

func (s *templateService) GetAll(ctx context.Context, req *dto.PaginationOptions) (*httpx.PaginatedData[model.Template], error) {
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
		slog.ErrorContext(ctx, "failed to list templates", "error", err)
		return nil, exceptions.NewInternalServerError(errcode.FailToGetTemplate)
	}

	res := httpx.NewPaginatedData(items, total, page, pageSize)
	return &res, nil
}

func (s *templateService) Delete(ctx context.Context, id string) error {
	if strings.TrimSpace(id) == "" {
		return exceptions.NewBadRequestError(response_code.InvalidParameter, nil)
	}

	if err := s.repository.Delete(ctx, id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return exceptions.NewNotFoundError(errcode.NotFoundTemplate)
		}
		return exceptions.NewInternalServerError(errcode.FailToDeleteTemplate)
	}

	return nil
}
