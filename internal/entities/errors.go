package entities

import "errors"

var (
	ErrNotFound       = errors.New("not found")
	ErrLoginConflict  = errors.New("login conflict")
	ErrLoginIncorrect = errors.New("wrong login or password")
	ErrUnauthorized   = errors.New("user is unauthorized")
)
