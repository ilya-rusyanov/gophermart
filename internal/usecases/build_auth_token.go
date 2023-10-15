package usecases

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/ilya-rusyanov/gophermart/internal/entities"
)

func buildAuthToken(
	expiration time.Duration, userID entities.Login, key string,
) (entities.AuthToken, error) {
	var result entities.AuthToken

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, entities.TokenClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(
				time.Now().Add(expiration)),
		},
		UserID: userID,
	})

	tokenString, err := token.SignedString([]byte(key))
	if err != nil {
		return result, fmt.Errorf("failed to sign token: %w", err)
	}
	result = entities.AuthToken(tokenString)
	return result, nil
}
