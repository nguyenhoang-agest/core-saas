package exceptions

import (
	"fmt"

	"github.com/example/core-backend/pkg/constants/response_code"
)

type CustomError interface {
	error
	Json() map[string]any
	FormattedLog() string
	StatusCode() int
}

type FieldErrorDetail struct {
	Field       string `json:"field,omitempty"`
	DataError   any    `json:"data_error,omitempty"`
	MessageCode string `json:"message_code,omitempty"`
}

type ResponseError struct {
	Status      string             `json:"status"`
	MessageCode string             `json:"message_code"`
	ErrorDetail []FieldErrorDetail `json:"error_detail,omitempty"`
}

func NewResponseError[T response_code.ResponseCode](code T, errorType int, errorDetail []FieldErrorDetail) ResponseError {
	return ResponseError{
		Status:      "error",
		MessageCode: fmt.Sprintf("E%02d%04d", errorType, int(code)),
		ErrorDetail: errorDetail,
	}
}
