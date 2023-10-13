package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ilya-rusyanov/gophermart/internal/entities"
)

type Register struct {
	errHandler ErrorHandler
}

func NewRegister(errHandler ErrorHandler) *Register {
	return &Register{
		errHandler: errHandler,
	}
}

func (reg *Register) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	var creds entities.AuthCredentials
	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		reg.errHandler.Handle(rw,
			fmt.Errorf("failed to parse credentials from JSON: %w", errParsing),
		)
		return
	}
}
