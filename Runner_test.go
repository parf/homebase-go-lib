package hb_test

import (
	"sync/atomic"
	"testing"
	"time"

	hb "github.com/parf/homebase-go-lib"
)

func TestParallelRunner(t *testing.T) {
	runner := hb.NewParallelRunner()

	var counter int32

	runner.Run("task1", func() {
		time.Sleep(100 * time.Millisecond)
		atomic.AddInt32(&counter, 1)
	})

	runner.Run("task2", func() {
		time.Sleep(100 * time.Millisecond)
		atomic.AddInt32(&counter, 1)
	})

	runner.Finish()

	if atomic.LoadInt32(&counter) != 2 {
		t.Errorf("Expected counter to be 2, got %d", counter)
	}
}

func TestSequentialRunner(t *testing.T) {
	runner := hb.NewSequentialRunner()

	var counter int32

	runner.Run("task1", func() {
		time.Sleep(50 * time.Millisecond)
		atomic.AddInt32(&counter, 1)
	})

	runner.Run("task2", func() {
		time.Sleep(50 * time.Millisecond)
		atomic.AddInt32(&counter, 1)
	})

	runner.Finish()

	if atomic.LoadInt32(&counter) != 2 {
		t.Errorf("Expected counter to be 2, got %d", counter)
	}
}

func TestMemReport(t *testing.T) {
	// Just ensure it doesn't panic
	hb.MemReport("test event")
	hb.MemReport("")
}
