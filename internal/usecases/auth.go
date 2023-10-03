package usecases

import (
	"context"
	"errors"

	"github.com/ilya-rusyanov/gophermart/internal/entities"
)

type Auth struct {
}

func NewAuth() *Auth {
	return &Auth{}
}

func (a *Auth) Register(context.Context, entities.AuthCredentials) (
	entities.Token, error,
) {
	return entities.Token{}, errors.New("TODO")
}
