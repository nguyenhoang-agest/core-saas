package kernel

import (
	"context"

	"github.com/go-chi/chi/v5"
)

// Module defines the lifecycle and HTTP registration contract implemented by every module.
// Wiring is explicit; no reflection or dynamic scanning is permitted.
type Module interface {
	Name() string
	Register(r chi.Router) error
	Migrate(ctx context.Context) error
}
