package app

import (
	"context"

	"github.com/example/core-backend/modules/audit/domain"
	"github.com/example/core-backend/pkg/errors"
	"github.com/example/core-backend/pkg/eventbus"
)

// Service provides audit logging functionality (STUB).
type Service struct {
	repo domain.Repository
	bus  *eventbus.Bus
}

func NewService(repo domain.Repository, bus *eventbus.Bus) *Service {
	return &Service{
		repo: repo,
		bus:  bus,
	}
}

func (s *Service) Record(ctx context.Context, entry *domain.LogEntry) error {
	if entry == nil {
		return errors.InvalidArgument("audit entry cannot be nil", nil)
	}
	return s.repo.Record(ctx, entry)
}

func (s *Service) ListByTenant(ctx context.Context, tenantID string) ([]domain.LogEntry, error) {
	if tenantID == "" {
		return nil, errors.InvalidArgument("tenant id is required", nil)
	}
	return s.repo.ListByTenant(ctx, tenantID)
}
