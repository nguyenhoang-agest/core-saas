package httpx

import (
	"context"
	"log/slog"
	"net/http"
	"sync"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// HealthCheck is a function that reports component health.
type HealthCheck func(ctx context.Context) error

// Server wraps chi.Mux and provides health check registration.
type Server struct {
	Router  *chi.Mux
	logger  *slog.Logger
	mu      sync.RWMutex
	checkers map[string]HealthCheck
}

// NewServer creates a chi router configured with standard middlewares and /healthz.
func NewServer(logger *slog.Logger) *Server {
	r := chi.NewRouter()

	s := &Server{
		Router:   r,
		logger:   logger,
		checkers: make(map[string]HealthCheck),
	}

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(RequestLogger(logger))
	r.Use(middleware.Recoverer)

	r.Get("/healthz", s.handleHealthz)

	return s
}

// RegisterHealthCheck adds a named health check probe.
func (s *Server) RegisterHealthCheck(name string, check HealthCheck) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.checkers[name] = check
}

func (s *Server) handleHealthz(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
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

	JSON(w, statusCode, map[string]any{
		"status": status,
		"checks": details,
	})
}

// RequestLogger creates a middleware logging HTTP requests via slog.
func RequestLogger(logger *slog.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

			next.ServeHTTP(ww, r)

			logger.InfoContext(r.Context(), "http_request",
				slog.String("method", r.Method),
				slog.String("path", r.URL.Path),
				slog.Int("status", ww.Status()),
				slog.Int("bytes", ww.BytesWritten()),
				slog.Duration("duration", time.Since(start)),
				slog.String("request_id", middleware.GetReqID(r.Context())),
			)
		})
	}
}
