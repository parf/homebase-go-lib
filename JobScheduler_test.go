package hb_test

import (
	"sync/atomic"
	"testing"
	"time"

	hb "github.com/parf/homebase-go-lib"
)

func TestJobScheduler(t *testing.T) {
	var counter int32

	jobFunc := func() {
		atomic.AddInt32(&counter, 1)
	}

	scheduler := hb.NewJobScheduler(1, jobFunc)

	// Start scheduler
	if err := scheduler.Start(); err != nil {
		t.Fatalf("Failed to start scheduler: %v", err)
	}

	if !scheduler.IsRunning() {
		t.Error("Scheduler should be running")
	}

	// Wait for job to run a few times
	time.Sleep(2500 * time.Millisecond)

	// Stop scheduler
	if err := scheduler.Stop(); err != nil {
		t.Fatalf("Failed to stop scheduler: %v", err)
	}

	if scheduler.IsRunning() {
		t.Error("Scheduler should not be running")
	}

	// Counter should be at least 2 (immediate + 2 intervals)
	count := atomic.LoadInt32(&counter)
	if count < 2 {
		t.Errorf("Expected at least 2 job runs, got %d", count)
	}
}

func TestJobSchedulerDoubleStart(t *testing.T) {
	scheduler := hb.NewJobScheduler(1, func() {})

	if err := scheduler.Start(); err != nil {
		t.Fatalf("First start failed: %v", err)
	}
	defer scheduler.Stop()

	// Try to start again - should fail
	if err := scheduler.Start(); err == nil {
		t.Error("Expected error when starting already running scheduler")
	}
}

func TestJobSchedulerStopBeforeStart(t *testing.T) {
	scheduler := hb.NewJobScheduler(1, func() {})

	// Try to stop before starting - should fail
	if err := scheduler.Stop(); err == nil {
		t.Error("Expected error when stopping non-running scheduler")
	}
}
