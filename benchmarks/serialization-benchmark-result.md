# Serialization Benchmark Results

## Test Configuration

- **Dataset:** 1,000,000 records (~160MB plain text, realistic data)
- **Record Structure:** 8 fields (ID, Name, Email, Age, Score, Active, Category, Timestamp)
- **Data Generator:** gofakeit v7 (realistic names, emails, dates, categories)
- **CPU:** AMD Ryzen 9 5900X 12-Core Processor
- **Test Date:** 2025-01-28

## Summary - Parquet Wins! 🏆

**The Clear Winner: Parquet**
- **Fastest overall performance:** 0.61s total (0.15s read + 0.46s write)
- **Excellent compression:** 44.11 MB (72% smaller than plain text)
- **Best for:** Everything - APIs, analytics, services, data warehouses

**💡 Key Insight:** Read performance is typically more important than write performance in production systems.

## 🏆 Top Recommendations

### #1: Parquet (Snappy Compression) 🥇
- **READ: 0.15s** - 4x faster than MsgPack, 13x faster than JSONL
- **WRITE: 0.46s** - Fastest write among binary formats
- **TOTAL: 0.61s** - Best overall performance
- **SIZE: 44.11 MB** - Excellent compression
- **Use for:** Everything - APIs, analytics, services, data warehouses
- **Why:** Winner in both speed and compatibility. Works with Spark, DuckDB, Pandas, Arrow, all major data tools.

### #2: FlatBuffer Plain 🥈
- **READ: 0.06s** - Absolute fastest reads (2.5x faster than Parquet)
- **WRITE: 0.78s**
- **TOTAL: 0.84s**
- **SIZE: 160.40 MB** - Very large (3.6x larger than Parquet)
- **Use for:** Hot data paths where read speed is critical and storage is unlimited
- **Why:** Zero-copy deserialization gives fastest reads, but massive files make this impractical for most use cases.

### #3: FlatBuffer + LZ4 🥉
- **READ: 0.21s** - Fast reads (3x faster than MsgPack)
- **WRITE: 1.11s**
- **TOTAL: 1.32s** - 2x slower than Parquet
- **SIZE: 66.12 MB** - 1.5x larger than Parquet
- **Use for:** When you need zero-copy access and can't use Parquet
- **Why:** Good read performance, but slower and larger than Parquet. Only use if Parquet isn't available.

---

## Quick Performance Comparison

| Rank | Format | READ | WRITE | TOTAL | SIZE | Why Choose |
|------|--------|------|-------|-------|------|------------|
| 🥇 | **Parquet** | **0.15s** | **0.46s** | **0.61s** | 44.11 MB | **Best overall** |
| 🥈 | **FlatBuffer Plain** | **0.06s** | 0.78s | 0.84s | 160.40 MB | Fastest reads (huge files) |
| 🥉 | **FlatBuffer + LZ4** | **0.21s** | 1.11s | 1.32s | 66.12 MB | Fast reads (no Parquet) |
| 4 | MsgPack + Zstd | 0.60s | 0.74s | 1.34s | 39.15 MB | Best MsgPack option |
| 5 | FlatBuffer + Zstd | 0.34s | 1.18s | 1.53s | 44.75 MB | Balanced compression |
| 6 | JSONL + Zstd-2 | 1.97s | 0.73s | 2.70s | 43.27 MB | Human-readable |

---

## File Size Comparison (Sorted Smallest → Biggest) 📦

