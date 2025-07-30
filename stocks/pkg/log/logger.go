package log

// Logger is the universal logger that can do everything.
type Logger interface {
	loggerStructured
	loggerFmt
}

type loggerStructured interface {
	// Trace logs at Trace log level using fields.
	Trace(msg string, fields ...Field)
	// Debug logs at Debug log level using fields.
	Debug(msg string, fields ...Field)
	// Info logs at Info log level using fields.
	Info(msg string, fields ...Field)
	// Warn logs at Warn log level using fields.
	Warn(msg string, fields ...Field)
	// Error logs at Error log level using fields.
	Error(msg string, fields ...Field)
	// Fatal logs at Fatal log level using fields.
	Fatal(msg string, fields ...Field)
}

type loggerFmt interface {
	// Tracef logs at Trace log level using fmt formatter.
	Tracef(format string, args ...interface{})
	// Debugf logs at Debug log level using fmt formatter.
	Debugf(format string, args ...interface{})
	// Infof logs at Info log level using fmt formatter.
	Infof(format string, args ...interface{})
	// Warnf logs at Warn log level using fmt formatter.
	Warnf(format string, args ...interface{})
	// Errorf logs at Error log level using fmt formatter.
	Errorf(format string, args ...interface{})
	// Fatalf logs at Fatal log level using fmt formatter.
	Fatalf(format string, args ...interface{})
}
