package response_code

import (
	"net/http"

	"github.com/example/core-backend/pkg/exceptions/err_type"
)

type UnauthorizedCode int

// Unauthorized response code: E01****
const (
	_ UnauthorizedCode = iota
	InvalidToken
	AccessTokenExpired
	RefreshTokenExpired
	WrongCredentials
	WrongJWTIssuer
	OTPExpired
	NotFoundSession
	WrongOTP
	VerifyTokenExpired
	NotFoundPersonalToken
	// ONLY APPEND ERROR CODE
)

func (c UnauthorizedCode) Code() string {
	return FormatCode(err_type.UnauthorizedError, int(c))
}

func (c UnauthorizedCode) Error() string {
	return c.Code()
}

func (c UnauthorizedCode) HTTPStatus() int {
	return http.StatusUnauthorized
}
