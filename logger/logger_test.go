package logger_test

import (
	"deckronomicon/logger"
	"testing"
)

func TestLogBuffer(t *testing.T) {
	log := &logger.LogBuffer{}
	log.Log("Test entry")

	if len(log.Entries) != 1 {
		t.Fatal("Expected 1 log entry")
	}

	if log.Entries[0] != "Test entry" {
		t.Fatalf("Unexpected log content: %s", log.Entries[0])
	}
}
