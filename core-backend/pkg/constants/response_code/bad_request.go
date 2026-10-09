package response_code

import "net/http"

type BadRequestCode int

const (
	_ BadRequestCode = iota
	BadRequestGeneral
	InvalidPayload
	InvalidParameter
	MissingRequiredField
	TenantNameRequired
	TenantAlreadyExists
)

func (c BadRequestCode) Code() string {
	return FormatCode(4, int(c))
}

func (c BadRequestCode) Error() string {
	return c.Code()
}

func (c BadRequestCode) HTTPStatus() int {
	return http.StatusBadRequest
}
