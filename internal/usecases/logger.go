package usecases

type Logger interface {
	Infof(string, ...any)
}
