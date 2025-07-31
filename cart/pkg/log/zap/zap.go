package zap

import (
	"cart/pkg/log"
	"fmt"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var _ log.Logger = &Logger{}

type Logger struct {
	L *zap.Logger
}

func NewLogger() (*Logger, func(), error) {
	// conn, err := net.Dial("tcp",)
	// if err != nil {
	// 	return nil, nil, err
	// }

	cfg := zap.NewProductionEncoderConfig()
	cfg.TimeKey = "timestamp"
	cfg.LevelKey = "severity"
	cfg.MessageKey = "message"

	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(cfg),
		zapcore.AddSync(zapcore.Lock(os.Stdout)), // Use os.Stdout
		zapcore.InfoLevel,
	)

	z := zap.New(core).With(
		zap.String("service", "logger-practice"),
		zap.String("env", "dev"))

	l := &Logger{
		L: z,
		// c: conn,
	}

	cleanup := func() {
		_ = l.L.Sync()
		// _ = l.c.Close()
	}

	return l, cleanup, nil
}

// Trace logs at Trace log level using fields
func (l *Logger) Trace(msg string, fields ...log.Field) {
	if ce := l.L.Check(zap.DebugLevel, msg); ce != nil {
		ce.Write(zapifyFields(fields...)...)
	}
}

// Tracef logs at Trace log level using fmt formatter
func (l *Logger) Tracef(msg string, args ...interface{}) {
	if ce := l.L.Check(zap.DebugLevel, ""); ce != nil {
		ce.Message = fmt.Sprintf(msg, args...)
		ce.Write()
	}
}

// Debug logs at Debug log level using fields
func (l *Logger) Debug(msg string, fields ...log.Field) {
	if ce := l.L.Check(zap.DebugLevel, msg); ce != nil {
		ce.Write(zapifyFields(fields...)...)
	}
}

// Debugf logs at Debug log level using fmt formatter
func (l *Logger) Debugf(msg string, args ...interface{}) {
	if ce := l.L.Check(zap.DebugLevel, ""); ce != nil {
		ce.Message = fmt.Sprintf(msg, args...)
		ce.Write()
	}
}

// Info logs at Info log level using fields
func (l *Logger) Info(msg string, fields ...log.Field) {
	if ce := l.L.Check(zap.InfoLevel, msg); ce != nil {
		ce.Write(zapifyFields(fields...)...)
	}
}

// Infof logs at Info log level using fmt formatter
func (l *Logger) Infof(msg string, args ...interface{}) {
	if ce := l.L.Check(zap.InfoLevel, ""); ce != nil {
		ce.Message = fmt.Sprintf(msg, args...)
		ce.Write()
	}
}

// Warn logs at Warn log level using fields
func (l *Logger) Warn(msg string, fields ...log.Field) {
	if ce := l.L.Check(zap.WarnLevel, msg); ce != nil {
		ce.Write(zapifyFields(fields...)...)
	}
}

// Warnf logs at Warn log level using fmt formatter
func (l *Logger) Warnf(msg string, args ...interface{}) {
	if ce := l.L.Check(zap.WarnLevel, ""); ce != nil {
		ce.Message = fmt.Sprintf(msg, args...)
		ce.Write()
	}
}

// Error logs at Error log level using fields
func (l *Logger) Error(msg string, fields ...log.Field) {
	if ce := l.L.Check(zap.ErrorLevel, msg); ce != nil {
		ce.Write(zapifyFields(fields...)...)
	}
}

// Errorf logs at Error log level using fmt formatter
func (l *Logger) Errorf(msg string, args ...interface{}) {
	if ce := l.L.Check(zap.ErrorLevel, ""); ce != nil {
		ce.Message = fmt.Sprintf(msg, args...)
		ce.Write()
	}
}

// Fatal logs at Fatal log level using fields
func (l *Logger) Fatal(msg string, fields ...log.Field) {
	if ce := l.L.Check(zap.FatalLevel, msg); ce != nil {
		ce.Write(zapifyFields(fields...)...)
	}
}

// Fatalf logs at Fatal log level using fmt formatter
func (l *Logger) Fatalf(msg string, args ...interface{}) {
	if ce := l.L.Check(zap.FatalLevel, ""); ce != nil {
		ce.Message = fmt.Sprintf(msg, args...)
		ce.Write()
	}
}
