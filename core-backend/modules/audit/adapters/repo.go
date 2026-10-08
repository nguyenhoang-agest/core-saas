package adapters

import (
	"context"

	"github.com/example/core-backend/modules/audit/domain"
	"github.com/example/core-backend/pkg/errors"
	"github.com/jackc/pgx/v5/pgxpool"
)

// PostgresAuditRepository implements domain.Repository (STUB).
type PostgresAuditRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresAuditRepository(pool *pgxpool.Pool) *PostgresAuditRepository {
	return &PostgresAuditRepository{pool: pool}
}

func (r *PostgresAuditRepository) Record(ctx context.Context, entry *domain.LogEntry) error {
	// STUB: intentionally omitted
	return nil
}

func (r *PostgresAuditRepository) ListByTenant(ctx context.Context, tenantID string) ([]domain.LogEntry, error) {
	// STUB: intentionally omitted
	return nil, errors.NotFound("no audit logs found", nil)
}
