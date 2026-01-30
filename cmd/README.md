# Universal Data Conversion Utilities

High-performance data format conversion and database import/export utilities with comprehensive SQL and compression support.

## Export Tools

Convert data from files or databases to various output formats:

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

# With additional compression (optional)
./any2parquet data.jsonl data.parquet.lz4   # → data.parquet.lz4
./any2parquet data.csv data.parquet.zst     # → data.parquet.zst

# Query MySQL/PostgreSQL databases to Parquet
./any2parquet --dsn="user:pass@localhost" --sql="SELECT * FROM users"
./any2parquet --dsn="user:pass@localhost" --table="mydb.users"
./any2parquet --driver=postgre --dsn="host=pg user=x password=y" --table="public.orders"
```

**Supported inputs:** JSONL, CSV, MsgPack, Parquet, **SQL databases** (MySQL, PostgreSQL)
**Performance:** 0.15s read, 0.46s write, 44MB for 1M records
**Best for:** Everything - APIs, analytics, data warehouses, ML pipelines, **database exports**

### any2jsonl
Convert any format to JSONL (human-readable debugging format).

```bash
# Build once
go build any2jsonl.go

# Convert files to JSONL
./any2jsonl data.parquet              # → data.jsonl
./any2jsonl data.csv -                # → stdout

# With compression (specify in output filename)
./any2jsonl data.parquet data.jsonl.zst   # → data.jsonl.zst (RECOMMENDED)
./any2jsonl data.csv data.jsonl.gz        # → data.jsonl.gz
./any2jsonl data.msgpack data.jsonl.lz4   # → data.jsonl.lz4

# Query MySQL/PostgreSQL databases
./any2jsonl --dsn="user:pass@localhost" --sql="SELECT * FROM users"
./any2jsonl --dsn="user:pass@localhost" --table="mydb.users"
./any2jsonl --driver=postgre --dsn="host=localhost user=x password=y" --table="public.orders"
```

**Supported inputs:** Parquet, JSONL, MsgPack, CSV, **SQL databases** (MySQL, PostgreSQL)
**Performance:** 1.91s read, 0.84s write, 43MB with Zstd
**Best for:** Debugging, data inspection, text processing with grep/jq, **database exports**

### any2csv
Convert any format to CSV (spreadsheet-compatible format).

```bash
# Build once
go build any2csv.go

# Convert files to CSV
./any2csv data.parquet                # → data.csv
./any2csv data.jsonl -                # → stdout
./any2csv data.msgpack output.csv     # → output.csv

# Query MySQL/PostgreSQL databases to CSV
./any2csv --dsn="user:pass@localhost" --sql="SELECT * FROM users" users.csv
./any2csv --dsn="root:pass@localhost" --table="mydb.orders" - | head
./any2csv --driver=postgre --dsn="user:pass@pghost" --table="public.logs" logs.csv
```

**Supported inputs:** Parquet, JSONL, MsgPack, CSV, **SQL databases** (MySQL, PostgreSQL)
**Column order:** Alphabetically sorted for consistency
**Best for:** Excel/Google Sheets, spreadsheet analysis, data exchange, reporting

### any2db
Import data from files or databases into database tables.

```bash
# Build once
go build any2db.go

# Import files to database
./any2db --dsn="root:pass@localhost/mydb" data.csv users
./any2db --dsn="root:pass@localhost/mydb" data.parquet orders
./any2db --driver=postgre --dsn="user:pass@pghost/mydb" data.jsonl public.events

# Copy/transform data between tables
./any2db --dsn="root:pass@localhost/mydb" --table="old_users" new_users
./any2db --dsn="root:pass@localhost/mydb" --sql="SELECT * FROM orders WHERE date>'2024-01-01'" orders_2024

# Import with custom batch size
./any2db --dsn="root:pass@localhost/mydb" --batch=5000 large_file.jsonl.zst events
```

**Features:**
- Automatically creates destination table if not exists
- Infers column types from data (BIGINT, DOUBLE, TEXT, BOOLEAN)
- Batch inserts for high performance (default: 1000 records)
- Auto-escapes values to prevent SQL injection
- Supports MySQL and PostgreSQL

**Supported inputs:** Parquet, JSONL, CSV, MsgPack, **SQL queries** (for table copying)
**Best for:** Database imports, ETL pipelines, table copying, data migration

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
./any2parquet mydata.csv.gz mydata.parquet.lz4
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
./any2jsonl --dsn="user:pass@host" --table="geo.zip" zipcodes.jsonl.zst
./any2parquet --dsn="user:pass@host" --sql="SELECT * FROM orders WHERE date > '2024-01-01'" recent_orders.parquet
```

### SQL Flags

- `--sql="SELECT * FROM table"` - SQL query to execute
- `--table="schema.table"` - Alternative to --sql (generates `SELECT * FROM table`)
- `--driver=mysql` - Database driver: `mysql` (default) or `postgre`
- `--dsn="connection-string"` - Database connection string

Output file name is a positional argument. If omitted or "-", outputs to stdout.

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
./any2parquet --dsn="root:password@localhost" --table="mydb.users" users.parquet

# PostgreSQL with complex query
./any2jsonl --driver=postgre \
  --dsn="host=pg.example.com port=5432 user=analyst password=secret dbname=analytics" \
  --sql="SELECT customer_id, SUM(amount) FROM orders GROUP BY customer_id" \
  customer_totals.jsonl.zst

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

### 📄 Use JSONL (any2jsonl) when:
- Debugging/inspecting data (human-readable)
- Need text processing tools (grep, jq, sed, awk)
- Use .zst extension for compression: 1.91s read, 43MB (vs 1.93s, 156MB plain)
- Never use for production - much slower than binary formats

### 📊 Use CSV (any2csv) when:
- Working with Excel or Google Sheets
- Need spreadsheet-compatible format
- Sharing data with non-technical users
- Generating reports for business users
- Compatible with all spreadsheet and BI tools

### 💾 Use any2db when:
- Loading data into MySQL or PostgreSQL databases
- Building ETL/data pipelines
- Migrating data between databases
- Creating database tables from files
- Bulk importing with automatic schema creation

## Performance Comparison (1M records)

| Format | Read | Write | Total | Size | Best For |
|--------|------|-------|-------|------|----------|
| **Parquet** 🏆 | **0.15s** | **0.46s** | **0.61s** | **44MB** | **Everything** |
| JSONL + Zstd | 1.91s | 0.84s | 2.75s | 43MB | Debugging |

Full benchmarks: [serialization-benchmark-result.md](../benchmarks/serialization-benchmark-result.md)

## Example Data

Pre-generated sample files in `examples/` directory:
- 100 records with realistic fake data (gofakeit v7)
- All formats: JSONL, CSV, Parquet
- All compressions: .gz, .zst, .lz4
- Total size: 108KB
- See `examples/README.md` for usage examples

## Notes

- **Schema-agnostic:** Converters automatically handle ANY schema structure
- **SQL support:** Direct export from MySQL and PostgreSQL databases
- **Input compression:** Auto-detected by file extension
- **Output filenames:** Auto-generated if not specified (from table name for SQL queries)
- **Backward compatible:** Existing file-based conversion continues to work unchanged
