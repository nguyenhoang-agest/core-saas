package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/example/core-backend/pkg/config"
	"github.com/example/core-backend/pkg/db"
	"github.com/example/core-backend/pkg/eventbus"
	"github.com/example/core-backend/pkg/httpx"
	"github.com/example/core-backend/pkg/logger"
)

func main() {
	cfg, err := config.Load("config.yaml")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load configuration: %v\n", err)
		os.Exit(1)
	}

	log := logger.New(logger.Options{
		Level:       cfg.Logger.Level,
		ServiceName: "monolith",
		Environment: cfg.Server.Env,
	})
	slog.SetDefault(log)

	log.Info("starting saas starter monolith", "port", cfg.Server.Port, "env", cfg.Server.Env)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Connect Database
	pool, err := db.Connect(ctx, cfg.Database.URL, cfg.Database.MaxConns)
	if err != nil {
		log.Warn("database connection failed (proceeding for offline/scaffold testing)", "error", err)
	} else {
		defer pool.Close()
		log.Info("connected to database pool")
	}

	// Shared in-process event bus
	bus := eventbus.New()

	// Register enabled modules explicitly from modules.go
	modules := enabledModules(pool, bus)

	// Run migrations if database is connected
	if pool != nil {
		for _, m := range modules {
			log.Info("running migrations for module", "module", m.Name())
			if err := m.Migrate(ctx); err != nil {
				log.Error("module migration failed", "module", m.Name(), "error", err)
				os.Exit(1)
			}
		}
	}

	// Setup HTTP server
	server := httpx.NewServer(log)
	if pool != nil {
		server.RegisterHealthCheck("database", func(ctx context.Context) error {
			return db.Ping(ctx, pool)
		})
	}

	// Register module routes
	for _, m := range modules {
		log.Info("registering routes for module", "module", m.Name())
		if err := m.Register(server.App); err != nil {
			log.Error("failed to register module routes", "module", m.Name(), "error", err)
			os.Exit(1)
		}
	}

	// Graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := server.App.Listen(":" + cfg.Server.Port); err != nil {
			log.Error("http server failed", "error", err)
		}
	}()

	log.Info("server listening", "port", cfg.Server.Port)
	<-sigChan
	log.Info("shutting down gracefully...")

	if err := server.App.ShutdownWithTimeout(10 * time.Second); err != nil {
		log.Error("server forced to shutdown", "error", err)
	}

	log.Info("server exited cleanly")
}
