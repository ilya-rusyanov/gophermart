package handlers

import (
	"net/http"

	"github.com/ilya-rusyanov/gophermart/internal/entities"
)

func processAuthToken(rw http.ResponseWriter, token entities.AuthToken) {
	c := http.Cookie{
		Name:  "access_token",
		Value: string(token),
	}

	http.SetCookie(rw, &c)
}
