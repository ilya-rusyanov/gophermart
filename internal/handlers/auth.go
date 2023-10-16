package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ilya-rusyanov/gophermart/internal/entities"
)

type AuthUsecase interface {
	Auth(context.Context, entities.AuthCredentials) (entities.AuthToken, error)
}

type Auth struct {
	errHandler ErrorHandler
	usecase    AuthUsecase
}

func NewAuth(usecase AuthUsecase, errHandler ErrorHandler) *Auth {
	return &Auth{
		errHandler: errHandler,
		usecase:    usecase,
	}
}

func (a *Auth) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	var creds entities.AuthCredentials
	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		a.errHandler.Handle(rw,
			fmt.Errorf("failed to parse credentials from JSON: %w",
				errParsing),
		)
		return
	}
	token, err := a.usecase.Auth(r.Context(), creds)
	if err != nil {
		a.errHandler.Handle(rw, fmt.Errorf("usecase failure: %w", err))
		return
	}
	processAuthToken(rw, token)
	rw.WriteHeader(http.StatusOK)
}

func processAuthToken(rw http.ResponseWriter, token entities.AuthToken) {
	c := http.Cookie{
		Name:  "access_token",
		Value: string(token),
		Path:  "/",
	}

	http.SetCookie(rw, &c)
}
