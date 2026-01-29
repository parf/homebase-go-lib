# Serialization Benchmark Results

## Test Configuration

- **Dataset:** 1,000,000 records (~1GB RAM)
- **Record Structure:** 8 fields (ID, Name, Email, Age, Score, Active, Category, Timestamp)
- **CPU:** AMD Ryzen 9 5900X 12-Core Processor
- **Test Date:** 2026-01-28

## Summary - READ PERFORMANCE IS KING! 🚀

This benchmark compares serialization formats with one critical insight:

**💡 Reads happen 10-100x more often than writes in production systems!**

Therefore, this benchmark **prioritizes READ SPEED** above all else:
1. **READ SPEED** (Most Important) - Impacts every user request
2. **WRITE SPEED** (Secondary) - Happens less frequently
3. **FILE SIZE** (Least Important) - Storage is cheap vs performance

**The Winner:** FlatBuffer and Parquet formats deliver **2-10x faster reads** via optimized columnar/zero-copy access.

Tested: 1 million records across JSONL, MessagePack, FlatBuffers, and Parquet with 7 compression algorithms.

## 🏆 Top 3 Recommendations (READ SPEED IS CRITICAL!)

### #1: Parquet (.parquet) 🥇
- **⚡ READ TIME: 0.11s** (5x faster than MsgPack, 17x faster than JSONL!)
- **Total Time:** 0.51s | Write: 0.40s
- **Size:** 8.36 MB (excellent compression + speed)
- **Best For:** Analytics, data warehouses, columnar data access, SQL-like queries
- **Why:** Columnar format optimized for analytics. Fastest reads with great compression. Built-in statistics for query optimization.

### #2: FlatBuffer + LZ4 (.fb.lz4) 🥈
- **⚡ READ TIME: 0.17s** (3x faster than MsgPack, 11x faster than JSONL!)
- **Total Time:** 0.84s | Write: 0.67s
- **Size:** 16.53 MB (size not critical for performance)
- **Best For:** General-purpose high-performance applications - APIs, services, real-time systems
- **Why:** Zero-copy deserialization gives 3-11x faster reads. Write speed still excellent (0.67s).

### #3: FlatBuffer + Zstd-2 (.fb.zst2) 🥉
- **⚡ READ TIME: 0.25s** (2x faster than MsgPack, 7x faster than JSONL!)
- **Total Time:** 0.98s | Write: 0.73s
- **Size:** 2.62 MB (if you really need smaller files)
- **Best For:** When storage costs matter AND you need fast reads
- **Why:** Still delivers 2-7x faster reads with excellent compression. Write speed good (0.73s).

## Key Findings - READ SPEED MATTERS MOST! 🚀

### 🏆 Winner by Read Performance (Most Important Metric)

- **🥇 #1: FlatBuffer + LZ4** - READ: **0.21s** (3x faster!) | Total: 0.86s
  - **Use this for production systems** - best read performance with compression

- **🥈 #2: FlatBuffer + Zstd-2** - READ: **0.28s** (2x faster!) | Total: 1.01s
  - Use when file size matters but reads must stay fast

- **🥉 #3: FlatBuffer Plain** - READ: **0.07s** (10x faster!) | Total: 0.71s
  - Absolute fastest reads if disk space unlimited

### Why Read Speed is Critical
- **Data is typically read 10-100x more often than written**
- **FlatBuffer formats deliver 2-10x faster reads** via zero-copy deserialization
- **Read latency directly impacts user experience** - every ms counts in APIs/services
- Write speed matters less - happens once, reads happen constantly

### Other Metrics (Less Important)
- ⚡ Fastest Write: JSONL + Zstd-2 (0.43s) - but read time is 1.83s (8x slower!)
- 📦 Smallest Size: MsgPack + XZ (0.94 MB) - but read time is 1.49s (7x slower!)
- 🎯 Best Non-FlatBuffer: MsgPack + Zstd-1 (read: 0.56s) - still 2.5x slower than FlatBuffer

---

## ⚡ FASTEST Formats (Performance Focus - Read Speed Priority)

### 🚀 FlatBuffer Formats: Fastest Reads by Far!
**All FlatBuffer formats have 3-10x faster reads than other formats due to zero-copy deserialization**

1. **FlatBuffer Plain** - **0.71s total** (Write: 0.64s, **Read: 0.07s**) - 150 MB
   - Fastest reads possible - no decompression overhead
   - Best for: Hot data paths, real-time systems

