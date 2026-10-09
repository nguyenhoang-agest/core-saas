package response_code

import "fmt"

type IResponseCode interface {
	error
	Code() string
	HTTPStatus() int
}

func FormatCode(group int, code int) string {
	return fmt.Sprintf("E%02d%04d", group, code)
}
