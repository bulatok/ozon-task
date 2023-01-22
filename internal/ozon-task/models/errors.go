package models

import (
	"net/http"
)

type CommonError struct {
	StatusCode int
	Err        string
}

func NewCommonErr(err string, statusCode int) CommonError {
	return CommonError{
		Err:        err,
		StatusCode: statusCode,
	}
}

func (c CommonError) Error() string {
	return c.Err
}

var (
	ErrNotFound       = NewCommonErr("not found", http.StatusGone)
	ErrInternalServer = NewCommonErr("internal server error", http.StatusInternalServerError)
)
