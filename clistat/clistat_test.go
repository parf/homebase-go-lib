package clistat_test

import (
	"testing"
	"time"

	"github.com/parf/homebase-go-lib/clistat"
)

func TestNew(t *testing.T) {
	timeout := int64(10)
	stat := clistat.New(timeout)

	if stat.Timeout != timeout {
		t.Errorf("Expected timeout %d, got %d", timeout, stat.Timeout)
	}

	if stat.Cnt != 0 {
		t.Errorf("Expected Cnt to be 0, got %d", stat.Cnt)
	}

	if stat.Start == 0 {
		t.Error("Start time should not be 0")
	}

	if stat.Ltime == 0 {
		t.Error("Ltime should not be 0")
	}
}

func TestHit(t *testing.T) {
	stat := clistat.New(1)

	// Test single hit
	stat.Hit()
	if stat.Cnt != 1 {
		t.Errorf("Expected Cnt to be 1, got %d", stat.Cnt)
	}

	// Test multiple hits
	for i := 0; i < 100; i++ {
		stat.Hit()
	}
	if stat.Cnt != 101 {
		t.Errorf("Expected Cnt to be 101, got %d", stat.Cnt)
	}
}

func TestHitWithTimeout(t *testing.T) {
	stat := clistat.New(1)

	// Hit 256 times to trigger logging (Cnt&255 == 0)
	for i := 0; i < 256; i++ {
		stat.Hit()
	}

	// Wait for timeout to elapse
	time.Sleep(2 * time.Second)

	// Hit 256 more times to trigger another log
	for i := 0; i < 256; i++ {
		stat.Hit()
	}

	if stat.Cnt != 512 {
		t.Errorf("Expected Cnt to be 512, got %d", stat.Cnt)
	}
}

func TestFinish(t *testing.T) {
	stat := clistat.New(10)

	// Do some hits
	for i := 0; i < 1000; i++ {
		stat.Hit()
	}

	// Finish should not panic
	stat.Finish()

	if stat.Cnt != 1000 {
		t.Errorf("Expected Cnt to be 1000, got %d", stat.Cnt)
	}
}
