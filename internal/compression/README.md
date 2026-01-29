# Internal Compression Package

**Note:** This is an internal package. For public API, use the `fileiterator` package which provides all file loading and iteration functionality.

The internal `compression` package contains legacy compression format loaders kept for backward compatibility within the module.

## Functions

### Binary File Loaders

#### LoadBinGzFile(filename string, dest *[]byte)

Load gzipped binary file explicitly.

```go
import "github.com/parf/homebase-go-lib/internal/compression"

var data []byte
compression.LoadBinGzFile("file.bin.gz", &data)
```

#### LoadBinZstdFile(filename string, dest *[]byte)

Load zstd-compressed binary file explicitly.

```go
var data []byte
compression.LoadBinZstdFile("file.bin.zst", &data)
```

#### LoadBinLz4File(filename string, dest *[]byte)

Load LZ4-compressed binary file explicitly.

```go
var data []byte
compression.LoadBinLz4File("file.bin.lz4", &data)
```

### Text File Loaders

#### LoadLinesGzFile(filename string, processor func(string))

Process lines in a gzipped text file.

```go
compression.LoadLinesGzFile("file.txt.gz", func(line string) {
    fmt.Println(line)
})
```

#### LoadLinesZstdFile(filename string, processor func(string))

Process lines in a zstd-compressed text file.

```go
compression.LoadLinesZstdFile("file.txt.zst", func(line string) {
    fmt.Println(line)
})
```

#### IterateLinesLz4(filename string, processor func(string))

Process lines in an LZ4-compressed text file.

```go
compression.IterateLinesLz4("file.txt.lz4", func(line string) {
    fmt.Println(line)
})
```

### Special Format Loaders

#### LoadIDTabGzFile(filename string, processor func(int32, string))

Process tab-separated ID-name pairs from a gzipped file. IDs are parsed as hexadecimal int32, names are converted to lowercase.

```go
compression.LoadIDTabGzFile("ids.tab.gz", func(id int32, name string) {
    fmt.Printf("ID: %d, Name: %s\n", id, name)
})
```

## Migration Notice

**All functions from this package are now available in `fileiterator` package.**

Use `github.com/parf/homebase-go-lib/fileiterator` instead:
- LoadBinGzFile, LoadBinZstdFile, LoadBinLz4File
- IterateLinesGz, IterateLinesZstd, IterateLinesLz4
- IterateIDTabFile (auto-detects all compression formats)
- IterateBinaryRecords and explicit format iterators

This internal package is kept for backward compatibility within the module only.

## Example

```go
package main

import (
    "fmt"
    "github.com/parf/homebase-go-lib/compression"
)

func main() {
    // Load binary gzip file
    var data []byte
    compression.LoadBinGzFile("data.bin.gz", &data)
    fmt.Printf("Loaded %d bytes\n", len(data))

    // Process text lines from zstd file
    compression.LoadLinesZstdFile("log.txt.zst", func(line string) {
        fmt.Println(line)
    })

    // Process ID-name pairs
    compression.LoadIDTabGzFile("names.tab.gz", func(id int32, name string) {
        fmt.Printf("%x: %s\n", id, name)
    })
}
```

## Binary Record Iterators

### IterateBinaryRecords(filename string, recordSize int, processor func([]byte))

Iterate over fixed-size binary records with **automatic compression detection** by extension (.gz, .zst, .zlib/.zz).

```go
compression.IterateBinaryRecords("data.bin.gz", 10, func(record []byte) {
    // process 10-byte record - automatically decompressed
})

compression.IterateBinaryRecords("data.bin.zst", 10, func(record []byte) {
    // process 10-byte record - automatically decompressed
})

compression.IterateBinaryRecords("data.bin.zlib", 10, func(record []byte) {
    // process 10-byte record - automatically decompressed
})
```

### Explicit Format Iterators

For explicit compression format control:

```go
// Gzip
compression.IterateGzipRecords("data.bin.gz", 10, processor)

// Zstd
compression.IterateZstdRecords("data.bin.zst", 10, processor)

// Zlib (RFC 1950)
compression.IterateZlibRecords("data.bin.zlib", 10, processor)
```

## Supported Formats

- **Gzip** (.gz) - Standard gzip compression (RFC 1952)
- **Zstd** (.zst) - Zstandard compression (modern, faster)
- **Zlib** (.zlib, .zz) - Zlib compression (RFC 1950)
- **LZ4** (.lz4) - Fast compression algorithm

## URL Support

All loaders support both local files and HTTP URLs:

```go
compression.LoadBinGzFile("http://example.com/data.bin.gz", &data)
compression.IterateBinaryRecords("http://example.com/records.bin.zst", 10, processor)
```
