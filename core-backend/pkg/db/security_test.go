package db_test

import (
	"testing"

	"github.com/example/core-backend/pkg/db"
)

func TestValidateDSN_Valid(t *testing.T) {
	validDSNs := []string{
		"postgres://user:pass@db.internal.example.com:5432/app",
		"postgresql://app_user:secret@rds-postgres.corp:5432/tenant_123",
		"postgres://tenant_db:5432/app",
	}

	for _, dsn := range validDSNs {
		if err := db.ValidateDSN(dsn); err != nil {
			t.Errorf("expected DSN %q to be valid, got error: %v", dsn, err)
		}
	}
}

func TestValidateDSN_BlockedHosts(t *testing.T) {
	blockedDSNs := []string{
		"postgres://user:pass@localhost:5432/app",
		"postgres://user:pass@127.0.0.1:5432/app",
		"postgres://user:pass@127.0.0.2:5432/app",
		"postgres://user:pass@0.0.0.0:5432/app",
		"postgres://user:pass@[::1]:5432/app",
		"postgres://user:pass@169.254.169.254:5432/app",
	}

	for _, dsn := range blockedDSNs {
		if err := db.ValidateDSN(dsn); err == nil {
			t.Errorf("expected DSN %q to be blocked, but got nil error", dsn)
		}
	}
}

func TestValidateDSN_InvalidFormat(t *testing.T) {
	invalidDSNs := []string{
		"",
		"not_a_dsn",
		"://invalid",
	}

	for _, dsn := range invalidDSNs {
		if err := db.ValidateDSN(dsn); err == nil {
			t.Errorf("expected invalid DSN %q to fail parsing, got nil error", dsn)
		}
	}
}
