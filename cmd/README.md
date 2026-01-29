# Conversion Tools

High-performance data format conversion utilities.

## Tools

### jsonl2fb-lz4
Convert JSONL to FlatBuffer LZ4 format (3x faster reads!)

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

## Use Case

Convert your data to FlatBuffer + LZ4 format for production systems where read performance is critical.

**Benefits:**
- 3x faster reads than JSONL
- 2x faster reads than MsgPack
- Zero-copy deserialization
- Good compression with LZ4

**When to use:**
- APIs serving data
- Real-time systems
- Caches
- Message queues
- Any read-heavy workload

See benchmarks in `/benchmarks/serialization-benchmark-result.md`
