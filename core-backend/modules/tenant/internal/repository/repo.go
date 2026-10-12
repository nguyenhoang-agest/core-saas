package repository

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/example/core-backend/modules/tenant/internal/model"
	"github.com/example/core-backend/modules/tenant/internal/service"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

var _ service.Repository = (*PostgresTenantRepository)(nil)

type PostgresTenantRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresTenantRepository(pool *pgxpool.Pool) *PostgresTenantRepository {
	return &PostgresTenantRepository{pool: pool}
}

func (r *PostgresTenantRepository) Create(ctx context.Context, tenant *model.Tenant) error {
	if r.pool == nil {
		return fmt.Errorf("database pool is not initialized")
	}

	if tenant.Id == "" {
		tenant.Id = uuid.NewString()
	}
	now := time.Now().UTC()
	if tenant.CreateAt.IsZero() {
		tenant.CreateAt = now
	}
	if tenant.UpdateAt.IsZero() {
		tenant.UpdateAt = now
	}
	if tenant.IsolateMode == "" {
		tenant.IsolateMode = model.IsolateModeShared
	}

	var dbConfigJSON []byte
	if tenant.DBConfig != nil {
		var err error
		dbConfigJSON, err = json.Marshal(tenant.DBConfig)
		if err != nil {
			return fmt.Errorf("marshaling db config: %w", err)
		}
	}

	query := `
		INSERT INTO tenants (id, name, description, status, isolate_mode, db_config, create_at, update_at, delete_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`
	_, err := r.pool.Exec(
		ctx,
		query,
		tenant.Id,
		tenant.Name,
		tenant.Description,
		string(tenant.Status),
		string(tenant.IsolateMode),
		dbConfigJSON,
		tenant.CreateAt,
		tenant.UpdateAt,
		tenant.DeleteAt,
	)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return model.ErrAlreadyExists
		}
		return fmt.Errorf("inserting tenant: %w", err)
	}

	return nil
}

func (r *PostgresTenantRepository) GetTenantByID(ctx context.Context, id string) (*model.Tenant, error) {
	if r.pool == nil {
		return nil, fmt.Errorf("database pool is not initialized")
	}

	query := `
		SELECT id, name, description, status, isolate_mode, db_config, create_at, update_at, delete_at
		FROM tenants
		WHERE id = $1 AND delete_at IS NULL
	`
	var t model.Tenant
	var statusStr string
	var isolateModeStr string
	var dbConfigJSON []byte

	err := r.pool.QueryRow(ctx, query, id).Scan(
		&t.Id,
		&t.Name,
		&t.Description,
		&statusStr,
		&isolateModeStr,
		&dbConfigJSON,
		&t.CreateAt,
		&t.UpdateAt,
		&t.DeleteAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, model.ErrNotFound
		}
		return nil, fmt.Errorf("querying tenant by id: %w", err)
	}

	t.Status = model.Status(statusStr)
	if isolateModeStr != "" {
		t.IsolateMode = model.IsolateMode(isolateModeStr)
	} else {
		t.IsolateMode = model.IsolateModeShared
	}

	if len(dbConfigJSON) > 0 {
		var dbCfg model.TenantDBConfig
		if err := json.Unmarshal(dbConfigJSON, &dbCfg); err == nil {
			t.DBConfig = &dbCfg
		}
	}

	return &t, nil
}

func (r *PostgresTenantRepository) GetTenantByName(ctx context.Context, name string) (*model.Tenant, error) {
	if r.pool == nil {
		return nil, fmt.Errorf("database pool is not initialized")
	}

	query := `
		SELECT id, name, description, status, isolate_mode, db_config, create_at, update_at, delete_at
		FROM tenants
		WHERE name = $1 AND delete_at IS NULL
	`
	var t model.Tenant
	var statusStr string
	var isolateModeStr string
	var dbConfigJSON []byte

	err := r.pool.QueryRow(ctx, query, name).Scan(
		&t.Id,
		&t.Name,
		&t.Description,
		&statusStr,
		&isolateModeStr,
		&dbConfigJSON,
		&t.CreateAt,
		&t.UpdateAt,
		&t.DeleteAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, model.ErrNotFound
		}
		return nil, fmt.Errorf("querying tenant by name: %w", err)
	}

	t.Status = model.Status(statusStr)
	if isolateModeStr != "" {
		t.IsolateMode = model.IsolateMode(isolateModeStr)
	} else {
		t.IsolateMode = model.IsolateModeShared
	}

	if len(dbConfigJSON) > 0 {
		var dbCfg model.TenantDBConfig
		if err := json.Unmarshal(dbConfigJSON, &dbCfg); err == nil {
			t.DBConfig = &dbCfg
		}
	}

	return &t, nil
}

func (r *PostgresTenantRepository) Update(ctx context.Context, tenant *model.Tenant) error {
	if r.pool == nil {
		return fmt.Errorf("database pool is not initialized")
	}

	tenant.UpdateAt = time.Now().UTC()

	var dbConfigJSON []byte
	if tenant.DBConfig != nil {
		var err error
		dbConfigJSON, err = json.Marshal(tenant.DBConfig)
		if err != nil {
			return fmt.Errorf("marshaling db config: %w", err)
		}
	}

	query := `
		UPDATE tenants
		SET name = $2, description = $3, status = $4, isolate_mode = $5, db_config = $6, update_at = $7
		WHERE id = $1 AND delete_at IS NULL
	`
	cmdTag, err := r.pool.Exec(
		ctx,
		query,
		tenant.Id,
		tenant.Name,
		tenant.Description,
		string(tenant.Status),
		string(tenant.IsolateMode),
		dbConfigJSON,
		tenant.UpdateAt,
	)
	if err != nil {
		return fmt.Errorf("updating tenant: %w", err)
	}

	if cmdTag.RowsAffected() == 0 {
		return model.ErrNotFound
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
		return model.ErrNotFound
	}

	return nil
}
