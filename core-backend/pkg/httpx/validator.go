package httpx

import (
	"log/slog"
	"reflect"
	"strings"

	"github.com/example/core-backend/pkg/helpers"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
)

var _ fiber.StructValidator = (*structValidator)(nil)

type structValidator struct {
	validate *validator.Validate
}

func NewStructValidator() fiber.StructValidator {
	v := validator.New()
	helpers.RegisterCustomValidations(v)

	v.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name, _, _ := strings.Cut(fld.Tag.Get("json"), ",")
		if name == "-" {
			return ""
		}
		return name
	})

	return &structValidator{validate: v}
}

func (sv *structValidator) Validate(out any) error {
	if err := sv.validate.Struct(out); err != nil {
		slog.Error("validation error", "error", err)
		return err
	}
	return nil
}
