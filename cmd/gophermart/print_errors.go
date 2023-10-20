package main

import (
	"context"
	"sync"
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

func fanInErrors(ctx context.Context, inputs ...<-chan error) <-chan error {
	result := make(chan error)

	var wg sync.WaitGroup

	for _, ch := range inputs {
		wg.Add(1)
		go func(ch <-chan error) {
			defer wg.Done()

			for err := range ch {
				select {
				case <-ctx.Done():
					return
				case result <- err:
				}
			}
		}(ch)
	}

	go func() {
		wg.Wait()

		close(result)
	}()

	return result
}