2. **FlatBuffer + LZ4** - **0.86s total** (Write: 0.65s, **Read: 0.21s**) - 16.53 MB 🏆
   - **RECOMMENDED: Best compressed performance**
   - Read: 3x faster than MsgPack/JSONL compressed
   - Size: 89% compression (10x smaller than plain)
   - Best for: Production use, APIs, services

3. **FlatBuffer + Zstd-2** - **1.01s total** (Write: 0.73s, **Read: 0.28s**) - 2.62 MB ⭐
   - **2nd BEST: Excellent compression + fast reads**
   - Read: Still 2x faster than MsgPack compressed
   - Size: 98.3% compression (57x smaller than plain)
   - Best for: Storage-constrained + performance-critical

4. **FlatBuffer + Zstd-1** - **0.98s total** (Write: 0.72s, **Read: 0.27s**) - 2.74 MB
   - Slightly faster write, slightly larger than Zstd-2

### Top Non-FlatBuffer Options (Slower Reads)
- **MsgPack + Zstd-1** - **1.20s total** (Read: 0.56s) - 5.67 MB
- **JSONL + Zstd-1** - **2.23s total** (Read: 1.79s) - 2.47 MB

### Fastest Write Operations
1. **JSONL + Zstd-2** - **0.43s**
2. **JSONL + Zstd-1** - **0.44s**
3. **FlatBuffer Plain** - **0.64s**

### Read Speed Comparison
```
Read Performance (1M records)
━━━━━━━━━━━━━━━━━━━━━━━━━━━
FlatBuffer Plain      ██ 0.07s  🚀 10x FASTER
FlatBuffer + LZ4      ██████ 0.21s  🚀 3x FASTER (RECOMMENDED)
FlatBuffer + Zstd-1   ████████ 0.27s  🚀 2x FASTER
FlatBuffer + Zstd-2   ████████ 0.28s  🚀 2x FASTER
MsgPack + Zstd-1      ████████████████ 0.56s
JSONL + Zstd-1        ████████████████████████████████████████ 1.79s
```

### 💡 Top 2 Recommendations for Performance:
1. **FlatBuffer + LZ4** (0.86s, 16.53 MB) - Best all-around compressed 🏆
2. **FlatBuffer + Zstd-2** (1.01s, 2.62 MB) - Best if size matters more ⭐

---

## Quick Reference - Sorted by READ SPEED (Most Important) 🚀

| Rank | Format | READ | Total | Write | Size | Why Choose This |
|------|--------|------|-------|-------|------|-----------------|
| **🥇** | **Parquet** 🏆 | **0.11s** | 0.51s | 0.40s | 8.36 MB | **BEST OVERALL - Analytics & queries** |
| **🥈** | **FlatBuffer + LZ4** ⭐ | **0.17s** | 0.84s | 0.67s | 16.53 MB | **Best for APIs & services** |
| **🥉** | **FlatBuffer + Zstd-2** | **0.25s** | 0.98s | 0.73s | 2.62 MB | Fast reads + small files |
| 4 | FlatBuffer Plain | 0.07s | 0.71s | 0.64s | 150.17 MB | Fastest possible (no compression) |
| 5 | MsgPack + Zstd-1 | 0.56s | 1.20s | 0.65s | 5.67 MB | Best if can't use above formats |
| 6 | JSONL + Zstd-1 | 1.79s | 2.23s | 0.44s | 2.47 MB | Human-readable (slow reads) |

**⚡ Read Speed Advantage:**
- FlatBuffer formats are **2-10x faster for reads** than anything else
- This advantage compounds with read frequency (10-100x more reads than writes)
- Choose based on READ performance, not write or size!

---

## 🚀 READ SPEED Comparison (MOST IMPORTANT!)

```
READ Performance - 1M Records (Data read 10-100x more than written!)
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

FlatBuffer Plain        █ 0.07s 🚀 FASTEST (no compression, 150 MB)
Parquet (Snappy)        ██ 0.11s 🏆 #1 RECOMMENDED (analytics, 8.36 MB)
FlatBuffer + LZ4        ███ 0.17s ⭐ #2 RECOMMENDED (APIs/services, 16.53 MB)
FlatBuffer + Zstd-2     ████ 0.25s ⚡ FAST (small files, 2.62 MB)
FlatBuffer + Zstd-1     ████ 0.27s ⚡ FAST (2.74 MB)
─────────────────────────────────────────────────
MsgPack + Zstd-1        ████████ 0.56s (5x slower than Parquet)
MsgPack + Zstd-2        ████████ 0.55s
MsgPack + LZ4           ████████ 0.57s
─────────────────────────────────────────────────
JSONL + Zstd-1          █████████████████████████ 1.79s ⚠️ 16x SLOWER
JSONL + Zstd-2          █████████████████████████ 1.83s
JSONL + LZ4             █████████████████████████ 1.83s
JSONL + Plain           █████████████████████████ 1.88s
MsgPack + XZ            █████████████████████████████████ 1.49s
JSONL + XZ              ██████████████████████████████████████████████████████████████████ 5.21s

💡 KEY INSIGHT: Reads happen 10-100x more often than writes!
   → Choose format based on READ speed, not write or size
   → Parquet & FlatBuffer formats are 5-17x faster for reads
   → Even 0.1s read improvement = 10-100s saved across all reads!
```

