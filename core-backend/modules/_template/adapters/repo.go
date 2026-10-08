package adapters

import (
	"context"
	"fmt"

	"github.com/example/core-backend/modules/_template/domain"
	"github.com/example/core-backend/pkg/errors"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Repository implements domain.Repository with pgx.
type Repository struct {
	pool *pgxpool.Pool
}

// NewRepository constructs a new Repository.
func NewRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{pool: pool}
}

// Create inserts a new entity.
func (r *Repository) Create(ctx context.Context, e *domain.Entity) error {
	// TODO: implement SQL query using pool
	return fmt.Errorf("not implemented")
}

// GetByID finds an entity by ID.
func (r *Repository) GetByID(ctx context.Context, id string) (*domain.Entity, error) {
	// TODO: implement SQL query using pool
	return nil, errors.NotFound("entity not found", nil)
}
