package response_code

import (
	"fmt"
	"net/http"
)

type SuccessCode int

const (
	_ SuccessCode = iota
	SuccessGeneral
	SuccessCreated
	SuccessUpdated
	SuccessDeleted
)

func (c SuccessCode) Code() string {
	return fmt.Sprintf("S00%04d", c)
}

func (c SuccessCode) Error() string {
	return c.Code()
}

func (c SuccessCode) HTTPStatus() int {
	if c == SuccessCreated {
		return http.StatusCreated
	}
	return http.StatusOK
}
