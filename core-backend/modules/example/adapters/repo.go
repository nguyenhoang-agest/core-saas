package adapters

import (
	"context"
	"errors"
	"fmt"

	"github.com/example/core-backend/modules/example/domain"
	pkgerrors "github.com/example/core-backend/pkg/errors"
	"github.com/example/core-backend/pkg/tenantctx"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// PostgresRepository implements domain.Repository with pgx and tenant strategy enforcement.
type PostgresRepository struct {
	pool     *pgxpool.Pool
	strategy tenantctx.Strategy
}

// NewPostgresRepository constructs a new PostgresRepository.
func NewPostgresRepository(pool *pgxpool.Pool, strategy tenantctx.Strategy) *PostgresRepository {
	return &PostgresRepository{
		pool:     pool,
		strategy: strategy,
	}
}

// Create inserts a new item record.
func (r *PostgresRepository) Create(ctx context.Context, item *domain.Item) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("starting tx: %w", err)
	}
	defer func() { _ = tx.Rollback(ctx) }()

	if r.strategy != nil {
		if err := r.strategy.Apply(ctx, tx); err != nil {
			return err
		}
	}

	query := `
		INSERT INTO example_items (id, tenant_id, name, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
	`
	_, err = tx.Exec(ctx, query, item.ID, item.TenantID, item.Name, item.CreatedAt, item.UpdatedAt)
	if err != nil {
		return fmt.Errorf("inserting item: %w", err)
	}

	return tx.Commit(ctx)
}

// GetByID retrieves a single item scoped to the tenant.
func (r *PostgresRepository) GetByID(ctx context.Context, id string) (*domain.Item, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("starting tx: %w", err)
	}
	defer func() { _ = tx.Rollback(ctx) }()

	if r.strategy != nil {
		if err := r.strategy.Apply(ctx, tx); err != nil {
			return nil, err
		}
	}

	query := `
		SELECT id, tenant_id, name, created_at, updated_at
		FROM example_items
		WHERE id = $1
	`
	var item domain.Item
	err = tx.QueryRow(ctx, query, id).Scan(
		&item.ID,
		&item.TenantID,
		&item.Name,
		&item.CreatedAt,
		&item.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, pkgerrors.NotFound("item not found", err)
		}
		return nil, fmt.Errorf("querying item: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}
	return &item, nil
}

// List returns items scoped to the active tenant.
func (r *PostgresRepository) List(ctx context.Context) ([]domain.Item, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("starting tx: %w", err)
	}
	defer func() { _ = tx.Rollback(ctx) }()

	if r.strategy != nil {
		if err := r.strategy.Apply(ctx, tx); err != nil {
			return nil, err
		}
	}

	query := `
		SELECT id, tenant_id, name, created_at, updated_at
		FROM example_items
		ORDER BY created_at DESC
	`
	rows, err := tx.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("querying items: %w", err)
	}
	defer rows.Close()

	var items []domain.Item
	for rows.Next() {
		var item domain.Item
		if err := rows.Scan(&item.ID, &item.TenantID, &item.Name, &item.CreatedAt, &item.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scanning item row: %w", err)
		}
		items = append(items, item)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterating items: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}
	return items, nil
}
