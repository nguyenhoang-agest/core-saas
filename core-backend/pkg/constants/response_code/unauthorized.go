package response_code

import "net/http"

type UnauthorizedCode int

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
)

func (c UnauthorizedCode) Code() string {
	return FormatCode(1, int(c))
}

func (c UnauthorizedCode) Error() string {
	return c.Code()
}

func (c UnauthorizedCode) HTTPStatus() int {
	return http.StatusUnauthorized
}