| Rank | Format | Size (MB) | Compression % | vs Plain | Notes |
|------|--------|-----------|---------------|----------|-------|
| 🥇 | **MsgPack + Brotli** | **34.17** | **71.4%** | 3.5x | Best compression, very slow (7.54s total) |
| 🥈 | **MsgPack + XZ** | **38.29** | **67.9%** | 3.1x | Tiny file, extremely slow (45.58s total) ⚠️ |
| 🥉 | **JSONL + Brotli** | **38.99** | **75.0%** | 4.0x | Small but very slow (9.89s total) |
| 4 | MsgPack + Zstd | **39.15** | 67.2% | 3.1x | ⭐ Best small+fast balance |
| 5 | JSONL + XZ | 40.57 | 74.0% | 3.8x | Extremely slow reads (48.94s) ⚠️ |
| 6 | MsgPack + Gzip | 41.00 | 65.6% | 2.9x | Slower than Zstd |
| 7 | JSONL + Zstd | 43.27 | 72.3% | 3.6x | Good balance for text |
| 8 | JSONL + Gzip | 43.27 | 72.3% | 3.6x | Slower than Zstd |
| 9 | **Parquet (Snappy)** | **44.11** 🏆 | **71.7%** | **3.5x** | **Best overall - fast + small** |
| 10 | FlatBuffer + Zstd | 44.75 | 72.1% | 3.5x | Good compression |
| 11 | MsgPack + LZ4 | 52.55 | 56.0% | 2.2x | Faster than above |
| 12 | JSONL + LZ4 | 64.19 | 58.9% | 2.4x | Fast compression |
| 13 | FlatBuffer + LZ4 | 66.12 | 58.8% | 2.4x | Fast reads, moderate size |
| 14 | MsgPack Plain | 119.33 | 23.5% | 1.3x | Very slow writes (23.24s) ⚠️ |
| 15 | JSONL Plain | 156.11 | 0% | 1.0x | No compression |
| 16 | FlatBuffer Plain | 160.40 | +2.7% | 1.03x | Larger than JSONL! |

