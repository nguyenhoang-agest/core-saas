package httpx

import "github.com/example/core-backend/pkg/exceptions"

type FieldErrorDetail = exceptions.FieldErrorDetail
type ResponseError = exceptions.ResponseError

type Response struct {
	Status string `json:"status,omitempty"`
	Data   any    `json:"data,omitempty"`
}

func NewResponse(data any) Response {
	return Response{
		Status: "success",
		Data:   data,
	}
}

func NewResponseError(messageCode string, details []FieldErrorDetail) ResponseError {
	return ResponseError{
		Status:      "error",
		MessageCode: messageCode,
		ErrorDetail: details,
	}
}