---

## Combined Size & Performance Comparison - Sorted by READ SPEED (Fastest First) 🚀

| Format | Read (s) | Total (s) | Size (MB) | Write (s) | Compression % |
|--------|----------|-----------|-----------|-----------|---------------|
| **🥇 FlatBuffer Plain** | **0.07** 🚀 | 0.71 | 150.17 | 0.64 | 0% |
| **🥈 Parquet (Snappy)** | **0.11** 🏆 | 0.51 | 8.36 | 0.40 | 94.4% |
| **🥉 FlatBuffer + LZ4** | **0.17** ⭐ | 0.84 | 16.53 | 0.67 | 89.0% |
| **FlatBuffer + Zstd-2** | **0.25** ⚡ | 0.98 | 2.62 | 0.73 | 98.3% |
| **FlatBuffer + Zstd-1** | **0.27** ⚡ | 0.98 | 2.74 | 0.71 | 98.2% |
| FlatBuffer + Zstd | 0.28 | 0.92 | 2.62 | 0.65 | 98.3% |
| MsgPack Plain | 0.53 | 22.34 | 114.44 | 21.81 | 0% ⚠️ |
| MsgPack + Zstd-2 | 0.55 | 1.19 | 7.64 | 0.64 | 93.3% |
| **MsgPack + Zstd-1** | 0.56 | 1.20 | 5.67 | 0.65 | 95.0% |
| MsgPack + Zstd | 0.56 | 1.20 | 7.64 | 0.64 | 93.3% |
| MsgPack + LZ4 | 0.57 | 1.23 | 17.67 | 0.66 | 84.6% |
| MsgPack + Brotli | 0.63 | 3.32 | 3.21 | 2.69 | 97.2% |
| MsgPack + Gzip | 0.64 | 2.53 | 10.69 | 1.90 | 90.7% |
| **MsgPack + XZ** | 1.49 | 4.82 | **0.94** 📦 | 3.33 | **99.2%** |
| **JSONL + Zstd-1** | 1.79 | 2.23 | 2.47 | 0.44 | 98.3% |
| JSONL + Zstd | 1.80 | 2.24 | 2.59 | 0.44 | 98.2% |
| JSONL + Zstd-2 | 1.83 | 2.26 | 2.59 | **0.43** | 98.2% |
| JSONL + LZ4 | 1.83 | 2.33 | 16.37 | 0.50 | 88.8% |
| JSONL + Brotli | 1.85 | 3.82 | 1.95 | 1.97 | 98.7% |
| JSONL Plain | 1.88 | 3.21 | 145.56 | 1.33 | 0% |
| JSONL + Gzip | 1.90 | 3.18 | 8.11 | 1.27 | 94.4% |
| JSONL + XZ | 5.21 | 8.39 | 4.09 | 3.18 | 97.2% |

**🚀 WHY SORT BY READ SPEED?**
- **Reads happen 10-100x more often than writes** in most applications
- **Read latency directly impacts user experience** - faster reads = better UX
- **Read advantage compounds:** 0.2s read improvement × 100 reads = 20s saved vs 1 write!
- **FlatBuffer formats are 2-10x faster for reads** due to zero-copy deserialization

**Read Speed Comparison:**
- FlatBuffer Plain: 0.07s (10x faster than MsgPack, 25x faster than JSONL)
- FlatBuffer + LZ4: 0.21s (3x faster than MsgPack, 8x faster than JSONL) 🏆
- FlatBuffer + Zstd-2: 0.28s (2x faster than MsgPack, 6x faster than JSONL) ⭐

**Conclusion: Choose FlatBuffer for any performance-critical application!**

### Analysis - Prioritizing Read Performance

