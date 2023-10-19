package handlers

type Logger interface {
	Infof(string, ...any)
	Info(...any)
	Error(...any)
}
