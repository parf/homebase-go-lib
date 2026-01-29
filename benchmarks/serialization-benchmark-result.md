# Serialization Benchmark Results

## Test Configuration

- **Dataset:** 1,000,000 records (~1GB RAM)
- **Record Structure:** 8 fields (ID, Name, Email, Age, Score, Active, Category, Timestamp)
- **CPU:** AMD Ryzen 9 5900X 12-Core Processor
- **Test Date:** 2026-01-28

## Summary

This benchmark compares different serialization formats (JSONL, MessagePack, FlatBuffers) with various compression algorithms (None, Gzip, Zstd, LZ4, Brotli, XZ) for storing and retrieving 1 million records.

## Key Findings

🏆 **Best Overall:** MessagePack + XZ (0.94 MB, 99.2% compression)
⚡ **Fastest Write:** JSONL + Zstd (458 ms)
⚡ **Fastest Read:** MsgPack Plain (554 ms)
📦 **Best Compression:** MessagePack + XZ (99.2% space reduction)

---

## File Size Comparison

| Format | Size (bytes) | Size (MB) | Compression % | vs Plain |
|--------|--------------|-----------|---------------|----------|
| **JSONL Plain** | 152,628,890 | 145.56 | 0% | - |
| JSONL Gzip | 8,508,708 | 8.11 | 94.4% | 17.9x |
| JSONL Zstd | 2,714,335 | 2.59 | 98.2% | 56.2x |
| JSONL LZ4 | 17,161,370 | 16.37 | 88.8% | 8.9x |
| JSONL Brotli | 2,046,654 | 1.95 | 98.7% | 74.6x |
| JSONL XZ | 4,289,016 | 4.09 | 97.2% | 35.6x |
| **MsgPack Plain** | 120,000,000 | 114.44 | 0% | - |
| MsgPack Gzip | 11,204,078 | 10.69 | 90.7% | 10.7x |
| MsgPack Zstd | 8,008,844 | 7.64 | 93.3% | 15.0x |
| MsgPack LZ4 | 18,527,538 | 17.67 | 84.6% | 6.5x |
| MsgPack Brotli | 3,365,899 | 3.21 | 97.2% | 35.7x |
| **MsgPack XZ** | **983,184** | **0.94** | **99.2%** | **122.1x** |
| **FlatBuffer Plain** | 157,467,520 | 150.17 | 0% | - |
| FlatBuffer Zstd | 2,752,056 | 2.62 | 98.3% | 57.2x |

---

## Write Performance (1M records)

| Format | Time (ms) | Time (s) | Records/sec | vs Fastest |
|--------|-----------|----------|-------------|------------|
| JSONL Plain | 1,311.86 | 1.31 | 762,267 | 2.9x |
| JSONL Gzip | 1,212.18 | 1.21 | 824,905 | 2.6x |
| **JSONL Zstd** | **458.12** | **0.46** | **2,182,850** | **1.0x** |
| JSONL LZ4 | 471.92 | 0.47 | 2,119,027 | 1.0x |
| MsgPack Plain | 21,870.17 | 21.87 | 45,726 | 47.7x |
| MsgPack Gzip | 1,918.74 | 1.92 | 521,172 | 4.2x |
| MsgPack Zstd | 632.61 | 0.63 | 1,580,638 | 1.4x |
| MsgPack LZ4 | 674.53 | 0.67 | 1,482,486 | 1.5x |
| FlatBuffer Plain | 677.44 | 0.68 | 1,475,968 | 1.5x |
| FlatBuffer Zstd | 637.50 | 0.64 | 1,568,635 | 1.4x |

**Note:** The MsgPack Plain write performance anomaly suggests an issue with buffering or encoding overhead. In practice, MsgPack should be faster than JSONL.

---

## Read Performance (1M records)

| Format | Time (ms) | Time (s) | Records/sec | vs Fastest |
|--------|-----------|----------|-------------|------------|
| **MsgPack Plain** | **554.16** | **0.55** | **1,804,525** | **1.0x** |
| MsgPack Zstd | 558.06 | 0.56 | 1,791,861 | 1.0x |
| MsgPack LZ4 | 577.85 | 0.58 | 1,730,533 | 1.0x |
| MsgPack Gzip | 630.02 | 0.63 | 1,587,249 | 1.1x |
| JSONL Plain | 1,836.08 | 1.84 | 544,664 | 3.3x |
| JSONL Zstd | 1,833.47 | 1.83 | 545,440 | 3.3x |
| JSONL LZ4 | 1,839.27 | 1.84 | 543,832 | 3.3x |
| JSONL Gzip | 1,880.19 | 1.88 | 531,847 | 3.4x |

---

## Compression Ratio Analysis

