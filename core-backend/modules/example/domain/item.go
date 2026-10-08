package domain

import (
	"context"
	"time"
)

// Item represents a tenant-scoped example entity.
type Item struct {
	ID        string    `json:"id"`
	TenantID  string    `json:"tenant_id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Repository defines data access operations for Item.
type Repository interface {
	Create(ctx context.Context, item *Item) error
	GetByID(ctx context.Context, id string) (*Item, error)
	List(ctx context.Context) ([]Item, error)
}
