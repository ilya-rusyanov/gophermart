package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func encodeJSON(
	rw http.ResponseWriter, errorHandler ErrorHandler, data any, status int,
) {
	rw.Header().Set("Content-Type", "application/json; charset=UTF-8")
	rw.WriteHeader(status)
	err := json.NewEncoder(rw).Encode(data)
	if err != nil {
		errorHandler(rw, fmt.Errorf("failed to encode JSON: %w", err))
		return
	}
}
