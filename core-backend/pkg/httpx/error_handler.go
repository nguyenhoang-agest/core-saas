package httpx

import (
	"errors"
	"log/slog"

	"github.com/example/core-backend/pkg/constants/response_code"
	"github.com/example/core-backend/pkg/exceptions"
	"github.com/example/core-backend/pkg/i18n"
	_ "github.com/example/core-backend/pkg/i18n/locales"
	"github.com/gofiber/fiber/v3"
)

func ErrorHandlerFnc(ctx fiber.Ctx, err error) error {
	var customErr exceptions.CustomError
	if errors.As(err, &customErr) {
		if logs := customErr.FormattedLog(); logs != "" {
			slog.ErrorContext(ctx.Context(), logs)
		}

		respMap := customErr.Json()
		msgCode := customErr.GetMessageCode()
		if _, hasMsg := respMap["message"]; !hasMsg || respMap["message"] == "" {
			if localizedMsg := i18n.T(ctx.Context(), msgCode); localizedMsg != "" {
				respMap["message"] = localizedMsg
			}
		}

		// Also localize each item in error_detail if present
		if details, ok := respMap["error_detail"].([]exceptions.FieldErrorDetail); ok {
			localizedDetails := make([]exceptions.FieldErrorDetail, len(details))
			for i, d := range details {
				if d.Message == "" {
					d.Message = i18n.T(ctx.Context(), d.MessageCode)
				}
				localizedDetails[i] = d
			}
			respMap["error_detail"] = localizedDetails
		}

		return ctx.Status(customErr.StatusCode()).JSON(respMap)
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
	respMap := internalErr.Json()
	if localizedMsg := i18n.T(ctx.Context(), internalErr.GetMessageCode()); localizedMsg != "" {
		respMap["message"] = localizedMsg
	}
	return ctx.Status(fiber.StatusInternalServerError).JSON(respMap)
}