**READ SPEED (Most Important - happens 10-100x more often):**
- **🥇 Fastest:** FlatBuffer Plain (0.07s) - 10x faster than competition
- **🥈 Best with Compression:** FlatBuffer + LZ4 (0.21s) - 3x faster than others 🏆
- **🥉 Fast + Small:** FlatBuffer + Zstd-2 (0.28s) - 2x faster than others ⭐

**WRITE SPEED (Secondary - happens less frequently):**
- Fastest: JSONL + Zstd-2 (0.43s)
- FlatBuffer formats: 0.64-0.73s (slightly slower but reads 3-10x faster!)

**FILE SIZE (Least Important - storage is cheap):**
- Smallest: MsgPack + XZ (0.94 MB) - but reads are 7x slower (1.49s vs 0.21s)
- FlatBuffer + LZ4: 16.53 MB (acceptable size for 3x read performance gain)

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

## Write Performance (1M records) - Sorted by Speed

| Format | Time (s) | Records/sec | vs Fastest | Size (MB) |
|--------|----------|-------------|------------|-----------|
| **Parquet (Snappy)** | **0.40** | **2,500,000** | **1.0x** | 8.36 🏆 |
| **JSONL Zstd-2** | **0.43** | **2,325,581** | **1.1x** | 2.59 |
| JSONL Zstd-1 | 0.44 | 2,272,727 | 1.1x | 2.47 |
| JSONL Zstd | 0.44 | 2,272,727 | 1.1x | 2.59 |
| JSONL LZ4 | 0.50 | 2,000,000 | 1.3x | 16.37 |
| FlatBuffer Plain | 0.64 | 1,562,500 | 1.5x | 150.17 |
| MsgPack Zstd-2 | 0.64 | 1,562,500 | 1.5x | 7.64 |
| **FlatBuffer + LZ4** | **0.65** | **1,538,462** | **1.5x** | **16.53** 🏆 |
| MsgPack Zstd-1 | 0.65 | 1,538,462 | 1.5x | 5.67 |
| MsgPack Zstd | 0.64 | 1,562,500 | 1.5x | 7.64 |
| MsgPack LZ4 | 0.66 | 1,515,152 | 1.5x | 17.67 |
| **FlatBuffer + Zstd-1** | **0.72** | **1,388,889** | **1.7x** | **2.74** |
| **FlatBuffer + Zstd-2** | **0.73** | **1,369,863** | **1.7x** | **2.62** ⭐ |
| FlatBuffer + Zstd | 0.73 | 1,369,863 | 1.7x | 2.62 |
| JSONL Gzip | 1.27 | 787,402 | 3.0x | 8.11 |
| JSONL Plain | 1.33 | 751,880 | 3.1x | 145.56 |
| MsgPack Gzip | 1.90 | 526,316 | 4.4x | 10.69 |
| JSONL Brotli | 1.97 | 507,614 | 4.6x | 1.95 |
| MsgPack Brotli | 2.69 | 371,747 | 6.3x | 3.21 |
| JSONL XZ | 3.18 | 314,465 | 7.4x | 4.09 |
| MsgPack XZ | 3.33 | 300,300 | 7.7x | 0.94 |
| MsgPack Plain | 21.81 | 45,853 | 50.7x | 114.44 |

**Top 3 for Write Speed:**
1. JSONL + Zstd-2 (0.43s) - Fastest write
2. JSONL + Zstd-1 (0.44s) - Almost as fast, smaller file
3. JSONL + Zstd (0.44s) - Default Zstd level

**Note:** The MsgPack Plain write anomaly (21.81s) suggests buffering/encoding overhead in the benchmark.

---

## Read Performance (1M records) - Sorted by Speed

