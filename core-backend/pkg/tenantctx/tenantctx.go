package tenantctx

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
)

type contextKey struct{}

var tenantKey = contextKey{}

// WithContext returns a new Context carrying tenantID.
func WithContext(ctx context.Context, tenantID string) context.Context {
	return context.WithValue(ctx, tenantKey, tenantID)
}

// FromContext extracts tenantID from Context.
func FromContext(ctx context.Context) (string, bool) {
	val := ctx.Value(tenantKey)
	if val == nil {
		return "", false
	}
	tenantID, ok := val.(string)
	return tenantID, ok && tenantID != ""
}

// Strategy defines how multi-tenancy is enforced on a database transaction.
//
// Documented Extension Points (not implemented in this starter kit):
//  1. SchemaPerTenantStrategy: switches search_path to the tenant schema
//     (e.g., SET LOCAL search_path = 'tenant_123, public')
//  2. DatabasePerTenantStrategy: dynamically routes connection pool lookups
//     by tenant ID to separate PostgreSQL databases.
type Strategy interface {
	Apply(ctx context.Context, tx pgx.Tx) error
}

// PostgresRLSStrategy enforces multi-tenancy via PostgreSQL Row-Level Security (RLS)
// and a discriminator tenant_id column. It sets a transaction-scoped local parameter
// `app.current_tenant_id` that RLS policies can evaluate using `current_setting('app.current_tenant_id')`.
type PostgresRLSStrategy struct{}

// NewPostgresRLSStrategy constructs a PostgresRLSStrategy.
func NewPostgresRLSStrategy() *PostgresRLSStrategy {
	return &PostgresRLSStrategy{}
}

// Apply sets the transaction-local `app.current_tenant_id` configuration.
func (s *PostgresRLSStrategy) Apply(ctx context.Context, tx pgx.Tx) error {
	tenantID, ok := FromContext(ctx)
	if !ok {
		return fmt.Errorf("tenant ID missing in context for RLS strategy")
	}

	// is_local=true ensures the setting is scoped strictly to the current transaction.
	_, err := tx.Exec(ctx, "SELECT set_config('app.current_tenant_id', $1, true)", tenantID)
	if err != nil {
		return fmt.Errorf("setting transaction rls tenant: %w", err)
	}

	return nil
}
