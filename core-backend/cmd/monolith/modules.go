package main

import (
	"github.com/example/core-backend/modules/audit"
	"github.com/example/core-backend/modules/example"
	"github.com/example/core-backend/modules/iam"
	"github.com/example/core-backend/modules/notification"
	"github.com/example/core-backend/modules/tenant"
	"github.com/example/core-backend/pkg/eventbus"
	"github.com/example/core-backend/pkg/kernel"
	"github.com/jackc/pgx/v5/pgxpool"
)

// enabledModules returns the list of all active modules wired into the monolith.
// When adding a new module via scripts/new-module.sh, it is registered here.
func enabledModules(pool *pgxpool.Pool, bus *eventbus.Bus) []kernel.Module {
	return []kernel.Module{
		example.New(pool, bus),
		iam.New(pool, bus),
		tenant.New(pool, bus),
		notification.New(pool, bus),
		audit.New(pool, bus),
	}
}
