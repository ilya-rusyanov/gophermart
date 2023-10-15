package handlers

import (
	"errors"
	"net/http"
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
	switch {
	case errors.Is(err, errParsing):
		http.Error(rw, err.Error(), http.StatusBadRequest)
	default:
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	}
	h.log.Error(err.Error())
}
