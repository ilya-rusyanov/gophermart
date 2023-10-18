package main

type ErrPrintLogger interface {
	Error(...any)
}

func printErrors(logger ErrPrintLogger, ch <-chan error) {
	for err := range ch {
		logger.Error(err.Error())
	}
}
