# homebase-go-lib

A Go library for [your project description].

## Installation

```bash
go get github.com/parf/homebase-go-lib
```

## Usage

```go
import hb "github.com/parf/homebase-go-lib"

func main() {
    // Use hb package
    // Example: hb.Version

    // Scale function - logarithmic scaling to 0-9 range
    scale := hb.Scale(1024)  // Returns 6
}
```

## Features

### General Utilities (`General.go`)
- **Scale**: Logarithmic base-4 scaling function (0-9 range)
- **Any2uint32**: Convert various integer types to uint32
- **DumpSortedMap**: Print maps in sorted key order

### Performance & Monitoring
- **Clistat**: CLI statistics tracker for hits per second monitoring
- **Runner**: Parallel and sequential task runners with memory/timing stats
- **JobScheduler**: Periodic job scheduler with start/stop control
- **MemReport**: Memory allocation reporting

### Database Utilities
- **BatchInserter**: Batch SQL insert operations for performance
- **SqlIterator**: Iterate over SQL queries with statistics
- **sql/SqlExtra**: Execute SQL queries and get results as maps (dynamic schema)

### File Processing
- **GZLoaders**: Load and process gzipped files (lines, binary, tab-separated)
- **Zlib**: Process zlib-compressed binary files
- **FUOpen**: Universal file/URL opener

### Debugging & Logging
- **Debug**: Configurable debug output (stderr or log)
- **Syslog**: System log utilities (notice/error)

## Development

### Prerequisites

- Go 1.25+
- Make (optional, for using Makefile commands)

### Building

```bash
make build
```

### Running Tests

```bash
make test
```

### Running Tests with Coverage

```bash
make test-coverage
```

### Formatting Code

```bash
make fmt
```

### Linting

```bash
make lint
```

## Project Structure

```
.
├── .github/          # GitHub Actions CI/CD workflows
├── docs/             # Additional documentation
├── examples/         # Example usage code
├── internal/         # Private packages (cannot be imported by external projects)
├── testdata/         # Test fixtures and data files
├── clistat/          # CLI statistics package
├── sql/              # SQL utilities package
├── homebase.go       # Main library code
├── homebase_test.go  # Tests
├── go.mod            # Go module definition
├── Makefile          # Build automation
└── README.md         # This file
```

## Contributing

Contributions are welcome. Please open an issue or submit a pull request.

## License

MIT License - see LICENSE file for details.
