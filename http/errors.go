package http

import (
	"net/http"

	"github.com/pkg/errors"
	"github.com/rpoletaev/huskyjam/internal"
)

type Error struct {
	Code    int
	Message string
}

func parseInternalError(err error) Error {
	switch errors.Cause(err) {
	case internal.ErrAlreadyExists:
		return Error{Message: err.Error(), Code: http.StatusBadRequest}
	case internal.ErrNotFound:
		return Error{Message: err.Error(), Code: http.StatusNotFound}
	default:
		return Error{Message: "internal error", Code: http.StatusInternalServerError}
	}
}
