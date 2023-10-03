package http

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/ilya-rusyanov/gophermart/internal/entities"
)

type Logger interface {
}

type Auth interface {
	Register(context.Context, entities.AuthCredentials) (
		entities.Token, error)
}

type Adapter struct {
	logger Logger
	auth   Auth
}

func New(logger Logger, auth Auth) *Adapter {
	return &Adapter{
		logger: logger,
		auth:   auth,
	}
}

func (a *Adapter) Register(rw http.ResponseWriter, r *http.Request) {
	var (
		authCred entities.AuthCredentials
	)

	err := json.NewDecoder(r.Body).Decode(&authCred)
	if err != nil {
		rw.WriteHeader(http.StatusUnauthorized)
		return
	}

	token, err := a.auth.Register(r.Context(), authCred)
	if err != nil {
		rw.WriteHeader(http.StatusUnauthorized)
		return
	}

	http.SetCookie(rw, a.buildAuthCookie(token))
	rw.WriteHeader(http.StatusOK)
}

func (a *Adapter) buildAuthCookie(token entities.Token) *http.Cookie {
	return &http.Cookie{
		Name:  "access_token",
		Value: token.Value,
	}
}
