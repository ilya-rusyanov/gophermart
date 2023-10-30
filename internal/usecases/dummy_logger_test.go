package usecases

type DummyLogger struct {
}

func (l *DummyLogger) Errorf(string, ...any) {
}

func (l *DummyLogger) Infof(string, ...any) {
}

func (l *DummyLogger) Info(...any) {
}
