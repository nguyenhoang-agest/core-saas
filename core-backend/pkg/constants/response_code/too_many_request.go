package response_code

import (
	"net/http"

	"github.com/example/core-backend/pkg/exceptions/err_type"
)

type TooManyRequestCode int

// TooManyRequest response code: E07****
const (
	_ TooManyRequestCode = iota
	RateLimitExceeded
	TooManyAttempts
	QuotaExceeded
	// ONLY APPEND ERROR CODE
)

func (c TooManyRequestCode) Code() string {
	return FormatCode(err_type.TooManyRequestError, int(c))
}

func (c TooManyRequestCode) Error() string {
	return c.Code()
}

func (c TooManyRequestCode) HTTPStatus() int {
	return http.StatusTooManyRequests
}
