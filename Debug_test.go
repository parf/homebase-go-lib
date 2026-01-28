package hb_test

import (
	"bytes"
	"log"
	"os"
	"strings"
	"testing"

	hb "github.com/parf/homebase-go-lib"
)

func TestDebug(t *testing.T) {
	// Capture stderr
	old := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	debug := hb.Debug("TEST", 2)

	// This should show (level 0 <= 2)
	debug(0, "error message")

	// This should show (level 1 <= 2)
	debug(1, "warning message")

	// This should show (level 2 <= 2)
	debug(2, "info message")

	// This should NOT show (level 3 > 2)
	debug(3, "debug message")

	w.Close()
	os.Stderr = old

	var buf bytes.Buffer
	buf.ReadFrom(r)
	output := buf.String()

	if !strings.Contains(output, "error message") {
		t.Error("Expected to see error message")
	}

	if !strings.Contains(output, "warning message") {
		t.Error("Expected to see warning message")
	}

	if strings.Contains(output, "debug message") {
		t.Error("Should not see debug message (level too high)")
	}
}

func TestDebugLog(t *testing.T) {
	// Capture log output
	var buf bytes.Buffer
	log.SetOutput(&buf)
	defer log.SetOutput(os.Stderr)

	debug := hb.DebugLog("TEST", 1)

	debug(0, "error message")
	debug(1, "warning message")
	debug(2, "should not appear")

	output := buf.String()

	if !strings.Contains(output, "error message") {
		t.Error("Expected to see error message")
	}

	if !strings.Contains(output, "warning message") {
		t.Error("Expected to see warning message")
	}

	if strings.Contains(output, "should not appear") {
		t.Error("Should not see message with level > debug level")
	}
}

func TestDebugFormatting(t *testing.T) {
	var buf bytes.Buffer
	log.SetOutput(&buf)
	defer log.SetOutput(os.Stderr)

	debug := hb.DebugLog("TEST", 3)

	debug(1, "formatted: %d %s", 42, "test")

	output := buf.String()

	if !strings.Contains(output, "formatted: 42 test") {
		t.Error("Expected formatted output")
	}
}
