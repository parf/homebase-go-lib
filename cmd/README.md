# Universal Format Converters

High-performance data format conversion utilities with comprehensive compression support.

## Tools

### any2parquet 🏆 (RECOMMENDED)
Convert any format to Parquet - the best overall format for everything.

```bash
# Build once
go build any2parquet.go

# Convert any format to Parquet
./any2parquet data.jsonl              # → data.parquet
./any2parquet data.csv.gz             # → data.parquet
./any2parquet data.msgpack.zst        # → data.parquet

# With LZ4 compression (optional, even smaller)
./any2parquet --lz4 data.jsonl        # → data.parquet.lz4
```

**Supported inputs:** JSONL, CSV, MsgPack, FlatBuffer (all with compression)
**Performance:** 0.15s read, 0.46s write, 44MB for 1M records
**Best for:** Everything - APIs, analytics, data warehouses, ML pipelines

### any2fb
Convert any format to FlatBuffer (fastest reads, but larger files).

```bash
# Build once
go build any2fb.go

# Convert any format to FlatBuffer
./any2fb data.jsonl                   # → data.fb
./any2fb data.parquet                 # → data.fb
./any2fb data.csv.gz                  # → data.fb

# With LZ4 compression (recommended)
./any2fb --lz4 data.jsonl             # → data.fb.lz4
```

**Supported inputs:** JSONL, CSV, MsgPack, Parquet (all with compression)
**Performance:** 0.06s read, 0.78s write, 160MB plain / 66MB with LZ4
**Best for:** Hot data paths where read speed is absolutely critical

### any2jsonl
Convert any format to JSONL (human-readable debugging format).

```bash
# Build once
go build any2jsonl.go

# Convert to plain JSONL
./any2jsonl data.parquet              # → data.jsonl

# With compression (recommended)
./any2jsonl --zst data.parquet        # → data.jsonl.zst (RECOMMENDED)
./any2jsonl --gz data.fb              # → data.jsonl.gz
./any2jsonl --lz4 data.msgpack        # → data.jsonl.lz4
```

**Supported inputs:** Parquet, FlatBuffer, MsgPack, CSV (all with compression)
**Performance:** 1.91s read, 0.84s write, 43MB with Zstd
**Best for:** Debugging, data inspection, text processing with grep/jq

## Quick Start

```bash
# 1. See example usage
cd examples/
cat README.md

# 2. Test with sample data (100 records)
cd ..
./any2parquet examples/sample-data.jsonl examples/test.parquet
./any2jsonl examples/test.parquet examples/output.jsonl

# 3. Convert your own data
./any2parquet --lz4 mydata.csv.gz mydata.parquet.lz4
```

## Compression Options

All converters auto-detect input compression and support output compression:

- **Gzip (.gz)** - Standard, widely supported, slow
- **Zstandard (.zst)** - RECOMMENDED: best balance of speed & compression
- **LZ4 (.lz4)** - Fastest compression, moderate compression ratio
- **Brotli (.br)** - Best compression, very slow
- **XZ (.xz)** - Excellent compression, extremely slow (avoid)

## Format Selection Guide

### 🏆 Use Parquet (any2parquet) for:
- **Everything** - APIs, analytics, data warehouses, ML pipelines
- Industry standard (Spark, DuckDB, Pandas, Arrow, all major tools)
- Best overall: 0.15s read, 0.46s write, 44MB for 1M records
- Columnar format: Extremely fast for queries, aggregations, filters

### ⚡ Use FlatBuffer (any2fb) when:
- Read speed is absolutely critical and storage is unlimited
- Hot data paths with very high-frequency reads
- 0.06s read (2.5x faster than Parquet), but 160MB uncompressed
- Use --lz4 flag: 0.21s read, 66MB (still 1.5x larger than Parquet)
- Only use if Parquet isn't available in your stack

### 📄 Use JSONL (any2jsonl) when:
- Debugging/inspecting data (human-readable)
- Need text processing tools (grep, jq, sed, awk)
- Always use --zst flag: 1.91s read, 43MB (vs 1.93s, 156MB plain)
- Never use for production - much slower than binary formats

## Performance Comparison (1M records)

| Format | Read | Write | Total | Size | Best For |
|--------|------|-------|-------|------|----------|
| **Parquet** 🏆 | **0.15s** | **0.46s** | **0.61s** | **44MB** | **Everything** |
| FlatBuffer Plain | 0.06s | 0.78s | 0.84s | 160MB | Fastest reads only |
| FlatBuffer + LZ4 | 0.21s | 1.11s | 1.32s | 66MB | Fast reads (no Parquet) |
| JSONL + Zstd | 1.91s | 0.84s | 2.75s | 43MB | Debugging |

Full benchmarks: [serialization-benchmark-result.md](../benchmarks/serialization-benchmark-result.md)

## Example Data

Pre-generated sample files in `examples/` directory:
- 100 records with realistic fake data (gofakeit v7)
- All formats: JSONL, CSV, Parquet, FlatBuffer
- All compressions: .gz, .zst, .lz4
- Total size: 108KB
- See `examples/README.md` for usage examples

## Building

```bash
# Build all converters
go build any2parquet.go
go build any2fb.go
go build any2jsonl.go

# Or build specific converter
go build any2parquet.go
```

Binaries are ~46MB each (includes all format libraries).

## Notes

- All converters assume TestRecord schema (id, name, email, age, score, active, category, timestamp)
- Input compression auto-detected by file extension
- Output filenames auto-generated if not specified
- For custom schemas, modify the converter Go files
