package app

import "net/http"

type StatusCodeError struct {
	code int
}

func (e StatusCodeError) Error() string {
	return http.StatusText(e.code)
}

func (e StatusCodeError) StatusCode() int {
	return e.code
}
