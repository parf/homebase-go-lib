# homebase-go-lib

A comprehensive Go library for file processing, database operations, and system utilities with extensive compression support for golang HomeBase framework.

## Installation

```bash
go get github.com/parf/homebase-go-lib
```

## Quick Start

```go
import (
    "github.com/parf/homebase-go-lib/fileiterator"
    "github.com/parf/homebase-go-lib/clistat"
    hb "github.com/parf/homebase-go-lib"
)

func main() {
    // Process compressed files automatically
    fileiterator.IterateLines("data.txt.gz", func(line string) {
        // Process each line
    })

    // Parse JSON Lines with type safety
    type User struct {
        Name  string `json:"name"`
        Email string `json:"email"`
    }
    fileiterator.IterateJSONLTyped("users.jsonl.zst", func(user User) error {
        // Process each user
        return nil
    })

    // Track processing statistics
    stat := clistat.New(10)
    for i := 0; i < 1000000; i++ {
        stat.Hit()  // Auto-reports progress
    }
    stat.Finish()

    // Scale values logarithmically
    scale := hb.Scale(1024)  // Returns 6
}
```

## Structured Data Format Examples

```go
import "github.com/parf/homebase-go-lib/fileiterator"

type Record struct {
    ID    int64  `json:"id"`
    Name  string `json:"name"`
    Email string `json:"email"`
}

// Process different formats (all support compression auto-detection)
func processData() {
    // CSV with automatic decompression
    // fileiterator.IterateCSVMap("data.csv.gz", func(row map[string]string) error {
    //     fmt.Printf("Name: %s\n", row["name"])
    //     return nil
    // })

    // JSON Lines with type safety
    // fileiterator.IterateJSONLTyped("records.jsonl.gz", func(rec Record) error {
    //     fmt.Printf("Processing: %s\n", rec.Name)
    //     return nil
    // })

    // MessagePack (more compact than JSON)
    fileiterator.IterateMsgPackTyped("data.msgpack.zst", func(rec Record) error {
        fmt.Printf("Record: %d - %s\n", rec.ID, rec.Name)
        return nil
    })

    // Apache Parquet (columnar format, best for analytics)
    fileiterator.IterateParquet("analytics.parquet", func(record map[string]interface{}) error {
        fmt.Printf("ID: %v, Name: %v\n", record["id"], record["name"])
        return nil
    })

    // FlatBuffer List (fastest reads, 3x faster than JSONL)
    fileiterator.IterateFlatBufferList("events.fb.lz4", func(data []byte) error {
        // Deserialize FlatBuffer bytes
        return nil
    })
}
```

**Format Selection Guide:**
- **CSV**: Human-readable, Excel compatible, simple structure
- **JSON Lines**: Widely supported, easy debugging, line-by-line processing
- **MessagePack**: 2x smaller than JSON, faster parsing
- **Parquet**: Best for analytics, columnar compression, SQL-like queries
- **FlatBuffer**: Fastest reads (3x), zero-copy deserialization, binary format

**Compression Support:** All formats work with `.gz`, `.zst`, `.lz4`, `.br`, `.xz` extensions automatically.

## Features

### General Utilities
- **Scale**: Logarithmic base-4 scaling (maps values to 0-9 range)
- **Any2uint32**: Type-safe integer conversion
- **DumpSortedMap**: Print maps in sorted key order

### Performance & Monitoring (`clistat/`)
- **CliStat**: Real-time statistics tracker with hits-per-second reporting
- **New(timeout)**: Create tracker with configurable timeout
- **Hit()**: Register event and auto-report progress
- **Finish()**: Print final statistics

### Task Management
- **Runner**: Parallel and sequential task execution with memory/timing stats
- **JobScheduler**: Periodic job scheduler with start/stop control
- **MemReport**: Memory allocation tracking and reporting

### Database Utilities (`sql/`)
- **BatchInserter**: Batch SQL insert operations with auto-escaping support
- **SqlIterator**: Iterate over SQL queries with statistics
- **SqlExtra**: Execute SQL queries and get results as maps (dynamic schema)

### File Iterator Package (`fileiterator/`)

**Unified file processing with automatic compression detection for 7 formats!**

#### Universal File Loaders (Auto-Detection)
- **FUOpen**: Universal file/URL opener with auto-decompression
- **LoadBinFile**: Load binary files with automatic decompression
- **IterateLines**: Process text files line-by-line with auto-decompression
- **IterateIDTabFile**: Process tab-separated hex ID-name pairs

#### Binary Record Iterators
- **IterateBinaryRecords**: Fixed-size binary records (auto-detects compression)
- **IterateZlibRecords** / **IterateGzipRecords** / **IterateZstdRecords**: Explicit format control

#### Structured Data Iterators
- **IterateJSONL**: Process JSON Lines files (untyped maps)
- **IterateJSONLTyped**: Process JSON Lines with type-safety (Go generics)
- **IterateCSV**: Process CSV files row-by-row (arrays)
- **IterateCSVMap**: Process CSV files with headers as maps
- **IterateMsgPackTyped**: Process MessagePack files with type-safety
- **IterateParquet**: Process Apache Parquet columnar files
- **IterateFlatBufferList**: Process FlatBuffer list files (fastest reads)

#### Compression Support
**7 formats with automatic detection:**
- Gzip (.gz)
- Zstd (.zst)
- Zlib (.zlib, .zz)
- LZ4 (.lz4)
- Brotli (.br)
- XZ (.xz)
- Plain files

#### Additional Features
- **URL Support**: HTTP/HTTPS URLs work with all functions
- **Custom Delimiters**: CSV with any delimiter (comma, tab, pipe, etc.)
- **Streaming**: Low memory usage for large files
- **Progress Tracking**: Automatic progress reporting

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
├── internal/         # Private packages (compression legacy loaders)
├── testdata/         # Test fixtures and data files
├── clistat/          # CLI statistics tracking package
├── fileiterator/     # File processing and iteration package
├── sql/              # SQL utilities and batch operations
├── homebase.go       # Main library code (utilities)
├── homebase_test.go  # Tests
├── go.mod            # Go module definition
├── Makefile          # Build automation
└── README.md         # This file
```

## Key Packages

| Package | Purpose | Key Functions |
|---------|---------|---------------|
| `fileiterator` | File processing with compression | `FUOpen`, `IterateLines`, `IterateJSONL`, `IterateCSV` |
| `clistat` | Statistics tracking | `New`, `Hit`, `Finish` |
| `sql` | Database utilities | `BatchInserter`, `SqlIterator`, `WildSqlQuery` |
| Main package | General utilities | `Scale`, `Any2uint32`, `Runner`, `JobScheduler` |

## Examples

See the `examples/` directory for complete working examples:
- `examples/fileiterator/` - File processing examples
- `examples/clistat/` - Statistics tracking
- `examples/sql/` - Database operations

## Contributing

Contributions are welcome. Please open an issue or submit a pull request.

## License

MIT License - see LICENSE file for details.
