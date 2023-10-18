package usecases

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/ilya-rusyanov/gophermart/internal/entities"
)

type RegisterStorage interface {
	FindUser(context.Context, entities.Login) error
	AddCredentials(context.Context, entities.AuthCredentials) error
}

type Register struct {
	storage    RegisterStorage
	logger     Logger
	expiration time.Duration
	key        string
}

func NewRegister(logger Logger, tokenExpiration time.Duration, signingKey string, storage RegisterStorage) *Register {
	return &Register{
		storage:    storage,
		logger:     logger,
		expiration: tokenExpiration,
		key:        signingKey,
	}
}

// register user
func (r *Register) Auth(ctx context.Context, creds entities.AuthCredentials) (
	entities.AuthToken, error,
) {
	var result entities.AuthToken

	err := r.storage.FindUser(ctx, creds.Login)
	switch {
	case errors.Is(err, entities.ErrNotFound):
		// all okay
		break
	case err == nil:
		return result,
			fmt.Errorf(
				"username already registered: %w",
				entities.ErrLoginConflict)
	default:
		return result, fmt.Errorf("unexpected storage error: %w", err)
	}

	// TODO: salt and hash password

	err = r.storage.AddCredentials(ctx, creds)
	if err != nil {
		return result, fmt.Errorf("failed to store credentials: %w", err)
	}

	result, err = buildAuthToken(r.expiration, creds.Login, r.key)
	if err != nil {
		return result, fmt.Errorf("failed to build auth token: %w", err)
	}

	r.logger.Infof("user %q registered", creds.Login)

	return result, nil
}
