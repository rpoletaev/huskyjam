package auth

import "errors"

// Main token errors
var (
	ErrUnauthorized = errors.New("unauthorized")
	ErrInvalidToken = errors.New("token is invalid")
	ErrTokenExpired = errors.New("token is expired")
)
