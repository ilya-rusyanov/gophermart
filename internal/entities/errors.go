package entities

import "errors"

var (
	ErrNotFound      = errors.New("not found")
	ErrLoginConflict = errors.New("login conflict")
)
