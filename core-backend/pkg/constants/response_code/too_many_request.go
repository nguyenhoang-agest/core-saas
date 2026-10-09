package response_code

import "net/http"

type TooManyRequestCode int

const (
	_ TooManyRequestCode = iota
	RateLimitExceeded
	TooManyAttempts
	QuotaExceeded
)

func (c TooManyRequestCode) Code() string {
	return FormatCode(7, int(c))
}

func (c TooManyRequestCode) Error() string {
	return c.Code()
}

func (c TooManyRequestCode) HTTPStatus() int {
	return http.StatusTooManyRequests
}
