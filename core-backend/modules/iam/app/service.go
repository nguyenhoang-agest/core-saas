package app

import (
	"context"

	"github.com/example/core-backend/modules/iam/domain"
	"github.com/example/core-backend/pkg/errors"
	"github.com/example/core-backend/pkg/eventbus"
)

// Service provides IAM authentication and user management use-cases.
type Service struct {
	repo   domain.Repository
	hasher domain.PasswordHasher
	tokens domain.TokenGenerator
	bus    *eventbus.Bus
}

// NewService constructs the IAM application service.
func NewService(
	repo domain.Repository,
	hasher domain.PasswordHasher,
	tokens domain.TokenGenerator,
	bus *eventbus.Bus,
) *Service {
	return &Service{
		repo:   repo,
		hasher: hasher,
		tokens: tokens,
		bus:    bus,
	}
}

// Register registers a new user (STUB).
func (s *Service) Register(ctx context.Context, email, password string) (*domain.User, error) {
	// STUB: business logic intentionally omitted; see README.md TODOs
	return nil, errors.Internal("iam.Register is not implemented", nil)
}

// Login authenticates a user and returns a token (STUB).
func (s *Service) Login(ctx context.Context, email, password string) (string, error) {
	// STUB: business logic intentionally omitted; see README.md TODOs
	return "", errors.Internal("iam.Login is not implemented", nil)
}
