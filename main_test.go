package main

import (
	"testing"
)

func TestBar(t *testing.T) {
	result := "bar"
	if result != "bar" {
		t.Errorf("expecting bar, got %s", result)
	}
}
