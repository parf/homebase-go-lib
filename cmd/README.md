# Universal Format Converters

High-performance data format conversion utilities with comprehensive compression support.

## Tools

### any2parquet 🏆 (RECOMMENDED)
Convert any format to Parquet - the best overall format for everything.

```bash
# Build once
go build any2parquet.go

# Convert files to Parquet
./any2parquet data.jsonl              # → data.parquet
./any2parquet data.csv.gz             # → data.parquet
./any2parquet data.msgpack.zst        # → data.parquet

# With LZ4 compression (optional, even smaller)
./any2parquet --lz4 data.jsonl        # → data.parquet.lz4

# Query MySQL/PostgreSQL databases to Parquet
./any2parquet --dsn="user:pass@localhost" --sql="SELECT * FROM users"
./any2parquet --dsn="user:pass@localhost" --table="mydb.users"
./any2parquet --driver=postgre --dsn="host=pg user=x password=y" --table="public.orders"
```

**Supported inputs:** JSONL, CSV, MsgPack, FlatBuffer, **SQL databases** (MySQL, PostgreSQL)
**Performance:** 0.15s read, 0.46s write, 44MB for 1M records
**Best for:** Everything - APIs, analytics, data warehouses, ML pipelines, **database exports**

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

# Convert files to plain JSONL
./any2jsonl data.parquet              # → data.jsonl

# With compression (recommended)
./any2jsonl --zst data.parquet        # → data.jsonl.zst (RECOMMENDED)
./any2jsonl --gz data.fb              # → data.jsonl.gz
./any2jsonl --lz4 data.msgpack        # → data.jsonl.lz4

# Query MySQL/PostgreSQL databases
./any2jsonl --dsn="user:pass@localhost" --sql="SELECT * FROM users"
./any2jsonl --dsn="user:pass@localhost" --table="mydb.users"
./any2jsonl --driver=postgre --dsn="host=localhost user=x password=y" --table="public.orders"
```

**Supported inputs:** Parquet, FlatBuffer, MsgPack, CSV, **SQL databases** (MySQL, PostgreSQL)
**Performance:** 1.91s read, 0.84s write, 43MB with Zstd
**Best for:** Debugging, data inspection, text processing with grep/jq, **database exports**

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

## SQL Database Support 🆕

Export data directly from MySQL or PostgreSQL databases to JSONL or Parquet formats.

### Basic Usage

```bash
# MySQL with simplified DSN (auto-expanded to full format)
./any2jsonl --dsn="user:pass@host" --sql="SELECT * FROM users"
./any2parquet --dsn="root:password@localhost" --table="mydb.orders"

# PostgreSQL
./any2jsonl --driver=postgre --dsn="host=pg user=x password=y dbname=db" --sql="SELECT * FROM logs"
./any2parquet --driver=postgre --dsn="user:pass@pghost:5432" --table="public.events"

# With custom output name and compression
./any2jsonl --dsn="user:pass@host" --table="geo.zip" --name="zipcode s.jsonl.zst"
./any2parquet --dsn="user:pass@host" --sql="SELECT * FROM orders WHERE date > '2024-01-01'" --name="recent_orders.parquet"
```

### SQL Flags

- `--sql="SELECT * FROM table"` - SQL query to execute
- `--table="schema.table"` - Alternative to --sql (generates `SELECT * FROM table`)
- `--name=output-file` - Output file name (auto-generated from table name if omitted)
- `--driver=mysql` - Database driver: `mysql` (default) or `postgre`
- `--dsn="connection-string"` - Database connection string

### DSN Formats

**MySQL:**
```bash
# Standard format
user:password@tcp(host:3306)/database

# Simplified format (auto-expanded)
user:password@host
```

**PostgreSQL:**
```bash
# Standard format
host=localhost port=5432 user=myuser password=mypass dbname=mydb sslmode=disable

# Simplified format (auto-expanded)
user:password@host:5432
```

### Examples

```bash
# Export 1000 ZIP codes to JSONL
./any2jsonl --dsn="parf:mv700@hdb3" --sql="select * from geo.zip limit 1000"

# Export users table to compressed Parquet (83% smaller than JSONL!)
./any2parquet --dsn="root:password@localhost" --table="mydb.users" --name="users.parquet"

# PostgreSQL with complex query
./any2jsonl --driver=postgre \
  --dsn="host=pg.example.com port=5432 user=analyst password=secret dbname=analytics" \
  --sql="SELECT customer_id, SUM(amount) FROM orders GROUP BY customer_id" \
  --name="customer_totals.jsonl.zst"

# Export and immediately analyze with jq
./any2jsonl --dsn="user:pass@host" --sql="SELECT * FROM logs LIMIT 100" | jq '.[] | select(.level == "ERROR")'
```

### Performance

Real-world test with `geo.zip` table (1000 records):
- **JSONL output:** 97.8 KB
- **Parquet output:** 16.9 KB (83% smaller!)
- Both formats preserve all data and support compression

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

- **Schema-agnostic:** Converters automatically handle ANY schema structure
- **SQL support:** Direct export from MySQL and PostgreSQL databases
- **Input compression:** Auto-detected by file extension
- **Output filenames:** Auto-generated if not specified (from table name for SQL queries)
- **Backward compatible:** Existing file-based conversion continues to work unchanged
