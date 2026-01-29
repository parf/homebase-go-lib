# 🏠 homebase-go-lib

[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)](https://go.dev/)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/parf/homebase-go-lib)](https://goreportcard.com/report/github.com/parf/homebase-go-lib)
[![GitHub Stars](https://img.shields.io/github/stars/parf/homebase-go-lib?style=social)](https://github.com/parf/homebase-go-lib/stargazers)

**High-performance Go library for data processing, file I/O, and system utilities** | Built for production workloads with automatic compression detection and format conversion

A comprehensive toolkit for building data pipelines, ETL workflows, and command-line tools in Go. Features universal file processing with 7 compression formats, structured data iteration (CSV, JSON Lines, Parquet, MsgPack), and production-ready utilities.

---

## 🌟 Key Features

<table>
<tr>
<td width="50%">

### 📁 Universal File Processing
- **7 compression formats** with auto-detection
- **5 structured formats**: CSV, JSONL, Parquet, MsgPack, FlatBuffer
- **HTTP/HTTPS URL support** for remote files
- **Streaming processing** for large files

</td>
<td width="50%">

### ⚡ High Performance
- **Zero-copy** operations where possible
- **Parallel processing** support
- **Memory efficient** streaming
- **Progress tracking** built-in

</td>
</tr>
<tr>
<td width="50%">

### 🔧 Production Ready
- **Type-safe** generic iterators
- **Error handling** throughout
- **Battle-tested** in production
- **Comprehensive tests**

</td>
<td width="50%">

### 🎯 Developer Friendly
- **Simple API** - consistent patterns
- **Auto-detection** - no manual config
- **Rich examples** included
- **Well documented**

</td>
</tr>
</table>

---

## 📦 Installation

```bash
go get github.com/parf/homebase-go-lib
```

---

## 🚀 Quick Start

### Process Compressed Files Automatically

```go
import "github.com/parf/homebase-go-lib/fileiterator"

// Works with .gz, .zst, .lz4, .br, .xz automatically!
fileiterator.IterateLines("access.log.gz", func(line string) error {
    // Process each line
    fmt.Println(line)
    return nil
})
```

### Type-Safe JSON Lines Processing

```go
type User struct {
    ID    int64  `json:"id"`
    Name  string `json:"name"`
    Email string `json:"email"`
}

// Generic type-safe iterator
fileiterator.IterateJSONLTyped("users.jsonl.zst", func(user User) error {
    fmt.Printf("User: %s <%s>\n", user.Name, user.Email)
    return nil
})
```

### Universal Schema Support

```go
// Read ANY schema from Parquet, CSV, JSONL, MsgPack
records, _ := fileiterator.ReadInput("data.parquet")
for _, record := range records {
    // record is map[string]any - works with any field structure!
    fmt.Printf("Record: %v\n", record)
}
```

### Progress Tracking

```go
import "github.com/parf/homebase-go-lib/clistat"

stat := clistat.New(10) // Report every 10 seconds
for i := 0; i < 1_000_000; i++ {
    stat.Hit()  // Auto-reports progress: "45.2K hits/sec (500K total)"
}
stat.Finish()   // Final summary
```

---

## 📋 Supported Formats

### Structured Data Formats

| Format | Extensions | Read | Write | Use Case | Performance |
|--------|-----------|------|-------|----------|-------------|
| 📄 **CSV** | `.csv`, `.tsv` | ✅ | ✅ | Excel compatibility, human-readable | Good |
| 📝 **JSON Lines** | `.jsonl`, `.ndjson` | ✅ | ✅ | Debugging, wide support | Moderate |
| 📊 **Apache Parquet** | `.parquet` | ✅ | ✅ | Analytics, columnar queries | **Excellent** |
| 🔧 **MessagePack** | `.msgpack`, `.mp` | ✅ | ✅ | Binary efficiency, 2x smaller than JSON | Very Good |
| ⚡ **FlatBuffer** | `.fb` | ✅ | ✅ | Zero-copy, fastest reads (3x faster) | **Fastest** |

### Compression Formats (Auto-Detected)

| Compression | Extension | Speed | Ratio | Use Case |
|-------------|-----------|-------|-------|----------|
| ⚡ **LZ4** | `.lz4` | **Fastest** | Good | Real-time processing |
| 🎯 **Zstandard** | `.zst` | Fast | **Excellent** | **Recommended** for most uses |
| 📦 **Gzip** | `.gz` | Moderate | Good | Universal compatibility |
| 🔥 **Brotli** | `.br` | Slow | **Best** | Maximum compression |
| ❄️ **XZ/LZMA** | `.xz` | Very Slow | Excellent | Archive storage |
| 📋 **Zlib** | `.zlib`, `.zz` | Moderate | Good | Legacy support |

**All formats work seamlessly with all compression types!** For example: `.jsonl.zst`, `.csv.gz`, `.parquet.lz4`

---

## 💡 Use Cases

### 🔄 Data Pipeline Processing

```go
// Convert between any formats with automatic compression
input, _ := fileiterator.ReadInput("raw-data.csv.gz")           // CSV + Gzip
fileiterator.WriteParquetAny("processed.parquet.zst", input)    // Parquet + Zstd
```

### 📊 Log Analysis

```go
stat := clistat.New(5)
fileiterator.IterateLines("access.log.gz", func(line string) error {
    if strings.Contains(line, "ERROR") {
        // Process error logs
    }
    stat.Hit()
    return nil
})
stat.Finish()  // "Processed 2.5M lines in 3.2s (781K lines/sec)"
```

### 🗄️ Database ETL

```go
// Extract from CSV, transform, load to database
fileiterator.IterateCSVMap("export.csv.zst", func(row map[string]string) error {
    // Transform data
    user := transformUser(row)

    // Load to database
    return db.Insert(user)
})
```

### 🚀 Batch Processing

```go
// Process millions of records efficiently
fileiterator.IterateParquetAny("events.parquet", func(event map[string]any) error {
    // Process each event with automatic memory management
    return processEvent(event)
})
```

---

## 📚 Core Packages

### 📁 `fileiterator` - Universal File Processing

The heart of homebase-go-lib. Process any file format with automatic compression detection.

#### Key Functions

```go
// Universal I/O
FUOpen(filename)              // Open any file/URL with auto-decompression
FUCreate(filename)            // Create file with auto-compression
ReadInput(filename)           // Read ANY schema to []map[string]any
WriteOutput(filename, data)   // Write ANY schema from []map[string]any

// Line-by-line Processing
IterateLines(filename, func(line string) error)

// Structured Data (Untyped)
IterateJSONL(filename, func(map[string]any) error)
IterateCSVMap(filename, func(map[string]string) error)
IterateMsgPack(filename, func(any) error)
IterateParquetAny(filename, func(map[string]any) error)

// Structured Data (Type-Safe Generics)
IterateJSONLTyped[T](filename, func(T) error)
IterateMsgPackTyped[T](filename, func(T) error)

// Binary Formats
IterateBinaryRecords(filename, recordSize, func([]byte) error)
IterateFlatBufferList(filename, func([]byte) error)
```

**Features:**
- ✅ Automatic compression detection from file extension
- ✅ HTTP/HTTPS URL support
- ✅ Streaming for memory efficiency
- ✅ Progress reporting integration
- ✅ Error handling with context

### 📊 `clistat` - Real-Time Statistics

Track processing progress with automatic hit-rate reporting.

```go
stat := clistat.New(10)  // Report every 10 seconds

for i := 0; i < 1_000_000; i++ {
    // Your processing logic
    stat.Hit()  // Automatically reports: "45.2K hits/sec (500K total)"
}

stat.Finish()  // Final: "Processed 1M items in 22.1s (45.2K/sec)"
```

**Features:**
- ✅ Automatic progress reporting
- ✅ Configurable intervals
- ✅ Hits-per-second calculation
- ✅ Total count tracking
- ✅ Elapsed time reporting

### 🗄️ `sql` - Database Utilities

Efficient database operations with batch processing.

```go
// Batch Insert
inserter := sql.NewBatchInserter(db, "users", []string{"id", "name", "email"}, 1000)
inserter.Add(1, "Alice", "alice@example.com")
inserter.Add(2, "Bob", "bob@example.com")
inserter.Flush()

// Query Iteration
sql.SqlIterator(db, "SELECT * FROM users WHERE active = true", func(row map[string]any) error {
    // Process each row
    return nil
})
```

---

## 🎯 Format Conversion Tools

### Universal Converters (Included)

Located in `cmd/` directory:

#### `any2parquet` - Convert to Apache Parquet

```bash
# Convert any format to Parquet
any2parquet data.jsonl                    # → data.parquet
any2parquet logs.csv.gz                   # → logs.parquet
any2parquet events.msgpack.zst            # → events.parquet
```

#### `any2jsonl` - Convert to JSON Lines

```bash
# Convert any format to human-readable JSONL
any2jsonl data.parquet                    # → data.jsonl
any2jsonl users.csv                       # → users.jsonl
any2jsonl metrics.parquet.zst             # → metrics.jsonl
```

**Standalone Tool:** [any-to-parquet](https://github.com/parf/any-to-parquet) - Optimized Parquet converter

---

## 📈 Performance Benchmarks

Based on 1 million records:

| Format | File Size | Read Time | Write Time | Compression | Best For |
|--------|-----------|-----------|------------|-------------|----------|
| **Parquet** | 44 MB | **0.15s** | **0.46s** | Excellent | **Everything** 🏆 |
| MsgPack.zst | 38 MB | 0.59s | 0.61s | Best | Binary efficiency |
| JSONL.zst | 43 MB | 1.91s | 0.84s | Excellent | Debugging |
| FlatBuffer.lz4 | 66 MB | **0.06s** | 0.42s | Good | Ultra-fast reads |
| CSV.gz | 52 MB | 2.1s | 1.2s | Good | Excel compatibility |
| Plain JSONL | 156 MB | 1.93s | 1.38s | None | Human-readable |

**Winner:** Parquet delivers the best balance of speed, compression, and compatibility.

[Full Benchmark Results](https://github.com/parf/homebase-go-lib/blob/main/benchmarks/serialization-benchmark-result.md) →

---

## 🧪 Examples & Tests

### Running Examples

```bash
# File processing examples
cd examples/fileiterator
go run main.go

# Statistics tracking
cd examples/clistat
go run main.go

# Schema examples (5 different data structures)
cd cmd/examples/schemas
./test-all-schemas.sh
```

### Test Different Schemas

The library works with **ANY schema structure**. See examples:

- [Products (E-commerce)](cmd/examples/schemas/products.jsonl)
- [Sensors (IoT)](cmd/examples/schemas/sensors.jsonl)
- [Users (CSV)](cmd/examples/schemas/users.csv)
- [Logs (Application)](cmd/examples/schemas/logs.jsonl)
- [Transactions (Finance)](cmd/examples/schemas/transactions.jsonl)

[View All Schema Examples →](cmd/examples/schemas/)

---

## 🛠️ Development

### Prerequisites

- Go 1.21 or higher
- Make (optional)

### Build

```bash
make build
```

### Run Tests

```bash
make test
```

### Test Coverage

```bash
make test-coverage
```

### Format & Lint

```bash
make fmt
make lint
```

---

## 📁 Project Structure

```
homebase-go-lib/
├── 📦 fileiterator/       # File processing & format conversion
│   ├── parquet.go         # Apache Parquet support
│   ├── jsonl.go           # JSON Lines processing
│   ├── csv.go             # CSV with auto-detection
│   ├── msgpack.go         # MessagePack binary format
│   ├── genericio.go       # Universal I/O functions
│   └── compression.go     # 7 compression formats
│
├── 📊 clistat/            # Real-time statistics tracking
│   └── clistat.go
│
├── 🗄️ sql/                # Database utilities
│   ├── batch.go           # Batch insert operations
│   └── iterator.go        # Query iteration
│
├── 🎯 cmd/                # Command-line tools
│   ├── any2parquet.go     # Universal → Parquet converter
│   ├── any2jsonl.go       # Universal → JSONL converter
│   └── examples/          # Usage examples
│       └── schemas/       # 5 different schema examples
│
├── 🧪 examples/           # Code examples
├── 📚 docs/               # Documentation
├── 🏗️ benchmarks/         # Performance benchmarks
└── 🧰 testdata/           # Test fixtures
```

---

## 🔑 Key Concepts

### Automatic Compression Detection

```go
// All these work automatically based on file extension:
fileiterator.IterateLines("file.txt")      // Plain text
fileiterator.IterateLines("file.txt.gz")   // Gzip compressed
fileiterator.IterateLines("file.txt.zst")  // Zstandard compressed
fileiterator.IterateLines("file.txt.lz4")  // LZ4 compressed
```

### Universal Schema Support

```go
// No schema definition needed - works with ANY structure!
records, _ := fileiterator.ReadInput("data.csv")
// records[0] might be: {"user_id": 1, "name": "Alice", "age": 28}

records2, _ := fileiterator.ReadInput("sensors.jsonl")
// records2[0] might be: {"sensor": "temp-01", "value": 23.5, "unit": "celsius"}
```

### Type-Safe Generics

```go
// Define your struct
type Product struct {
    ID    int     `json:"id"`
    Name  string  `json:"name"`
    Price float64 `json:"price"`
}

// Get type-safe iteration with Go generics
fileiterator.IterateJSONLTyped("products.jsonl", func(p Product) error {
    fmt.Printf("%s: $%.2f\n", p.Name, p.Price)
    return nil
})
```

---

## 🤝 Contributing

Contributions welcome! Please:

1. 🍴 Fork the repository
2. 🌿 Create a feature branch
3. ✅ Add tests for new functionality
4. 📝 Update documentation
5. 🚀 Submit a pull request

[Report Bug](https://github.com/parf/homebase-go-lib/issues) · [Request Feature](https://github.com/parf/homebase-go-lib/issues)

---

## 📄 License

MIT License - see [LICENSE](LICENSE) file for details

---

## 🔗 Related Projects

- [any-to-parquet](https://github.com/parf/any-to-parquet) - Standalone Parquet converter tool
- [Apache Parquet](https://parquet.apache.org/) - Columnar storage format
- [Apache Arrow](https://arrow.apache.org/) - In-memory columnar data

---

## 🏷️ Keywords

go library, golang, file processing, data pipeline, ETL, compression, gzip, zstd, lz4, parquet, json lines, csv processing, msgpack, data engineering, batch processing, streaming, apache parquet, columnar format, data conversion, format converter, structured data, log processing, statistics tracking, progress reporting, database utilities, sql batch insert, type-safe iterators, go generics, high performance, production ready

---

## ⭐ Star History

If you find this library useful, please [give it a star](https://github.com/parf/homebase-go-lib/stargazers)! ⭐

[![Star History Chart](https://api.star-history.com/svg?repos=parf/homebase-go-lib&type=Date)](https://star-history.com/#parf/homebase-go-lib&Date)

---

<div align="center">

**Built with ❤️ for the Go and data engineering community**

[Documentation](https://github.com/parf/homebase-go-lib/tree/main/docs) · [Examples](https://github.com/parf/homebase-go-lib/tree/main/examples) · [Benchmarks](https://github.com/parf/homebase-go-lib/tree/main/benchmarks)

</div>
