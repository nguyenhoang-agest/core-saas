package httpx

import (
	"encoding/json"
	"net/http"

	"github.com/example/core-backend/pkg/errors"
)

// Response represents a standard JSON response wrapper.
type Response struct {
	Data any `json:"data,omitempty"`
}

// ErrorResponse represents a standardized JSON error response.
type ErrorResponse struct {
	Error ErrorDetail `json:"error"`
}

type ErrorDetail struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// JSON writes a JSON response with status code.
func JSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}

// Error formats and writes an error response mapped from an application error.
func Error(w http.ResponseWriter, err error) {
	if appErr, ok := errors.AsAppError(err); ok {
		JSON(w, appErr.HTTPStatusCode(), ErrorResponse{
			Error: ErrorDetail{
				Code:    string(appErr.Code),
				Message: appErr.Message,
			},
		})
		return
	}

	JSON(w, http.StatusInternalServerError, ErrorResponse{
		Error: ErrorDetail{
			Code:    string(errors.CodeInternal),
			Message: "internal server error",
		},
	})
}

// DecodeJSON decodes request JSON body into target.
func DecodeJSON(r *http.Request, target any) error {
	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(target); err != nil {
		return errors.InvalidArgument("invalid json request payload", err)
	}
	return nil
}