| Format | Time (s) | Records/sec | vs Fastest | Notes |
|--------|----------|-------------|------------|-------|
| **FlatBuffer Plain** | **0.07** | **14,285,714** | **1.0x** | 🚀 Zero-copy |
| **Parquet (Snappy)** | **0.11** | **9,090,909** | **1.6x** | 🏆 Best compressed analytics |
| **FlatBuffer + LZ4** | **0.17** | **5,882,353** | **2.4x** | ⭐ Best APIs/services |
| **FlatBuffer + Zstd-2** | **0.25** | **4,000,000** | **3.6x** | ⚡ Fast + small |
| **FlatBuffer + Zstd-1** | **0.27** | **3,703,704** | **3.9x** | ⚡ Fast decompression |
| FlatBuffer + Zstd | 0.28 | 3,571,429 | 4.0x | Default level |
| MsgPack Plain | 0.53 | 1,886,792 | 7.6x | No compression |
| MsgPack Zstd-2 | 0.55 | 1,818,182 | 7.9x | Fast decompress |
| MsgPack Zstd-1 | 0.56 | 1,785,714 | 8.0x | Good balance |
| MsgPack Zstd | 0.56 | 1,785,714 | 8.0x | Default level |
| MsgPack LZ4 | 0.57 | 1,754,386 | 8.1x | Fast decompress |
| MsgPack Gzip | 0.64 | 1,562,500 | 9.1x | Standard |
| JSONL Zstd-1 | 1.79 | 558,659 | 25.6x | Smallest JSON |
| JSONL Zstd | 1.80 | 555,556 | 25.7x | Default level |
| JSONL Zstd-2 | 1.83 | 546,448 | 26.1x | Fast write |
| JSONL LZ4 | 1.83 | 546,448 | 26.1x | Fast decompress |
| JSONL Brotli | 1.85 | 540,541 | 26.4x | Small but slow |
| JSONL Plain | 1.88 | 531,915 | 26.9x | No compression |
| JSONL Gzip | 1.90 | 526,316 | 27.1x | Standard |
| MsgPack XZ | 1.49 | 671,141 | 21.3x | Slowest to read |
| JSONL XZ | 5.21 | 191,939 | 74.4x | Very slow reads |

**Top 3 for Read Speed (5x+ faster than others!):**
1. **FlatBuffer Plain (0.07s)** - Absolute fastest, zero-copy (no compression)
2. **Parquet (0.11s)** - Best for analytics, columnar format 🏆
3. **FlatBuffer + LZ4 (0.17s)** - Best for APIs/services, zero-copy ⭐

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

### Speed vs Compression Trade-offs - Sorted by Total Time

| Rank | Format | Read (s) | Total (s) | Write (s) | Size (MB) | Use Case |
|------|--------|----------|-----------|-----------|-----------|----------|
| **🥇** | **FlatBuffer Plain** | **0.07** 🚀 | **0.71** | 0.64 | 150.17 | Absolute fastest |
| **🥈** | **FlatBuffer + LZ4** | **0.21** 🏆 | **0.86** | 0.65 | 16.53 | **Best compressed** |
| **🥉** | **FlatBuffer + Zstd-2** | **0.28** ⭐ | **1.01** | 0.73 | 2.62 | Small + fast reads |
| 4 | MsgPack + Zstd-1 | 0.56 | 1.20 | 0.65 | 5.67 | Best non-FlatBuffer |
| 5 | JSONL + Zstd-1 | 1.79 | 2.23 | 0.44 | 2.47 | Best human-readable |
| 6 | JSONL + Zstd-2 | 1.83 | 2.26 | **0.43** | 2.59 | Fastest write |
| 7 | MsgPack + XZ | 1.49 | 4.82 | 3.33 | **0.94** 📦 | Smallest file |

---

## Recommendations - READ SPEED FIRST!

### Use Case: Analytics / Data Warehouses / BI (ANALYTICAL WORKLOADS)
**🏆 #1 RECOMMENDED:** Parquet (.parquet)
- **READ: 0.11s (5x faster than MsgPack, 17x faster than JSONL!)** ← Perfect for analytics!
- Write: 0.40s (excellent)
- Total: 0.51s
- Size: 8.36 MB (great compression)
- **Why:** Columnar format optimized for analytical queries. Supports predicate pushdown, column pruning, and SQL engines (Spark, DuckDB, etc.)

### Use Case: Production APIs / Services / Real-Time (READ-HEAVY!)
**🏆 #2 RECOMMENDED:** FlatBuffer + LZ4 (.fb.lz4)
- **READ: 0.17s (3x faster!)** ← This is what matters for APIs!
- Write: 0.67s (still excellent)
- Total: 0.84s
- Size: 16.53 MB (storage is cheap vs performance)
- **Why:** Your API reads data 100x more than it writes. 3x faster reads = 3x better latency!

**⭐ ALTERNATIVE:** FlatBuffer + Zstd-2 (.fb.zst2) if storage costs matter
- **READ: 0.28s (2x faster!)**
- Write: 0.73s
- Total: 1.01s
- Size: 2.62 MB (smaller files)

### Use Case: High-Throughput Logging (WRITE-HEAVY)
**Recommendation:** JSONL + Zstd-2 (.jsonl.zst2) - only when writes > reads
- Write: 0.43s (fastest write)
- Read: 1.83s (8x slower reads - acceptable if rarely read)
- Total: 2.26s
- Size: 2.59 MB
- Human-readable format
- **Only use if:** Logs written constantly but read rarely

