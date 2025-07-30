package logger

var _ Logger = (*logger)(nil)

type Logger interface {
}

type logger struct {
}
