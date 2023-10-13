package handlers

import "net/http"

type Register struct {
}

func NewRegister() *Register {
	return &Register{}
}

func (reg *Register) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
}
