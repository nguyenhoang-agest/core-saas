package httpx

import (
	"encoding/json"
	"net/http"

	"github.com/example/core-backend/pkg/constants/response_code"
	pkgerrors "github.com/example/core-backend/pkg/errors"
	"github.com/example/core-backend/pkg/exceptions"
)

func JSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}

func Success[T any](w http.ResponseWriter, status int, data T) {
	JSON(w, status, NewResponse(data))
}

func Error(w http.ResponseWriter, err error) {
	if customErr, ok := err.(exceptions.CustomError); ok {
		JSON(w, customErr.StatusCode(), customErr.Json())
		return
	}

	if appErr, ok := pkgerrors.AsAppError(err); ok {
		JSON(w, appErr.HTTPStatusCode(), NewResponseError(string(appErr.Code), nil))
		return
	}

	JSON(w, http.StatusInternalServerError, NewResponseError(response_code.InternalServerGeneral.Code(), nil))
}

func DecodeJSON(r *http.Request, target any) error {
	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(target); err != nil {
		return exceptions.NewBadRequestError(response_code.BadRequestValidation, nil)
	}
	return nil
}
