package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ilya-rusyanov/gophermart/internal/entities"
)

type RegisterUsecase interface {
	Register(context.Context, entities.AuthCredentials) (entities.AuthToken, error)
}

type Register struct {
	errHandler ErrorHandler
	usecase    RegisterUsecase
}

func NewRegister(usecase RegisterUsecase, errHandler ErrorHandler) *Register {
	return &Register{
		errHandler: errHandler,
		usecase:    usecase,
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
	token, err := reg.usecase.Register(r.Context(), creds)
	if err != nil {
		reg.errHandler.Handle(rw, fmt.Errorf("usecase failure: %w", err))
		return
	}
	processAuthToken(rw, token)
	rw.WriteHeader(http.StatusOK)
}
