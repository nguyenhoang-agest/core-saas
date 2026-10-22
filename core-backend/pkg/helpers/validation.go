package helpers

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/example/core-backend/pkg/constants/response_code"
	"github.com/example/core-backend/pkg/exceptions"
	"github.com/example/core-backend/pkg/exceptions/err_type"
	"github.com/go-playground/validator/v10"
)

func RegisterCustomValidations(validate *validator.Validate) {
	_ = validate
}

func GetValidationErrorCode(tag string) response_code.ValidationCode {
	switch tag {
	case "required":
		return response_code.RequiredField
	case "email":
		return response_code.NotEmail
	case "url":
		return response_code.NotUrl
	case "alphanum":
		return response_code.NotAlphanum
	case "alpha":
		return response_code.NotAlpha
	case "numeric":
		return response_code.NotNumeric
	case "number":
		return response_code.NotNumber
	case "uuid", "uuid3", "uuid4", "uuid5":
		return response_code.NotUuid
	case "min":
		return response_code.MinLength
	case "max":
		return response_code.MaxLength
	case "len":
		return response_code.LengthExact
	case "ne":
		return response_code.NotEqual
	case "eqfield":
		return response_code.NotEqualField
	case "gt":
		return response_code.MustGreatThan
	case "gte":
		return response_code.MustGreatThanEqual
	case "lt":
		return response_code.MustLessThan
	case "lte":
		return response_code.MustLessThanEqual
	case "datetime":
		return response_code.WrongDatetime
	case "ip", "ipv4", "ipv6":
		return response_code.NotIP
	case "hexcolor":
		return response_code.NotHexColor
	case "oneof":
		return response_code.NotOneOf
	default:
		return response_code.CommonError
	}
}

func FormatValidationErrors(
	verrs validator.ValidationErrors,
	structType reflect.Type,
) []exceptions.FieldErrorDetail {
	fieldErrors := make([]exceptions.FieldErrorDetail, 0, len(verrs))

	if structType != nil && structType.Kind() == reflect.Pointer {
		structType = structType.Elem()
	}

	for _, e := range verrs {
		jsonTag := e.Field()
		validationErrCode := GetValidationErrorCode(e.Tag())
		messageCodeString := fmt.Sprintf("E%02d%04d", err_type.ValidationError, validationErrCode)

		fieldErrors = append(fieldErrors, exceptions.FieldErrorDetail{
			Field:       jsonTag,
			MessageCode: messageCodeString,
		})
	}

	return fieldErrors
}

func ValidationErrorFromBind(err error, request any) error {
	var verrs validator.ValidationErrors
	if errors.As(err, &verrs) {
		return exceptions.NewValidationError(
			response_code.CommonError,
			FormatValidationErrors(verrs, reflect.TypeOf(request)),
		)
	}

	return exceptions.NewBadRequestError(response_code.BadRequestValidation, nil)
}
