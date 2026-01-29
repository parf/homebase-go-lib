# Conversion Tools

High-performance data format conversion utilities for optimal read performance.

## Parquet Conversion Tools (Analytics Format)

### jsonl2parquet ⭐
Convert JSONL to Parquet format (best for analytics, 17x faster reads than JSONL!)

```bash
go run jsonl2parquet.go -input data.jsonl.gz -output data.parquet

# Or use the bash wrapper:
./jsonl2parquet.sh -input data.jsonl.gz -output data.parquet
```

**Performance:** 17x faster reads than JSONL, excellent for analytics
**Best for:** Data warehouses, BI tools, SQL queries, Apache Spark, DuckDB

### csv2parquet
Convert CSV to Parquet format

```bash
go run csv2parquet.go -input data.csv.gz -output data.parquet

# Or use the bash wrapper:
./csv2parquet.sh -input data.csv.gz -output data.parquet

# With options:
./csv2parquet.sh -input data.tsv -output data.parquet -delimiter=tab -header=true
```

**Supports:**
- Automatic compression detection (.gz, .zst, .lz4, .br, .xz)
- Custom delimiters (comma, tab, pipe, semicolon)
- Header detection
- Automatic type inference (int64, float64, bool, string)

## FlatBuffer Conversion Tools (High-Performance APIs)

### jsonl2fb-lz4
Convert JSONL to FlatBuffer LZ4 format (10x faster reads than JSONL!)

```bash
go run jsonl2fb-lz4.go -input data.jsonl -output data.fb.lz4

# Or use the bash wrapper:
./jsonl2fb-lz4.sh data.jsonl data.fb.lz4
```

**Performance:** 3x faster reads than JSONL, good compression

### fb-lz42jsonl
Convert FlatBuffer LZ4 back to JSONL

```bash
go run fb-lz42jsonl.go -input data.fb.lz4 -output data.jsonl
```

### parquet2fb-lz4
Convert Parquet to FlatBuffer LZ4 format

```bash
go run parquet2fb-lz4.go -input data.parquet -output data.fb.lz4
```

**Why convert from Parquet?**
- Parquet: Great for columnar analytics, slow for row-oriented reads
- FlatBuffer + LZ4: 3-10x faster for row-oriented access (APIs, services)

## Format Selection Guide

### Use Parquet when:
- Running analytical queries (SQL, aggregations, filters)
- Working with BI tools (Tableau, Looker, Power BI)
- Processing with big data frameworks (Spark, DuckDB, Pandas)
- Need columnar access patterns
- **Read performance: 0.11s for 1M records (17x faster than JSONL)**

### Use FlatBuffer + LZ4 when:
- Building APIs serving data
- Real-time systems requiring low latency
- Row-oriented access patterns
- Message queues, caches, RPCs
- **Read performance: 0.19s for 1M records (10x faster than JSONL)**

### Use JSONL when:
- Need human-readable format
- Debugging/inspecting data
- Simple line-by-line processing
- Text-based tools (grep, sed, awk)
- **Read performance: 1.82s for 1M records**

## Performance Comparison

| Format | Read (1M records) | Best For |
|--------|------------------|----------|
| **Parquet** | **0.11s** 🏆 | Analytics, data warehouses, SQL queries |
| **FlatBuffer + LZ4** | **0.19s** ⭐ | APIs, services, real-time systems |
| MsgPack + Zstd | 0.57s | Binary serialization |
| JSONL + Zstd | 1.82s | Human-readable logs |

See full benchmarks in `/benchmarks/serialization-benchmark-result.md`
