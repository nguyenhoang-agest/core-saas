package httpx

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
