package adapters

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/example/core-backend/modules/tenant/domain"
	"github.com/example/core-backend/modules/tenant/domain/exception"
	"github.com/example/core-backend/modules/tenant/repository"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

var _ repository.ITenantRepository = (*PostgresTenantRepository)(nil)

type PostgresTenantRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresTenantRepository(pool *pgxpool.Pool) *PostgresTenantRepository {
	return &PostgresTenantRepository{pool: pool}
}

func (r *PostgresTenantRepository) Create(ctx context.Context, tenant *domain.Tenant) error {
	if r.pool == nil {
		return fmt.Errorf("database pool is not initialized")
	}

	query := `
		INSERT INTO tenants (id, name, description, status, create_at, update_at, delete_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`
	_, err := r.pool.Exec(
		ctx,
		query,
		tenant.Id,
		tenant.Name,
		tenant.Description,
		string(tenant.Status),
		tenant.CreateAt,
		tenant.UpdateAt,
		tenant.DeleteAt,
	)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return exception.ErrAlreadyExists
		}
		return fmt.Errorf("inserting tenant: %w", err)
	}

	return nil
}

func (r *PostgresTenantRepository) GetTenantByID(ctx context.Context, id string) (*domain.Tenant, error) {
	if r.pool == nil {
		return nil, fmt.Errorf("database pool is not initialized")
	}

	query := `
		SELECT id, name, description, status, create_at, update_at, delete_at
		FROM tenants
		WHERE id = $1 AND delete_at IS NULL
	`
	var t domain.Tenant
	var statusStr string

	err := r.pool.QueryRow(ctx, query, id).Scan(
		&t.Id,
		&t.Name,
		&t.Description,
		&statusStr,
		&t.CreateAt,
		&t.UpdateAt,
		&t.DeleteAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, exception.ErrNotFound
		}
		return nil, fmt.Errorf("querying tenant by id: %w", err)
	}

	t.Status = domain.Status(statusStr)
	return &t, nil
}

func (r *PostgresTenantRepository) GetTenantByName(ctx context.Context, name string) (*domain.Tenant, error) {
	if r.pool == nil {
		return nil, fmt.Errorf("database pool is not initialized")
	}

	query := `
		SELECT id, name, description, status, create_at, update_at, delete_at
		FROM tenants
		WHERE name = $1 AND delete_at IS NULL
	`
	var t domain.Tenant
	var statusStr string

	err := r.pool.QueryRow(ctx, query, name).Scan(
		&t.Id,
		&t.Name,
		&t.Description,
		&statusStr,
		&t.CreateAt,
		&t.UpdateAt,
		&t.DeleteAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, exception.ErrNotFound
		}
		return nil, fmt.Errorf("querying tenant by name: %w", err)
	}

	t.Status = domain.Status(statusStr)
	return &t, nil
}

func (r *PostgresTenantRepository) Update(ctx context.Context, tenant *domain.Tenant) error {
	if r.pool == nil {
		return fmt.Errorf("database pool is not initialized")
	}

	tenant.UpdateAt = time.Now().UTC()

	query := `
		UPDATE tenants
		SET name = $2, description = $3, status = $4, update_at = $5
		WHERE id = $1 AND delete_at IS NULL
	`
	cmdTag, err := r.pool.Exec(
		ctx,
		query,
		tenant.Id,
		tenant.Name,
		tenant.Description,
		string(tenant.Status),
		tenant.UpdateAt,
	)
	if err != nil {
		return fmt.Errorf("updating tenant: %w", err)
	}

	if cmdTag.RowsAffected() == 0 {
		return exception.ErrNotFound
	}

	return nil
}

func (r *PostgresTenantRepository) Delete(ctx context.Context, id string) error {
	if r.pool == nil {
		return fmt.Errorf("database pool is not initialized")
	}

	now := time.Now().UTC()
	query := `
		UPDATE tenants
		SET delete_at = $2, update_at = $2
		WHERE id = $1 AND delete_at IS NULL
	`
	cmdTag, err := r.pool.Exec(ctx, query, id, now)
	if err != nil {
		return fmt.Errorf("deleting tenant: %w", err)
	}

	if cmdTag.RowsAffected() == 0 {
		return exception.ErrNotFound
	}

	return nil
}
