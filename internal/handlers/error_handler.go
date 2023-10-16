package handlers

import "net/http"

type ErrorHandler func(http.ResponseWriter, error)
