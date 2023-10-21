package usecases

type Logger interface {
	Info(...any)
	Infof(string, ...any)
	Errorf(string, ...any)
}
