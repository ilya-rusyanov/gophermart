package main

import (
	"context"
)

type ErrPrintLogger interface {
	Error(...any)
}

func printErrors(ctx context.Context, logger ErrPrintLogger, ch <-chan error) {
	for {
		select {
		case <-ctx.Done():
			return
		case err := <-ch:
			logger.Error(err.Error())
		}
	}
}
