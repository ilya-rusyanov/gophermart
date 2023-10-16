package handlers

import (
	"errors"
	"net/http"

	"github.com/ilya-rusyanov/gophermart/internal/entities"
)

type Logger interface {
	Error(...any)
}

type DefaultErrorHandler struct {
	log Logger
}

func NewDefaultErrorHandler(logger Logger) *DefaultErrorHandler {
	return &DefaultErrorHandler{
		log: logger,
	}
}

func (h *DefaultErrorHandler) Handle(rw http.ResponseWriter, err error) {
	var statusCode int
	switch {
	case errors.Is(err, errParsing):
		statusCode = http.StatusBadRequest
	case errors.Is(err, entities.ErrLoginConflict):
		statusCode = http.StatusConflict
	case errors.Is(err, entities.ErrLoginIncorrect):
		statusCode = http.StatusUnauthorized
	case errors.Is(err, entities.ErrUnauthorized):
		statusCode = http.StatusUnauthorized
	default:
		statusCode = http.StatusInternalServerError
	}

	http.Error(rw, err.Error(), statusCode)

	h.log.Error(err.Error())
}
