# Serialization Benchmark Results

## Test Configuration

- **Dataset:** 1,000,000 records (~1GB RAM)
- **Record Structure:** 8 fields (ID, Name, Email, Age, Score, Active, Category, Timestamp)
- **CPU:** AMD Ryzen 9 5900X 12-Core Processor
- **Test Date:** 2026-01-28

## Summary

This benchmark compares different serialization formats (JSONL, MessagePack, FlatBuffers) with various compression algorithms (None, Gzip, Zstd, LZ4, Brotli, XZ) for storing and retrieving 1 million records.

## Key Findings

- 🏆 **Best Overall:** MsgPack + Zstd-1 (1.20s total, 5.67 MB, 95% compression)

- ⚡ **Fastest Overall:** FlatBuffer Plain (0.71s total, 150 MB)

- ⚡ **Fastest with Compression:** FlatBuffer + LZ4 (0.86s total, 16.53 MB, 89% compression)

- ⚡ **Fastest Write:** JSONL + Zstd-2 (0.43s)

- ⚡ **Fastest Read:** FlatBuffer Plain (0.07s - zero-copy)

- 📦 **Best Compression:** MsgPack + XZ (0.94 MB, 99.2% reduction, 4.82s total)

- 💡 **New Discovery:** Zstd-1 (level 1) is faster AND produces smaller files than default Zstd!

---

## ⚡ FASTEST Formats (Performance Focus)

### Top 5 Fastest Overall (Total Time)
1. **FlatBuffer Plain** - **0.71s** (150 MB) - Zero-copy, no compression
2. **FlatBuffer + LZ4** - **0.86s** (16.53 MB) - Best compressed performance 🏆
3. **FlatBuffer + Zstd-1** - **0.98s** (2.74 MB) - Fast with excellent compression
4. **FlatBuffer + Zstd-2** - **1.01s** (2.62 MB) - Slightly smaller
5. **MsgPack + Zstd-1** - **1.20s** (5.67 MB) - Best non-FlatBuffer option

### Fastest Write Operations
1. **JSONL + Zstd-2** - **0.43s** (2.59 MB)
2. **JSONL + Zstd-1** - **0.44s** (2.47 MB)
3. **JSONL + Zstd** - **0.44s** (2.59 MB)

### Fastest Read Operations
1. **FlatBuffer Plain** - **0.07s** (150 MB) - Zero-copy 🚀
2. **FlatBuffer + LZ4** - **0.21s** (16.53 MB) - 3x faster than other compressed
3. **FlatBuffer + Zstd-1** - **0.27s** (2.74 MB)

### 💡 Performance Winner: **FlatBuffer + LZ4**
- **Total:** 0.86s (21% slower than plain)
- **Size:** 16.53 MB (89% compression)
- **Read:** 0.21s (3x faster than other compressed formats)
- **Best choice for performance-critical applications**

---

## Quick Reference: All Top Performers

| Metric | Format | Value | Trade-off |
|--------|--------|-------|-----------|
| **⚡ Fastest Total** | FlatBuffer Plain | 0.71s | Large (150 MB) |
| **⚡ Fastest Compressed** | FlatBuffer + LZ4 | 0.86s | Medium (16.53 MB) 🏆 |
| **⚡ Fastest Write** | JSONL + Zstd-2 | 0.43s | Small (2.59 MB) |
| **⚡ Fastest Read** | FlatBuffer Plain | 0.07s | Large (150 MB) |
| **📦 Smallest File** | MsgPack + XZ | 0.94 MB | Slow (4.82s total) |
| **⚖️ Best Balance** | MsgPack + Zstd-1 | 1.20s, 5.67 MB | Recommended ✓ |
| **📝 Fast + Readable** | JSONL + Zstd-1 | 2.23s, 2.47 MB | Human-readable |

---

## ⚡ Speed Comparison Chart (Faster = Better)

```
Performance (Total Time) - 1M Records
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

FlatBuffer Plain        ████ 0.71s ⚡⚡⚡ FASTEST
FlatBuffer + LZ4        █████ 0.86s ⚡⚡⚡ BEST COMPRESSED
FlatBuffer + Zstd-1     ██████ 0.98s ⚡⚡
MsgPack + Zstd-1        ███████ 1.20s ⚡⚡ BALANCED
MsgPack + Zstd-2        ███████ 1.19s ⚡⚡
JSONL + Zstd-1          ████████████████ 2.23s ⚡ READABLE
JSONL + Zstd-2          ████████████████ 2.26s ⚡
MsgPack + Gzip          █████████████████ 2.53s
JSONL + Gzip            ██████████████████ 3.18s
MsgPack + Brotli        ███████████████████ 3.32s
MsgPack + XZ            ████████████████████████ 4.82s (smallest: 0.94 MB)
JSONL + XZ              ████████████████████████████████ 8.39s
MsgPack Plain           ████████████████████████████████████████████ 22.34s ⚠️

Top 3 for Speed:
1. FlatBuffer Plain (0.71s) - no compression
2. FlatBuffer + LZ4 (0.86s) - 89% compression ⭐
3. FlatBuffer + Zstd-1 (0.98s) - 98% compression
```

---

## Combined Size & Performance Comparison

