package repository

import (
	"context"
)

type IRepository[T any] interface {
	Create(ctx context.Context, item *T) error
	Update(ctx context.Context, item *T) error
	Delete(ctx context.Context, id string) error
}
