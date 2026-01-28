# Clistat - CLI Statistics Tracker

Package clistat provides a simple statistics tracker for command-line applications to monitor hits per second (HPS) and track progress.

## Features

- Track hit counts over time
- Automatic periodic progress reporting
- Calculate hits per second (HPS)
- Final statistics on completion

## Usage

```go
package main

import (
    "github.com/parf/homebase-go-lib/clistat"
)

func main() {
    // Create a new tracker with 5 second reporting timeout
    stat := clistat.New(5)

    // Track hits in your processing loop
    for i := 0; i < 1000000; i++ {
        // Do some work...

        stat.Hit()  // Record each hit
    }

    // Print final statistics
    stat.Finish()
}
```

## How It Works

- **New(timeout int64)**: Creates a new Clistat tracker with the specified timeout in seconds
- **Hit()**: Records a hit. Progress is logged every 256 hits if the timeout has elapsed
- **Finish()**: Prints final statistics including total hits and elapsed time

## Performance

The Hit() method is optimized to only check the timer every 256 hits (using bitwise AND with 255), minimizing overhead for high-frequency operations.

## Output Example

```
Cnt=256 K. HPS=51 K
Cnt=512 K. HPS=48 K
DONE, Hits: 1000000, Seconds: 20
```

## TODO

- Add prefix support for customizing log output
