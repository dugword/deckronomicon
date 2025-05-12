package logger

import (
	"fmt"
	"io"
)

type Logger struct {
	stdout io.Writer
	stderr io.Writer
}

func NewLogger(stdout io.Writer, stderr io.Writer) *Logger {
	return &Logger{
		stdout: stdout,
		stderr: stderr,
	}
}

func (l *Logger) Log(message string) {
	fmt.Fprintln(l.stdout, message)
}

func (l *Logger) Error(message string) {
	fmt.Fprintln(l.stderr, message)
}

type LogBuffer struct {
	Entries []string
}

func (lb *LogBuffer) Log(message string) {
	lb.Entries = append(lb.Entries, message)
}

func (lb *LogBuffer) Error(message string) {
	lb.Entries = append(lb.Entries, "ERROR: "+message)
}
