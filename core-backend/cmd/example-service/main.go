package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/example/core-backend/modules/example"
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

	// Service-specific port override default
	if os.Getenv("PORT") == "" {
		cfg.Server.Port = "8081"
	}

	log := logger.New(logger.Options{
		Level:       cfg.Logger.Level,
		ServiceName: "example-service",
		Environment: cfg.Server.Env,
	})
	slog.SetDefault(log)

	log.Info("starting standalone example service (microservice extraction demonstration)",
		"port", cfg.Server.Port,
		"env", cfg.Server.Env,
	)

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

	// In a standalone microservice, eventbus can be backed by Kafka/RabbitMQ or an in-process bus.
	bus := eventbus.New()

	// Instantiate ONLY example module
	exampleMod := example.New(pool, bus)

	if pool != nil {
		log.Info("running migrations for module", "module", exampleMod.Name())
		if err := exampleMod.Migrate(ctx); err != nil {
			log.Error("module migration failed", "module", exampleMod.Name(), "error", err)
			os.Exit(1)
		}
	}

	server := httpx.NewServer(log)
	if pool != nil {
		server.RegisterHealthCheck("database", func(ctx context.Context) error {
			return db.Ping(ctx, pool)
		})
	}

	if err := exampleMod.Register(server.App); err != nil {
		log.Error("failed to register module routes", "module", exampleMod.Name(), "error", err)
		os.Exit(1)
	}


	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := server.App.Listen(":" + cfg.Server.Port); err != nil {
			log.Error("http server failed", "error", err)
		}
	}()

	log.Info("example-service listening", "port", cfg.Server.Port)
	<-sigChan
	log.Info("shutting down gracefully...")

	if err := server.App.ShutdownWithTimeout(10 * time.Second); err != nil {
		log.Error("server forced to shutdown", "error", err)
	}

	log.Info("example-service exited cleanly")
}