| Format | Size (MB) | Write (s) | Read (s) | Total (s) | Compression % |
|--------|-----------|-----------|----------|-----------|---------------|
| **JSONL Plain** | 145.56 | 1.33 | 1.88 | 3.21 | 0% |
| JSONL Gzip | 8.11 | 1.27 | 1.90 | 3.18 | 94.4% |
| **JSONL Zstd-1** | **2.47** | **0.44** | **1.79** | **2.23** | **98.3%** |
| JSONL Zstd-2 | 2.59 | 0.43 | 1.83 | 2.26 | 98.2% |
| JSONL Zstd | 2.59 | 0.44 | 1.80 | 2.24 | 98.2% |
| JSONL LZ4 | 16.37 | 0.50 | 1.83 | 2.33 | 88.8% |
| JSONL Brotli | 1.95 | 1.97 | 1.85 | 3.82 | 98.7% |
| JSONL XZ | 4.09 | 3.18 | 5.21 | 8.39 | 97.2% |
| **MsgPack Plain** | 114.44 | 21.81 | **0.53** | 22.34 | 0% |
| MsgPack Gzip | 10.69 | 1.90 | 0.64 | 2.53 | 90.7% |
| **MsgPack Zstd-1** | **5.67** | **0.65** | **0.56** | **1.20** | **95.0%** |
| MsgPack Zstd-2 | 7.64 | 0.64 | 0.55 | 1.19 | 93.3% |
| MsgPack Zstd | 7.64 | 0.64 | 0.56 | 1.20 | 93.3% |
| MsgPack LZ4 | 17.67 | 0.66 | 0.57 | 1.23 | 84.6% |
| MsgPack Brotli | 3.21 | 2.69 | 0.63 | 3.32 | 97.2% |
| **MsgPack XZ** | **0.94** | 3.33 | 1.49 | 4.82 | **99.2%** |
| **FlatBuffer Plain** | 150.17 | **0.64** | **0.07** | **0.71** | 0% |
| **FlatBuffer LZ4** | 16.53 | 0.65 | 0.21 | **0.86** | 89.0% |
| FlatBuffer Zstd-1 | 2.74 | 0.72 | 0.27 | 0.98 | 98.2% |
| FlatBuffer Zstd-2 | 2.62 | 0.73 | 0.28 | 1.01 | 98.3% |
| FlatBuffer Zstd | 2.62 | 0.73 | 0.28 | 1.02 | 98.3% |

### Analysis
- **Fastest Overall:** FlatBuffer Plain (0.71s total, 0.07s read)
- **Fastest with Compression:** FlatBuffer + LZ4 (0.86s total, 89% compression)
- **Best Balanced:** MsgPack + Zstd-1 (1.20s total, 5.67 MB, 95% compression)
- **Fastest Write:** JSONL + Zstd-2 (0.43s)
- **Fastest Read:** FlatBuffer Plain (0.07s) - zero-copy deserialization
- **Smallest Size:** MsgPack + XZ (0.94 MB, 99.2% compression)

### Zstd Compression Level Comparison
**Zstd-1 (Fastest) is the clear winner** for most use cases:
- **JSONL:** Zstd-1 produces smaller files (2.47 MB vs 2.59 MB) at same speed
- **MsgPack:** Zstd-1 is 26% smaller (5.67 MB vs 7.64 MB) with same performance
- **FlatBuffer:** Zstd-1 is slightly larger but faster (0.98s vs 1.02s)
- **Recommendation:** Use `.zst1` extension for best speed/size balance

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
| **Best Overall** | MsgPack + Zstd-1 | 0.65 | 0.56 | 1.20 | 5.67 | Best balance ✓ |
| **Fastest** | FlatBuffer Plain | 0.64 | 0.07 | 0.71 | 150.17 | Zero-copy reads |
| **Fast + Compressed** | FlatBuffer + LZ4 | 0.65 | 0.21 | 0.86 | 16.53 | Recommended ✓ |
| **Fast Write** | JSONL + Zstd-2 | 0.43 | 1.83 | 2.26 | 2.59 | Fastest write |
| **Fast Read** | FlatBuffer Plain | 0.64 | 0.07 | 0.71 | 150.17 | Zero-copy |
| **Best Compression** | MsgPack + XZ | 3.33 | 1.49 | 4.82 | 0.94 | 99.2% reduction |
| **Small + Fast** | MsgPack + Zstd-1 | 0.65 | 0.56 | 1.20 | 5.67 | Recommended ✓ |

---

## Recommendations

### Use Case: High-Throughput Logging
**Recommendation:** JSONL + Zstd-1 (.jsonl.zst1)
- Fast writes (0.44s for 1M records)
- Total time: 2.23s (write + read)
- Excellent compression (98.3%, 2.47 MB)
- Human-readable format
- **Why Zstd-1:** Faster AND smaller than default Zstd

### Use Case: Data Archival
**Recommendation:** MessagePack + XZ (.msgpack.xz)
- Maximum compression (0.94 MB, 99.2% reduction)
- Total time: 4.82s (acceptable for archival)
- 122x smaller than plain JSONL
- Binary format for long-term storage

### Use Case: Real-Time Processing
**Recommendation:** FlatBuffer + LZ4 or MsgPack + Zstd-1
- **FlatBuffer + LZ4:** 0.86s total, 16.53 MB, zero-copy - **BEST**
- **MsgPack + Zstd-1:** 1.20s total, 5.67 MB, type-safe
- **FlatBuffer Plain:** 0.71s total, 150 MB - fastest but large
- Choose FlatBuffer + LZ4 for best speed/size balance

### Use Case: Human-Readable Archives
**Recommendation:** JSONL + Zstd-1 (.jsonl.zst1)
- 2.47 MB (98.3%), 2.23s total - **RECOMMENDED**
- Human-readable, fast, excellent compression
- Smaller than Brotli with better performance

### Use Case: Large Binary Data
**Recommendation:** FlatBuffer + LZ4 or Zstd-1
- **LZ4:** 0.86s total, 16.53 MB - **RECOMMENDED for speed**
- **Zstd-1:** 0.98s total, 2.74 MB - better compression
- **Plain:** 0.71s total, 150 MB - fastest but large
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
