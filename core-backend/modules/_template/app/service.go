package app

import (
	"context"

	"github.com/example/core-backend/modules/_template/domain"
	"github.com/example/core-backend/pkg/errors"
	"github.com/example/core-backend/pkg/eventbus"
)

// Service provides use cases for Entity.
type Service struct {
	repo domain.Repository
	bus  *eventbus.Bus
}

// NewService constructs an application Service.
func NewService(repo domain.Repository, bus *eventbus.Bus) *Service {
	return &Service{
		repo: repo,
		bus:  bus,
	}
}

// GetByID fetches an entity by ID.
func (s *Service) GetByID(ctx context.Context, id string) (*domain.Entity, error) {
	if id == "" {
		return nil, errors.InvalidArgument("id cannot be empty", nil)
	}
	return s.repo.GetByID(ctx, id)
}
