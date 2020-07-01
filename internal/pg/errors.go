package pg

import (
	"strings"

	"github.com/pkg/errors"
)

func uniqueViolation(err error) bool {
	if err == nil {
		return false
	}
	err = errors.Cause(err)
	errMsg := strings.ToLower(err.Error())
	return strings.Contains(errMsg, "unique constraint failed") ||
		strings.Contains(errMsg, "23505") ||
		strings.Contains(errMsg, "duplicate key") ||
		strings.Contains(errMsg, "duplicate entry") ||
		strings.Contains(errMsg, "error 1062")
}

func notFound(err error) bool {
	if err == nil {
		return false
	}
	err = errors.Cause(err)
	return strings.Contains(err.Error(), "not found")
}

func deadlock(err error) bool {
	if err == nil {
		return false
	}
	err = errors.Cause(err)
	return strings.Contains(err.Error(), "deadlock")
}

func accessDB(err error) error {
	return errors.Wrap(err, "access db")
}
