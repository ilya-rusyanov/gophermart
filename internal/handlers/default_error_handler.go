package handlers

import (
	"errors"
	"net/http"
)

type Logger interface {
	Info(...any)
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
	}
}