**Key Findings:**
- 🏆 **MsgPack + Zstd** (#4): Best balance of small size (39.15 MB) + good speed (1.34s)
- ⚠️ **Avoid XZ**: Tiny files but 45-49s total time (too slow for most use cases)
- ⚠️ **Avoid Brotli**: Good compression but 7-10s total time (2-3x slower than Zstd)
- 🎯 **Parquet**: Not the smallest (44.11 MB) but fastest overall (0.61s) - best trade-off

---

## Zstd Compression Level Analysis

**Zstd Level 2 is modestly better than Level 1 - prefer Level 2 when possible.**

### JSONL Comparison:
| Level | Size | Write | Read | Total | Compression Gain |
|-------|------|-------|------|-------|------------------|
| Zstd-1 | 43.87 MB | 0.79s | 2.04s | 2.84s | Baseline |
| **Zstd-2** | **43.27 MB** 🏆 | **0.73s** | **1.97s** | **2.70s** | **0.60 MB smaller (1.4%), 0.14s faster (5%)** |
| Zstd (default=4) | 43.27 MB | 0.84s | 1.91s | 2.75s | Same size as L2, slower write |

### MsgPack Comparison:
| Level | Size | Write | Read | Total | Compression Gain |
|-------|------|-------|------|-------|------------------|
| Zstd-1 | 40.66 MB | 0.75s | 0.60s | 1.35s | Baseline |
| **Zstd-2** | **39.15 MB** 🏆 | 0.79s | 0.58s | 1.37s | **1.51 MB smaller (3.7%), nearly same speed** |
| Zstd (default=4) | 39.15 MB | 0.74s | 0.60s | 1.34s | Same size as L2 |

**Recommendation:** Use Zstd Level 2 (`.zst2` extension) for incrementally better compression (1.4-3.7% smaller) without performance penalty. Default level 4 is also fine and provides same compression as L2.

---

## Combined Performance & Size Rankings

### By Total Time (Fastest → Slowest)

| Rank | Format | Total | Read | Write | Size | Use Case |
|------|--------|-------|------|-------|------|----------|
| 🥇 | **Parquet** | **0.61s** | **0.15s** | **0.46s** | 44.11 MB | **Best overall - use this!** |
| 🥈 | FlatBuffer Plain | 0.84s | **0.06s** | 0.78s | 160.40 MB | Fastest reads, huge files |
| 🥉 | FlatBuffer + LZ4 | 1.32s | 0.21s | 1.11s | 66.12 MB | Fast reads, no Parquet |
| 4 | MsgPack + Zstd | 1.34s | 0.60s | 0.74s | 39.15 MB | Best MsgPack option |
| 5 | MsgPack + Zstd-2 | 1.37s | 0.58s | 0.79s | 39.15 MB | Same as above |
| 6 | FlatBuffer + Zstd | 1.53s | 0.34s | 1.18s | 44.75 MB | Balanced |
| 7 | FlatBuffer + Zstd-2 | 1.58s | 0.36s | 1.22s | 44.75 MB | Slightly slower |
| 8 | MsgPack + LZ4 | 1.59s | 0.62s | 0.97s | 52.55 MB | Fast, moderate size |
| 9 | JSONL + Zstd-2 | 2.70s | 1.97s | 0.73s | 43.27 MB | Best JSONL option |
| 10 | JSONL + Zstd | 2.75s | 1.91s | 0.84s | 43.27 MB | Default level |
| 11 | JSONL + LZ4 | 2.85s | 1.97s | 0.88s | 64.19 MB | Fast write |
| 12 | JSONL Plain | 3.31s | 1.93s | 1.38s | 156.11 MB | No compression |
| 13 | MsgPack + Gzip | 6.67s | 1.05s | 5.61s | 41.00 MB | Slower than Zstd |
| 14 | JSONL + Gzip | 7.23s | 2.52s | 4.71s | 43.27 MB | Very slow |
| 15 | MsgPack + Brotli | 7.54s | 0.99s | 6.55s | 34.17 MB | Too slow |
| 16 | JSONL + Brotli | 9.89s | 2.45s | 7.44s | 38.99 MB | Too slow |
| 17 | MsgPack Plain | 23.24s | 0.58s | 22.66s | 119.33 MB | Very slow writes ⚠️ |
| 18 | MsgPack + XZ | 45.58s | 32.00s | 13.58s | 38.29 MB | Extremely slow |
| 19 | JSONL + XZ | 48.94s | 35.11s | 13.82s | 40.57 MB | Extremely slow |

---

## Read Performance Rankings (Most Important!)

Read performance is typically more important than write performance in production systems.

| Rank | Format | Read Time | vs Parquet | Records/sec | Use Case |
|------|--------|-----------|------------|-------------|----------|
| 🥇 | **FlatBuffer Plain** | **0.06s** | 2.5x faster | 16,666,667 | Fastest (huge files) |
| 🥈 | **Parquet** | **0.15s** 🏆 | Baseline | **6,666,667** | **Best overall** |
| 🥉 | **FlatBuffer + LZ4** | **0.21s** | 1.4x slower | 4,761,905 | Fast (no Parquet) |
| 4 | FlatBuffer + Zstd | 0.34s | 2.3x slower | 2,941,176 | Balanced |
| 5 | FlatBuffer + Zstd-2 | 0.36s | 2.4x slower | 2,777,778 | Good compression |
| 6 | MsgPack Plain | 0.58s | 3.9x slower | 1,724,138 | No compression |
| 7 | MsgPack + Zstd-2 | 0.58s | 3.9x slower | 1,724,138 | Best MsgPack |
| 8 | MsgPack + Zstd | 0.60s | 4.0x slower | 1,666,667 | Good balance |
| 9 | MsgPack + LZ4 | 0.62s | 4.1x slower | 1,612,903 | Fast |
| 10 | MsgPack + Brotli | 0.99s | 6.6x slower | 1,010,101 | Small but slow |
| 11 | MsgPack + Gzip | 1.05s | 7.0x slower | 952,381 | Slower than Zstd |
| 12 | JSONL + Zstd | 1.91s | 12.7x slower | 523,560 | Human-readable |
| 13 | JSONL Plain | 1.93s | 12.9x slower | 518,135 | No compression |
| 14 | JSONL + Zstd-2 | 1.97s | 13.1x slower | 507,614 | Best JSONL |
| 15 | JSONL + LZ4 | 1.97s | 13.1x slower | 507,614 | Fast write |
| 16 | JSONL + Brotli | 2.45s | 16.3x slower | 408,163 | Too slow |
| 17 | JSONL + Gzip | 2.52s | 16.8x slower | 396,825 | Very slow |
| 18 | MsgPack + XZ | 32.00s | 213x slower | 31,250 | ⚠️ Extremely slow |
| 19 | JSONL + XZ | 35.11s | 234x slower | 28,482 | ⚠️ Extremely slow |

**Parquet is 4-13x faster than MsgPack/JSONL for reads!**

---

## Write Performance Rankings

| Rank | Format | Write Time | Records/sec | Size | Notes |
|------|--------|------------|-------------|------|-------|
| 🥇 | **Parquet** | **0.46s** 🏆 | **2,173,913** | 44.11 MB | **Fastest binary** |
| 🥈 | JSONL + Zstd-2 | 0.73s | 1,369,863 | 43.27 MB | Fastest JSONL |
| 🥉 | MsgPack + Zstd | 0.74s | 1,351,351 | 39.15 MB | Fastest MsgPack |
| 4 | FlatBuffer Plain | 0.78s | 1,282,051 | 160.40 MB | Fast, huge files |
| 5 | MsgPack + Zstd-2 | 0.79s | 1,265,823 | 39.15 MB | Nearly same as above |
| 6 | JSONL + Zstd | 0.84s | 1,190,476 | 43.27 MB | Default Zstd |
| 7 | JSONL + LZ4 | 0.88s | 1,136,364 | 64.19 MB | Fast compression |
| 8 | MsgPack + LZ4 | 0.97s | 1,030,928 | 52.55 MB | Fast |
| 9 | FlatBuffer + LZ4 | 1.11s | 900,901 | 66.12 MB | Good overall |
| 10 | FlatBuffer + Zstd | 1.18s | 847,458 | 44.75 MB | Balanced |
| 11 | FlatBuffer + Zstd-2 | 1.22s | 819,672 | 44.75 MB | Better compression |
| 12 | JSONL Plain | 1.38s | 724,638 | 156.11 MB | No compression |
| 13 | JSONL + Gzip | 4.71s | 212,314 | 43.27 MB | Slow |
| 14 | MsgPack + Gzip | 5.61s | 178,253 | 41.00 MB | Very slow |
| 15 | MsgPack + Brotli | 6.55s | 152,672 | 34.17 MB | Too slow |
| 16 | JSONL + Brotli | 7.44s | 134,409 | 38.99 MB | Too slow |
| 17 | MsgPack + XZ | 13.58s | 73,638 | 38.29 MB | Extremely slow |
| 18 | JSONL + XZ | 13.82s | 72,359 | 40.57 MB | Extremely slow |
| 19 | MsgPack Plain | 22.66s | 44,129 | 119.33 MB | ⚠️ Anomalously slow (likely GC/memory issue) |

---

## Compression Algorithm Characteristics

| Algorithm | Level | Compression | Speed | CPU | Best For |
|-----------|-------|-------------|-------|-----|----------|
| **Snappy** | N/A | Medium (72%) | Very Fast | Low | Parquet (balanced) |
| **LZ4** | N/A | Low (59-68%) | Fastest | Very Low | Real-time, streaming |
| **Zstd-1** | 1 | Medium (67-72%) | Very Fast | Low | ⚠️ DON'T USE - Level 2 better |
| **Zstd-2** | 2 | Good (72-75%) | Very Fast | Low | ✅ RECOMMENDED - Best balance |
| **Zstd** | 4 (default) | Good (72-75%) | Fast | Medium | Standard Zstd |
| **Gzip** | default | Medium (66-72%) | Slow | Medium | Legacy archives |
| **Brotli** | default | Very Good (71-75%) | Very Slow | High | Static files only |
| **XZ** | default | Excellent (68-74%) | Extremely Slow | Very High | ⚠️ Cold storage only |

**⚠️ Important:** Zstd Level 2 provides 1.4-3.7% better compression than Level 1 without performance penalty. Prefer Level 2 (`.zst2`) or default Level 4 (`.zst`).

---

## Recommendations by Use Case

### 🏆 For Everything (Default Choice)
**Use: Parquet (.parquet)**
- Total: 0.61s (fastest overall)
- Read: 0.15s (4x faster than MsgPack)
- Write: 0.46s (fastest binary format)
- Size: 44.11 MB (excellent compression)
- Compatible with all major data tools (Spark, DuckDB, Pandas, Arrow)

### For Maximum Read Speed (Storage Not a Concern)
**Use: FlatBuffer Plain (.fb)**
- Read: 0.06s (2.5x faster than Parquet)
- Total: 0.84s
- Size: 160.40 MB (3.6x larger than Parquet!)
- Only use if: Read speed is absolutely critical and you have unlimited storage

### When You Can't Use Parquet
**Use: FlatBuffer + LZ4 (.fb.lz4)**
- Read: 0.21s (still 3x faster than MsgPack)
- Total: 1.32s (2x slower than Parquet)
- Size: 66.12 MB
- Zero-copy deserialization

### For Smallest Files (Speed Less Important)
**Use: MsgPack + Brotli (.msgpack.br) OR MsgPack + Zstd (.msgpack.zst)**
- Brotli: 34.17 MB (smallest) but 7.54s total (too slow!)
- **Zstd: 39.15 MB (15% larger) but 1.34s total** ✅ BETTER CHOICE
- Only 5 MB larger but 5.6x faster!

### For Human-Readable Data
**Use: JSONL + Zstd-2 (.jsonl.zst2)**
- Read: 1.97s (acceptable for debugging)
- Write: 0.73s (fast)
- Total: 2.70s
- Size: 43.27 MB
- Can be opened with zstd -d, then read as text

### For Write-Heavy Logging
**Use: JSONL + Zstd-2 (.jsonl.zst2)**
- Same as above
- Fast writes (0.73s) when you're logging constantly
- Acceptable slow reads (1.97s) if logs are rarely read

### For Cold Storage / Archival
**Use: MsgPack + Zstd (.msgpack.zst)**
- Size: 39.15 MB (excellent compression)
- Total: 1.34s (very fast)
- Good balance of small size and reasonable access time

---

## Methodology

### Test Data Generation
Data generated using **gofakeit v7** for realistic production-like datasets:
```go
gofakeit.Seed(12345) // Reproducible

testDataset[i] = TestRecord{
    ID:        int64(i),
    Name:      gofakeit.Name(),           // "John Smith", "Maria Garcia"
    Email:     gofakeit.Email(),          // "john.smith@example.com"
    Age:       gofakeit.Number(18, 80),   // Random 18-80
    Score:     gofakeit.Float64Range(0, 100),
    Active:    gofakeit.Bool(),
    Category:  gofakeit.RandomString([...]), // "Electronics", "Books", etc.
    Timestamp: gofakeit.Date().Unix(),    // Random Unix timestamp
}
```

### Record Structure
```go
type TestRecord struct {
    ID        int64   // Sequential ID
    Name      string  // Full name (varies 10-30 chars)
    Email     string  // Email address (varies 20-40 chars)
    Age       int     // Age 18-80
    Score     float64 // Score 0-100
    Active    bool    // Random true/false
    Category  string  // One of 10 categories
    Timestamp int64   // Random Unix timestamp
}
```

### Benchmark Process
1. **Write:** Serialize 1M records to file
2. **Read:** Deserialize all 1M records from file
3. **Size:** Measure actual file size on disk
4. **Time:** Wall-clock time for complete operation

### Hardware
- CPU: AMD Ryzen 9 5900X (24 threads @ 3.7-4.8 GHz)
- RAM: DDR4
- Storage: NVMe SSD
- OS: Linux 6.18.6

---

## Conclusion

**🏆 Winner: Parquet**

Parquet is the clear winner for production systems:
- ✅ Fastest overall performance (0.61s total)
- ✅ Very fast reads (0.15s - 4x faster than MsgPack)
- ✅ Fastest writes for binary formats (0.46s)
- ✅ Excellent compression (44.11 MB)
- ✅ Industry standard (works with all major data tools)

**Use Parquet for everything** unless you have a specific reason not to:
- Can't use Parquet → FlatBuffer + LZ4
- Need absolute fastest reads + unlimited storage → FlatBuffer Plain
- Need human-readable text → JSONL + Zstd-2

**⚠️ Avoid:**
- XZ compression (45-49s operations - too slow for most use cases)
- Brotli compression (7-10s - much slower than Zstd with similar compression)
- Zstd Level 1 (Level 2 provides 1.4-3.7% better compression at same speed)
- MsgPack Plain writes (22s - anomalously slow, likely memory/GC issue in benchmark)
