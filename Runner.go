package hb

/*
	Parallel and Sequential Task Runner with logging and stats

    log Task start/finish - ellapsed time, memory diff & total  (no diff for parallel mode)
    log Grand Finish -  ellapsed time, memory diff & total

	runner := hb.NewParallelRunner()
	// runner := hb.NewSequentialRunner()  - DROP-IN debug replacement
	runner.Run(func())
	runner.Finish()

*/

import (
	"fmt"
	"log"
	"runtime"
	"sync"
	"time"
)

// ParallelRunner runs tasks in parallel with performance tracking
type ParallelRunner struct {
	start    time.Time
	startMem float64
	wg       *sync.WaitGroup
}

// SequentialRunner is a debug option to use instead of ParallelRunner
type SequentialRunner struct {
	start    time.Time
	startMem float64
}

func memStatsAllocMB() float64 {
	memStats := &runtime.MemStats{}
	runtime.ReadMemStats(memStats)
	return float64(memStats.Alloc) / 0x100000
}

var (
	g_memAllocated uint64
)

// MemReport reports allocated memory to STDOUT
func MemReport(event string) {
	memStats := &runtime.MemStats{}
	runtime.ReadMemStats(memStats)
	diff := int64(memStats.Alloc) - int64(g_memAllocated)
	if g_memAllocated == 0 {
		diff = 0
	}
	if event != "" {
		fmt.Print(event + " ")
	}
	fmt.Printf(" - MEM(MB) allocated:%.1f diff:%.1f\n", float64(memStats.Alloc)/0x100000, float64(diff)/0x100000)
	g_memAllocated = memStats.Alloc
}

// NewParallelRunner creates a new parallel task runner
func NewParallelRunner() ParallelRunner {
	var wg sync.WaitGroup
	return ParallelRunner{time.Now(), memStatsAllocMB(), &wg}
}

// Run starts a task in parallel and logs start/finish with elapsed time and memory total
func (p *ParallelRunner) Run(name string, f func()) {
	p.wg.Add(1)
	go func() {
		defer p.wg.Done()
		start := time.Now()
		log.Println("* Started: " + name)
		f()
		log.Printf("* Finished: %v %.1f seconds, mem(MB) total: %.1f\n", name, time.Since(start).Seconds(), memStatsAllocMB())
	}()
}

// Finish waits for all parallel tasks to complete and logs final statistics
func (p *ParallelRunner) Finish() {
	p.wg.Wait()
	endMem := memStatsAllocMB()
	log.Printf("--- Finished RUNNER %.1f seconds, mem(MB){ diff: %.1f total: %.1f}\n", time.Since(p.start).Seconds(), endMem-p.startMem, endMem)
}

// NewSequentialRunner creates a drop-in debug replacement for ParallelRunner
// Usage:
//
//	runner := hb.NewSequentialRunner()
//	runner.Run(func())
//	runner.Finish()
func NewSequentialRunner() SequentialRunner {
	return SequentialRunner{time.Now(), memStatsAllocMB()}
}

// Run executes a task sequentially and logs start/finish with elapsed time, memory diff & total
func (p *SequentialRunner) Run(name string, f func()) {
	start := time.Now()
	startMem := memStatsAllocMB()
	log.Println("* Started: " + name)
	f()
	endMem := memStatsAllocMB()
	log.Printf("* Finished %v %.1f seconds, mem(MB){ diff: %.1f total: %.1f}\n", name, time.Since(start).Seconds(), endMem-startMem, endMem)
}

// Finish logs final statistics for the sequential runner
func (p *SequentialRunner) Finish() {
	endMem := memStatsAllocMB()
	log.Printf("--- Finished RUNNER %.1f seconds, mem(MB){ diff: %.1f total: %.1f}\n", time.Since(p.start).Seconds(), endMem-p.startMem, endMem)
}
