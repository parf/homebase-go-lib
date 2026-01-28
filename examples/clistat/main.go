package main

import (
	"fmt"
	"time"

	"github.com/parf/homebase-go-lib/clistat"
)

func main() {
	// Create a new clistat with 5 second timeout for progress reporting
	stat := clistat.New(5)

	fmt.Println("Starting hit simulation...")
	fmt.Println("Progress will be logged every 5 seconds")

	// Simulate hits
	for i := 0; i < 100000; i++ {
		stat.Hit()

		// Simulate some work
		if i%1000 == 0 {
			time.Sleep(10 * time.Millisecond)
		}
	}

	// Print final statistics
	stat.Finish()
}