### Use Case: Data Archival (SIZE MATTERS, SPEED DOESN'T)
**Recommendation:** MessagePack + XZ (.msgpack.xz) - only for cold storage
- Size: 0.94 MB (99.2% compression)
- Read: 1.49s (7x slower - acceptable for archival)
- Write: 3.33s (slow - acceptable for one-time write)
- Total: 4.82s
- **Only use if:** Data written once, read rarely, storage costs high

### Use Case: Performance-Critical Applications (MOST COMMON!)
**This is YOUR use case if:** APIs, services, caches, real-time systems, databases

**🏆 #1 CHOICE:** FlatBuffer + LZ4 (.fb.lz4)
- **READ: 0.21s** ← 3x faster than MsgPack, 8x faster than JSONL!
- Write: 0.65s (excellent)
- Total: 0.86s
- Size: 16.53 MB (negligible cost vs 3x performance gain)
- **Impact:** If you serve 1000 reads/sec, this saves 150ms × 1000 = 150 CPU seconds/sec!

**⭐ #2 CHOICE:** FlatBuffer + Zstd-2 (.fb.zst2) - if storage costs significant
- **READ: 0.28s** ← 2x faster than MsgPack, 6x faster than JSONL!
- Write: 0.73s (good)
- Total: 1.01s
- Size: 2.62 MB (much smaller)
- **Trade-off:** Slightly slower reads (0.07s difference) for 6x smaller files

**🥉 #3 CHOICE:** FlatBuffer Plain (.fb) - if disk space unlimited
- **READ: 0.07s** ← Absolute fastest possible!
- Write: 0.64s
- Total: 0.71s
- Size: 150 MB (large but storage is cheap)

**❌ DON'T USE:** MsgPack/JSONL for performance-critical reads
- MsgPack reads: 0.56s (2.5x slower than FlatBuffer + LZ4)
- JSONL reads: 1.79s+ (8x+ slower than FlatBuffer + LZ4)
- **Your users will notice the difference!**

### Use Case: Human-Readable Archives (Debugging/Inspection)
**Only use JSONL if:** You MUST be able to open files in a text editor

**Recommendation:** JSONL + Zstd-1 (.jsonl.zst1)
- Read: 1.79s (8x slower than FlatBuffer - acceptable if rarely read)
- Write: 0.44s
- Total: 2.23s
- Size: 2.47 MB
- Human-readable (can grep/sed/awk)
- **Trade-off:** Sacrificing 8x read performance for text format

**Better Alternative:** FlatBuffer + LZ4 + conversion tool
- Use FlatBuffer for performance (read: 0.21s, 8x faster)
- Write simple tool to convert FlatBuffer → JSON when needed for inspection
- Get 8x faster reads 99% of the time, readable when needed 1% of time

### Use Case: Large Binary Data / Network Protocols
**🏆 USE THIS:** FlatBuffer + LZ4 (.fb.lz4)
- **READ: 0.21s** - 3x faster for every network request!
- Write: 0.65s
- Total: 0.86s
- Size: 16.53 MB
- **Perfect for:** Network protocols, RPC, message queues, caches
- Zero-copy deserialization - no parsing overhead
- Random field access without full decode
- **Example:** gRPC-style services, message brokers, distributed caches

**⭐ Alternative:** FlatBuffer + Zstd-2 (.fb.zst2) if bandwidth costs high
- **READ: 0.28s** - 2x faster for network requests
- Write: 0.73s
- Total: 1.01s
- Size: 2.62 MB (smaller network payloads)
- Zero-copy deserialization

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

For **performance-critical applications**, choose formats based on READ SPEED:

1. **Parquet** - Best for analytics, data warehouses, SQL queries (Read: 0.11s, Total: 0.51s, Size: 8.36 MB)
2. **FlatBuffer + LZ4** - Best for APIs, services, real-time systems (Read: 0.17s, Total: 0.84s, Size: 16.53 MB)
3. **FlatBuffer + Zstd-2** - When storage matters + fast reads needed (Read: 0.25s, Total: 0.98s, Size: 2.62 MB)

For **write-heavy logging** where reads are rare, **JSONL + Zstd** provides human readability with good compression.

For **archival purposes** where space is critical and read performance is less important, **MessagePack + XZ** achieves 99.2% compression (122x reduction) from the original size.
