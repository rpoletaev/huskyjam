package internal

import "errors"

var (
	ErrNotFound      = errors.New("not found")
	ErrInternalError = errors.New("internal server error")
	ErrAlreadyExists = errors.New("already exists")
	ErrBadRequest    = errors.New("bad request")
)
