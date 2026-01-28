package hb

/*

JobScheduler is a periodic job scheduler that executes a given job function at a specified interval (seconds).
You can start and stop the scheduler, and it will run the job function immediately upon starting.
The scheduler uses a goroutine and a ticker to manage the timing of job execution

Usage:
	scheduler := hb.NewJobScheduler(5, jobFunc)
	if err := scheduler.Start(); err != nil {
	... handle error
	}
	...
	if err := scheduler.Stop(); err != nil {
	... handle error
	}

*/

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"
)

// JobScheduler represents a periodic job scheduler
type JobScheduler struct {
	interval  time.Duration // seconds
	jobFunc   func()
	ctx       context.Context
	cancel    context.CancelFunc
	wg        sync.WaitGroup
	isRunning bool
	mu        sync.Mutex
}

// NewJobScheduler creates a new job scheduler
func NewJobScheduler(intervalSeconds int, jobFunc func()) *JobScheduler {
	ctx, cancel := context.WithCancel(context.Background())

	return &JobScheduler{
		interval:  time.Duration(intervalSeconds) * time.Second,
		jobFunc:   jobFunc,
		ctx:       ctx,
		cancel:    cancel,
		isRunning: false,
	}
}

// Start begins the job scheduler
func (js *JobScheduler) Start() error {
	js.mu.Lock()
	defer js.mu.Unlock()

	if js.isRunning {
		return fmt.Errorf("job scheduler is already running")
	}

	js.isRunning = true
	js.wg.Add(1)

	go func() {
		defer js.wg.Done()

		// Run job immediately at start
		js.jobFunc()

		ticker := time.NewTicker(js.interval)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				js.jobFunc()
			case <-js.ctx.Done():
				log.Println("Job scheduler stopped")
				return
			}
		}
	}()

	log.Printf("Job scheduler started with interval of %v\n", js.interval)
	return nil
}

// Stop halts the job scheduler
func (js *JobScheduler) Stop() error {
	js.mu.Lock()
	defer js.mu.Unlock()

	if !js.isRunning {
		return fmt.Errorf("job scheduler is not running")
	}

	log.Println("Stopping job scheduler...")
	js.cancel()
	js.wg.Wait()
	js.isRunning = false

	log.Println("Job scheduler has been stopped gracefully")
	return nil
}

// IsRunning returns whether the scheduler is currently running
func (js *JobScheduler) IsRunning() bool {
	js.mu.Lock()
	defer js.mu.Unlock()
	return js.isRunning
}
