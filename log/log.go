// Package log provides a logger object with standard logging level methods.
// Each log level method has a unique prefix added to the start of each message.
package log

import (
	"io"
	"log"
)

// Logger holds the underlying loggers for each log level
type Logger struct {
	error *log.Logger
	info  *log.Logger
	log   *log.Logger
	warn  *log.Logger
}

// Error always prints the values to stderr
func (l *Logger) Error(v ...any) {
	l.error.Println(v...)
}

// Errorf always prints the format string to stderr
func (l *Logger) Errorf(format string, v ...any) {
	l.error.Printf(format, v...)
}

// Info prints the values to stdout when verbose
func (l *Logger) Info(v ...any) {
	l.info.Println(v...)
}

// Infof prints the format string to stdout when verbose
func (l *Logger) Infof(format string, v ...any) {
	l.info.Printf(format, v...)
}

// Log always prints the values to stdout
func (l *Logger) Log(v ...any) {
	l.log.Println(v...)
}

// Logf always prints the format string to stdout
func (l *Logger) Logf(format string, v ...any) {
	l.log.Printf(format, v...)
}

// Warn prints the values to stderr when verbose
func (l *Logger) Warn(v ...any) {
	l.warn.Println(v...)
}

// Warnf prints the format string to stderr when verbose
func (l *Logger) Warnf(format string, v ...any) {
	l.warn.Printf(format, v...)
}

// NewLogger returns a new logger
func NewLogger(prefix string, stdout, stderr io.Writer, verbose bool) *Logger {
	l := &Logger{}
	l.error = log.New(stderr, "🔥 ERROR: ", log.Ldate|log.Lmicroseconds|log.Ltime|log.Lmsgprefix)
	l.info = log.New(stdout, "✨ INFO: ", log.Ldate|log.Lmicroseconds|log.Ltime|log.Lmsgprefix)
	l.log = log.New(stdout, prefix+": ", log.Ldate|log.Lmicroseconds|log.Ltime|log.Lmsgprefix)
	l.warn = log.New(stderr, "⚠️ WARN: ", log.Ldate|log.Lmicroseconds|log.Ltime|log.Lmsgprefix)
	if !verbose {
		l.info = log.New(io.Discard, "", 0)
		l.warn = log.New(io.Discard, "", 0)
	}
	return l
}
