package log_test

import (
	"bytes"
	"errors"
	"io"
	"strings"
	"testing"

	"github.com/nike-data-quality-engine/halfpipe/packages/log"
)

func TestNew(t *testing.T) {
	t.Run("create new logger", func(t *testing.T) {
		logger := log.NewLogger("📝", io.Discard, io.Discard, false)
		if logger == nil {
			t.Error("failed to create logger")
		}
	})
	t.Run("log an error", func(t *testing.T) {
		in := "this is an error"
		want := "🔥 ERROR: this is an error\n"
		stderr := bytes.Buffer{}
		logger := log.NewLogger("📝", io.Discard, &stderr, false)
		logger.Error(errors.New(in))
		got := stderr.String()
		if !strings.HasSuffix(got, want) {
			t.Errorf("logger.Error(%q) = %q; want: %q", in, got, want)
		}
	})
	t.Run("log a formatted error message", func(t *testing.T) {
		in := "this is an error"
		want := "🔥 ERROR: my error: this is an error\n"
		stderr := bytes.Buffer{}
		logger := log.NewLogger("📝", io.Discard, &stderr, false)
		logger.Errorf("my error: %s", errors.New(in))
		got := stderr.String()
		if !strings.HasSuffix(got, want) {
			t.Errorf(`logger.Errorf("my error: %%s", %q) = %q; want: %q`, in, got, want)
		}
	})
	t.Run("log an info message when verbose", func(t *testing.T) {
		in := "hey look here"
		want := "✨ INFO: hey look here\n"
		stdout := bytes.Buffer{}
		logger := log.NewLogger("📝", &stdout, io.Discard, true)
		logger.Info(in)
		got := stdout.String()
		if !strings.HasSuffix(got, want) {
			t.Errorf("logger.Info(%q) = %q; want: %q", in, got, want)
		}
	})
	t.Run("log a formatted info message when verbose", func(t *testing.T) {
		in := "hey look here"
		want := "✨ INFO: you: hey look here\n"
		stdout := bytes.Buffer{}
		logger := log.NewLogger("📝", &stdout, io.Discard, true)
		logger.Infof("you: %s", in)
		got := stdout.String()
		if !strings.HasSuffix(got, want) {
			t.Errorf(`logger.Infof("you: %%s", %q) = %q; want: %q`, in, got, want)
		}
	})
	t.Run("don't log an info message when not verbose", func(t *testing.T) {
		in := "hey look here"
		want := ""
		stdout := bytes.Buffer{}
		logger := log.NewLogger("📝", &stdout, io.Discard, false)
		logger.Info(in)
		got := stdout.String()
		if !strings.HasSuffix(got, want) {
			t.Errorf("logger.Info(%q) = %q; want: %q", in, got, want)
		}
	})
	t.Run("log a message", func(t *testing.T) {
		in := "this is it"
		want := "📝: this is it\n"
		stdout := bytes.Buffer{}
		logger := log.NewLogger("📝", &stdout, io.Discard, true)
		logger.Log(in)
		got := stdout.String()
		if !strings.HasSuffix(got, want) {
			t.Errorf("logger.Log(%q) = %q; want: %q", in, got, want)
		}
	})
	t.Run("log a formatted message", func(t *testing.T) {
		in := "this is it"
		want := "📝: really: this is it\n"
		stdout := bytes.Buffer{}
		logger := log.NewLogger("📝", &stdout, io.Discard, true)
		logger.Logf("really: %s", in)
		got := stdout.String()
		if !strings.HasSuffix(got, want) {
			t.Errorf(`logger.Logf("really: %%s", %q) = %q; want: %q`, in, got, want)
		}
	})
	t.Run("log a warning when verbose", func(t *testing.T) {
		in := "uh oh"
		want := "⚠️ WARN: uh oh\n"
		stderr := bytes.Buffer{}
		logger := log.NewLogger("📝", io.Discard, &stderr, true)
		logger.Warn(in)
		got := stderr.String()
		if !strings.HasSuffix(got, want) {
			t.Errorf("logger.Warn(%q) = %q; want: %q", in, got, want)
		}
	})
	t.Run("log a formatted warning when verbose", func(t *testing.T) {
		in := "uh oh"
		want := "⚠️ WARN: oh no: uh oh\n"
		stderr := bytes.Buffer{}
		logger := log.NewLogger("📝", io.Discard, &stderr, true)
		logger.Warnf("oh no: %s", in)
		got := stderr.String()
		if !strings.HasSuffix(got, want) {
			t.Errorf(`logger.Warnf("oh no: %%s", %q) = %q; want: %q`, in, got, want)
		}
	})
	t.Run("don't log a warning when not verbose", func(t *testing.T) {
		in := "uh oh"
		want := ""
		stderr := bytes.Buffer{}
		logger := log.NewLogger("📝", io.Discard, &stderr, false)
		logger.Warn(in)
		got := stderr.String()
		if !strings.HasSuffix(got, want) {
			t.Errorf("logger.Warn(%q) = %q; want: %q", in, got, want)
		}
	})
}
