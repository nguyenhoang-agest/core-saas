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
	GetMessageCode() string
}

type FieldErrorDetail struct {
	Field       string `json:"field,omitempty"`
	DataError   any    `json:"data_error,omitempty"`
	MessageCode string `json:"message_code,omitempty"`
	Message     string `json:"message,omitempty"`
}

type ResponseError struct {
	Status      string             `json:"status"`
	MessageCode string             `json:"message_code"`
	Message     string             `json:"message,omitempty"`
	ErrorDetail []FieldErrorDetail `json:"error_detail,omitempty"`
}

func (r ResponseError) GetMessageCode() string {
	return r.MessageCode
}

func NewResponseError[T response_code.ResponseCode](code T, errorType int, errorDetail []FieldErrorDetail) ResponseError {
	return ResponseError{
		Status:      "error",
		MessageCode: fmt.Sprintf("E%02d%04d", errorType, int(code)),
		ErrorDetail: errorDetail,
	}
}
