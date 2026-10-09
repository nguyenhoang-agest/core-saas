package httpx

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/example/core-backend/pkg/constants/response_code"
	pkgerrors "github.com/example/core-backend/pkg/errors"
)

func JSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}

func Success[T any](w http.ResponseWriter, status int, data T) {
	JSON(w, status, NewResponse(data))
}

func ErrorCode(w http.ResponseWriter, code response_code.IResponseCode) {
	JSON(w, code.HTTPStatus(), NewResponseError(code.Code(), nil))
}

func ErrorWithDetails(w http.ResponseWriter, code response_code.IResponseCode, details []FieldErrorDetail) {
	JSON(w, code.HTTPStatus(), NewResponseError(code.Code(), details))
}

func Error(w http.ResponseWriter, err error) {
	var respCode response_code.IResponseCode
	if errors.As(err, &respCode) {
		ErrorCode(w, respCode)
		return
	}

	if appErr, ok := pkgerrors.AsAppError(err); ok {
		JSON(w, appErr.HTTPStatusCode(), NewResponseError(string(appErr.Code), nil))
		return
	}

	ErrorCode(w, response_code.InternalServerGeneral)
}

func DecodeJSON(r *http.Request, target any) error {
	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(target); err != nil {
		return response_code.InvalidPayload
	}
	return nil
}
