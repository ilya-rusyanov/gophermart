package entities

import (
	"errors"
	"time"
)

var (
	ErrNotFound                    = errors.New("not found")
	ErrLoginConflict               = errors.New("login conflict")
	ErrLoginIncorrect              = errors.New("wrong login or password")
	ErrUnauthorized                = errors.New("user is unauthorized")
	ErrInvalidOrder                = errors.New("invalid order ID")
	ErrAlreadyUploaded             = errors.New("order already uploaded")
	ErrAlreadyUploadedOtherUser    = errors.New("order was alredy uploaded by other user")
	ErrAccrualOrderIsNotRegistered = errors.New("order is not registered in accrual")
)

type AccrualTooManyRequestsError struct {
	err    error
	Period time.Duration
}

func (e *AccrualTooManyRequestsError) Unwrap() error {
	return e.err
}

func (e *AccrualTooManyRequestsError) Error() string {
	return e.err.Error()
}
