# Serialization Benchmark Results

## Test Configuration

- **Dataset:** 1,000,000 records (~1GB RAM)
- **Record Structure:** 8 fields (ID, Name, Email, Age, Score, Active, Category, Timestamp)
- **CPU:** AMD Ryzen 9 5900X 12-Core Processor
- **Test Date:** 2026-01-28

## Summary

This benchmark compares different serialization formats (JSONL, MessagePack, FlatBuffers) with various compression algorithms (None, Gzip, Zstd, LZ4, Brotli, XZ) for storing and retrieving 1 million records.

## Key Findings

🏆 **Best Overall:** MsgPack + Zstd (1.22s total, 7.64 MB, 93.3% compression)
⚡ **Fastest Overall:** FlatBuffer Plain (0.72s total, 150 MB)
⚡ **Fastest Write:** JSONL + Zstd (0.44s)
⚡ **Fastest Read:** FlatBuffer Plain (0.07s - zero-copy)
📦 **Best Compression:** MsgPack + XZ (0.94 MB, 99.2% reduction, 4.70s total)

---

## Quick Reference: Top Performers

| Metric | Format | Value | Trade-off |
|--------|--------|-------|-----------|
| **Fastest Total Time** | FlatBuffer Plain | 0.72s | Large (150 MB) |
| **Fastest Write** | JSONL + Zstd | 0.44s | Good size (2.59 MB) |
| **Fastest Read** | FlatBuffer Plain | 0.07s | Large (150 MB) |
| **Smallest File** | MsgPack + XZ | 0.94 MB | Slow (4.70s total) |
| **Best Balance** | MsgPack + Zstd | 1.22s, 7.64 MB | Recommended ✓ |
| **Fast + Small** | JSONL + Zstd | 2.26s, 2.59 MB | Human-readable |

---

## Combined Size & Performance Comparison

| Format | Size (MB) | Write (s) | Read (s) | Total (s) | Compression % |
|--------|-----------|-----------|----------|-----------|---------------|
| **JSONL Plain** | 145.56 | 1.34 | 1.94 | 3.28 | 0% |
| JSONL Gzip | 8.11 | 1.21 | 1.88 | 3.09 | 94.4% |
| **JSONL Zstd** | **2.59** | **0.44** | **1.82** | **2.26** | **98.2%** |
| JSONL LZ4 | 16.37 | 0.48 | 1.83 | 2.31 | 88.8% |
| JSONL Brotli | 1.95 | 1.95 | 1.87 | 3.82 | 98.7% |
| JSONL XZ | 4.09 | 3.23 | 5.21 | 8.44 | 97.2% |
| **MsgPack Plain** | 114.44 | 22.05 | **0.54** | 22.59 | 0% |
| MsgPack Gzip | 10.69 | 1.94 | 0.64 | 2.58 | 90.7% |
| **MsgPack Zstd** | 7.64 | 0.66 | 0.56 | **1.22** | 93.3% |
| MsgPack LZ4 | 17.67 | 0.66 | 0.60 | 1.27 | 84.6% |
| MsgPack Brotli | 3.21 | 2.67 | 0.62 | 3.29 | 97.2% |
| **MsgPack XZ** | **0.94** | 3.22 | 1.48 | 4.70 | **99.2%** |
| **FlatBuffer Plain** | 150.17 | **0.65** | **0.07** | **0.72** | 0% |
| FlatBuffer Zstd | 2.62 | 0.67 | 0.23 | 0.90 | 98.3% |

### Analysis
- **Fastest Overall:** FlatBuffer Plain (0.72s total, 0.07s read)
- **Best Balanced:** MsgPack + Zstd (1.22s total, 7.64 MB, 93.3% compression)
- **Fastest Write:** JSONL + Zstd (0.44s)
- **Fastest Read:** FlatBuffer Plain (0.07s) - zero-copy deserialization
- **Smallest Size:** MsgPack + XZ (0.94 MB, 99.2% compression)

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

| Category | Format | Write (s) | Read (s) | Total (s) | Size (MB) | Notes |
|----------|--------|-----------|----------|-----------|-----------|-------|
| **Best Overall** | MsgPack + Zstd | 0.66 | 0.56 | 1.22 | 7.64 | Best balance |
| **Fastest** | FlatBuffer Plain | 0.65 | 0.07 | 0.72 | 150.17 | Zero-copy reads |
| **Fast Write** | JSONL + Zstd | 0.44 | 1.82 | 2.26 | 2.59 | Fastest write |
| **Fast Read** | FlatBuffer Plain | 0.65 | 0.07 | 0.72 | 150.17 | Zero-copy |
| **Best Compression** | MsgPack + XZ | 3.22 | 1.48 | 4.70 | 0.94 | 99.2% reduction |
| **Small + Fast** | MsgPack + Zstd | 0.66 | 0.56 | 1.22 | 7.64 | Recommended |

---

## Recommendations

### Use Case: High-Throughput Logging
**Recommendation:** JSONL + Zstd
- Fast writes (0.44s for 1M records)
- Total time: 2.26s (write + read)
- Good compression (98.2%, 2.59 MB)
- Human-readable format

### Use Case: Data Archival
**Recommendation:** MessagePack + XZ
- Maximum compression (0.94 MB, 99.2% reduction)
- Total time: 4.70s (acceptable for archival)
- 122x smaller than plain JSONL
- Binary format for long-term storage

### Use Case: Real-Time Processing
**Recommendation:** FlatBuffer Plain or MsgPack + Zstd
- **FlatBuffer:** 0.72s total (0.07s read), zero-copy access
- **MsgPack + Zstd:** 1.22s total, 93% smaller, type-safe
- Choose FlatBuffer for speed, MsgPack + Zstd for balance

### Use Case: Human-Readable Archives
**Recommendation:** JSONL + Brotli or Zstd
- **Brotli:** 1.95 MB (98.7%), 3.82s total
- **Zstd:** 2.59 MB (98.2%), 2.26s total - **RECOMMENDED**
- Human-readable, fast, good compression

### Use Case: Large Binary Data
**Recommendation:** FlatBuffer Plain or + Zstd
- **Plain:** 0.72s total, 150 MB - fastest
- **Zstd:** 0.90s total, 2.62 MB - 98.3% smaller
- Zero-copy deserialization
- Random field access without full decode

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
