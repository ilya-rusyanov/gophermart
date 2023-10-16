package middleware

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/ilya-rusyanov/gophermart/internal/entities"
	"github.com/ilya-rusyanov/gophermart/internal/handlers"

	"github.com/golang-jwt/jwt/v4"
)

type ErrorHandler interface {
	Handle(http.ResponseWriter, error)
}

type Auth struct {
	signingKey       string
	excludeEndpoints map[string]struct{}
	errorHandler     ErrorHandler
}

func NewAuth(key string, errorHandler ErrorHandler, excludeEndpoints ...string) *Auth {
	exclude := make(map[string]struct{}, len(excludeEndpoints))
	for _, e := range excludeEndpoints {
		exclude[e] = struct{}{}
	}

	return &Auth{
		signingKey:       key,
		excludeEndpoints: exclude,
		errorHandler:     errorHandler,
	}
}

func (m *Auth) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		if m.applicable(r.URL) {
			cookie, err := r.Cookie("access_token")
			if err != nil {
				m.errorHandler.Handle(rw,
					fmt.Errorf(
						"failed to read access token from cookie: %w",
						entities.ErrUnauthorized))
				return
			}
			var login *entities.Login
			if !m.valid(*cookie, &login) || login == nil {
				m.errorHandler.Handle(rw,
					fmt.Errorf(
						"cookie error: %w",
						entities.ErrUnauthorized))
				return
			}

			ctx := context.WithValue(r.Context(),
				handlers.ContextKeyLogin, *login)
			r = r.WithContext(ctx)
		}

		next.ServeHTTP(rw, r)
	})
}

func (a *Auth) valid(cookie http.Cookie, username **entities.Login) bool {
	claims := entities.TokenClaims{}
	_, err := jwt.ParseWithClaims(cookie.Value, &claims,
		func(t *jwt.Token) (interface{}, error) {
			return []byte(a.signingKey), nil
		})

	if err != nil {
		return false
	}

	*username = &claims.UserID

	return true
}

func (a *Auth) applicable(url *url.URL) bool {
	_, ok := a.excludeEndpoints[url.Path]
	return !ok
}
