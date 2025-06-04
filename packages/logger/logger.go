package logger

import (
	"encoding/json"
	"fmt"
	"time"
)

type LogLevel int

const (
	DebugLevel LogLevel = iota
	InfoLevel
	WarningLevel
	ErrorLevel
	CriticalLevel
)

var logLevelNames = map[LogLevel]struct {
	label string
	emoji string
	color string
}{
	DebugLevel:    {"DEBUG", "🐛", "\033[36m"},    // Cyan
	InfoLevel:     {"INFO", "✨", "\033[32m"},     // Green
	WarningLevel:  {"WARNING", "⚠️", "\033[33m"}, // Yellow
	ErrorLevel:    {"ERROR", "🔥", "\033[31m"},    // Red
	CriticalLevel: {"CRITICAL", "💀", "\033[41m"}, // Red background
}

const reset = "\033[0m"

type Logger struct {
	IncludeContext bool
	ContextFunc    func() any
}

func NewLogger() *Logger {
	return &Logger{}
}

func BeautifyContext(context any) string {
	out, err := json.MarshalIndent(context, "", "  ")
	if err != nil {
		panic(fmt.Sprintf("Failed to beautify context: %v", err))
	}
	return string(out)
}

func PrintContext(context any) {
	if context == nil {
		return
	}
	fmt.Printf(
		"%s [%s] %s %s%s\n",
		"\033[37m", // Light Grey
		"CONTEXT",
		"📎",
		BeautifyContext(context),
		reset,
	)
}

func (l *Logger) logf(level LogLevel, format string, args ...any) {
	ll := logLevelNames[level]
	ts := time.Now().Format("15:04:05")
	msg := fmt.Sprintf(format, args...)
	fmt.Printf(
		"%s%s [%s] %s %s%s\n",
		ll.color,
		ts,
		ll.label,
		ll.emoji,
		msg,
		reset,
	)
	if !l.IncludeContext || l.ContextFunc == nil {
		return
	}
	PrintContext(l.ContextFunc())
}

func argsToString(args ...any) string {
	var out string
	for _, arg := range args {
		out = fmt.Sprintf("%s %v", out, arg)
	}
	return out
}

func (l *Logger) Debug(args ...any) {
	l.logf(DebugLevel, "%s", argsToString(args...))
}

func (l *Logger) Debugf(format string, args ...any) {
	l.logf(DebugLevel, format, args...)
}

func (l *Logger) Info(args ...any) {
	l.logf(InfoLevel, "%s", argsToString(args...))
}

func (l *Logger) Infof(format string, args ...any) {
	l.logf(InfoLevel, format, args...)
}

func (l *Logger) Warn(args ...any) {
	l.logf(WarningLevel, "%s", argsToString(args...))
}

func (l *Logger) Warnf(format string, args ...any) {
	l.logf(WarningLevel, format, args...)
}

func (l *Logger) Error(args ...any) {
	l.logf(ErrorLevel, "%s", argsToString(args...))
}

func (l *Logger) Errorf(format string, args ...any) {
	l.logf(ErrorLevel, format, args...)
}

func (l *Logger) Critical(args ...any) {
	l.logf(CriticalLevel, "%s", argsToString(args...))
}

func (l *Logger) Criticalf(format string, args ...any) {
	l.logf(CriticalLevel, format, args...)
}
