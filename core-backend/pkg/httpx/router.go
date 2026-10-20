package httpx

import (
	"context"
	"log/slog"
	"net/http"
	"sync"
	"time"

	"github.com/example/core-backend/pkg/i18n"
	"github.com/example/core-backend/pkg/tenantctx"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/recover"
)

type HealthCheck func(ctx context.Context) error

type Server struct {
	App      *fiber.App
	logger   *slog.Logger
	mu       sync.RWMutex
	checkers map[string]HealthCheck
}

func NewServer(logger *slog.Logger) *Server {
	s := &Server{
		logger:   logger,
		checkers: make(map[string]HealthCheck),
	}

	app := fiber.New(fiber.Config{
		ErrorHandler:    ErrorHandlerFnc,
		StructValidator: NewStructValidator(),
	})
	app.Use(recover.New())
	app.Use(tenantctx.FiberMiddleware())
	app.Use(i18n.LocaleMiddleware())
	app.Use(RequestLogger(logger))

	app.Get("/healthz", s.handleHealthz)

	s.App = app
	return s
}

func (s *Server) RegisterHealthCheck(name string, check HealthCheck) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.checkers[name] = check
}

func (s *Server) handleHealthz(c fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.Context(), 3*time.Second)
	defer cancel()

	s.mu.RLock()
	defer s.mu.RUnlock()

	status := "ok"
	statusCode := http.StatusOK
	details := make(map[string]string)

	for name, check := range s.checkers {
		if err := check(ctx); err != nil {
			status = "degraded"
			statusCode = http.StatusServiceUnavailable
			details[name] = err.Error()
		} else {
			details[name] = "ok"
		}
	}

	return c.Status(statusCode).JSON(map[string]any{
		"status": status,
		"checks": details,
	})
}

func RequestLogger(logger *slog.Logger) fiber.Handler {
	return func(c fiber.Ctx) error {
		start := time.Now()
		err := c.Next()
		duration := time.Since(start)

		if logger != nil {
			logger.InfoContext(c.Context(), "http_request",
				slog.String("method", c.Method()),
				slog.String("path", c.Path()),
				slog.Int("status", c.Response().StatusCode()),
				slog.Duration("duration", duration),
			)
		}

		return err
	}
}
