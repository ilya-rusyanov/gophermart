package handlers

import "net/http"

type ErrorHandler interface {
	Handle(http.ResponseWriter, error)
}
