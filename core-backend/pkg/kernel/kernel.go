package kernel

import (
	"context"

	"github.com/gofiber/fiber/v3"
)

type Module interface {
	Name() string
	Register(r fiber.Router) error
	Migrate(ctx context.Context) error
}
