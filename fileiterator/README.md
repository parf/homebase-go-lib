# File Iterator Package

The `fileiterator` package provides high-level iterators for structured file formats (JSONL, CSV) with automatic compression detection.

## Features

- **JSONL (JSON Lines)** support with typed and untyped parsing
- **CSV** support with flexible options and map-based iteration
- **Automatic compression** detection (.gz, .zst)
- **URL support** - works with both local files and HTTP/HTTPS URLs
- **Error handling** with detailed error messages (line/row numbers)
- **Progress tracking** - prints row/line counts

## JSONL (JSON Lines) Support

### IterateJSONL - Untyped

Process JSONL files with generic `map[string]any` objects:

```go
import "github.com/parf/homebase-go-lib/fileiterator"

err := fileiterator.IterateJSONL("data.jsonl.gz", func(obj map[string]any) error {
    fmt.Printf("ID: %v, Name: %v\n", obj["id"], obj["name"])
    return nil
})
```

### IterateJSONLTyped - Type-Safe

Process JSONL files with strongly-typed structs:

```go
type User struct {
    ID   int    `json:"id"`
    Name string `json:"name"`
    Email string `json:"email"`
}

err := fileiterator.IterateJSONLTyped("users.jsonl", func(user User) error {
    fmt.Printf("User: %s <%s>\n", user.Name, user.Email)
    return nil
})
```

**Supports:**
- Plain JSONL files: `data.jsonl`
- Gzipped: `data.jsonl.gz`
- Zstd: `data.jsonl.zst`
- URLs: `http://example.com/data.jsonl.gz`

## CSV Support

### CSVOptions

Configure CSV parsing behavior:

```go
opts := fileiterator.DefaultCSVOptions()
opts.Comma = '\t'              // Tab-separated (TSV)
opts.SkipHeader = true          // Skip first row
opts.TrimLeadingSpace = true    // Trim spaces
opts.Comment = '#'              // Comment character
```

### IterateCSV - Array-based

Process CSV files row by row as string slices:

```go
opts := fileiterator.DefaultCSVOptions()
opts.SkipHeader = true

err := fileiterator.IterateCSV("data.csv.gz", opts, func(row []string) error {
    fmt.Printf("Name: %s, Age: %s\n", row[0], row[1])
    return nil
})
```

### IterateCSVMap - Map-based

Process CSV files with header row, get each row as a map:

```go
err := fileiterator.IterateCSVMap("users.csv", fileiterator.DefaultCSVOptions(), func(row map[string]string) error {
    fmt.Printf("Name: %s, Email: %s\n", row["name"], row["email"])
    return nil
})
```

**Supports:**
- Plain CSV files: `data.csv`
- Gzipped: `data.csv.gz`
- Zstd: `data.csv.zst`
- URLs: `http://example.com/data.csv.gz`
- TSV (tab-separated) - set `Comma` to `'\t'`
- Custom delimiters

## Examples

### JSONL from URL

```go
err := fileiterator.IterateJSONL("https://example.com/logs.jsonl.gz", func(log map[string]any) error {
    fmt.Printf("Timestamp: %v, Level: %v\n", log["timestamp"], log["level"])
    return nil
})
```

### CSV with Custom Delimiter

```go
opts := fileiterator.DefaultCSVOptions()
opts.Comma = '|'  // Pipe-separated
opts.SkipHeader = true

err := fileiterator.IterateCSV("data.psv", opts, func(row []string) error {
    // Process pipe-separated values
    return nil
})
```

### Typed JSONL Processing

```go
type LogEntry struct {
    Timestamp string `json:"timestamp"`
    Level     string `json:"level"`
    Message   string `json:"message"`
}

err := fileiterator.IterateJSONLTyped("logs.jsonl.zst", func(entry LogEntry) error {
    if entry.Level == "ERROR" {
        fmt.Printf("[%s] %s\n", entry.Timestamp, entry.Message)
    }
    return nil
})
```

### CSV Map with Progress

```go
count := 0
err := fileiterator.IterateCSVMap("users.csv.gz", fileiterator.DefaultCSVOptions(), func(row map[string]string) error {
    count++
    if count%1000 == 0 {
        fmt.Printf("Processed %d users...\n", count)
    }
    return nil
})
fmt.Printf("Total users: %d\n", count)
```

## Error Handling

All iterator functions return detailed errors with line/row numbers:

```go
err := fileiterator.IterateJSONL("data.jsonl", func(obj map[string]any) error {
    if obj["id"] == nil {
        return fmt.Errorf("missing required field: id")
    }
    return nil
})

if err != nil {
    // Error format: "line 42: processor error: missing required field: id"
    log.Fatal(err)
}
```

## Performance

- Streaming processing - low memory footprint
- Automatic buffering for optimal performance
- Progress messages to stdout
- Suitable for large files (GB+)

## Compression Support

All functions automatically detect compression by file extension:
- **Gzip** (.gz) - Standard gzip compression
- **Zstd** (.zst) - Modern, faster compression

No special code needed - just use compressed files directly.

## URL Support

All functions work with HTTP/HTTPS URLs:
- Automatically fetches and streams content
- Works with compressed URLs
- No temporary files created

## When to Use

**Use this package when:**
- Processing structured data files (JSONL, CSV)
- Files can be large (streaming processing)
- Need automatic compression handling
- Want type-safe JSONL parsing
- Processing data from URLs

**Use other packages when:**
- Need random access to file contents
- Working with binary formats (use `compression.IterateBinaryRecords`)
- Custom parsing requirements
