package httpx

import (
	"errors"
	"log/slog"

	"github.com/example/core-backend/pkg/constants/response_code"
	"github.com/example/core-backend/pkg/exceptions"
	"github.com/gofiber/fiber/v3"
)

// ErrorHandlerFnc is the centralized Fiber error handler formatting CustomError into JSON responses.
func ErrorHandlerFnc(ctx fiber.Ctx, err error) error {
	var customErr exceptions.CustomError
	if errors.As(err, &customErr) {
		if logs := customErr.FormattedLog(); logs != "" {
			slog.ErrorContext(ctx.Context(), logs)
		}
		return ctx.Status(customErr.StatusCode()).JSON(customErr.Json())
	}

	var fiberErr *fiber.Error
	if errors.As(err, &fiberErr) {
		return ctx.Status(fiberErr.Code).JSON(fiber.Map{
			"status":  "error",
			"message": fiberErr.Message,
		})
	}

	slog.ErrorContext(ctx.Context(), "internal server error", "error", err)

	internalErr := exceptions.NewInternalServerError(response_code.InternalServerGeneral)
	return ctx.Status(fiber.StatusInternalServerError).JSON(internalErr.Json())
}