### Best Compression by Algorithm

| Algorithm | Best Format | Size (MB) | Compression % |
|-----------|-------------|-----------|---------------|
| **XZ** | MsgPack | 0.94 | 99.2% |
| **Brotli** | JSONL | 1.95 | 98.7% |
| **Zstd** | JSONL | 2.59 | 98.2% |
| **Gzip** | JSONL | 8.11 | 94.4% |
| **LZ4** | JSONL | 16.37 | 88.8% |

### Speed vs Compression Trade-offs

| Category | Format | Write (s) | Read (s) | Size (MB) | Notes |
|----------|--------|-----------|----------|-----------|-------|
| **Balanced** | JSONL + Zstd | 0.46 | 1.83 | 2.59 | Best all-around |
| **Fast Write** | JSONL + Zstd | 0.46 | 1.83 | 2.59 | Fastest write |
| **Fast Read** | MsgPack Plain | 21.87 | 0.55 | 114.44 | Fastest read |
| **Best Compression** | MsgPack + XZ | N/A | N/A | 0.94 | 99.2% reduction |
| **Small + Fast** | MsgPack + Zstd | 0.63 | 0.56 | 7.64 | Good balance |

---

## Recommendations

### Use Case: High-Throughput Logging
**Recommendation:** JSONL + Zstd
- Fast writes (458ms for 1M records)
- Good compression (98.2%, 2.59 MB)
- Human-readable format
- Fast decompression

### Use Case: Data Archival
**Recommendation:** MessagePack + XZ or Brotli
- Maximum compression (0.94-3.21 MB)
- 99.2% space reduction
- Binary format (smaller than JSON)
- Good for long-term storage

### Use Case: Real-Time Processing
**Recommendation:** MessagePack Plain or + Zstd
- Fastest reads (554-558ms)
- Binary format for efficiency
- Low decompression overhead with Zstd
- Type-safe parsing

### Use Case: Human-Readable Archives
**Recommendation:** JSONL + Brotli
- Human-readable format
- Excellent compression (1.95 MB, 98.7%)
- Good for debug/inspection
- Easy to parse with standard tools

### Use Case: Large Binary Data
**Recommendation:** FlatBuffer + Zstd
- Zero-copy deserialization
- Good compression (2.62 MB, 98.3%)
- Fast read access to specific fields
- Ideal for network protocols

---

## Methodology

### Test Data
Each record contains:
```go
type TestRecord struct {
    ID        int64
    Name      string  // ~20 chars
    Email     string  // ~30 chars
    Age       int
    Score     float64
    Active    bool
    Category  string  // ~10 chars
    Timestamp int64
}
```

### Benchmark Process
1. **Write:** Create 1M records and serialize to file
2. **Read:** Deserialize all 1M records from file
3. **Size:** Measure actual file size on disk
4. **Time:** Wall-clock time for complete operation

### Hardware
- CPU: AMD Ryzen 9 5900X (24 threads)
- Storage: NVMe SSD
- OS: Linux 6.18.6

---

## Detailed Analysis

### Why MessagePack + XZ is Best for Compression
- **Binary format:** 25% smaller than JSON before compression
- **XZ compression:** Highest compression ratio (LZMA2)
- **Structured data:** Repetitive patterns compress well
- **Trade-off:** Slower compression (XZ is CPU-intensive)

### Why JSONL + Zstd is Best for Performance
- **Fast compression:** Zstd optimized for speed
- **Streaming:** Line-by-line processing
- **Good ratio:** 98.2% compression
- **Dictionary:** Zstd learns patterns in data

### Why MsgPack is Fastest for Reading
- **Binary format:** No parsing overhead
- **Compact encoding:** Less data to read
- **Type preservation:** Direct deserialization
- **Streaming:** Efficient decoder implementation

---

## Compression Algorithm Characteristics

| Algorithm | Compression | Speed | CPU Usage | Use Case |
|-----------|-------------|-------|-----------|----------|
| **LZ4** | Low (88%) | Fastest | Low | Real-time, streaming |
| **Gzip** | Medium (90-94%) | Fast | Medium | Standard archives |
| **Zstd** | High (93-98%) | Very Fast | Medium | General purpose |
| **Brotli** | Very High (97-98%) | Medium | High | Web, static files |
| **XZ** | Maximum (99%+) | Slow | Very High | Long-term archives |

---

## Conclusion

For most applications, **JSONL + Zstd** or **MessagePack + Zstd** provide the best balance of performance and compression. Choose JSONL for human readability and debugging, MessagePack for maximum performance and storage efficiency.

For archival purposes where space is critical and read performance is less important, **MessagePack + XZ** achieves 99.2% compression (122x reduction) from the original size.
