package usecases

import (
	"context"
	"errors"

	"github.com/ilya-rusyanov/gophermart/internal/entities"
)

type Register struct {
}

func NewRegister() *Register {
	return &Register{}
}

func (a *Register) Register(context.Context, entities.AuthCredentials) (
	entities.AuthToken, error,
) {
	return entities.AuthToken(""), errors.New("TODO")
}
