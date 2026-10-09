package response_code

import "net/http"

type UnprocessableEntityCode int

const (
	_ UnprocessableEntityCode = iota
	UnprocessableGeneral
	BusinessRuleViolation
	StateTransitionInvalid
)

func (c UnprocessableEntityCode) Code() string {
	return FormatCode(8, int(c))
}

func (c UnprocessableEntityCode) Error() string {
	return c.Code()
}

func (c UnprocessableEntityCode) HTTPStatus() int {
	return http.StatusUnprocessableEntity
}
