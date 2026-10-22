package response_code

import (
	"net/http"

	"github.com/example/core-backend/pkg/exceptions/err_type"
)

type ValidationCode int

// Validation response code: E05****
const (
	CommonError ValidationCode = iota
	RequiredField
	NotEmail
	NotUrl
	NotAlphanum
	NotAlpha
	NotNumeric
	NotNumber
	NotUuid
	NotIP
	NotHexColor
	NotOneOf
	MinLength
	MaxLength
	LengthExact
	NotEqual
	MustGreatThan
	MustGreatThanEqual
	MustLessThan
	MustLessThanEqual
	WrongDatetime
	NotEqualField
	// ONLY APPEND ERROR CODE
)

func (c ValidationCode) Code() string {
	return FormatCode(err_type.ValidationError, int(c))
}

func (c ValidationCode) Error() string {
	return c.Code()
}

func (c ValidationCode) HTTPStatus() int {
	return http.StatusBadRequest
}
