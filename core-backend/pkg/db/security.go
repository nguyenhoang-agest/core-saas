package db

import (
	"errors"
	"strings"

	"github.com/jackc/pgx/v5/pgconn"
)

const DefaultTenantMaxConns int32 = 10

const DefaultTenantSSLMode string = "require"

func ValidateDSN(dsn string) error {
	dsn = strings.TrimSpace(dsn)
	if dsn == "" {
		return errors.New("dsn is empty")
	}

	cfg, err := pgconn.ParseConfig(dsn)
	if err != nil {
		return err
	}

	host := strings.ToLower(strings.TrimSpace(cfg.Host))
	if host == "" {
		return errors.New("missing database host in dsn")
	}

	// Block loopback, local interfaces, and cloud metadata services to prevent SSRF
	blockedHosts := map[string]struct{}{
		"localhost":       {},
		"127.0.0.1":       {},
		"::1":             {},
		"0.0.0.0":         {},
		"169.254.169.254": {},
	}
	if _, blocked := blockedHosts[host]; blocked {
		return errors.New("database host is prohibited")
	}

	if strings.HasPrefix(host, "127.") {
		return errors.New("loopback database host is prohibited")
	}

	return nil
}
