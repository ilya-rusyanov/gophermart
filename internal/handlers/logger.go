package handlers

type Logger interface {
	Infof(string, ...any)
	Error(...any)
}
