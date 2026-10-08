package app

import (
	"context"
	"time"

	"github.com/example/core-backend/modules/example/domain"
	"github.com/example/core-backend/pkg/errors"
	"github.com/example/core-backend/pkg/eventbus"
	"github.com/example/core-backend/pkg/tenantctx"
	"github.com/google/uuid"
)

const ItemCreatedTopic = "example.item.created"

// ItemCreatedEvent is published when an item is created.
type ItemCreatedEvent struct {
	ItemID    string    `json:"item_id"`
	TenantID  string    `json:"tenant_id"`
	Name      string    `json:"name"`
	Timestamp time.Time `json:"timestamp"`
}

// CreateItemInput contains parameters for creating an item.
type CreateItemInput struct {
	Name string `json:"name"`
}

// Service coordinates application logic for items.
type Service struct {
	repo domain.Repository
	bus  *eventbus.Bus
}

// NewService constructs an Item application service.
func NewService(repo domain.Repository, bus *eventbus.Bus) *Service {
	return &Service{
		repo: repo,
		bus:  bus,
	}
}

// Create validates input, extracts tenant context, stores item, and emits an event.
func (s *Service) Create(ctx context.Context, input CreateItemInput) (*domain.Item, error) {
	if input.Name == "" {
		return nil, errors.InvalidArgument("item name cannot be empty", nil)
	}

	tenantID, ok := tenantctx.FromContext(ctx)
	if !ok {
		return nil, errors.Unauthorized("missing tenant context", nil)
	}

	now := time.Now().UTC()
	item := &domain.Item{
		ID:        uuid.NewString(),
		TenantID:  tenantID,
		Name:      input.Name,
		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := s.repo.Create(ctx, item); err != nil {
		return nil, err
	}

	if s.bus != nil {
		_ = eventbus.Publish(ctx, s.bus, ItemCreatedTopic, ItemCreatedEvent{
			ItemID:    item.ID,
			TenantID:  item.TenantID,
			Name:      item.Name,
			Timestamp: now,
		})
	}

	return item, nil
}

// GetByID retrieves a single item by ID.
func (s *Service) GetByID(ctx context.Context, id string) (*domain.Item, error) {
	if id == "" {
		return nil, errors.InvalidArgument("id cannot be empty", nil)
	}
	return s.repo.GetByID(ctx, id)
}

// List returns all items within the active tenant scope.
func (s *Service) List(ctx context.Context) ([]domain.Item, error) {
	return s.repo.List(ctx)
}
