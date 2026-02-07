# cache - Universal Map-to-Parquet File Caching

File-based caching for Go maps using Parquet format with precise scalar type preservation.

## Usage

```go
import "github.com/parf/homebase-go-lib/cache"
```

### Basic pattern

```go
data, ok := cache.Map("users.parquet")
if !ok {
    // cache miss — compute data
    data = expensiveComputation()
    cache.WriteMap("users.parquet", data)
}
result := data.(map[string]any)
```

### String-keyed maps

```go
// map[string]any — each key becomes a parquet column
m := map[string]any{
    "name":  "Alice",
    "age":   uint32(30),
    "score": float64(95.5),
    "active": true,
}
cache.WriteMap("profile.parquet", m)

data, ok := cache.Map("profile.parquet")
// data.(map[string]any)["age"] is uint32(30) — type preserved
```

### Numeric-keyed maps

```go
// map[uint32]string — stored as key-value columns
lookup := map[uint32]string{
    1:   "one",
    2:   "two",
    100: "hundred",
}
cache.WriteMap("lookup.parquet", lookup)

data, ok := cache.Map("lookup.parquet")
result := data.(map[uint32]any)
// result[uint32(1)] == "one"
```

### Any map key/value combination

```go
cache.WriteMap("ids.parquet", map[int]float64{1: 1.1, 2: 2.2})
cache.WriteMap("flags.parquet", map[uint32]uint64{10: 1000, 20: 2000})
```

## Type Preservation

All scalar types survive the parquet round-trip with their exact Go type:

| Go type | Parquet type | Read-back type |
|---------|-------------|---------------|
| string | STRING | string |
| bool | BOOLEAN | bool |
| int8 | INT8 | int8 |
| int16 | INT16 | int16 |
| int32 | INT32 | int32 |
| int64 | INT64 | int64 |
| int | INT64 | int64 |
| uint8 | UINT8 | uint8 |
| uint16 | UINT16 | uint16 |
| uint32 | UINT32 | uint32 |
| uint64 | UINT64 | uint64 |
| uint | UINT64 | uint64 |
| float32 | FLOAT | float32 |
| float64 | DOUBLE | float64 |

`int` and `uint` are promoted to `int64`/`uint64` (platform-dependent sizes).

## Storage Format

- **String-keyed maps** (`map[string]any`): single-row parquet, each key = column name
- **Numeric-keyed maps** (`map[uint32]string`, etc.): multi-row parquet with `_key` and `_val` columns
- Internal Snappy compression
- Cache metadata stored in parquet key-value metadata

## API

```go
func Map(filename string) (any, bool)
```
Read cached map from parquet file. Returns `(data, true)` on hit, `(nil, false)` on miss.
Corrupted files are treated as cache misses.

```go
func WriteMap(filename string, data any)
```
Write any map to parquet cache file. Errors are silently ignored (cache is best-effort).
