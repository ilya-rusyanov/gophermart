package usecases

type DummyLogger struct {
}

func (l *DummyLogger) Error(...any) {
}

func (l *DummyLogger) Errorf(string, ...any) {
}

func (l *DummyLogger) Infof(string, ...any) {
}

func (l *DummyLogger) Info(...any) {
}
