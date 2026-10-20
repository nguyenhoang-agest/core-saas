package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/example/core-backend/modules/audit"
	"github.com/example/core-backend/modules/example"
	"github.com/example/core-backend/modules/notification"
	"github.com/example/core-backend/modules/tenant"
	"github.com/example/core-backend/modules/user"
	"github.com/example/core-backend/pkg/config"
	"github.com/example/core-backend/pkg/db"
	"github.com/example/core-backend/pkg/eventbus"
	"github.com/example/core-backend/pkg/kernel"
)

func main() {
	cfg, err := config.Load("config.yaml")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load configuration: %v\n", err)
		os.Exit(1)
	}

	ctx := context.Background()
	pool, err := db.Connect(ctx, cfg.Database.URL, cfg.Database.MaxConns)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer pool.Close()

	bus := eventbus.New()

	modules := []kernel.Module{
		tenant.New(pool, bus),
		user.New(pool, bus),
		notification.New(pool, bus),
		audit.New(pool, bus),
		example.New(pool, bus),
	}

	for _, m := range modules {
		slog.Info("migrating module", "module", m.Name())
		if err := m.Migrate(ctx); err != nil {
			slog.Error("migration failed", "module", m.Name(), "error", err)
			os.Exit(1)
		}
	}

	slog.Info("all migrations applied successfully")
}
